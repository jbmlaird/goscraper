package main

import (
	"log"
	"sort"
	"sync"
)

type SitemapGenerator struct {
	sitemapLinks []string
	mu sync.Mutex
}

// sitemapLinks is sorted afterwards so that sort.Search() can be called which must be called on a sorted slice
func (s *SitemapGenerator) addToSitemap(link string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.contains(link) {
		log.Printf("link already exists in sitemap: %v", link)
	} else {
		s.sitemapLinks = append(s.sitemapLinks, link)
		sort.Strings(s.sitemapLinks)
	}
}

func (s *SitemapGenerator) returnSitemap() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sitemapLinks
}

// This method is not directly tested because it's provided by the Golang documentation:
// https://golang.org/pkg/sort/#SearchStrings calls https://golang.org/pkg/sort/#Search
func (s *SitemapGenerator) contains(string string) bool {
	// TODO: Not sure why mutex locking here causes the app to hang. Figure out why
	// s.mu.Lock()
	// defer s.mu.Unlock()
	i := sort.SearchStrings(s.sitemapLinks, string)
	return i < len(s.sitemapLinks) && s.sitemapLinks[i] == string
}