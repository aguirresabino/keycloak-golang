# Keycloak

Exemplo de um cliente escrito em Golang.

## Iniciando servidor do Keycloak

```sh
cd keycloak

docker compose up -d

# caso sua versão do docker ainda não possua o compose integrado, executar o comando abaixo:
# docker-compose up -d
```
Quando o servidor ficar disponível em http://localhost:8080 já podemos importar o arquivo realm-export.json para criar um realm (Development) e um client (myclient) para testes. Em seguida, basta cadastrar os usuários no menu "User".

## Iniciando cliente Go

```sh
go run main.go
```

O cliente ficará disponível em http://localhost:8081.
