package main

import (
	"fmt"
	"bufio"
	"net"
	"encoding/json"
	"time"
)

func newUsuario() *Topico {
	return &Topico{
		Acao : "sub", // Definido inicialmente como sub
		Tipo : "usuario",
		TipoId : "456",
		Comando : "",
		Valor : 0.0,
		Estado : false,
	}
}

func publicarTopico(conn net.Conn, topico *Topico) {

	topico.Acao = "pub";

	data, err := json.Marshal(topico);
	if err != nil {
		fmt.Println("(Usuario) Erro ao enviar comando")
	}

	conn.Write(data);

}

func assinarTopico(conn net.Conn, topico *Topico, stop chan bool){

	reader := bufio.NewReader(conn)
	var comando Topico

	topico.Acao = "sub"

	// Envia topico para se inscrever em outro
	data, err := json.Marshal(topico)
	if err != nil{
		fmt.Println("(Usuario) Erro ao enviar comando")
		return
	}
	conn.Write(data)
	conn.Write([]byte("\n"))

	for {
		select {
		case <-stop:
			fmt.Println("(Usuario) Encerrando assinatura...")
			return

		default:
			// Define timeout para não travar
			conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

			data, err := reader.ReadBytes('\n')
			if err != nil {

				// timeout esperado → continua loop
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}

				fmt.Printf("[%s] (Usuario): Erro\n", timeStamp())
				return
			}

			// Desserialização do JSON
			if err := json.Unmarshal(data, &comando); err != nil {
				continue
			}

			fmt.Printf("[%s](Usuario): Topico Recebido:\n%s/%s/%s/%t\nValor: %.2f\n",
				timeStamp(), comando.Tipo, comando.TipoId,
				comando.Comando, comando.Estado, comando.Valor)
		}
	}
}