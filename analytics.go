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

	currentDate := time.Now().Format("2006-01-02")
	currentTime := time.Now().Format("15:04:05")

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
			Date:         currentDate,            // Datum hinzugefügt
			Time:         currentTime,            // Uhrzeit hinzugefügt
		}
		nodes = append(nodes, node)
	}

	return SensorData{
		Date: currentDate,
		Time: currentTime,
		Data: Data{Nodes: nodes},
	}
}

func insertSampleData(client *mongo.Client) error {
	collection := client.Database("HoneyMesh").Collection("data")

	for i := 0; i < 90; i++ { // Für die letzten 90 Tage
		currentDate := time.Now().AddDate(0, 0, -i).Format("2006-01-02") // Datum für den jeweiligen Tag setzen

		for j := 0; j < 24; j++ { // Für jede Stunde eines Tages
			sensorData := generateRandomSensorData()
			sensorData.Date = currentDate                  // Datum einsetzen
			sensorData.Time = fmt.Sprintf("%02d:00:00", j) // Stunde einsetzen

			// Aktualisieren Sie das Datum und die Zeit für jeden Node
			for k := range sensorData.Data.Nodes {
				sensorData.Data.Nodes[k].Date = currentDate
				sensorData.Data.Nodes[k].Time = sensorData.Time
			}

			_, err := collection.InsertOne(context.Background(), sensorData)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type NodeAnalyzedData struct {
	NodeID      string
	LastHour    []NodeData
	Last24Hours []NodeData
	Last7Days   []NodeData
	Last30Days  []NodeData
	Last90Days  []NodeData
}

type AllNodesAnalyzedData struct {
	DataByNode map[string]*NodeAnalyzedData
}

func analyzeSensorData(client *mongo.Client) (*AllNodesAnalyzedData, error) {
	collection := client.Database("HoneyMesh").Collection("data")
	allNodesData := &AllNodesAnalyzedData{DataByNode: make(map[string]*NodeAnalyzedData)}

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

	intervals := []time.Duration{1 * time.Hour, 24 * time.Hour, 7 * 24 * time.Hour, 30 * 24 * time.Hour, 90 * 24 * time.Hour}
	for _, interval := range intervals {
		sensorData, err := fetchDataForInterval(interval)
		if err != nil {
			return nil, err
		}

		for _, data := range sensorData {
			for _, nodeData := range data.Data.Nodes {
				nodeAnalysis := allNodesData.DataByNode[nodeData.NodeID]
				if nodeAnalysis == nil {
					nodeAnalysis = &NodeAnalyzedData{NodeID: nodeData.NodeID}
					allNodesData.DataByNode[nodeData.NodeID] = nodeAnalysis
				}

				switch interval {
				case 1 * time.Hour:
					nodeAnalysis.LastHour = append(nodeAnalysis.LastHour, nodeData)
				case 24 * time.Hour:
					nodeAnalysis.Last24Hours = append(nodeAnalysis.Last24Hours, nodeData)
				case 7 * 24 * time.Hour:
					nodeAnalysis.Last7Days = append(nodeAnalysis.Last7Days, nodeData)
				case 30 * 24 * time.Hour:
					nodeAnalysis.Last30Days = append(nodeAnalysis.Last30Days, nodeData)
				case 90 * 24 * time.Hour:
					nodeAnalysis.Last90Days = append(nodeAnalysis.Last90Days, nodeData)
				}
			}
		}
	}

	return allNodesData, nil
}
