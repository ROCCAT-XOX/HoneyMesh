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

	router := Router(client)

	router.Run(":8080")

}
