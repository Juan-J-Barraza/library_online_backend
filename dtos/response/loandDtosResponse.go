package response

import (
	"libraryOnline/models"
	"time"
)

type LoanResponse struct {
	ID                 uint         `json:"id"`
	User               UserResponse `json:"user"`
	Book               BookResponse `json:"book"`
	Status             string       `json:"status"`
	Quantity           int          `json:"quantity"`
	ReservationDate    time.Time    `json:"reservation_date"`
	BorrowedDate       *time.Time   `json:"borrowed_date"`
	ExpectedReturnDate *time.Time   `json:"expected_return_date"`
	ActualReturnDate   *time.Time   `json:"actual_return_date"`
}

func ToLoanResponse(l models.Loand) LoanResponse {
	return LoanResponse{
		ID: l.ID,
		User: UserResponse{
			ID: l.User.ID, Name: l.User.Name,
			LastName: l.User.LastName, Email: l.User.Email, Role: l.User.Role,
		},
		Book:               ToBookResponse(l.Book),
		Status:             l.Status,
		Quantity:           l.Quantity,
		ReservationDate:    l.ReservationDate,
		BorrowedDate:       l.BorrowedDate,
		ExpectedReturnDate: l.ExpectedReturnDate,
		ActualReturnDate:   l.ActualReturnDate,
	}
}
