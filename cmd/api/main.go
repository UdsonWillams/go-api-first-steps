package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"

	"go-api-first-steps/internal/handlers"
	"go-api-first-steps/internal/middleware"
	"go-api-first-steps/internal/product"
	"go-api-first-steps/internal/storage"
	"go-api-first-steps/pkg/logger"

	// IMPORTANTE: Importe a pasta docs gerada pelo swag (mesmo que ainda não exista)
	_ "go-api-first-steps/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           API de Produtos Go
// @version         1.0
// @description     API com Autenticação Keycloak, GORM, Slog e Docker.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Suporte
// @contact.email   suporte@exemplo.com

// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// --- 1. Configuração de Logs ---
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	ctxHandler := logger.NewContextHandler(jsonHandler)
	slog.SetDefault(slog.New(ctxHandler))

	// --- 2. Variáveis de Ambiente ---
	if err := godotenv.Load(); err != nil {
		slog.Warn("Nenhum arquivo .env encontrado")
	}

	// 2. Configura Logs (Terminal + Azure)
	// Toda a complexidade foi para pkg/logger/setup.go
	logger.Setup()

	// 3. Configurações
	port := getEnv("PORT", ":8080")
	dbUrl := getEnv("DB_URL", "meubanco.db")

	keycloakStr := os.Getenv("KEYCLOAK_PUBLIC_KEY")
	if keycloakStr == "" {
		slog.Error("ERRO CRITICO: KEYCLOAK_PUBLIC_KEY ausente")
		os.Exit(1)
	}
	rsaPublicKey, err := parseKeycloakKey(keycloakStr)
	if err != nil {
		slog.Error("Erro na chave Keycloak", "error", err)
		os.Exit(1)
	}

	// 4. Injeção de Dependências
	repo := storage.NewRepository(dbUrl)
	service := &product.Service{Repo: repo}
	handler := &handlers.ProductHandler{Service: service}

	// 5. Server
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger()) // Middleware gera trace_id

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api")
	{
		api.GET("/products", middleware.Auth(rsaPublicKey, ""), handler.List)
		api.POST("/products", middleware.Auth(rsaPublicKey, "admin"), handler.Create)
		api.PUT("/products/:id", middleware.Auth(rsaPublicKey, "manager"), handler.Update)
		api.DELETE("/products/:id", middleware.Auth(rsaPublicKey, "admin"), handler.Delete)
	}

	slog.Info("Servidor iniciado", "port", port, "azure_enabled", os.Getenv("APPINSIGHTS_CONNECTION_STRING") != "")
	if err := r.Run(port); err != nil {
		slog.Error("O servidor parou inesperadamente", "error", err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func parseKeycloakKey(base64Key string) (*rsa.PublicKey, error) {
	pemStr := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", base64Key)
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("bloco PEM inválido")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}
