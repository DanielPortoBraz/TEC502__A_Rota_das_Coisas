# TEC502__A_Rota_das_Coisas
Projeto avaliativo do Componente Curricular TEC502: MI - Concorrência e Conectividade.

## Definição do Projeto

Este projeto propõe o desenvolvimento de uma **rede IoT (Internet of Things)** composta por sensores, atuadores e uma interface de interação com o usuário.

A arquitetura é baseada no princípio de **desacoplamento entre dispositivos**, sendo mediada por um **Broker central** responsável por:

- Receber dados enviados pelos sensores  
- Processar e distribuir essas informações  
- Receber comandos de usuários  
- Gerenciar o acionamento dos atuadores  

Essa abordagem proporciona maior **escalabilidade, modularidade e flexibilidade** ao sistema.

### Protocolos de Comunicação

O sistema utiliza diferentes protocolos conforme o tipo de comunicação:

- **UDP (User Datagram Protocol)**  
  Utilizado para envio de dados dos sensores ao Broker, priorizando baixa latência e menor overhead.

- **TCP/IP (Transmission Control Protocol / Internet Protocol)**  
  Utilizado na comunicação entre usuários, atuadores e o Broker, garantindo confiabilidade e integridade na transmissão dos dados.

---

## Arquitetura do Sistema

A arquitetura segue um modelo inspirado no paradigma **Publish/Subscribe (Pub/Sub)**, no qual:

- Sensores atuam como **publicadores de dados**
- Usuários enviam **comandos e requisições**
- Atuadores executam **ações com base nas decisões do sistema**
- O Broker atua como **núcleo central de comunicação e controle**

Essa organização permite que os componentes operem de forma independente, facilitando manutenção e expansão do sistema.

### Diagrama (Esboço Inicial)

O diagrama abaixo representa uma visão inicial da arquitetura do sistema:

<img width="640" height="380" alt="Diagrama em branco" src="https://github.com/user-attachments/assets/7d9af048-76f1-435e-810b-3f656ec4a65e" />

  > **Observação:** Este diagrama é um rascunho inicial que não segue qualquer padrão.
