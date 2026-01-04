# Guia do Swagger

Este projeto utiliza [swaggo/swag](https://github.com/swaggo/swag) para gerar automaticamente a documentação da API.

## Pré-requisitos

Certifique-se de ter o `swag` instalado. Você pode instalar usando o Makefile:

```bash
make install-tools
```

Ou manualmente:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Como Gerar a Documentação

Sempre que você alterar anotações do Swagger (os comentários que começam com `// @...`) em `main.go` ou nos handlers, você deve regenerar a documentação.

Execute o comando:

```bash
make swag
```

Isso irá:
1. Ler as anotações do código.
2. Gerar os arquivos JSON e YAML na pasta `cmd/api/swagger`.
3. Atualizar o `cmd/api/swagger/docs.go`.

## Acessando o Swagger UI

1. Inicie a aplicação:
   ```bash
   make run
   ```
2. Acesse no navegador:
   http://localhost:8080/swagger/index.html

## Principais Anotações Utilizadas

- **@Summary**: Resumo curto do endpoint.
- **@Description**: Descrição detalhada.
- **@Param**: Define parâmetros de rota, query ou body.
  - Ex: `// @Param page query int false "Página" default(1)`
- **@Success**: Define a resposta de sucesso.
- **@Failure**: Define respostas de erro.
- **@Router**: Define o caminho e o método HTTP.
  - Ex: `// @Router /products [get]`
