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

	data, err := json.Marshal(usuario)
	if err != nil {
		fmt.Println("(Usuario) Erro ao enviar comando")
		return
	}

	conn.Write(data)
	conn.Write([]byte("\n"))

	// Guarda último valor por tópico
	ultimos := make(map[string]Topico)

	// Garante tratamento concorrente de leitura exibida no terminal a cada 1 s
	ticker := time.NewTicker(time.Second);
	defer ticker.Stop();

	for {
		select {

		case <-stop:
			fmt.Println("(Usuario) Encerrando assinatura...");
			return

		// Recebe dados continuamente (não bloqueia fluxo)
		case topico := <-topicoChan:
			chave := fmt.Sprintf("%s/%s", topico.Tipo, topico.TipoId);
			ultimos[chave] = topico;

		// Atualiza tela a cada 1 segundo
		case <-ticker.C:

			if len(ultimos) == 0 {
				continue;
			}

			for _, t := range ultimos {
				fmt.Printf("[%s](Usuario): Topico Recebido - %s/%s/%s/%t\nValor: %.2f\n",
				timeStamp(), t.Tipo, t.TipoId,
				t.Comando, t.Estado, t.Valor);
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