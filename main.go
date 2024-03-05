package main

import (
	"context"
	"log"
)

func main() {
	client, err := connectToMongo()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	//setupPeriodicTask(client)

	// Beispieldaten einfügen
	/*
		err = insertSampleData(client)
		if err != nil {
			fmt.Println("Fehler beim Einfügen von Sample-Daten:", err)
			return
		}

	*/

	// Erstellen eines neuen Nutzers
	/*
		err = createUser(client, "admin@admin.de", "admin")
		if err != nil {
			log.Fatal(err)
		}
	*/

	router := Router(client)

	router.Run(":8080")

}
