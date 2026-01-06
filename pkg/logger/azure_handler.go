package logger

import (
	"context"
	"log/slog"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
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

	// Adiciona attrs pré-configurados (via WithAttrs)
	for _, attr := range h.attrs {
		key := attr.Key
		if h.group != "" {
			key = h.group + "." + key
		}
		trace.Properties[key] = attr.Value.String()
	}

	// Adiciona attrs do record
	r.Attrs(func(a slog.Attr) bool {
		key := a.Key
		if h.group != "" {
			key = h.group + "." + key
		}
		trace.Properties[key] = a.Value.String()
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
	newAttrs := make([]slog.Attr, len(h.attrs), len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	newAttrs = append(newAttrs, attrs...)
	return &AzureHandler{Client: h.Client, attrs: newAttrs, group: h.group}
}

func (h *AzureHandler) WithGroup(name string) slog.Handler {
	newGroup := name
	if h.group != "" {
		newGroup = h.group + "." + name
	}
	return &AzureHandler{Client: h.Client, attrs: h.attrs, group: newGroup}
}
