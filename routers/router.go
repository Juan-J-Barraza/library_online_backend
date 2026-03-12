package routers

import (
	"libraryOnline/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/gorm"
)

func SetRouters(db *gorm.DB, app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:4200", "http://localhost:51022"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	app.Use(middleware.PaginationMiddleware)

	apiV1 := app.Group("/api/v1")
	SetHealthCheckRouter(apiV1)
	SetLoginRouter(db, apiV1)

	protected := apiV1.Use("/", middleware.AuthRequired())
	SetUserRouter(db, protected)
	SetAuthorRouter(db, protected)
	SetEditorialRouter(db, protected)
	SetBooRouter(db, apiV1)
}
