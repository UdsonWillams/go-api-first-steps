package domain

import "time"

// Product representa a entidade de domínio Produto.
// Esta struct é agnóstica de framework e pode ser usada em qualquer camada.
type Product struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Price     float64    `json:"price"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
