package main

import (
	"fmt"
	"os"

	"github.com/prmsrswt/gophercises/adventure"
)

func main() {
	file, err := os.Open("gopher.json")
	if err != nil {
		panic("Cannot read file : gopher.json")
	}
	defer file.Close()

	story, err := adventure.ParseStory(file)
	if err != nil {
		panic("Cannot read file : gopher.json")
	}

	fmt.Println(story["intro"].Title)
}
