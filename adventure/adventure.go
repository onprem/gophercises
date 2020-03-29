package adventure

import (
	"encoding/json"
	"io"
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
