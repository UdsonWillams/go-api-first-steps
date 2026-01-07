package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api-first-steps/internal/api"
	"go-api-first-steps/internal/config"
	"go-api-first-steps/internal/dependencies"
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
	// 1. Carregar Configurações (Carrega .env e variáveis)
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Falha ao carregar configurações", "error", err)
		os.Exit(1)
	}

	// 2. Configuração Básica de Logs (Com App Insights se configurado)
	logger.Setup(cfg.AppInsightsConnectionString, cfg.DevMode)

	// 3. Injeção de Dependências
	// Usamos o container para não poluir o main com construções complexas
	ctn := dependencies.NewContainer(cfg)

	// 4. Configuração do Servidor (Rotas)
	r := api.NewRouter(cfg, ctn)

	// 5. Configurar servidor HTTP
	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	// 6. Iniciar servidor em goroutine
	go func() {
		slog.Info("Servidor iniciado", "port", cfg.Port, "azure_enabled", cfg.AppInsightsConnectionString != "")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Erro ao iniciar servidor", "error", err)
			os.Exit(1)
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Desligando servidor graciosamente...")

	// Timeout de 10 segundos para conexões ativas finalizarem
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Erro ao desligar servidor", "error", err)
		os.Exit(1)
	}

	slog.Info("Servidor desligado com sucesso")
}
