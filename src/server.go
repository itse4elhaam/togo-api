package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/itse4elhaam/togo-api.git/src/handlers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectDb(connectionString string) *mongo.Client {
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
	return client
}

func main() {
	// in golang := is used as declaration + assignment operator
	// it will infer types automatically and not let you do the wrong assignment
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := os.Getenv("DB_URL")
	if connectionString == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	dbClient := connectDb(connectionString)
	http.HandleFunc("/api/todos/", func(w http.ResponseWriter, r *http.Request) {
		todoHandler.TodosController(w, r, dbClient, "")
	})
	http.HandleFunc("/api/todos/{userID}", func(w http.ResponseWriter, r *http.Request) {
		todoHandler.TodosController(w, r, dbClient, mux.Vars(r)["userID"])
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}
	fmt.Println("\nThe server is up and running on: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	// routes

}
