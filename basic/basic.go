// Package basic declares the functions and default expansions lists for expanding a token using the
// Basic algorithm.
package basic

import (
	"regexp"
	"strings"

	"github.com/eroatta/token/expansion"
	"github.com/eroatta/token/lists"
)

// DefaultExpansions contains the set of possible expansions included on the default configuration for Basic.
var DefaultExpansions = expansion.NewSetBuilder().AddList(lists.Dictionary).Build()

// Expand on Basic receives a token and returns an array of possible expansions.
//
// The Basic expansion algorithm builds a regular expression for the given token and
// runs it against several lists built from the source code and natural words from
// stop lists and dictionaries. It was proposed by Lawrie, Feild and Binkley.
func Expand(token string, srcWords expansion.Set, phrases map[string]string, defaultWords expansion.Set) []string {
	token = strings.ToLower(token)

	// build the search regex
	var pattern strings.Builder
	pattern.WriteString("\\b")
	for _, char := range token {
		pattern.WriteString("[")
		pattern.WriteRune(char)
		pattern.WriteString("]\\w*")
	}
	exp := regexp.MustCompile(pattern.String())

	// stage 1: should look on the words from the source code and then phrases lists
	expansions := exp.FindAllString(srcWords.String(), -1)
	if len(expansions) > 0 {
		return expansions
	}

	if phrase := phrases[token]; phrase != "" {
		return []string{strings.ReplaceAll(phrase, "-", " ")}
	}

	// stage 2: should look on the dictionary and stop lists
	expansions = exp.FindAllString(defaultWords.String(), -1)

	return expansions
}
