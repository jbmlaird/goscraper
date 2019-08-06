package main

import "testing"

func TestVerifyUrl(t *testing.T) {
	successCases := []struct {
		Name     string
		URL      string
		Hostname string
	}{
		{
			"verify valid URL without protocol",
			"monzo.com",
			"monzo.com",
		},
		{
			"verify valid URL http protocol",
			"http://monzo.com",
			"monzo.com",
		},
		{
			"verify valid URL https protocol",
			"https://monzo.com",
			"monzo.com",
		},
		{
			"verify valid URL ftp protocol",
			"ftp://monzo.com",
			"monzo.com",
		},
		{
			"verify valid URL smtp protocol",
			"smtp://monzo.com",
			"monzo.com",
		},
		{
			"verify with double domain extension",
			"https://monzo.co.uk",
			"monzo.co.uk",
		},
		{
			"verify, ignoring case",
			"HTTPS://monzo.co.uk",
			"monzo.co.uk",
		},
	}

	for _, test := range successCases {
		t.Run(test.Name, func(t *testing.T) {
			got, err := verifyUrl(test.URL)
			assertNoError(t, err)
			assertOutput(t, got, test.Hostname)
		})
	}

	errorCases := []struct {
		Name     string
		URL      string
		Hostname error
	}{
		{
			"fail URL with made up protocol",
			"monzo://monzo.com",
			errInvalidRegex,
		},
		{
			"fail URL without domain extension",
			"monzo.",
			errInvalidRegex,
		},
		{
			"fail URL that's not the root",
			"monzo.com/help",
			errInvalidRegex,
		},
	}

	for _, test := range errorCases {
		t.Run(test.Name, func(t *testing.T) {
			_, err := verifyUrl(test.URL)
			assertErrorMessage(t, err, errInvalidRegex.Error())
		})
	}
}
