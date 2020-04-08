package main

import (
	"fmt"

	"github.com/prmsrswt/gophercises/sitemap"
)

func main() {
	s, err := sitemap.BuildSitemap("http://calhoun.io/")
	if err != nil {
		panic(err)
	}

	data, err := s.GetXML()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
