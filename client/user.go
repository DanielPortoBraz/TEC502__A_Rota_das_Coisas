package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// Estado local
type User struct {
	ID string
}

func newUser() *User {
	return &User{
		ID: os.Getenv("HOSTNAME"),
	}
}

func enviarComando(conn net.Conn, target string, estado bool) {

	msg := Mensagem{
		Type :   "command",
		Target : target,
		Estado : estado,
	}

	jsonMsg, _ := json.Marshal(msg)

	conn.Write(jsonMsg)

	fmt.Println("Comando enviado:", string(jsonMsg))
}

func receberDados(conn net.Conn) {

	reader := bufio.NewReader(conn)

	for {
		msgBytes, err := reader.ReadBytes('}')
		if err != nil {
			fmt.Println("Erro leitura:", err)
			return
		}

		var msg Mensagem
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			fmt.Println("Usuário: Erro JSON:", err)
			continue
		}

		// Só processa dados de sensor
		if msg.Type == "data" {
			fmt.Printf("Usuário: Dado recebido de [%s]: %f\n", msg.Target, msg.Valor)
		}
	}
}