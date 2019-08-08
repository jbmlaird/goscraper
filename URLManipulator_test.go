package main

import (
	"testing"
)

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
		{
			"verify, removing trailing slash",
			"HTTPS://monzo.co.uk/",
			"HTTPS://monzo.co.uk",
		},
	}

	for _, test := range successCases {
		t.Run(test.Name, func(t *testing.T) {
			urlManipulator := NewUrlManipulator()
			got, err := urlManipulator.verifyHostname(test.URL)
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
			urlManipulator := NewUrlManipulator()
			_, err := urlManipulator.verifyHostname(test.URL)
			assertErrorMessage(t, err, errInvalidUrl.Error())
		})
	}
}

func TestAddHttpsIfNecessary(t *testing.T) {
	cases := []struct {
		Name  string
		Input string
		Want  string
	}{
		{
			"adds https:// to blank protocol",
			"google.com",
			"https://google.com",
		},
		{
			"doesn't add https:// to https://",
			"https://google.com",
			"https://google.com",
		},
		{
			"doesn't add https:// to http://",
			"http://google.com",
			"http://google.com",
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := addHttpsIfNecessary(test.Input)
			assertStringOutput(t, got, test.Want)
		})
	}
}

func TestIsSameDomain(t *testing.T) {
	isSameDomainNoErrorCases := []struct {
		Name       string
		Hostname   string
		UrlToCheck string
	}{
		{"checkSameDomain absolute path returns true",
			"https://www.monzo.com",
			"https://www.monzo.com/help",
		},
		{"checkSameDomain relative path returns true",
			"https://www.monzo.com",
			"/help",
		},
	}

	for _, test := range isSameDomainNoErrorCases {
		t.Run(test.Name, func(t *testing.T) {
			urlManipulator := NewUrlManipulator()
			assertNoError(t, urlManipulator.checkSameDomain(test.UrlToCheck, test.Hostname))
		})
	}

	isSameDomainErrorCases := []struct {
		Name       string
		Hostname   string
		UrlToCheck string
	}{
		{"checkSameDomain different domain returns false",
			"https://www.monzo.com",
			"https://www.monzo.co.uk/help",
		},
		{"checkSameDomain homepage returns false",
			"https://www.monzo.com",
			"/",
		},
		{"checkSameDomain empty returns false",
			"https://www.monzo.com",
			"",
		},
	}

	for _, test := range isSameDomainErrorCases {
		t.Run(test.Name, func(t *testing.T) {
			urlManipulator := NewUrlManipulator()
			assertErrorMessage(t, urlManipulator.checkSameDomain(test.UrlToCheck, test.Hostname), errDifferentDomain.Error())
		})
	}
}
