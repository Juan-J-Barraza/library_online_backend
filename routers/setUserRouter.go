package routers

import (
	"libraryOnline/controllers"
	"libraryOnline/middleware"
	"libraryOnline/repository"
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetUserRouter(db *gorm.DB, apiv1 fiber.Router) {
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHanlder := controllers.NewUserHandler(userService)

	allowedRole := middleware.RoleMiddleware("ADMIN")

	apiv1.Get("/users", allowedRole, userHanlder.GetAll)
	apiv1.Get("/users/:id", userHanlder.FindByID)
	apiv1.Put("/users/:id", userHanlder.Update)
	apiv1.Delete("/users/:id", userHanlder.Delete)
}
