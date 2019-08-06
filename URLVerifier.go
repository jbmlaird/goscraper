package main

import (
	"errors"
	"fmt"
	"regexp"
)
const hostname = "hostname"
var (
	errInvalidRegex = errors.New("invalid regex")
	validUrlRegex = fmt.Sprintf(`(?i)^(?:(https?|ftp|smtp)\:\/\/)?(?P<%v>[[:alnum:]]+\.[[:alnum:]]+(?:\.[[:alnum:]]+)?)$`, hostname)
)

func verifyHostname(url string) (string, error) {
	regex := regexp.MustCompile(validUrlRegex)
	res := regex.FindStringSubmatch(url)

	for i, name := range regex.SubexpNames() {
		if name == hostname && res != nil && i < len(res) {
			return res[i], nil
		}
	}
	return "", errInvalidRegex
}

// TODO: Extract hostname function?
func isSubdomain(hostname, subdomain string) bool {
	return false
}