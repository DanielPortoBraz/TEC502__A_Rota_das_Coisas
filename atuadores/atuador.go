package main

import (
	"fmt"
	"net"
	"bufio"
	"encoding/json"
	"time"
	"os"
)


func newAtuador() *Topico {
	return &Topico{
		Acao : "sub",
		Tipo : "atuador",
		TipoId : os.Getenv("HOSTNAME"),
		Comando : "",
		Valor : 0.0,
		Estado: false,
	}
}

// Heartbeat para indicar que a conexão está viva
func heartbeat(conn net.Conn, lastPong *time.Time) {
	for {
		ping_Top := Topico{Acao : "ping"};
		data, _ := json.Marshal(ping_Top);

		_, err := conn.Write(data) 
		if err != nil {
			return // conexão morreu
		}
		conn.Write([]byte("\n"));

		time.Sleep(5 * time.Second)

		// Verifica se a conexão com Broker caiu nos últimos 10 segundos
		if time.Since(*lastPong) >= 10 * time.Second{
			fmt.Printf("\n\n[%s] (Atuador): Broker desconectou!\n\n", timeStamp());
			return
		}
	}
}

// Assina o Tópico e atualiza o estado do atuador
func assinarComando(atuador *Topico, conn net.Conn, lastPong *time.Time) {
	
	// 1. Publica tópico do atuador
	// 2. Espera continuamente até que alguém assine
	// 3. Atualiza o estado do atuador conforme aquele do tópico assinado
	
	// Serializa JSON - Tópico a ser publicado
	atuador_json, err := json.Marshal(atuador);
	if err != nil {
		return;
	}

	// Publica Tópico do atuador
	conn.Write(atuador_json);
	conn.Write([]byte("\n"));

	
	reader := bufio.NewReader(conn)

	// Assina próprio tópico continuamente
	var topico Topico;

	for {
		fmt.Printf("[%s] (Atuador): ID- %s | Estado: %t\n", timeStamp(), atuador.TipoId, atuador.Estado);
		data, err := reader.ReadBytes('\n');
		if err != nil {
			return;
		}
		
		// Desserializar JSON
		json.Unmarshal(data, &topico);
		
		// Resposta do Broker para Heartbeat
		if topico.Acao == "pong" {
			*lastPong = time.Now();
			fmt.Printf("[%s] (Atuador) Acao: %s\n", timeStamp(), topico.Acao);

			
		} else{ // Atualiza estado
			if topico.Comando == "true" {
				atuador.Estado = true;
			} else {
				atuador.Estado = false;
			}
			fmt.Printf("[%s] (Atuador): Topico assinado - %s/%s\nComando: %s\n", timeStamp(), topico.Tipo, topico.TipoId, topico.Comando);
		}
	}
}