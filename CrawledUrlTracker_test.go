package main

import "testing"

func TestCrawledUrlTracker(t *testing.T) {
	alreadyCrawledFailCases := []struct {
		Name        string
		CrawledUrls map[string]struct{}
		UrlToCheck  string
	}{
		{"already crawled detects repeated URL",
			map[string]struct{}{
				"https://monzo.com/help": {},
			},
			"https://monzo.com/help",
		},
	}

	for _, test := range alreadyCrawledFailCases {
		t.Run(test.Name, func(t *testing.T) {
			crawledUrlTracker := NewCrawlerUrlTracker()
			crawledUrlTracker.crawledUrls = test.CrawledUrls
			assertErrorMessage(t, crawledUrlTracker.isAlreadyCrawled(test.UrlToCheck), errAlreadyCrawled.Error())
		})
	}

	alreadyCrawledSuccessCases := []struct {
		Name        string
		CrawledUrls map[string]struct{}
		UrlToCheck  string
	}{
		{
			"already crawled doesn't flag new URL",
			map[string]struct{}{
				"https://monzo.com/": {},
			},
			"https://monzo.com/help",
		},
	}

	for _, test := range alreadyCrawledSuccessCases {
		t.Run(test.Name, func(t *testing.T) {
			crawledUrlTracker := NewCrawlerUrlTracker()
			crawledUrlTracker.crawledUrls = test.CrawledUrls
			assertNoError(t, crawledUrlTracker.isAlreadyCrawled(test.UrlToCheck))
		})
	}
}
