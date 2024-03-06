package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Router(client *mongo.Client, redirectToFirstLogin bool) *gin.Engine {
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
	// 404 Error Handler
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "Seite nicht gefunden",
		})
	})
	// API
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
	// User-Control
	router.GET("/user", AuthRequired(), func(c *gin.Context) {
		users, err := getAllUsers(client) // `client` ist dein MongoDB-Client
		if err != nil {
			// Behandle den Fehler angemessen
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"message": "Fehler beim Laden der Nutzer",
			})
			return
		}
		message := c.Query("message") // Erhält den Query-Parameter "message"
		c.HTML(http.StatusOK, "user.html", gin.H{
			"title":   "User",
			"message": message,
			"users":   users, // Übergebe die Nutzerliste an das Template
		})
	})

	router.POST("/create-user", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		confirmPassword := c.PostForm("confirmPassword")
		origin := c.PostForm("origin")

		// Überprüfen, ob die Passwörter übereinstimmen
		if password != confirmPassword {
			message := "Die Passwörter stimmen nicht überein."
			// Weiterleitung mit Fehlermeldung
			redirectURL := "/user"
			if origin == "first-login" {
				redirectURL = "/first-login"
			}
			c.Redirect(http.StatusSeeOther, redirectURL+"?error="+url.QueryEscape(message))
			return
		}

		// Versuch, den Nutzer zu erstellen
		err := createUser(client, username, password)
		if err != nil {
			fmt.Println("Fehler beim Erstellen des Nutzers:", err)
			// Weiterleitung mit Fehlermeldung
			redirectURL := "/user"
			if origin == "first-login" {
				redirectURL = "/first-login"
			}
			c.Redirect(http.StatusSeeOther, redirectURL+"?error="+url.QueryEscape("Fehler beim Erstellen des Nutzers"))
			return
		}

		// Weiterleitung nach dem Erstellen des Nutzers
		if origin == "first-login" {
			// Umleitung zu /login nach der Erstellung des ersten Nutzers
			c.Redirect(http.StatusSeeOther, "/login")
		} else {
			// Umleitung zu /user mit einer Erfolgsmeldung
			c.Redirect(http.StatusSeeOther, "/user?message="+url.QueryEscape("User erfolgreich erstellt"))
		}
	})
	router.GET("/first-login", func(c *gin.Context) {
		exists, err := userExists(client)
		if err != nil {
			// Fehlerbehandlung, z.B. Logging und Anzeigen eines Fehlers
			log.Printf("Fehler beim Überprüfen der Benutzerexistenz: %v", err)
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Interner Serverfehler"})
			return
		}
		if exists {
			// Wenn bereits Benutzer existieren, Umleitung zu /login oder einer anderen Seite
			c.Redirect(http.StatusFound, "/login")
		} else {
			// Ansonsten die Seite first-login anzeigen
			c.HTML(http.StatusOK, "first-login.html", gin.H{"title": "Erste Anmeldung"})
		}
	})
	// Normal Route
	router.GET("/", func(c *gin.Context) {
		// Überprüfen, ob Benutzer in der Datenbank existieren
		userExists, err := userExists(client)
		if err != nil {
			// Fehlerbehandlung, z.B. Loggen des Fehlers und Senden einer Fehlermeldung
			log.Printf("Fehler bei der Überprüfung auf vorhandene Benutzer: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Interner Serverfehler"})
			return
		}

		if !userExists {
			// Keine Benutzer vorhanden, umleiten auf /first-login
			c.Redirect(http.StatusTemporaryRedirect, "/first-login")
		} else if isUserAuthenticated(c) {
			// Benutzer ist authentifiziert und es existieren Benutzer, umleiten auf /home
			c.Redirect(http.StatusTemporaryRedirect, "/home")
		} else {
			// Benutzer ist nicht authentifiziert, umleiten auf /login
			c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
	})
	router.GET("/dashboard", AuthRequired(), func(c *gin.Context) {
		weights, err := getLatestWeightForEachNode(client) // Beispiel: Daten der letzten 24 Stunden
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title":   "Dashboard",
			"weights": weights,
		})
	})
	router.GET("/home", AuthRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"title": "Home",
		})
	})
	// Login
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
		})
	})
	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("user")
		password := c.PostForm("password")

		var user User
		collection := client.Database("HoneyMesh").Collection("users")
		err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			// Hier setzen wir die Fehlermeldung und rendern die Login-Seite erneut
			c.HTML(http.StatusOK, "login.html", gin.H{
				"title": "Login",
				"error": "Benutzername oder Passwort falsch",
			})
			return
		}

		sessionToken, err := GenerateSecureToken(32)
		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"title": "Login",
				"error": "Fehler bei der Token-Generierung",
			})
			return
		}

		globalSessionToken = sessionToken
		maxAge := 3600
		c.SetCookie("session_token", sessionToken, maxAge, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/home")
	})

	// Logout
	router.GET("/logout", AuthRequired(), func(c *gin.Context) {
		// Löschen des Session-Cookies, Anpassung für lokale Entwicklung
		c.SetCookie("session_token", "", -1, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/login")
	})
	// Add Data to DB over ESP32S
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
