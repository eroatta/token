package splitters

import (
	"regexp"
)

// Conserv represents a conservative splitter.
type Conserv struct {
	Splitter
}

// Split on Conserv receives a token and returns an array of hard/soft words,
// split by:
// * Underscores
// * Numbers
// * CamelCase.
func (c Conserv) Split(token string) ([]string, error) {
	processedToken := addMarkersOnDigits(token)
	processedToken = addMarkersOnLowerToUpperCase(processedToken)
	processedToken = addMarkersOnUpperToLowerCase(processedToken)

	regex := regexp.MustCompile("_")
	splitToken := regex.Split(processedToken, -1)

	return splitToken, nil
}
