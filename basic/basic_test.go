package basic

import (
	"fmt"
	"testing"

	"github.com/eroatta/token/expansion"

	"github.com/stretchr/testify/assert"
)

func TestExpand_OnBasic_ShouldReturnExpansion(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected []string
	}{
		{"no_expansion", "noExpansion", []string{}},
		{"in_words", "parser", []string{"parser"}},
		{"phrase", "JSON", []string{"java script object notation"}},
		{"in_lower_case", "case", []string{"case"}},
		{"case_unsensitive", "Case", []string{"case"}},
	}

	srcWords := expansion.NewSetBuilder().AddStrings("parser").Build()
	phraseList := map[string]string{
		"json": "java-script-object-notation",
	}
	stopList := expansion.NewSetBuilder().AddStrings("case").Build()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Expand(tt.token, srcWords, phraseList, stopList)

			assert.ElementsMatch(t, tt.expected, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func BenchmarkBasicExpansion(b *testing.B) {
	srcWords := expansion.NewSetBuilder().AddStrings("parser").Build()
	phraseList := map[string]string{
		"json": "java-script-object-notation",
	}

	for i := 0; i < b.N; i++ {
		Expand("rdy", srcWords, phraseList, DefaultExpansions)
	}
}
