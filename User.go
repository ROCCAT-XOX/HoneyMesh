package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func createUser(client *mongo.Client, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	collection := client.Database("HoneyMesh").Collection("users")
	_, err = collection.InsertOne(context.Background(), bson.M{"username": username, "password": string(hashedPassword)})
	if err != nil {
		return err
	}

	fmt.Println("Nutzer erfolgreich erstellt")
	return nil
}

func userExists(client *mongo.Client) (bool, error) {
	collection := client.Database("HoneyMesh").Collection("users")
	count, err := collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func getAllUsers(client *mongo.Client) ([]User, error) {
	var users []User
	collection := client.Database("HoneyMesh").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
