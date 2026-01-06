package api

import (
	"context"
	v1 "go-api-first-steps/internal/api/v1"
	"go-api-first-steps/internal/config"
	"go-api-first-steps/internal/handlers"
	"go-api-first-steps/internal/middleware"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter configura e retorna o motor do Gin com todas as rotas, middlewares e configurações necessárias.
//
// Ele registra:
//   - Middleware de Logger e Recovery.
//   - Rotas do Swagger UI.
//   - Rotas de Health Check.
//   - Grupos de API versionados (ex: /api/v1).
func NewRouter(cfg *config.Config, productHandler *handlers.ProductHandler) *gin.Engine {
	// Logger Configuration
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	// Initialize OIDC Authenticator
	var authenticator *middleware.Authenticator
	var err error

	if cfg.DevMode {
		slog.WarnContext(
			context.Background(),
			"⚠️  DevMode ativo - usando autenticador de desenvolvimento (sem validação real)")
		authenticator = middleware.NewDevAuthenticator()
	} else {
		authenticator, err = middleware.NewAuthenticator(cfg)
		if err != nil {
			slog.WarnContext(
				context.Background(),
				"Failed to initialize Authenticator",
				"error", err)
			// Log error but allow startup. Auth middleware will return 500 if verifier is missing.
			// In production, you might want to os.Exit(1) here.
		}
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", handlers.HealthCheck)

	// API V1 Config
	apiV1 := r.Group("/api/v1")
	{
		// Pass dependencies to V1 router
		v1.RegisterRoutes(apiV1, authenticator, productHandler)
	}

	return r
}
