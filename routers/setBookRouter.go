package routers

import (
	"libraryOnline/controllers"
	"libraryOnline/middleware"
	"libraryOnline/repository"
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetBooRouter(db *gorm.DB, apiv1 fiber.Router) {
	bookRepo := repository.NewBookRepository(db)
	authorRepo := repository.NewAuthorRepository(db)
	editorialRepo := repository.NewEditorialRepository(db)
	paginationRepo := repository.NewPaginationRepository(db)
	bookService := services.NewBookService(bookRepo, authorRepo, editorialRepo, paginationRepo)
	bookHanlder := controllers.NewBookHandler(bookService)

	allowedRole := middleware.RoleMiddleware("ADMIN")

	apiv1.Post("/books", allowedRole, bookHanlder.Create)
	apiv1.Get("/books", bookHanlder.GetAll)
	apiv1.Get("/books/:id", bookHanlder.FindByID)
	apiv1.Put("/books/:id", allowedRole, bookHanlder.Update)
	apiv1.Delete("/books/:id", allowedRole, bookHanlder.Delete)
}
