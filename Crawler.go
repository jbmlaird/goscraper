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
	chanErr chan error
}

func NewCrawler(hostname string) *CrawlerImpl {
	return &CrawlerImpl{
		addHttpsIfNecessary(hostname),
		NewUrlManipulator(),
		NewRetryHttpClient(3, 0, 1, 10),
		NewCrawlerUrlTracker(),
		SitemapBuilder{},
		sync.RWMutex{},
		make(chan error),
	}
}

var errDifferentDomain = errors.New("URL belongs to another domain")
var errAlreadyCrawled = errors.New("already crawled URL")
var errSingleCharacter = errors.New("URL is only a single character")

func (c *CrawlerImpl) buildSitemap(urlToCrawl string) ([]string, error) {
	hostname, err := c.urlManipulator.verifyBaseUrl(urlToCrawl)
	if err != nil {
		if err == errInvalidBaseUrl {
			return nil, errors.Wrapf(err, "URL supplied is in the incorrect format: %v", urlToCrawl)
		}
		return nil, errors.Wrapf(err, "Error parsing given URL %v", urlToCrawl)
	}

	// What if a goroutine fails against a certain URL? Remove it from the sitemap?
	var wg sync.WaitGroup
	c.request(hostname, &wg)
	wg.Wait()
	for err := range c.chanErr {
		log.Printf("goroutine failed with: %v", err)
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

func (c *CrawlerImpl) request(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	cleanedUrl, err := cleanUrl(url, c.hostnameWithProtocol)
	if err != nil {
		c.chanErr<-errors.Wrapf(err, "invalid URL passed to clean URL: %v", url)
		return
		//return errors.Wrapf(err, "invalid URL passed to clean URL: %v", url)
	}
	// TODO: WaitGroup crashes if this code is here rather than in the addToCrawledUrlsIfUncrawled method
	// no idea why
	/*err = c.urlManipulator.checkSameDomain(url, c.hostnameWithProtocol)
	//if err != nil {
	//	return errors.Wrapf(err, "%v is a different domain, original URL: %v", cleanedUrl, url)
	}*/
	err = c.addToCrawledUrlsIfUncrawled(cleanedUrl)
	if err != nil {
		c.chanErr<-errors.Wrapf(err, "skipping cleaned url %v, original url %v", cleanedUrl, url)
		return
		//return errors.Wrapf(err, "skipping cleaned url %v, original url %v", cleanedUrl, url)
	}
	log.Printf("crawling cleaned URL: %v, original URL: %v", cleanedUrl, url)
	responseBody, err := c.getResponseBody(cleanedUrl)
	if err != nil {
		c.chanErr<-errors.Wrapf(err, "skipping cleaned url %v, original url %v", cleanedUrl, url)
		return
		//return errors.Wrapf(err, "unable to get response body for cleaned URL %v, original URL %v", cleanedUrl, url)
	}
	c.sitemapBuilder.addToSitemap(cleanedUrl)
	urls, err := findUrls(responseBody)
	if err != nil {
		c.chanErr<-errors.Wrapf(err, "unable to find any URLs for cleaned URL %v, original URL: %v", cleanedUrl, url)
		return
		//return errors.Wrapf(err, "unable to find any URLs for cleaned URL %v, original URL: %v", cleanedUrl, url)
	}
	err = responseBody.Close()
	if err != nil {
		c.chanErr<-errors.Wrapf(err, "unable to close response body from cleaned URL %v, original URL %v", cleanedUrl, url)
		return
		//return errors.Wrapf(err, "unable to close response body from cleaned URL %v, original URL %v", cleanedUrl, url)
	}
	log.Printf("adding waitgroup: %d", len(urls))
	for _, url := range urls {
		// TODO: Handle error for Goroutine
		go c.request(url, wg)
	}
}

func (c *CrawlerImpl) addToCrawledUrlsIfUncrawled(url string) error {
	err := c.urlManipulator.checkSameDomain(url, c.hostnameWithProtocol)
	if err != nil {
		return errors.Wrapf(err, "%v is a different domain", url)
	}
	err = c.isAlreadyCrawled(url)
	if err != nil {
		return errors.Wrapf(err, "%v has already been crawled", url)
	}
	return nil
}
