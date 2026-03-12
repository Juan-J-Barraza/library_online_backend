package validators

import (
	"libraryOnline/dtos/request"
	"time"
)

func ValidateConfirmLoan(req request.ConfirmLoanRequest) string {
	if req.ExpectedReturnDate.IsZero() {
		return "La fecha de devolución es requerida"
	}
	if req.ExpectedReturnDate.Before(time.Now()) {
		return "La fecha de devolución no puede ser en el pasado"
	}
	return ""
}

func ValidateCreateDirectLoan(req request.CreateDirectLoanRequest) string {
	if req.UserId == 0 {
		return "El usuario es requerido"
	}
	if req.BookId == 0 {
		return "El libro es requerido"
	}
	if req.Quantity <= 0 {
		return "La cantidad debe ser mayor a 0"
	}
	if req.ExpectedReturnDate != nil && req.ExpectedReturnDate.Before(time.Now()) {
		return "La fecha de devolución no puede ser en el pasado"
	}
	return ""
}