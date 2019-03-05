package splitters

import (
	"math"
	"regexp"
	"strings"
)

const (
	closeToZeroProbability = 0.000000000000001
)

// GenTest represents a Generation and Tests splitting algorithm, proposed by Lawrie, Binkley and Morrell.
type GenTest struct {
	// TODO review
	list          string
	simCalculator SimCalculator
	context       []string
	dicctionary   map[string]interface{}
}

// NewGenTest creates a new GenTest splitter.
func NewGenTest() *GenTest {
	return &GenTest{}
}

// SimCalculator is the interface that wraps the basic Sim method.
//
// Sim calculates the probability for two words to be co-located on the same sentence.
type SimCalculator interface {
	Sim(firstWord string, secondWord string) (prob float64)
}

// Split on GenTest receives a token and returns an array of hard/soft words,
// split by the Generation and Test algorithm proposed by Lawrie, Binkley and Morrell.
//
// The algorithm splits the token by its markers. Then, for each new token, checks if
// it is a dicctionary word or not. If it is a dicctionary word, then no more processing is done.
// But if the token is not a dicctionary word, then the generate and testing algorithm starts.
//
// The first step is to retrieve the list of potential splits, each potential split composed by a
// list of softwords.
// For each softword on the potential split, the algorithm looks for a set of expansions, and then
// calculates the cohesion.
// Cohesion is calculated for every expansion on every softword of each potential split, comparing the expansion
// with the rest of the expansions on every softword belonging to the same potential split.
//
// The potential split with the highest score is the selected split.
func (g *GenTest) Split(token string) ([]string, error) {
	preprocessedToken := addMarkersOnDigits(token)
	preprocessedToken = addMarkersOnLowerToUpperCase(preprocessedToken)

	splitToken := make([]string, 0, 10)
	for _, tok := range splitOnMarkers(preprocessedToken) {
		// discard one letter tokens and dictionary words
		if len(tok) == 1 || (g.dicctionary)[strings.ToLower(tok)] != nil {
			splitToken = append(splitToken, tok)
			continue
		}

		potentialSplits := make([]potentialSplit, 0)
		for _, pSplit := range generatePotentialSplits(tok) {
			// for each potential split softword find the list of expansions
			for i := 0; i < len(pSplit.softwords); i++ {
				// if no expansion is found, the input itself must be considered an expansion
				expansions := g.findExpansions(pSplit.softwords[i].word)
				if len(expansions) == 0 {
					expansions = append(expansions, pSplit.softwords[i].word)
				}

				for _, translation := range expansions {
					pSplit.softwords[i].expansions = append(pSplit.softwords[i].expansions, expansion{translation, 0})
				}
			}

			// for each expansion of every softword of the potential split, calculate the cohesion with the rest
			// of the softwords on the potential split
			for i := 0; i < len(pSplit.softwords); i++ {
				for j := 0; j < len(pSplit.softwords[i].expansions); j++ {
					cohesion := g.cohesion(pSplit, pSplit.softwords[i].expansions[j].translation, i)
					pSplit.softwords[i].expansions[j].cohesion = cohesion
				}
			}

			// calculate the score
			pSplit.score = g.score(pSplit)

			potentialSplits = append(potentialSplits, pSplit)
		}

		tokenBestSplit := findBestSplit(potentialSplits)
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

	expansions := make([]string, 0)
	for _, candidate := range exp.FindAllString(g.list, -1) {
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

// similarityScore returns the Log of the similarity computed by the SimCalculator.
//
// For two equal words, the similarity score is zero. If the probability is zero, we use
// a close-to-zero value to avoid issues with -Inf.
func (g *GenTest) similarityScore(word string, anotherWord string) float64 {
	w1 := strings.ToLower(word)
	w2 := strings.ToLower(anotherWord)
	if w1 == w2 {
		return 0
	}

	prob := g.simCalculator.Sim(w1, w2)
	if prob == 0 {
		prob = closeToZeroProbability
	}

	return math.Log(prob)
}

// cohesion computes the similarity score between an expansion and the other soft words on the potential splittings
// list, but for the same word (k not i).
func (g *GenTest) cohesion(potentialSplit potentialSplit, expansion string, index int) (cohesion float64) {
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

// score calculates the score for a split. The score is the average similarities computed over all
// of pairs of expanded words and each expanded word paired with each context word.
// An average is used to avoid biasing the results toward excesive splitting.
func (g *GenTest) score(split potentialSplit) float64 {
	var expansionsScore float64
	expandedWords := splitOnMarkers(split.bestExpansion())
	for i, w1 := range expandedWords {
		var wordScore float64
		// add expansions similarities
		for j, w2 := range expandedWords {
			if i == j {
				continue
			}

			wordScore += g.similarityScore(w1, w2)
		}

		// add context similarities
		for _, contextWord := range g.context {
			wordScore += g.similarityScore(w1, contextWord)
		}

		expansionsScore += wordScore
	}

	n := float64(len(split.softwords))
	c := float64(len(g.context))

	return expansionsScore / (n * (n + c))
}
