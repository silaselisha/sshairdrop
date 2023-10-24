package main

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/silaselisha/fiber-api/middleware"
	"github.com/silaselisha/fiber-api/util"
)

var validate *validator.Validate
func main() {
	config, err := util.Load(".")
	if err != nil {
		log.Fatal(err)
	}

	client, database, err := util.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	defer database.Client().Disconnect(context.Background())
	store := NewStore(client, database)

	app := fiber.New()

	validate = validator.New()
	validate.RegisterValidation("email", util.EmailValidator)
	

	api := app.Group("/api")
	v1 := api.Group("/v1", func(ctx *fiber.Ctx) error {
		ctx.Set("version", "v1")
		return ctx.Next()
	})
	v1.Post("/users", store.createUser)
	v1.Post("/login", store.login)
	v1.Get("/users/:id", store.getUserById)

	v1.Get("/products", middleware.Protected(), store.getProducts)

	app.Listen(config.ServerAddress)
}