package main

import (
	"fmt"
	"net"
	"log"
)

type SensorData struct {
	ID    string  `json:"id"`
	Valor float64 `json:"valor"`
}

// Inicia Servidor TCP
func StartServerTCP() {
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

		// Gerenciamento de Clientes TCP
		go handleConnectionTCP(conn);
	}
}

// Inicia Servidor UDP
func StartServerUDP(sensores chan SensorData) {

	address, err := net.ResolveUDPAddr("udp", "localhost:9000");
	if err != nil {
		log.Fatal(err);
	}

	conn, err := net.ListenUDP("udp", address);
	if err != nil {
		log.Fatal(err);
	}
	defer conn.Close();

	// Gerenciamento de Clientes UDP
	handleConnectionUDP(conn, sensores);
}

func main() {
	sensores := make(chan SensorData, 5); 

	// Iniciando Servidores TCP UDP - abre a porta, aceita e gerencia conexões
	go StartServerTCP();
	go StartServerUDP(sensores);

	// Tratativa dos sensores - Recebe dados do canal de sensores
	for sensor := range sensores {
		fmt.Printf("Broker: Dado Recebido\nID: %s Valor: %.2f\n\n", sensor.ID, sensor.Valor);
	}

}