package main

import (
	"github.com/pkg/errors"
	"sort"
	"sync"
)

type SitemapBuilder struct {
	crawledUrls      map[string]struct{}
	sitemapUrls      map[string]struct{}
	crawledUrlsMutex sync.RWMutex
	sitemapUrlsMutex sync.RWMutex
}

func NewSitemapBuilder() *SitemapBuilder {
	return &SitemapBuilder{
		crawledUrls:      map[string]struct{}{},
		sitemapUrls:      map[string]struct{}{},
		crawledUrlsMutex: sync.RWMutex{},
		sitemapUrlsMutex: sync.RWMutex{},
	}
}

var errAlreadyCrawled = errors.New("already crawled URL")
var errAlreadyInSitemap = errors.New("already added URL to sitemap")

func (s *SitemapBuilder) AddToCrawledUrls(url string) error {
	//s.crawledUrlsMutex.Lock()
	//defer s.crawledUrlsMutex.Unlock()
	//return s.addToMap(url, s.crawledUrls, errAlreadyCrawled)

	s.crawledUrlsMutex.RLock()
	_, exists := s.crawledUrls[url]
	s.crawledUrlsMutex.RUnlock()
	if !exists {
		s.crawledUrlsMutex.Lock()
		s.crawledUrls[url] = struct{}{}
		s.crawledUrlsMutex.Unlock()
		return nil
	}
	return errAlreadyCrawled
}

func (s *SitemapBuilder) AddToSitemap(url string) error {
	//s.sitemapUrlsMutex.Lock()
	//defer s.sitemapUrlsMutex.Unlock()
	//return s.addToMap(url, s.sitemapUrls, errAlreadyInSitemap)

	s.sitemapUrlsMutex.RLock()
	_, exists := s.sitemapUrls[url]
	s.sitemapUrlsMutex.RUnlock()
	if !exists {
		s.sitemapUrlsMutex.Lock()
		s.sitemapUrls[url] = struct{}{}
		s.sitemapUrlsMutex.Unlock()
		return nil
	}
	return errAlreadyCrawled
}

func (s *SitemapBuilder) addToMap(url string, urlMap map[string]struct{}, err error) error {
	_, exists := urlMap[url]
	if !exists {
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
