package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/silaselisha/fiber-api/types"
)

const (
	SECRET_KEY_SIZE = 32
)

type JwtMaker struct {
	SecretKey string `json:"secret_key"`
}

func NewJwtMaker(secretkey string) (Maker, error) {
	if len(secretkey) < SECRET_KEY_SIZE {
		return nil, fmt.Errorf("invalid secret key %+v", secretkey)
	}

	return &JwtMaker{
		SecretKey: secretkey,
	}, nil
}

func (jm *JwtMaker) CreateToken(email string, duration time.Duration) (string, error) {
	payload := types.NewPayload(email, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	result, err := token.SignedString([]byte(jm.SecretKey))
	if err != nil {
		return "", err
	}

	return result, nil
}

func (jm *JwtMaker) VerifyToken(token string) (*types.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %+v", token.Header["alg"])
		}
		return []byte(jm.SecretKey), nil
	}

	t, err := jwt.ParseWithClaims(token, &types.Payload{}, keyFunc)
	if err != nil {
		return nil, err
	}

	payload, ok := t.Claims.(*types.Payload)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return payload, nil
}