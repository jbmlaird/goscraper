package main

import (
	"github.com/pkg/errors"
	"io"
	"log"
	"strings"
	"sync"
)

type Crawler interface {
	buildSitemap(hostname string) ([]string, error)
	getResponseBody(url string) (io.ReadCloser, error)
	alreadyCrawled(url string) bool
	isSameDomain(url string) bool
}

type CrawlerImpl struct {
	hostname       string
	client         *RetryHttpClient
	*CrawlerUrlChecker
	sitemapBuilder SitemapBuilder
	mu             sync.RWMutex
}

func NewCrawler(hostname string) *CrawlerImpl {
	return &CrawlerImpl{
		hostname,
		NewRetryHttpClient(3, 0, 1, 10),
		NewCrawlerUrlTracker(),
		SitemapBuilder{},
		sync.RWMutex{},
	}
}

var errDifferentDomain = errors.New("URL belongs to another domain")
var errAlreadyCrawled = errors.New("already crawled URL")

func (c *CrawlerImpl) buildSitemap(hostname string) ([]string, error) {
	// What if a goroutine fails against a certain URL? Remove it from the sitemap?
	_ = c.request(hostname)

	//select {
	//// only add when it's finished
	//case crawledUrl := <-crawledUrls:
	//	// further handling
	//	c.sitemapBuilder.addToSitemap(crawledUrl)
	//case <-time.After(time.Second * 10):
	//	return nil, nil
	//}
	return c.sitemapBuilder.returnSitemap(), nil
}

func (c *CrawlerImpl) getResponseBody(url string) (io.ReadCloser, error) {
	c.addToCrawledUrls(url)

	response, err := c.client.getResponse(url)
	if err != nil {
		// TODO: This shouldn't be fatal
		// Just log that this URL failed and then retry?
		return nil, errors.WithMessagef(err, "failed to fetch URL: %v", c.hostname)
	}
	if response != nil {
		return response.Body, nil
	}
	return nil, errors.New("unable to read response body")
}

func (c *CrawlerImpl) request(url string) error {
	if c.alreadyCrawled(url) {
		log.Printf("skipping url %v as already been crawled", url)
		return errAlreadyCrawled
	}
	if !c.isSameDomain(url) {
		log.Printf("skipping url %v as different domain", url)
		return errDifferentDomain
	}
	log.Printf("crawling URL: %v", url)
	responseBody, err := c.getResponseBody(url)
	if err != nil {
		// TODO: Is this all I need?
		if err == errDifferentDomain || err == errAlreadyCrawled {
			// ignore. Some URLs won't be required
		} else {
			return errors.WithMessagef(err, "unable to get response body for %v", url)
		}
	} else if responseBody != nil {
		// TODO: Handle error
		urls, _ := findUrls(responseBody)
		responseBody.Close()
		// check URLs are valid
		// add to sitemap then begin request
		for _, url := range urls {
			// This is definitely wrong as it will add on this URL to external URLs
			c.request(c.hostname + url)
		}
	}
	return nil
}

func (c *CrawlerImpl) isSameDomain(url string) bool {
	// split after then check prefix?
	if (len(url) > 0 && url[0] == '/' && len(url) > 1) || strings.Contains(url, c.hostname) {
		// this needs to be tested for when hostname = monzo.com and url = community.monzo.com
		// you would need to ensure that the start of the string is empty
		return true
	}
	return false
}
