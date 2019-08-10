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
			"verify valid URL without protocolCapGroup",
			"monzo.com",
			"monzo.com",
		},
		{
			"verify valid URL http protocolCapGroup",
			"http://monzo.com",
			"http://monzo.com",
		},
		{
			"verify valid URL https protocolCapGroup",
			"https://monzo.com",
			"https://monzo.com",
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
			urlManipulator := NewUrlParser()
			err := urlManipulator.VerifyBaseUrl(test.URL)
			assertNoError(t, err)
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
			errInvalidBaseUrl,
		},
		{
			"fail URL with smtp protocol",
			"smtp://monzo.com",
			errInvalidBaseUrl,
		},
		{
			"fail URL with ftp protocol",
			"ftp://monzo.com",
			errInvalidBaseUrl,
		},
		{
			"fail URL without domain extension",
			"monzo.",
			errInvalidBaseUrl,
		},
		{
			"fail URL that's empty",
			"",
			errInvalidBaseUrl,
		},
		{
			"fail URL that's not the root",
			"monzo.com/help",
			errInvalidBaseUrl,
		},
		{
			"fail URL that's not alpha numeric",
			"sd0_93hj$£%^.com/help",
			errInvalidBaseUrl,
		},
	}

	for _, test := range errorCases {
		t.Run(test.Name, func(t *testing.T) {
			urlManipulator := NewUrlParser()
			err := urlManipulator.VerifyBaseUrl(test.URL)
			assertErrorMessage(t, err, errInvalidBaseUrl.Error())
		})
	}
}

func TestIsSameDomain(t *testing.T) {
	isSameDomainNoErrorCases := []struct {
		Name       string
		Hostname   string
		UrlToCheck string
		Want       string
	}{
		{
			"CheckSameDomain absolute path returns no error",
			"https://www.monzo.com",
			"https://www.monzo.com/help",
			"https://www.monzo.com/help",
		},
	}

	for _, test := range isSameDomainNoErrorCases {
		t.Run(test.Name, func(t *testing.T) {
			urlManipulator := NewUrlParser()
			err := urlManipulator.CheckSameDomain(test.UrlToCheck, test.Hostname)
			assertNoError(t, err)
		})
	}

	isSameDomainErrorCases := []struct {
		Name       string
		Hostname   string
		UrlToCheck string
	}{
		{
			"CheckSameDomain different domain returns error",
			"https://www.monzo.com",
			"https://www.monzo.co.uk/help",
		},
		{
			"CheckSameDomain relative homepage returns error",
			"https://www.monzo.com",
			"/",
		},
		{
			"CheckSameDomain empty returns error",
			"https://www.monzo.com",
			"",
		},
		{
			"CheckSameDomain relative path returns no error",
			"https://www.monzo.com",
			"/help",
		},
	}

	for _, test := range isSameDomainErrorCases {
		t.Run(test.Name, func(t *testing.T) {
			urlManipulator := NewUrlParser()
			err := urlManipulator.CheckSameDomain(test.UrlToCheck, test.Hostname)
			assertErrorMessage(t, err, errDifferentDomain.Error())
		})
	}
}
