package repository

import (
	"libraryOnline/dtos/response"

	"gorm.io/gorm"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetStats() (*response.DashboardResponse, error) {
	var stats response.DashboardResponse

	// Total de libros (suma de todos los ejemplares)
	r.db.Table("books").
		Where("deleted_at IS NULL").
		Select("COALESCE(SUM(total_quantity), 0)").
		Scan(&stats.TotalBooks)

	// Libros disponibles
	r.db.Table("books").
		Where("deleted_at IS NULL").
		Select("COALESCE(SUM(available_quantity), 0)").
		Scan(&stats.AvailableBooks)

	// Libros prestados (ACTIVE)
	r.db.Table("loands").
		Where("deleted_at IS NULL AND status = ?", "ACTIVE").
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&stats.BorrowedBooks)

	// Libros reservados (RESERVED)
	r.db.Table("loands").
		Where("deleted_at IS NULL AND status = ?", "RESERVED").
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&stats.ReservedBooks)

	return &stats, nil
}
