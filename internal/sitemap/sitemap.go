package sitemap

// SiteMap DataStructure containing urls and connections
// SiteMap is just a graph
type SiteMap struct {
	URLs        map[string]struct{}
	Connections map[string]map[string]struct{}
}

// Gives a New SiteMap
func New() *SiteMap {
	return &SiteMap{
		URLs:        make(map[string]struct{}),
		Connections: make(map[string]map[string]struct{}),
	}
}

// AddURL add a new url in the sitemap
// If url already exist, then it returns false
func (s *SiteMap) AddURL(url string) bool {
	if _, ok := s.URLs[url]; ok {
		return false
	}
	s.URLs[url] = struct{}{}
	return true
}

// AddConnection add a new connection from u to v in the sitemap
// If connection already exist, then it returns true,false
// It also returns false,false if any node doesn't already exists
func (s *SiteMap) AddConnection(u, v string) (bool, bool) {
	if _, ok := s.URLs[u]; !ok {
		return false, false
	}
	if _, ok := s.URLs[v]; !ok {
		return false, false
	}

	if _, ok := s.Connections[u]; !ok {
		s.Connections[u] = make(map[string]struct{})
	}
	if _, ok := s.Connections[u][v]; ok {
		return true, false
	}
	s.Connections[u][v] = struct{}{}
	return true, true
}
