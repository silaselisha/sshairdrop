package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/silaselisha/fiber-api/middleware"
	"github.com/silaselisha/fiber-api/utils"
)

const (
	SERVER_ADDRESS = ":8080"
)

func main() {
	client, database, err := utils.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	defer database.Client().Disconnect(context.Background())
	store := NewStore(client, database)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())

	app.Post("/users", store.createUser)
	app.Post("/login", store.login)
	app.Get("/users/:id", store.getUserById)

	app.Get("/products", middleware.Protected(), store.getProducts)

	app.Listen(SERVER_ADDRESS)
}
