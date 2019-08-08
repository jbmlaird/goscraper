package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	start := time.Now()
	urlToCrawl := "https://monzo.com"

	crawler := NewCrawler(urlToCrawl)
	sitemap, err := crawler.buildSitemap(urlToCrawl)
	if err != nil {
		log.Fatalf("unable to crawl base URL: %v, err: %v", urlToCrawl, err)
	}
	for _, value := range sitemap {
		fmt.Println(value)
	}
	fmt.Printf("crawling took: %s", time.Since(start))
}
