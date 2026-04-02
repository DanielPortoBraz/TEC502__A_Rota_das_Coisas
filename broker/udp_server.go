package main

import (
	"encoding/json"
	"fmt"
	"net"
)

func handleConnectionUDP (conn *net.UDPConn, broker *Broker) {

	// Buffer para receber os dados
	buffer := make([]byte, 1024);
	var topico Topico;

	go broker.monitorarTopicos();
	
	for {
		n, _, err := conn.ReadFromUDP(buffer) // Recebe dados UDP
		if err != nil {
			fmt.Printf("[%s] (Broker) (UDP): Error: %v\n\n", timeStamp(), err);
			continue
		}
		
		// Desserializa JSON
		json.Unmarshal(buffer[:n], &topico);
		
		// Publica Tópico (Sensor) 
		broker.publicar(topico);
		
		fmt.Printf("[%s] (Broker): Publicação - %s/%s/%s\nValor: %.2f\n", timeStamp(), topico.Tipo, topico.TipoId, topico.Comando, topico.Valor);

	}
}