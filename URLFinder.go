package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
)

func findUrls(responseBody io.ReadCloser) (string, error) {
	defer responseBody.Close()
	document, err := goquery.NewDocumentFromReader(responseBody)
	if err != nil {
		return "", fmt.Errorf("Unable to parse response body with error %v", err)
	}
	document.Find("a[href]").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		fmt.Println(href)
	})
	return "", nil
}
