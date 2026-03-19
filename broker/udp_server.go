package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func handleConnectionUDP (conn *net.UDPConn, ch chan SensorData) {

	// Buffer para receber os dados
	buffer := make([]byte, 1024);
	var sensor SensorData;

	for {
		n, _, err := conn.ReadFromUDP(buffer) // Recebe dados UDP
		if err != nil {
			fmt.Printf("UDP: Error: %v\n\n", err)
			continue
		}

		// Desserializa JSON
		json.Unmarshal(buffer[:n], &sensor);

		// Envia dados de leitura para o canal de Sensores
		ch <- sensor;

		fmt.Println("UDP: Dado recebido\n");
	}
}