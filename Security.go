package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Diese Funktion ist nur ein Platzhalter. In einer echten Anwendung würden Sie hier eine tatsächliche Authentifizierungsprüfung durchführen.
func isUserAuthenticated(c *gin.Context) bool {
	// Beispiel: Prüfen, ob ein Session-Cookie vorhanden ist
	// _, err := c.Cookie("session_token")
	// return err == nil

	// Für dieses Beispiel geben wir einfach true zurück, um anzudeuten, dass der Benutzer authentifiziert ist.
	// Ersetzen Sie dies durch Ihre eigene Logik.
	return true
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
