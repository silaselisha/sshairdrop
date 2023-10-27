package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/silaselisha/fiber-api/mail"
	"github.com/silaselisha/fiber-api/token"
	"github.com/silaselisha/fiber-api/types"
	"github.com/silaselisha/fiber-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	Login(ctx *fiber.Ctx) error
	CreateUser(ctx *fiber.Ctx) error
	GetUserById(ctx *fiber.Ctx) error
	VerifyAccount(ctx *fiber.Ctx) error
}

type MDBStore struct {
	cl *mongo.Client
	db *mongo.Database
}

var t10n *types.RandToken
var Validate *validator.Validate

func NewStore(cl *mongo.Client, db *mongo.Database) Store {
	return &MDBStore{
		cl: cl,
		db: db,
	}
}

type UserCreateParams struct {
	UserName  string    `json:"username" bson:"username" validate:"required"`
	Email     string    `json:"email" bson:"email" validate:"required,email"`
	Password  string    `json:"password" bson:"password" validate:"required"`
	Gender    string    `json:"gender" bson:"gender" validate:"required"`
	Verified  bool      `json:"verified" bson:"verified"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type UserRequestParams struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Gender    string    `json:"gender" bson:"gender"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func (s *MDBStore) CreateUser(ctx *fiber.Ctx) error {
	data := new(UserCreateParams)

	if err := ctx.BodyParser(data); err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "internal server error")
	}

	if err := Validate.Struct(data); err != nil {
		return util.ErrorHandler(ctx, http.StatusBadRequest, err, "invalid request")
	}

	hashedPassword, err := util.EncryptPassword(data.Password)
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "internal server error")
	}

	data.Password = hashedPassword
	data.Verified = false
	data.CreatedAt = time.Now()

	collection := s.db.Collection("users")
	record, err := collection.InsertOne(ctx.Context(), data)

	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "failed to create new user")
	}

	t10n = types.NewRandToken(util.RandTokenGenerator(8), data.Email)
	verification_link := fmt.Sprintf("http://localhost:3000/verify?token=%v", t10n.Token)

	content, err := mail.ParseMailTemplate(data.UserName, verification_link)
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "internal server error")
	}

	var user types.User
	filter := bson.D{{Key: "_id", Value: record.InsertedID}}
	if err := collection.FindOne(ctx.Context(), filter).Decode(&user); err != nil {
		return util.ErrorHandler(ctx, http.StatusNotFound, err, "user not found")
	}

	config, err := util.Load(".")
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "invalid request")
	}


	jwtmaker, err := token.NewJwtMaker(config.TokenSecretKey)
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "internal server error")
	}

	token, err := jwtmaker.CreateToken(user.Email, 15*time.Minute)
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "internal server error")
	}

	sender := mail.NewGmailSender("ssh file drop", "elishasilas87@gmail.com", config.SenderEmailPassword)
	sender.SendEmail([]string{data.Email}, nil, "Account activation", content)

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func (s *MDBStore) GetUserById(ctx *fiber.Ctx) error {
	params := ctx.Params("id")

	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusBadRequest, err, "invalid request")
	}

	filter := bson.D{{Key: "_id", Value: _id}}
	var user UserRequestParams

	collection := s.db.Collection("users")
	if err := collection.FindOne(ctx.Context(), filter).Decode(&user); err != nil {
		return util.ErrorHandler(ctx, http.StatusBadRequest, err, "invalid request")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

type LoginParams struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func (s *MDBStore) Login(ctx *fiber.Ctx) error {
	var loginCred LoginParams

	if err := ctx.BodyParser(&loginCred); err != nil {
		return util.ErrorHandler(ctx, http.StatusBadRequest, err, "invalid request")
	}

	var user types.User
	filter := bson.D{{Key: "email", Value: loginCred.Email}}
	if err := s.db.Collection("users").FindOne(ctx.Context(), filter).Decode(&user); err != nil {
		return util.ErrorHandler(ctx, http.StatusBadRequest, err, "invalid request")
	}

	if err := util.DecryptPassword(user.Password, loginCred.Password); err != nil {
		return util.ErrorHandler(ctx, http.StatusBadRequest, err, "invalid request")
	}

	config,
		err := util.Load(".")
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "internl server error")
	}

	jwtMaker, err := token.NewJwtMaker(config.TokenSecretKey)
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "internal server error")
	}

	token, err := jwtMaker.CreateToken(user.Email, 15*time.Minute)
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, err, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}

func (s *MDBStore) VerifyAccount(ctx *fiber.Ctx) error {
	token := ctx.Query("token")

	if token != t10n.Token {
		return util.ErrorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid token"), "invalid token")
	}

	time_checker := t10n.IssuedAt.Add(t10n.ExpiresIn)
	if time.Now().After(time_checker) {
		return util.ErrorHandler(ctx, http.StatusBadRequest, fmt.Errorf("invalid token"), "invalid token")
	}

	filter := bson.M{"email": t10n.Email}
	update := bson.M{"$set": bson.D{{Key: "verified", Value: true}}}

	var user *types.User
	err := s.db.Collection("users").FindOneAndUpdate(ctx.Context(), filter, update).Decode(&user)
	if err != nil {
		return util.ErrorHandler(ctx, http.StatusInternalServerError, fmt.Errorf("internal server error"), "internal server error")
	}

	t10n.ExpiresIn = 1 * time.Microsecond
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"user": *user,
	})
}