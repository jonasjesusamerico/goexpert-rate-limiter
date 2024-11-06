# Desafio Go: Rate Limiter

Este projeto implementa um sistema de Rate Limiting em Go, que controla a quantidade de requisições permitidas a um serviço em um período específico. O objetivo é gerenciar o tráfego, proteger o servidor contra sobrecarga e melhorar a segurança e estabilidade geral do sistema.

## Visão Geral

O Rate Limiter faz uso do Redis para armazenar dados sobre as requisições, garantindo alta performance e facilidade de escalabilidade. Ele suporta duas estratégias principais de limitação:

- **Limitação por IP**: Restringe o número de requisições por segundo para um endereço IP específico.
- **Limitação por Token**: Permite definir limites de requisições por segundo para tokens de acesso personalizados.

## Configuração e Execução

### Pré-requisitos

- **Docker**: Certifique-se de que o Docker está instalado no seu sistema.
- **Docker Compose**: É necessário para orquestrar a aplicação e o Redis.

### Configuração do Docker Compose

Crie um arquivo `docker-compose.yml` com o seguinte conteúdo:

```yaml
version: '3.8'

services:
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
    networks:
      - app-network

  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - redis
    environment:
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

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

- Use ferramentas como `curl` para enviar requisições em rápida sucessão e validar o comportamento do Rate Limiter.
- Exemplo de teste:
  ```sh
  curl -X GET http://localhost:8080/
  ```
- Verifique se as requisições são bloqueadas ao exceder o limite configurado.