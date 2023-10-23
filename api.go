package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/silaselisha/fiber-api/middleware"
	"github.com/silaselisha/fiber-api/types"
	"github.com/silaselisha/fiber-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	cl *mongo.Client
	db *mongo.Database
}

func NewStore(cl *mongo.Client, db *mongo.Database) *Store {
	return &Store{
		cl: cl,
		db: db,
	}
}

type UserCreateParams struct {
	UserName  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	Gender    string    `json:"gender" bson:"gender"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type UserRequestParams struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Gender    string    `json:"gender" bson:"gender"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func (s *Store) createUser(ctx *fiber.Ctx) error {
	data := new(UserCreateParams)

	if err := ctx.BodyParser(data); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err":     err.Error(),
			"message": "invalid request",
		})
	}

	hashedPassword, err := utils.EncryptPassword(data.Password)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err":     err.Error(),
			"message": "internal server error",
		})
	}
	data.Password = hashedPassword
	data.CreatedAt = time.Now()

	collection := s.db.Collection("users")
	record, err := collection.InsertOne(ctx.Context(), data)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err":     err.Error(),
			"message": "failed to create a user",
		})
	}

	var user types.User
	filter := bson.D{{Key: "_id", Value: record.InsertedID}}
	if err := collection.FindOne(ctx.Context(), filter).Decode(&user); err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"err":     err.Error(),
			"message": "user not found",
		})
	}

	jwtmaker, err := utils.NewJwtMaker("abcdefghijklmnopqrstuvwxyzabcdefghijklmonpqrsturvwxyz")
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err":     err.Error(),
			"message": "failed to create a user",
		})
	}

	token, err := jwtmaker.CreateToken(user.Email, 15*time.Minute)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err":     err.Error(),
			"message": "failed to create a user",
		})
	}

	fmt.Println(token)
	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func (s *Store) getUserById(ctx *fiber.Ctx) error {
	params := ctx.Params("id")

	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err":     err.Error(),
			"message": "invalid user id",
		})
	}

	filter := bson.D{{Key: "_id", Value: _id}}
	var user UserRequestParams

	collection := s.db.Collection("users")
	if err := collection.FindOne(ctx.Context(), filter).Decode(&user); err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"err":     err.Error(),
			"message": "user not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

type LoginParams struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (s *Store) login(ctx *fiber.Ctx) error {
	var loginCred LoginParams

	if err := ctx.BodyParser(&loginCred); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err":     err.Error(),
			"message": "bad request",
		})
	}

	var user types.User
	filter := bson.D{{Key: "email", Value: loginCred.Email}}
	if err := s.db.Collection("users").FindOne(ctx.Context(), filter).Decode(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
			"message": "inavlid request",
		})
	}
	
	if err := utils.DecryptPassword(user.Password, loginCred.Password); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err": err.Error(),
			"message": "invalid request",
		})
	}

	jwtMaker, err := utils.NewJwtMaker("abcdefghijklmnopqrstuvwxyzabcdefghijklmonpqrsturvwxyz")
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
			"message": "internal server error",
		})
	}
	token, err := jwtMaker.CreateToken(user.Email, 15 * time.Minute)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
			"message": "internal server error",
		})
		
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func(s *Store) getProducts(ctx *fiber.Ctx) error {
	user := ctx.Locals(middleware.User).(types.Payload)
	fmt.Println(user.Email)
	return nil
}


// ** error handler D_R_Y
// func errorHandler(ctx *fiber.Ctx, status int, err error, message string) error {
// 	return ctx.Status(status).JSON(fiber.Map{
// 		"err": err.Error(),
// 		"message": message,
// 	})
// }