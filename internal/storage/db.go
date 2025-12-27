package storage

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model // ID, CreatedAt, UpdatedAt, DeletedAt

	// Tags controlam o banco:
	// unique: não deixa repetir nome
	// not null: obriga a ter valor
	// size: limita caracteres (varchar)
	Name  string  `json:"name" gorm:"type:text;unique;not null"`
	Price float64 `json:"price" gorm:"default:0"` // Exemplo de campo novo
}

type Repository struct {
	Db *gorm.DB
}

func NewRepository(dbPath string) *Repository {
	// Se dbPath for ":memory:", o banco roda na RAM (para testes)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("falha ao conectar no banco")
	}
	db.AutoMigrate(&Product{})
	return &Repository{Db: db}
}

func (r *Repository) Save(name string) (*Product, error) {
	p := Product{Name: name}
	result := r.Db.Create(&p)
	return &p, result.Error
}

func (r *Repository) FindAll() ([]Product, error) {
	var products []Product
	result := r.Db.Find(&products)
	return products, result.Error
}

func (r *Repository) Update(id uint, newName string) error {
	var p Product
	// Primeiro busca, depois atualiza
	if err := r.Db.First(&p, id).Error; err != nil {
		return err
	}
	p.Name = newName
	return r.Db.Save(&p).Error
}

func (r *Repository) Delete(id uint) error {
	return r.Db.Delete(&Product{}, id).Error
}
