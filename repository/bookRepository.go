package repository

import (
	"libraryOnline/dtos/filters"
	"libraryOnline/models"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) baseQuery() *gorm.DB {
	return r.db.Model(&models.Book{}).
		Preload("Editorial").
		Preload("Authors")
}

func (r *BookRepository) GetAll(f filters.FiltersBook) (*gorm.DB, []models.Book, error) {
	query := r.baseQuery()
	var books []models.Book
	if f.Title != "" {
		query = query.Where("title ILIKE ?", "%"+f.Title+"%")
	}
	if f.EditorialId != 0 {
		query = query.Where("editorial_id = ?", f.EditorialId)
	}
	if f.AuthorId != 0 {
		query = query.Joins("JOIN book_authors ON book_authors.book_id = books.id").
			Where("book_authors.author_id = ?", f.AuthorId)
	}

	err := query.Find(&books).Error
	return query, books, err
}

func (r *BookRepository) FindByID(id uint) (*models.Book, error) {
	var book models.Book
	err := r.baseQuery().Where("books.id = ?", id).First(&book).Error
	return &book, err
}

func (r *BookRepository) FindByTitle(title string) (*models.Book, error) {
	var book models.Book
	err := r.db.Where("title = ?", title).First(&book).Error
	return &book, err
}

func (r *BookRepository) ActiveLoansBooks(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Loand{}).
		Where("book_id = ? AND status IN ?", id, []string{"ACTIVE", "RESERVED"}).
		Count(&count).Error
	return count > 0, err
}

func (r *BookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) UpdateAuthors(book *models.Book, authors []models.Author) error {
	return r.db.Model(book).Association("Authors").Replace(authors)
}

func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}
