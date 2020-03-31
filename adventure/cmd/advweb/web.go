package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/prmsrswt/gophercises/adventure"
)

var (
	port      int
	storyPath string
)

func init() {
	flag.StringVar(&storyPath, "story", "gopher.json", "path of the JSON file containing story")
	flag.IntVar(&port, "port", 8080, "the port to start the server on")
}

func main() {
	flag.Parse()

	file, err := os.Open(storyPath)
	if err != nil {
		panic("Cannot read story file")
	}
	defer file.Close()

	story, err := adventure.ParseStory(file)
	if err != nil {
		panic("Error parsing json")
	}

	tplByte, err := ioutil.ReadFile("story.html")
	if err != nil {
		panic("Error reading story template")
	}

	tpl, err := template.New("").Parse(string(tplByte))
	if err != nil {
		panic("Error parsing story template")
	}

	h := adventure.NewHandler(story, adventure.WithTemplate(tpl))
	mux := http.NewServeMux()

	mux.Handle("/", h)

	fmt.Printf("Server starting on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
