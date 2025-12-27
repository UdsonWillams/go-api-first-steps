package handlers

import (
	"encoding/json"
	"go-api-first-steps/internal/product"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	Service *product.Service
}

// POST /products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}
	name, err := h.Service.CreateProduct(body.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Criado: " + name))
}

// GET /products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, _ := h.Service.ListProducts()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// PUT /products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // Go 1.22+ feature
	id, _ := strconv.Atoi(idStr)

	var body struct {
		Name string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := h.Service.UpdateProduct(uint(id), body.Name); err != nil {
		http.Error(w, "Erro ao atualizar", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Atualizado com sucesso"))
}

// DELETE /products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.Atoi(idStr)

	h.Service.DeleteProduct(uint(id))
	w.Write([]byte("Deletado com sucesso"))
}
