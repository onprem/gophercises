package linkparser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an <a> tag
type Link struct {
	Href string
	Text string
}

func (l Link) String() string {
	return fmt.Sprintf("<a href='%s'>%s</a>", l.Href, l.Text)
}

// ParseLinks parses a document for link tags in it
func ParseLinks(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := searchLinks(doc)

	return links, nil
}

func searchLinks(node *html.Node) []Link {
	links := make([]Link, 0)

	if node.Type == html.ElementNode && node.Data == "a" {
		l := extractLinkFromNode(node)
		return append(links, l)
	}

	for n := node.FirstChild; n != nil; n = n.NextSibling {
		links = append(links, searchLinks(n)...)
	}

	return links
}

func extractLinkFromNode(node *html.Node) Link {
	var href string
	for _, v := range node.Attr {
		if v.Key == "href" {
			href = v.Val
			break
		}
	}

	text := getLinkText(node)
	return Link{href, text}
}

func getLinkText(node *html.Node) string {
	if node.Type == html.TextNode {
		return strings.Join(strings.Fields(node.Data), " ")
	}

	var str string

	if node.Type == html.ElementNode {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			str += getLinkText(c)
		}
	}

	return strings.Join(strings.Fields(str), " ") // Trims unnecessary whitespace from the whole string
}
