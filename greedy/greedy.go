// Package greedy. Esta técnica utiliza un diccionario, una lista de abreviaturas conocidas
// y una lista de corte, la cual incluye identificadores predefinidos, librerías, funciones,
// nombres de variables comunes y letras individuales.
// Después de retornar cada hard word encontrada en alguna de las tres listas,
// como una soft word simple, el resto de las hard words se consideran para división.
// A partir de ahí, recursivamente, se analizan los sufijos y prefijos de las palabras
// hasta que se encuentren en alguna de las listas, prefiriendo siempre las palabras de mayor longitud.
package greedy

import (
	"strings"

	"github.com/eroatta/token-splitex/marker"
)

var defaultDictionary map[string]interface{}
var defaultKnownAbbreviations map[string]interface{}
var defaultStopList map[string]interface{}

// DefaultList contains the words included on the default configuration for Greedy,
// defined on Field, Binkley and Lawrie's paper.
var DefaultList List

func init() {
	var empty []string
	DefaultList = NewListBuilder().
		Dicctionary(empty).
		KnownAbbreviations(empty).
		StopList(empty).
		Build()
}

// List TODO
type List interface {
	Contains(string) bool
}

type list struct {
	elements map[string]bool
}

func (l list) Contains(element string) bool {
	return l.elements[strings.ToLower(element)]
}

// ListBuilder TODO
type ListBuilder interface {
	Dicctionary([]string) ListBuilder
	KnownAbbreviations([]string) ListBuilder
	StopList([]string) ListBuilder
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

// Split on Greedy receives a token and returns an array of hard/soft words,
// split by the Greedy algorithm proposed by TODO.
func Split(token string, list List) []string {
	preprocessedToken := marker.OnDigits(token)
	preprocessedToken = marker.OnLowerToUpperCase(preprocessedToken)

	splitToken := make([]string, 0, 10)
	for _, s := range marker.SplitBy(preprocessedToken) {
		if list.Contains(s) {
			splitToken = append(splitToken, s)
		} else {
			preffixSplittings := marker.SplitBy(findPreffix(s, "", list))
			suffixSplittings := marker.SplitBy(findSuffix(s, "", list))
			chosenSplittings := compare(preffixSplittings, suffixSplittings, list)

			splitToken = append(splitToken, chosenSplittings...)
		}
	}

	return splitToken
}

// findPreffix looks for the longest preffix that exists on any list.
// If the token exists on any list, the process continues to look for the longest
// preffix within the remaining token. If not, then the process continues the search
// with a smaller token.
func findPreffix(token string, splitToken string, list List) string {
	if len(token) == 0 {
		return splitToken
	}

	if list.Contains(token) {
		return token + "_" + findPreffix(splitToken, "", list)
	}

	sToken := string(token[len(token)-1]) + splitToken
	s := token[:len(token)-1]

	return findPreffix(s, sToken, list)
}

// findSuffix looks for the longest suffix that exists on any list.
// If the token exists on any list, the process continues to look for the longest
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

// Compare calculates the ratio between found words on any list and
// the total number of splittings. The string with the highest ratio
// is chosen as the proper splitting.
func compare(preffixSplittings []string, suffixSplittings []string, list List) []string {
	if inListRatio(preffixSplittings, list) > inListRatio(suffixSplittings, list) {
		return preffixSplittings
	}

	return suffixSplittings
}

// inListRatio calculates the ratio between the total words
// passed as parameter vs. the total of those words found on any list.
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
