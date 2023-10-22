package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func InitDB() (*mongo.Client, *mongo.Database, error) {
	DB_NAME := "trivia"
	DB_URI := "mongodb://localhost:27017"
	ctx, clear := context.WithTimeout(context.Background(), 10 * time.Second)
	defer clear()

	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Print(err)
		return nil, nil, err
	}

	db = client.Database(DB_NAME)
	return client, db, nil
}