package request

type CreateOrUpdatedUserRequest struct {
	Name     string `json:"name"      validate:"required,min=2,max=50"`
	LastName string `json:"last_name" validate:"required,min=2,max=50"`
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8"`
	Role     string `json:"role"      validate:"required,oneof=PROFESOR ESTUDIANTE ADMIN"`
}


type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
