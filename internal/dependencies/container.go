package dependencies

import (
	"go-api-first-steps/internal/config"
	"go-api-first-steps/internal/handlers"
	"go-api-first-steps/internal/services/product"
	sqliteRepo "go-api-first-steps/internal/storage/sqlite"
)

// Container mantém todas as dependências da aplicação inicializadas.
// Centraliza a criação de objetos (Wiring) para manter o main.go limpo.
type Container struct {
	ProductHandler *handlers.ProductHandler
}

// NewContainer inicializa todas as dependências do projeto.
// Aqui é o único lugar onde o acoplamento concreto deve acontecer.
func NewContainer(cfg *config.Config) *Container {
	// Repositories
	repo := sqliteRepo.NewRepository(cfg.DBUrl)

	// Services
	service := product.NewService(repo)

	// Handlers
	productHandler := &handlers.ProductHandler{Service: service}

	return &Container{
		ProductHandler: productHandler,
	}
}
