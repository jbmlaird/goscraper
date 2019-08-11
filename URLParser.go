package main

import (
	"errors"
	"fmt"
	"regexp"
)

type URLParser struct {
	urlRegex *regexp.Regexp
}

func NewUrlParser() *URLParser {
	return &URLParser{
		regexp.MustCompile(validUrlRegex),
	}
}

const (
	protocolCapGroup       = "protocolCapGroup"
	hostnameCapGroup       = "hostnameCapGroup"
	hostnameOnwardsCapGoup = "hostnameOnwardsCapGoup"
)

var (
	errInvalidBaseUrl    = errors.New("invalid base url. Only http/https and base URLs are supported")
	protocolRegex        = fmt.Sprintf(`(?P<%v>https?\:\/\/)?`, protocolCapGroup)
	hostnameRegex        = fmt.Sprintf(`(?P<%v>(www.)?[[:alnum:]]+\.[[:alnum:]]+(?:\.[[:alnum:]]+)?)`, hostnameCapGroup)
	hostnameOnwardsRegex = fmt.Sprintf(`(?P<%v>.*)`, hostnameOnwardsCapGoup)
	validUrlRegex        = fmt.Sprintf(`(?i)^%v%v\/?%v?$`, protocolRegex, hostnameRegex, hostnameOnwardsRegex)
)

func (u *URLParser) VerifyBaseUrl(url string) error {
	baseUrlRes := u.urlRegex.FindStringSubmatch(url)
	validBaseUrl := false

	if baseUrlRes != nil {
		for i, name := range u.urlRegex.SubexpNames() {
			// Shouldn't need the final check. I think my regex is misbehaving
			if name == hostnameOnwardsCapGoup && i < len(baseUrlRes) && baseUrlRes[i] != "" {
				return errInvalidBaseUrl
			}
			if name == hostnameCapGroup && i < len(baseUrlRes) {
				validBaseUrl = true
			}
		}
	}
	if validBaseUrl {
		return nil
	}
	return errInvalidBaseUrl
}

func (u *URLParser) VerifyProtocol(url string) error {
	baseUrlRes := u.urlRegex.FindStringSubmatch(url)
	validProtocol := false

	for i, name := range u.urlRegex.SubexpNames() {
		if name == protocolCapGroup && baseUrlRes != nil && i < len(baseUrlRes) {
			validProtocol = true
		}
	}
	if validProtocol {
		return nil
	}
	return errUnsupportedProtocol
}

func (u *URLParser) CheckSameDomain(url, baseUrl string) error {
	urlHostname := u.findHostname(url)
	baseUrlHostname := u.findHostname(baseUrl)

	if urlHostname == baseUrlHostname {
		return nil
	}
	return errDifferentDomain
}

func (u *URLParser) findHostname(url string) (extractedHostname string) {
	urlRes := u.urlRegex.FindStringSubmatch(url)
	for i, name := range u.urlRegex.SubexpNames() {
		if (name == hostnameCapGroup) && urlRes != nil && i < len(urlRes) {
			return urlRes[i]
		}
	}
	return
}
