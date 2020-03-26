package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(urlMap map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		dest, ok := urlMap[path]
		if ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}

	paths := buildMap(parsedYaml)
	return MapHandler(paths, fallback), err
}

func buildMap(pathURLs []pathURL) map[string]string {
	paths := make(map[string]string)
	for _, pu := range pathURLs {
		paths[pu.Path] = pu.URL
	}

	return paths
}

func parseYaml(yamlBytes []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(yamlBytes, &pathURLs)
	if err != nil {
		return nil, err
	}
	return pathURLs, nil
}

type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// JSONHandler will parse the provided JSON and tries
// to map any Path provides with it's URL
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsonBytes)
	if err != nil {
		return nil, err
	}

	paths := buildMap(parsedJSON)
	return MapHandler(paths, fallback), err
}

func parseJSON(jsonBytes []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := json.Unmarshal(jsonBytes, &pathURLs)
	if err != nil {
		return nil, err
	}

	return pathURLs, nil
}

// Hello says hello to world
func Hello(name string) {
	fmt.Printf("Hello, %s!\n", name)
}
