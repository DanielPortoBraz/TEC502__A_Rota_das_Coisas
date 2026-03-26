package main

import (
	"fmt"
	"net"
	"log"
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
func StartServerTCP(topicos chan Topico) {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	fmt.Println("Broker: TCP: Porta Aberta 9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		fmt.Println("Broker: TCP: Novo dispositivo conectado:", conn.RemoteAddr())

		// Gerenciamento de Clientes TCP
		go handleConnectionTCP(conn, topicos);
	}
}

// Inicia Servidor UDP
func StartServerUDP(topicos chan Topico) {

	address, err := net.ResolveUDPAddr("udp", ":9000");
	if err != nil {
		log.Fatal(err);
	}

	conn, err := net.ListenUDP("udp", address);
	if err != nil {
		log.Fatal(err);
	}
	defer conn.Close();

	// Gerenciamento de Clientes UDP
	handleConnectionUDP(conn, topicos);
}

func main() {

	// Canal de Tópicos
	topicos := make(chan Topico, 100); 

	// Iniciando Servidores TCP UDP - abre a porta, aceita e gerencia conexões
	go StartServerTCP(topicos);
	go StartServerUDP(topicos);

	
	for {
		
	}
}