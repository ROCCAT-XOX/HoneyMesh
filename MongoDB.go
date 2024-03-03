package main

import (
	"context"
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
