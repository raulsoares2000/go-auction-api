# Sistema de Leilões em Go

## Visão Geral do Projeto

Este projeto consiste numa API RESTful para um sistema de leilões, desenvolvida integralmente em Go (Golang). A aplicação utiliza o framework Gin para a gestão de rotas e requisições HTTP e o MongoDB como banco de dados para a persistência de dados.

Uma funcionalidade central do sistema é o **encerramento automático dos leilões**. Após um leilão ser criado, uma rotina concorrente (goroutine) é iniciada para monitorizar o seu tempo de vida. Uma vez que o tempo predefinido expira, o status do leilão é automaticamente atualizado para "Concluído", impedindo novos lances.

## Funcionalidades Principais

- Criação de Leilões com tempo de duração configurável
- Encerramento automático de leilões expirados
- Criação de Lances (Bids) para leilões ativos
- Consulta de Leilões por status (Ativo, Concluído), categoria e nome do produto
- Consulta de todos os Lances de um determinado leilão
- Consulta do Lance vencedor de um leilão
- Consulta de Utilizadores por ID

## Tecnologias Utilizadas

- Go (Golang)
- Gin Web Framework
- MongoDB
- Docker
- Docker Compose

## Pré-requisitos

Para executar este projeto, certifique-se de que tem as seguintes ferramentas instaladas no seu ambiente de desenvolvimento:

- Docker Engine
- Docker Compose

## Como Executar o Projeto

Siga os passos abaixo para configurar e executar a aplicação no seu ambiente local.

### 1. Clone o Repositório

Navegue até ao diretório onde deseja guardar o projeto e clone o repositório.

### 2. Defina as Variáveis de Ambiente

Dentro do projeto, navegue até à pasta `cmd/auction` e modifique o arquivo chamado `.env`. Um exemplo de configuração seria:

```env
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=20s
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

**Nota sobre as variáveis:** A variável `AUCTION_INTERVAL` define o tempo de duração de um leilão. O valor `20s` (20 segundos) é ideal para testes rápidos.

### 3. Construa e Inicie os Contêineres

Volte para a raiz do projeto (o diretório onde se encontra o ficheiro `docker-compose.yml`) e execute o seguinte comando no seu terminal:

```bash
docker-compose up --build
```

Este comando irá construir a imagem Docker da aplicação Go, descarregar a imagem do MongoDB e iniciar ambos os serviços, conectando-os numa rede partilhada. Aguarde até que os logs indiquem que ambos os serviços estão a rodar e saudáveis.

## Exemplos de Utilização da API

Com a aplicação a rodar, pode interagir com os endpoints utilizando uma ferramenta como o `curl` ou o Postman.

### Criar um novo leilão

Abra um novo terminal e execute o comando abaixo:

```bash
curl --location 'http://localhost:8080/auction' \
--header 'Content-Type: application/json' \
--data '{
    "product_name": "Playstation 5",
    "category": "Eletronicos",
    "description": "Video game de ultima geracao, em perfeito estado.",
    "condition": 1
}'
```

### Procurar leilões concluídos (status=1)

Após esperar o tempo definido em `AUCTION_INTERVAL`, o leilão criado será movido para o estado "concluído":

```bash
curl --location 'http://localhost:8080/auction?status=1'
```

### Procurar leilões ativos (status=0)

```bash
curl --location 'http://localhost:8080/auction?status=0'
```

## Como Executar os Testes Automatizados

O projeto inclui testes automatizados para validar a lógica de negócio, incluindo o encerramento automático dos leilões.

### 1. Certifique-se de que o Ambiente Docker está em Execução

Os testes precisam de se conectar à instância do MongoDB que está a rorar no contêiner.

### 2. Execute o Comando de Teste

Na raiz do projeto, execute o seguinte comando no seu terminal:

```bash
go test -v ./...
```

Este comando irá procurar e executar todos os ficheiros de teste no projeto, mostrando um relatório detalhado da execução no terminal.


