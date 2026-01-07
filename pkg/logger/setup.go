package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

// Setup configura o logger global (JSON + Azure opcional)
func Setup(connectionString string, debug bool) {
	// 1. Handler Básico (Terminal)
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)

	// 2. Handler Azure (Se tiver config)
	var azureHandler slog.Handler

	if connectionString != "" {
		iKey, endpoint := parseConnectionString(connectionString)

		// Se falhar o parse (ex: for só a key antiga), tenta usar direto
		if iKey == "" {
			iKey = connectionString
		}

		config := appinsights.NewTelemetryConfiguration(iKey)

		// Configura endpoint regional se existir
		if endpoint != "" {
			config.EndpointUrl = endpoint + "v2/track"
		}

		// Configurações para envio mais rápido (Dev / Teste)
		config.MaxBatchSize = 10
		config.MaxBatchInterval = 2 * time.Second

		client := appinsights.NewTelemetryClientFromConfig(config)

		// DIAGNÓSTICO: Ver erros internos do SDK
		if debug {
			appinsights.NewDiagnosticsMessageListener(func(msg string) error {
				fmt.Printf("🔴 [AppInsights Internal] %s\n", msg)
				return nil
			})
		}

		client.Context().CommonProperties["service"] = "api-go-template" // TODO change to you own app
		azureHandler = NewAzureHandler(client)
	}

	// 3. Unificação (Fanout)
	var finalHandler slog.Handler
	if azureHandler != nil {
		finalHandler = NewFanoutHandler(jsonHandler, azureHandler)
	} else {
		finalHandler = jsonHandler
	}

	// 4. Trace ID Context Wrapper
	ctxHandler := NewContextHandler(finalHandler)

	// 5. Define como Global
	slog.SetDefault(slog.New(ctxHandler))
}

// parseConnectionString extrai InstrumentationKey e IngestionEndpoint
func parseConnectionString(cs string) (string, string) {
	parts := strings.Split(cs, ";")
	var iKey, endpoint string

	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		if strings.EqualFold(key, "InstrumentationKey") {
			iKey = val
		} else if strings.EqualFold(key, "IngestionEndpoint") {
			endpoint = val
		}
	}
	return iKey, endpoint
}
