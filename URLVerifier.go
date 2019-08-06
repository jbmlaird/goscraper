package main

import "regexp"

const validUrlRegex = `(?i)^(?:(https?|ftp|smtp)\:\/\/)?([[:alnum:]]+\.[[:alnum:]]+(?:\.[[:alnum:]]+)?)$`

// This could be expanded to contain verification on the string and return appropriate error messages
// depending on how it has been malformed
func verifyUrl(url string) (bool, error) {
	return regexp.MatchString(validUrlRegex, url)
}
