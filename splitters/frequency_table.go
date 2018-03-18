package splitters

import (
	"errors"
)

// FrequencyTable is a lookup table that stores the number of occurrences
// of each unique string in a set of strings.
type FrequencyTable struct {
	occurrences      map[string]int
	totalOccurrences int
}

// NewFrequencyTable creates and initializes an emtpy frequency table.
func NewFrequencyTable() *FrequencyTable {
	return &FrequencyTable{
		occurrences:      make(map[string]int),
		totalOccurrences: 0,
	}
}

// Frequency determines how frequently a token occurs in a set of strings.
func (f *FrequencyTable) Frequency(token string) float64 {
	if f.totalOccurrences == 0 {
		return 0.0
	}

	return float64(f.occurrences[token]) / float64(f.totalOccurrences)
}

// SetCount sets how many times a token occurs in a set of strings.
func (f *FrequencyTable) SetCount(token string, count int) error {
	if count < 0 {
		return errors.New("Count must be greater or equal than 0")
	}

	f.totalOccurrences = f.totalOccurrences + (count - f.occurrences[token])
	f.occurrences[token] = count
	return nil
}

// TotalOccurrences provides the total number of string occurrences on the
// frequency table.
func (f *FrequencyTable) TotalOccurrences() int {
	return f.totalOccurrences
}
