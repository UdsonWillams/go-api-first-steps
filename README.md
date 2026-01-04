# Go API Project (First Steps)

![Go Gopher](https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher-front.png)

API REST moderna escrita em Go, focada em boas práticas, arquitetura limpa e alta performance.

## 📚 Documentação
Toda a documentação detalhada foi movida para a pasta `docs/`.

- **[Arquitetura](docs/ARCHITECTURE.md)**: Entenda a estrutura de pastas (Clean Architecture).
- **[Stack Tecnológico](docs/STACK.md)**: Gin, GORM, SQLite, Slog, Keycloak.
- **[Guia: De Python para Go](docs/GOLANG_BASICS.md)**: Se você vem do Python, comece por aqui.
- **[Concorrência: Goroutines](docs/CONCURRENCY.md)**: O superpoder do Go explicado.
- **[Testes e Benchmarks](docs/TESTING.md)**: Como garantir qualidade e medir nanosegundos.
- **[Swagger Guide](docs/SWAGGER_GUIDE.md)**: Como gerar a documentação da API.
- **[Fun Facts 🐹](docs/FUN_FACTS.md)**: Curiosidades sobre a linguagem e os criadores.

---

## 🚀 Como Rodar

### Pré-requisitos
- Go 1.22+
- Make (Opcional, mas recomendado)

### Comandos Rápidos

```bash
# Rodar a aplicação
make run

# Rodar Testes
make test

# Gerar Documentação Swagger
make swag

# Verificar Linters e Qualidade
make lint
```

## 🔐 Autenticação (OIDC / Keycloak)
Este projeto usa **OpenID Connect**.
Para rodar localmente, configure o `.env` (use `.env.example` como base) apontando para sem Keycloak.

```env
KEYCLOAK_URL=http://localhost:8080/realms/meurealm
KEYCLOAK_CLIENT_ID=meu-client
```

## 🛠 Features Implementadas
- [x] API Versioning (`/api/v1`)
- [x] Paginação de Resultados
- [x] Autenticação Stateless com JWKS (Singleton)
- [x] Validação de Roles (AND/OR Logic)
- [x] Logging Estruturado (JSON)
- [x] Graceful Shutdown
