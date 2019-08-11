package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSitemapGenerator(t *testing.T) {
	t.Run("inputting string that already exists in sitemap URL map throws already in sitemap error", func(t *testing.T) {
		link := "https://www.monzo.com"
		sitemapBuilder := NewSitemapBuilder()

		err := sitemapBuilder.AddToCrawledUrls(link)
		assertNoError(t, err)
		err = sitemapBuilder.AddToCrawledUrls(link)
		assertErrorMessage(t, err, errAlreadyCrawled.Error())
	})

	t.Run("inputting string that already exists in crawled URLs map throws already crawled error", func(t *testing.T) {
		link := "https://www.monzo.com"
		sitemapBuilder := NewSitemapBuilder()

		err := sitemapBuilder.AddToSitemap(link)
		assertNoError(t, err)
		err = sitemapBuilder.AddToSitemap(link)
		assertErrorMessage(t, err, errAlreadyInSitemap.Error())
	})

	cases := []struct{
		Name string
		Input []string
		Want []string
	}{
		{
			"output sorted sitemap from links",
			[]string{
				"/help",
				"/help/faq",
				"/help/faq/question-one",
				"/search",
				"/about-us",
				"/contact-us",
			},
			[]string{
				"/about-us",
				"/contact-us",
				"/help",
				"/help/faq",
				"/help/faq/question-one",
				"/search",
			},
		},
		{
			"don't input the same link",
			[]string{
				"/help",
				"/help/faq",
				"/help/faq/question-one",
				"/help",
				"/help/faq",
				"/help/faq/question-two",
				"/help",
				"/help/faq",
				"/help/faq/question-three",
				"/help",
				"/help/faq",
				"/help/faq/question-four",
			},
			[]string{
				"/help",
				"/help/faq",
				"/help/faq/question-four",
				"/help/faq/question-one",
				"/help/faq/question-three",
				"/help/faq/question-two",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			sitemapBuilder := NewSitemapBuilder()

			for _, link := range test.Input {
				_ = sitemapBuilder.AddToSitemap(link)
			}

			got := sitemapBuilder.BuildSitemap()
			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("sitemap not in expected format. Got %v, wanted %v", got, test.Want)
			}
		})
	}

	// Specifically execute this with go test -race
	t.Run("multiple goroutines do not cause race conditions reading and writing to sitemap and crawled URLs", func(t *testing.T) {
		sitemapGenerator := NewSitemapBuilder()

		links := []string{
			"/help",
			"/help/faq",
			"/help/faq/question-one",
			"/help/faq/question-two",
			"/help/faq/question-three",
			"/help/faq/question-four",
		}

		for i:=0; i<1000; i++ {
			go func() {
				for _, link := range links {
					_ = sitemapGenerator.AddToSitemap(link)
				}
			}()

			go func() {
				for _, link := range links {
					_ = sitemapGenerator.AddToCrawledUrls(link)
				}
			}()
		}

		got := sitemapGenerator.BuildSitemap()
		want := []string{
			"/help",
			"/help/faq",
			"/help/faq/question-four",
			"/help/faq/question-one",
			"/help/faq/question-three",
			"/help/faq/question-two",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("sitemap not in expected format. Got %v, wanted %v", got, want)
		}
	})
}

// This was so I could test 3 cases:
//  1. Using the same mutex for both maps - 2747832 ns/op
//  2. Using two separate mutexes for each of the maps - 1872559 ns/op
//  3. Using two separate RWMutexes. RLocking/RUnlocking when checking the map and Locking/Unlocking when writing - 1893284 ns/op
//
// Case 2 came out the fastest
//
// Case 3 was narrowly longer so that may actually be the best implementation if one wanted to be very clear about locking
func BenchmarkSitemapBuilder_AddingToMapsWithMutexes(b *testing.B) {
	sitemapBuilder := NewSitemapBuilder()
	for n := 0; n < b.N; n++ {
		for i:=0; i<1000; i++ {
			link := fmt.Sprintf("link%v", i)
			go func(link string) {
				_ = sitemapBuilder.AddToSitemap(link)
			}(link)
			go func(link string) {
				_ = sitemapBuilder.AddToCrawledUrls(link)
			}(link)
		}
		for i:=1000; i>0; i-- {
			link := fmt.Sprintf("link%v", i)
			go func(link string) {
				_ = sitemapBuilder.AddToSitemap(link)
			}(link)
			go func(link string) {
				_ = sitemapBuilder.AddToCrawledUrls(link)
			}(link)
		}
	}
}
