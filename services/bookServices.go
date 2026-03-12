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
)

type BookService struct {
	repo           *repository.BookRepository
	authorRepo     *repository.AuthorRepository
	editorialRepo  *repository.EditorialRepository
	paginationrepo *repository.PaginationRepository
}

func NewBookService(
	repo *repository.BookRepository,
	authorRepo *repository.AuthorRepository,
	editorialRepo *repository.EditorialRepository,
	paginationrepo *repository.PaginationRepository,
) *BookService {
	return &BookService{repo: repo,
		authorRepo:     authorRepo,
		editorialRepo:  editorialRepo,
		paginationrepo: paginationrepo,
	}
}

func (s *BookService) GetAll(f filters.FiltersBook, p *utils.Pagination) (*utils.Pagination, error) {
	query, books, err := s.repo.GetAll(f)
	if err != nil {
		return nil, fmt.Errorf("error al obtener libros")
	}
	result, err := s.paginationrepo.GetPaginatedResults(query, p, &books)
	if err != nil {
		return nil, err
	}

	responses := make([]response.BookResponse, len(books))
	for i, b := range books {
		responses[i] = response.ToBookResponse(b)
	}
	result.Data = responses
	return result, nil
}

func (s *BookService) FindByID(id uint) (*response.BookResponse, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("book not found")
	}
	res := response.ToBookResponse(*book)
	return &res, nil
}

func (s *BookService) Create(req request.CreateOrUpdateBookRequest) (*response.BookResponse, error) {
	// Verifica que no exista un libro con el mismo título
	existing, _ := s.repo.FindByTitle(req.Title)
	if existing != nil && existing.ID != 0 {
		return nil, errors.New("libro existente con el mismo titulo")
	}

	// Verifica que la editorial exista
	_, err := s.editorialRepo.FindByID(req.EditorialId)
	if err != nil {
		return nil, errors.New("editorial no encontrado")
	}

	// Verifica y obtiene los autores
	authors, err := s.authorRepo.FindByIds(req.AuthorIds)
	if err != nil || len(authors) == 0 {
		return nil, errors.New("no se encontro autor")
	}

	book := models.Book{
		Title:             req.Title,
		AvailableQuantity: req.AvailableQuantity,
		TotalQuantity:     req.TotalQuantity,
		Image:             req.Image,
		EditorialId:       req.EditorialId,
		Authors:           authors,
	}

	if err := s.repo.Create(&book); err != nil {
		return nil, err
	}

	// Recarga con relaciones
	created, err := s.repo.FindByID(book.ID)
	if err != nil {
		return nil, err
	}

	res := response.ToBookResponse(*created)
	return &res, nil
}

func (s *BookService) Update(id uint, req request.CreateOrUpdateBookRequest) (*response.BookResponse, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("book not found")
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.AvailableQuantity != 0 {
		book.AvailableQuantity = req.AvailableQuantity
	}
	if req.TotalQuantity != 0 {
		book.TotalQuantity = req.TotalQuantity
	}
	if req.Image != "" {
		book.Image = req.Image
	}
	if req.EditorialId != 0 {
		_, err := s.editorialRepo.FindByID(req.EditorialId)
		if err != nil {
			return nil, errors.New("editorial no enncontrado")
		}
		book.EditorialId = req.EditorialId
	}
	if len(req.AuthorIds) > 0 {
		authors, err := s.authorRepo.FindByIds(req.AuthorIds)
		if err != nil || len(authors) == 0 {
			return nil, errors.New("autores no encontrados")
		}
		if err := s.repo.UpdateAuthors(book, authors); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(book); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(id)
	res := response.ToBookResponse(*updated)
	return &res, nil
}

func (s *BookService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("libro no entoncrado")
	}
	activeBookLoand, err := s.repo.ActiveLoansBooks(id)
	if err != nil {
		return fmt.Errorf("error al obtener estado")
	}

	if activeBookLoand {
		return fmt.Errorf("No se puede eliminar el libro porque esta activo o reservado")
	}

	return s.repo.Delete(id)
}
