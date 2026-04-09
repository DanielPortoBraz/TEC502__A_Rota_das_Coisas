# A Rota das Coisas 

## 📌 Visão Geral

Este projeto implementa um **middleware de integração para redes IoT** baseado no paradigma **Publish/Subscribe**, desenvolvido em **Golang** como parte de uma abordagem **PBL (Problem-Based Learning)**. O projeto segue como componente avaliativo da disciplina **TEC 502 MI - Concorrência e Conectividade** no curso de Engenharia de Computação na Universidade Estadual de Feira de Santana (UEFS).

A solução resolve problemas clássicos de **alto acoplamento** em arquiteturas ponto-a-ponto, utilizando um **Broker central** responsável por intermediar a comunicação entre sensores, atuadores e usuários.

---
## Estrutura do Projeto

```bash

├── broker/
│   ├── main.go
│   ├── tcp_server.go
│   ├── udp_server.go
│   └── Dockerfile
├── sensor/
│   ├── main.go
│   ├── sensor.go
│   └── Dockerfile
├── atuador/
│   ├── main.go
│   ├── atuador.go
│   └── Dockerfile
├── client/
│   ├── main.go
│   ├── user.go
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```
---

## 🧠 Conceitos-Chave

- Arquitetura **Publish/Subscribe (inspirada no MQTT)**
- Separação de tráfego:
  - **UDP** → Telemetria (sensores)
  - **TCP** → Comandos críticos (atuadores e usuários)
- **Concorrência com goroutines**
- **Comunicação via sockets**
- **Containerização com Docker Compose**

---

## 🏗️ Arquitetura do Sistema

O sistema segue uma **topologia em estrela**, em que conforme o diagrama abaixo, têm-se as seguintes entidades:


<img width="600" height="560" alt="PBL1_GO - Arquitetura Geral" src="https://github.com/user-attachments/assets/4b598279-8926-4555-8a2a-0473a790a2dd" />

- **Broker**
  - Núcleo do sistema
  - Gerencia tópicos e assinantes
  - Possui servidores TCP e UDP

- **Sensores**
  - Publicam dados via UDP
  - Alta frequência (~100ms)
  - Stateless

- **Atuadores**
  - Recebem comandos via TCP
  - Implementam heartbeat e reconexão

- **Usuário**
  - Interface interativa via terminal
  - Permite:
    - Assinar sensores
    - Controlar atuadores
    - Testar concorrência

---

## 📡 Estrutura de Tópicos

Formato padrão:

```
tipo/tipoId/comando
```

Exemplos:
```
sensor/10/-
atuador/15/on
```

Estrutura interna (JSON):

```json
{
  "acao": "pub|sub|unsub|ping",
  "tipo": "sensor|atuador|usuario",
  "tipoId": "id",
  "comando": "true|false",
  "valor": 0.0,
  "estado": false
}
```

---

## ⚙️ Funcionalidades

### ✔️ Broker
- Gerenciamento de tópicos com `map + mutex`
- Suporte a wildcard (`#`)
- Heartbeat (ping/pong)
- Remoção de sensores inativos
- Alta escalabilidade

### ✔️ Usuário (Interface CLI)
- Publicar comandos em atuadores
- Assinar sensores (inclusive com `#`)
- Atualização de tela a cada 1s
- Teste de concorrência
- Reconexão automática

### ✔️ Atuadores
- Execução confiável via TCP
- Monitoramento de conexão
- Reconexão automática

### ✔️ Sensores
- Publicação contínua
- Baixo overhead (UDP)

---

## 🖥️ Interface do Usuário

Menu interativo:

```
[c] Publicar comando
[s] Assinar sensor
[t] Teste de concorrência
```

### Exemplo de uso
```bash
----- Painel de Controle -----
[c] Publicar comando
[s] Assinar sensor
[t] Teste de Concorrência p/ Atuadores
Digite o comando: c
Digite o ID do dispositivo ou "#" para Visualizar Disponíveis: #

[9-4-2026 17:57:8](Usuario): Topico Listado - atuador/57d35f72496f
[9-4-2026 17:57:8](Usuario): Topico Listado - atuador/99073dd03805
[9-4-2026 17:57:9](Usuario): Topico Listado - atuador/57d35f72496f
[9-4-2026 17:57:9](Usuario): Topico Listado - atuador/99073dd03805

p
(Usuario) Encerrando assinatura...

----- Painel de Controle -----
[c] Publicar comando
[s] Assinar sensor
[t] Teste de Concorrência p/ Atuadores
Digite o comando: s
Digite o ID do dispositivo: cc38a6224bcc

[9-4-2026 17:57:34](Usuario): Topico Recebido - sensor/cc38a6224bcc/
Estado: false
Valor: 71.92

[9-4-2026 17:57:35](Usuario): Topico Recebido - sensor/cc38a6224bcc/
Estado: false
Valor: 2.32

[9-4-2026 17:57:36](Usuario): Topico Recebido - sensor/cc38a6224bcc/
Estado: false
Valor: 99.13

p
(Usuario) Encerrando assinatura...

----- Painel de Controle -----
[c] Publicar comando
[s] Assinar sensor
[t] Teste de Concorrência p/ Atuadores
Digite o comando: c
Digite o ID do dispositivo ou "#" para Visualizar Disponíveis: 99073dd03805
Digite o Comando (on/off): on

Comando enviado.
```

---

## 🐳 Execução com Docker

### 1. Configurar IP do Broker

Definir variável de ambiente:

```
BROKER_HOST=<IP_DA_MAQUINA>
```

### 2. Subir containers

```
docker compose up --build
```

---

## 🌐 Execução em Múltiplas Máquinas

1. Rodar o **Broker** em uma máquina
2. Descobrir o IP da máquina
3. Configurar `BROKER_HOST` nos clientes
4. Rodar sensores, usuários e atuadores em outras máquinas

---

## 📊 Resultados de Testes

> Os testes foram feitos com sensores lendo a cada 100ms
- **301 containers simultâneos**
- CPU: **2% ~ 5% (pico ~7%)**
- Memória: **~13.48 MB**
- Tráfego: **36.4 MB**
- Sistema estável (sem deadlocks)

---

## 🔄 Mecanismos Importantes

### 🔁 Heartbeat
- Ping a cada 5s
- Timeout de 10s
- Detecta falhas de conexão

### 🔌 Reconexão
- Tentativa a cada 3 segundos
- Automática

### 🧵 Concorrência
- Goroutines para:
  - Dispatcher
  - Heartbeat
  - Monitoramento
  - Terminal

---

## 🚧 Melhorias Futuras

- Implementação de **QoS**
- Suporte a **TLS**
- Persistência de mensagens
- Autenticação de dispositivos

---

## 👨‍💻 Autor

**Daniel Porto Braz**  
Engenharia de Computação - UEFS

---

## 📄 Licença

Projeto acadêmico desenvolvido para fins educacionais.
