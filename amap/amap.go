package amap

import (
	"regexp"
	"sort"
	"strings"

	porterstemmer "github.com/reiver/go-porterstemmer"
)

type searchExpansion func(pattern, []string, string, string, string, string) []string

var (
	consonants, _ = regexp.Compile("[a-z][^aeiou]+")
	manyVowels, _ = regexp.Compile("[a-z][aeiou][aeiou]+")
	searchers     = map[string]searchExpansion{
		singleWordGroup: searchSingleWordExpansion,
		multiWordGroup:  searchMultiWordExpansion,
	}
)

// Amap represents an Automatically Mining Abbreviations in Programs expander.
type Amap struct {
	variableDeclarations []string
	methodName           string
	methodBodyText       string
	methodComments       string
	packageComments      string
	text                 []string
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
	patterns := []pattern{
		(&patternBuilder{}).kind(acronymType).shortForm(token).build(),
		(&patternBuilder{}).kind(prefixType).shortForm(token).build(),
		(&patternBuilder{}).kind(droppedLettersType).shortForm(token).build(),
		(&patternBuilder{}).kind(wordCombinationType).shortForm(token).build(),
	}

	// TODO: refactor
	varDeclarations := a.variableDeclarations
	methodName := a.methodName
	methodBodyText := a.methodBodyText
	methodComments := a.methodComments
	packageComments := a.packageComments
	referenceText := a.text

	var expansion string
	for _, pttrn := range patterns {
		search := searchers[pttrn.group]
		longForms := search(pttrn, varDeclarations, methodName, methodBodyText,
			methodComments, packageComments)
		if len(longForms) == 1 {
			expansion = longForms[0]
			break
		}

		if len(longForms) > 1 {
			expansion = findMostFrequentLongForm(pttrn, longForms, referenceText)
			break
		}
	}

	var expansions []string
	if expansion != "" {
		expansions = append(expansions, expansion)
	}

	return expansions
}

// searchSingleWordExpansion looks for candidate long forms for a given pattern, focusing on single word expansions.
func searchSingleWordExpansion(pttrn pattern, variableDeclarations []string, methodName string,
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

// searchMultiWordExpansion looks for candidate long forms for a given pattern, focusing on single word expansions.
func searchMultiWordExpansion(pttrn pattern, variableDeclarations []string, methodName string,
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

// findMostFrequentLongForm selects a long form between the available long forms.
// The process follows several steps. On the first step, it uses the long form that most frequently matches the
// short form’s pattern in this scope.
// On the second step, words with the same stem are grouped and the frequencies updated accordingly.
// Finally, if the previous steps fail, the MFE process is used.
func findMostFrequentLongForm(pttrn pattern, longForms []string, referenceText []string) string {
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
	mfw = mostFrequentExpansion(pttrn, referenceText)

	return mfw
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
		var sw string
		for _, c := range strings.Split(w, " ") {
			sw += " " + porterstemmer.StemString(c)
		}
		stemmedWords[i] = strings.TrimSpace(sw)
	}

	return stemmedWords
}

// mostFrequentExpansion calculates the most frequent expansion based on the number of matches
// between the short form and the given text.
func mostFrequentExpansion(pttrn pattern, text []string) string {
	var totalMatches int
	results := make(map[string]int, 0)

	matcher, _ := regexp.Compile(pttrn.regex)
	for _, t := range text {
		matches := matcher.FindAllString(t, -1)
		totalMatches += len(matches)

		for _, match := range matches {
			results[match]++
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
