package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

func newUsuario() *Topico {
	return &Topico{
		Acao : "sub", // Definido inicialmente como sub
		Tipo : "usuario",
		TipoId : "456",
		Comando : "",
		Valor : 0.0,
		Estado : false,
	}
}

func publicarTopico(conn net.Conn, topico *Topico) {

	topico.Acao = "pub";

	data, err := json.Marshal(topico);
	if err != nil {
		fmt.Println("(Usuario) Erro ao enviar comando")
	}

	conn.Write(data);

}

func assinarTopico(conn net.Conn, topico *Topico){

	reader := bufio.NewReader(conn);
	var comando Topico;

	topico.Acao = "sub";

	// Envia topico para se inscrever em outro
	data, err := json.Marshal(topico);
	if err != nil{
		fmt.Println("(Usuario) Erro ao enviar comando)")
	}
	conn.Write(data);


	// Assinatura do topico
	data, err = reader.ReadBytes('}') // ReadBytes é bloqueante
	if err != nil {
		fmt.Printf("[%s] (Usuario): Erro", timeStamp());
	}

	// Desserialização do JSON
	json.Unmarshal(data, &comando)
	fmt.Printf("[%s](Usuario): Topico Recebido:\n%s/%s/%s/%t\nValor: %.2f\n", timeStamp(), comando.Tipo, comando.TipoId, comando.Comando, comando.Estado, comando.Valor);

}