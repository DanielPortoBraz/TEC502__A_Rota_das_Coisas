package main

import (
	"fmt"
	"encoding/json"
	"net"
	"math/rand"
	"os"
)

type Sensor struct {
	ID string `json:"id"`
    Valor float64 `json:"valor"`
}

func newSensor() *Sensor{
	return &Sensor {
		ID : os.Getenv("HOSTNAME"), // Atribui id ao hostname do Conteiner Docker
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
	return 100 * rand.Float64(); // Valor entre 0 e 100
} 

func enviarDado(s *Sensor, conn net.Conn, value float64) error {
	// Atualiza Valor de Leitura do Sensor
	s.setValor(value); 

	// Serialização do JSON
    data, err := json.Marshal(s)
    if err != nil {
        return err
    }

	// Envio de dado para a Rede
    _, err = conn.Write(data)

	fmt.Println("Device: Sensor", s);
    return err
}