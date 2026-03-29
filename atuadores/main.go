package main

import (
	"fmt"
	"net"
	"time"
)

// ============ Atuador =============
/*

1. Conecta ao Broker via TCP
2. Cria Atuador: Fornece atuadorID ao Broker
3. Assina tópico: atuador/ *atuadorID* / comando (Aguarda mensagem por ReadBytes e decodifica o tópico)
4. Executa ação (mudança de estado)
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


func main() {

	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close();

	fmt.Println("Atuador conectado ao broker");

	atuador := newAtuador();
	lastPong := time.Now();

	go heartbeat(conn, &lastPong);

	go assinarComando(atuador, conn, &lastPong);

	for {
		fmt.Printf("[%s] (Atuador): ID- %s | Estado: %t\n", timeStamp(), atuador.TipoId, atuador.Estado);
		time.Sleep(time.Second);
	}
}