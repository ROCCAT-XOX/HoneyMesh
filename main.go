package main

import (
	"context"
	"log"
	"os"
	"strings"
)

func main() {
	// Use environment variable for MongoDB URI if set
	mongoURI := "mongodb://localhost:27017"
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		mongoURI = uri
	}

	// Connect to MongoDB using the URI
	client, err := connectToMongoWithURI(mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Check if we should generate sample data
	if genData := os.Getenv("GENERATE_SAMPLE_DATA"); strings.ToLower(genData) == "true" {
		log.Println("Generating sample data...")
		err = insertSampleData(client)
		if err != nil {
			log.Printf("Error inserting sample data: %v", err)
		} else {
			log.Println("Sample data successfully generated")
		}
	}

	// Check if the User-Collection exists and contains users
	userColExists, err := ensureUserCollectionExists(client)
	if err != nil {
		log.Fatalf("Error checking user collection: %v", err)
	}
	var redirectToFirstLogin bool
	if userColExists {
		userExists, err := userExists(client)
		if err != nil {
			log.Fatalf("Error checking for existing users: %v", err)
		}
		redirectToFirstLogin = !userExists // Redirect to /first-login if no users exist
	} else {
		redirectToFirstLogin = true // Redirect to /first-login if the collection was newly created
	}

	// Set up the router and start the web server
	router := Router(client, redirectToFirstLogin)
	router.Run(":8080")
}
