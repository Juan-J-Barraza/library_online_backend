package routers

import (
	"libraryOnline/controllers"
	"libraryOnline/repository"
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetGalleryRouter(db *gorm.DB, apiv1 fiber.Router) {
	bookRepo := repository.NewBookRepository(db)
	authorRepo := repository.NewAuthorRepository(db)
	editorialRepo := repository.NewEditorialRepository(db)
	paginationRepo := repository.NewPaginationRepository(db)
	bookService := services.NewBookService(bookRepo, authorRepo, editorialRepo, paginationRepo)
	bookHanlder := controllers.NewBookHandler(bookService)

	apiv1.Get("/books", bookHanlder.GetAll)

}
