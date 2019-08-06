package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"io"
)

func findUrls(responseBody io.ReadCloser) (urls []string, err error) {
	document, err := goquery.NewDocumentFromReader(responseBody)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse response body")
	}
	document.Find("a[href]").Each(func(i int, selection *goquery.Selection) {
		// if valid link (absolute path?)
		// not contained in sitemap already
		// get response body and find URLs
		href, _ := selection.Attr("href")
		urls = append(urls, href)
	})
	return urls, nil
}
