package types

import (
	"fmt"
	"time"
)

type User struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Gender    string    `json:"gender" bson:"gender"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Payload struct {
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresIn time.Time `json:"expires_in"`
}

func NewPayload(email string, duration time.Duration) *Payload {
	return &Payload{
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiresIn: time.Now().Add(duration),
	}
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresIn) {
		return fmt.Errorf("token expires at: %+v", p.ExpiresIn)
	}
	return nil
}