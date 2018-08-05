package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/khrm/smap/internal/crawler"
	"github.com/khrm/smap/internal/parser"
	"github.com/khrm/smap/internal/sitemap"
)

var (
	logger     = log.New(os.Stdout, "logger: ", log.Lshortfile)
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}
)

var wg sync.WaitGroup

func main() {
	domain := flag.String("domain", "goharbor.io", "domain to crawl")
	depth := flag.Int("depth", 3, "depth to crawl")
	root := flag.Bool("root", true, "restrict to only domain given")
	debug := flag.Bool("debug", false, "whether to print everything")
	scheme := flag.String("scheme", "https",
		"scheme of the domain like http")
	stdXMLSiteMap := flag.Bool("stdsmap", true, "whether to"+
		" print standard sitemap xml")

	flag.Parse()

	p := parser.New(httpClient, logger, *debug)

	s := sitemap.New()
	c := crawler.NewConfig(*root, *depth)
	u, err := url.Parse(*domain)
	if err != nil {
		log.Println("Failed")

	}

	if *depth < 1 {
		logger.Println("Invalid depth entered")
	}

	if u.Scheme == "" {
		u.Scheme = *scheme
	}

	if u.Host == "" {
		u.Host = u.Path
		u.Path = ""
	}

	cl := crawler.New(u, p, logger, s, &wg, *debug)

	wg.Add(1)
	go cl.Crawl(u, c)
	wg.Wait()
	data, err := json.MarshalIndent(cl.SM, "  ", "    ")
	if err != nil {
		log.Println("error marshaling data to json", err)
	}
	fmt.Println(string(data))

	if *stdXMLSiteMap {
		xsm, err := cl.SM.ToXMLSTDSiteMap()
		if err != nil {
			log.Println("error marshaling data to xml", err)
		}
		fmt.Println("StdSiteMap:\n", string(xsm))
	}
}
