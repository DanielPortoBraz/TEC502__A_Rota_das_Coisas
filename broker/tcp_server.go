package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn, ch chan SensorData) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	var sensor SensorData

	for {
		data, err := reader.ReadBytes('}') // ReadBytes é bloqueante
		if err != nil {
			return
		}

		// Desserialização do JSON
		json.Unmarshal(data, &sensor)

		ch <- sensor; // Envia o sensor recebido no canal de sensores do Broker

		fmt.Println("TCP: Dado Recebido")

	}
}
