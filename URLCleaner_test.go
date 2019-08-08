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
	}

	for _, test := range stripAfterSeparatorCases {
		t.Run(test.Name, func(t *testing.T) {
			got := stripAfterSeparator(test.Input, test.Separator)
			assertStringOutput(t, got, test.Want)
		})
	}
}

func TestCleanUrl(t *testing.T) {
	testCleanUrlCases := []struct {
		Name    string
		BaseUrl string
		Input   string
		Want    string
	}{
		{"cleanUrl adds https:// to BaseUrl",
			"google.co.uk",
			"google.co.uk",
			"https://google.co.uk",
		},
		{"cleanUrl adds base url to relative path",
			"condenastint.com",
			"/help",
			"https://condenastint.com/help",
		},
		{"cleanUrl doesn't add https:// to https protocol",
			"https://monzo.com",
			"https://monzo.com",
			"https://monzo.com",
		},
		{"cleanUrl doesn't add https:// to http protocol",
			"https://monzo.com",
			"http://monzo.com/help",
			"http://monzo.com/help",
		},
	}

	for _, test := range testCleanUrlCases {
		t.Run(test.Name, func(t *testing.T) {
			got, err := cleanUrl(test.Input, test.BaseUrl)
			assertNoError(t, err)
			assertStringOutput(t, got, test.Want)
		})
	}

	t.Run("pass single character, throw error", func(t *testing.T) {
		got, err := cleanUrl("/", "www.monzo.com")
		assertErrorMessage(t, err, errSingleCharacter.Error())
		assertStringOutput(t, got, "")
	})

	t.Run("pass empty string, throw error", func(t *testing.T) {
		got, err := cleanUrl("", "www.monzo.com")
		assertErrorMessage(t, err, errSingleCharacter.Error())
		assertStringOutput(t, got, "")
	})
}
