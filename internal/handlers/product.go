package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"go-api-first-steps/internal/services/product"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service *product.Service
}

// Create cria um novo produto
// @Summary      Cria um produto
// @Description  Cria um novo produto no banco de dados
// @Tags         produtos
// @Accept       json
// @Produce      json
// @Param        request body     handlers.CreateProductRequest true "Dados do Produto"
// @Success      201     {object} handlers.MessageResponse
// @Failure      400     {object} handlers.ErrorResponse
// @Failure      500     {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.WarnContext(c.Request.Context(), "JSON inválido", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "JSON inválido"})
		return
	}

	slog.InfoContext(c.Request.Context(), "Criando produto", "name", req.Name)

	name, err := h.Service.CreateProduct(req.Name)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Erro ao criar", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, MessageResponse{Message: "Criado: " + name})
}

// List lista todos os produtos
// @Summary      Lista produtos
// @Description  Retorna a lista completa de produtos cadastrados
// @Tags         produtos
// @Produce      json
// @Success      200  {array}   handlers.ProductResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	products, err := h.Service.ListProducts()
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "Erro ao listar", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Erro ao listar"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// Update atualiza um produto
// @Summary      Atualiza um produto
// @Description  Atualiza o nome de um produto existente pelo ID
// @Tags         produtos
// @Accept       json
// @Produce      json
// @Param        id      path     int                          true "ID do Produto"
// @Param        request body     handlers.UpdateProductRequest true "Novos dados"
// @Success      200     {object} handlers.MessageResponse
// @Failure      400     {object} handlers.ErrorResponse
// @Failure      500     {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "JSON inválido"})
		return
	}

	if err := h.Service.UpdateProduct(uint(id), req.Name); err != nil {
		slog.ErrorContext(c.Request.Context(), "Erro ao atualizar", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Erro ao atualizar"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Atualizado com sucesso"})
}

// Delete remove um produto
// @Summary      Deleta um produto
// @Description  Remove um produto do banco pelo ID (Soft Delete)
// @Tags         produtos
// @Produce      json
// @Param        id   path      int  true  "ID do Produto"
// @Success      200  {object}  handlers.MessageResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	if err := h.Service.DeleteProduct(uint(id)); err != nil {
		slog.ErrorContext(c.Request.Context(), "Erro ao deletar", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Erro ao deletar"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Deletado com sucesso"})
}
