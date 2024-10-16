# API em Golang

Este projeto é uma API desenvolvida em Go. Abaixo estão as instruções para configurar e executar a API.

## Requisitos

- Go 1.18 ou superior
- Git

## Instalação

1. **Clone o repositório:**

    ```bash
    git clone https://github.com/IntegradorUFFS/back
    cd back
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

## Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).
