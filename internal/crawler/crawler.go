package crawler

import (
	"log"
	"net/url"
	"strings"

	"github.com/khrm/smap/internal/parser"
	"github.com/khrm/smap/internal/sitemap"
)

type config struct {
	RootOnly bool
	Depth    int
}

type service struct {
	root   *url.URL
	parser parser.Parse
	log    *log.Logger
	SM     *sitemap.SiteMap
}

func NewConfig(r bool, d int) *config {
	return &config{
		RootOnly: r,
		Depth:    d,
	}
}

func New(r *url.URL, p parser.Parse, l *log.Logger,
	s *sitemap.SiteMap) *service {
	return &service{
		root:   r,
		parser: p,
		log:    l,
		SM:     s,
	}
}

func (s *service) Crawl(u *url.URL, c *config) {
	if c == nil {
		s.log.Println("config passed is nil")
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
		s.log.Println("Crawler encountered an error", err,
			"while crawling", clink)
	}

	for i := range urls {
		l, err := url.Parse(urls[i])
		if err != nil {
			s.log.Println("link:", urls[i], " isn't valid, err",
				err)
			continue
		}
		if l.Scheme == "" {
			l.Scheme = u.Scheme
		}

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
			cond.Depth = cond.Depth - 1
			s.Crawl(l, &cond)
			s.SM.AddConnection(clink, link)
		}
	}
}
