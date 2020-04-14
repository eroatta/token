// Package greedy declares the functions and list builders for splitting a token using the
// Greedy algorithm.
package greedy

import (
	"strings"

	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/marker"
)

// Separator specifies the current separator.
var Separator string = " "

// DefaultList contains the words included on the default configuration for Greedy,
// defined on Field, Binkley and Lawrie's paper.
// This list includes words from:
// * a dictionary
// * a known abbreviations list
// * a stop list
var DefaultList = lists.NewBuilder().Add(lists.Dictionary.Elements()...).
	Add(lists.KnownAbbreviations.Elements()...).
	Add(lists.Stop.Elements()...).
	Build()

// Split on Greedy receives a token and returns an array of hard and soft words,
// split by the Greedy algorithm proposed by Field, Binkley and Lawrie.
// This technique splits a token into hard words and checks for a greedy splitting those hard words
// that cannot be matched to any word on the list.
// The process evaluates prefixes and suffixes recursively until any of them are found on the list,
// preferring longer words.
func Split(token string, list lists.List) string {
	preprocessedToken := marker.OnDigits(token)
	preprocessedToken = marker.OnLowerToUpperCase(preprocessedToken)
	preprocessedToken = strings.ToLower(preprocessedToken)

	splitToken := make([]string, 0, 10)
	for _, s := range marker.SplitBy(preprocessedToken) {
		if list.Contains(s) {
			splitToken = append(splitToken, s)
		} else {
			preffixSplittings := marker.SplitBy(findPrefix(s, "", list))
			suffixSplittings := marker.SplitBy(findSuffix(s, "", list))
			chosenSplittings := chooseSplittings(preffixSplittings, suffixSplittings, list)

			splitToken = append(splitToken, chosenSplittings...)
		}
	}

	return strings.Join(splitToken, Separator)
}

// findPrefix looks for the longest prefix exinsting on the list.
// If the token exists on the list, the process continues to look for the longest
// prefix within the remaining token. If not, then the process continues the search
// with a smaller token.
func findPrefix(token string, splitToken string, list lists.List) string {
	if len(token) == 0 {
		return splitToken
	}

	if list.Contains(token) {
		return token + "_" + findPrefix(splitToken, "", list)
	}

	sToken := string(token[len(token)-1]) + splitToken
	s := token[:len(token)-1]

	return findPrefix(s, sToken, list)
}

// findSuffix looks for the longest suffix existing on the list.
// If the token exists on the list, the process continues to look for the longest
// suffix within the remaining token. If not, the the process continues the search
// with a smaller token.
func findSuffix(token string, splitToken string, list lists.List) string {
	if len(token) == 0 {
		return splitToken
	}

	if list.Contains(token) {
		return findSuffix(splitToken, "", list) + "_" + token
	}

	sToken := splitToken + string(token[0])
	s := token[1:len(token)]

	return findSuffix(s, sToken, list)
}

// chooseSplittings calculates the ratio between found words on the list and
// the total number of splittings and chooses the proper splitting.
func chooseSplittings(preffixSplittings []string, suffixSplittings []string, list lists.List) []string {
	if inListRatio(preffixSplittings, list) > inListRatio(suffixSplittings, list) {
		return preffixSplittings
	}

	return suffixSplittings
}

// inListRatio calculates the ratio between the total words
// passed as parameter vs. the total of those words found the list.
func inListRatio(words []string, list lists.List) float64 {
	found := 0
	for _, word := range words {
		if list.Contains(word) {
			found++
		}
	}

	ratio := 0.0
	if len(words) > 0 {
		ratio = float64(found) / float64(len(words))
	}

	return ratio
}
