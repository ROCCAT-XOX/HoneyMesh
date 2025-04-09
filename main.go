package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	client, err := connectToMongo()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	err = insertSampleData(client)
	if err != nil {
		fmt.Println("Fehler beim Einfügen von Sample-Daten:", err)
		return
	}

	// Überprüfen, ob die User-Collection existiert und Benutzer enthält
	userColExists, err := ensureUserCollectionExists(client)
	if err != nil {
		log.Fatalf("Fehler bei der Prüfung der User-Collection: %v", err)
	}
	var redirectToFirstLogin bool
	if userColExists {
		userExists, err := userExists(client)
		if err != nil {
			log.Fatalf("Fehler bei der Überprüfung auf vorhandene Benutzer: %v", err)
		}
		redirectToFirstLogin = !userExists // Umleitung zu /first-login, wenn keine Benutzer vorhanden sind
	} else {
		redirectToFirstLogin = true // Umleitung zu /first-login, wenn die Collection neu erstellt wurde
	}

	router := Router(client, redirectToFirstLogin)
	router.Run(":8080")

	//setupPeriodicTask(client)

	// Beispieldaten einfügen

}
