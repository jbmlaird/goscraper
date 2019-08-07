package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)
const (
	protocol = "protocol"
	hostname = "hostname"
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
