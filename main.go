package main

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/django/v3"
	"github.com/silaselisha/fiber-api/handlers"
	"github.com/silaselisha/fiber-api/util"
)

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
	store := handlers.NewStore(client, database)

	engine := django.New("./templates", ".django")
	
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("static","./static")
	app.Use(cors.New())
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{
			"Title": "fiber api",
		}, "layouts/main")
	})

	handlers.Validate = validator.New()
	handlers.Validate.RegisterValidation("email", util.EmailValidator)

	api := app.Group("/api")
	v1 := api.Group("/v1", func(ctx *fiber.Ctx) error {
		ctx.Set("version", "v1")
		return ctx.Next()
	})
	v1.Post("/users", store.CreateUser)
	v1.Post("/login", store.Login)
	v1.Get("/users/:id", store.GetUserById)

	app.Listen(config.ServerAddress)
}
