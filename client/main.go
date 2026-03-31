package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"net/http"
)

// ============= Usuários ==============

/*

1. Conecta ao Broker via TCP
2. Cria Usuário: Fornece usuarioID ao Broker
3. Executa Terminal
  L> Aguarda usuário definir qual tópico irá assinar - Deve informar no modelo:
	"tipo/tipoId/comando"
4. Envia comando, ou recebe valor/estado de sensor/atuador

*/

/*
Padrão de Tópico: tipo/tipoId/comando. ("Valor" não será utilizado dentro do tópico, mas como dado fornecido por sensores. O mesmo vale para "Estado" para o caso de atuadores)
Exemplo: sensor/sensor_1/-
Exemplo: atuador/atuador_1/off
*/
type Topico struct {
	Acao string `json:"acao"`
	Tipo string `json:"tipo"`
	TipoId string `json:"tipoId"`
	Comando string `json:"comando"`
	Valor float64 `json:"valor"`
	Estado bool `json:"estado"`
}

// Retorna timeStamp
func timeStamp() string{
	currentTime := time.Now()

	return (fmt.Sprintf("%d-%d-%d %d:%d:%d",
		currentTime.Day(),
		currentTime.Month(),
		currentTime.Year(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second()))
}


var conn net.Conn
var usuario *Topico

var topicoChan = make(chan Topico, 100)

func main() {
	var err error

	conn, err = net.Dial("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	usuario = newUsuario()

	pongChan := make(chan struct{}, 1)

	go heartbeat(conn)
	go dispatcher(conn, pongChan, topicoChan)
	go monitorConexao(pongChan)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/pub", handlePublicar)
	http.HandleFunc("/sub", handleSub)
	http.HandleFunc("/stream", handleStream)

	fmt.Println("Servidor em http://localhost:8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func handlePublicar(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Comando string `json:"comando"`
		TipoId  string `json:"tipoId"`
	}

	json.NewDecoder(r.Body).Decode(&data)

	usuario.Acao = "pub"
	usuario.Tipo = "atuador"
	usuario.TipoId = data.TipoId
	usuario.Comando = data.Comando

	publicarTopico(conn, usuario)

	w.Write([]byte("OK"))
}

func handleSub(w http.ResponseWriter, r *http.Request) {
	tipoId := r.URL.Query().Get("id")

	usuario.Acao = "sub"
	usuario.Tipo = "sensor"
	usuario.TipoId = tipoId

	data, _ := json.Marshal(usuario)
	conn.Write(data)
	conn.Write([]byte("\n"))

	w.Write([]byte("Subscrito"))
}

// STREAM DE DADOS (SSE)
func handleStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()

	for {
		<-ticker.C

		var topico Topico

		// pega o dado mais recente
		select {
		case topico = <-topicoChan:
			for len(topicoChan) > 0 {
				topico = <-topicoChan
			}
		default:
			continue
		}

		jsonData, _ := json.Marshal(topico)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		w.(http.Flusher).Flush()
	}
}