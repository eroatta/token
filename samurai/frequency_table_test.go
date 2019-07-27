package samurai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrequency_OnEmptyFrequencyTable_ShouldReturnZero(t *testing.T) {
	ft := NewFrequencyTable()
	got := ft.Frequency("test")

	assert.Equal(t, 0.0, got, "frequency for any token on an empty table should be 0")
}

func TestTotalOccurrences_OnEmptyFrequencyTable_ShouldReturnZero(t *testing.T) {
	ft := NewFrequencyTable()
	got := ft.TotalOccurrences()

	assert.Equal(t, 0, got, "total number of occurrences on an empty table should be 0")
}

func TestSetOccurrences_WithNegativeValue_ShouldReturnError(t *testing.T) {
	ft := NewFrequencyTable()
	err := ft.SetOccurrences("token", -1)

	assert.Error(t, err, "an error should occur when trying to set a count with a negative value")
}

func TestSetOccurrences_WithNewAndOnlyToken_ShouldAddTheTokenAndSetTotal(t *testing.T) {
	ft := NewFrequencyTable()
	ft.SetOccurrences("token", 3)

	freq := ft.Frequency("token")
	total := ft.TotalOccurrences()

	assert.Equal(t, 1.0, freq, "frequency for the given token should match")
	assert.Equal(t, 3, total, "total number of occurrences should match")
}

func TestSetOccurrences_WithNewTokenOnExistingTable_ShouldAddTheTokenAndUpdateTotal(t *testing.T) {
	ft := NewFrequencyTable()
	ft.SetOccurrences("token", 3)

	ft.SetOccurrences("new-token", 2)
	freq := ft.Frequency("new-token")
	total := ft.TotalOccurrences()

	assert.Equal(t, 0.4, freq, "frequency for the given token should match")
	assert.Equal(t, 5, total, "total number of occurrences should match")
}

func TestSetOccurrences_WithExistingTokenOnTable_ShouldIncreaseTokenAndTotalValues(t *testing.T) {
	ft := NewFrequencyTable()
	ft.SetOccurrences("token", 3)

	ft.SetOccurrences("existing-token", 1)
	freq := ft.Frequency("existing-token")
	total := ft.TotalOccurrences()

	assert.Equal(t, 0.25, freq, "frequency for the given token should match")
	assert.Equal(t, 4, total, "total number of occurrences should match")

	ft.SetOccurrences("existing-token", 2)
	freq = ft.Frequency("existing-token")
	total = ft.TotalOccurrences()

	assert.Equal(t, 0.4, freq, "frequency for the given token should match")
	assert.Equal(t, 5, total, "total number of occurrences should match")
}

func TestSetOccurrences_WithExistingTokenOnTable_ShouldDecreaseTokenAndTotalValues(t *testing.T) {
	ft := NewFrequencyTable()
	ft.SetOccurrences("token", 3)

	ft.SetOccurrences("existing-token", 2)
	freq := ft.Frequency("existing-token")
	total := ft.TotalOccurrences()

	assert.Equal(t, 0.4, freq, "frequency for the given token should match")
	assert.Equal(t, 5, total, "total number of occurrences should match")

	ft.SetOccurrences("existing-token", 1)
	freq = ft.Frequency("existing-token")
	total = ft.TotalOccurrences()

	assert.Equal(t, 0.25, freq, "frequency for the given token should match")
	assert.Equal(t, 4, total, "total number of occurrences should match")
}
