package repository

import (
	"libraryOnline/models"

	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) GetAll() ([]models.Author, error) {
	var authors []models.Author
	err := r.db.Find(&authors).Error
	return authors, err
}

func (r *AuthorRepository) FindByIds(ids []uint) ([]models.Author, error) {
	var authors []models.Author
	err := r.db.Where("id IN ?", ids).Find(&authors).Error
	return authors, err
}