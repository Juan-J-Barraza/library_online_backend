package models

import "gorm.io/gorm"

type Editorial struct {
	gorm.Model
	Name string 
}
