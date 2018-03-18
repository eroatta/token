package splitters

import (
	"testing"

	"github.com/deckarep/golang-set"

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
		{1, "getString", []string{"get", "String"}},
		{2, "GPSstate", []string{"GPS", "state"}},
		{3, "ASTVisitor", []string{"AST", "Visitor"}},
		{4, "notype", []string{"no", "type"}},
		{5, "astnotype", []string{"ast", "no", "type"}},
	}

	freqTable := createTestFrequencyTable()
	globalFreqTable := createTestGlobalFrequencyTable()
	prefixes := createTestPrefixes()
	suffixes := createTestSuffixes()
	samurai := NewSamurai(freqTable, globalFreqTable, prefixes, suffixes)
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
	ft.SetCount("get", 3)
	ft.SetCount("string", 10)
	ft.SetCount("gets", 1)
	ft.SetCount("ring", 3)

	ft.SetCount("gps", 12)
	ft.SetCount("gp", 1)
	ft.SetCount("state", 22)

	ft.SetCount("ast", 2)
	ft.SetCount("visitor", 1)

	ft.SetCount("no", 5)
	ft.SetCount("not", 4)
	ft.SetCount("type", 5)

	return ft
}

func createTestGlobalFrequencyTable() *FrequencyTable {
	ft := NewFrequencyTable()
	ft.SetCount("get", 100)
	ft.SetCount("string", 200)
	ft.SetCount("gets", 150)
	ft.SetCount("ring", 15)

	ft.SetCount("gps", 98)
	ft.SetCount("gp", 13)
	ft.SetCount("state", 224)

	ft.SetCount("no", 22)
	ft.SetCount("not", 63)
	ft.SetCount("type", 112)

	return ft
}

func createTestPrefixes() *mapset.Set {
	return &defaultPrefixes
}

func createTestSuffixes() *mapset.Set {
	return &defaultSuffixes
}

func BenchmarkSamuraiSplitting(b *testing.B) {
	freqTable := createTestFrequencyTable()
	globalFreqTable := createTestGlobalFrequencyTable()
	prefixes := createTestPrefixes()
	suffixes := createTestSuffixes()

	samurai := NewSamurai(freqTable, globalFreqTable, prefixes, suffixes)
	for i := 0; i < b.N; i++ {
		samurai.Split("notype")
	}
}
