package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Tunnel struct {
	W    io.Writer
	Done chan struct{}
}

var Tunnels map[string]chan Tunnel

func (s *MDBStore) FileShare(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	tun, ok := Tunnels[token]

	fmt.Println(token)

	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err":     fmt.Errorf("ssh session id not found: -> %v", token),
			"message": "internal server error",
		})
	}

	done := make(chan struct{})
	tun <- Tunnel{
		W:    os.Stdout,
		Done: done,
	}

	<-done
	return ctx.Status(http.StatusOK).JSON(fiber.Map{})
}