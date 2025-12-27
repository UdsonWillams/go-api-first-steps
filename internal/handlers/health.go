package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck responde se a API está online (Versão Gin)
// @Summary      Verifica saúde da API
// @Description  Retorna status 200 se a API estiver rodando
// @Tags         sistema
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"msg":    "API rodando liso!",
	})
}
