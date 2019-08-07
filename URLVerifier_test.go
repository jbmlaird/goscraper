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
			"http://monzo.com",
		},
		{
			"verify valid URL https protocol",
			"https://monzo.com",
			"https://monzo.com",
		},
		{
			"verify valid URL ftp protocol",
			"ftp://monzo.com",
			"ftp://monzo.com",
		},
		{
			"verify valid URL smtp protocol",
			"smtp://monzo.com",
			"smtp://monzo.com",
		},
		{
			"verify with double domain extension",
			"https://monzo.co.uk",
			"https://monzo.co.uk",
		},
		{
			"verify, ignoring case",
			"HTTPS://monzo.co.uk",
			"HTTPS://monzo.co.uk",
		},
	}

	for _, test := range successCases {
		t.Run(test.Name, func(t *testing.T) {
			got, err := verifyHostname(test.URL)
			assertNoError(t, err)
			assertStringOutput(t, got, test.Hostname)
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
			errInvalidUrl,
		},
		{
			"fail URL without domain extension",
			"monzo.",
			errInvalidUrl,
		},
		{
			"fail URL that's not the root",
			"monzo.com/help",
			errInvalidUrl,
		},
	}

	for _, test := range errorCases {
		t.Run(test.Name, func(t *testing.T) {
			_, err := verifyHostname(test.URL)
			assertErrorMessage(t, err, errInvalidUrl.Error())
		})
	}
}
