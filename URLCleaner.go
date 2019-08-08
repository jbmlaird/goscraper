package main

import "strings"

// Adds hostnameWithProtocol with protocolCapGroup to relative URLs
func addHostnameAndProtocolToRelativeUrls(url, hostnameWithProtocol string) string {
	if len(url) > 1 && url[0] == '/' {
		url = hostnameWithProtocol + url
	}
	return url
}

func stripTrailingSlash(url string) string {
	if len(url) > 1 && url[len(url)-1] == '/' {
		url = url[0 : len(url)-1]
	}
	return url
}

func stripAfterSeparator(url, separator string) string {
	return strings.Split(url, separator)[0]
}

func addHttpsIfNecessary(url string) string {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return "https://" + url
	}
	return url
}

func cleanUrl(url, baseUrl string) (string, error) {
	if len(url) <= 1 {
		return "", errSingleCharacter
	}
	url = stripTrailingSlash(url)
	url = stripAfterSeparator(url, "?") // strip queries
	url = stripAfterSeparator(url, "#") // strip anchors
	url = addHostnameAndProtocolToRelativeUrls(url, baseUrl)
	url = addHttpsIfNecessary(url)
	return url, nil
}
