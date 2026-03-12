package controllers

import (
	"libraryOnline/dtos/filters"
	"libraryOnline/dtos/request"
	"libraryOnline/services"
	"libraryOnline/utils"
	"libraryOnline/utils/validators"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type BookHandler struct {
	service *services.BookService
}

func NewBookHandler(service *services.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetAll(c fiber.Ctx) error {
	p := c.Locals("pagination").(*utils.Pagination)
	var editInt int
	var authorId int
	var err error
	editStr := c.Query("editorial_id")
	if editStr != "" {
		editInt, err = strconv.Atoi(editStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid id",
			})
		}

	}
	authorStr := c.Query("author_id")
	if authorStr != "" {
		authorId, err = strconv.Atoi(authorStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid id",
			})
		}
	}

	f := filters.FiltersBook{
		Title:       c.Query("title"),
		EditorialId: uint(editInt),
		AuthorId:    uint(authorId),
	}

	result, err := h.service.GetAll(f, p)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(result)
}

func (h *BookHandler) FindByID(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	book, err := h.service.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(book)
}

func (h *BookHandler) Create(c fiber.Ctx) error {
	var req request.CreateOrUpdateBookRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	vallidators := validators.ValidateCreateBook(req)
	if vallidators != "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": vallidators})
	}

	book, err := h.service.Create(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(book)
}

func (h *BookHandler) Update(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var req request.CreateOrUpdateBookRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	validator := validators.ValidateUpdateBook(req)
	if validator != "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": validator})
	}

	book, err := h.service.Update(uint(id), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(book)
}

func (h *BookHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
