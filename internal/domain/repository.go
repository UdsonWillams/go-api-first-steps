package domain

// ProductRepository define o contrato para persistência de produtos.
// Qualquer implementação (SQLite, PostgreSQL, MongoDB) deve seguir esta interface.
type ProductRepository interface {
	// Save persiste um novo produto e retorna o produto criado.
	Save(name string) (*Product, error)

	// FindAll retorna uma lista paginada de produtos.
	FindAll(page, pageSize int) ([]Product, error)

	// FindByID busca um produto pelo ID.
	FindByID(id uint) (*Product, error)

	// Update atualiza o nome de um produto existente.
	Update(id uint, name string) error

	// Delete remove um produto (soft delete).
	Delete(id uint) error
}
