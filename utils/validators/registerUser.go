package validators

import "libraryOnline/dtos/request"

func ValidatorUser(req request.CreateOrUpdatedUserRequest) string {
	if req.Name == "" {
		return "El nombre es requerido"
	}

	if req.LastName == "" {
		return "El apellido es requerido"
	}

	if req.Email == "" {
		return "El email es requerido"
	}

	if req.Password == "" {
		return "La contraseña es requerida"
	}

	if req.Role == "" {
		return "El rol es requerido"
	}

	return ""
}
