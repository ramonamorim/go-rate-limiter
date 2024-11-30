# go-rate-limite: Desafio Rate Limiter

Este projeto desenvolve um sistema de Rate Limiting em Go para limitar a quantidade de requisições que um serviço pode receber em um determinado período. A ideia é ajudar a controlar o tráfego, evitar sobrecarga no servidor e garantir mais segurança e estabilidade para o sistema.

## Visão Geral

O Redis é usado para guardar os dados das requisições, oferecendo alto desempenho e fácil escalabilidade. Ele permite trabalhar com duas estratégias principais de rate limiting:

- **Limitação por IP**: Restringe o número de requisições por segundo para um endereço IP específico.
- **Limitação por Token**: Permite definir limites de requisições por segundo para tokens de acesso personalizados.

### Pré-requisitos

- **Docker**: Certifique-se de que o Docker está instalado no seu sistema.
- **Docker Compose**: É necessário para orquestrar a aplicação e o Redis.


### Estrutura :

O projeto está estruturado da seguinte forma:

- **Middleware**: Responsável por interceptar requisições e aplicar as regras de Rate Limiting.
- **Serviço de Rate Limiter**: Implementa a lógica de controle de requisições.
- **Armazenamento Redis**: Utilizado para manter contagens e controlar o tempo de expiração.


### Iniciando o Projeto

1. Na raiz do projeto, execute o comando:
   ```sh
   docker-compose up --build
   ```
2. O serviço estará disponível em: `http://localhost:8080`.

## Como Usar

Para testar o Rate Limiter, utilize um cliente HTTP como `curl` ou ferramentas como Postman:

- **Requisições com Limitação por IP**: Todas as requisições feitas do mesmo endereço IP serão monitoradas e limitadas.
- **Requisições com Token de Acesso**: Inclua um cabeçalho `API_KEY` na requisição, por exemplo:
  ```sh
  curl -H "API_KEY: seu_token" http://localhost:8080/
  ```

## Estrutura do Projeto

O projeto está estruturado da seguinte forma:

- **Middleware**: Responsável por interceptar requisições e aplicar as regras de Rate Limiting.
- **Serviço de Rate Limiter**: Implementa a lógica de controle de requisições.
- **Armazenamento Redis**: Utilizado para manter contagens e controlar o tempo de expiração.

## Testes

- Use o `curl` para enviar requisições e validar o comportamento do Rate Limiter.
- Exemplo de teste:
  ```sh
  curl -X GET http://localhost:8080/
  ```
- Verifique se as requisições são bloqueadas ao exceder o limite configurado.
- Testes:
```sh
go test ./internal/application
go test ./internal/domain/service 
go test ./internal/infra/config 
go test ./internal/infra/limiter
go test ./internal/infra/middleware
```