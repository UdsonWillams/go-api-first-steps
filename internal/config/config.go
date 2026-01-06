package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                        string
	DBUrl                       string
	AppInsightsConnectionString string

	// OIDC Configurations
	KeycloakURL string // Ex: http://localhost:8080/realms/myrealm
	ClientID    string // Ex: my-backend

	// Development Mode
	// Se true, permite rodar sem autenticação (apenas para desenvolvimento local)
	DevMode bool
}

func Load() (*Config, error) {
	// Carrega .env se existir
	if err := godotenv.Load(); err != nil {
		slog.Warn("Nenhum arquivo .env encontrado")
	}

	devMode := strings.ToLower(os.Getenv("DEV_MODE")) == "true"

	cfg := &Config{
		Port:                        getEnv("PORT", ":8080"),
		DBUrl:                       getEnv("DB_URL", "meubanco.db"),
		AppInsightsConnectionString: os.Getenv("APPINSIGHTS_CONNECTION_STRING"),
		KeycloakURL:                 os.Getenv("KEYCLOAK_URL"),
		ClientID:                    os.Getenv("KEYCLOAK_CLIENT_ID"),
		DevMode:                     devMode,
	}

	// Validação básica (apenas em modo produção)
	if !cfg.DevMode {
		if cfg.KeycloakURL == "" {
			return nil, fmt.Errorf("ERRO CRITICO: KEYCLOAK_URL ausente (use DEV_MODE=true para desenvolvimento)")
		}
		if cfg.ClientID == "" {
			return nil, fmt.Errorf("ERRO CRITICO: KEYCLOAK_CLIENT_ID ausente (use DEV_MODE=true para desenvolvimento)")
		}
	} else {
		slog.Warn("⚠️  Modo de desenvolvimento ativado - Autenticação pode estar desabilitada!")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
