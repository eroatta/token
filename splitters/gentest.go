package splitters

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

// GenTest represents a Generation and Tests splitting algorithm, proposed by Lawrie, Binkley and Morrell.
type GenTest struct {
	// TODO review
	list        string
	simProvider SimilarityProvider
	context     []string
	dicctionary *map[string]interface{}
}

// NewGenTest creates a new GenTest splitter.
func NewGenTest() *GenTest {
	return &GenTest{}
}

// SimilarityProvider defines the contract that a provider for calculating similarities should implement.
type SimilarityProvider interface {
	Sim(firstWord string, secondWord string) float64
}

// Split on GenTest receives a token and returns an array of hard/soft words,
// split by the Generation and Test algorithm proposed by Lawrie, Binkley and Morrell.
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
//
// computation for cohesion sums over all the other soft-words besides the current one
// e.g. cohesion(string, str-len) = log(sim(string, len))
//
// the maximal cohesion is chosen among the possible expansions
// e.g. cohesion(string, str-len) = -5.9431
// e.g. cohesion(steer, str-len) = -16.1222
//
// score is the maximal cohesion with the selected translation as default expansion
// e.g. score(str) = -5.9431 with "string" as expansion
func (g *GenTest) Split(token string) ([]string, error) {
	preprocessedToken := addMarkersOnDigits(token)
	preprocessedToken = addMarkersOnLowerToUpperCase(preprocessedToken)

	splitToken := make([]string, 0, 10)
	for _, tok := range splitOnMarkers(preprocessedToken) {
		// discard dictionary words
		// TODO: review this pointer call
		if (*g.dicctionary)[tok] != nil {
			// TODO remove
			fmt.Sprintln("Discarded diccionary word: " + tok)
			splitToken = append(splitToken, tok)
		}

		potentialSplits := generatePotentialSplits(tok)
		for _, pSplit := range potentialSplits {
			fmt.Sprintln(pSplit.split)

			// for each potential split find the list of expansions
			for _, softword := range pSplit.softwords {
				for _, translation := range g.findExpansions(softword.word) {
					softword.expansions = append(softword.expansions, expansion{translation, 0})
				}
			}

			// for each expansion of every softword of the potential split, calculate the cohesion with the rest
			// of the softwords on the potential split
			for index, softword := range pSplit.softwords {
				for _, expansion := range softword.expansions {
					cohesion := g.cohesion(pSplit, expansion.translation, index)
					expansion.cohesion = cohesion
				}
			}
		}

		tokenBestSplit := findBestSplit(potentialSplits)
		// TODO: remove
		fmt.Sprintln("Best Split: " + tokenBestSplit.split)
		splitToken = append(splitToken, splitOnMarkers(tokenBestSplit.split)...)
	}

	return splitToken, nil
}

// generatePotentialSplits generates every possible splitting for a given token.
func generatePotentialSplits(token string) []potentialSplit {
	potentialSplits := []potentialSplit{newPotentialSplit(token)}
	for i := 1; i < len(token); i++ {
		leading := token[:i]
		trailing := token[i:]

		potentialSplit := leading + "_" + trailing
		potentialSplits = append(potentialSplits, newPotentialSplit(potentialSplit))

		for j := 1; j < len(trailing); j++ {
			potentialSplit = leading + "_" + trailing[:j] + "_" + trailing[j:]
			potentialSplits = append(potentialSplits, newPotentialSplit(potentialSplit))
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
	// TODO add rules!
	return exp.FindAllString(g.list, -1)
}

// similarityScore returns the Log of the similarity computed by the SimilarityProvider.
// For two equal words, the similarity score is zero.
func (g *GenTest) similarityScore(word string, anotherWord string) float64 {
	if word == anotherWord {
		return 0
	}

	logSim := math.Log(g.simProvider.Sim(word, anotherWord))
	if math.IsNaN(logSim) || math.IsInf(logSim, 0) {
		return 0
	}

	return logSim
}

// cohesion computes the similarity score between an expansion and the other soft words on the potential splittings
// list, but for the same word (k not i).
func (g *GenTest) cohesion(potentialSplit potentialSplit, expansion string, index int) float64 {
	var cohesion float64
	for i, softword := range potentialSplit.softwords {
		if index == i {
			continue
		}

		// compute cohesion between the proposed expansion and the expansions for the others soft words
		for _, translation := range softword.expansions {
			cohesion += g.similarityScore(expansion, translation.translation)
		}
	}

	// add context cohesion
	for _, contextWord := range g.context {
		cohesion += g.similarityScore(expansion, contextWord)
	}

	return cohesion
}
