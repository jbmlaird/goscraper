package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
)

func main() {
	urlToCrawl := "https://monzo.com"
	hostname, err := verifyHostname(urlToCrawl)

	// Could be using errors.Wrap here. Explore later.
	if err != nil {
		if err == errInvalidRegex {
			log.Fatalf("URL supplied is in the incorrect format: %v, err: %v", urlToCrawl, err)
		}
		log.Fatalf("Error parsing given URL %v, err: %v", urlToCrawl, err)
	}

	crawler := NewCrawler()
	sitemap, err := crawler.crawlUrl(hostname)
	if err != nil {
		panic(errors.WithMessage(err, "unable to crawl URL"))
	}
	fmt.Println(sitemap)
}
