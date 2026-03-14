package repository

import (
	"libraryOnline/dtos/filters"
	"libraryOnline/models"
	"time"

	"gorm.io/gorm"
)

type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) baseQuery() *gorm.DB {
	return r.db.Model(&models.Loand{}).
		Preload("User").
		Preload("Book").
		Preload("Book.Editorial").
		Preload("Book.Authors").
		Where("status = ?", "RESERVED")
}

func (r *ReservationRepository) GetAll(f filters.FilterReservation) (*gorm.DB, []models.Loand, error) {
	var loans []models.Loand
	query := r.baseQuery()
	if f.BookName != "" {
		query = query.
			Joins("JOIN books ON books.id = loands.book_id").
			Where("books.title ILIKE ?", "%"+f.BookName+"%")
	}

	err := query.Find(&loans).Error
	return query, loans, err
}

func (r *ReservationRepository) GetByUserID(userID uint, f filters.FilterReservation) (*gorm.DB, []models.Loand, error) {
	query := r.baseQuery()
	if f.BookName != "" {
		query = query.
			Joins("JOIN books ON books.id = loands.book_id").
			Where("books.title ILIKE ?", "%"+f.BookName+"%")
	}

	var loans []models.Loand
	err := query.Where("user_id = ?", userID).Find(&loans).Error
	return query, loans, err
}

func (r *ReservationRepository) FindByID(id uint) (*models.Loand, error) {
	var loan models.Loand
	err := r.db.
		Preload("User").
		Preload("Book").
		Preload("Book.Editorial").
		Preload("Book.Authors").
		Where("id = ? AND status = ?", id, "RESERVED").
		First(&loan).Error
	return &loan, err
}

func (r *ReservationRepository) Create(loan *models.Loand) error {
	return r.db.Create(loan).Error
}

func (r *ReservationRepository) Update(loan *models.Loand) error {
	return r.db.Save(loan).Error
}

// FindExpired busca reservas que superaron su fecha límite
func (r *ReservationRepository) FindExpired() ([]models.Loand, error) {
	var loans []models.Loand
	err := r.db.
		Where("status = ? AND expected_return_date < ?", "RESERVED", time.Now()).
		Find(&loans).Error
	return loans, err
}

func (r *ReservationRepository) CountActiveByUserID(userID uint) (int, error) {
	var total int
	err := r.db.Model(&models.Loand{}).
		Select("COALESCE(SUM(quantity), 0)").
		Where("user_id = ? AND status IN ?", userID, []string{"RESERVED", "ACTIVE"}).
		Scan(&total).Error
	return total, err
}
