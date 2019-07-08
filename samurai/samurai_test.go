package samurai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSamurai_WithNilTables_ShouldReturnDefaultTables(t *testing.T) {
	samurai := NewSamurai(nil, nil, nil, nil)

	assert.NotNil(t, samurai.localFreqTable, "the local frequency table should be using a default table")
	assert.NotNil(t, samurai.globalFreqTable, "the global frequency table should be using a default table")
	assert.NotNil(t, samurai.prefixes, "the common prefixes set should be using a default set")
	assert.NotNil(t, samurai.suffixes, "the common suffixes set should be using a default set")
}

func TestSplit_OnSamurai_ShouldReturnValidSplits(t *testing.T) {
	cases := []struct {
		ID       int
		token    string
		expected []string
	}{
		{0, "car", []string{"car"}},
		{1, "getString", []string{"get", "String"}},
		{2, "GPSstate", []string{"GPS", "state"}},
		{3, "ASTVisitor", []string{"AST", "Visitor"}},
		{4, "notype", []string{"no", "type"}},
		{5, "astnotype", []string{"ast", "no", "type"}},
	}

	freqTable := createTestFrequencyTable()
	globalFreqTable := createTestGlobalFrequencyTable()
	// TODO: use prefixes and suffixes
	samurai := NewSamurai(freqTable, globalFreqTable, nil, nil)
	for _, c := range cases {
		got, err := samurai.Split(c.token)
		if err != nil {
			assert.Fail(t, "we shouldn't get any errors at this point", err)
		}

		assert.Equal(t, c.expected, got, "elements should match in number and order for identifier number")
	}
}

func createTestFrequencyTable() *FrequencyTable {
	ft := NewFrequencyTable()
	ft.SetOccurrences("get", 3)
	ft.SetOccurrences("string", 10)
	ft.SetOccurrences("gets", 1)
	ft.SetOccurrences("ring", 3)

	ft.SetOccurrences("gps", 12)
	ft.SetOccurrences("gp", 1)
	ft.SetOccurrences("state", 22)

	ft.SetOccurrences("ast", 2)
	ft.SetOccurrences("visitor", 1)

	ft.SetOccurrences("no", 5)
	ft.SetOccurrences("not", 4)
	ft.SetOccurrences("type", 5)

	return ft
}

func createTestGlobalFrequencyTable() *FrequencyTable {
	ft := NewFrequencyTable()
	ft.SetOccurrences("get", 100)
	ft.SetOccurrences("string", 200)
	ft.SetOccurrences("gets", 150)
	ft.SetOccurrences("ring", 15)

	ft.SetOccurrences("gps", 98)
	ft.SetOccurrences("gp", 13)
	ft.SetOccurrences("state", 224)

	ft.SetOccurrences("no", 22)
	ft.SetOccurrences("not", 63)
	ft.SetOccurrences("type", 112)

	return ft
}

func BenchmarkSamuraiSplitting(b *testing.B) {
	freqTable := createTestFrequencyTable()
	globalFreqTable := createTestGlobalFrequencyTable()

	// TODO: use prefixes and suffixes
	samurai := NewSamurai(freqTable, globalFreqTable, nil, nil)
	for i := 0; i < b.N; i++ {
		samurai.Split("notype")
	}
}
