package todoHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	m "github.com/itse4elhaam/togo-api.git/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// todo create a mapper to simplify this logic
func TodosController(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client, todoId string) {
	if r.Method == "POST" {
		createTodo(w, r, dbClient)
		return
	} else if r.Method == "GET" {
		getTodos(w, r, dbClient)
		return
	} else if r.Method == "PATCH" {
		deleteTodo(w, r, dbClient, todoId)

	} else if r.Method == "DELETE" {
		updateTodo(w, r, dbClient, todoId)

	} else {
		fmt.Println("method not supported", r.Method)
		return
	}
}

func getTodos(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	collection := dbClient.Database("togoDb").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{}) // Find all (empty filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx) // Important: Close cursor when done

	var todos []m.Todo // Slice to store results
	for cursor.Next(ctx) {
		var todo m.Todo
		if err := cursor.Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		// Handle encoding errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createTodo(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {

	var t m.Todo

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if t.Title == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// why not just get the collection in props instead ?
	collection := dbClient.Database("togoDb").Collection("todos")
	insertResult, err := collection.InsertOne(context.TODO(), t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct{ ID interface{} }{insertResult.InsertedID})
	// insert a new todo object into the database
	fmt.Println(t)
}

func updateTodo(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client, todoId string) {
	fmt.Println(todoId)

}

func deleteTodo(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client, todoId string) {
	objectID, err := primitive.ObjectIDFromHex(todoId) 
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	collection := dbClient.Database("yourDatabaseName").Collection("todoCollection")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Example: DeleteOne
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted documents: %v\n", result.DeletedCount)

}
