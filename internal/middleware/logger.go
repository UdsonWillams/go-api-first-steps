package middleware

import (
	"context"
	"log/slog"
	"time"

	"go-api-first-steps/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 1. Gera ID único
		traceID := uuid.NewString()

		// 2. Injeta no Contexto do Go
		// Isso permite que o slog.InfoContext pegue esse valor depois
		ctx := context.WithValue(c.Request.Context(), logger.TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		// 3. Devolve no Header para o cliente saber qual é o ID
		c.Header("X-Trace-ID", traceID)

		c.Next()

		// 4. Log final da requisição
		duration := time.Since(start)
		status := c.Writer.Status()

		slog.InfoContext(c.Request.Context(), "Requisição finalizada",
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("duration", duration.String()),
		)
	}
}
