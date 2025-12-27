package logger

import (
	"context"
	"log/slog"
)

// TraceIDKey é a chave usada no Contexto
const TraceIDKey = "trace_id"

type ContextHandler struct {
	slog.Handler
}

// NewContextHandler cria o interceptador
func NewContextHandler(h slog.Handler) *ContextHandler {
	return &ContextHandler{Handler: h}
}

// Handle intercepta o log, olha o Contexto e injeta o ID se existir
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		r.AddAttrs(slog.String("trace_id", id))
	}
	return h.Handler.Handle(ctx, r)
}
