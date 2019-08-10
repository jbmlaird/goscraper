package main

import (
	"github.com/pkg/errors"
	"sort"
	"sync"
)

type SitemapBuilder struct {
	crawledUrls map[string]struct{}
	sitemapUrls map[string]struct{}
	mutex           sync.Mutex
}

func NewSitemapBuilder() *SitemapBuilder {
	return &SitemapBuilder{
		crawledUrls: map[string]struct{}{},
		sitemapUrls: map[string]struct{}{},
		mutex: sync.Mutex{},
	}
}

var errAlreadyCrawled = errors.New("already crawled URL")
var errAlreadyInSitemap = errors.New("already added URL to sitemap")

func (s *SitemapBuilder) AddToCrawledUrls(url string) error {
	return s.addToMap(url, s.crawledUrls, errAlreadyCrawled)
}

func (s *SitemapBuilder) AddToSitemap(url string) error {
	return s.addToMap(url, s.sitemapUrls, errAlreadyInSitemap)
}

func (s *SitemapBuilder) addToMap(url string, urlMap map[string]struct{}, err error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := urlMap[url]
	if !ok {
		urlMap[url] = struct{}{}
		return nil
	}
	return err
}

func (s *SitemapBuilder) BuildSitemap() []string {
	sitemap := make([]string, len(s.sitemapUrls))
	i := 0
	for key := range s.sitemapUrls {
		sitemap[i] = key
		i++
	}
	sort.Slice(sitemap, func(firstElement, secondElement int) bool {
		return sitemap[firstElement] < sitemap[secondElement]
	})
	return sitemap
}
