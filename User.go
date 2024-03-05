package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func createUser(client *mongo.Client, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	collection := client.Database("HoneyMesh").Collection("user")
	_, err = collection.InsertOne(context.Background(), bson.M{"username": username, "password": string(hashedPassword)})
	if err != nil {
		return err
	}

	fmt.Println("Nutzer erfolgreich erstellt")
	return nil
}
