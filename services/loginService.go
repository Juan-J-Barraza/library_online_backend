package services

import (
	"errors"
	"fmt"
	"libraryOnline/dtos/request"
	"libraryOnline/dtos/response"
	"libraryOnline/models"
	"libraryOnline/repository"
	"libraryOnline/utils"

	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	userRepo *repository.UserRepository
}

func NewLoginService(userRepo *repository.UserRepository) *LoginService {
	return &LoginService{userRepo: userRepo}
}

func (s *LoginService) CreateUser(u *request.CreateOrUpdatedUserRequest) (*response.LoginResponse, error) {
	existing, _ := s.userRepo.FindByEmail(u.Email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(hash)

	user := &models.User{
		Name: u.Name, LastName: u.LastName,
		Email: u.Email, Password: u.Password,
		Role: u.Role,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("error al crear el usuario")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("error al generar token")
	}

	return &response.LoginResponse{
		Token: token,
		User: response.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

func (s *LoginService) Login(req request.LoginRequest) (*response.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, errors.New("error al generar el token")
	}

	return &response.LoginResponse{
		Token: token,
		User: response.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}
