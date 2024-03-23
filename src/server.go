package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/itse4elhaam/togo-api.git/src/handlers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDb(connectionString string) {
	clientOptions := options.Client().ApplyURI(connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func main() {
	// in golang := is used as declaration + assignment operator
	// it will infer types automatically and not let you do the wrong assignment

	http.HandleFunc("/api/todos", todoHandler.TodosController)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := os.Getenv("DB_URL")
	if connectionString == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	connectDb(connectionString)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}
	fmt.Println("\nThe server is up and running on: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
