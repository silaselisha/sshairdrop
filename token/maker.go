package token

import (
	"time"

	"github.com/silaselisha/fiber-api/types"
)

type Maker interface {
	CreateToken(email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*types.Payload, error)
}
