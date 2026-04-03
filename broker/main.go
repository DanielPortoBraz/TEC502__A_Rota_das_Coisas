package main

import (
	"fmt"
	"net"
	"sync"
	"log"
	"encoding/json"
	"time"
	"strings"
)

// ============== Broker ===========
/*

1. Inicia servidores TCP e UDP
2. Registra as listas de sensores, usuários e atuadores
3. Separa os tópicos
4. Roteia Sub para tópico Pub

* Fornece Logs de conexões

*/


// STRUCTS

/*
Padrão de Tópico: tipo/tipoId/comando. ("Valor" não será utilizado dentro do tópico, mas como dado fornecido por topicos. O mesmo vale para "Estado")
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
};

// Guarda os assinantes de cada tópico. Assume que não tem assinantes UDP
type Broker struct {
	assinantes map[string][]net.Conn; // Map de Assinantes (Guarda conns TCP para cada tópico)- [Topico]: [conn1, conn2, ...]
    ultimaAtividade map[string]time.Time // [Topico]: Timestamp - Utilizado para monitorar conexões UDP (Sensores)
	mu sync.RWMutex;
};

// Remove conexões TCP que falharam ou foram interrompidas
func (broker *Broker) removerConn(conn net.Conn) {

    // Percorre todos os tópicos 
    for topico, lista := range broker.assinantes {
        
        indice := -1;
        // Busca a conexão dentro do slice deste tópico
        for i, c := range lista {
            if c == conn {
                indice = i
                break
            }
        }

        // Se encontrou a conexão neste tópico, remove-a
        if indice != -1 {
            // Reajusta assinantes "pulando" a conexão encontrada na lista associada ao tópico
            broker.assinantes[topico] = append(lista[:indice], lista[indice+1:]...);
            
            // Se o tópico ficou sem ninguém, limpa a chave do mapa
            if len(broker.assinantes[topico]) == 0 {
                delete(broker.assinantes, topico);
            }
        }
    }
}



// FUNCOES - Pub/Sub

// Publica tópico no Broker
func (broker *Broker) publicar(topico Topico){
	
	// Apenas rota do tópico
	chave := fmt.Sprintf("%s/%s", topico.Tipo, topico.TipoId)

	// Recebe a lista de conns do Tópico a ser publicado
	broker.mu.Lock();

	if _, ok := broker.assinantes[chave]; !ok {
		broker.assinantes[chave] = []net.Conn{} // Cria uma lista de conns
	}
	
	conns := broker.assinantes[chave];
	broker.ultimaAtividade[chave] = time.Now(); // Útil para UDP (sensores) somente
	broker.mu.Unlock();

	// Serializa o topico 
	data, err := json.Marshal(topico);
	if err != nil {
		return;
	}

	// Publica Tópico para todos os conns inscritos
	for _, c := range conns {
		go func(conn net.Conn) {
			_, err := conn.Write(data);
			if err != nil {
				return;
			}
			conn.Write([]byte("\n"));
		}(c)
	}
}

// Assina tópico no Broker
func (broker *Broker) assinar(topico Topico, conn net.Conn){
	chave := fmt.Sprintf("%s/%s", topico.Tipo, topico.TipoId)

	// Adiciona um assinante (conn) ao Tópico
	broker.mu.Lock();

	// Verifica se o assinante irá assinar todos os tópicos de um determinado tipo (uso do WildCard "#")
	if topico.TipoId != "#"{
		broker.assinantes[chave] = append(broker.assinantes[chave], conn);
	} else {

		for key := range broker.assinantes {
			
			// Assina todos os tópicos daquele tipo
			if strings.HasPrefix(key, topico.Tipo+"/") {
				broker.assinantes[key] = append(broker.assinantes[key], conn)
			}
		}
	}
	broker.mu.Unlock();
}

// Monitora tópicos visando eliminar aqueles via UDP que não estão atualizando há mais de 10 segundos
func (broker *Broker) monitorarTopicos() {
    for {
        time.Sleep(5 * time.Second)
        broker.mu.Lock()

        for topico, lastPub := range broker.ultimaAtividade {

			// Verifica para excluir tópico: se o tópico é de sensor, não atualiza há +10s 
            if strings.HasPrefix(topico, "sensor") && time.Since(lastPub) > 10 * time.Second{
                fmt.Printf("[%s] (Broker) ALERTA: Tópico Removido por inatividade (+10s inativo)- %s\n", timeStamp(), topico);
                
                // Exclui tópico 
                delete(broker.assinantes, topico) 
				delete(broker.ultimaAtividade, topico); 

            }
        }
        broker.mu.Unlock()
    }
}

// FUNCOES GERAIS

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

// Inicia Servidor TCP
func StartServerTCP(broker *Broker) {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%s] (Broker) (TCP): Porta Aberta 9000\n", timeStamp());

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		fmt.Printf("[%s] (Broker) (TCP): Novo dispositivo conectado: %v\n", timeStamp(), conn.RemoteAddr())

		// Gerenciamento de Clientes TCP
		go handleConnectionTCP(conn, broker);
	}
}

// Inicia Servidor UDP
func StartServerUDP(broker *Broker) {

	address, err := net.ResolveUDPAddr("udp", ":9000");
	if err != nil {
		log.Fatal(err);
	}

	conn, err := net.ListenUDP("udp", address);
	if err != nil {
		log.Fatal(err);
	}

	fmt.Printf("[%s] (Broker) (UDP): Porta Aberta 9000\n", timeStamp());

	defer conn.Close();

	// Gerenciamento de Clientes UDP
	handleConnectionUDP(conn, broker);
}

func main() {

	// Canal de Tópicos
	broker := Broker{assinantes: make(map[string][]net.Conn), ultimaAtividade: make(map[string]time.Time)};

	// Iniciando Servidores TCP UDP - abre a porta, aceita e gerencia conexões
	go StartServerTCP(&broker);
	go StartServerUDP(&broker);
	
	for {		
		time.Sleep(time.Second);
	}
}