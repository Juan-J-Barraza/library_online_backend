package repository

import (
	"libraryOnline/dtos/filters"
	"libraryOnline/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepository) GetAll(filters filters.FiltersUser) ([]models.User, error) {
	var users []models.User
	query := r.db.Model(&models.User{})

	if filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.LastName != "" {
		query = query.Where("last_name ILIKE ?", "%"+filters.LastName+"%")
	}
	if filters.Role != "" {
		query = query.Where("role = ?", filters.Role)
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) HasActiveLoans(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Loand{}).
		Where("user_id = ? AND status IN ?", id, []string{"ACTIVE", "RESERVED"}).
		Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(u *models.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
