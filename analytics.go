package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"
)

func generateRandomSensorData() SensorData {
	rand.Seed(time.Now().UnixNano())
	nodeCount := rand.Intn(10) + 1 // 1 bis 10 Nodes

	var nodes []NodeData
	for i := 0; i < nodeCount; i++ {
		node := NodeData{
			NodeID:       fmt.Sprintf("%d", i+1),
			Weight:       rand.Float64() * 100,   // 0 bis 100
			TempOut:      rand.Float64()*50 - 10, // -10 bis 40
			TempIn:       rand.Float64() * 50,    // 0 bis 50
			Humidity:     rand.Float64()*40 + 20, // 20 bis 60
			AirQuality:   rand.Float64() * 500,   // 0 bis 500
			WifiStrength: rand.Float64()*80 + 20, // 20 bis 100
		}
		nodes = append(nodes, node)
	}

	return SensorData{
		Date: time.Now().Format("2006-01-02"),
		Time: time.Now().Format("15:04:05"),
		Data: Data{Nodes: nodes},
	}
}

func insertSampleData(client *mongo.Client) error {
	collection := client.Database("HoneyMesh").Collection("data")
	now := time.Now()

	for i := 0; i < 90; i++ { // Für die letzten 90 Tage
		for j := 0; j < 24; j++ { // Für jede Stunde eines Tages
			sensorData := generateRandomSensorData()
			sensorData.Date = now.AddDate(0, 0, -i).Format("2006-01-02") // Datum setzen
			sensorData.Time = fmt.Sprintf("%02d:00:00", j)               // Stunde setzen

			_, err := collection.InsertOne(context.Background(), sensorData)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type AnalyzedData struct {
	LastHour    []SensorData
	Last24Hours []SensorData
	Last7Days   []SensorData
	Last30Days  []SensorData
	Last90Days  []SensorData
}

func analyzeSensorData(client *mongo.Client) (*AnalyzedData, error) {
	collection := client.Database("HoneyMesh").Collection("data")
	analyzed := &AnalyzedData{}

	// Hilfsfunktion, um Daten für ein bestimmtes Intervall abzufragen
	fetchDataForInterval := func(duration time.Duration) ([]SensorData, error) {
		var results []SensorData
		start := time.Now().Add(-duration)
		filter := bson.M{"date": bson.M{"$gte": start.Format("2006-01-02")}}

		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var elem SensorData
			if err := cursor.Decode(&elem); err != nil {
				return nil, err
			}
			results = append(results, elem)
		}

		return results, nil
	}

	// Abrufen der Daten für jedes Intervall
	var err error
	if analyzed.LastHour, err = fetchDataForInterval(1 * time.Hour); err != nil {
		return nil, err
	}
	if analyzed.Last24Hours, err = fetchDataForInterval(24 * time.Hour); err != nil {
		return nil, err
	}
	if analyzed.Last7Days, err = fetchDataForInterval(7 * 24 * time.Hour); err != nil {
		return nil, err
	}
	if analyzed.Last30Days, err = fetchDataForInterval(30 * 24 * time.Hour); err != nil {
		return nil, err
	}
	if analyzed.Last90Days, err = fetchDataForInterval(90 * 24 * time.Hour); err != nil {
		return nil, err
	}

	return analyzed, nil
}
