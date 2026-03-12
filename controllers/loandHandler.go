package controllers

import (
	"libraryOnline/dtos/request"
	"libraryOnline/services"
	"libraryOnline/utils"
	"libraryOnline/utils/validators"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type LoanHandler struct {
	service *services.LoanService
}

func NewLoanHandler(service *services.LoanService) *LoanHandler {
	return &LoanHandler{service: service}
}

func (h *LoanHandler) GetAll(c fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.Claims)
	p := c.Locals("pagination").(*utils.Pagination)

	var result *utils.Pagination
	var err error

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

func (h *LoanHandler) FindByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id inválido"})
	}

	claims := c.Locals("claims").(*utils.Claims)
	loan, err := h.service.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	if claims.Role != "ADMIN" && loan.User.ID != claims.UserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no tienes permisos"})
	}

	return c.JSON(loan)
}

// ConfirmLoan — RESERVED → ACTIVE
func (h *LoanHandler) ConfirmLoan(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id inválido"})
	}

	var req request.ConfirmLoanRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "body inválido"})
	}

	if msg := validators.ValidateConfirmLoan(req); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	loan, err := h.service.ConfirmLoan(uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(loan)
}

// ReturnLoan — ACTIVE → RETURNED
func (h *LoanHandler) ReturnLoan(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id inválido"})
	}

	loan, err := h.service.ReturnLoan(uint(id))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(loan)
}

// CreateDirect — solo ADMIN
func (h *LoanHandler) CreateDirect(c fiber.Ctx) error {
	var req request.CreateDirectLoanRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "body inválido"})
	}

	if msg := validators.ValidateCreateDirectLoan(req); msg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
	}

	loan, err := h.service.CreateDirect(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(loan)
}
