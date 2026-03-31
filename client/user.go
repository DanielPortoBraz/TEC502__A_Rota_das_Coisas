package main

import (
	"fmt"
	"bufio"
	"net"
	"encoding/json"
	"os"
	"time"
)

func newUsuario() *Topico {
	return &Topico{
		Acao : "sub", // Definido inicialmente como sub
		Tipo : "usuario",
		TipoId : os.Getenv("HOSTNAME"),	
		Comando : "",
		Valor : 0.0,
		Estado : false,
	}
}

// Heartbeat para indicar que a conexão está viva
func heartbeat(conn net.Conn) {
	for {
		ping_Top := Topico{Acao : "ping"};
		data, _ := json.Marshal(ping_Top);
		
		_, err := conn.Write(data);
		if err != nil {
			return // conexão morreu
		}

		conn.Write([]byte("\n"));
		time.Sleep(5 * time.Second);
	}
}

// Monitora mensagens "pong" para indicar conexão ativa do Broker
func monitorConexao(pongChan chan struct{}) {
	timeout := 10 * time.Second; // A cada 10 segundos
	timer := time.NewTimer(timeout)

	for {
		select {
		case <-pongChan:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(timeout) // Reinicia o cronômetro sempre que o pong chega

		case <-timer.C:
			fmt.Printf("[%s] (Usuario) Broker desconectou por timeout!\n", timeStamp())
			return
		}
	}
}

// Recebe o pacote enviado pelo Broker e identifica se é um heartbeat ("pong") ou tópico a ser assinado
func dispatcher(conn net.Conn, pongChan chan struct{}, topicoChan chan Topico) {
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("[%s](Usuario) Broker Desconectou\n", timeStamp());
			close(pongChan)
			close(topicoChan)
			return
		}

		var topico Topico
		if err := json.Unmarshal(data, &topico); err != nil {
			continue
		}

		// Roteamento
		if topico.Acao == "pong" {
			pongChan <- struct{}{};
		} else {
			topicoChan <- topico;
		}
	}
}



func publicarTopico(conn net.Conn, usuario *Topico) {

	usuario.Acao = "pub";

	data, err := json.Marshal(usuario);
	if err != nil {
		fmt.Println("(Usuario) Erro ao enviar comando");
	}

	conn.Write(data)
	conn.Write([]byte("\n"))
}

func assinarTopico(conn net.Conn, usuario *Topico, topicoChan chan Topico, stop chan bool) {

	usuario.Acao = "sub"

	// Envia topico para assinar
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

		case topico := <-topicoChan:
			// Esvazia FIFO evitando obter um dado muito atrasado
			for len(topicoChan) > 0 {
                topico = <-topicoChan
            }

			fmt.Printf("[%s](Usuario): Topico Recebido - %s/%s/%s/%t\nValor: %.2f\n",
				timeStamp(), topico.Tipo, topico.TipoId,
				topico.Comando, topico.Estado, topico.Valor)

			// Recebe a leitura do sensor a cada 1 segundo (Evita sobrecarga no terminal)
			time.Sleep(time.Second);
		}
	}
}

func desassinarTopico(conn net.Conn, usuario *Topico) {
	usuario.Acao = "unsub"

	data, _ := json.Marshal(usuario)
	conn.Write(data)
	conn.Write([]byte("\n"))
}