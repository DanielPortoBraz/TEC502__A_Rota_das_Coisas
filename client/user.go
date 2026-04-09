package main

import (
	"fmt"
	"bufio"
	"net"
	"encoding/json"
	"os"
	"time"
	"bytes"
	"sync"
)

var writeMu sync.Mutex

// Cria struct de novo usuário (tópico)
func newUsuario() *Topico {
	return &Topico{
		Acao : "sub",
		Tipo : "usuario",
		TipoId : os.Getenv("HOSTNAME"),
		Comando : "",
		Valor : 0.0,
		Estado : false,
	}
}

// Heartbeat para manter conexão com broker utilizando "ping-pong"
func heartbeat(conn net.Conn, done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			ping_Top := Topico{Acao : "ping"}
			data, _ := json.Marshal(ping_Top)

			writeMu.Lock()
			_, err := conn.Write(data)
			if err != nil {
				writeMu.Unlock()
				return
			}
			conn.Write([]byte("\n"))
			writeMu.Unlock()

			time.Sleep(5 * time.Second)
		}
	}
}

// Montiora a conexão, caso o "pong" do broker não seja retornado
func monitorConexao(pongChan chan struct{}, done chan struct{}) {
	timeout := 10 * time.Second
	timer := time.NewTimer(timeout)

	for {
		select {
		case <-done:
			return

		case <-pongChan:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(timeout)

		case <-timer.C:
			fmt.Printf("[%s] (Usuario) Broker desconectou por timeout!\n", timeStamp())
			return
		}
	}
}

// Dispatcher que distribui as mensagens recebidas do broker em: pongs e tópicos
func dispatcher(conn net.Conn, pongChan chan struct{}, topicoChan chan Topico, done chan struct{}) {
	reader := bufio.NewReader(conn)

	for {
		select {
		case <-done:
			return
		default:
			data, err := reader.ReadBytes('\n')
			if err != nil { // Verifica erro de desconexão
				fmt.Printf("[%s](Usuario) Broker Desconectou\n", timeStamp())
				return
			}

			data = bytes.TrimSpace(data)
			if len(data) == 0 {
				continue
			}

			var topico Topico
			if err := json.Unmarshal(data, &topico); err != nil {
				continue
			}

			if topico.Acao == "pong" {
				select {
				case pongChan <- struct{}{}:
				default:
				}
				continue
			}

			if topico.Tipo == "" || topico.TipoId == "" {
				continue
			}

			select {
			case topicoChan <- topico:
			case <-done:
				return
			}
		}
	}
}

func publicarTopico(conn net.Conn, usuario *Topico) {
	usuario.Acao = "pub"

	data, err := json.Marshal(usuario)
	if err != nil {
		return
	}

	writeMu.Lock()
	conn.Write(data)
	conn.Write([]byte("\n"))
	writeMu.Unlock()
}

// Assina tópico: primeiro envia tópico a ser assinado e então aguarda continuamente leituras daquele tópico
func assinarTopico(conn net.Conn, usuario *Topico, topicoChan chan Topico, stop chan bool, done chan struct{}) {

	// Esvazia buffer antigo
	for len(topicoChan) > 0 {
		<-topicoChan
	}

	usuario.Acao = "sub"

	data, err := json.Marshal(usuario)
	if err != nil {
		return
	}

	writeMu.Lock()
	conn.Write(data)
	conn.Write([]byte("\n"))
	writeMu.Unlock()

	// Buffer para manter os dados mais recentes em exibição
	ultimos := make(map[string]Topico)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop() // Ticker, utilizado para atualizar terminal a cada 1s 

	for {
		select {

		case <-done:
			return;

		case <-stop:
			fmt.Println("(Usuario) Encerrando assinatura...")
			return

		case topico, ok := <-topicoChan:
			if !ok {
				return
			}
			chave := fmt.Sprintf("%s/%s", topico.Tipo, topico.TipoId)
			ultimos[chave] = topico

		case <-ticker.C:

			if len(ultimos) == 0 {
				fmt.Printf("[%s](Usuario): Nenhum tópico encontrado\n", timeStamp())
				continue
			}

			for _, t := range ultimos {

				if t.Valor == -1 {
					fmt.Printf("[%s](Usuario): Topico Listado - %s/%s\n",
						timeStamp(), t.Tipo, t.TipoId)

				} else {
					fmt.Printf("[%s](Usuario): Topico Recebido - %s/%s/\nEstado: %t\nValor: %.2f\n",
						timeStamp(), t.Tipo, t.TipoId,
						t.Estado, t.Valor)
				}
			}
		}
	}
}

func desassinarTopico(conn net.Conn, usuario *Topico) {
	usuario.Acao = "unsub"

	data, _ := json.Marshal(usuario)

	writeMu.Lock()
	conn.Write(data)
	conn.Write([]byte("\n"))
	writeMu.Unlock()
}