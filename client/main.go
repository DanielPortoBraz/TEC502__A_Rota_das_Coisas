package main

import (
	"fmt"
	"net"
	"time"
	"bufio"
	"os"
)

// Obtém o IP da máquina que roda o Broker Servidor
func getBrokerHost() string {
	host := os.Getenv("BROKER_HOST")
	if host == "" {
		host = "broker"
	}
	return host
}

type Topico struct {
	Acao string `json:"acao"`
	Tipo string `json:"tipo"`
	TipoId string `json:"tipoId"`
	Comando string `json:"comando"`
	Valor float64 `json:"valor"`
	Estado bool `json:"estado"`
}

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

func runUsuario(conn net.Conn, c int) error {
	topico := newUsuario()

	pongChan := make(chan struct{}, 1)
	topicoChan := make(chan Topico, 100)
	errChan := make(chan error, 4)
	done := make(chan struct{})

	go func() {
		heartbeat(conn, done)
		errChan <- fmt.Errorf("heartbeat morreu")
	}()

	go func() {
		dispatcher(conn, pongChan, topicoChan, done)
		errChan <- fmt.Errorf("dispatcher morreu")
	}()

	go func() {
		monitorConexao(pongChan, done)
		errChan <- fmt.Errorf("timeout de conexão")
	}()

	go func() {
		errChan <- terminal(conn, topico, topicoChan, c, done)
	}()

	err := <-errChan
	close(done)
	return err
}

func terminal(conn net.Conn, topico *Topico, topicoChan chan Topico, c int, done chan struct{}) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Terminal Pub/Sub Iniciado ===")

	for {
		select {
		case <-done:
			return fmt.Errorf("finalizado")
		default:
		}

		fmt.Print("\n----- Painel de Controle -----\n[c] Publicar comando\n[s] Assinar sensor\n[t] Teste de Concorrência p/ Atuadores\nDigite o comando: ")

		if !scanner.Scan() {
			return fmt.Errorf("entrada encerrada")
		}
		opcao := scanner.Text()

		if opcao == "c" {

			fmt.Print("Digite o ID do dispositivo ou \"#\" para Visualizar Disponíveis: ")
			if !scanner.Scan() {
				return fmt.Errorf("entrada encerrada")
			}
			tipoId := scanner.Text()

			if tipoId == "#" {
				subTopico := *topico
				subTopico.Acao = "sub"
				subTopico.Tipo = "atuador"
				subTopico.TipoId = "#"

				stop := make(chan bool)

				go func() {
					for {
						select {
						case <-done:
							return
						default:
							if !scanner.Scan() {
								return
							}
							if scanner.Text() == "p" {
								select {
									case stop <- true:
									case <-done:
								}
								return
							}
						}
					}
				}()

				assinarTopico(conn, &subTopico, topicoChan, stop, done)

				subTopico.Acao = "unsub"
				desassinarTopico(conn, &subTopico)

			} else {

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

			go func() {
				for {
					select {
					case <-done:
						return
					default:
						if !scanner.Scan() {
							return
						}
						if scanner.Text() == "p" {
							select {
								case stop <- true:
								case <-done:
							}
							return
						}
					}
				}
			}()

			assinarTopico(conn, &subTopico, topicoChan, stop, done)

			subTopico.Acao = "unsub"
			desassinarTopico(conn, &subTopico)
		}

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
						return
					case <-done:
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

	c := 1

	for {
		fmt.Println("Tentando conectar ao broker...")

		conn, err := net.Dial("tcp", fmt.Sprintf("%s:9000", getBrokerHost()))
		if err != nil {
			fmt.Println("Erro:", err)
			time.Sleep(3 * time.Second)
			c++
			continue
		}

		fmt.Println("Conectado!")

		err = runUsuario(conn, c)

		fmt.Println("Conexão perdida:", err)

		conn.Close()
		time.Sleep(3 * time.Second)
	}
}