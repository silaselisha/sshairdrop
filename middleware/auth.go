package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/silaselisha/fiber-api/utils"
)

var User struct{} = struct{}{}

func Protected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		authorization, ok := headers["Authorization"]
		if !ok {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"err":     errors.New("unauthorized 1").Error(),
				"message": "unauthorized",
			})
		}

		result := strings.Split(authorization[0], " ")

		if result[0] != "Bearer" {
			fmt.Println(result[0])
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"err":     errors.New("unauthorized 2").Error(),
				"message": "unauthorized",
			})
		}

		maker, err := utils.NewJwtMaker("abcdefghijklmnopqrstuvwxyzabcdefghijklmonpqrsturvwxyz")
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"err":     errors.New("unauthorized 3").Error(),
				"message": "unauthorized",
			})
		}

		payload, err := maker.VerifyToken(result[1])
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"err":     errors.New("unauthorized 4").Error(),
				"message": "unauthorized",
			})
		}

		ctx.Locals(User, *payload)
		return ctx.Next()
	}
}