# Tech Stack

Tecnologias e bibliotecas escolhidas para este projeto.

## Core
- **Go 1.22+**: Linguagem moderna, tipada e compilada.

## Web Framework
- **Gin Web Framework** (`github.com/gin-gonic/gin`)
  - Motivo: Alta performance (baseado em HttpRouter), fácil de usar e ecossistema gigante.

## Database (ORM)
- **GORM** (`gorm.io/gorm`)
- **SQLite** (`github.com/glebarez/sqlite`)
  - O GORM facilita o CRUD e Migrations automáticas. O SQLite foi escolhido pela simplicidade (arquivo local), mas o GORM permite trocar para Postgres/MySQL mudando apenas 1 linha.

## Autenticação & Segurança
- **Keycloak (OIDC)**
- **go-oidc** (`github.com/coreos/go-oidc/v3`)
  - Implementamos Autenticação robusta usando OpenID Connect. A API não gerencia usuários/senhas, ela apenas valida tokens assinados pelo Keycloak, garantindo segurança de nível empresarial.

## Observabilidade
- **log/slog** (Tech Nativa do Go 1.21+)
  - Logging estruturado (JSON) nativo.
- **Azure Application Insights** (Opcional)
  - Preparado para integração via env var.

## Documentação
- **Swagger / OpenAPI** (`github.com/swaggo/swag`)
  - Gera documentação interativa automaticamente a partir de anotações no código.
