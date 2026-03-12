package services

import (
	"libraryOnline/dtos/response"
	"libraryOnline/repository"
)

type AuthorService struct {
	repo *repository.AuthorRepository
}

func NewAuthorService(repo *repository.AuthorRepository) *AuthorService {
	return &AuthorService{repo: repo}
}

func (s *AuthorService) GetAll() ([]response.AuthorResponse, error) {
	authors, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]response.AuthorResponse, len(authors))
	for i, a := range authors {
		responses[i] = response.AuthorResponse{ID: a.ID, Name: a.Name, LastName: a.LastName}
	}
	return responses, nil
}
