package middleware

import (
	"libraryOnline/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AuthRequired() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "el token es requerido",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token inválido o expirado",
			})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

func RoleMiddleware(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*utils.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "no autenticado",
			})
		}

		for _, role := range roles {
			if claims.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "no tienes permisos para acceder a este recurso",
		})
	}
}
