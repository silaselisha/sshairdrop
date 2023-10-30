package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/django/v3"
	"github.com/silaselisha/fiber-api/handlers"
	"github.com/silaselisha/fiber-api/middleware"
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
	app.Static("static", "./static")
	app.Use(cors.New())

	handlers.Validate = validator.New()
	handlers.Validate.RegisterValidation("email", util.EmailValidator)

	app.Get("/verify/:token?", store.VerifyAccount)
	api := app.Group("/api")
	v1 := api.Group("/v1", func(ctx *fiber.Ctx) error {
		ctx.Set("version", "v1")
		return ctx.Next()
	})
	v1.Post("/users", store.CreateUser)
	v1.Post("/login", store.Login)
	v1.Get("/users/:id", store.GetUserById)

	// ** SSH SESSION
	var mu sync.Mutex
	handlers.Tunnels = make(map[string]chan handlers.Tunnel)

	go func() {
		ssh.Handle(func(s ssh.Session) {
			token := util.RandTokenGenerator(24)
			fmt.Println("token: -> ", token)
			mu.Lock()
			handlers.Tunnels[token] = make(chan handlers.Tunnel)
			mu.Unlock()

			fmt.Println("token: -> ", token)
			mu.Lock()
			tunnelCh := handlers.Tunnels[token]
			mu.Unlock()
			fmt.Println("Tunnel is ready!")

			tunnel := <-tunnelCh
			_, err := io.Copy(tunnel.W, s)
			if err != nil {
				log.Fatal(err)
			}
			close(tunnel.Done)
			s.Write([]byte("tunneling is done..."))
		})

		log.Fatal(ssh.ListenAndServe(":2222", nil))
	}()

	app.Use(middleware.Protected())
	v1.Get("/sshairdrop/:token?", store.FileShare)

	app.Listen(config.ServerAddress)
}