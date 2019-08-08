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
	urlManipulator := NewUrlManipulator()
	hostname, err := urlManipulator.verifyBaseUrl(urlToCrawl)

	// Could be using errors.Wrap here. Explore later.
	if err != nil {
		if err == errInvalidBaseUrl {
			log.Fatalf("URL supplied is in the incorrect format: %v, err: %v", urlToCrawl, err)
		}
		log.Fatalf("Error parsing given URL %v, err: %v", urlToCrawl, err)
	}

	crawler := NewCrawler(hostname, urlManipulator)
	sitemap, err := crawler.buildSitemap(hostname)
	if err != nil {
		panic(errors.WithMessage(err, "unable to crawl URL"))
	}
	for _, value := range sitemap {
		fmt.Println(value)
	}
	fmt.Printf("crawling took: %s", time.Since(start))
}
