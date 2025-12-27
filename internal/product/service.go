package product

import (
	"fmt"
	"go-api-first-steps/internal/storage"
)

type Service struct {
	Repo *storage.Repository
}

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

func (s *Service) ListProducts() ([]storage.Product, error) {
	return s.Repo.FindAll()
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
