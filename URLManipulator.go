package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type URLManipulator struct {
	urlRegex *regexp.Regexp
}

func NewUrlManipulator() *URLManipulator {
	return &URLManipulator{
		regexp.MustCompile(validUrlRegex),
	}
}

const (
	protocol = "protocol"
	hostname = "hostnameWithProtocol"
)

var (
	errInvalidUrl = errors.New("invalid url")
	protocolRegex = fmt.Sprintf(`(?i)^(?P<%v>(https?|ftp|smtp)\:\/\/)?`, protocol)
	hostnameRegex = fmt.Sprintf(`(?P<%v>[[:alnum:]]+\.[[:alnum:]]+(?:\.[[:alnum:]]+)?)`, hostname)
	validUrlRegex = fmt.Sprintf(`%v%v\/?$`, protocolRegex, hostnameRegex)
)

func (u *URLManipulator) verifyHostname(url string) (string, error) {
	res := u.urlRegex.FindStringSubmatch(url)
	validUrl := false

	var sb strings.Builder
	for i, name := range u.urlRegex.SubexpNames() {
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
func (u *URLManipulator) checkSameDomain(url, hostname string) error {
	urlRes := u.urlRegex.FindStringSubmatch(url)

	var urlHostname string
	var hostnameHostname string

	for i, name := range u.urlRegex.SubexpNames() {
		if (name == hostname) && urlRes != nil && i < len(urlRes) {
			if name == hostname {
				urlHostname = urlRes[i]
			}
		}
	}

	hostnameRes := u.urlRegex.FindStringSubmatch(hostname)
	for i, name := range u.urlRegex.SubexpNames() {
		if (name == hostname) && urlRes != nil && i < len(urlRes) {
			if name == hostname {
				hostnameHostname = hostnameRes[i]
			}
		}
	}

	if urlHostname == hostnameHostname {
		return nil
	}
	return errDifferentDomain
}
