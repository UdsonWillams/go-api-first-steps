# 🚀 Go API - Product Service

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)
![Azure](https://img.shields.io/badge/Azure-App%20Insights-0078D4?style=for-the-badge&logo=microsoftazure&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

> **Microsserviço de Alta Performance para Gerenciamento de Produtos**

Este projeto é um exemplo robusto de uma API REST desenvolvida em **Go (Golang)**, projetada para ambientes **Cloud Native**. Ele implementa as melhores práticas de mercado, incluindo arquitetura limpa, observabilidade distribuída e segurança JWT avançada.

---

## ⚡ Por que este projeto é incrível?

- **🚀 Performance Extrema:** Compilado para código de máquina, sem VM, com gerenciamento de memória eficiente.
- **🏗 Clean Architecture:** Código desacoplado, testável e fácil de manter.
- **🔐 Segurança Enterprise:** Autenticação via **Keycloak** (JWT RS256) com controle de acesso baseado em cargos (RBAC).
- **🐳 Container Native:** Imagens Docker **Alpine** otimizadas (< 20MB) prontas para Kubernetes.
- **🔍 Observabilidade Híbrida:** Integração nativa com **Azure Application Insights** e Logs JSON estruturados.

---

## 🛠 Stack Tecnológica

| Tech        | Função        | Descrição                             |
| :---------- | :------------ | :------------------------------------ |
| **Go**      | Linguagem     | Versão 1.22+                          |
| **Gin**     | Framework Web | Alta performance e middleware robusto |
| **GORM**    | ORM           | Manipulação de dados e AutoMigrate    |
| **Slog**    | Logging       | Logs estruturados com `trace_id`      |
| **Swagger** | Docs          | Documentação automática via código    |
| **Docker**  | Container     | Multi-stage build                     |

---

## 📂 Estrutura do Projeto

O projeto segue o **Standard Go Project Layout**:

```bash
.
├── cmd/
│   ├── api/            # 🏁 Entrypoint (Main)
│   └── mock_token/     # 🛠 Gerador de Tokens (Dev Tools)
├── internal/           # 🔒 Código Privado (Core Business)
│   ├── handlers/       # 🎮 Controladores HTTP
│   ├── middleware/     # 🚦 Autenticação, Logs, CORS
│   ├── product/        # 📦 Regras de Negócio (Service)
│   └── storage/        # 💾 Camada de Dados (Repository)
├── pkg/                # 📦 Bibliotecas Compartilhadas
│   └── logger/         # 📝 Configuração avançada de Logs (Fanout/Azure)
├── docs/               # 📄 Arquivos OpenAPI/Swagger
├── docker-compose.yml  # 🐳 Orquestração Local
└── .env                # 🔑 Variáveis de Ambiente
```

---

## 🚀 Como Rodar

### 1️⃣ Pré-requisitos

- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) (Opcional)

### 2️⃣ Configuração

Crie um arquivo `.env` na raiz:

```env
PORT=:8080
DB_URL=meubanco.db
# Cole a chave pública gerada pelo passo 3 abaixo:
KEYCLOAK_PUBLIC_KEY=...
# (Opcional) Connection String do Azure App Insights
APPINSIGHTS_CONNECTION_STRING=...
```

### 3️⃣ Gerando Acessos (Mock)

Como não temos um Keycloak rodando, use nossa ferramenta interna para gerar credenciais:

```bash
go run cmd/mock_token/main.go
```

- ✅ Copie a **Public Key** para o `.env`.
- ✅ Copie o **Token Bearer** para usar nas requisições.

### 4️⃣ Executando

**Modo Dev (Local):**

```bash
go mod tidy
go run cmd/api/main.go
```

**Modo Produção (Docker):**

```bash
docker-compose up --build
```

---

## 📖 Documentação Interativa

Acesse o Swagger UI para testar os endpoints visualmente:

👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

1.  Clique em **Authorize** 🔓
2.  Cole seu Token Bearer
3.  Teste os endpoints (`GET`, `POST`, `PUT`, `DELETE`)

---

## ☁️ Observabilidade (Azure)

O sistema possui um **Logger Híbrido (Fanout)**. Se configurado, ele envia logs para:

1.  **Console (Stdout):** Em formato JSON para o Docker/K8s.
2.  **Azure App Insights:** Envio assíncrono via SDK.

> **Dica:** O campo `Operation Id` no Azure é sincronizado com o `trace_id` dos logs da aplicação.

---

## 👨‍💻 Autor

<table align="center">
    <tr>
        <td align="center">
            <a href="https://github.com/udsonwillams">
                <img src="https://github.com/udsonwillams.png" width="100px;" alt="Foto do Udson Willams"/>
                <br />
                <sub><b>Udson Willams</b></sub>
            </a>
        </td>
    </tr>
</table>

<p align="center">
  Feito com 💜 e Go 🐹
</p>
