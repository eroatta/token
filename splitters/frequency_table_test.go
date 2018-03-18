package splitters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFrequencyZeroForAnyTokenOnEmptyTable(t *testing.T) {
	ft := NewFrequencyTable()
	got := ft.Frequency("test")

	assert.Equal(t, 0.0, got, "frequency for any token on an empty table should be 0")
}

func TestTotalOccurrencesZeroOnEmptyTable(t *testing.T) {
	ft := NewFrequencyTable()
	got := ft.TotalOccurrences()

	assert.Equal(t, 0, got, "total number of occurrences on an empty table should be 0")
}

func TestErrorWhenSetCountWithNegativeValue(t *testing.T) {
	ft := NewFrequencyTable()
	err := ft.SetCount("token", -1)

	assert.Error(t, err, "an error should occur when trying to set a count with a negative value")
}

func TestSetCountForNewStringOnEmptyTable(t *testing.T) {
	ft := NewFrequencyTable()
	ft.SetCount("token", 3)

	freq := ft.Frequency("token")
	total := ft.TotalOccurrences()

	assert.Equal(t, 1.0, freq, "frequency for the given token should match")
	assert.Equal(t, 3, total, "total number of occurrences should match")
}

func TestSetCountForNewStringOnTable(t *testing.T) {
	ft := NewFrequencyTable()
	ft.SetCount("token", 3)

	ft.SetCount("new-token", 2)
	freq := ft.Frequency("new-token")
	total := ft.TotalOccurrences()

	assert.Equal(t, 0.4, freq, "frequency for the given token should match")
	assert.Equal(t, 5, total, "total number of occurrences should match")
}

func TestIncreasingSetCountForExistingStringOnTable(t *testing.T) {
	ft := NewFrequencyTable()
	ft.SetCount("token", 3)

	ft.SetCount("existing-token", 1)
	freq := ft.Frequency("existing-token")
	total := ft.TotalOccurrences()

	assert.Equal(t, 0.25, freq, "frequency for the given token should match")
	assert.Equal(t, 4, total, "total number of occurrences should match")

	ft.SetCount("existing-token", 2)
	freq = ft.Frequency("existing-token")
	total = ft.TotalOccurrences()

	assert.Equal(t, 0.4, freq, "frequency for the given token should match")
	assert.Equal(t, 5, total, "total number of occurrences should match")
}

func TestDecreasingSetCountForExistingStringOnTable(t *testing.T) {
	ft := NewFrequencyTable()
	ft.SetCount("token", 3)

	ft.SetCount("existing-token", 2)
	freq := ft.Frequency("existing-token")
	total := ft.TotalOccurrences()

	assert.Equal(t, 0.4, freq, "frequency for the given token should match")
	assert.Equal(t, 5, total, "total number of occurrences should match")

	ft.SetCount("existing-token", 1)
	freq = ft.Frequency("existing-token")
	total = ft.TotalOccurrences()

	assert.Equal(t, 0.25, freq, "frequency for the given token should match")
	assert.Equal(t, 4, total, "total number of occurrences should match")
}
