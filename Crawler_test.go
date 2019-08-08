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
	isSameDomainCases := []struct {
		Name       string
		Hostname   string
		UrlToCheck string
		Want       bool
	}{
		{"isSameDomain absolute path returns true",
			"https://www.monzo.com",
			"https://www.monzo.com/help",
			true,
		},
		{"isSameDomain relative path returns true",
			"https://www.monzo.com",
			"/help",
			true,
		},
		{"isSameDomain different domain returns false",
			"https://www.monzo.com",
			"https://www.monzo.co.uk/help",
			false,
		},
		{"isSameDomain homepage returns false",
			"https://www.monzo.com",
			"/",
			false,
		},
		{"isSameDomain empty returns false",
			"https://www.monzo.com",
			"",
			false,
		},
	}

	for _, test := range isSameDomainCases {
		t.Run(test.Name, func(t *testing.T) {
			crawler := NewCrawler(test.Hostname)
			assertBoolean(t, crawler.isSameDomain(test.UrlToCheck), test.Want)
		})
	}

	addHttpsIfNecessaryCases := []struct {
		Name  string
		Input string
		Want  string
	}{
		{"addHttpsIfNecessary adds https:// to hostname",
			"monzo.com",
			"https://monzo.com",
		},
		{"addHttpsIfNecessary adds https:// to relative path",
			"/help",
			"https://monzo.com/help",
		},
	}

	for _, test := range addHttpsIfNecessaryCases {
		t.Run(test.Name, func(t *testing.T) {
			got := addHttpsIfNecessary(test.Input)
			assertStringOutput(t, got, test.Want)
		})
	}
}
