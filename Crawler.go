package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"sync"
)

type Crawler struct {
	hostnameWithProtocol string
	urlParser            *URLParser
	client               *RetryHttpClient
	sitemapBuilder       *SitemapBuilder
	mu                   sync.RWMutex
	crawlerErrors        []error
}

func NewCrawler(hostname string) *Crawler {
	return &Crawler{
		addHttpsIfNecessary(hostname),
		NewUrlParser(),
		NewRetryHttpClient(3, 5, 10),
		NewSitemapBuilder(),
		sync.RWMutex{},
		[]error{},
	}
}

var errDifferentDomain = errors.New("URL belongs to another domain")
var errInvalidUrl = errors.New("URL is either a single character or empty")

func (c *Crawler) BuildSitemap(baseUrl string) ([]string, error) {
	err := c.urlParser.VerifyBaseUrl(baseUrl)
	if err != nil {
		if err == errInvalidBaseUrl {
			return nil, errors.Wrapf(err, "URL supplied is in the incorrect format: %v", baseUrl)
		}
		return nil, errors.Wrapf(err, "Error parsing given URL %v", baseUrl)
	}

	var wg sync.WaitGroup
	err = c.crawlUrl(baseUrl, baseUrl, &wg)
	wg.Wait()
	c.writeErrorsToFile()
	if err != nil {
		return nil, errors.Wrapf(err, "problem trying to crawl base URL: %v", baseUrl)
	}
	return c.sitemapBuilder.BuildSitemap(), nil
}

// parentUrl is passed in so that when the crawler logs, you will be able to see the page that this link is located on
// that is causing this error. This would help with debugging the website.
func (c *Crawler) crawlUrl(parentUrl, url string, wg *sync.WaitGroup) error {
	cleanedUrl, err := CleanUrl(url, c.hostnameWithProtocol)
	if err != nil {
		return errors.Wrapf(err, "invalid URL passed to clean URL: %v", url)
	}
	err = c.urlParser.CheckSameDomain(cleanedUrl, c.hostnameWithProtocol)
	if err != nil {
		return errors.Wrapf(err, "%v is a different domain, parent URL: %v, original URL: %v", cleanedUrl, parentUrl, url)
	}
	err = c.sitemapBuilder.AddToCrawledUrls(cleanedUrl)
	if err != nil {
		log.Printf("%v has already been crawled", url)
		return errors.Wrapf(err, "skipping cleaned URL %v, parent URL: %v, original URL: %v", cleanedUrl, parentUrl, url)
	}
	log.Printf("crawling cleaned URL: %v, original URL: %v", cleanedUrl, url)
	responseBody, err := c.getPageContents(cleanedUrl)
	if err != nil {
		return errors.Wrapf(err, "unable to get response body for cleaned URL %v, parent URL: %v, original URL: %v", parentUrl, cleanedUrl, url)
	}
	err = c.sitemapBuilder.AddToSitemap(cleanedUrl)
	if err != nil {
		return errors.Wrapf(err, "URL already in sitemap. Cleaned URL: %v, parent URL: %v, URL: %v", cleanedUrl, parentUrl, url)
	}
	urls, err := findUrls(responseBody)
	if err != nil {
		return errors.Wrapf(err, "unable to find any URLs for cleaned URL %v, parent URL: %v, original URL: %v", cleanedUrl, parentUrl, url)
	}
	err = responseBody.Close()
	if err != nil {
		return errors.Wrapf(err, "unable to close response body from cleaned URL %v, parent URL: %v, original URL %v", cleanedUrl, parentUrl, url)
	}
	wg.Add(len(urls))
	for _, linkedUrl := range urls {
		go func(parentUrl, linkedUrl string, wg *sync.WaitGroup) {
			defer wg.Done()
			err = c.crawlUrl(parentUrl, linkedUrl, wg)
			if err != nil {
				c.crawlerErrors = append(c.crawlerErrors, err)
			}
		}(url, linkedUrl, wg)
	}
	return nil
}

// Whack the URL and return the response body
func (c *Crawler) getPageContents(url string) (io.ReadCloser, error) {
	response, err := c.client.getResponse(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch URL: %v", c.hostnameWithProtocol)
	}
	if response != nil {
		return response.Body, nil
	}
	return nil, errors.Errorf("unable to read response body for URL %v", url)
}

// I don't want goroutine errors to crash the program, only the base URL.
//
// In reality, one would read the logs in a UI so I have output errors to a file.
func (c *Crawler) writeErrorsToFile() {
	f, err := os.Create("goroutineErrors.txt")
	defer f.Close()
	if err != nil {
		fmt.Printf("unable to write goroutine errors to a file with error: %v", err)
		return
	}
	for _, goroutineError := range c.crawlerErrors {
		_, _ = f.WriteString(fmt.Sprintf("%v\n", goroutineError))
	}
}
