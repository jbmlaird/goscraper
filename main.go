package main

import (
	"log"
)

func main() {
	urlToCrawl := "https://monzo.com"
	validUrl, err := verifyUrl(urlToCrawl)

	// Could be using errors.Wrap here. Explore later.
	if err != nil {
		log.Fatalf("Error parsing given URL %v, err: %v", urlToCrawl, err)
	} else if !validUrl {
		log.Fatalf("URL supplied is in the incorrect format: %v, err: %v", urlToCrawl, err)
	}

	httpClient := NewHttpClient(3, 0, 1, 1)

	//var wg sync.WaitGroup
	//wg.Add()

	response, err := httpClient.getUrl(urlToCrawl)
	if err != nil {
		log.Fatalf("failed to fetch URL: %v", urlToCrawl)
	}
	if response != nil {
		findUrls(response.Body)
		// do some swag shit with the links when returned
	}
}
