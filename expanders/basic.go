package expanders

import (
	"errors"
	"regexp"
	"strings"
)

// Basic represents the Basic expansion algorithm, proposed by Lawrie, Feild and Binkley.
type Basic struct {
	srcWords    map[string]interface{}
	srcPhrases  map[string]string
	stopList    map[string]interface{}
	dicctionary map[string]interface{}
}

// NewBasic creates a new Basic expander with the given lists.
func NewBasic(srcWords map[string]interface{}, srcPhrases map[string]string, stopList map[string]interface{}, dicc map[string]interface{}) *Basic {
	return &Basic{
		srcWords:    srcWords,
		srcPhrases:  srcPhrases,
		stopList:    stopList,
		dicctionary: dicc,
	}
}

// Expand on Basic receives a token and returns an array of possible expansions.
func (b Basic) Expand(token string) ([]string, error) {
	token = strings.ToLower(token)
	if ok := b.stopList[token]; ok != nil {
		return []string{token}, nil
	}

	if phrase := b.srcPhrases[token]; phrase != "" {
		return strings.Split(phrase, "-"), nil
	}

	if ok := b.srcWords[token]; ok != nil {
		return []string{token}, nil
	}

	// build the search regex
	var pattern strings.Builder
	pattern.WriteString("\\b")
	for _, char := range token {
		pattern.WriteString("[")
		pattern.WriteRune(char)
		pattern.WriteString("]\\w")
	}
	exp := regexp.MustCompile(pattern.String())

	// TODO complete, use defined errors
	expansions := exp.FindAllString("", -1)
	if len(expansions) > 1 {
		return nil, errors.New("Multiple matches")
	}

	if len(expansions) == 0 {
		return nil, errors.New("No match")
	}

	return expansions, nil
}
