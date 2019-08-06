package main

import "testing"

func TestVerifyUrl(t *testing.T) {
	cases := []struct {
		Name    string
		URL     string
		IsValid bool
	}{
		{
			"verify valid URL without protocol",
			"monzo.com",
			true,
		},
		{
			"verify valid URL http protocol",
			"http://monzo.com",
			true,
		},
		{
			"verify valid URL https protocol",
			"https://monzo.com",
			true,
		},
		{
			"verify valid URL ftp protocol",
			"ftp://monzo.com",
			true,
		},
		{
			"verify valid URL smtp protocol",
			"smtp://monzo.com",
			true,
		},
		{
			"verify with double domain extension",
			"https://monzo.co.uk",
			true,
		},
		{
			"verify, ignoring case",
			"HTTPS://monzo.co.uk",
			true,
		},
		{
			"fail URL with made up protocol",
			"monzo://monzo.com",
			false,
		},
		{
			"fail URL without domain extension",
			"monzo.",
			false,
		},
		{
			"fail URL that's not the hostname",
			"monzo.com/help",
			false,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got, err := verifyUrl(test.URL)
			assertNoError(t, err)
			assertOutput(t, got, test.IsValid)
		})
	}
}
