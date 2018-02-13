package splitters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreedySplitting(t *testing.T) {
	cases := []struct {
		ID       int
		token    string
		expected []string
	}{}

	greedy := new(Greedy)
	for _, c := range cases {
		got, err := greedy.Split(c.token)
		if err != nil {
			assert.Fail(t, "we shouldn't get any errors at this point", err)
		}

		assert.Equal(t, c.expected, got, "elements should match in number and order for identifier number")
	}
}

func BenchmarkGreedySplitting(b *testing.B) {
	greedy := new(Greedy)
	for i := 0; i < b.N; i++ {
		greedy.Split("spongebob_squarePants")
	}
}
