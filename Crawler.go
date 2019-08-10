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

func (c *Crawler) BuildSitemap(urlToCrawl string) ([]string, error) {
	err := c.urlParser.VerifyBaseUrl(urlToCrawl)
	if err != nil {
		if err == errInvalidBaseUrl {
			return nil, errors.Wrapf(err, "URL supplied is in the incorrect format: %v", urlToCrawl)
		}
		return nil, errors.Wrapf(err, "Error parsing given URL %v", urlToCrawl)
	}

	var wg sync.WaitGroup
	err = c.crawlUrl(urlToCrawl, &wg)
	wg.Wait()
	c.writeErrorsToFile()
	if err != nil {
		return nil, errors.Wrapf(err, "problem trying to crawl base URL: %v", urlToCrawl)
	}
	return c.sitemapBuilder.BuildSitemap(), nil
}

func (c *Crawler) crawlUrl(url string, wg *sync.WaitGroup) error {
	cleanedUrl, err := CleanUrl(url, c.hostnameWithProtocol)
	if err != nil {
		return errors.Wrapf(err, "invalid URL passed to clean URL: %v", url)
	}
	err = c.urlParser.CheckSameDomain(cleanedUrl, c.hostnameWithProtocol)
	if err != nil {
		return errors.Wrapf(err, "%v is a different domain, original URL: %v", cleanedUrl, url)
	}
	err = c.sitemapBuilder.AddToCrawledUrls(url)
	if err != nil {
		log.Printf("%v has already been crawled", url)
		return errors.Wrapf(err, "skipping cleaned url %v, original url %v", cleanedUrl, url)
	}
	log.Printf("crawling cleaned URL: %v, original URL: %v", cleanedUrl, url)
	responseBody, err := c.getPageContents(cleanedUrl)
	if err != nil {
		return errors.Wrapf(err, "unable to get response body for cleaned URL %v, original URL %v", cleanedUrl, url)
	}
	err = c.sitemapBuilder.AddToSitemap(cleanedUrl)
	if err != nil {
		return errors.Wrapf(err, "URL already in sitemap. Cleaned URL: %v, URL: %v", cleanedUrl, url)
	}
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
			defer wg.Done()
			err = c.crawlUrl(url, wg)
			if err != nil {
				c.crawlerErrors = append(c.crawlerErrors, err)
			}
		}(url, wg)
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

// I don't want goroutine errors to crash the program
//
// In reality, one would read the logs in a UI so I have output errors to a file just for this exercise
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
