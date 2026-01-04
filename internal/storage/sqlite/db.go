package storage

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Product representa a entidade Produto no banco de dados.
type Product struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt

	// Tags controlam o comportamento do GORM:
	// unique: constraint de unicidade.
	// not null: campo obrigatório.
	// type:text: define o tipo da coluna no SQLite.
	Name  string  `json:"name" gorm:"type:text;unique;not null"`
	Price float64 `json:"price" gorm:"default:0"`
}

// Repository gerencia a persistência de produtos usando GORM.
type Repository struct {
	DB *gorm.DB
}

// NewRepository inicializa a conexão com o banco de dados e roda migrações.
//
// Parâmetros:
//   - DBPath: Caminho para o arquivo arquivo.db ou ":memory:" para testes.
func NewRepository(DBPath string) *Repository {
	// Se DBPath for ":memory:", o banco roda na RAM (para testes)
	DB, err := gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		panic("falha ao conectar no banco")
	}
	if err := DB.AutoMigrate(&Product{}); err != nil {
		panic("Falha ao rodar migration: " + err.Error())
	}

	return &Repository{DB: DB}
}

func (r *Repository) Save(name string) (*Product, error) {
	p := Product{Name: name}
	result := r.DB.Create(&p)
	return &p, result.Error
}

// FindAll recupera produtos com suporte a paginação.
func (r *Repository) FindAll(page, pageSize int) ([]Product, error) {
	var products []Product
	offset := (page - 1) * pageSize
	result := r.DB.Offset(offset).Limit(pageSize).Find(&products)
	return products, result.Error
}

func (r *Repository) Update(id uint, newName string) error {
	var p Product
	// Primeiro busca, depois atualiza
	if err := r.DB.First(&p, id).Error; err != nil {
		return err
	}
	p.Name = newName
	return r.DB.Save(&p).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&Product{}, id).Error
}
