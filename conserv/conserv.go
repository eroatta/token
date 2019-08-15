// Package conserv provides the methods to split a token with the conservative approach,
// which is the benchmark algorithm.
package conserv

import (
	"github.com/eroatta/token/marker"
)

// Split on Conserv receives a token and returns an array of hard/soft words,
// split by:
// * Underscores
// * Numbers
// * CamelCase.
func Split(token string) []string {
	processedToken := marker.OnDigits(token)
	processedToken = marker.OnLowerToUpperCase(processedToken)
	processedToken = marker.OnUpperToLowerCase(processedToken)

	return marker.SplitBy(processedToken)
}
