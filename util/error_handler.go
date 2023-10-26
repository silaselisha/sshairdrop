package util

import "github.com/gofiber/fiber/v2"

func ErrorHandler(ctx *fiber.Ctx, status int, err error, message string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"err": err.Error(),
		"message": message,
	})
}
