package validators

import "libraryOnline/dtos/request"

func ValidateCreateBook(req request.CreateOrUpdateBookRequest) string {
	if req.Title == "" {
		return "El título es requerido"
	}
	if req.TotalQuantity <= 0 {
		return "La cantidad total debe ser mayor a 0"
	}
	if req.AvailableQuantity < 0 {
		return "La cantidad disponible no puede ser negativa"
	}
	if req.AvailableQuantity > req.TotalQuantity {
		return "La cantidad disponible no puede ser mayor a la cantidad total"
	}
	if req.EditorialId == 0 {
		return "La editorial es requerida"
	}
	if len(req.AuthorIds) == 0 {
		return "Al menos un autor es requerido"
	}
	return ""
}

func ValidateUpdateBook(req request.CreateOrUpdateBookRequest) string {
	if req.TotalQuantity < 0 {
		return "La cantidad total no puede ser negativa"
	}
	if req.AvailableQuantity < 0 {
		return "La cantidad disponible no puede ser negativa"
	}
	if req.AvailableQuantity > req.TotalQuantity && req.TotalQuantity != 0 {
		return "La cantidad disponible no puede ser mayor a la cantidad total"
	}
	return ""
}
