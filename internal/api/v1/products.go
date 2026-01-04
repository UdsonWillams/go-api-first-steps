package v1

import (
	"go-api-first-steps/internal/handlers"
	"go-api-first-steps/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerProductRoutes(router *gin.RouterGroup, auth *middleware.Authenticator, h *handlers.ProductHandler) {
	products := router.Group("/products")
	{
		products.GET("", auth.CheckMiddleware("OR"), h.List)
		products.POST("", auth.CheckMiddleware("OR", "admin"), h.Create)
		products.PUT("/:id", auth.CheckMiddleware("OR", "manager"), h.Update)
		products.DELETE("/:id", auth.CheckMiddleware("OR", "admin"), h.Delete)
	}
}
