package main

import (
	"context"
	"log"
	"mongodb-golang/controllers"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Entry point of the application
func main() {
	r := httprouter.New()

	// Initialize the user controller with a MongoDB client
	uc := controllers.NewUserController(getMongoClient())

	r.GET("/", uc.GetUser)
	r.POST("/user", uc.CreateUser)       // Use uc methods for user operations
	r.DELETE("/user/:id", uc.DeleteUser) // Use uc methods for deleting users

	log.Fatal(http.ListenAndServe(":8080", r)) // Start the server
	log.Println("Server is running on port 8080")
}

// Function to get a MongoDB client
func getMongoClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
