package controllers

import (
	"libraryOnline/dtos/request"
	"libraryOnline/services"

	"github.com/gofiber/fiber/v3"
)

type LoginHandler struct {
	loginService *services.LoginService
}

func NewLoginHandler(loginService *services.LoginService) *LoginHandler {
	return &LoginHandler{loginService: loginService}
}

func (h *LoginHandler) RegisterUser(c fiber.Ctx) error {
	var req request.CreateOrUpdatedUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	response, err := h.loginService.CreateUser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
func (h *LoginHandler) Login(c fiber.Ctx) error {
	var req request.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Datos inválidos",
		})
	}

	res, err := h.loginService.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
