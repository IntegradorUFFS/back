# API em Golang

Este projeto é uma API desenvolvida em Go. Abaixo estão as instruções para configurar e executar a API.

## Requisitos

- Go 1.18 ou superior
- Git

## Instalação

1. **Clone o repositório:**

   ```bash
   git clone https://github.com/ImaKrp/golang-api.git
   cd golang-api
   ```

2. **Instale as dependências:**

   Certifique-se de ter o Go instalado e configure seu ambiente de desenvolvimento. Em seguida, instale as dependências do projeto usando o comando:

   ```bash
   go mod tidy
   ```

## Gerar Código

O projeto pode usar geração automática de código. Para gerar o código necessário, execute:

    go generate ./...

Este comando irá buscar por todos os arquivos `//go:generate` no projeto e executar os comandos especificados neles.

## Executar a API

Para iniciar a API, execute o comando:

    go run cmd/router/main.go

Isso irá iniciar o servidor e expor a API na porta configurada.

## Configurando o Backend (API em Go)

1. Acesse o diretório backend:

```bash
cd backend
```

2. Instale as dependências:

```bash
go mod tidy
```

2. Rode as migrações:

```bash
go install github.com/jackc/tern/v2@latest

go generate ./...
```

3. Configure o PostgreSQL e defina as variáveis de ambiente no arquivo `.env`:

```bash
DATABASE_PORT =
DATABASE_HOST =
DATABASE_USER =
DATABASE_PASSWORD =
DATABASE =
JWT_SECRET =
```

4. Execute a API:

```bash
 go run cmd/router/main.go
```

A API estará disponível em `http://localhost:8080`.

## Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).
