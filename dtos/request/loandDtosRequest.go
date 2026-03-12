package request

import "time"

type ConfirmLoanRequest struct {
	ExpectedReturnDate time.Time `json:"expected_return_date" validate:"required"`
}

type CreateDirectLoanRequest struct {
	UserId             uint       `json:"user_id"              validate:"required"`
	BookId             uint       `json:"book_id"              validate:"required"`
	Quantity           int        `json:"quantity"             validate:"required,min=1"`
	ExpectedReturnDate *time.Time `json:"expected_return_date"`
}