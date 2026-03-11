package middleware

import (
	"libraryOnline/utils"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func PaginationMiddleware(c fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	pagination := &utils.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	c.Locals("pagination", pagination)

	return c.Next()
}
