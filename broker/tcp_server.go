package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func handleConnectionTCP(conn net.Conn, broker *Broker) {
	
	// Caso ocorra uma desconexão (timeout ou falha na rede), a conexão é removida do map de assinantes
	defer func() {
		fmt.Printf("[%s] (Broker) (TCP) Encerrando conexão: %v\n", timeStamp(), conn.RemoteAddr());
		broker.removerConn(conn);
		conn.Close();
	}()

	reader := bufio.NewReader(conn);
	var topico Topico;

	for {
		// Timeout de 10 segundos
		conn.SetReadDeadline(time.Now().Add(10 * time.Second));

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
			fmt.Printf("[%s] (Broker): Publicação - %s/%s/%s\n", timeStamp(), topico.Tipo, topico.TipoId, topico.Comando);

		case "sub":
			// Assina tópico
			broker.assinar(topico, conn);
			fmt.Printf("[%s] (Broker): Assinatura - %s/%s/%s\n", timeStamp(), topico.Tipo, topico.TipoId, topico.Comando);

		case "unsub":
			// Retira assinatura do tópico
			broker.removerConn(conn);
			fmt.Printf("[%s] (Broker): Retira Assinatura\nTópico - %s/%s/%s\nDispositivo - %v\n", timeStamp(), 
				topico.Tipo, topico.TipoId, topico.Comando,
				conn.RemoteAddr());
			
		case "ping":
			// Responde a um ping
			ping_Top := Topico{Acao : "pong"};
			data, _ := json.Marshal(ping_Top);
			
			_, err := conn.Write(data) 
			if err != nil {
				return // conexão morreu
			}
			conn.Write([]byte("\n"));
		}

		fmt.Printf("[%s] (Broker) (TCP): Acao - %s\nID: %v\nTópico: %s/%s/%s\n", timeStamp(), topico.Acao, conn.RemoteAddr(),
			topico.Tipo, topico.TipoId, topico.Comando);
	}
}
