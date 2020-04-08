package main

import (
	"flag"
	"fmt"

	"github.com/prmsrswt/gophercises/sitemap"
)

func main() {
	var url string
	flag.StringVar(&url, "url", "http://calhoun.io", "URL to build SiteMap for")

	flag.Parse()

	s, err := sitemap.BuildSitemap(url)
	if err != nil {
		panic(err)
	}

	data, err := s.GetXML()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
