package logger

import (
	"context"
	"log/slog"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts" // <--- IMPORT NOVO
)

type AzureHandler struct {
	Client appinsights.TelemetryClient
	attrs  []slog.Attr
	group  string
}

func NewAzureHandler(client appinsights.TelemetryClient) *AzureHandler {
	return &AzureHandler{Client: client}
}

func (h *AzureHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *AzureHandler) Handle(ctx context.Context, r slog.Record) error {
	msg := r.Message
	severity := mapSeverity(r.Level)

	// Cria o Trace com o nível correto
	trace := appinsights.NewTraceTelemetry(msg, severity)
	trace.Timestamp = r.Time

	r.Attrs(func(a slog.Attr) bool {
		trace.Properties[a.Key] = a.Value.String()
		return true
	})

	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		trace.Tags.Operation().SetId(id)
		trace.Properties["trace_id"] = id
	}

	h.Client.Track(trace)
	return nil
}

// mapSeverity corrigido: Retorna contracts.SeverityLevel
func mapSeverity(l slog.Level) contracts.SeverityLevel {
	switch {
	case l >= slog.LevelError:
		return contracts.Error
	case l >= slog.LevelWarn:
		return contracts.Warning
	case l >= slog.LevelInfo:
		return contracts.Information
	default:
		return contracts.Verbose
	}
}

func (h *AzureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &AzureHandler{Client: h.Client, attrs: append(h.attrs, attrs...)}
}

func (h *AzureHandler) WithGroup(name string) slog.Handler {
	return &AzureHandler{Client: h.Client, group: name}
}
