package logger

import (
	"log/slog"
	"os"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

// Setup configura o logger global (JSON + Azure opcional)
func Setup() {
	// 1. Handler Básico (Terminal)
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)

	// 2. Handler Azure (Se tiver config)
	var azureHandler slog.Handler
	appInsightsString := os.Getenv("APPINSIGHTS_CONNECTION_STRING")

	if appInsightsString != "" {
		client := appinsights.NewTelemetryClient(appInsightsString)

		// Configurações de contexto (Isso funciona e é útil)
		client.Context().CommonProperties["service"] = "go-api-products"
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
