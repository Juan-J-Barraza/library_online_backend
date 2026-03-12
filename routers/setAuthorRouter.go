package routers

import (
	"libraryOnline/controllers"
	"libraryOnline/middleware"
	"libraryOnline/repository"
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetAuthorRouter(db *gorm.DB, apiv1 fiber.Router) {
	authorRepo := repository.NewAuthorRepository(db)
	authorService := services.NewAuthorService(authorRepo)
	authorhandler := controllers.NewAuthorHandler(authorService)
	allowedRole := middleware.RoleMiddleware("ADMIN")

	apiv1.Get("/authors", allowedRole, authorhandler.GetAll)
}
