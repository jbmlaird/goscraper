package main

import "sync"

type CrawlerUrlChecker struct {
	mu sync.RWMutex
	crawledUrls    map[string]struct{}
}

func NewCrawlerUrlTracker() *CrawlerUrlChecker {
	return &CrawlerUrlChecker{
		mu:          sync.RWMutex{},
		crawledUrls: make(map[string]struct{}),
	}
}

func (c *CrawlerUrlChecker) alreadyCrawled(url string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.crawledUrls[url]
	if !ok {
		return false
	}
	return true
}

func (c *CrawlerUrlChecker) addToCrawledUrls(url string) {
	c.mu.Lock()
	c.crawledUrls[url] = struct{}{}
	c.mu.Unlock()
}
