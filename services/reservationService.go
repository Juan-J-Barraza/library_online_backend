package services

import (
	"errors"
	"fmt"
	"libraryOnline/dtos/filters"
	"libraryOnline/dtos/request"
	"libraryOnline/dtos/response"
	"libraryOnline/models"
	"libraryOnline/repository"
	"libraryOnline/utils"
	"time"
)

// Días límite para recoger la reserva antes de cancelarse
const (
	ReservationDeadlineDays = 3
	MaxBooksPerUser         = 5
)

type ReservationService struct {
	repo             *repository.ReservationRepository
	bookRepo         *repository.BookRepository
	userRepo         *repository.UserRepository
	paginatetionRepo *repository.PaginationRepository
}

func NewReservationService(
	repo *repository.ReservationRepository,
	bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository,
	paginatetionRepo *repository.PaginationRepository,
) *ReservationService {
	return &ReservationService{repo: repo, bookRepo: bookRepo, userRepo: userRepo, paginatetionRepo: paginatetionRepo}
}

func (s *ReservationService) GetAll(f filters.FilterReservation, p *utils.Pagination) (*utils.Pagination, error) {
	query, loans, err := s.repo.GetAll(f)
	if err != nil {
		return nil, err
	}
	result, err := s.paginatetionRepo.GetPaginatedResults(query, p, &loans)
	if err != nil {
		return nil, err
	}
	result.Data = toReservationResponses(loans)
	return result, nil
}

func (s *ReservationService) GetByUserID(userID uint, f filters.FilterReservation, p *utils.Pagination) (*utils.Pagination, error) {
	query, loans, err := s.repo.GetByUserID(userID, f)
	if err != nil {
		return nil, err
	}
	result, err := s.paginatetionRepo.GetPaginatedResults(query, p, &loans)
	if err != nil {
		return nil, err
	}
	result.Data = toReservationResponses(loans)
	return result, nil
}

func (s *ReservationService) FindByID(id uint) (*response.ReservationResponse, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("reserva no encontrada")
	}
	res := response.ToReservationResponse(*loan)
	return &res, nil
}

func (s *ReservationService) Create(req request.CreateReservationRequest, claims *utils.Claims) (*response.ReservationResponse, error) {
	// Determina el userID: admin puede especificar otro usuario
	userID := claims.UserID
	if req.UserId != 0 && claims.Role == "ADMIN" {
		userID = req.UserId
	}

	// Verifica que el usuario exista
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	totalReserved, err := s.repo.CountActiveByUserID(userID)
	if err != nil {
		return nil, err
	}
	if totalReserved+req.Quantity > MaxBooksPerUser {
		return nil, fmt.Errorf("límite alcanzado: tienes %d/%d libros reservados y/o prestados", totalReserved, MaxBooksPerUser)
	}

	// Verifica que el libro exista y tenga disponibilidad
	book, err := s.bookRepo.FindByID(req.BookId)
	if err != nil {
		return nil, errors.New("libro no encontrado")
	}
	if book.AvailableQuantity < req.Quantity {
		return nil, errors.New("no hay suficientes ejemplares disponibles")
	}

	// Fecha límite para recoger la reserva
	deadline := time.Now().AddDate(0, 0, ReservationDeadlineDays)
	if req.ExpectedReturnDate != nil {
		deadline = *req.ExpectedReturnDate
	}
	loan := models.Loand{
		UserId:             userID,
		BookId:             req.BookId,
		Status:             "RESERVED",
		Quantity:           req.Quantity,
		ReservationDate:    time.Now(),
		ExpectedReturnDate: &deadline,
	}

	// Descuenta disponibilidad
	book.AvailableQuantity -= req.Quantity
	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	if err := s.repo.Create(&loan); err != nil {
		// Revierte disponibilidad si falla
		book.AvailableQuantity += req.Quantity
		s.bookRepo.Update(book)
		return nil, err
	}

	created, _ := s.repo.FindByID(loan.ID)
	res := response.ToReservationResponse(*created)
	return &res, nil
}

func (s *ReservationService) Cancel(id uint, claims *utils.Claims) error {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("reserva no encontrada")
	}

	// Estudiante/Profesor solo puede cancelar las suyas
	if claims.Role != "ADMIN" && loan.UserId != claims.UserID {
		return errors.New("no tienes permisos para cancelar esta reserva")
	}

	// Devuelve disponibilidad al libro
	book, err := s.bookRepo.FindByID(loan.BookId)
	if err != nil {
		return errors.New("libro no encontrado")
	}
	book.AvailableQuantity += loan.Quantity
	if err := s.bookRepo.Update(book); err != nil {
		return err
	}

	loan.Status = "CANCELLED"
	return s.repo.Update(loan)
}

// CancelExpired cancela automáticamente las reservas vencidas
func (s *ReservationService) CancelExpired() error {
	expired, err := s.repo.FindExpired()
	if err != nil {
		return err
	}

	for _, loan := range expired {
		book, err := s.bookRepo.FindByID(loan.BookId)
		if err != nil {
			continue
		}
		book.AvailableQuantity += loan.Quantity
		s.bookRepo.Update(book)

		loan.Status = "CANCELLED"
		s.repo.Update(&loan)
	}
	return nil
}

func toReservationResponses(loans []models.Loand) []response.ReservationResponse {
	responses := make([]response.ReservationResponse, len(loans))
	for i, l := range loans {
		responses[i] = response.ToReservationResponse(l)
	}
	return responses
}
