package samurai

import (
	"errors"
	"strings"
)

// errOccurrencesGreatherOrEqualThanZero indicates that occurrences on a frequency table must be always
// greater or equal than zero.
var errOccurrencesGreatherOrEqualThanZero = errors.New("Occurrences must be greater or equal than 0")

// FrequencyTable is a lookup table that stores the number of occurrences
// of each unique string in a set of strings.
type FrequencyTable struct {
	occurrences      map[string]int
	totalOccurrences int
}

// NewFrequencyTable creates and initializes an empty frequency table.
func NewFrequencyTable() *FrequencyTable {
	return &FrequencyTable{
		occurrences:      make(map[string]int),
		totalOccurrences: 0,
	}
}

// SetOccurrences sets how many times a token appeared in a set of strings.
func (f *FrequencyTable) SetOccurrences(token string, occurrences int) error {
	if occurrences < 0 {
		return errOccurrencesGreatherOrEqualThanZero
	}

	key := strings.ToLower(token)
	f.totalOccurrences = f.totalOccurrences + (occurrences - f.occurrences[key])
	f.occurrences[key] = occurrences

	return nil
}

// TotalOccurrences provides the total number of occurrences on the frequency table.
func (f FrequencyTable) TotalOccurrences() int {
	return f.totalOccurrences
}

// Frequency determines how frequently a token occurs in a set of strings.
func (f FrequencyTable) Frequency(token string) float64 {
	if f.totalOccurrences == 0 {
		return 0.0
	}

	return float64(f.occurrences[strings.ToLower(token)]) / float64(f.totalOccurrences)
}
