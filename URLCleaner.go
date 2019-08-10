package main

import (
	"github.com/pkg/errors"
	"strings"
)

var errPathOrQuery = errors.New("URL is either a path or a query")
var errUnsupportedProtocol = errors.New("only http and https is supported")

func CleanUrl(url, baseUrl string) (string, error) {
	if len(url) <= 1 {
		return "", errInvalidUrl
	} else if url[0] == '#' || url[0] == '?' {
		return "", errPathOrQuery
	} else if strings.HasPrefix(url, "ftp://") || strings.HasPrefix(url, "smtp://") {
		return "", errUnsupportedProtocol
	}
	url = stripTrailingSlash(url)
	url = stripAfterSeparator(url, "?") // strip queries
	url = stripAfterSeparator(url, "#") // strip anchors
	url = addHostnameAndProtocolToRelativeUrls(url, baseUrl)
	url = addHttpsIfNecessary(url)
	return url, nil
}

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

// When comparing links, the hostname is compared with the protocol ignored to prevent duplicates from being added
func addHttpsIfNecessary(url string) string {
	if strings.HasPrefix(url, "http://") {
		return strings.Replace(url, "http://", "https://", 1)
	}
	if !strings.HasPrefix(url, "https://") && !strings.Contains(url, "://") {
		return "https://" + url
	}
	return url
}
