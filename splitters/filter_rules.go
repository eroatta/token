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
	if abbr == "" {
		return false
	}

	if len(word)-len(abbr) != 1 {
		return false
	}

	if abbr[0] != word[0] {
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

	r := strings.NewReplacer("a", "", "e", "", "i", "", "o", "", "u", "")
	removedVowels := r.Replace(strings.ToLower(word))

	return abrr == removedVowels
}
