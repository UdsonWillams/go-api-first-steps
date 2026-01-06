package storage

import (
	"time"

	"go-api-first-steps/internal/domain"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// ProductModel representa a entidade Produto no banco de dados (GORM Model).
// É separado do domain.Product para manter o domínio agnóstico de ORM.
type ProductModel struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt

	// Tags controlam o comportamento do GORM:
	// unique: constraint de unicidade.
	// not null: campo obrigatório.
	// type:text: define o tipo da coluna no SQLite.
	Name  string  `json:"name" gorm:"type:text;unique;not null"`
	Price float64 `json:"price" gorm:"default:0"`
}

// TableName define o nome da tabela no banco (mantém compatibilidade)
func (ProductModel) TableName() string {
	return "products"
}

// toDomain converte o model GORM para domain.Product
func (p *ProductModel) toDomain() *domain.Product {
	var deletedAt *time.Time
	if p.DeletedAt.Valid {
		deletedAt = &p.DeletedAt.Time
	}
	return &domain.Product{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

// Repository gerencia a persistência de produtos usando GORM.
// Implementa domain.ProductRepository.
type Repository struct {
	DB *gorm.DB
}

// Garantia em tempo de compilação que Repository implementa a interface
var _ domain.ProductRepository = (*Repository)(nil)

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
	if err := DB.AutoMigrate(&ProductModel{}); err != nil {
		panic("Falha ao rodar migration: " + err.Error())
	}

	return &Repository{DB: DB}
}

func (r *Repository) Save(name string) (*domain.Product, error) {
	p := ProductModel{Name: name}
	result := r.DB.Create(&p)
	if result.Error != nil {
		return nil, result.Error
	}
	return p.toDomain(), nil
}

// FindAll recupera produtos com suporte a paginação.
func (r *Repository) FindAll(page, pageSize int) ([]domain.Product, error) {
	var models []ProductModel
	offset := (page - 1) * pageSize
	result := r.DB.Offset(offset).Limit(pageSize).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	products := make([]domain.Product, len(models))
	for i, m := range models {
		products[i] = *m.toDomain()
	}
	return products, nil
}

// FindByID busca um produto pelo ID.
func (r *Repository) FindByID(id uint) (*domain.Product, error) {
	var p ProductModel
	if err := r.DB.First(&p, id).Error; err != nil {
		return nil, err
	}
	return p.toDomain(), nil
}

func (r *Repository) Update(id uint, newName string) error {
	var p ProductModel
	// Primeiro busca, depois atualiza
	if err := r.DB.First(&p, id).Error; err != nil {
		return err
	}
	p.Name = newName
	return r.DB.Save(&p).Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&ProductModel{}, id).Error
}
