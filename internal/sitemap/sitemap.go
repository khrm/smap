package sitemap

import (
	"encoding/xml"
	"sync"
)

// SiteMap DataStructure containing urls and connections
// SiteMap is just a graph
type SiteMap struct {
	URLs        map[string]struct{}
	Connections map[string]map[string]struct{}
	sync.Mutex
}

// New Gives an instance of SiteMap
func New() *SiteMap {
	return &SiteMap{
		URLs:        make(map[string]struct{}),
		Connections: make(map[string]map[string]struct{}),
	}
}

// AddURL add a new url in the sitemap
// If url already exist, then it returns false
func (s *SiteMap) AddURL(url string) bool {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.URLs[url]; ok {
		return false
	}
	s.URLs[url] = struct{}{}
	return true
}

// AddConnection add a new connection from u to v in the sitemap
// If connection already exist, then it returns true,false
// It also returns false,false if any node doesn't already exists
func (s *SiteMap) AddConnection(u, v string) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.Connections[u]; !ok {
		s.Connections[u] = make(map[string]struct{})
	}
	s.Connections[u][v] = struct{}{}
}

// ToXMLSTDSiteMap gives you standardise sitemap give root url
func (s *SiteMap) ToXMLSTDSiteMap() ([]byte, error) {
	xsm := struct {
		XMLName   xml.Name `xml:"urlset"`
		XMLnsAttr string   `xml:"xmlns,attr"`
		URL       []struct {
			Loc string `xml:"loc"`
		} `xml:"url"`
	}{
		XMLnsAttr: "https://www.sitemaps.org/schemas/sitemap/0.9",
	}
	for i := range s.URLs {
		xsm.URL = append(xsm.URL, struct {
			Loc string `xml:"loc"`
		}{i})
	}

	return xml.Marshal(xsm)

}
