package repository

import (
	"libraryOnline/models"

	"gorm.io/gorm"
)

type EditorialRepository struct {
	db *gorm.DB
}

func NewEditorialRepository(db *gorm.DB) *EditorialRepository {
	return &EditorialRepository{db: db}
}

func (r *EditorialRepository) GetAll() ([]models.Editorial, error) {
	var editorials []models.Editorial
	err := r.db.Find(&editorials).Error
	return editorials, err
}

func (r *EditorialRepository) FindByID(id uint) (*models.Editorial, error) {
	var editorial models.Editorial
	err := r.db.First(&editorial, id).Error
	return &editorial, err
}
