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
}

// NewConfig gives an instance of config
func NewConfig(r bool, d int) *CondConfig {
	return &CondConfig{
		rootOnly: r,
		depth:    d,
	}
}

// Service contains detail needed for crawler Service
type Service struct {
	root   *url.URL
	parser parser.ServiceParse
	log    *log.Logger
	SM     *sitemap.SiteMap
	wg     *sync.WaitGroup
	debug  bool
}

// New gives an instance of crawler.service needed to crawl documents
// and put them into sitemap graph
func New(r *url.URL, p parser.ServiceParse, l *log.Logger,
	s *sitemap.SiteMap, wg *sync.WaitGroup, debug bool) *Service {
	return &Service{
		root:   r,
		parser: p,
		log:    l,
		SM:     s,
		wg:     wg,
		debug:  debug,
	}
}

// Crawl service get the urls and determine their links
// it save them in sitemap graph
func (s *Service) Crawl(u *url.URL, c *CondConfig) {
	defer s.wg.Done()

	if c == nil {
		if s.debug {
			s.log.Println("config passed is nil")
		}
		return
	}

	//current link
	clink := u.String()
	ok := s.SM.AddURL(clink)
	if !ok {
		return
	}

	if c.depth <= 0 {
		return
	}

	urls, err := s.parser.ExtractURLs(clink)
	if err != nil {
		if s.debug {
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
			s.wg.Add(1)
			go s.Crawl(l, &cond)
			s.SM.AddConnection(clink, link)
		}
	}
}

func (s *Service) urlParse(r *url.URL, path string) (*url.URL, error) {
	l, err := url.Parse(path)
	if err != nil {
		if s.debug {
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
