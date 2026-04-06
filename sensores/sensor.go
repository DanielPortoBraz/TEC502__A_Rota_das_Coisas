package main

import (
	"fmt"
	"encoding/json"
	"net"
	"math/rand"
)


func newSensor() *Topico{
	return &Topico {
		Acao : "pub",
		Tipo : "sensor",
		TipoId : fmt.Sprintf("%d", rand.Intn(100)),
		Comando : "",
		Valor : 0.0,
		Estado : false,
	}	
}


func lerDado() float64 {
	return 100 * rand.Float64(); // Valor entre 0 e 100
} 

func enviarDado(conn net.Conn, sensor *Topico) error {

	msg := sensor; // Sensor em formato de tópico

	// Serialização para JSON
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)

	//fmt.Printf("[%s] (Sensor):\nID: %s\nValor: %.2f\n", timeStamp(), sensor.TipoId, sensor.Valor);

	return err
}