package handlers_test

import (
	"bytes"
	"encoding/json"
	"go-api-first-steps/internal/handlers"
	"go-api-first-steps/internal/services/product"
	storage "go-api-first-steps/internal/storage/sqlite"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter prepara o ambiente de teste com banco em memória
func setupRouter() *gin.Engine {
	// 1. Banco em memória (SQLite)
	repo := storage.NewRepository(":memory:")

	// 2. Service
	svc := &product.Service{Repo: repo}

	// 3. Handler
	handler := &handlers.ProductHandler{Service: svc}

	// 4. Gin Router (Modo Teste)
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Registra rotas (simplificado, sem Auth para focar no Handler)
	r.POST("/products", handler.Create)
	r.GET("/products", handler.List)

	return r
}

func TestCreateProduct_Success(t *testing.T) {
	router := setupRouter()

	// Payload JSON
	product := map[string]interface{}{
		"name":  "Test Product",
		"price": 99.90,
	}
	jsonValue, _ := json.Marshal(product)

	// Request
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	// Response Recorder
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Asserts
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Test Product")
}

func TestCreateProduct_InvalidJSON(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer([]byte("{invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListProducts(t *testing.T) {
	router := setupRouter()

	// Inserir um produto via request (ou direto no banco) para testar a listagem
	// Aqui vamos inserir via request para testar o fluxo completo
	productBody := []byte(`{"name":"List Item", "price":10}`)
	setupReq, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(productBody))
	setupReq.Header.Set("Content-Type", "application/json")
	setupW := httptest.NewRecorder()
	router.ServeHTTP(setupW, setupReq)
	assert.Equal(t, http.StatusCreated, setupW.Code)

	// Agora lista
	req, _ := http.NewRequest("GET", "/products?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "List Item")
}
