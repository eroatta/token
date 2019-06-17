package expanders

import (
	"regexp"
	"sort"
	"strings"

	porterstemmer "github.com/reiver/go-porterstemmer"
)

var consonants *regexp.Regexp
var manyVowels *regexp.Regexp

func init() {
	consonants, _ = regexp.Compile("[a-z][^aeiou]+")
	manyVowels, _ = regexp.Compile("[a-z][aeiou][aeiou]+")
}

// Amap represents an Automatically Mining Abbreviations in Programs expander.
type Amap struct {
}

// NewAmap creates an AMAP expander.
func NewAmap() *Amap {
	return &Amap{}
}

// Expand on AMAP receives a token and returns and array of possible expansions.
//
// The AMAP expansion algorithm handles single-word and multi-word abbreviations.
// For each type of abbreviation AMAP creates and applies a pattern to look for possible
// expansions. AMAP is capable of select the more appropiate expansions based on available
// information on the given context.
func (a Amap) Expand(token string) []string {
	return []string{}
}

// singleWordExpansion looks for candidate long forms for a given pattern, focusing on single word expansions.
func (a Amap) singleWordExpansion(pttrn pattern, variableDeclarations []string, methodName string,
	methodBodyText string, methodComments string, packageComments string) []string {
	var longForms []string

	// restricts the search to prefix or dropped letters to those short forms longer than 3 letters or
	// composed of all consonants letters with an optional leading vowel
	if (pttrn.kind == prefixType || consonants.MatchString(pttrn.shortForm) || len(pttrn.shortForm) > 3) &&
		!manyVowels.MatchString(pttrn.shortForm) {

		// 9: Search TypeNames and corresponding declared variable names for “pattern sf”
		matcher, _ := regexp.Compile(pttrn.regex + "[ ]" + pttrn.shortForm)
		for _, v := range variableDeclarations {
			if matcher.MatchString(v) {
				// append only the matching name to the candidate expansions
				longForms = append(longForms, strings.Split(v, " ")[0])
			}
		}
		if len(longForms) == 1 {
			return longForms
		}

		// 10: Search MethodName for “pattern”
		matcher, _ = regexp.Compile(pttrn.regex)
		if matcher.MatchString(methodName) {
			longForms = append(longForms, methodName)
			if len(longForms) == 1 {
				return longForms
			}
		}

		// 11: Search Statements for “pattern sf” and “sf pattern” (ignored)

		if len(pttrn.shortForm) != 2 {
			// 13: Search method words for “pattern”
			matcher, _ := regexp.Compile(pttrn.regex)
			longForms = append(longForms, matcher.FindAllString(methodBodyText, -1)...)
			if len(longForms) == 1 {
				return longForms
			}

			// 14: Search method comment words for “pattern”
			longForms = append(longForms, matcher.FindAllString(methodComments, -1)...)
			if len(longForms) == 1 {
				return longForms
			}
		}
		if pttrn.kind == prefixType && len(pttrn.shortForm) > 1 {
			// 17: Search class comment words for “pattern”
			matcher, _ := regexp.Compile(pttrn.regex)
			longForms = append(longForms, matcher.FindAllString(packageComments, -1)...)
			if len(longForms) == 1 {
				return longForms
			}
		}
	}

	// output: long form candidates
	return longForms
}

// multiWordExpansion looks for candidate long forms for a given pattern, focusing on single word expansions.
func (a Amap) multiWordExpansion(pttrn pattern, variableDeclarations []string, methodName string,
	methodBodyText string, methodComments string, packageComments string) []string {
	var longForms []string

	if pttrn.kind == acronymType || len(pttrn.shortForm) > 3 {
		// 9: Search TypeNames and corresponding declared variable names for “pattern sf”
		matcher, _ := regexp.Compile(pttrn.regex + "[ ]" + pttrn.shortForm)
		for _, v := range variableDeclarations {
			if matcher.MatchString(v) {
				// append only the matching name to the candidate expansions
				longForms = append(longForms, strings.TrimSpace(strings.TrimSuffix(v, pttrn.shortForm)))
			}
		}
		if len(longForms) == 1 {
			return longForms
		}

		// 10: Search MethodName for “pattern”
		matcher, _ = regexp.Compile(pttrn.regex)
		if matcher.MatchString(methodName) {
			longForms = append(longForms, methodName)
			if len(longForms) == 1 {
				return longForms
			}
		}

		// 11: Search all identifiers in the method for “pattern” (ignored)

		// 12: Search string literals for “pattern”
		longForms = append(longForms, matcher.FindAllString(methodBodyText, -1)...)
		if len(longForms) == 1 {
			return longForms
		}

		// 13: Search method comment words for “pattern”
		longForms = append(longForms, matcher.FindAllString(methodComments, -1)...)
		if len(longForms) == 1 {
			return longForms
		}

		// 15: If acronym, search class comment words for “pattern”
		if pttrn.kind == acronymType {
			longForms = append(longForms, matcher.FindAllString(packageComments, -1)...)
		}
	}

	return longForms
}

// filterMultipleLongForms handles the TODO
// Step 1. Use the long form that most frequently matches the
// short form’s pattern in this scope. For example, if ‘value’
// matched the prefix pattern for ‘val’ three times and ‘valid’
// only once, return ‘value’.
// Step 2. Group words with the same stem [15] and update
// the frequencies accordingly. For example, if the words ‘default’ (2 matches), ‘defaults’ (2 matches), and ‘define’ (2
// matches) all match the prefix pattern for ‘def’, group ‘default’ and ‘defaults’ to be the shortest long form, ‘default’
// (4 matches), and return the long form with the highest frequency.
// Step 3. If there is still no clear winner, continue searching
// for the pattern at broader scope levels. For example, if both
// ‘string buffer’ and ‘sound byte’ match the acronym pattern
// for ‘sb’ at the method identifier level, continue to search for
// the acronym pattern in string literals and comments. We
// store the frequencies of the tied long forms so that the most
// frequently occurring long form candidates are favored when
// searching the broader scope.
// Step 4. If all else fails, abandon the search and let MFE
// select the long form. At this point we stop searching for
// long form candidates of different abbreviation types. For
// example, if a prefix pattern has already found long form
// candidates, we avoid finding dropped letter long form candidates by halting the search for a given short form within
// a method.
func (a Amap) filterMultipleLongForms(longForms []string) string {
	// step 1: use the long form that most frequently matches the short form's pattern in this scope
	mfw := mostFrequentWord(longForms)
	if mfw != "" {
		return mfw
	}

	// step 2: group words with the same stem
	mfw = mostFrequentWord(stemmedWords(longForms))
	if mfw != "" {
		return mfw
	}

	// step 3: skipped, because we continue searching at broader levels if multiple matches are found

	// step 4: use MFE

	return ""
}

func mostFrequentWord(words []string) string {
	if len(words) == 0 {
		return ""
	}

	counts := make(map[string]int, len(words))
	for _, w := range words {
		counts[w]++
	}

	var countArr []wordCount
	for k, v := range counts {
		countArr = append(countArr, wordCount{k, v})
	}

	if len(countArr) == 1 {
		return countArr[0].key
	}

	sort.Slice(countArr, func(i, j int) bool {
		return countArr[i].value > countArr[j].value
	})

	var mostFrequentWord string
	if countArr[0].value > countArr[1].value {
		mostFrequentWord = countArr[0].key
	}

	return mostFrequentWord
}

type wordCount struct {
	key   string
	value int
}

func stemmedWords(words []string) []string {
	stemmedWords := make([]string, len(words))
	for i, w := range words {
		stemmedWords[i] = porterstemmer.StemString(w)
	}

	return stemmedWords
}

// most frequent expansion: match short form (using the same pattern) to the whole
// program/packages
func mostFrequentExpansion(pttrn pattern, text []string) string {
	var totalMatches int
	results := make(map[string]int, 0)

	matcher, _ := regexp.Compile(pttrn.regex)
	for _, t := range text {
		matches := matcher.FindAllString(t, -1)
		totalMatches += len(matches)

		for _, match := range matches {
			results[porterstemmer.StemString(match)]++
		}
	}

	type relativeFreq struct {
		word  string
		value float64
	}

	relativeFrequencies := make([]relativeFreq, len(results))
	for word, count := range results {
		rf := float64(count) / float64(totalMatches)
		if count >= 3 && rf > 0.5 {
			relativeFrequencies = append(relativeFrequencies, relativeFreq{word: word, value: rf})
		}
	}

	sort.Slice(relativeFrequencies, func(i, j int) bool {
		return relativeFrequencies[i].value > relativeFrequencies[j].value
	})

	var mostFrequentExpansion string
	if len(relativeFrequencies) > 0 {
		mostFrequentExpansion = relativeFrequencies[0].word
	}

	return mostFrequentExpansion
}
