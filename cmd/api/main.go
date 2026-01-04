package main

import (
	"log/slog"
	"os"

	"go-api-first-steps/internal/api"
	"go-api-first-steps/internal/config"
	"go-api-first-steps/internal/handlers"
	"go-api-first-steps/internal/services/product"
	storage "go-api-first-steps/internal/storage/sqlite"
	"go-api-first-steps/pkg/logger"

	// IMPORTANTE: Importe a pasta docs gerada pelo swag
	_ "go-api-first-steps/cmd/api/swagger"
)

// @title           API de Produtos Go
// @version         1.0
// @description     API com Autenticação Keycloak, GORM, Slog e Docker.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Suporte
// @contact.email   suporte@exemplo.com

// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Configuração Básica de Logs (Antes de tudo)
	logger.Setup()

	// 2. Carregar Configurações
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Falha ao carregar configurações", "error", err)
		os.Exit(1)
	}

	// 3. Injeção de Dependências
	repo := storage.NewRepository(cfg.DBUrl)
	service := &product.Service{Repo: repo}
	handler := &handlers.ProductHandler{Service: service}

	// 4. Configuração do Servidor (Rotas)
	r := api.NewRouter(cfg, handler)

	// 5. Iniciar Servidor
	slog.Info("Servidor iniciado", "port", cfg.Port, "azure_enabled", cfg.AppInsightsConnectionString != "")
	if err := r.Run(cfg.Port); err != nil {
		slog.Error("O servidor parou inesperadamente", "error", err)
	}
}
