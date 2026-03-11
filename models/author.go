package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	Name     string `gorm:"not null"`
	LastName string `gorm:"not null"`
}
