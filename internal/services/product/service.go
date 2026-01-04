package product

import (
	"fmt"
	storage "go-api-first-steps/internal/storage/sqlite"
)

// Service encapsula a lógica de negócio relacionada a produtos.
// Ele interage com o Repositório para persistência de dados.
type Service struct {
	Repo *storage.Repository
}

// CreateProduct valida e cria um novo produto.
// Retorna erro se o nome estiver vazio.
func (s *Service) CreateProduct(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("nome vazio")
	}
	p, err := s.Repo.Save(name)
	if err != nil {
		return "", err
	}
	return p.Name, nil
}

// ListProducts retorna uma lista paginada de produtos.
//
// Parâmetros:
//   - page: Número da página (inicia em 1).
//   - pageSize: Quantidade de itens por página (default 10, max 100).
func (s *Service) ListProducts(page, pageSize int) ([]storage.Product, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.Repo.FindAll(page, pageSize)
}

func (s *Service) UpdateProduct(id uint, name string) error {
	if name == "" {
		return fmt.Errorf("nome vazio")
	}
	return s.Repo.Update(id, name)
}

func (s *Service) DeleteProduct(id uint) error {
	return s.Repo.Delete(id)
}
