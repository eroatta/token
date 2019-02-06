package splitters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHighestCohesion_OnSoftwordWithoutExpansions_ShouldReturnZero(t *testing.T) {
	softword := softword{
		word:       "foo",
		expansions: []expansion{},
	}

	got := softword.highestCohesion()

	assert.Equal(t, float64(0), got)
}

func TestHighestCohesion_OnSoftwordWithOnlyOneExpansion_ShouldReturnTheOnlyCohesionValue(t *testing.T) {
	softword := softword{
		word: "foo",
		expansions: []expansion{
			{"floor", 1.2345},
		},
	}

	got := softword.highestCohesion()

	assert.Equal(t, 1.2345, got)
}

func TestHighestCohesion_OnSoftword_ShouldReturnTheHighestCohesionOnAnyExpansion(t *testing.T) {
	softword := softword{
		word: "foo",
		expansions: []expansion{
			{"floor", 1.2345},
			{"foot", 2.1123},
			{"football", -0.1123},
		},
	}

	got := softword.highestCohesion()

	assert.Equal(t, 2.1123, got)
}

func TestHighestCohesion_OnEmptyPotentialSplit_ShouldReturnZero(t *testing.T) {
	potentialSplit := potentialSplit{}

	got := potentialSplit.highestCohesion()

	assert.Equal(t, float64(0), got)
}

func TestHighestCohesion_OnPotentialSplitWithOneSoftword_ShouldReturnTheSoftwordCohesion(t *testing.T) {
	uniqueSoftword := softword{
		word: "bar",
		expansions: []expansion{
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
		expansions: []expansion{
			{"bar", 1.2345},
		},
	}

	secondSoftword := softword{
		word: "bum",
		expansions: []expansion{
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

/*
func TestNewPotentialSplit_OnMarkedHardword_ShouldReturnPotentialSplit(t *testing.T) {
	got := newPotentialSplit("foo_bar")

	assert.Equal(t, "foo_bar", got.split, "split should match de input")
	assert.ElementsMatch(t, []string{"foo", "bar"}, got.softwords, "elements should match")
	//assert.Equal(t, 0, len(got.expansions), "expansions map should be empty")
}

func TestNewPotentialSplit_OnEmptyHardword_ShouldReturnEmptyPotentialSplit(t *testing.T) {
	got := newPotentialSplit("")

	assert.Equal(t, "", got.split, "split should be empty")
	assert.ElementsMatch(t, []string{}, got.softwords, "there should be no softwords")
	//assert.Equal(t, 0, len(got.expansions), "there should be no elements")
}
*/
