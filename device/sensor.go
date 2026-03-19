package main

import (
	"fmt"
	"encoding/json"
	"net"
	"math/rand"
)

type Sensor struct {
	ID string `json:"id"`
    Valor float64 `json:"valor"`
}

func newSensor() *Sensor{
	return &Sensor {
		ID : "id123",
		Valor : 0.0,
	}	
}


func (s *Sensor) getId() string {
	return s.ID;
}

func (s *Sensor) getValor() float64 {
	return s.Valor;
}

func (s *Sensor) setValor(valor float64) {
	s.Valor = valor;
}

func lerDado() float64 {
	return 100 * rand.Float64(); // Valor entre 0 e 100 em 
} 

func enviarDado(conn net.Conn, value float64) error {
    sensor := newSensor();

	// Atualiza Valor de Leitura do Sensor
	sensor.setValor(value); 

	// Serialização do JSON
    data, err := json.Marshal(sensor)
    if err != nil {
        return err
    }

	// Envio de dado para a Rede
    _, err = conn.Write(data)

	fmt.Println("Dado:", sensor);
    return err
}