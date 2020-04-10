package conserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit_OnConserv_ShouldReturnValidSplits(t *testing.T) {
	cases := []struct {
		token    string
		expected string
	}{
		{"spongebob_squarePants", "spongebob square Pants"},
		{"Extraordinaire", "Extraordinaire"},
		{"extraordinairE", "extraordinair E"},
		{"extraordinaire", "extraordinaire"},
		{"extra_ordinaire", "extra ordinaire"},
		{"leto2nd", "leto 2 nd"},
		{"brooklyn99", "brooklyn 99"},
		{"mySQL", "my SQL"},
		{"mySql", "my Sql"},
		{"mySQl", "my S Ql"},
		{"9999", "9999"},
		{"", ""},
	}

	for _, c := range cases {
		got := Split(c.token)

		assert.Equal(t, c.expected, got, "elements should match in number and order")
	}
}

func BenchmarkConservSplitting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Split("spongebob_squarePants")
	}
}
