package splitters

import (
	"regexp"
)

// Conserv represents a conservative splitter.
type Conserv struct {
	Splitter
}

// Split on Conserv receives a token and returns an array of hard/soft words,
// splitted by:
// * Underscores
// * Numbers
// * CamelCase.
func (c Conserv) Split(token string) ([]string, error) {
	//add markers between letters and numbers
	leadingNumRegex := regexp.MustCompile("([a-zA-Z])([0-9])")
	processedToken := leadingNumRegex.ReplaceAllString(token, "${1}_$2")

	//add markers between numbers and letters
	trailingNumRegex := regexp.MustCompile("([0-9])([a-zA-Z])")
	processedToken = trailingNumRegex.ReplaceAllString(processedToken, "${1}_$2")

	//add markers for lower to upper case camel-case combination
	lowerToUpperRegex := regexp.MustCompile("([a-z])([A-Z])")
	processedToken = lowerToUpperRegex.ReplaceAllString(processedToken, "${1}_$2")

	//add markers for upper to lower case camel-case combination
	upperToLowerRegex := regexp.MustCompile("([A-Z]+)([A-Z])([a-z])")
	processedToken = upperToLowerRegex.ReplaceAllString(processedToken, "${1}_$2$3")

	regex := regexp.MustCompile("_")
	splitToken := regex.Split(processedToken, -1)

	return splitToken, nil
}
