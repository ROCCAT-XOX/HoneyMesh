package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// fetchData wird erweitert, um den MongoDB-Client als Parameter zu akzeptieren
func fetchData(client *mongo.Client) {
	url := "http://192.168.1.110/"
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Fehler beim Abrufen der Daten:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Fehler beim Lesen des Response Body:", err)
		return
	}

	var data SensorData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Fehler beim Dekodieren der JSON Daten:", err)
		return
	}

	// Speichern der Daten in MongoDB
	collection := client.Database("HoneyMesh").Collection("data")
	_, err = collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Println("Fehler beim Schreiben in die Datenbank:", err)
		return
	}

	log.Println("Daten erfolgreich gespeichert:", data)
}

// setupPeriodicTask wird erweitert, um den MongoDB-Client zu akzeptieren
func setupPeriodicTask(client *mongo.Client) {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			fetchData(client) // Client wird hier übergeben
		}
	}()
}
