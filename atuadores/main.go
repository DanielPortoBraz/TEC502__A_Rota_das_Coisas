package main

import (
	"fmt"
	"net"
	"time"
)

const PORTA = "172.16.201.5"

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

// Função que encapsula toda a lógica do atuador com assinatura contínua
func runAtuador(conn net.Conn) error {
	atuador := newAtuador()
	lastPong := time.Now()

	errChan := make(chan error, 2)

	// Goroutine do heartbeat
	go func() {
		heartbeat(conn, &lastPong)
		errChan <- fmt.Errorf("heartbeat morreu")
	}()

	// Goroutine principal (assinatura)
	go func() {
		assinarComando(atuador, conn, &lastPong)
		errChan <- fmt.Errorf("leitura morreu")
	}()

	// Espera qualquer erro
	return <-errChan
}

// Main - Roda atuador e faz tentativa de Reconexão a cada 3 segundos, caso o Broker caia
func main() {

	for {
		fmt.Println("Tentando conectar...")

		conn, err := net.Dial("tcp", "localhost:9000")
		if err != nil {
			fmt.Println("Erro ao conectar:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		fmt.Println("Conectado ao broker!")

		err = runAtuador(conn)

		fmt.Println("Conexão perdida:", err)

		conn.Close()

		time.Sleep(3 * time.Second)
	}
}