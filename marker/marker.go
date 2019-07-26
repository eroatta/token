// Package marker defines functions to add and splits tokens by markers.
package marker

import (
	"regexp"
	"strings"
)

var (
	leadingNumRegex   = regexp.MustCompile("([a-zA-Z])([0-9])")
	trailingNumRegex  = regexp.MustCompile("([0-9])([a-zA-Z])")
	lowerToUpperRegex = regexp.MustCompile("([a-z])([A-Z])")
	upperToLowerRegex = regexp.MustCompile("([A-Z]+)([A-Z])([a-z])")
)

// OnDigits applies markers between letters and numbers, and also between numbers
// and letters. The default marker is the underscore character.
func OnDigits(token string) string {
	//add markers between letters and numbers
	processedToken := leadingNumRegex.ReplaceAllString(token, "${1}_$2")

	//add markers between numbers and letters
	return trailingNumRegex.ReplaceAllString(processedToken, "${1}_$2")
}

// OnLowerToUpperCase applies markers on each lower-to-upper case combination.
func OnLowerToUpperCase(token string) string {
	return lowerToUpperRegex.ReplaceAllString(token, "${1}_$2")
}

// OnUpperToLowerCase applies markers on each upper-to-lower case combination, when
// there are at least two or more upper case letters and then a lower case letter.
func OnUpperToLowerCase(token string) string {
	return upperToLowerRegex.ReplaceAllString(token, "${1}_$2$3")
}

// SplitBy splits the given token by its markers.
func SplitBy(token string) []string {
	return strings.Split(strings.Trim(token, "_"), "_")
}
