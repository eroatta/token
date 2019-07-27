package gentest

import "strings"

type filterFunc func(abbr string, word string) bool

// isTruncation checks if the abbreviation is a truncation of the word.
func isTruncation(abbr string, word string) bool {
	if abbr == "" || word == "" {
		return false
	}

	if len(abbr) > 1 && len(word) > 1 && abbr[len(abbr)-1:] == word[len(word)-1:] && abbr[len(abbr)-1:] == "s" {
		abbr = abbr[0 : len(abbr)-2]
		word = word[0 : len(word)-2]
	}

	return strings.HasPrefix(word, abbr)
}

// hasRemovedChar checks if the abbreviation matches the word when one of its characters are removed.
func hasRemovedChar(abbr string, word string) bool {
	if abbr == "" || len(word)-len(abbr) != 1 || abbr[0] != word[0] {
		return false
	}

	removedChars := 0
	var i, j int
	for i < len(abbr) && j < len(word) {
		if abbr[i] == word[j] {
			i++
			j++
		} else {
			removedChars++
			j++
		}
	}

	return removedChars == 1
}

// hasRemovedVowels checks if the abbreviation matches the word when all of its vowels are removed.
func hasRemovedVowels(abrr string, word string) bool {
	if abrr == "" {
		return false
	}

	return abrr == removeVowels(strings.ToLower(word))
}

func removeVowels(word string) string {
	r := strings.NewReplacer("a", "", "e", "", "i", "", "o", "", "u", "")
	return r.Replace(strings.ToLower(word))
}

// hasRemovedCharAfterRemovedVowels checks if an abbreviation matches a word with previously removed
// vowels and a char.
func hasRemovedCharAfterRemovedVowels(abbr string, word string) bool {
	if abbr == "" {
		return false
	}

	return hasRemovedChar(abbr, removeVowels(strings.ToLower(word)))
}
