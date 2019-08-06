package main

import (
	"log"
)

func main() {
	urlToCrawl := "https://monzo.com"
	hostname, err := verifyHostname(urlToCrawl)

	// Could be using errors.Wrap here. Explore later.
	if err != nil {
		if err == errInvalidRegex {
			log.Fatalf("URL supplied is in the incorrect format: %v, err: %v", urlToCrawl, err)
		}
		log.Fatalf("Error parsing given URL %v, err: %v", urlToCrawl, err)
	}

	httpClient := NewHttpClient(3, 0, 1, 10)

	response, err := httpClient.getResponse(urlToCrawl)
	if err != nil {
		log.Fatalf("failed to fetch URL: %v", urlToCrawl)
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
}
