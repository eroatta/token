package expanders

import "strings"

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

	//var expansions []string

	// build the search regex

	//it := b.dicc.Iterator()
	/*for elem := range it.C {

	}*/

	return []string{}, nil
}

/*
func (g *GenTest) findExpansions(input string) []string {
	if input == "" || strings.TrimSpace(input) == "" {
		return []string{}
	}

	// build the regexp
	var pattern strings.Builder
	pattern.WriteString("\\b")
	for _, char := range input {
		pattern.WriteString("[")
		pattern.WriteRune(char)
		pattern.WriteString("]\\w*")
	}
	exp := regexp.MustCompile(pattern.String())

	expansions := make([]string, 0)
	for _, candidate := range exp.FindAllString(g.possibleExpansions, -1) {
		if any(input, candidate, isTruncation, hasRemovedChar, hasRemovedVowels, hasRemovedCharAfterRemovedVowels) {
			expansions = append(expansions, candidate)
		}
	}

	return expansions
}*/
