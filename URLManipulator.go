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
	protocolCapGroup = "protocolCapGroup"
	hostnameCapGroup = "hostnameWithProtocol"
	pathCapGroup     = "pathCapGroup"
)

var (
	errInvalidBaseUrl = errors.New("invalid base url")
	protocolRegex     = fmt.Sprintf(`(?P<%v>(https?|ftp|smtp)\:\/\/)?`, protocolCapGroup)
	hostnameRegex     = fmt.Sprintf(`(?P<%v>[[:alnum:]]+\.[[:alnum:]]+(?:\.[[:alnum:]]+)?)`, hostnameCapGroup)
	pathRegex         = fmt.Sprintf(`(?P<%v>.*)`, pathCapGroup)
	validUrlRegex     = fmt.Sprintf(`(?i)^%v%v\/?%v$`, protocolRegex, hostnameRegex, pathRegex)
)

func (u *URLManipulator) verifyBaseUrl(url string) (string, error) {
	baseUrlRes := u.urlRegex.FindStringSubmatch(url)
	validBaseUrl := false

	var sb strings.Builder
	for i, name := range u.urlRegex.SubexpNames() {
		// Shouldn't need the final check. I think my regex is misbehaving
		if name == pathCapGroup && baseUrlRes != nil && i < len(baseUrlRes) && baseUrlRes[i] != "" {
			return "", errInvalidBaseUrl
		}
		if (name == hostnameCapGroup || name == protocolCapGroup) && baseUrlRes != nil && i < len(baseUrlRes) {
			if name == hostnameCapGroup {
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

func (u *URLManipulator) checkSameDomain(url, baseUrl string) error {
	urlHostname := u.findHostname(url)
	baseUrlHostname := u.findHostname(baseUrl)

	if urlHostname == baseUrlHostname {
		return nil
	}
	return errDifferentDomain
}

func (u *URLManipulator) findHostname(url string) (extractedHostname string) {
	urlRes := u.urlRegex.FindStringSubmatch(url)
	for i, name := range u.urlRegex.SubexpNames() {
		if (name == hostnameCapGroup) && urlRes != nil && i < len(urlRes) {
			extractedHostname = urlRes[i]
		}
	}
	return
}
