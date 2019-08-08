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
		//regexp.MustCompile(oldValidUrlRegex),
	}
}

const (
	protocol = "protocol"
	hostname = "hostnameWithProtocol"
	path     = "path"
)

var (
	errInvalidBaseUrl = errors.New("invalid base url")
	protocolRegex     = fmt.Sprintf(`(?P<%v>(https?|ftp|smtp)\:\/\/)?`, protocol)
	hostnameRegex     = fmt.Sprintf(`(?P<%v>[[:alnum:]]+\.[[:alnum:]]+(?:\.[[:alnum:]]+)?)`, hostname)
	pathRegex         = fmt.Sprintf(`(?P<%v>\/.*)?`, path)
	validUrlRegex     = fmt.Sprintf(`(?i)^%v%v\/?%v$`, protocolRegex, hostnameRegex, pathRegex)
)

func (u *URLManipulator) verifyBaseUrl(url string) (string, error) {
	baseUrlRes := u.urlRegex.FindStringSubmatch(url)
	validBaseUrl := false

	var sb strings.Builder
	for i, name := range u.urlRegex.SubexpNames() {
		// Shouldn't need the final check. I think my regex is misbehaving
		if name == path && baseUrlRes != nil && i < len(baseUrlRes) && baseUrlRes[i] != "" {
			return "", errInvalidBaseUrl
		}
		if (name == hostname || name == protocol) && baseUrlRes != nil && i < len(baseUrlRes) {
			if name == hostname {
				validBaseUrl = true
			}
			sb.WriteString(baseUrlRes[i])
		}
	}
	if validBaseUrl {
		return sb.String(), nil
	}
	return "", errInvalidBaseUrl
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
func (u *URLManipulator) checkSameDomain(url, baseUrl string) (modifiedUrl string, err error) {
	modifiedUrl = addHostnameAndProtocolToRelativeUrls(url, baseUrl)
	modifiedUrl = addHttpsIfNecessary(modifiedUrl)

	urlHostname := u.findHostname(modifiedUrl)
	baseUrlHostname := u.findHostname(baseUrl)

	if urlHostname == baseUrlHostname {
		return modifiedUrl, nil
	}
	return "", errDifferentDomain
}

func (u *URLManipulator) findHostname(url string) (extractedHostname string) {
	urlRes := u.urlRegex.FindStringSubmatch(url)
	for i, name := range u.urlRegex.SubexpNames() {
		if (name == hostname) && urlRes != nil && i < len(urlRes) {
			extractedHostname = urlRes[i]
		}
	}
	return
}
