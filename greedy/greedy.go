// Package greedy declares the functions and list builders for splitting a token using the
// Greedy algorithm.
package greedy

import (
	"strings"

	"github.com/eroatta/token-splitex/lists"
	"github.com/eroatta/token-splitex/marker"
)

// DefaultList contains the words included on the default configuration for Greedy,
// defined on Field, Binkley and Lawrie's paper.
var DefaultList = NewListBuilder().Dicctionary(lists.Dicctionary).
	KnownAbbreviations(lists.KnownAbbreviations).
	StopList(lists.Stop).
	Build()

// List declares the contract for a list.
type List interface {
	// Contains checks if a word is contained on the list.
	Contains(string) bool
}

type list struct {
	elements map[string]bool
}

func (l list) Contains(element string) bool {
	return l.elements[strings.ToLower(element)]
}

// ListBuilder is used to identify and set the lists that are used while processing a token with the Greedy algorithm.
type ListBuilder interface {
	// Dicctionary sets the dicctionary words.
	Dicctionary([]string) ListBuilder
	// KnownAbbreviations sets the list of common and known abbreviations.
	KnownAbbreviations([]string) ListBuilder
	// StopList sets the stop list words and expressions.
	StopList([]string) ListBuilder
	// Build creates a list with all the given lists.
	Build() List
}

type listBuilder struct {
	dictionary         []string
	knownAbbreviations []string
	stopList           []string
}

// NewListBuilder creates a ListBuilder for Greedy.
func NewListBuilder() ListBuilder {
	return &listBuilder{}
}

func (lb *listBuilder) Dicctionary(dicctionary []string) ListBuilder {
	lb.dictionary = dicctionary
	return lb
}

func (lb *listBuilder) KnownAbbreviations(knownAbbreviations []string) ListBuilder {
	lb.knownAbbreviations = knownAbbreviations
	return lb
}
func (lb *listBuilder) StopList(stopList []string) ListBuilder {
	lb.stopList = stopList
	return lb
}

func (lb *listBuilder) Build() List {
	elements := make(map[string]bool,
		len(lb.dictionary)+len(lb.knownAbbreviations)+len(lb.stopList))

	iter := func(list []string) {
		for _, e := range list {
			elements[strings.ToLower(e)] = true
		}
	}

	for _, l := range [][]string{lb.dictionary, lb.knownAbbreviations, lb.stopList} {
		iter(l)
	}

	return &list{
		elements: elements,
	}
}

// Split on Greedy receives a token and returns an array of hard and soft words,
// split by the Greedy algorithm proposed by Field, Binkley and Lawrie.
// This technique splits a token into hard words and checks for a greedy splitting those hard words
// that cannot be matched to any word on the list.
// The process evaluates prefixes and suffixes recursively until any of them are found on the list,
// prefering longer words.
func Split(token string, list List) []string {
	preprocessedToken := marker.OnDigits(token)
	preprocessedToken = marker.OnLowerToUpperCase(preprocessedToken)

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

	return splitToken
}

// findPrefix looks for the longest prefix exinsting on the list.
// If the token exists on the list, the process continues to look for the longest
// prefix within the remaining token. If not, then the process continues the search
// with a smaller token.
func findPrefix(token string, splitToken string, list List) string {
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
func findSuffix(token string, splitToken string, list List) string {
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
func chooseSplittings(preffixSplittings []string, suffixSplittings []string, list List) []string {
	if inListRatio(preffixSplittings, list) > inListRatio(suffixSplittings, list) {
		return preffixSplittings
	}

	return suffixSplittings
}

// inListRatio calculates the ratio between the total words
// passed as parameter vs. the total of those words found the list.
func inListRatio(words []string, list List) float64 {
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
