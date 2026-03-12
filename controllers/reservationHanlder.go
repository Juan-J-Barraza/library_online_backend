package controllers

import (
	"libraryOnline/dtos/request"
	"libraryOnline/services"
	"libraryOnline/utils"
	"libraryOnline/utils/validators"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type ReservationHandler struct {
	service *services.ReservationService
}

func NewReservationHandler(service *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) GetAll(c fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.Claims)
	p := c.Locals("pagination").(*utils.Pagination)

	var result *utils.Pagination
	var err error

	// Admin ve todas, estudiante/profesor solo las suyas
	if claims.Role == "ADMIN" {
		result, err = h.service.GetAll(p)
	} else {
		result, err = h.service.GetByUserID(claims.UserID, p)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *ReservationHandler) FindByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id inválido"})
	}

	claims := c.Locals("claims").(*utils.Claims)
	reservation, err := h.service.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	// Estudiante/Profesor solo puede ver las suyas
	if claims.Role != "ADMIN" && reservation.User.ID != claims.UserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no tienes permisos"})
	}

	return c.JSON(reservation)
}

func (h *ReservationHandler) Create(c fiber.Ctx) error {
	var req request.CreateReservationRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "body inválido"})
	}

	if msg := validators.ValidateCreateReservation(req); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	claims := c.Locals("claims").(*utils.Claims)
	reservation, err := h.service.Create(req, claims)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(reservation)
}

func (h *ReservationHandler) Cancel(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id inválido"})
	}

	claims := c.Locals("claims").(*utils.Claims)
	if err := h.service.Cancel(uint(id), claims); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
