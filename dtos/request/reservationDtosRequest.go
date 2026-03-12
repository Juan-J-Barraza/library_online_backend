package request

import "time"

type CreateReservationRequest struct {
	BookId   uint `json:"book_id"   validate:"required"`
	Quantity int  `json:"quantity"  validate:"required,min=1"`
	// Solo el admin/librarian puede especificar el usuario
	// Si viene vacío se toma del token
	UserId             uint       `json:"user_id"`
	ExpectedReturnDate *time.Time `json:"expected_return_date"`
}

type CancelReservationRequest struct {
	Reason string `json:"reason"`
}
