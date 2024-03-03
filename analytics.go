package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"math/rand"
	"time"
)

func generateRandomSensorData(currentDate, currentTime string) SensorData {
	rand.Seed(time.Now().UnixNano())

	var nodes []NodeData
	for i := 0; i < 5; i++ { // Immer 10 Nodes
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

	for day := 0; day < 30; day++ {
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

// Abfrage 24h
func getSensorDataFromLast24Hours(client *mongo.Client) ([]SensorData, error) {
	collection := client.Database("HoneyMesh").Collection("data")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Berechnen der Daten für die letzten 24 Stunden
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	// Filter für Dokumente des aktuellen und des vorherigen Tages
	filter := bson.M{
		"date": bson.M{
			"$gte": yesterday.Format("2006-01-02"),
		},
	}
	options := options.Find().SetSort(bson.D{{"date", 1}, {"time", 1}})

	var results []SensorData
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var sensorData SensorData
		if err := cursor.Decode(&sensorData); err != nil {
			return nil, err
		}
		// Programmatische Filterung der Daten, um sicherzustellen, dass sie innerhalb der letzten 24 Stunden liegen
		sensorDateTimeStr := sensorData.Date + "T" + sensorData.Time
		sensorDateTime, err := time.Parse("2006-01-02T15:04:05", sensorDateTimeStr)
		if err != nil {
			log.Printf("Fehler beim Parsen des Datums/der Zeit: %v", err)
			continue
		}
		if sensorDateTime.After(yesterday) {
			results = append(results, sensorData)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	fmt.Println(results)
	return results, nil
}

// Dynamische Abfrage, Abfrage durch Stunden
func getFilteredSensorData(client *mongo.Client, hours int) ([]SensorData, error) {
	collection := client.Database("HoneyMesh").Collection("data")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Berechnen der Zeit für die letzten X Stunden, basierend auf dem Parameter
	targetTime := time.Now().Add(time.Duration(-hours) * time.Hour)

	filter := bson.M{
		"date": bson.M{
			"$gte": targetTime.Format("2006-01-02"),
		},
	}
	options := options.Find().SetSort(bson.D{{"date", 1}, {"time", 1}})

	var weights []SensorData
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var sensorData SensorData
		if err := cursor.Decode(&sensorData); err != nil {
			return nil, err
		}
		sensorDateTimeStr := sensorData.Date + "T" + sensorData.Time
		sensorDateTime, err := time.Parse("2006-01-02T15:04:05", sensorDateTimeStr)
		if err != nil {
			log.Printf("Fehler beim Parsen des Datums/der Zeit: %v", err)
			continue
		}
		if sensorDateTime.After(targetTime) {
			weights = append(weights, sensorData)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return weights, nil
}

func getLatestWeightForEachNode(client *mongo.Client, hours int) ([]bson.M, error) {
	collection := client.Database("HoneyMesh").Collection("data")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Berechnen der Zeit für die letzten X Stunden, basierend auf dem Parameter
	targetTime := time.Now().Add(time.Duration(-hours) * time.Hour)

	// Aggregations-Pipeline
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"date": bson.M{"$gte": targetTime.Format("2006-01-02")}}}},
		{{"$unwind", "$nodes"}},
		{{"$sort", bson.D{{"nodes.nodeid", 1}, {"date", -1}, {"time", -1}}}},
		{{"$group", bson.M{
			"_id":          "$nodes.nodeid",
			"latestWeight": bson.M{"$first": "$nodes.weight"},
			"latestDate":   bson.M{"$first": "$date"},
			"latestTime":   bson.M{"$first": "$time"},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}}, // Sortiert die Ergebnisse nach der NodeID in aufsteigender Reihenfolge
	}

	var results []bson.M
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
