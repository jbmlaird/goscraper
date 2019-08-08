package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"time"
)

func main() {
	start := time.Now()
	urlToCrawl := "https://monzo.com"
	hostname, err := verifyHostname(urlToCrawl)

	// Could be using errors.Wrap here. Explore later.
	if err != nil {
		if err == errInvalidUrl {
			log.Fatalf("URL supplied is in the incorrect format: %v, err: %v", urlToCrawl, err)
		}
		log.Fatalf("Error parsing given URL %v, err: %v", urlToCrawl, err)
	}

	crawler := NewCrawler(hostname)
	sitemap, err := crawler.buildSitemap(hostname)
	time.Sleep(time.Second * 20)
	if err != nil {
		panic(errors.WithMessage(err, "unable to crawl URL"))
	}
	fmt.Println(sitemap)
	fmt.Printf("crawling took: %s", time.Since(start))
}
