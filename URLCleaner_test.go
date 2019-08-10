package main

import "testing"

func TestStripTrailingSlash(t *testing.T) {
	cases := []struct {
		Name  string
		Input string
		Want  string
	}{
		{
			"trailing slash, strip trailing slash",
			"https://www.monzo.com/",
			"https://www.monzo.com",
		},
		{
			"non-trailing slash returns",
			"https://www.monzo.com",
			"https://www.monzo.com",
		},
		{
			"empty string doesn't throw an error and returns",
			"",
			"",
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			got := stripTrailingSlash(test.Input)
			assertStringOutput(t, got, test.Want)
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
			"adds https:// to blank protocolCapGroup",
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

func TestAddHostnameAndProtocolToRelativeUrls(t *testing.T) {
	successCases := []struct {
		Name                 string
		HostnameWithProtocol string
		URL                  string
		Want                 string
	}{
		{
			"relative path home, returns slash",
			"https://www.google.co.uk",
			"/",
			"/",
		},
		{
			"relative path, returns appended hostname",
			"https://www.google.co.uk",
			"/home",
			"https://www.google.co.uk/home",
		},
	}

	for _, test := range successCases {
		t.Run(test.Name, func(t *testing.T) {
			got := addHostnameAndProtocolToRelativeUrls(test.URL, test.HostnameWithProtocol)
			assertStringOutput(t, got, test.Want)
		})
	}
}

func TestStripAfterSeparator(t *testing.T) {
	stripAfterSeparatorCases := []struct {
		Name      string
		Separator string
		Input     string
		Want      string
	}{
		{
			"strip query",
			"?",
			"https://www.monzo.co.uk?search=test",
			"https://www.monzo.co.uk",
		},
		{
			"strip anchor",
			"#",
			"https://www.monzo.co.uk#over-here",
			"https://www.monzo.co.uk",
		},
		{
			"strip non-existent character, nothing",
			"%",
			"https://www.monzo.co.uk/yeah-dawg",
			"https://www.monzo.co.uk/yeah-dawg",
		},
		{
			"strip single character with character, nothing",
			"/",
			"/",
			"",
		},
	}

	for _, test := range stripAfterSeparatorCases {
		t.Run(test.Name, func(t *testing.T) {
			got := stripAfterSeparator(test.Input, test.Separator)
			assertStringOutput(t, got, test.Want)
		})
	}
}

func TestCleanUrl(t *testing.T) {
	cleanUrlSuccessCases := []struct {
		Name    string
		BaseUrl string
		Input   string
		Want    string
	}{
		{"CleanUrl adds https:// to BaseUrl",
			"google.co.uk",
			"google.co.uk",
			"https://google.co.uk",
		},
		{"CleanUrl adds base url to relative path",
			"condenastint.com",
			"/help",
			"https://condenastint.com/help",
		},
		{"CleanUrl doesn't add https:// to https protocol",
			"https://monzo.com",
			"https://monzo.com",
			"https://monzo.com",
		},
		{"CleanUrl doesn't add https:// to http protocol",
			"https://monzo.com",
			"http://monzo.com/help",
			"http://monzo.com/help",
		},
	}

	for _, test := range cleanUrlSuccessCases {
		t.Run(test.Name, func(t *testing.T) {
			got, err := CleanUrl(test.Input, test.BaseUrl)
			assertNoError(t, err)
			assertStringOutput(t, got, test.Want)
		})
	}

	cleanUrlFailureCases := []struct {
		Name    string
		Input   string
		Error error
	}{
		{"pass single character, throw invalid URL error",
			"/",
			errInvalidUrl,
		},
		{"pass empty string, throw invalid URL error",
			"",
			errInvalidUrl,
		},
		{"pass anchor URL, throw path or query error",
			"#some-stuff",
			errPathOrQuery,
		},
		{"pass query string, throw error",
			"?s=fintech_is_cool",
			errPathOrQuery,
		},
		{"pass SMTP, throw unsupported protocol error",
			"smtp://some-website.com",
			errUnsupportedProtocol,
		},
		{"pass FTP, throw unsupported protocol error",
			"ftp://some-other-website.com",
			errUnsupportedProtocol,
		},
	}

	for _, test := range cleanUrlFailureCases {
		t.Run(test.Name, func(t *testing.T) {
			got, err := CleanUrl(test.Input, "www.monzo.com")
			assertErrorMessage(t, err, test.Error.Error())
			assertStringOutput(t, got, "")
		})
	}
}
