package v1

import (
	"go-api-first-steps/internal/handlers"
	"go-api-first-steps/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, auth *middleware.Authenticator, productHandler *handlers.ProductHandler) {
	// Register Product Routes
	registerProductRoutes(router, auth, productHandler)
}
