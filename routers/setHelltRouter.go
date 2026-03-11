package routers

import (
	"libraryOnline/controllers"

	"github.com/gofiber/fiber/v3"
)

func SetHealthCheckRouter(apiV1 fiber.Router) {

	apiV1.Get("/health", controllers.CheckHealth)
}
