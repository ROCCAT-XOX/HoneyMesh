package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Globale Variable für Demonstration; In Produktionscode durch Datenbankspeicherung ersetzen
var globalSessionToken string

func GenerateSecureToken(length int) (string, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		// Ein Fehler trat während der Generierung auf
		return "", err
	}

	// Konvertiert den Byte-Slice in eine Base64-kodierte Zeichenkette
	return base64.URLEncoding.EncodeToString(token), nil
}

func isUserAuthenticated(c *gin.Context) bool {
	sessionToken, err := c.Cookie("session_token")
	return err == nil && sessionToken == globalSessionToken
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isUserAuthenticated(c) {
			// Für serverseitig gerenderte Seiten: Umleitung zur Login-Seite
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
