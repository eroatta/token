package expanders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBasicWithNilListsReturnsNewBasicWithDefaultLists(t *testing.T) {

}

func TestBasicExpansion(t *testing.T) {
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
}

func BenchmarkBasicExpansion(b *testing.B) {

}
