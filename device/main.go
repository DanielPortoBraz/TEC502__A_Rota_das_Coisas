package main

import (
	"fmt"
	"net"
	"time"
)

// 1. Conecta a Porta por UDP
// 2. Cria Sensor
// 3. Envia leitura em loop contínuo (formato UDP)

func main() {

	// Conexão por UDP
    conn, err := net.Dial("udp", "localhost:9000");
	if err != nil {
		panic(err);
	}
	
	fmt.Println("Porta Aberta")

	for {
		enviarDado(conn, lerDado());
		time.Sleep(1 * time.Second); // Envia uma leitura a cada segundo
	}
}