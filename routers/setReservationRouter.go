package routers

import (
	"libraryOnline/controllers"
	"libraryOnline/repository"
	"libraryOnline/services"
	"libraryOnline/utils/scheduler"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetReservationRouter(db *gorm.DB, apiv1 fiber.Router) {
	reservationRepo := repository.NewReservationRepository(db)
	bookRepo := repository.NewBookRepository(db)
	userRepo := repository.NewUserRepository(db)
	pagiantionRepo := repository.NewPaginationRepository(db)

	reservationService := services.NewReservationService(reservationRepo, bookRepo, userRepo, pagiantionRepo)
	reservationHandler := controllers.NewReservationHandler(reservationService)

	// cron job
	scheduler.StartReservationScheduler(reservationService)

	apiv1.Get("/reservations", reservationHandler.GetAll)
	apiv1.Get("/reservations/:id", reservationHandler.FindByID)
	apiv1.Post("/reservations", reservationHandler.Create)
	apiv1.Patch("/reservations/:id/cancel", reservationHandler.Cancel)
}
