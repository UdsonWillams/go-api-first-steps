package product

import (
	"go-api-first-steps/internal/storage"
	"testing"
)

func TestCreateAndListProduct(t *testing.T) {
	// 1. Setup: Cria banco em memória (:memory:)
	// Isso garante que o teste não suje o seu banco real 'meubanco.db'
	repo := storage.NewRepository(":memory:")
	service := Service{Repo: repo}

	// 2. Teste de CRIAÇÃO
	createdName, err := service.CreateProduct("Mouse Gamer")

	// Validações (Asserts)
	if err != nil {
		t.Fatalf("Erro inesperado ao criar: %v", err)
	}
	if createdName != "Mouse Gamer" {
		t.Errorf("Esperava 'Mouse Gamer', recebeu '%s'", createdName)
	}

	// 3. Teste de LISTAGEM
	products, err := service.ListProducts()
	if err != nil {
		t.Fatalf("Erro ao listar: %v", err)
	}

	if len(products) != 1 {
		t.Errorf("Esperava 1 produto na lista, encontrou %d", len(products))
	}

	if products[0].Name != "Mouse Gamer" {
		t.Errorf("O nome do produto salvo está errado: %s", products[0].Name)
	}
}

func TestValidateEmptyName(t *testing.T) {
	repo := storage.NewRepository(":memory:")
	service := Service{Repo: repo}

	_, err := service.CreateProduct("")

	if err == nil {
		t.Error("Deveria ter dado erro ao criar produto sem nome, mas não deu.")
	}
}
