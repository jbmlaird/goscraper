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

// This could be expanded to contain verification on the string and return appropriate error messages
// depending on how it has been malformed
func verifyUrl(url string) (string, error) {
	regex := regexp.MustCompile(validUrlRegex)
	res := regex.FindStringSubmatch(url)

	for i, name := range regex.SubexpNames() {
		if name == hostname && res != nil && i < len(res) {
			return res[i], nil
		}
	}
	return "", errInvalidRegex
}
