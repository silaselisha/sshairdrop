package util

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database

func InitDB() (*mongo.Client, *mongo.Database, error) {
	config, err := Load("./..")
	if err != nil {
		return nil, nil, err
	}
	ctx, clear := context.WithTimeout(context.Background(), 10*time.Second)
	defer clear()

	cl, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DbUri))
	if err != nil {
		return nil, nil, err
	}

	if err = cl.Ping(ctx, &readpref.ReadPref{}); err != nil {
		return nil, nil, err
	}

	db = cl.Database(config.DbName)
	return cl, db, nil
}
