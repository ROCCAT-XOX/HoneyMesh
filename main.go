package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

func loadSensorData(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sensorData []SensorData
		collection := client.Database("HoneyMesh").Collection("data")

		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Abfragen der Datenbank"})
			c.Abort()
			return
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var elem SensorData
			err := cursor.Decode(&elem)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Dekodieren der Daten"})
				c.Abort()
				return
			}
			sensorData = append(sensorData, elem)
		}
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler bei der Cursor-Verarbeitung"})
			c.Abort()
			return
		}

		// Speichere die Daten im Kontext von Gin
		c.Set("sensorData", sensorData)
		c.Next()
	}
}

func main() {
	client, err := connectToMongo()
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	// Beispieldaten einfügen

	if err := insertSampleData(client); err != nil {
		log.Fatal("Fehler beim Einfügen von Beispieldaten: ", err)
	}

	router := Router(client)

	router.Run(":8080")
}
