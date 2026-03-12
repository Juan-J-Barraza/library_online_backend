package routers

import (
	"libraryOnline/controllers"
	"libraryOnline/middleware"
	"libraryOnline/repository"
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetLoanRouter(db *gorm.DB, apiv1 fiber.Router) {
	loanRepo := repository.NewLoanRepository(db)
	reservationRepo := repository.NewReservationRepository(db)
	bookRepo := repository.NewBookRepository(db)
	userRepo := repository.NewUserRepository(db)
	paginationRepo := repository.NewPaginationRepository(db)

	loanService := services.NewLoanService(loanRepo, reservationRepo, bookRepo, userRepo, paginationRepo)
	loanHandler := controllers.NewLoanHandler(loanService)

	onlyAdmin := middleware.RoleMiddleware("ADMIN")

	apiv1.Get("/loans", loanHandler.GetAll)
	apiv1.Get("/loans/:id", loanHandler.FindByID)

	apiv1.Post("/loans/direct", onlyAdmin, loanHandler.CreateDirect)
	apiv1.Patch("/loans/:id/confirm", onlyAdmin, loanHandler.ConfirmLoan)
	apiv1.Patch("/loans/:id/return", onlyAdmin, loanHandler.ReturnLoan)
}
