package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

// Diese Funktion filtert SensorDaten basierend auf einem Zeitfenster und sortiert sie nach NodeID und Zeit.
func getFilteredSensorData(client *mongo.Client, hours int) ([]SensorData, error) {
	collection := client.Database("HoneyMesh").Collection("data")

	// Berechne das Startdatum/-zeit für die Filterung
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	// Filter für Daten, die nach dem Startzeitpunkt erstellt wurden
	filter := bson.M{
		"date": startTime.Format("2006-01-02"),               // Nur das Datum ohne Uhrzeit
		"time": bson.M{"$gte": startTime.Format("15:04:05")}, // Uhrzeit ab dem Startzeitpunkt
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []SensorData
	for cursor.Next(context.Background()) {
		var result SensorData
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
