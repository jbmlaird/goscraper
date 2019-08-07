package main

import (
	"log"
)

type Crawler struct {
	Client *RetryHttpClient
}

func NewCrawler() *Crawler {
	return &Crawler{
		NewRetryHttpClient(3, 0, 1, 10),
	}
}

func (c *Crawler) crawlUrl(hostname string) ([]string, error) {
	response, err := c.Client.getResponse(hostname)
	if err != nil {
		log.Fatalf("failed to fetch URL: %v", hostname)
	}
	if response != nil {
		defer response.Body.Close()
		urls, err := findUrls(response.Body)
		if err != nil {
			log.Printf("unable to parse response body, err: %v", err)
		}
		sitemapGenerator := SitemapGenerator{}
		for _, href := range urls {
			// TODO: will need to strip protocol and/or hostname
			if isSubdomain(hostname, href) && !sitemapGenerator.contains(href) {

			} else {
				log.Printf("Ignoring '%v' as it is not a subdomain of %v", href, hostname)
			}
		}
		// do some swag shit with the links when returned
	}

	return nil, nil
}
