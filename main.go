package main

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func main() {
	start := time.Now()
	urlToCrawl := "https://monzo.com"

	crawler := NewCrawler(urlToCrawl)
	sitemap, err := crawler.buildSitemap(urlToCrawl)
	if err != nil {
		panic(errors.WithMessage(err, "unable to crawl URL"))
	}
	for _, value := range sitemap {
		fmt.Println(value)
	}
	fmt.Printf("crawling took: %s", time.Since(start))
}
