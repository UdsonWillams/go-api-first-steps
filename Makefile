# Nome do binário final
APP_NAME=product-api

# Atalho para o arquivo principal
MAIN_FILE=cmd/api/main.go

# --- COMANDOS PRINCIPAIS ---

# Padrão: se digitar só 'make', ele roda a aplicação
all: run

# 🚀 Roda a aplicação (Hot Reload se usar air, ou go run normal)
run:
	@echo "🔥 Rodando a API..."
	go run $(MAIN_FILE)

# 🛠 Instala as ferramentas necessárias (O pedido principal!)
install-tools:
	@echo "📦 Instalando Linter e Swag..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# 🧹 Roda o Linter (Verifica erros e estilo)
lint:
	@echo "🔍 Verificando código..."
	golangci-lint run

# 📄 Atualiza a documentação do Swagger
swag:
	@echo "📄 Gerando Swagger..."
	swag init -g $(MAIN_FILE) --output docs

# 🧪 Roda os testes
test:
	@echo "🧪 Rodando testes..."
	go test -v ./...

# 🔑 Gera o Token Mock (atalho pro script que criamos)
mock:
	@echo "🔑 Gerando Token de Teste..."
	go run cmd/mock_token/main.go

# 🏗 Builda o binário para produção
build:
	@echo "🏗 Compilando..."
	go build -o bin/$(APP_NAME).exe $(MAIN_FILE)

# 🐳 Roda tudo no Docker
docker-up:
	docker-compose up --build

# 🧹 Limpa dependências não usadas
tidy:
	go mod tidy

# ⚡ Atalho Full: Formata, Gera Doc, Linta e Roda
dev: swag tidy lint run
