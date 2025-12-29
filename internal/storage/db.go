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
	DB *gorm.DB
}

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

func (r *Repository) FindAll() ([]Product, error) {
	var products []Product
	result := r.DB.Find(&products)
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
