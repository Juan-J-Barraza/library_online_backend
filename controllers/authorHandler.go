package controllers

import (
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
)

type AuthorHandler struct {
	service *services.AuthorService
}

func NewAuthorHandler(service *services.AuthorService) *AuthorHandler {
	return &AuthorHandler{service: service}
}

func (h *AuthorHandler) GetAll(c fiber.Ctx) error {
	authors, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(authors)
}
