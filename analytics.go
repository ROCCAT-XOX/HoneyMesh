package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"

	"math/rand"
	"time"
)

func generateRandomSensorData(currentDate, currentTime string) SensorData {
	rand.Seed(time.Now().UnixNano())

	var nodes []NodeData
	for i := 0; i < 10; i++ { // Immer 10 Nodes
		node := NodeData{
			NodeID:       fmt.Sprintf("%d", i+1),
			Weight:       rand.Float64() * 100,
			TempOut:      rand.Float64()*50 - 10,
			TempIn:       rand.Float64() * 50,
			Humidity:     rand.Float64()*40 + 20,
			AirQuality:   rand.Float64() * 500,
			WifiStrength: rand.Float64()*80 + 20,
			Date:         currentDate,
			Time:         fmt.Sprintf("%sT%s", currentDate, currentTime[:8]), // Nur die ersten 8 Zeichen von currentTime (HH:MM:SS)
		}
		nodes = append(nodes, node)
	}

	return SensorData{
		Date:  currentDate,
		Time:  currentTime[:8], // Nur die ersten 8 Zeichen von currentTime (HH:MM:SS)
		Nodes: nodes,
	}
}

func insertSampleData(client *mongo.Client) error {
	collection := client.Database("HoneyMesh").Collection("data")

	for day := 0; day < 14; day++ {
		for hour := 0; hour < 24; hour++ {
			dataTime := time.Now().Add(time.Duration(-day) * 24 * time.Hour).Add(time.Duration(-hour) * time.Hour)
			currentDate := dataTime.Format("2006-01-02")
			currentTime := dataTime.Format("15:04:05") + fmt.Sprintf("%09d", time.Now().UnixNano()%1e9)

			sensorData := generateRandomSensorData(currentDate, currentTime)

			_, err := collection.InsertOne(context.Background(), sensorData)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
