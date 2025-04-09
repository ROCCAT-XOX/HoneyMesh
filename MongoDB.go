package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Original function kept for backward compatibility
func connectToMongo() (*mongo.Client, error) {
	return connectToMongoWithURI("mongodb://localhost:27017")
}

// New function that accepts a URI parameter
func connectToMongoWithURI(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB.")
	return client, nil
}

func ensureUserCollectionExists(client *mongo.Client) (bool, error) {
	db := client.Database("HoneyMesh")
	collections, err := db.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		return false, err
	}
	for _, collectionName := range collections {
		if collectionName == "users" {
			return true, nil // Collection already exists
		}
	}
	// Collection doesn't exist, so create it
	err = db.CreateCollection(context.Background(), "users")
	if err != nil {
		return false, err
	}
	return false, nil // Collection was newly created
}
