package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"strings"
)

func Router(client *mongo.Client) *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/assets", "./assets")

	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			c.Header("Cache-Control", "public, max-age=86400")
		}
	})

	router.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})

	router.Use(loadSensorData(client))

	router.GET("/", func(c *gin.Context) {
		analyzedData, err := analyzeSensorData(client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler bei der Datenanalyse"})
			return
		}

		jsonData, err := json.Marshal(analyzedData)
		if err != nil {
			// Fehlerbehandlung
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Konvertieren der Daten"})
			return
		}

		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":      "HoneyMesh",
			"sensorData": template.JS(jsonData), // Konvertierte JSON-Daten
		})
	})

	router.POST("/submit-data", func(c *gin.Context) {
		var sensorData SensorData
		if err := c.BindJSON(&sensorData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		collection := client.Database("HoneyMesh").Collection("data")
		_, err := collection.InsertOne(context.Background(), sensorData) // Verwenden Sie context.Background()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fehler beim Schreiben in die Datenbank"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Daten erfolgreich gespeichert"})
	})

	return router
}
