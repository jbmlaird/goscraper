package main

import (
	"errors"
	"io"
	"log"
	"sync"
)

type Crawler struct {
	Client           *RetryHttpClient
	crawledUrls      map[string]struct{}
	sitemapGenerator SitemapGenerator
	mu               sync.RWMutex
}

func NewCrawler() *Crawler {
	return &Crawler{
		NewRetryHttpClient(3, 0, 1, 10),
		make(map[string]struct{}),
		SitemapGenerator{},
		sync.RWMutex{},
	}
}

func (c *Crawler) crawlWebsite(hostname string) ([]string, error) {
	//crawledUrls := make(chan string)
	// TODO: Handle error
	responseBody, _ := c.getResponseBody(hostname)

	if responseBody != nil {
		// TODO: Handle error
		urls, _ := findUrls(responseBody)
		responseBody.Close()
		for _, url := range urls {
			c.sitemapGenerator.addToSitemap(url)
		}
	}

	//select {
	//// only add when it's finished
	//case crawledUrl := <-crawledUrls:
	//	// further handling
	//	c.sitemapGenerator.addToSitemap(crawledUrl)
	//case <-time.After(time.Second * 10):
	//	return nil, nil
	//}
	return c.sitemapGenerator.returnSitemap(), nil
}

func (c *Crawler) getResponseBody(url string) (io.ReadCloser, error) {
	if c.alreadyCrawled(url) {
		log.Printf("skipping url %v as already been crawled", url)
		return nil, errors.New("already crawled")
	}
	c.mu.Lock()
	c.crawledUrls[url] = struct{}{}
	c.mu.Unlock()

	response, err := c.Client.getResponse(hostname)
	if err != nil {
		// TODO: This shouldn't be fatal
		// Just log that this URL failed and then retry?
		log.Fatalf("failed to fetch URL: %v", hostname)
	}
	if response != nil {
		return response.Body, nil
	}
	return nil, errors.New("unable to read response body")
}

func (c *Crawler) alreadyCrawled(url string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.crawledUrls[url]
	if !ok {
		return false
	}
	return true
}
