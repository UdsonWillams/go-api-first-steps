package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                        string
	DBUrl                       string
	AppInsightsConnectionString string

	// OIDC Configurations
	KeycloakURL string // Ex: http://localhost:8080/realms/myrealm
	ClientID    string // Ex: my-backend
}

func Load() (*Config, error) {
	// Carrega .env se existir
	if err := godotenv.Load(); err != nil {
		slog.Warn("Nenhum arquivo .env encontrado")
	}

	cfg := &Config{
		Port:                        getEnv("PORT", ":8080"),
		DBUrl:                       getEnv("DB_URL", "meubanco.db"),
		AppInsightsConnectionString: os.Getenv("APPINSIGHTS_CONNECTION_STRING"),
		KeycloakURL:                 os.Getenv("KEYCLOAK_URL"),
		ClientID:                    os.Getenv("KEYCLOAK_CLIENT_ID"),
	}

	// Validação básica
	if cfg.KeycloakURL == "" {
		return nil, fmt.Errorf("ERRO CRITICO: KEYCLOAK_URL ausente")
	}
	if cfg.ClientID == "" {
		return nil, fmt.Errorf("ERRO CRITICO: KEYCLOAK_CLIENT_ID ausente")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
