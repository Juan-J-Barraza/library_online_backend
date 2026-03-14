package services

import (
	"errors"
	"libraryOnline/dtos/filters"
	"libraryOnline/dtos/request"
	"libraryOnline/dtos/response"
	"libraryOnline/models"
	"libraryOnline/repository"
	"libraryOnline/utils"
	"time"
)

const LoanDeadlineDays = 7

type LoanService struct {
	repo            *repository.LoanRepository
	reservationRepo *repository.ReservationRepository
	bookRepo        *repository.BookRepository
	userRepo        *repository.UserRepository
	paginationRepo  *repository.PaginationRepository
}

func NewLoanService(
	repo *repository.LoanRepository,
	reservationRepo *repository.ReservationRepository,
	bookRepo *repository.BookRepository,
	userRepo *repository.UserRepository,
	paginationRepo *repository.PaginationRepository,

) *LoanService {
	return &LoanService{
		repo:            repo,
		reservationRepo: reservationRepo,
		bookRepo:        bookRepo,
		userRepo:        userRepo,
		paginationRepo:  paginationRepo,
	}
}

func (s *LoanService) GetAll(f filters.FilterLoan, p *utils.Pagination) (*utils.Pagination, error) {
	query, loans, err := s.repo.GetAll(f)
	if err != nil {
		return nil, err
	}
	result, err := s.paginationRepo.GetPaginatedResults(query, p, &loans)
	if err != nil {
		return nil, err
	}
	result.Data = toLoanResponses(loans)
	return result, nil
}

func (s *LoanService) GetByUserID(userID uint, f filters.FilterLoan, p *utils.Pagination) (*utils.Pagination, error) {
	query, loans, err := s.repo.GetByUserID(userID, f)
	if err != nil {
		return nil, err
	}
	result, err := s.paginationRepo.GetPaginatedResults(query, p, &loans)
	if err != nil {
		return nil, err
	}
	result.Data = toLoanResponses(loans)
	return result, nil
}

func (s *LoanService) FindByID(id uint) (*response.LoanResponse, error) {
	loan, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("préstamo no encontrado")
	}
	res := response.ToLoanResponse(*loan)
	return &res, nil
}

// ConfirmLoan — transiciona RESERVED → ACTIVE (confirmar entrega)
func (s *LoanService) ConfirmLoan(reservationID uint, req request.ConfirmLoanRequest) (*response.LoanResponse, error) {
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		return nil, errors.New("reserva no encontrada")
	}

	now := time.Now()
	reservation.Status = "ACTIVE"
	reservation.BorrowedDate = &now
	reservation.ExpectedReturnDate = &req.ExpectedReturnDate

	if err := s.reservationRepo.Update(reservation); err != nil {
		return nil, err
	}

	// Recarga con relaciones para la respuesta
	updated, _ := s.reservationRepo.FindActiveByID(reservationID)
	res := response.ToLoanResponse(*updated)
	return &res, nil
}

// ReturnLoan — transiciona ACTIVE → RETURNED (recibir devolución)
func (s *LoanService) ReturnLoan(loanID uint) (*response.LoanResponse, error) {
	loan, err := s.repo.FindByID(loanID)
	if err != nil {
		return nil, errors.New("préstamo no encontrado")
	}
	if loan.Status != "ACTIVE" {
		return nil, errors.New("solo se pueden devolver préstamos activos")
	}

	now := time.Now()
	loan.Status = "RETURNED"
	loan.ActualReturnDate = &now

	// Devuelve disponibilidad al libro
	book, err := s.bookRepo.FindByID(loan.BookId)
	if err != nil {
		return nil, errors.New("libro no encontrado")
	}
	book.AvailableQuantity += loan.Quantity
	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	if err := s.repo.Update(loan); err != nil {
		return nil, err
	}

	res := response.ToLoanResponse(*loan)
	return &res, nil
}

// CreateDirect — crea préstamo directo sin reserva (solo ADMIN)
func (s *LoanService) CreateDirect(req request.CreateDirectLoanRequest) (*response.LoanResponse, error) {
	_, err := s.userRepo.FindByID(req.UserId)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	book, err := s.bookRepo.FindByID(req.BookId)
	if err != nil {
		return nil, errors.New("libro no encontrado")
	}
	if book.AvailableQuantity < req.Quantity {
		return nil, errors.New("no hay suficientes ejemplares disponibles")
	}

	deadline := time.Now().AddDate(0, 0, LoanDeadlineDays)
	if req.ExpectedReturnDate != nil {
		deadline = *req.ExpectedReturnDate
	}

	now := time.Now()
	loan := models.Loand{
		UserId:             req.UserId,
		BookId:             req.BookId,
		Status:             "ACTIVE",
		Quantity:           req.Quantity,
		ReservationDate:    now,
		BorrowedDate:       &now,
		ExpectedReturnDate: &deadline,
	}

	book.AvailableQuantity -= req.Quantity
	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	if err := s.repo.Create(&loan); err != nil {
		book.AvailableQuantity += req.Quantity
		s.bookRepo.Update(book)
		return nil, err
	}

	created, _ := s.repo.FindByID(loan.ID)
	res := response.ToLoanResponse(*created)
	return &res, nil
}

func toLoanResponses(loans []models.Loand) []response.LoanResponse {
	responses := make([]response.LoanResponse, len(loans))
	for i, l := range loans {
		responses[i] = response.ToLoanResponse(l)
	}
	return responses
}
