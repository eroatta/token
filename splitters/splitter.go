package splitters

import (
	"regexp"
	"strings"
)

// Splitter defines the required behavior for any token splitter.
type Splitter interface {
	Split(token string) ([]string, error)
}

func addMarkersOnDigits(token string) string {
	//add markers between letters and numbers
	leadingNumRegex := regexp.MustCompile("([a-zA-Z])([0-9])")
	processedToken := leadingNumRegex.ReplaceAllString(token, "${1}_$2")

	//add markers between numbers and letters
	trailingNumRegex := regexp.MustCompile("([0-9])([a-zA-Z])")
	processedToken = trailingNumRegex.ReplaceAllString(processedToken, "${1}_$2")

	return processedToken
}

func addMarkersOnLowerToUpperCase(token string) string {
	//add markers for lower to upper case camel-case combination
	lowerToUpperRegex := regexp.MustCompile("([a-z])([A-Z])")
	processedToken := lowerToUpperRegex.ReplaceAllString(token, "${1}_$2")

	return processedToken
}

func addMarkersOnUpperToLowerCase(token string) string {
	//add markers for upper to lower case camel-case combination
	upperToLowerRegex := regexp.MustCompile("([A-Z]+)([A-Z])([a-z])")
	processedToken := upperToLowerRegex.ReplaceAllString(token, "${1}_$2$3")

	return processedToken
}

func splitOnMarkers(token string) []string {
	regex := regexp.MustCompile("_")
	splitToken := regex.Split(strings.Trim(token, "_"), -1)

	return splitToken
}
