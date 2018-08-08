package crawler

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/khrm/smap/internal/parser"
	"github.com/khrm/smap/internal/sitemap"
)

var (
	errInvalidURL = errors.New("link isn't valid")
)

// CondConfig is used to determine certain conditional
// configs like whether only root url are going to be
// extracted or depth of the query
type CondConfig struct {
	rootOnly bool
	depth    int
	debug    bool
}

// NewConfig gives an instance of config
func NewConfig(r bool, d int, debug bool) *CondConfig {
	return &CondConfig{
		rootOnly: r,
		depth:    d,
		debug:    debug,
	}
}

// Service contains detail needed for crawler Service
type Service struct {
	root   *url.URL
	parser parser.ServiceParse
	log    *log.Logger
	sm     *sitemap.SiteMap
	c      *CondConfig
}

// New gives an instance of crawler.service needed to crawl documents
// and put them into sitemap graph
func New(r *url.URL, p parser.ServiceParse, l *log.Logger,
	c *CondConfig) *Service {
	return &Service{
		root:   r,
		parser: p,
		log:    l,
		sm:     sitemap.New(),
		c:      c,
	}
}

// Start the crawl service and save them the links find
// in the sitemap graph
func (s *Service) Start() (sm *sitemap.SiteMap) {
	if s.c == nil {
		s.log.Println("config passed is nil")
		return nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	s.crawl(s.root, s.c, wg)
	wg.Wait()
	return s.sm
}

// crawl service get the urls and determine their links
// it save them in sitemap graph
func (s *Service) crawl(u *url.URL, c *CondConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	//current link
	clink := u.String()
	ok := s.sm.AddURL(clink)
	if !ok {
		return
	}

	if c.depth == 0 {
		return
	}

	urls, err := s.parser.ExtractURLs(clink)
	if err != nil {
		if s.c.debug {
			s.log.Println("Crawler encountered an error", err,
				"while crawling", clink)
		}
	}

	for i := range urls {
		l, err := s.urlParse(u, urls[i])
		if err != nil {
			continue
		}
		link := l.String()
		if c.rootOnly && strings.Contains(link, s.root.Host) {
			cond := *c
			// Reducing the depth
			cond.depth = cond.depth - 1
			wg.Add(1)
			go s.crawl(l, &cond, wg)
			s.sm.AddConnection(clink, link)
		}
	}
}

func (s *Service) urlParse(r *url.URL, path string) (*url.URL, error) {
	l, err := url.Parse(path)
	if err != nil {
		if s.c.debug {
			s.log.Println("link:", path, " isn't valid, err",
				err)
		}
		return nil, errInvalidURL
	}

	// If scheme isn't there we put current link scheme
	if l.Scheme == "" {
		l.Scheme = r.Scheme
	}

	// if host isn't there, we put current url's host
	if l.Host == "" {
		l.Host = r.Host
	}

	// Fragment means they are same url
	l.Fragment = ""

	// We remove trailing / as url with and without it are same
	l.Path = strings.TrimRight(l.Path, "/")
	return l, nil
}
