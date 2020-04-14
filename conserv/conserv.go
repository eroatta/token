// Package conserv provides the methods to split a token with the conservative approach,
// which is the benchmark algorithm.
package conserv

import (
	"strings"

	"github.com/eroatta/token/marker"
)

// Separator specifies the current separator.
var Separator string = " "

// Split on Conserv receives a token and returns an array of hard/soft words,
// split by:
// * Underscores
// * Numbers
// * CamelCase.
func Split(token string) string {
	processedToken := marker.OnDigits(token)
	processedToken = marker.OnLowerToUpperCase(processedToken)
	processedToken = marker.OnUpperToLowerCase(processedToken)
	processedToken = strings.ToLower(processedToken)

	return strings.Join(marker.SplitBy(processedToken), Separator)
}
