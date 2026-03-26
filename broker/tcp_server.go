package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

func handleConnectionTCP(conn net.Conn, ch chan Topico) {
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

		switch topico.Acao {

		case "pub":
			fmt.Println("Publicação")
			ch <- topico;

		case "sub":

			assinado := (<- ch);

			if topico.TipoId == assinado.TipoId {
				fmt.Println("Assinatura")
				data, _ = json.Marshal(assinado);

				// Envia dados da assinatura
				conn.Write(data);
			}
		}

		fmt.Printf("[%s] (Broker) (TCP):\n%s/%s/%s\n", timeStamp(), topico.Tipo, topico.TipoId, topico.Comando);
	}
}
