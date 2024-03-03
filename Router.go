package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"strconv"
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

	// Route-Handler für verschiedene Zeitfenster
	router.GET("/data", func(c *gin.Context) {
		// Abfrageparameter "hours" aus der URL abrufen
		hoursStr := c.Query("hours")

		// Standardwert für die Stunden festlegen, falls nicht angegeben
		hours := 24 // Standardwert

		// Versuchen, die Stunden aus dem Abfrageparameter zu parsen
		if hoursStr != "" {
			parsedHours, err := strconv.Atoi(hoursStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for 'hours' parameter"})
				return
			}
			hours = parsedHours
		}

		// Daten mit der entsprechenden Anzahl von Stunden abrufen
		data, err := getFilteredSensorData(client, hours)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	router.GET("/test", func(c *gin.Context) {

		weights, err := getLatestWeightForEachNode(client, 24) // Beispiel: Daten der letzten 24 Stunden
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "test.html", gin.H{
			"title":   "Meine Testseite",
			"weights": weights,
		})
	})

	router.GET("/login", func(c *gin.Context) {

		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
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
