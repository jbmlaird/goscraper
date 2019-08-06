package main

import (
	"fmt"
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

	httpClient := NewHttpClient(3, 0, 1, 10)
	fmt.Println(httpClient)
	// a, b := httpClient.getUrl(urlToCrawl)
	// do some swag shit with the links when returned
}
