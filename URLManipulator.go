package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	protocol = "protocol"
	hostname = "hostnameWithProtocol"
)

var (
	errInvalidUrl = errors.New("invalid url")
	validUrlRegex = fmt.Sprintf(`(?i)^(?P<%v>(https?|ftp|smtp)\:\/\/)?(?P<%v>[[:alnum:]]+\.[[:alnum:]]+(?:\.[[:alnum:]]+)?)$`, protocol, hostname)
)

func verifyHostname(url string) (string, error) {
	regex := regexp.MustCompile(validUrlRegex)
	res := regex.FindStringSubmatch(url)
	validUrl := false

	var sb strings.Builder
	for i, name := range regex.SubexpNames() {
		if (name == hostname || name == protocol) && res != nil && i < len(res) {
			if name == hostname {
				validUrl = true
			}
			sb.WriteString(res[i])
		}
	}
	if validUrl {
		return sb.String(), nil
	}
	return "", errInvalidUrl
}

func addHttpsIfNecessary(url string) string {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		return "https://" + url
	}
	return url
}

// Adds hostnameWithProtocol with protocol to relative URLs
func addHostnameAndProtocolToRelativeUrls(url, hostnameWithProtocol string) string {
	if len(url) > 1 && url[0] == '/' {
		url = hostnameWithProtocol + url
	}
	return url
}

// URLs passed into this will always have a hostnameWithProtocol prefix
func checkSameDomain(url, hostname string) error {
	url = addHostnameAndProtocolToRelativeUrls(url, hostname)
	// TODO: strip protocols on both strings for comparison
	//url = stripProtocol(url)
	if strings.HasPrefix(url, hostname) {
		return nil
	}
	return errDifferentDomain
}
