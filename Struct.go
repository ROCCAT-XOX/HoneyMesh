package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type SensorData struct {
	Date  string     `bson:"date"`
	Time  string     `bson:"time"`
	Nodes []NodeData `bson:"nodes"`
}

type NodeData struct {
	NodeID       string  `bson:"nodeid"`
	Weight       float64 `bson:"weight"`
	TempOut      float64 `bson:"tempout"`
	TempIn       float64 `bson:"tempin"`
	Humidity     float64 `bson:"humidity"`
	AirQuality   float64 `bson:"airquality"`
	WifiStrength float64 `bson:"wifistrength"`
	Date         string  `bson:"date"`
	Time         string  `bson:"time"`
}

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
