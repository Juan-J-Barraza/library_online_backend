package repository

import (
	"libraryOnline/models"

	"gorm.io/gorm"
)

type LoanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) baseQuery() *gorm.DB {
	return r.db.Model(&models.Loand{}).
		Preload("User").
		Preload("Book").
		Preload("Book.Editorial").
		Preload("Book.Authors").
		Where("status IN ?", []string{"ACTIVE", "RETURNED"})
}

func (r *LoanRepository) GetAll() (*gorm.DB, []models.Loand, error) {
	query := r.baseQuery()
	var loands []models.Loand

	err := query.Find(&loands).Error
	return query, loands, err
}

func (r *LoanRepository) GetByUserID(userID uint) (*gorm.DB, []models.Loand, error) {
	query := r.baseQuery()
	var loands []models.Loand
	err := query.Where("user_id = ?", userID).Error

	return query, loands, err
}

func (r *LoanRepository) FindByID(id uint) (*models.Loand, error) {
	var loan models.Loand
	err := r.db.
		Preload("User").
		Preload("Book").
		Preload("Book.Editorial").
		Preload("Book.Authors").
		Where("id = ? AND status IN ?", id, []string{"ACTIVE", "RETURNED"}).
		First(&loan).Error
	return &loan, err
}

func (r *LoanRepository) Update(loan *models.Loand) error {
	return r.db.Save(loan).Error
}

func (r *LoanRepository) Create(loan *models.Loand) error {
	return r.db.Create(loan).Error
}
