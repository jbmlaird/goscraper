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
	crawlerErrors  []error
}

func NewCrawler(hostname string) *CrawlerImpl {
	return &CrawlerImpl{
		addHttpsIfNecessary(hostname),
		NewUrlManipulator(),
		NewRetryHttpClient(3, 0, 1, 10),
		NewCrawlerUrlTracker(),
		SitemapBuilder{},
		sync.RWMutex{},
		[]error{},
	}
}

var errDifferentDomain = errors.New("URL belongs to another domain")
var errAlreadyCrawled = errors.New("already crawled URL")
var errSingleCharacter = errors.New("URL is only a single character")

func (c *CrawlerImpl) buildSitemap(urlToCrawl string) ([]string, error) {
	err := c.urlManipulator.verifyBaseUrl(urlToCrawl)
	if err != nil {
		if err == errInvalidBaseUrl {
			return nil, errors.Wrapf(err, "URL supplied is in the incorrect format: %v", urlToCrawl)
		}
		return nil, errors.Wrapf(err, "Error parsing given URL %v", urlToCrawl)
	}

	var wg sync.WaitGroup
	err = c.request(urlToCrawl, &wg)
	wg.Wait()
	for _, value := range c.crawlerErrors {
		log.Printf("a goroutine failed with error: %v", value)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "problem trying to crawl base URL: %v", urlToCrawl)
	}
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
	cleanedUrl, err := cleanUrl(url, c.hostnameWithProtocol)
	if err != nil {
		return errors.Wrapf(err, "invalid URL passed to clean URL: %v", url)
	}
	err = c.urlManipulator.checkSameDomain(cleanedUrl, c.hostnameWithProtocol)
	if err != nil {
		return errors.Wrapf(err, "%v is a different domain, original URL: %v", cleanedUrl, url)
	}
	err = c.addToCrawledUrlsIfUncrawled(cleanedUrl)
	if err != nil {
		return errors.Wrapf(err, "skipping cleaned url %v, original url %v", cleanedUrl, url)
	}
	log.Printf("crawling cleaned URL: %v, original URL: %v", cleanedUrl, url)
	responseBody, err := c.getResponseBody(cleanedUrl)
	if err != nil {
		return errors.Wrapf(err, "unable to get response body for cleaned URL %v, original URL %v", cleanedUrl, url)
	}
	c.sitemapBuilder.addToSitemap(cleanedUrl)
	urls, err := findUrls(responseBody)
	if err != nil {
		return errors.Wrapf(err, "unable to find any URLs for cleaned URL %v, original URL: %v", cleanedUrl, url)
	}
	err = responseBody.Close()
	if err != nil {
		return errors.Wrapf(err, "unable to close response body from cleaned URL %v, original URL %v", cleanedUrl, url)
	}
	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string, wg *sync.WaitGroup) {
			err = c.request(url, wg)
			if err != nil {
				c.crawlerErrors = append(c.crawlerErrors, err)
			}
		}(url, wg)
	}
	return nil
}

func (c *CrawlerImpl) addToCrawledUrlsIfUncrawled(url string) error {
	err := c.isAlreadyCrawled(url)
	if err != nil {
		log.Printf("%v has already been crawled", url)
		return errAlreadyCrawled
	}
	return nil
}
