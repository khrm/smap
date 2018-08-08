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
	"time"

	"github.com/khrm/smap/internal/crawler"
	"github.com/khrm/smap/internal/parser"
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

func main() {
	// Getting configuration
	domain := flag.String("domain", "goharbor.io", "domain to crawl")
	depth := flag.Int("depth", -3, "depth to crawl")
	root := flag.Bool("root", true, "restrict to only domain given")
	debug := flag.Bool("debug", false, "whether to print everything")
	scheme := flag.String("scheme", "https",
		"scheme of the domain like http")
	stdXMLSiteMap := flag.Bool("stdsmap", true, "whether to"+
		" print standard sitemap xml")

	flag.Parse()

	p := parser.New(httpClient, logger, *debug)

	c := crawler.NewConfig(*root, *depth, *debug)
	u, err := url.Parse(*domain)
	if err != nil {
		log.Println("Failed")

	}

	if u.Scheme == "" {
		u.Scheme = *scheme
	}

	if u.Host == "" {
		u.Host = u.Path
		u.Path = ""
	}

	crawl := crawler.New(u, p, logger, c)

	sm := crawl.Start()

	data, err := json.MarshalIndent(sm, "  ", "    ")
	if err != nil {
		log.Println("error marshaling data to json", err)
	}
	fmt.Println(string(data))

	if *stdXMLSiteMap {
		xsm, err := sm.ToXMLSTDSiteMap()
		if err != nil {
			log.Println("error marshaling data to xml", err)
		}
		fmt.Println("StdSiteMap:\n", string(xsm))
	}
}
