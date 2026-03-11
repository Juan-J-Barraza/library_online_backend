package repository

import (
	"fmt"
	"libraryOnline/utils"

	"gorm.io/gorm"
)

type PaginationRepository struct {
	db *gorm.DB
}

func NewPaginationRepository(db *gorm.DB) *PaginationRepository {
	return &PaginationRepository{db: db}
}

func (r *PaginationRepository) GetPaginatedResults(query *gorm.DB, pagination *utils.Pagination, result interface{}) (*utils.Pagination, error) {
	var totalItems int64

	if query == nil {
		return nil, fmt.Errorf("query is nil")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	pagination.TotalItems = int(totalItems)

	offset := (pagination.Page - 1) * pagination.PageSize

	if err := query.Limit(pagination.PageSize).Offset(offset).Find(result).Error; err != nil {
		return nil, err
	}

	pagination.Calculate()

	pagination.Data = result

	return pagination, nil
}
