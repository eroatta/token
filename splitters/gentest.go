package splitters

import (
	"fmt"
	"regexp"
	"strings"
)

// GenTest represents a Generation and Tests splitting algorithm, proposed by Lawrie, Binkley and Morrell.
type GenTest struct {
	// TODO review
	list string
}

// NewGenTest creates a new GenTest splitter.
func NewGenTest() *GenTest {
	return &GenTest{}
}

// PotentialSplit represents a GenTest potential split. It holds data related to the split, the softwords
// and their expansions.
type PotentialSplit struct {
	split      string
	softwords  []string
	expansions map[string][]string
}

// NewPotentialSplit creates and initializes a new potential split for the given hardword.
func NewPotentialSplit(hardword string) PotentialSplit {
	var softwords []string
	if hardword != "" {
		softwords = splitOnMarkers(hardword)
	}

	return PotentialSplit{
		split:      hardword,
		softwords:  softwords,
		expansions: make(map[string][]string, len(softwords)),
	}
}

// Split on GenTest receives a token and returns an array of hard/soft words,
// split by the Generation and Test algorithm proposed by Lawrie, Binkley and Morrell.
func (g *GenTest) Split(token string) ([]string, error) {
	preprocessedToken := addMarkersOnDigits(token)
	preprocessedToken = addMarkersOnLowerToUpperCase(preprocessedToken)

	splitToken := make([]string, 0, 10)
	for _, word := range splitOnMarkers(preprocessedToken) {
		splitToken = append(splitToken, word)
	}

	for _, tok := range splitToken {
		for _, pSplit := range generatePotentialSplits(tok) {
			fmt.Sprintln(pSplit.split)
		}
	}

	// generate potential splitings
	// e.g. splits("strlen") = {st-rlen, str-len}
	// find the list of expansions/translations for each potential splitting
	// e.g. E(st) = {stop, string, set}
	// e.g. E(rlen) = {riflemen}
	// e.g. E(str) = {steer, string}
	// e.g. E(len) = {lender, length}
	// define the similarity score between an expansion and the other soft words on the potential splittings
	// list, but for the same word (k not i)
	// uses the Google Data Set, and sim(Ei,j; Sk) = sum sim(Ei,j, e) <- e = E(Sk)
	// e.g. sim(string, len) = sim(string, lender) + sim(string, length)

	// computation for cohesion sums over all the other soft-words besides the current one
	// e.g. cohesion(string, str-len) = log(sim(string, len))

	// the maximal cohesion is chosen among the possible expansions
	// e.g. cohesion(string, str-len) = -5.9431
	// e.g. cohesion(steer, str-len) = -16.1222

	// score is the maximal cohesion with the selected translation as default expansion
	// e.g. score(str) = -5.9431 with "string" as expansion

	return splitToken, nil
}

// generatePotentialSplits generates every possible splitting for a given token.
func generatePotentialSplits(token string) []PotentialSplit {
	potentialSplits := []PotentialSplit{NewPotentialSplit(token)}
	for i := 1; i < len(token); i++ {
		leading := token[:i]
		trailing := token[i:]

		potentialSplit := leading + "_" + trailing
		potentialSplits = append(potentialSplits, NewPotentialSplit(potentialSplit))

		for j := 1; j < len(trailing); j++ {
			potentialSplit = leading + "_" + trailing[:j] + "_" + trailing[j:]
			potentialSplits = append(potentialSplits, NewPotentialSplit(potentialSplit))
		}
	}

	return potentialSplits
}

// findExpansions retrieves a set of words that could be considered expansions for the input string.
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

	// find on a list (what should this list include?)
	return exp.FindAllString(g.list, -1)
}

func testSplit(split string) {

}

func similarity() {

}
