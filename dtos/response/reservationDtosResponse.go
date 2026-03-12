package response

import (
	"libraryOnline/models"
	"time"
)

type ReservationResponse struct {
	ID                 uint         `json:"id"`
	User               UserResponse `json:"user"`
	Book               BookResponse `json:"book"`
	Status             string       `json:"status"`
	Quantity           int          `json:"quantity"`
	ReservationDate    time.Time    `json:"reservation_date"`
	ExpectedReturnDate *time.Time   `json:"expected_return_date"`
}

func ToReservationResponse(l models.Loand) ReservationResponse {
	return ReservationResponse{
		ID: l.ID,
		User: UserResponse{
			ID:       l.User.ID,
			Name:     l.User.Name,
			LastName: l.User.LastName,
			Email:    l.User.Email,
			Role:     l.User.Role,
		},
		Book:               ToBookResponse(l.Book),
		Status:             l.Status,
		Quantity:           l.Quantity,
		ReservationDate:    l.ReservationDate,
		ExpectedReturnDate: l.ExpectedReturnDate,
	}
}
