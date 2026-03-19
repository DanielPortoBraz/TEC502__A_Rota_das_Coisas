package main

import (
	"fmt"
	"net"
)

type SensorData struct {
	ID    string  `json:"id"`
	Valor float64 `json:"valor"`
}

func StartServer(sensores chan SensorData) {
	ln, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}

	fmt.Println("Broker: Porta Aberta 9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		fmt.Println("Novo dispositivo conectado:", conn.RemoteAddr())
		go handleConnection(conn, sensores)
	}
}

func main() {
	sensores := make(chan SensorData, 5); 

	// Iniciando Servidor - abre a porta, aceita e gerencia conexões
	go StartServer(sensores);

	// Tratativa dos sensores - Recebe dados do canal de sensores
	for sensor := range sensores {
		fmt.Println("Broker: Dado Recebido\n", sensor);
	}

}