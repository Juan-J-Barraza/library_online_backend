package routers

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetRouters(db *gorm.DB, app *fiber.App) {

	apiV1 := app.Group("/api/v1")
	SetHealthCheckRouter(apiV1)

}
