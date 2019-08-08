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
	hostname string
	client   *RetryHttpClient
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
	var wg sync.WaitGroup
	wg.Add(1)
	_ = c.request(hostname, &wg)
	wg.Done()
	wg.Wait()
	return c.sitemapBuilder.returnSitemap(), nil
}

func (c *CrawlerImpl) getResponseBody(url string) (io.ReadCloser, error) {
	response, err := c.client.getResponse(url)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to fetch URL: %v", c.hostname)
	}
	if response != nil {
		return response.Body, nil
	}
	return nil, errors.Errorf("unable to read response body for URL %v", url)
}

func (c *CrawlerImpl) request(url string, wg *sync.WaitGroup) error {
	defer wg.Done()
	url, err := c.addToCrawledUrlsIfUncrawled(url)
	if err != nil {
		return errors.Wrapf(err, "skipping url", url)
	}
	log.Printf("crawling URL: %v", url)
	responseBody, err := c.getResponseBody(url)
	if err != nil {
		return errors.WithMessagef(err, "unable to get response body for %v", url)
	}
	if responseBody != nil {
		c.sitemapBuilder.addToSitemap(url)
		// TODO: Handle error
		urls, _ := findUrls(responseBody)
		responseBody.Close()
		// check URLs are valid
		// add to sitemap then begin request
		for _, url := range urls {
			wg.Add(1)
			go c.request(url, wg)
		}
	}
	return nil
}

func (c *CrawlerImpl) addToCrawledUrlsIfUncrawled(url string) (string, error) {
	if len(url) > 1 && url[0] == '/' {
		url = c.hostname + url
	}
	if !c.isSameDomain(url) {
		log.Printf("%v is a different domain", url)
		return "", errDifferentDomain
	}
	if c.alreadyCrawled(url) {
		log.Printf("%v has already been crawled", url)
		return "", errAlreadyCrawled
	}
	return url, nil
}

func addHttpsIfNecessary(url string) string {
	if !strings.HasPrefix(url, "https://") || !strings.HasPrefix(url, "http://") {
		return "https://" + url
	}
	return url
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
