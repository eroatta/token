// Package gentest provides the functions to split and expand a token using the GenTest algorithm.
package gentest

import (
	"math"
	"regexp"
	"strings"

	"github.com/eroatta/token/expansion"
	"github.com/eroatta/token/lists"
	"github.com/eroatta/token/marker"
)

const (
	closeToZeroProbability = 0.000000000000001
)

// SimilarityCalculator is the interface that wraps the basic Sim method.
type SimilarityCalculator interface {
	// Similarity defines the probability for two words to be co-located on the same sentence.
	Similarity(firstWord string, secondWord string) (prob float64)
}

// Split on GenTest receives a token and returns an array of hard/soft words,
// split by the Generation and Test algorithm proposed by Lawrie, Binkley and Morrell.
//
// The algorithm splits the token by its markers. Then, for each new token, checks if
// it is a dictionary word or not. If it is a dictionary word, then no more processing is done.
// But if the token is not a dictionary word, then the generate and testing algorithm starts.
//
// The first step is to retrieve the list of potential splits, each potential split composed by a
// list of softwords.
// For each softword on the potential split, the algorithm looks for a set of expansions, and then
// calculates the cohesion.
// Cohesion is calculated for every expansion on every softword of each potential split, comparing the expansion
// with the rest of the expansions on every softword belonging to the same potential split.
//
// The potential split with the highest score is the selected split, and all the combined selected splits
// form the splitted token.
func Split(token string, simCalc SimilarityCalculator, context lists.List, peSet expansion.Set) []string {
	parts := generateAndTest(token, simCalc, context, peSet)

	splits := make([]string, 0, len(parts))
	for _, part := range parts {
		splits = append(splits, marker.SplitBy(part.split)...)
	}

	return splits
}

// Expand on GenTest receives a token and returns an array of hard and expanded softwords,
// based on the Generation and Test algorithm proposed by Lawrie, Binkley and Morrell.
//
// The algorithm splits the token by its markers. Then, for each new token, checks if
// it is a dictionary word or not. If it is a dictionary word, then no more processing is done.
// But if the token is not a dictionary word, then the generate and testing algorithm starts.
//
// The first step is to retrieve the list of potential splits, each potential split composed by a
// list of softwords.
// For each softword on the potential split, the algorithm looks for a set of expansions, and then
// calculates the cohesion.
// Cohesion is calculated for every expansion on every softword of each potential split, comparing the expansion
// with the rest of the expansions on every softword belonging to the same potential split.
//
// The potential split with the highest score is the selected split, and all the combined best expansions
// form the expanded token.
func Expand(token string, simCalc SimilarityCalculator, context lists.List, peSet expansion.Set) []string {
	parts := generateAndTest(token, simCalc, context, peSet)

	expansions := make([]string, 0, len(parts))
	for _, part := range parts {
		expansions = append(expansions, marker.SplitBy(part.bestExpansion())...)
	}

	return expansions
}

func generateAndTest(token string, simCalc SimilarityCalculator, context lists.List, peSet expansion.Set) []potentialSplit {
	similarity := func(w1 string, w2 string) float64 {
		return similarityScore(simCalc, w1, w2)
	}
	contextWords := context.Elements()

	preprocessedToken := marker.OnDigits(token)
	preprocessedToken = marker.OnLowerToUpperCase(preprocessedToken)

	selectedSplits := make([]potentialSplit, 0, 10)
	for _, tok := range marker.SplitBy(preprocessedToken) {
		// discard one letter tokens and dictionary words
		if len(tok) == 1 || peSet.Contains(tok) {
			selectedSplits = append(selectedSplits, hardwordAsPotentialSplit(tok))
			continue
		}

		potentialSplits := make([]potentialSplit, 0)
		for _, pSplit := range generatePotentialSplits(tok) {
			// for each potential split softword find the list of expansions
			for i := 0; i < len(pSplit.softwords); i++ {
				// if no expansion is found, the input itself must be considered an expansion
				expansions := findExpansions(pSplit.softwords[i].word, peSet)
				if len(expansions) == 0 {
					expansions = append(expansions, pSplit.softwords[i].word)
				}

				for _, translation := range expansions {
					pSplit.softwords[i].expansions = append(pSplit.softwords[i].expansions, possibleExpansion{translation, 0})
				}
			}

			// for each expansion of every softword of the potential split, calculate the cohesion with the rest
			// of the softwords on the potential split
			for i := 0; i < len(pSplit.softwords); i++ {
				for j := 0; j < len(pSplit.softwords[i].expansions); j++ {
					cohesion := cohesion(similarity, pSplit, pSplit.softwords[i].expansions[j].translation, i, contextWords)
					pSplit.softwords[i].expansions[j].cohesion = cohesion
				}
			}

			// calculate the score considering the context too
			pSplit.score = score(similarity, pSplit, contextWords)

			potentialSplits = append(potentialSplits, pSplit)
		}

		tokenBestSplit := findBestSplit(potentialSplits)
		selectedSplits = append(selectedSplits, tokenBestSplit)
	}

	return selectedSplits
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

func hardwordAsPotentialSplit(hardword string) potentialSplit {
	return potentialSplit{
		split: hardword,
		softwords: []softword{{
			word: hardword,
			expansions: []possibleExpansion{
				{translation: hardword},
			},
		}},
	}
}

// findExpansions retrieves a set of words that could be considered expansions for the input string.
func findExpansions(input string, possibleExpansions expansion.Set) []string {
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
	for _, candidate := range exp.FindAllString(possibleExpansions.String(), -1) {
		if any(input, candidate, isTruncation, hasRemovedChar, hasRemovedVowels, hasRemovedCharAfterRemovedVowels) {
			expansions = append(expansions, candidate)
		}
	}

	return expansions
}

func any(abbr string, word string, filters ...filterFunc) bool {
	for _, filter := range filters {
		if filter(abbr, word) {
			return true
		}
	}

	return false
}

type similarityFunc func(string, string) float64

// similarityScore returns the Log of the similarity computed by the SimilarityCalculator.
//
// For two equal words, the similarity score is zero. If the probability is zero, we use
// a close-to-zero value to avoid issues with -Inf.
func similarityScore(calculator SimilarityCalculator, word string, anotherWord string) float64 {
	w1 := strings.ToLower(word)
	w2 := strings.ToLower(anotherWord)
	if w1 == w2 {
		return 0
	}

	prob := calculator.Similarity(w1, w2)
	if prob == 0 {
		prob = closeToZeroProbability
	}

	return math.Log(prob)
}

// cohesion computes the similarity score between an expansion and the other soft words on the potential splittings
// list, but for the same word (k not i).
func cohesion(similarity similarityFunc, ps potentialSplit, expansion string, idx int, context []string) float64 {
	var cohesion float64
	for i, softword := range ps.softwords {
		if idx == i {
			continue
		}

		// compute cohesion between the proposed expansion and the expansions for the others soft words
		for _, translation := range softword.expansions {
			cohesion += similarity(expansion, translation.translation)
		}
	}

	// add context cohesion
	for _, contextWord := range context {
		cohesion += similarity(expansion, contextWord)
	}

	return cohesion
}

// score calculates the score for a split. The score is the average similarities computed over all
// of pairs of expanded words and each expanded word paired with each context word.
// An average is used to avoid biasing the results toward excesive splitting.
func score(similarity similarityFunc, split potentialSplit, context []string) float64 {
	var expansionsScore float64
	expandedWords := marker.SplitBy(split.bestExpansion())
	for i, w1 := range expandedWords {
		var wordScore float64
		// add expansions similarities
		for j, w2 := range expandedWords {
			if i == j {
				continue
			}

			wordScore += similarity(w1, w2)
		}

		// add context similarities
		for _, contextWord := range context {
			wordScore += similarity(w1, contextWord)
		}

		expansionsScore += wordScore
	}

	n := float64(len(split.softwords))
	c := float64(len(context))

	return expansionsScore / (n * (n + c))
}
