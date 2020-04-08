package sitemap

import (
	"encoding/xml"
	"net/http"
	"net/url"

	"github.com/prmsrswt/gophercises/linkparser"
)

// SiteMap represents a Website SiteMap
type SiteMap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urlset  []URL    `xml:"url"`
}

// URL represents an URL entry in the SiteMap
type URL struct {
	Loc string `xml:"loc"`
}

// GetXML converts the SiteMap to XML spec
func (s *SiteMap) GetXML() ([]byte, error) {
	data, err := xml.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, err
	}

	data = append([]byte(xml.Header), data...) // Add XML header

	return data, nil
}

// BuildSitemap takes a domain and returns a generated
// XML SiteMap for that domain
func BuildSitemap(baseURL string) (*SiteMap, error) {
	smap := &SiteMap{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}

	urls := make(map[string]struct{})

	links, err := sitemapBuilder(baseURL, urls)
	if err != nil {
		return nil, err
	}

	for key := range links {
		url := URL{key}
		smap.Urlset = append(smap.Urlset, url)
	}

	return smap, nil
}

// sitemapBuilder recursively parses given URL for links in it
// uses fillterLinks to get related URLs and then repeats this
// until no new URL is found.
func sitemapBuilder(url string, urls map[string]struct{}) (map[string]struct{}, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	links, err := linkparser.ParseLinks(res.Body)
	res.Body.Close()

	filteredURLs, err := filterLinks(links, url)

	// var newURLs map[string]struct{}
	newURLs := urls

	for _, v := range filteredURLs {
		if _, ok := newURLs[v]; !ok {
			newURLs[v] = struct{}{}
			urls, err = sitemapBuilder(v, urls)
			if err != nil {
				return nil, err
			}
		}
	}

	return urls, nil
}

// filterLinks returns links related to currentURL.
// ie. the URL have same scheme and host or is relative.
func filterLinks(links []linkparser.Link, currentURL string) ([]string, error) {
	var urls []string
	current, err := url.Parse(currentURL)
	if err != nil {
		return nil, err
	}

	for _, v := range links {
		link := v.Href
		u, err := url.Parse(link)
		if err != nil {
			return nil, err
		}

		u = current.ResolveReference(u)

		if u.Host == current.Host && u.Scheme == current.Scheme {
			u.Fragment = "" // Remove the #fragment part
			urls = append(urls, u.String())
		}
	}

	return urls, nil
}
