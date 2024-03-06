package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Funktion zum Herstellen einer Verbindung mit MongoDB
func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Überprüfen der Verbindung
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Erfolgreich mit MongoDB verbunden.")
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
			return true, nil // Collection existiert bereits
		}
	}
	// Collection existiert nicht, also erstellen
	err = db.CreateCollection(context.Background(), "users")
	if err != nil {
		return false, err
	}
	return false, nil // Collection wurde neu erstellt
}
