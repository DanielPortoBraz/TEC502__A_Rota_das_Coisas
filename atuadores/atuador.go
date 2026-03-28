package main

import (
	"fmt"
	"net"
	"bufio"
	"encoding/json"
	"math/rand"
)


func newAtuador() *Topico {
	return &Topico{
		Acao : "sub",
		Tipo : "atuador",
		TipoId : fmt.Sprintf("%d", rand.Intn(100)),
		Comando : "-",
		Valor : 0.0,
		Estado: false,
	}
}


// Assina o Tópico e atualiza o estado do atuador
func assinarComando(atuador *Topico, conn net.Conn) {
	
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

	// Assina tópico continuamente
	var topico Topico;
	for {
		fmt.Println("Assinando")
		data, err := reader.ReadBytes('\n');
		if err != nil {
			return;
		}
		
		// Desserializar JSON
		json.Unmarshal(data, &topico);
		fmt.Println("(Atuador): Topico assinado -", topico);

		// Atualiza estado
		if topico.Comando == "true" {
			atuador.Estado = true;
		} else {
			atuador.Estado = false;
		}
	}
}