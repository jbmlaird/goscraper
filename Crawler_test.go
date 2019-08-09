package main

import (
	"strings"
	"testing"
)

func TestCrawler(t *testing.T) {
	t.Run("integration test crawler output", func(t *testing.T) {
		hostname := "https://www.monzo.com"
		crawler := NewCrawler(hostname)
		sitemapLinks, err := crawler.buildSitemap(hostname)
		assertNoError(t, err)

		duplicates := returnDuplicates(sitemapLinks)
		for _, duplicate := range duplicates {
			t.Errorf("duplicate URL found in sitemap: %v", duplicate)
		}

		notSameHostname := ensureSameHostname(hostname, sitemapLinks)
		for _, incorrectUrl := range notSameHostname {
			t.Errorf("URL added to sitemap which belongs to a different domain: %v", incorrectUrl)
		}
	})
}

func returnDuplicates(urls []string) []string {
	linkMap := make(map[string]struct{})
	var duplicates []string

	for _, value := range urls {
		_, exist := linkMap[value]
		if exist {
			duplicates = append(duplicates, value)
		} else {
			linkMap[value] = struct{}{}
		}
	}

	return duplicates
}

// I could use my URLManipulator but this would be using the same code that the crawler had used to add to the sitemap
func ensureSameHostname(hostname string, sitemapLinks []string) []string {
	var incorrectLinks []string
	hostnameStripped := strings.Replace(hostname, "https://", "", -1)

	for _, sitemapLink := range sitemapLinks {
		sitemapLinkStripped := strings.Replace(sitemapLink, "https://", "", -1)
		sitemapLinkStripped = strings.Replace(sitemapLinkStripped, "http://", "", -1)
		sitemapLinkStripped = strings.Replace(sitemapLinkStripped, "smtp://", "", -1)
		sitemapLinkStripped = strings.Replace(sitemapLinkStripped, "ftp://", "", -1)
		if !strings.HasPrefix(sitemapLinkStripped, hostnameStripped) {
			incorrectLinks = append(incorrectLinks, sitemapLink)
		}
	}
	return incorrectLinks
}

func BenchmarkNewCrawler(b *testing.B) {
	hostname := "www.monzo.com"
	crawler := NewCrawler(hostname)
	for n := 0; n < b.N; n++ {
		crawler.buildSitemap(hostname)
	}
}
