package services

import (
	"fmt"
	"libraryOnline/dtos/response"
	"libraryOnline/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetStats() (*response.DashboardResponse, error) {
	dashboar, err := s.repo.GetStats()
	if err != nil {
		return nil, fmt.Errorf("error al obtener dashboard")
	}

	return dashboar, err
}
