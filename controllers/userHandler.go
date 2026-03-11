package controllers

import (
	"libraryOnline/dtos/filters"
	"libraryOnline/dtos/request"
	"libraryOnline/services"
	"libraryOnline/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(c fiber.Ctx) error {
	pagination := c.Locals("pagination").(*utils.Pagination)

	f := filters.FiltersUser{
		Name:     c.Query("name"),
		LastName: c.Query("last_name"),
		Role:     c.Query("role"),
	}
	paginatedUsers, err := h.service.GetAll(f, pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(paginatedUsers)
}

func (h *UserHandler) FindByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	claims, ok := c.Locals("claims").(*utils.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "no autenticado",
		})
	}

	if claims.Role != "ADMIN" && claims.UserID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "No es tu usuario",
		})
	}
	user, err := h.service.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	claims, ok := c.Locals("claims").(*utils.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "no autenticado",
		})
	}

	if claims.Role != "ADMIN" && claims.UserID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "No es tu usuario",
		})
	}

	var req request.CreateOrUpdatedUserRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	if err := h.service.Update(uint(id), req); err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user updated successfully",
	})
}

func (h *UserHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}
	claims, ok := c.Locals("claims").(*utils.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "no autenticado",
		})
	}

	if claims.Role != "ADMIN" && claims.UserID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "No es tu usuario",
		})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}
