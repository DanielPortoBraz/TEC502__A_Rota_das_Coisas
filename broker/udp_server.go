package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func handleConnectionUDP (conn *net.UDPConn, pub chan Topico) {

	// Buffer para receber os dados
	buffer := make([]byte, 1024);
	var sensor Topico;

	for {
		n, _, err := conn.ReadFromUDP(buffer) // Recebe dados UDP
		if err != nil {
			fmt.Printf("Broker: UDP: Error: %v\n\n", err)
			continue
		}

		// Desserializa JSON
		json.Unmarshal(buffer[:n], &sensor);

		// Envia Tópico do Sensor para o canal de Topicos
		pub <- sensor;

		fmt.Printf("[%s] (Broker) (UDP):\n%s/%s/%s\n", timeStamp(), sensor.Tipo, sensor.TipoId, sensor.Comando);
	}
}