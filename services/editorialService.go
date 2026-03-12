package services

import (
	"libraryOnline/dtos/response"
	"libraryOnline/repository"
)

type EditorialService struct {
	repo *repository.EditorialRepository
}

func NewEditorialService(repo *repository.EditorialRepository) *EditorialService {
	return &EditorialService{repo: repo}
}

func (s *EditorialService) GetAll() ([]response.EditorialResponse, error) {
	editorials, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]response.EditorialResponse, len(editorials))
	for i, e := range editorials {
		responses[i] = response.EditorialResponse{ID: e.ID, Name: e.Name}
	}
	return responses, nil
}
