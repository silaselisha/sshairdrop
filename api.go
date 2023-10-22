package main

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
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

	var user User
	filter := bson.D{{Key: "_id", Value: record.InsertedID}}
	if err := collection.FindOne(ctx.Context(), filter).Decode(&user); err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"err": err.Error(),
			"message": "user not found",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(user)
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

	return ctx.Status(http.StatusOK).JSON(user)
}