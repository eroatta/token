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
	}{
		{0, "car", []string{"car"}},
		{1, "getString", []string{"get", "String"}},
		{2, "GPSstate", []string{"GPS", "state"}},
		{3, "ASTVisitor", []string{"AST", "Visitor"}},
		{4, "notype", []string{"no", "type"}},
	}

	dicc := createTestWords()
	emtpyList := make(map[string]interface{})
	greedy := NewGreedy(&dicc, &emtpyList, &emtpyList)
	for _, c := range cases {
		got, err := greedy.Split(c.token)
		if err != nil {
			assert.Fail(t, "we shouldn't get any errors at this point", err)
		}

		assert.Equal(t, c.expected, got, "elements should match in number and order for identifier number")
	}
}

func createTestWords() map[string]interface{} {
	dicc := make(map[string]interface{})
	dicc["get"] = 1
	dicc["string"] = 1
	dicc["gps"] = 1
	dicc["state"] = 1
	dicc["ast"] = 1
	dicc["visitor"] = 1
	dicc["no"] = 1
	dicc["type"] = 1

	return dicc
}

func BenchmarkGreedySplitting(b *testing.B) {
	greedy := new(Greedy)
	for i := 0; i < b.N; i++ {
		greedy.Split("spongebob_squarePants")
	}
}
