package validators

import (
	"libraryOnline/dtos/request"
	"time"
)

func ValidateCreateReservation(req request.CreateReservationRequest) string {
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
