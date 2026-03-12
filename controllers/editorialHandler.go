package controllers

import (
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
)

type EditorialHandler struct {
	service *services.EditorialService
}

func NewEditorialHandler(service *services.EditorialService) *EditorialHandler {
	return &EditorialHandler{service: service}
}

func (h *EditorialHandler) GetAll(c fiber.Ctx) error {
	editorials, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(editorials)
}
