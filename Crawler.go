package main

import (
	"github.com/pkg/errors"
	"io"
	"log"
	"sync"
)

type Crawler interface {
	buildSitemap(hostname string) ([]string, error)
	getResponseBody(url string) (io.ReadCloser, error)
	alreadyCrawled(url string) bool
	isSameDomain(url string) bool
}

type CrawlerImpl struct {
	hostnameWithProtocol string
	urlManipulator       *URLManipulator
	client               *RetryHttpClient
	*CrawlerUrlChecker
	sitemapBuilder SitemapBuilder
	mu             sync.RWMutex
}

func NewCrawler(hostname string, manipulator *URLManipulator) *CrawlerImpl {
	return &CrawlerImpl{
		addHttpsIfNecessary(hostname),
		manipulator,
		NewRetryHttpClient(3, 0, 1, 10),
		NewCrawlerUrlTracker(),
		SitemapBuilder{},
		sync.RWMutex{},
	}

}

var errDifferentDomain = errors.New("URL belongs to another domain")
var errAlreadyCrawled = errors.New("already crawled URL")
var errSingleCharacter = errors.New("URL is only a single character")

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
		return nil, errors.WithMessagef(err, "failed to fetch URL: %v", c.hostnameWithProtocol)
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
		return errors.Wrapf(err, "skipping url %v", url)
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
		for _, url := range urls {
			wg.Add(1)
			go c.request(url, wg)
		}
	}
	return nil
}

func (c *CrawlerImpl) addToCrawledUrlsIfUncrawled(url string) (string, error) {
	if len(url) <= 1 {
		return "", errSingleCharacter
	}
	url = addHostnameAndProtocolToRelativeUrls(url, c.hostnameWithProtocol)
	url = addHttpsIfNecessary(url)
	err := c.urlManipulator.checkSameDomain(url, c.hostnameWithProtocol)
	if err != nil {
		log.Printf("%v is a different domain", url)
		return "", errDifferentDomain
	}
	err = c.isAlreadyCrawled(url)
	if err != nil {
		log.Printf("%v has already been crawled", url)
		return "", errAlreadyCrawled
	}
	return url, nil
}
