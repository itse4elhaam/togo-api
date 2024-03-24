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

// some general notes : 
/*
- bson is used alot because binary json is the standard format for mongodb data base, hence in order to do CRUD we often need to convert info into that format
- functions with capital first character is public and can be used in other files but the lowercase ones are private to the they're in
- in order to convert anything into json we use this function(err := json.NewDecoder(r.Body).Decode(&t)), it takes a reference to object we want to store it in 

*/

func TodosController(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client, todoId string) {
	if r.Method == "POST" {
		createTodo(w, r, dbClient)
		return
	} else if r.Method == "GET" {
		getTodos(w, r, dbClient)
		return
	} else if r.Method == "PATCH" {
		updateTodo(w, r, dbClient, todoId)
	} else if r.Method == "DELETE" {
		deleteTodo(w, r, dbClient, todoId)
	} else {
		log.Fatal("method not supported", r.Method)
		return
	}
}

func getTodos(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client) {
	collection := dbClient.Database("togoDb").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// we get the data in a cursor because the data could be large, this is a good way to optimize that
	cursor, err := collection.Find(ctx, bson.D{}) // Find all (empty filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx) // Important: Close cursor when done

	// looping over the cursor and decoding it to store it in the todos array
	var todos []m.Todo 
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
}

func updateTodo(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client, todoId string) {
	var t m.Todo
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// &t.Completed is used because that is the only way to check if it is null or not
	if t.Title == "" && &t.Completed == nil {
		http.Error(w, "Nothing to update", http.StatusBadRequest)
		return
	}

	collection := dbClient.Database("togoDb").Collection("todos")
	id, _ := primitive.ObjectIDFromHex(todoId)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.D{}
	if t.Title != "" {
		update = append(update, bson.E{Key: "title", Value: t.Title})
	}
	if &t.Completed != nil {
		update = append(update, bson.E{Key: "completed", Value: t.Completed})
	}

	if len(update) == 0 {
		return
	}

	update = bson.D{{Key: "$set", Value: update}}
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if result.ModifiedCount == 0{
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	fmt.Println("Updated Documents: ", result.ModifiedCount)
	var foundTodo m.Todo
	err = collection.FindOne(context.TODO(), filter).Decode(&foundTodo)
	if err != nil {
		log.Fatal("ok")
	}
}

func deleteTodo(w http.ResponseWriter, r *http.Request, dbClient *mongo.Client, todoId string) {
	if todoId == "" {
		http.Error(w, "Todo id cannot be null", http.StatusBadRequest)
		return
	}
	objectID, err := primitive.ObjectIDFromHex(todoId)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.D{{Key: "_id", Value: objectID}}
	collection := dbClient.Database("togoDb").Collection("todos")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Example: DeleteOne
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted documents: %v\n", result.DeletedCount)
}
