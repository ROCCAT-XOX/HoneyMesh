package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type NodeData struct {
	NodeID       string  `json:"nodeID"`
	Weight       float64 `json:"weight"`
	TempOut      float64 `json:"tempOut"`
	TempIn       float64 `json:"tempIn"`
	Humidity     float64 `json:"humidity"`
	AirQuality   float64 `json:"airQuality"`
	WifiStrength float64 `json:"wifiStrength"`
	Date         string  `json:"date"`
	Time         string  `json:"time"`
}

type Data struct {
	Nodes []NodeData `json:"nodes"`
}

type SensorData struct {
	Date string `json:"date"`
	Time string `json:"time"`
	Data Data   `json:"data"`
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
