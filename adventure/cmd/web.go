package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prmsrswt/gophercises/adventure"
)

func main() {
	file, err := os.Open("gopher.json")
	port := 8080
	if err != nil {
		panic("Cannot read file : gopher.json")
	}
	defer file.Close()

	story, err := adventure.ParseStory(file)
	if err != nil {
		panic("Error parsing json")
	}

	h := adventure.NewHandler(story)
	mux := http.NewServeMux()

	mux.Handle("/", h)

	fmt.Printf("Server starting on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
