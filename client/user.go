package main

import (
	"fmt"
	"bufio"
	"net"
	"encoding/json"
	"math/rand"
	"time"
)

func newUsuario() *Topico {
	return &Topico{
		Acao : "sub", // Definido inicialmente como sub
		Tipo : "usuario",
		TipoId : fmt.Sprintf("%d", rand.Intn(100)),
		Comando : "",
		Valor : 0.0,
		Estado : false,
	}
}

// Heartbeat para indicar que a conexão está viva
func heartbeat(conn net.Conn) {
	for {
		ping_Top := Topico{Acao : "ping"}
		data, _ := json.Marshal(ping_Top)
		
		_, err := conn.Write(data)
		if err != nil {
			return // conexão morreu
		}

		conn.Write([]byte("\n"))
		time.Sleep(5 * time.Second)
	}
}

func publicarTopico(conn net.Conn, usuario *Topico) {

	usuario.Acao = "pub"

	data, err := json.Marshal(usuario)
	if err != nil {
		fmt.Println("(Usuario) Erro ao enviar comando")
	}

	conn.Write(data)
	conn.Write([]byte("\n"))
}

func assinarTopico(conn net.Conn, usuario *Topico, stop chan bool){

	reader := bufio.NewReader(conn)
	var topico Topico

	usuario.Acao = "sub"

	// Envia topico para se inscrever em outro
	data, err := json.Marshal(usuario)
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
			if err := json.Unmarshal(data, &topico); err != nil {
				continue
			}

			// Resposta do Broker para Heartbeat
			if topico.Acao == "pong" {
				fmt.Printf("[%s] (Usuario): %v\n", timeStamp(), topico)

			} else { // Tópico assinado
				fmt.Printf("[%s](Usuario): Topico Recebido - %s/%s/%s/%t\nValor: %.2f\n",
					timeStamp(), topico.Tipo, topico.TipoId,
					topico.Comando, topico.Estado, topico.Valor)
			}
		}
	}
}

func desassinarTopico(conn net.Conn, usuario *Topico) {
	usuario.Acao = "unsub"

	data, _ := json.Marshal(usuario)
	conn.Write(data)
	conn.Write([]byte("\n"))
}