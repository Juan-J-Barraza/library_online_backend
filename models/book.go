package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title             string `gorm:"not null"`
	AvailableQuantity int    `gorm:"not null"`
	TotalQuantity     int    `gorm:"not null"`
	Image             string
	EditorialId       uint      `gorm:"index:idx_books_editorial"`
	Editorial         Editorial `gorm:"foreignKey:EditorialId;"`
	Authors           []Author  `gorm:"many2many:book_authors;"`
}
