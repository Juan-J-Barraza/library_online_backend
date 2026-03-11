package services

import (
	"errors"
	"fmt"
	"libraryOnline/dtos/filters"
	"libraryOnline/dtos/request"
	"libraryOnline/dtos/response"
	"libraryOnline/models"
	"libraryOnline/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(u *request.CreateOrUpdatedUserRequest) error {
	existing, _ := s.repo.FindByEmail(u.Email)
	if existing != nil {
		return errors.New("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	if err := s.repo.Create(&models.User{
		Name: u.Name, LastName: u.LastName,
		Email: u.Email, Password: u.Password,
		Role: u.Role,
	}); err != nil {
		return fmt.Errorf("error al crear el usuario")
	}

	return nil
}

func (s *UserService) GetAll(f filters.FiltersUser) ([]response.UserResponse, error) {
	users, err := s.repo.GetAll(f)
	if err != nil {
		return nil, fmt.Errorf("error to get users")
	}

	responseUser := []response.UserResponse{}
	for _, user := range users {
		responseUser = append(responseUser, response.UserResponse{ID: user.ID,
			Name: user.Name, LastName: user.LastName, Email: user.Email, Role: user.Role})
	}

	if users == nil {
		responseUser = []response.UserResponse{}
	}

	return responseUser, nil

}

func (s *UserService) FindByID(id uint) (*response.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &response.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		LastName: user.LastName,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *UserService) Update(id uint, userUpdated request.CreateOrUpdatedUserRequest) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if userUpdated.Name != "" {
		user.Name = userUpdated.Name
	}

	if userUpdated.LastName != "" {
		user.LastName = userUpdated.LastName
	}
	if userUpdated.Email != "" {
		user.Email = userUpdated.Email
	}
	if userUpdated.Role != "" {
		user.Role = userUpdated.Role
	}

	if userUpdated.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(userUpdated.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}

	return s.repo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	hasActiveLoand, err := s.repo.HasActiveLoans(id)
	if hasActiveLoand {
		return fmt.Errorf("No se puede eliminar el usuario por tener prestamos activos")
	}
	return s.repo.Delete(id)
}
