package main

import (
	"fmt"
	"net"
	"time"
	"bufio"
	"os"
)

const PORTA = "172.16.201.5"

// ============= Usuários ==============

/*

1. Conecta ao Broker via TCP
2. Cria Usuário: Fornece usuarioID ao Broker
3. Executa Terminal
  L> Aguarda usuário definir qual tópico irá assinar - Deve informar no modelo:
	"tipo/tipoId/comando"
4. Envia comando, ou recebe valor/estado de sensor/atuador

*/

/*
Padrão de Tópico: tipo/tipoId/comando. ("Valor" não será utilizado dentro do tópico, mas como dado fornecido por sensores. O mesmo vale para "Estado" para o caso de atuadores)
Exemplo: sensor/sensor_1/-
Exemplo: atuador/atuador_1/off
*/
type Topico struct {
	Acao string `json:"acao"`
	Tipo string `json:"tipo"`
	TipoId string `json:"tipoId"`
	Comando string `json:"comando"`
	Valor float64 `json:"valor"`
	Estado bool `json:"estado"`
}

// Retorna timeStamp
func timeStamp() string{
	currentTime := time.Now()

	return (fmt.Sprintf("%d-%d-%d %d:%d:%d",
		currentTime.Day(),
		currentTime.Month(),
		currentTime.Year(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second()))
}

func runUsuario(conn net.Conn) error {
	topico := newUsuario()

	pongChan := make(chan struct{}, 1)
	topicoChan := make(chan Topico, 100)

	errChan := make(chan error, 3)

	// HEARTBEAT
	go func() {
		heartbeat(conn)
		errChan <- fmt.Errorf("heartbeat morreu")
	}()

	// DISPATCHER
	go func() {
		dispatcher(conn, pongChan, topicoChan)
		errChan <- fmt.Errorf("dispatcher morreu")
	}()

	// MONITOR
	go func() {
		monitorConexao(pongChan)
		errChan <- fmt.Errorf("timeout de conexão")
	}()

	// TERMINAL (loop principal do usuário)
	go func() {
		errChan <- terminal(conn, topico, topicoChan)
	}()

	// Espera qualquer erro
	return <-errChan
}

func terminal(conn net.Conn, topico *Topico, topicoChan chan Topico) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Terminal Pub/Sub Iniciado ===")

	for {
		fmt.Print("\n----- Painel de Controle -----\n[c] Publicar comando\n[s] Assinar sensor\n[t] Teste de Concorrência p/ Atuadores\nDigite o comando: ")

		if !scanner.Scan() {
			return fmt.Errorf("entrada encerrada")
		}
		opcao := scanner.Text()

		// ================== PUBLICAR ATUADOR ==================
		if opcao == "c" {

			fmt.Print("Digite o ID do dispositivo ou \"#\" para Visualizar Disponíveis: ");
			if !scanner.Scan() {
				return fmt.Errorf("entrada encerrada")
			}
			tipoId := scanner.Text()

			// Trata para o caso do usuário querer visualizar (assinar) os tópicos disponíveis de atuadores
			if tipoId == "#" {
				subTopico := *topico
				subTopico.Acao = "sub"
				subTopico.Tipo = "atuador"
				subTopico.TipoId = "#"

				stop := make(chan bool)

				go assinarTopico(conn, &subTopico, topicoChan, stop);

				for {
					if !scanner.Scan() {
						return fmt.Errorf("entrada encerrada")
					}

					if scanner.Text() == "p" {
						stop <- true
						subTopico.Acao = "unsub"
						desassinarTopico(conn, &subTopico)
						break
					}
				}
			} else { // Trata o caso que o usuário irá enviar um comando para atuador (publicar)

				fmt.Print("Digite o Comando (on/off): ")
				if !scanner.Scan() {
					return fmt.Errorf("entrada encerrada")
				}

				comando := "false"
				if scanner.Text() == "on" {
					comando = "true"
				}

				pubTopico := *topico
				pubTopico.Acao = "pub"
				pubTopico.Tipo = "atuador"
				pubTopico.TipoId = tipoId
				pubTopico.Comando = comando

				go publicarTopico(conn, &pubTopico)

				fmt.Println("Comando enviado.")
			}
		}

		// ================== ASSINAR SENSOR ==================
		if opcao == "s" {
			fmt.Print("Digite o ID do dispositivo: ")

			if !scanner.Scan() {
				return fmt.Errorf("entrada encerrada")
			}
			tipoId := scanner.Text()

			subTopico := *topico
			subTopico.Acao = "sub"
			subTopico.Tipo = "sensor"
			subTopico.TipoId = tipoId

			stop := make(chan bool)

			go assinarTopico(conn, &subTopico, topicoChan, stop)

			fmt.Println("Assinando... (digite 'p' para parar)")

			for {
				if !scanner.Scan() {
					return fmt.Errorf("entrada encerrada")
				}

				if scanner.Text() == "p" {
					stop <- true
					subTopico.Acao = "unsub"
					desassinarTopico(conn, &subTopico)
					break
				}
			}
		}

		// ================== TESTE ==================
		if opcao == "t" {

			fmt.Print("Digite o ID do dispositivo: ")
			if !scanner.Scan() {
				return fmt.Errorf("entrada encerrada")
			}
			tipoId := scanner.Text()

			fmt.Println("Publicando... (digite 'p' para parar)")

			stop := make(chan struct{})

			go func() {
				toggleCmd := true

				for {
					select {
					case <-stop:
						fmt.Println("Parando publicação")
						return

					default:
						pubTopico := *topico
						pubTopico.Acao = "pub"
						pubTopico.Tipo = "atuador"
						pubTopico.TipoId = tipoId

						if toggleCmd {
							pubTopico.Comando = "true"
						} else {
							pubTopico.Comando = "false"
						}

						toggleCmd = !toggleCmd

						publicarTopico(conn, &pubTopico)

						time.Sleep(time.Second)
					}
				}
			}()

			for {
				if !scanner.Scan() {
					return fmt.Errorf("entrada encerrada")
				}

				if scanner.Text() == "p" {
					close(stop)
					break
				}
			}
		}
	}
}

func main() {
	for {
		fmt.Println("Tentando conectar ao broker...")

		conn, err := net.Dial("tcp", ":9000")
		if err != nil {
			fmt.Println("Erro:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		fmt.Println("Conectado!")

		err = runUsuario(conn)

		fmt.Println("Conexão perdida:", err)

		conn.Close()

		time.Sleep(3 * time.Second)
	}
}