package splitters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSamuraiWithNilTablesReturnNewSamuraiWithDefaultTables(t *testing.T) {
	samurai := NewSamurai(nil, nil, nil, nil)

	assert.NotNil(t, samurai.localFreqTable, "the local frequency table should be using a default table")
	assert.NotNil(t, samurai.globalFreqTable, "the global frequency table should be using a default table")
	assert.NotNil(t, samurai.prefixes, "the common prefixes set should be using a default set")
	assert.NotNil(t, samurai.suffixes, "the common suffixes set should be using a default set")
}

func TestSamuraiSplitting(t *testing.T) {
	cases := []struct {
		ID       int
		token    string
		expected []string
	}{
		{0, "car", []string{"car"}},
	}

	freqTable := createTestFrequencyTable()
	globalFreqTable := createTestGlobalFrequencyTable()
	prefixes := createTestPrefixes()
	suffixes := createTestSuffixes()
	samurai := NewSamurai(&freqTable, &globalFreqTable, &prefixes, &suffixes)
	for _, c := range cases {
		got, err := samurai.Split(c.token)
		if err != nil {
			assert.Fail(t, "we shouldn't get any errors at this point", err)
		}

		assert.Equal(t, c.expected, got, "elements should match in number and order for identifier number")
	}
}

func createTestFrequencyTable() frequencyTable {
	return nil
}

func createTestGlobalFrequencyTable() frequencyTable {
	return nil
}

func createTestPrefixes() set {
	return nil
}

func createTestSuffixes() set {
	return nil
}

func BenchmarkSamuraiSplitting(b *testing.B) {

}
