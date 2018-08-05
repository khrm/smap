package crawler

import (
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/khrm/smap/internal/parser"
	"github.com/khrm/smap/internal/sitemap"
)

// Config is used to determine certain configs like
// whether only root url are going to be extracted
// depth of the query
type config struct {
	RootOnly bool
	Depth    int
}

// service contains detail needed for crawler service
type service struct {
	root   *url.URL
	parser parser.Parse
	log    *log.Logger
	SM     *sitemap.SiteMap
	wg     *sync.WaitGroup
	debug  bool
}

// NewConfig gives an instance of config
func NewConfig(r bool, d int) *config {
	return &config{
		RootOnly: r,
		Depth:    d,
	}
}

// New gives an instance of crawler.service needed to crawl documents
// and put them into sitemap graph
func New(r *url.URL, p parser.Parse, l *log.Logger,
	s *sitemap.SiteMap, wg *sync.WaitGroup, debug bool) *service {
	return &service{
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
func (s *service) Crawl(u *url.URL, c *config) {
	defer s.wg.Done()
	//	fmt.Println("Inside Crawl")
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

	if c.Depth <= 0 {
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
		l, err := url.Parse(urls[i])
		if err != nil {
			if s.debug {

				s.log.Println("link:", urls[i], " isn't valid, err",
					err)
			}
			continue
		}

		// If scheme isn't there we put current link scheme
		if l.Scheme == "" {
			l.Scheme = u.Scheme
		}

		// if host isn't there, we put current url's host
		if l.Host == "" {
			l.Host = u.Host
		}

		// Fragment means they are same url
		l.Fragment = ""

		// We remove trailing / as url with and without it are same
		l.Path = strings.TrimRight(l.Path, "/")

		link := l.String()
		if c.RootOnly && strings.Contains(link, s.root.Host) {
			cond := *c
			// Reducing the depth
			cond.Depth = cond.Depth - 1
			s.wg.Add(1)
			go s.Crawl(l, &cond)
			s.SM.AddConnection(clink, link)
		}
	}
}
