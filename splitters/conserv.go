package splitters

import (
	"regexp"
)

// Conserv represents a conservative splitter.
type Conserv struct {
	Splitter
}

// Split on Conserv receives a token and returns an array of hard-words.
func (c Conserv) Split(token string) ([]string, error) {
	//remove numbers
	numRegex := regexp.MustCompile("([0-9]+)")
	t := numRegex.ReplaceAllString(token, "_${1}_")

	//separate camelcase through underscores
	ccRegex := regexp.MustCompile("([a-z])([A-Z])")
	cc := ccRegex.ReplaceAllString(t, "${1}_$2")

	endCcRegex := regexp.MustCompile("([A-Z])([a-z])")
	rcc := endCcRegex.ReplaceAllString(cc, "_${1}${2}")

	regex := regexp.MustCompile("_")
	s := regex.Split(rcc, -1)

	//remove the empty items
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}

	return r, nil
}
