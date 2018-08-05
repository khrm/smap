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

type transportClient interface {
	Get(url string) (resp *http.Response, err error)
}

// Parse is the interface which satisfy the service of extracting
// URLS from HTML
type Parse interface {
	ExtractURLs(url string) ([]string, error)
}

type parser struct {
	log    *log.Logger
	client transportClient
}

func New(client transportClient, l *log.Logger) *parser {
	return &parser{
		client: client,
		log:    l}
}

func (p *parser) ExtractURLs(url string) ([]string, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		p.log.Printf("Error :%s encountered crawling link: %s", err, url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		p.log.Printf("URL %s gives 404", url)
		return nil, errLink404
	}

	h := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(h, "text/html") {
		return nil, errInvalidContentTypeHeader
	}

	return p.linksInBody(resp.Body), nil
}

func (p *parser) linksInBody(body io.ReadCloser) []string {
	t := html.NewTokenizer(body)

	urls := []string{}
	for tt := t.Next(); tt != html.ErrorToken; tt = t.Next() {
		token := t.Token()
		if tt == html.StartTagToken && token.Data == "a" {
			attr := token.Attr
			for _, a := range attr {
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
