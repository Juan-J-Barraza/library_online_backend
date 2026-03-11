package routers

import (
	"libraryOnline/controllers"
	"libraryOnline/repository"
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetLoginRouter(db *gorm.DB, apiv1 fiber.Router) {
	userRepo := repository.NewUserRepository(db)
	loginService := services.NewLoginService(userRepo)
	hanlerLogin := controllers.NewLoginHandler(loginService)

	apiv1.Post("/auth/register", hanlerLogin.RegisterUser)
	apiv1.Post("/auth/login", hanlerLogin.Login)
}
