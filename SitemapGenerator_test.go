package main

import (
	"reflect"
	"testing"
)

func TestSitemapGenerator(t *testing.T) {
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
			sitemapGenerator := SitemapGenerator{}

			for _, link := range test.Input {
				sitemapGenerator.addToSitemap(link)
			}
			got := sitemapGenerator.returnSitemap()

			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("sitemap not in expected format. Got %v, wanted %v", got, test.Want)
			}
		})
	}
}
