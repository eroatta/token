package basic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBasic_WithLists_ShouldReturnBasicWithGivenLists(t *testing.T) {
	srcWords := make(map[string]interface{}, 0)
	srcPhrases := make(map[string]string, 0)
	stopList := make(map[string]interface{}, 0)
	dicc := make(map[string]interface{}, 0)

	got := NewBasic(srcWords, srcPhrases, stopList, dicc)

	assert.Equal(t, "", *got.words)
	assert.Equal(t, "", *got.stopAndDicc)
	assert.Equal(t, srcPhrases, got.srcPhrases)
}

func TestExpand_OnBasic_ShouldReturnExpansion(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected []string
	}{
		{"no_expansion", "noExpansion", []string{}},
		{"in_words", "parser", []string{"parser"}},
		{"phrase", "JSON", []string{"java", "script", "object", "notation"}},
		{"in_lower_case", "case", []string{"case"}},
		{"case_unsensitive", "Case", []string{"case"}},
	}

	words := map[string]interface{}{
		"parser": true,
	}

	phraseList := map[string]string{
		"json": "java-script-object-notation",
	}

	stopList := map[string]interface{}{
		"case": true,
	}

	basic := NewBasic(words, phraseList, stopList, nil)
	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			got := basic.Expand(fixture.token)

			assert.ElementsMatch(t, fixture.expected, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func BenchmarkBasicExpansion(b *testing.B) {
	// TODO: add expansion tests...
}
