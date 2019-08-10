package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	urlToCrawl := "https://monzo.com"

	crawler := NewCrawler(urlToCrawl)
	sitemap, err := crawler.BuildSitemap(urlToCrawl)
	if err != nil {
		log.Fatalf("unable to crawl base URL: %v, err: %v", urlToCrawl, err)
	}
	fmt.Printf("crawling took: %s", time.Since(start))
	writeSitemapToFile(sitemap)
}

func writeSitemapToFile(sitemap []string) {
	f, err := os.Create("sitemap.txt")
	defer f.Close()
	if err != nil {
		fmt.Printf("unable to write sitemap to a file with error: %v", err)
		return
	}
	for _, goroutineError := range sitemap {
		_, _ = f.WriteString(fmt.Sprintf("%v\n", goroutineError))
	}
}
