package routers

import (
	"libraryOnline/controllers"
	"libraryOnline/middleware"
	"libraryOnline/repository"
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetEditorialRouter(db *gorm.DB, apiv1 fiber.Router) {
	editRepo := repository.NewEditorialRepository(db)
	editService := services.NewEditorialService(editRepo)
	editHandler := controllers.NewEditorialHandler(editService)
	allowedRole := middleware.RoleMiddleware("ADMIN")

	apiv1.Get("/editorials", allowedRole, editHandler.GetAll)
}
