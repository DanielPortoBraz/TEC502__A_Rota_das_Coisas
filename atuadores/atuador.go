package main

import (
	"fmt"
	"net"
	"bufio"
	"encoding/json"
	"os"
)

type Atuador struct {
	ID     string
	Estado bool
}

// Comando
 
type Comando struct {
	Target string `json:"target"`
	Estado bool   `json:"estado"`
}

func newAtuador() *Atuador {
	return &Atuador{
		ID:     os.Getenv("HOSTNAME"),
		Estado: false,
	}
}

func (a *Atuador) setEstado(estado bool) {
	a.Estado = estado
}

func receberComando(a *Atuador, conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		msgBytes, err := reader.ReadBytes('}')
		if err != nil {
			fmt.Println("Erro leitura:", err)
			return
		}

		var cmd Comando
		err = json.Unmarshal(msgBytes, &cmd)
		if err != nil {
			fmt.Println("Erro JSON:", err)
			continue
		}

		if cmd.Target != a.ID {
			continue
		}

		// Atualiza estado
		a.setEstado(cmd.Estado)

		fmt.Printf("Atuador [%s] atualizado → Estado: %t\n", a.ID, a.Estado)
	}
}