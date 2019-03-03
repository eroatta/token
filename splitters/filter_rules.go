package splitters

import "strings"

type filterFunc func(abbr string, word string) bool

// isTruncation checks if the abbreviation is a truncation of the word.
func isTruncation(abbr string, word string) bool {
	if abbr == "" {
		return false
	}

	return strings.HasPrefix(word, abbr)
}

// hasRemovedChar checks if the abbreviation matches the word when one of its characters are removed.
func hasRemovedChar(abbr string, word string) bool {
	if len(word)-len(abbr) != 1 {
		return false
	}

	return true
}

// hasRemovedVowels checks if the abbreviation matches the word when all of its vowels are removed.
func hasRemovedVowels(abrr string, word string) bool {
	r := strings.NewReplacer("a", "", "e", "", "i", "", "o", "", "u", "")
	removedVowels := r.Replace(strings.ToLower(word))

	return abrr == removedVowels
}
