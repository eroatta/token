package expanders

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAmap_ShouldReturnNewAmapInstance(t *testing.T) {
	// TODO: complete
}

func TestExpand_OnAmap_ShouldReturnExpansion(t *testing.T) {
	cases := []struct {
		name     string
		token    string
		expected []string
	}{
		{"no_expansion", "noExpansion", []string{}},
	}

	amap := NewAmap()
	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			got := amap.Expand(fixture.token)

			assert.ElementsMatch(t, fixture.expected, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}