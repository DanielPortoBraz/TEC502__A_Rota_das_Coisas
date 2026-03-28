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
	scanner := bufio.NewScanner(os.Stdin);

	fmt.Println("=== Terminal Pub/Sub Iniciado ===")

	for {
		fmt.Print("\nPublicar comando ou Assinar sensor [c / s]: ")
		if !scanner.Scan() { break }
		opcao := scanner.Text()

		fmt.Print("Digite o ID do dispositivo: ")
		if !scanner.Scan() { break }
		IdSensor := scanner.Text()

		if opcao == "c" {
			topico.Tipo = "atuador";
			topico.Comando = IdSensor;
			go publicarTopico(conn, topico);
			fmt.Println("comando enviado.");

		} else if opcao == "s" {
			topico.Tipo = "sensor"
			topico.TipoId = IdSensor

			stop := make(chan bool)

			go assinarTopico(conn, topico, stop)

			fmt.Println("Assinando... (digite 'p' para parar)")

			for {
				if !scanner.Scan() { break }
				cmd := scanner.Text()

				if cmd == "p" {
					stop <- true
					break
				}
			}
		}
	}
}
