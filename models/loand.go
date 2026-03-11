package models

import (
	"time"

	"gorm.io/gorm"
)

type Loand struct {
	gorm.Model
	UserId             uint
	User               User `gorm:"foreignKey:UserId;"`
	BookId             uint
	Book               Book      `gorm:"foreignKey:BookId;"`
	Status             string    `gorm:"not null;default:'RESERVED'"`
	Quantity           int       `gorm:"not null;default:1"`
	ReservationDate    time.Time `gorm:"not null"`
	BorrowedDate       *time.Time
	ExpectedReturnDate *time.Time
	ActualReturnDate   *time.Time
}
