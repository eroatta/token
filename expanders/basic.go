package expanders

import (
	"regexp"
	"strings"
)

// Basic represents the Basic expansion algorithm, proposed by Lawrie, Feild and Binkley.
type Basic struct {
	words       *string
	stopAndDicc *string
	srcPhrases  map[string]string
}

// NewBasic creates a new Basic expander with the given lists.
func NewBasic(srcWords map[string]interface{}, srcPhrases map[string]string, stopList map[string]interface{}, dicc map[string]interface{}) *Basic {
	arrWords := make([]string, len(srcWords))
	for k := range srcWords {
		arrWords = append(arrWords, k)
	}
	words := strings.Join(arrWords, " ")

	arrStopAndDicc := make([]string, len(dicc)+len(stopList))
	for k := range dicc {
		arrStopAndDicc = append(arrStopAndDicc, k)
	}

	// merge lists to avoid duplication
	for k := range stopList {
		if _, found := dicc[k]; !found {
			arrStopAndDicc = append(arrStopAndDicc, k)
		}
	}
	stopAndDicc := strings.Join(arrStopAndDicc, " ")

	return &Basic{
		words:       &words,
		srcPhrases:  srcPhrases,
		stopAndDicc: &stopAndDicc,
	}
}

// Expand on Basic receives a token and returns an array of possible expansions.
func (b Basic) Expand(token string) []string {
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
	expansions := exp.FindAllString(*b.words, -1)
	if len(expansions) > 0 {
		return expansions
	}

	if phrase := b.srcPhrases[token]; phrase != "" {
		return strings.Split(phrase, "-")
	}

	// stage 2: should look on the dicctionary and stop lists
	expansions = exp.FindAllString(*b.stopAndDicc, -1)

	return expansions
}
