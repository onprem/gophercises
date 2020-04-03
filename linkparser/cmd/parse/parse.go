package main

import (
	"fmt"
	"os"

	"github.com/prmsrswt/gophercises/linkparser"
)

func main() {
	file, err := os.Open("test.html")
	if err != nil {
		panic("error opening file")
	}
	defer file.Close()

	fmt.Println("parsing html...")

	links, err := linkparser.ParseLinks(file)
	if err != nil {
		panic("error parsing links")
	}

	fmt.Printf("Links found: %d\n\n", len(links))
	for _, v := range links {
		fmt.Println(v)
	}
}
