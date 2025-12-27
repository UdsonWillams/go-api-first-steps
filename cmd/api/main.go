package main

import (
	"log"
	"net/http"

	"go-api-first-steps/internal/handlers"
	"go-api-first-steps/internal/product"
	"go-api-first-steps/internal/storage"
)

func main() {
	// Passamos o caminho do arquivo do banco
	repo := storage.NewRepository("meubanco.db")
	service := &product.Service{Repo: repo}
	handler := &handlers.ProductHandler{Service: service}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", handler.Create)
	mux.HandleFunc("GET /products", handler.List)
	mux.HandleFunc("PUT /products/{id}", handler.Update)
	mux.HandleFunc("DELETE /products/{id}", handler.Delete)

	log.Println("API CRUD rodando na porta :8080...")
	http.ListenAndServe(":8080", mux)
}
