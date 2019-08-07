package main

import "testing"

func TestCrawledUrlTracker(t *testing.T) {
	alreadyCrawledCases := []struct{
		Name        string
		CrawledUrls map[string]struct{}
		UrlToCheck string
		Want        bool
	}{
		{"already crawled detects repeated URL",
			map[string]struct{}{
				"https://monzo.com/help": {},
			},
			"https://monzo.com/help",
			true,
		},
		{
			"already crawled doesn't flag new URL",
			map[string]struct{}{
				"https://monzo.com/": {},
			},
			"https://monzo.com/help",
			false,
		},
	}

	for _, test := range alreadyCrawledCases {
		t.Run(test.Name, func(t *testing.T) {
			crawledUrlTracker := NewCrawlerUrlTracker()
			crawledUrlTracker.crawledUrls = test.CrawledUrls
			assertBoolean(t, crawledUrlTracker.alreadyCrawled(test.UrlToCheck), test.Want)
		})
	}
}
