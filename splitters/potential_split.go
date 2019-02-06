package splitters

import (
	"sort"
	"strings"
)

// potentialSplit represents a GenTest potential split. It holds data related to the split, the softwords
// and their expansions.
type potentialSplit struct {
	split     string
	softwords []softword
}

// softword represents a potential word and holds a set of related expansions.
type softword struct {
	word       string
	expansions []expansion
}

// expansion is a posible translation that holds an specific cohesion score.
type expansion struct {
	translation string
	cohesion    float64
}

// highestCohesion on a potential split returns the sum of the highest available expansion for each softword.
func (p potentialSplit) highestCohesion() float64 {
	var cohesion float64
	for _, softword := range p.softwords {
		cohesion += softword.highestCohesion()
	}

	return cohesion
}

// bestExpansion on a potential split returns the best expansion for each softword, combined and joined
// with underscores.
func (p potentialSplit) bestExpansion() string {
	expansions := make([]string, len(p.softwords))
	for i, softword := range p.softwords {
		expansions[i] = softword.bestExpansion()
	}

	return strings.Join(expansions, "_")
}

// highestCohesion on a softword returns the highest cohesion of any of its available translations.
func (s softword) highestCohesion() float64 {
	if len(s.expansions) == 0 {
		return 0
	}

	sort.Slice(s.expansions, func(i, j int) bool {
		return s.expansions[i].cohesion > s.expansions[j].cohesion
	})

	return s.expansions[0].cohesion
}

// bestExpansion on a softword returns the best expansion based on the cohesion value.
func (s softword) bestExpansion() string {
	if len(s.expansions) == 0 {
		return s.word
	}

	sort.Slice(s.expansions, func(i, j int) bool {
		return s.expansions[i].cohesion > s.expansions[j].cohesion
	})

	return s.expansions[0].translation
}

// findBestSplit looks for the potential split with the highest cohesion and selects it as the best
// potential split available.
func findBestSplit(potentialSplits []potentialSplit) potentialSplit {
	if len(potentialSplits) == 0 {
		return potentialSplit{}
	}

	sort.Slice(potentialSplits, func(i, j int) bool {
		return potentialSplits[i].highestCohesion() > potentialSplits[j].highestCohesion()
	})

	return potentialSplits[0]
}

// TODO check if we should have it or not...
// NewPotentialSplit creates and initializes a new potential split for the given hardword.
func newPotentialSplit(hardword string) potentialSplit {
	var softwords []softword
	if hardword != "" {
		for _, word := range splitOnMarkers(hardword) {
			softwords = append(softwords, softword{
				word:       word,
				expansions: make([]expansion, 0),
			})
		}
	}

	return potentialSplit{
		split:     hardword,
		softwords: softwords,
	}
}
