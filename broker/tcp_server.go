package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

func handleConnectionTCP(conn net.Conn, broker *Broker) {
	defer conn.Close()

	reader := bufio.NewReader(conn);
	var topico Topico;

	for {
		data, err := reader.ReadBytes('}') // ReadBytes é bloqueante
		if err != nil {
			return
		}

		// Desserialização do JSON
		json.Unmarshal(data, &topico)

		// Verifica se é um tópico de publicação ou assinatura
		switch topico.Acao {

		case "pub":
			// Publica tópico
			broker.publicar(topico);
			
		case "sub":
			// Assina tópico
			broker.assinar(topico, conn);
		}

		fmt.Printf("[%s] (Broker) (TCP):\n%s/%s/%s\n", timeStamp(), topico.Tipo, topico.TipoId, topico.Comando);
	}
}
