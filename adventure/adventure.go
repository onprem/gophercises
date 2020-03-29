package adventure

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

// StoryArc is a part of story, kind of like a chapter
// which leads to other arcs
type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func init() {
	tpl = template.Must(template.New("").Parse(templateString))
}

var tpl *template.Template

// Story represents a Choose Your Own Adventure story.
// Each key is the name of a story chapter (aka "arc"), and
// each value is a Chapter.
type Story map[string]StoryArc

// ParseStory parses the incoming Reader and returns the story
func ParseStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story

	if err := d.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

type handler struct {
	s         Story
	t         *template.Template
	arcParser func(*http.Request) string
}

// HandlerOption is used with the NewHandler function to
// configure the http.Handler returned
type HandlerOption func(*handler)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.arcParser(r)

	chapter, ok := h.s[path]
	if !ok {
		http.Error(w, "Chapter not found", http.StatusNotFound)
		return
	}

	err := h.t.Execute(w, chapter)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

func defaultArcParser(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

// NewHandler will construct an http.Handler that will render
// the story provided.
// The default handler will use the full path (minus the / prefix)
// as the arc name, defaulting to "intro" if the path is
// empty. The default template creates option links that follow
// this pattern.
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl, defaultArcParser}

	for _, opt := range opts {
		opt(&h)
	}
	return h
}

// WithTemplate is an option to provide a custom template to
// be used when rendering stories.
func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// WithArcParser is an option to provide a custom function
// for processing the story arc from the incoming request
func WithArcParser(fn func(*http.Request) string) HandlerOption {
	return func(h *handler) {
		h.arcParser = fn
	}
}

var templateString = `
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Choose your own Adventure</title>
</head>
<body>
  <section class="page">
    <h1>{{ .Title }}</h1>
    <div class="story">
      {{ range .Story }}
        <p>{{ . }}</p>
      {{ end }}
    </div>
    {{ if .Options }}
      <div class="options">
        {{ range .Options }}
          <a href="{{ .Arc }}">{{ .Text }}</a>
        {{ end }}
      </div>
    {{ else }}
      <b>The End</b>
    {{ end }}
  </section>
  <style>
    body {
      min-height: 100vh;
      font-family: Helvetica, sans-serif;
    }
    .page {
      display: flex;
      flex-direction: column;
      align-items: center;
      width: 100%;
      box-sizing: border-box;
      padding: 50px 100px;
    }
    .story {
      width: 80%;
    }
    .options {
      display: flex;
      flex-direction: column;
      align-items: center;
    }
    a {
      text-decoration: none;
      font-size: 1.1rem;
    }
  </style>
</body>
</html>
`
