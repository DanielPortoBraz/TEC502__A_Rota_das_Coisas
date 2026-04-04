package main

import (
	"fmt"
	"net"
	"time"
)

// IP do Computador Servidor
const IP_SERVER = "172.16.201.5"

// ============ Sensor =============
/*

1. Conecta ao Broker via UDP
2. Cria Sensor: Fornece sensorID ao Broker
3. Publica leitura em tópico: sensor/ *sensorID* /- (Valor dentro da struct tópico) 

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


func main() {

	// Conexão por UDP
    conn, err := net.Dial("udp", fmt.Sprintf("%s:9000", IP_SERVER));
	if err != nil {
		panic(err);
	}
	
	fmt.Println("Device: UDP: Porta Aberta")

	// Sensor Criado (Tópico)
	sensor := newSensor()

	for {
		sensor.Valor = lerDado();

		enviarDado(conn, sensor);

		time.Sleep(100 * time.Millisecond); // Envia leitura contínua a cada 5 segundos
	}
}