package main

import (
	"testing"
)

//type MockCrawler struct {
//	client           *RetryHttpClient
//	crawledUrls      map[string]struct{}
//	sitemapBuilder SitemapBuilder
//	mu               sync.RWMutex
//}
//
//func (m *MockCrawler) buildSitemap(hostname string) ([]string, error) {
//	return nil, nil
//}
//
//func (m *MockCrawler) getResponseBody(url string) (io.ReadCloser, error) {
//	return nil, nil
//}
//
//func (m *MockCrawler) alreadyCrawled(url string) bool {
//	return false
//}

func TestCrawler(t *testing.T) {
	cases := []struct{
		Name string
		Hostname string
		UrlToCheck string
		Want bool
	}{
		{"is same domain absolute path returns true",
			"https://www.monzo.com",
			"https://www.monzo.com/help",
			true,
		},
		{"is same domain relative path returns true",
			"https://www.monzo.com",
			"/help",
			true,
		},
		{"is different domain returns false",
			"https://www.monzo.com",
			"https://www.monzo.co.uk/help",
			false,
		},
		{"is homepage returns false",
			"https://www.monzo.com",
			"/",
			false,
		},
		{"is empty. returns false",
			"https://www.monzo.com",
			"",
			false,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			crawler := NewCrawler(test.Hostname)
			assertBoolean(t, crawler.isSameDomain(test.UrlToCheck), test.Want)
		})
	}
}
