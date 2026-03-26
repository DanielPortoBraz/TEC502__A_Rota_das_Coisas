package main

import (
	"fmt"
	"net"
	"time"
	"bufio"
	"os"
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
Padrão de Tópico: tipo/tipoId/comando. ("Valor" não será utilizado dentro do tópico, mas como dado fornecido por sensores. O mesmo vale para "Estado")
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

// 1. Conecta ao Broker via TCP
// 2. Exibe terminal
// 3. Recebe comando 
// 4. Envia comando para o Broker

func main() {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close();

	topico := newUsuario();
	scanner := bufio.NewScanner(os.Stdin);

	fmt.Println("=== Terminal Pub/Sub Iniciado ===")

	for {
		fmt.Print("\nPublicar comando ou Assinar sensor [c / s]: ")
		if !scanner.Scan() { break }
		opcao := scanner.Text()

		fmt.Print("Digite o ID do sensor (ex: sensor/1/on): ")
		if !scanner.Scan() { break }
		IdSensor := scanner.Text()

		// Lógica de decisão
		if opcao == "c" {
			topico.Tipo = "atuador";
			topico.Comando = IdSensor;
			// Publicar costuma ser uma operação rápida, 
			// mas se quiser concorrência total, mantenha o 'go'
			go publicarTopico(conn, topico);
			fmt.Println("comando enviado.");

		} else if opcao == "s" {
			topico.Tipo = "sensor";
			topico.TipoId = IdSensor ;
			
			for {
				assinarTopico(conn, topico)
				time.Sleep(time.Second);
			}
		}
	}
}
