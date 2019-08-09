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
	sitemapBuilder *SitemapBuilder
	channelManager *ChannelManager
}

func NewCrawler(hostname string) *CrawlerImpl {
	return &CrawlerImpl{
		addHttpsIfNecessary(hostname),
		NewUrlManipulator(),
		NewRetryHttpClient(3, 0, 1, 10),
		NewCrawlerUrlTracker(),
		NewSitemapBuilder(),
		NewChannelManager(),
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

	var wg sync.WaitGroup
	c.channelManager.StartListening()
	wg.Add(1)
	go c.request(hostname, &wg)
	wg.Wait()
	c.channelManager.CloseChannels()
	//for _, err := range c.goRoutineErrors {
	//	log.Printf("goroutine failed with: %v", err)
	//}
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
		//c.chanErr <- errors.Wrapf(err, "invalid URL passed to clean URL: %v", url)
		return
	}
	err = c.urlManipulator.checkSameDomain(cleanedUrl, c.hostnameWithProtocol)
	if err != nil {
		//c.chanErr <- errors.Wrapf(err, "%v is a different domain, original URL: %v", cleanedUrl, url)
		return
	}
	err = c.addToCrawledUrlsIfUncrawled(cleanedUrl)
	if err != nil {
		//c.chanErr <- errors.Wrapf(err, "skipping cleaned url %v, original url %v", cleanedUrl, url)
		return
	}
	log.Printf("crawling cleaned URL: %v, original URL: %v", cleanedUrl, url)
	responseBody, err := c.getResponseBody(cleanedUrl)
	if err != nil {
		//c.chanErr <- errors.Wrapf(err, "skipping cleaned url %v, original url %v", cleanedUrl, url)
		return
	}
	c.sitemapBuilder.addToSitemap(cleanedUrl)
	urls, err := findUrls(responseBody)
	if err != nil {
		//c.chanErr <- errors.Wrapf(err, "unable to find any URLs for cleaned URL %v, original URL: %v", cleanedUrl, url)
		return
	}
	err = responseBody.Close()
	if err != nil {
		//c.chanErr <- errors.Wrapf(err, "unable to close response body from cleaned URL %v, original URL %v", cleanedUrl, url)
		return
	}
	wg.Add(len(urls))
	for _, url := range urls {
		go c.request(url, wg)
	}
}

func (c *CrawlerImpl) addToCrawledUrlsIfUncrawled(cleanedUrl string) error {
	err := c.isAlreadyCrawled(cleanedUrl)
	if err != nil {
		return errors.Wrapf(err, "%v has already been crawled", cleanedUrl)
	}
	return nil
}

func (c *CrawlerImpl) saveErrChan() {
	//for {
	//	c.goRoutineErrors = append(c.goRoutineErrors, <-c.chanErr)
	//}
}
