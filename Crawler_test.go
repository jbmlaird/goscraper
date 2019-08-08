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
//func (m *MockCrawler) buildSitemap(hostnameWithProtocol string) ([]string, error) {
//	return nil, nil
//}
//
//func (m *MockCrawler) getResponseBody(url string) (io.ReadCloser, error) {
//	return nil, nil
//}
//
//func (m *MockCrawler) isAlreadyCrawled(url string) bool {
//	return false
//}

func TestCrawler(t *testing.T) {
	addToCrawledUrlsCases := []struct {
		Name                 string
		HostnameWithProtocol string
		Input                string
		Want                 string
	}{
		{"addToCrawledUrlsCases adds https:// to hostnameWithProtocol",
			"google.co.uk",
			"google.co.uk",
			"https://google.co.uk",
		},
		{"addToCrawledUrlsCases adds hostnameWithProtocol to relative path",
			"condenastint.com",
			"/help",
			"https://condenastint.com/help",
		},
		{"addToCrawledUrlsCases doesn't add https:// to https protocol",
			"https://monzo.com",
			"https://monzo.com",
			"https://monzo.com",
		},
		{"addToCrawledUrlsCases doesn't add https:// to http protocol",
			"https://monzo.com",
			"http://monzo.com/help",
			"http://monzo.com/help",
		},
	}

	for _, test := range addToCrawledUrlsCases {
		t.Run(test.Name, func(t *testing.T) {
			crawler := NewCrawler(test.HostnameWithProtocol)
			got, err := crawler.addToCrawledUrlsIfUncrawled(test.Input)
			assertNoError(t, err)
			assertStringOutput(t, got, test.Want)
		})
	}
}
