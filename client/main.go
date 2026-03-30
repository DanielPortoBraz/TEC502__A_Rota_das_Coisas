package main

import (
	"fmt"
	"net"
	"time"
	"bufio"
	"os"
)

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


func main() {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close();

	topico := newUsuario();

	pongChan := make(chan struct{}, 1); // Canal que recebe "pong" do Broker
	topicoChan := make(chan Topico, 100); // Canal que recebe tópicos assinados

	go heartbeat(conn);
	
	go dispatcher(conn, pongChan, topicoChan)
	go monitorConexao(pongChan);


	scanner := bufio.NewScanner(os.Stdin);

	fmt.Println("=== Terminal Pub/Sub Iniciado ===")

	for {
		fmt.Print("\n----- Painel de Controle -----\n[c] Publicar comando\n[s] Assinar sensor\n[t] Teste de Concorrência p/ Atuadores\nDigite o comando: ")
		if !scanner.Scan() { break }
		opcao := scanner.Text()

		fmt.Print("Digite o ID do dispositivo: ")
		if !scanner.Scan() { break }
		tipoId := scanner.Text()

		// Publicar Tópico de Comando para Atuador
		if opcao == "c" {
			fmt.Print("Digite o Comando (on/off): ")
			if !scanner.Scan() { break } 	
			var comando string;

			if scanner.Text() == "on" {
				comando = "true";
			}else {
				comando = "false";
			}

			topico.Acao = "pub";
			topico.Tipo = "atuador";
			topico.TipoId = tipoId;
			topico.Comando = comando;

			go publicarTopico(conn, topico);
			fmt.Println("Comando enviado.");

		// Assinar Tópico de Sensor
		} else if opcao == "s" {
			topico.Acao = "sub"
			topico.Tipo = "sensor"
			topico.TipoId = tipoId

			stop := make(chan bool)

			go assinarTopico(conn, topico, topicoChan, stop);

			fmt.Println("Assinando... (digite 'p' para parar)")

			for {
				if !scanner.Scan() { break }
				cmd := scanner.Text()

				if cmd == "p" {
					stop <- true
					topico.Acao = "unsub";
					desassinarTopico(conn, topico)
					break
				}
			}
		} else if opcao == "t" {

			for {
				topico.Acao = "pub";
				topico.Tipo = "atuador";
				topico.TipoId = "123";

				fmt.Println("Publicando... (digite 'p' para parar)");

				stop := make(chan struct{});
				
				// Publica a cada segundo
				go func() {
					toggleCmd := true; // Variável para teste de concorrência - Alterna comando de estado para atuador

					for {
						select {
						case <-stop:
							fmt.Println("Parando publicação")
							return

						default:
							if toggleCmd {
								topico.Comando = "true"
							} else {
								topico.Comando = "false"
							}

							toggleCmd = !toggleCmd

							publicarTopico(conn, topico)
							time.Sleep(time.Second)
						}
					}
				}()

				for {
					if !scanner.Scan() { break }
					cmd := scanner.Text()

					if cmd == "p" {
						close(stop);
						break
					}
				}

				break;
			}

		}


	}
}