package parser

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

var (
	errInvalidContentTypeHeader = errors.New("unsupported content" +
		" type header for crawling")
	errLink404 = errors.New("url gives 404")
)

// transportClient defines the interface needed to get content from url
type transportClient interface {
	Get(url string) (resp *http.Response, err error)
}

// ParserService is the interface which satisfy the service of extracting
// URLS from HTML
type Service interface {
	ExtractURLs(url string) ([]string, error)
}

type parser struct {
	log    *log.Logger
	client transportClient
	debug  bool
}

// New returns a new parser used for links extractions
func New(client transportClient, l *log.Logger, debug bool) Service {
	return &parser{
		client: client,
		log:    l,
		debug:  debug,
	}
}

// ExtractURLs make request to url and fetch response
// response is used to get links if it is html
func (p *parser) ExtractURLs(url string) ([]string, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		if p.debug {
			p.log.Printf("Error :%s encountered crawling link: %s",
				err, url)
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		if p.debug {
			p.log.Printf("URL %s gives 404", url)
		}
		return nil, errLink404
	}

	h := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(h, "text/html") {
		return nil, errInvalidContentTypeHeader
	}

	return p.linksInBody(resp.Body), nil
}

// linksInBody get all links present in a html document
func (p *parser) linksInBody(body io.ReadCloser) []string {
	t := html.NewTokenizer(body)

	urls := []string{}
	for tt := t.Next(); tt != html.ErrorToken; tt = t.Next() {
		token := t.Token()
		if tt == html.StartTagToken && token.Data == "a" {
			attr := token.Attr
			for _, a := range attr {
				// Removing spaces as they are valid in html

				v := strings.TrimSpace(a.Val)
				key := strings.TrimSpace(a.Key)
				if key == "href" {
					urls = append(urls, v)
				}
			}
		}
	}
	return urls
}
