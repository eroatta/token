package samurai

import (
	"testing"

	"github.com/eroatta/token/lists"
	"github.com/stretchr/testify/assert"
)

func TestSplit_ShouldReturnValidSplits(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{"no_split", "car", "car"},
		{"by_lower_to_upper_case", "getString", "get string"},
		{"by_upper_to_lower_case", "GPSstate", "gps state"},
		{"with_upper_case_and_softword_starting_with_upper_case", "ASTVisitor", "ast visitor"},
		{"lowercase_softword", "notype", "no type"},
		{"multiple_lowercase_softword", "astnotype", "ast no type"},
	}

	tCtx := NewTokenContext(createTestFrequencyTable(), createTestGlobalFrequencyTable())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.token, tCtx, lists.Prefixes, lists.Suffixes)

			assert.Equal(t, tt.want, got, "elements should match in number and order")
		})
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
	tCtx := NewTokenContext(createTestFrequencyTable(), createTestGlobalFrequencyTable())

	for i := 0; i < b.N; i++ {
		Split("notype", tCtx, lists.Prefixes, lists.Suffixes)
	}
}
