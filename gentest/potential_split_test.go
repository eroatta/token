package gentest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPotentialSplit_OnMarkedHardword_ShouldReturnPotentialSplit(t *testing.T) {
	got := newPotentialSplit("foo_bar")

	assert.Equal(t, "foo_bar", got.split)
	assert.Equal(t, 2, len(got.softwords))

	expectedSoftwords := []softword{
		{"foo", make([]possibleExpansion, 0)},
		{"bar", make([]possibleExpansion, 0)},
	}
	assert.ElementsMatch(t, expectedSoftwords, got.softwords)
}

func TestNewPotentialSplit_OnEmptyHardword_ShouldReturnEmptyPotentialSplit(t *testing.T) {
	got := newPotentialSplit("")

	assert.Equal(t, "", got.split)
	assert.Equal(t, 0, len(got.softwords))
}
func TestHighestCohesion_OnSoftwordWithoutExpansions_ShouldReturnZero(t *testing.T) {
	softword := softword{
		word:       "foo",
		expansions: []possibleExpansion{},
	}

	got := softword.highestCohesion()

	assert.Equal(t, float64(0), got)
}

func TestHighestCohesion_OnSoftwordWithOnlyOneExpansion_ShouldReturnTheOnlyCohesionValue(t *testing.T) {
	softword := softword{
		word: "foo",
		expansions: []possibleExpansion{
			{"floor", 1.2345},
		},
	}

	got := softword.highestCohesion()

	assert.Equal(t, 1.2345, got)
}

func TestHighestCohesion_OnSoftword_ShouldReturnTheHighestCohesionOnAnyExpansion(t *testing.T) {
	softword := softword{
		word: "foo",
		expansions: []possibleExpansion{
			{"floor", 1.2345},
			{"foot", 2.1123},
			{"football", -0.1123},
		},
	}

	got := softword.highestCohesion()

	assert.Equal(t, 2.1123, got)
}

func TestBestExpansion_OnSoftwordWithNoExpansions_ShouldReturnTheSoftword(t *testing.T) {
	softword := softword{
		word:       "foo",
		expansions: []possibleExpansion{},
	}

	got := softword.bestExpansion()

	assert.Equal(t, "foo", got)
}

func TestBestExpansion_OnSoftwordWithOnlyOneExpansion_ShouldReturnTheOnlyExpansion(t *testing.T) {
	softword := softword{
		word: "foo",
		expansions: []possibleExpansion{
			{"floor", 1.2345},
		},
	}

	got := softword.bestExpansion()

	assert.Equal(t, "floor", got)
}

func TestBesExpansion_OnSoftwordWithSeveralExpansions_ShouldReturnTheExpansionWithHighestCohesion(t *testing.T) {
	softword := softword{
		word: "foo",
		expansions: []possibleExpansion{
			{"floor", 1.2345},
			{"foot", 2.1123},
			{"football", -0.1123},
		},
	}

	got := softword.bestExpansion()

	assert.Equal(t, "foot", got)
}

func TestHighestCohesion_OnEmptyPotentialSplit_ShouldReturnZero(t *testing.T) {
	potentialSplit := potentialSplit{}

	got := potentialSplit.highestCohesion()

	assert.Equal(t, float64(0), got)
}

func TestHighestCohesion_OnPotentialSplitWithOneSoftword_ShouldReturnTheSoftwordCohesion(t *testing.T) {
	uniqueSoftword := softword{
		word: "bar",
		expansions: []possibleExpansion{
			{"bar", 1.2345},
		},
	}

	potentialSplit := potentialSplit{
		split:     "bar",
		softwords: []softword{uniqueSoftword},
	}

	got := potentialSplit.highestCohesion()

	assert.Equal(t, 1.2345, got)
}

func TestHighestCohesion_OnPotentialSplitWithSeveralSoftwords_ShouldReturnTheirHighestCohesion(t *testing.T) {
	firstSoftword := softword{
		word: "bar",
		expansions: []possibleExpansion{
			{"bar", 1.2345},
		},
	}

	secondSoftword := softword{
		word: "bum",
		expansions: []possibleExpansion{
			{"bump", 0.5551},
			{"bumpy", 0.1999},
		},
	}

	potentialSplit := potentialSplit{
		split:     "bar_bum",
		softwords: []softword{firstSoftword, secondSoftword},
	}

	got := potentialSplit.highestCohesion()

	assert.Equal(t, 1.7896, got)
}

func TestBestExpansion_OnEmptyPotentialSplit_ShouldReturnEmptyString(t *testing.T) {
	potentialSplit := potentialSplit{}

	got := potentialSplit.bestExpansion()

	assert.Equal(t, "", got)
}

func TestBestExpansion_OnPotentialSplitWithOneSoftword_ShouldReturnTheSoftwordBestExpansion(t *testing.T) {
	uniqueSoftword := softword{
		word: "bar",
		expansions: []possibleExpansion{
			{"bar", 1.2345},
		},
	}

	potentialSplit := potentialSplit{
		split:     "bar",
		softwords: []softword{uniqueSoftword},
	}

	got := potentialSplit.bestExpansion()

	assert.Equal(t, "bar", got)
}

func TestBestExpansion_OnPotentialSplitWithSeveralSoftwords_ShouldReturnTheSplitExpansion(t *testing.T) {
	firstSoftword := softword{
		word: "bar",
		expansions: []possibleExpansion{
			{"bar", 1.2345},
		},
	}

	secondSoftword := softword{
		word: "bum",
		expansions: []possibleExpansion{
			{"bump", 0.5551},
			{"bumpy", 0.1999},
		},
	}

	potentialSplit := potentialSplit{
		split:     "bar_bum",
		softwords: []softword{firstSoftword, secondSoftword},
	}

	got := potentialSplit.bestExpansion()

	assert.Equal(t, "bar_bump", got)
}

func TestFindBestSplit_OnEmptyPotentialSplitsList_ShouldReturnEmptySplit(t *testing.T) {
	var potentialSplits []potentialSplit

	got := findBestSplit(potentialSplits)

	assert.Equal(t, potentialSplit{}, got)
}

func TestFindBestSplit_OnPotentialSplitsListWithOneItem_ShouldReturnTheOnlySplit(t *testing.T) {
	pSplit := potentialSplit{
		split: "bar",
		softwords: []softword{
			{"bar", []possibleExpansion{{"bar", 1.2345}}},
		},
		score: 1.0012,
	}

	got := findBestSplit([]potentialSplit{pSplit})

	assert.Equal(t, pSplit, got)
}

func TestFindBestSplit_OnPotentialSplitsList_ShouldReturnTheSplitWithHighestCohesion(t *testing.T) {
	bestSplit := potentialSplit{
		split: "str_len",
		softwords: []softword{
			{"str", []possibleExpansion{{"string", 2.3432}}},
			{"len", []possibleExpansion{{"length", 2.0011}}},
		},
		score: 1.2131,
	}

	notSoBadSplit := potentialSplit{
		split: "st_rlen",
		softwords: []softword{
			{"st", []possibleExpansion{{"string", 1.9432}}},
			{"rlen", []possibleExpansion{{"riflemen", 0.9011}}},
		},
		score: 1.0013,
	}

	badSplit := potentialSplit{
		split: "s_trlen",
		softwords: []softword{
			{"s", []possibleExpansion{}},
			{"trlen", []possibleExpansion{}},
		},
		score: 0.5541,
	}

	got := findBestSplit([]potentialSplit{badSplit, bestSplit, notSoBadSplit})

	assert.Equal(t, bestSplit, got)
}
