package todoHandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	m "github.com/itse4elhaam/togo-api.git/src/models"
)

func TodosController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		CreateTodos(w,r)
		return
	}else if r.Method == "GET"{
		GetTodos(w,r)
	}
}

// GET /api/todo/
func GetTodos(w http.ResponseWriter, r *http.Request) {
	todos := []m.Todo{
		{
			ID:        1,
			Title:     "Write a blog post",
			Completed: false,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			Title:     "Learn a new programming concept",
			Completed: false,
			CreatedAt: time.Now(),
		},
		{
			ID:        3,
			Title:     "Plan a weekend trip",
			Completed: false,
			CreatedAt: time.Now(),
		},
		// ... Add more todo items as needed
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(todos)
	if err != nil {
		// Handle encoding errors
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateTodos(w http.ResponseWriter, r*http.Request){

	var t m.Todo;

    err := json.NewDecoder(r.Body).Decode(&t)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	fmt.Println(t);
}
