package main

import "time"

// "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Gender    string    `json:"gender" bson:"gender"`
	CreatedAt time.Time `json:"created_at"`
}
