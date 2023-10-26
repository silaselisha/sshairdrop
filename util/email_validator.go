package util

import (
	"context"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/silaselisha/fiber-api/types"
	"go.mongodb.org/mongo-driver/bson"
)

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)

	truthy := regex.MatchString(email)
	if truthy {
		var user types.User

		filter := bson.M{"email": email}
		if err := db.Collection("users").FindOne(context.Background(), filter).Decode(&user); err != nil {
			return true
		}
	}
	return false
}

func EmailValidator(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	return validateEmail(email)
}
