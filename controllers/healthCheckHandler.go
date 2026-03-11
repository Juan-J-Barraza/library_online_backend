package controllers

import "github.com/gofiber/fiber/v3"

func CheckHealth(ctx fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{"message": "Is Ok"})
}
