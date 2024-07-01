package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Db *mongo.Client

func ConnectToDatabase() error {
	// fmt.Println("Connecting to the database")

	mongoDbUri := os.Getenv("MONGODB")
	fmt.Println("mongo db url: ", mongoDbUri)

	if mongoDbUri == "" {
		return errors.New("mongo db url not found")
	}

	var err error
	Db, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		return err
	}

	fmt.Println("Connected to the database")
	return nil
}

func GetMongoCollection(collName string) (*mongo.Collection, error) {
	fmt.Println("getting ", collName, " collection")
	collection := Db.Database("JwtAuthentication").Collection(collName)
	return collection, nil
}
