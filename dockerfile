# --- Estágio 1: Build (Compilação) ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Instala dependências do sistema necessárias para compilar com CGo (SQLite precisa)
RUN apk add --no-cache gcc musl-dev

# Baixa dependências do Go
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário
# CGO_ENABLED=0 cria um binário estático puro (melhor pra containers)
# Mas como usamos SQLite "pure go" (modernc/glebarez), funciona liso.
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# --- Estágio 2: Run (Execução) ---
FROM alpine:latest

WORKDIR /app

# Copia o binário do estágio anterior
COPY --from=builder /app/main .

# Copia o arquivo .env (opcional, idealmente envs vêm do docker-compose)
COPY .env .

# Expõe a porta
EXPOSE 8080

# Comando para rodar
CMD ["./main"]
