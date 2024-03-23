package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}

func main() {
	// in golang := is used as declaration + assignment operator
	// it will infer types automatically and not let you do the wrong assignment

	http.HandleFunc("/", helloHandler)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}
	fmt.Println("\nThe server is up and running on: http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
