package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

func handleConnectionTCP(conn net.Conn) {
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

		fmt.Println("TCP: Dado Recebido\n\n")

	}
}
