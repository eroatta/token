package expanders

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

	assert.Equal(t, srcWords, got.srcWords)
	assert.Equal(t, srcPhrases, got.srcPhrases)
	assert.Equal(t, stopList, got.stopList)
	assert.Equal(t, dicc, got.dicctionary)
}

func TestExpand_OnBasic_ShouldReturnExpansion(t *testing.T) {
	tests := []struct {
		token    string
		expected []string
	}{
		{"noExpansion", []string{}},
		{"case", []string{"case"}},
		{"Case", []string{"case"}},
		{"JSON", []string{"java", "script", "object", "notation"}},
		{"parser", []string{"parser"}},
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
		got, err := basic.Expand(fixture.token)
		if err != nil {
			assert.Fail(t, "we shouldn't get any errors at this point", err)
		}

		assert.ElementsMatch(t, fixture.expected, got, fmt.Sprintf("found elements: %v", got))
	}
}

func BenchmarkBasicExpansion(b *testing.B) {

}
