# 📘 Guia de Estudos e Decisões do Projeto Go

Este documento resume o que aprendemos e validamos na criação da estrutura inicial da API em Go. O foco é entender **por que** as coisas são feitas de determinada maneira e os comandos essenciais.

---

## 1. Comandos Essenciais do Terminal

Aqui está o "dicionário" básico para operar o Go.

### `go mod init <nome-do-projeto>`

- **O que faz:** Cria a "identidade" do projeto. Gera o arquivo `go.mod`.
- **Para que serve:** É como o `requirements.txt` do Python ou `package.json` do Node. Ele diz ao Go: "Eu sou um módulo e minhas dependências são essas".

### `go run cmd/api/main.go`

- **O que faz:** Compila o código na memória RAM e executa imediatamente.
- **Para que serve:** Usado durante o desenvolvimento para testar rápido. Não gera arquivo executável no disco.

### `go build -o meu-app.exe cmd/api/main.go`

- **O que faz:** Compila o código e **gera um binário** (um arquivo `.exe` no Windows ou executável no Linux/Mac).
- **Para que serve:** É o que você manda para o servidor de produção. Esse arquivo não precisa do Go instalado para rodar.

### `go get <link-do-pacote>`

- **O que faz:** Baixa uma biblioteca da internet.
- **Para que serve:** Instalar drivers de banco, frameworks, etc.

### `go mod tidy`

- **O que faz:** A "faxina".
- **Para que serve:** Remove dependências que você não usa mais e baixa as que faltam no `go.mod`.

---

## 2. Regras da Linguagem (A "Pegadinha" da Visibilidade)

Em Go, não existem palavras como `public` ou `private`. A visibilidade é definida pela **primeira letra** do nome da função, struct ou variável.

| Sintaxe                | Visibilidade             | Explicação                                                                    |
| :--------------------- | :----------------------- | :---------------------------------------------------------------------------- |
| **`func HealthCheck`** | **Pública** (Exported)   | Outros pacotes (pastas) conseguem ver e importar essa função.                 |
| **`func healthCheck`** | **Privada** (Unexported) | Só funciona dentro da pasta onde foi criada. Ninguém de fora consegue chamar. |

> **Regra de Ouro:** Quer usar em outro arquivo que está em outra pasta? **Comece com Letra Maiúscula.**

---

## 3. Estrutura de Pastas (Padrão de Mercado)

### `cmd/`

- É a porta de entrada.
- Cada pasta aqui dentro vira um executável diferente.
- **Exemplo:** `cmd/api/main.go` inicia a API. `cmd/worker/main.go` iniciaria um processo de background.

### `internal/`

- É a área "VIP" e protegida do seu código.
- **Regra do Go:** O Go proíbe que projetos externos importem qualquer coisa que esteja dentro de uma pasta chamada `internal`.
- Serve para garantir que a lógica do seu negócio (`internal/handlers`, `internal/storage`) seja usada apenas pelo seu próprio projeto.

---

## 4. O Problema do SQLite (CGo vs Pure Go)

Validamos a conexão com banco de dados SQLite e encontramos um erro comum.

- **O Erro:** `Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo`.
- **A Causa:** O driver mais famoso (`github.com/mattn/go-sqlite3`) usa código em **C** por baixo dos panos. Para funcionar, exige que você tenha um compilador C (GCC) instalado e configurado no Windows, o que é complexo.
- **A Solução:** Trocamos para um driver **"Pure Go"** (`modernc.org/sqlite`).
- **Vantagem:** Esse driver foi reescrito 100% em Go. Ele compila em qualquer máquina sem precisar instalar nada extra.

---

## 5. Fluxo de Dados (Arquitetura Simples)

Criamos 3 camadas que se comunicam via injeção de dependência (passando um objeto para dentro do outro):

1.  **Main (`cmd/api`)**:

    - O chefe. Ele inicia tudo.
    - Cria o banco (`Repo`), cria o serviço (`Service`) e conecta os dois.

2.  **Service (`internal/product`)**:

    - As regras de negócio.
    - Ele recebe o `Repo` e diz: "Validei os dados, agora salva aí".

3.  **Repository (`internal/storage`)**:
    - O operário do banco.
    - Ele só sabe falar SQL (`INSERT`, `SELECT`). Não sabe regra de negócio.

---

# 📘 Guia de Estudos: API Go com GORM e Testes

Este documento resume a evolução do projeto, saindo de um código básico para uma API profissional com CRUD completo, banco de dados gerenciado via ORM e testes automatizados.

---

## 1. Novos Comandos Essenciais

Além dos comandos básicos, agora usamos estes:

### `go get gorm.io/gorm`

- **O que faz:** Baixa a biblioteca do GORM (nosso ORM).

### `go test ./...`

- **O que faz:** O comando mágico de testes.
- **O detalhe:** O `./...` diz ao Go: "Rode os testes desta pasta **e de todas as subpastas** recursivamente".
- **Saída:** Mostra `ok` (passou) ou `FAIL` (quebrou) para cada pacote.

---

## 2. O Que é ORM (GORM)?

Antes, escrevíamos SQL manual (`INSERT INTO...`). Agora usamos ORM (_Object-Relational Mapping_).

- **Conceito:** O ORM mapeia suas **Structs** (classes do Go) para **Tabelas** do banco.
- **AutoMigrate:** O GORM olha para sua struct `Product` e cria a tabela automaticamente. Se você adicionar um campo novo no código, ele atualiza o banco sozinho.
- **Model:** Ao colocar `gorm.Model` dentro da sua struct, você ganha de graça os campos:
  - `ID` (Chave primária)
  - `CreatedAt` (Data de criação)
  - `UpdatedAt` (Data de atualização)
  - `DeletedAt` (Soft Delete - o dado não é apagado, apenas escondido).

---

## 3. A Solução do Driver (CGo vs Pure Go)

Tivemos problemas de compilação no Windows porque o driver padrão do SQLite exige um compilador C (GCC).

- **Solução:** Usamos o driver `github.com/glebarez/sqlite`.
- **Por que?** Ele é **"Pure Go"**. Foi reescrito do zero usando apenas Go, eliminando a necessidade de instalar ferramentas externas no Windows.

---

## 4. Hierarquia de Pastas (Evitando o "Import Cycle")

Aprendemos que o Go é rigoroso com dependências circulares. A regra é: **As setas de importação só apontam para baixo.**

**A Ordem Correta:**

1.  🟦 **Main** (`cmd/api`) → _Importa Handlers, Service e Repo_
2.  ⬇️
3.  🟩 **Handlers** (`internal/handlers`) → _Importa Service_
4.  ⬇️
5.  🟨 **Service** (`internal/product`) → _Importa Storage_
6.  ⬇️
7.  🟥 **Storage** (`internal/storage`) → _Não importa ninguém do projeto_

> **Erro Comum:** Se o `Service` tentar importar o `Handler`, o Go trava, pois cria um loop infinito (A chama B, que chama A).

---

## 5. Estrutura do CRUD

Implementamos as 4 operações básicas mapeadas para verbos HTTP:

| Verbo HTTP | Rota             | Função no Código | Ação                      |
| :--------- | :--------------- | :--------------- | :------------------------ |
| **POST**   | `/products`      | `Create`         | Cria novo item.           |
| **GET**    | `/products`      | `List`           | Busca todos os itens.     |
| **PUT**    | `/products/{id}` | `Update`         | Altera um item existente. |
| **DELETE** | `/products/{id}` | `Delete`         | Remove um item.           |

_Dica: No Go 1.22+, usamos `r.PathValue("id")` para pegar o ID direto da URL, sem precisar de bibliotecas externas de roteamento._

---

## 6. Testes Automatizados

Criamos arquivos com o final `_test.go` (ex: `service_test.go`).

### A Estratégia do Banco em Memória

Para testar, não queremos sujar o arquivo real `meubanco.db`.

- **Truque:** Passamos a string `":memory:"` para o GORM.
- **Resultado:** O SQLite cria um banco inteiramente na memória RAM. Ele é super rápido, isolado e desaparece assim que o teste acaba.

### Exemplo de Teste

```go
func TestCreateProduct(t *testing.T) {
    // 1. Arrange (Prepara): Banco fake na memória
    repo := storage.NewRepository(":memory:")
    service := Service{Repo: repo}

    // 2. Act (Age): Tenta criar
    nome, _ := service.CreateProduct("Teclado")

    // 3. Assert (Valida): Verifica se deu certo
    if nome != "Teclado" {
        t.Errorf("Esperava Teclado, veio %s", nome)
    }
}
```

---

# 🚀 API Go (Gin + GORM + SQLite)

Projeto inicial de uma API RESTful robusta usando as melhores práticas do ecossistema Go.

## 🛠 Tecnologias

- **Linguagem:** Go (Golang) 1.22+
- **Framework Web:** [Gin Web Framework](https://github.com/gin-gonic/gin) (Alta performance e produtividade)
- **ORM:** [GORM](https://gorm.io/) (Manipulação de banco de dados)
- **Database:** SQLite (Driver Pure Go - sem dependência de CGo)
- **Config:** Godotenv (Variáveis de ambiente)
- **Testes:** Go Testing + Banco em memória

## 📂 Estrutura (Clean Architecture Simplificada)

O projeto segue o padrão `Standard Go Project Layout`:

- `cmd/api`: Ponto de entrada (Main).
- `internal/handlers`: Camada HTTP (Gin Controllers).
- `internal/product`: Regra de Negócio (Service).
- `internal/storage`: Acesso a Dados (Repository/SQL).

## ⚡ Como Rodar

### Pré-requisitos

- Go instalado

### Passo a Passo

1.  **Clone o repo:**

    ```bash
    git clone [https://github.com/UdsonWillams/go-api-first-steps.git](https://github.com/UdsonWillams/go-api-first-steps.git)
    cd go-api-first-steps
    ```

2.  **Instale as dependências:**

    ```bash
    go mod tidy
    ```

3.  **Configure o ambiente:**
    Crie um arquivo `.env` na raiz:

    ```env
    PORT=:8080
    DB_URL=meubanco.db
    ```

4.  **Execute:**
    ```bash
    go run cmd/api/main.go
    ```

## 🧪 Testes

Para rodar os testes unitários (que usam banco em memória):

```bash
go test ./...
```
