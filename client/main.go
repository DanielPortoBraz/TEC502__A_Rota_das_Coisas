package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
)

// 1. Conectar ao servidor do broker
// 2. Esperar tipo de dispositivo que o usuário deseja se conectar
// 3. Enviar comando/mensagem para o broker


// Mensagem de rede
type Mensagem struct {
	Type   string `json:"type"`
	Target string `json:"target"`
	Estado bool `json:"estado"` // usado para definir estado do atuador
	Valor  float64 `json:"valor,omitempty"` // usado para dados de sensor
}


func main() {

	conn, err := net.Dial("tcp", "broker:9000")
	if err != nil {
		panic(err)
	}

	user := newUser()

	fmt.Println("Usuário conectado:", user.ID)

	go func() {
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("Digite comando (on/off): ")
			text, _ := reader.ReadString('\n')

			if text == "on\n" {
				enviarComando(conn, "atuador_1", true)
			} else if text == "off\n" {
				enviarComando(conn, "atuador_1", false)
			}
		}
	}()

	select {} // mantém rodando
}