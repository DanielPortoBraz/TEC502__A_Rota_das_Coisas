package main

import (
	"fmt"
	"net"
	"sync"
	"log"
	"encoding/json"
	"time"
)

// ============== Broker ===========
/*

1. Inicia servidores TCP e UDP
2. Registra as listas de sensores, usuários e atuadores
3. Separa os tópicos
4. Roteia Sub para tópico Pub

* Fornece Logs de conexões

*/


// STRUCTS

/*
Padrão de Tópico: tipo/tipoId/comando. ("Valor" não será utilizado dentro do tópico, mas como dado fornecido por topicos. O mesmo vale para "Estado")
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
};

// Guarda os assinantes de cada tópico. Assume que não tem assinantes UDP
type Broker struct {
	assinantes map[string][]net.Conn; // Map de Assinantes (Guarda conns TCP para cada tópico)- [Topico]: [conn1, conn2, ...]
	mu sync.RWMutex;
};


// FUNCOES - Pub/Sub



// Publica tópico no Broker
func (broker *Broker) publicar(topico Topico){
	
	// Apenas rota do tópico
	chave := fmt.Sprintf("%s/%s", topico.Tipo, topico.TipoId)

	// Recebe a lista de conns do Tópico a ser publicado
	broker.mu.RLock();
	conns := broker.assinantes[chave];
	broker.mu.RUnlock();

	// Serializa o topico 
	data, err := json.Marshal(topico);
	if err != nil {
		return;
	}

	// Publica Tópico para todos os conns inscritos
	for _, c := range conns {
		go func(conn net.Conn) {
			_, err := conn.Write(data);
			if err != nil {
				return;
			}
			conn.Write([]byte("\n"));
		}(c)
	}
}

// Assina tópico no Broker
func (broker *Broker) assinar(topico Topico, conn net.Conn){
	chave := fmt.Sprintf("%s/%s", topico.Tipo, topico.TipoId)

	// Adiciona um assinante (conn) ao Tópico
	broker.mu.Lock();
	broker.assinantes[chave] = append(broker.assinantes[chave], conn);
	broker.mu.Unlock();
}


// FUNCOES GERAIS

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

// Inicia Servidor TCP
func StartServerTCP(broker *Broker) {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%s] (Broker) (TCP): Porta Aberta 9000\n", timeStamp());

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		fmt.Println("Broker: TCP: Novo dispositivo conectado:", conn.RemoteAddr())

		// Gerenciamento de Clientes TCP
		go handleConnectionTCP(conn, broker);
	}
}

// Inicia Servidor UDP
func StartServerUDP(broker *Broker) {

	address, err := net.ResolveUDPAddr("udp", ":9000");
	if err != nil {
		log.Fatal(err);
	}

	conn, err := net.ListenUDP("udp", address);
	if err != nil {
		log.Fatal(err);
	}

	fmt.Printf("[%s] (Broker) (UDP): Porta Aberta 9000\n", timeStamp());

	defer conn.Close();

	// Gerenciamento de Clientes UDP
	handleConnectionUDP(conn, broker);
}

func main() {

	// Canal de Tópicos
	broker := Broker{assinantes: make(map[string][]net.Conn)};
	

	// Iniciando Servidores TCP UDP - abre a porta, aceita e gerencia conexões
	go StartServerTCP(&broker);
	go StartServerUDP(&broker);
	
	for {		
		time.Sleep(time.Second);
	}
}