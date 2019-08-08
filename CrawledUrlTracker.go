package main

import "sync"

type CrawlerUrlChecker struct {
	mu          sync.Mutex
	crawledUrls map[string]struct{}
}

func NewCrawlerUrlTracker() *CrawlerUrlChecker {
	return &CrawlerUrlChecker{
		mu:          sync.Mutex{},
		crawledUrls: make(map[string]struct{}),
	}
}

func (c *CrawlerUrlChecker) isAlreadyCrawled(url string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.crawledUrls[url]
	if !ok {
		c.crawledUrls[url] = struct{}{}
		return nil
	}
	return errAlreadyCrawled
}
