package expanders

import (
	"regexp"
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

// singleWordExpansion looks for candidate long forms for a given pattern
func (a Amap) singleWordExpansion(pttrn pattern) []string {
	// input: short form
	// input: regular expression, pattern

	var methodName string
	// input: method body text
	var methodBodyText string
	// input: method comments
	var methodComments string
	// input: class comments (prefix only)
	var packageComments string

	var longForms []string

	// restricts the search to prefix or dropped letters to those short forms longer than 3 letters or
	// composed of all consonants letters with an optional leading vowel
	if (pttrn.kind == prefixType || consonants.MatchString(pttrn.shortForm) || len(pttrn.shortForm) > 3) &&
		!manyVowels.MatchString(pttrn.shortForm) {

		// 9: Search TypeNames and corresponding declared variable names for “pattern sf”
		// TODO: complete

		// 10: Search MethodName for “pattern”
		matcher, _ := regexp.Compile(pttrn.regex)
		if matcher.MatchString(methodName) {
			longForms = append(longForms, methodName)
			if len(longForms) == 1 {
				return longForms
			}
		}
		// 11: Search Statements for “pattern sf” and “sf pattern”
		// TODO: complete

		if len(pttrn.shortForm) != 2 {
			// 13: Search method words for “pattern”
			matcher, _ := regexp.Compile(pttrn.regex)
			longForms = matcher.FindAllString(methodBodyText, -1)
			if len(longForms) == 1 {
				return longForms
			}

			// 14: Search method comment words for “pattern”
			longForms = matcher.FindAllString(methodComments, -1)
			if len(longForms) == 1 {
				return longForms
			}
		}
		if pttrn.kind == prefixType && len(pttrn.shortForm) > 1 {
			// 17: Search class comment words for “pattern”
			matcher, _ := regexp.Compile(pttrn.regex)
			longForms = matcher.FindAllString(packageComments, -1)
			if len(longForms) == 1 {
				return longForms
			}
		}
	}

	// output: long form candidates
	return longForms
}
