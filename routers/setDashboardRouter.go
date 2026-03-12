package routers

import (
	"libraryOnline/controllers"
	"libraryOnline/middleware"
	"libraryOnline/repository"
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetDashboardRouter(db *gorm.DB, apiv1 fiber.Router) {
	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardService := services.NewDashboardService(dashboardRepo)
	dashboardHandler := controllers.NewDashboardHandler(dashboardService)
	allowedAdmin := middleware.RoleMiddleware("ADMIN")

	apiv1.Get("/dashboard", allowedAdmin, dashboardHandler.GetStats)
}
