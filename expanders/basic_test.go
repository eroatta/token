package expanders

import (
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

/*func TestBasicExpansion(t *testing.T) {
	cases := []struct {
		token    string
		expected []string
	}{}

	basic := NewBasic(nil, nil)
	for _, c := range cases {
		got, err := basic.Expand(c.token)
		if err != nil {
			assert.Fail(t, "we shouldn't get any errors at this point", err)
		}

		assert.ElementsMatch(t, c.expected, got, "elements should match")
	}
}*/

func BenchmarkBasicExpansion(b *testing.B) {

}
