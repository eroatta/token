package lists

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains_OnList_ShouldRetrieveElement(t *testing.T) {
	type args struct {
		elements []string
	}
	tests := []struct {
		name  string
		args  args
		token string
		want  bool
	}{
		{"empty_list", args{elements: []string{}}, "any", false},
		{"with_element", args{elements: []string{"word"}}, "word", true},
		{"with_element_case_insensitive", args{elements: []string{"WoRd"}}, "wOrD", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewBuilder().Add(tt.args.elements...).Build()
			got := list.Contains(tt.token)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSize_OnList_ShouldRetrieveListSize(t *testing.T) {
	type args struct {
		elements []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"empty_list", args{elements: []string{}}, 0},
		{"with_element", args{elements: []string{"word"}}, 1},
		{"with_repetead_element", args{elements: []string{"WoRd", "word"}}, 1},
		{"with_several_elements", args{elements: []string{"word", "Word", "diff"}}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewBuilder().Add(tt.args.elements...).Build()
			got := list.Size()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestElements_OnList_ShouldRetrieveElements(t *testing.T) {
	type args struct {
		elements []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty_list", args{elements: []string{}}, []string{}},
		{"with_element", args{elements: []string{"word"}}, []string{"word"}},
		{"with_repetead_element", args{elements: []string{"WoRd", "word"}}, []string{"word"}},
		{"with_several_elements", args{elements: []string{"word", "Word", "diff"}}, []string{"word", "diff"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewBuilder().Add(tt.args.elements...).Build()
			got := list.Elements()
			fmt.Println(got)

			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestNewBuilder_ShouldReturnListBuilder(t *testing.T) {
	got := NewBuilder()

	assert.NotNil(t, got, "ListBuilder shouldn't be nil")
	assert.IsType(t, new(listBuilder), got)
}

func TestBuild_OnListBuilder_ShouldIncludeProvidedElements(t *testing.T) {
	type args struct {
		elems     []string
		moreElems []string
	}
	tests := []struct {
		name  string
		args  args
		token string
		want  bool
	}{
		{"empty_list", args{}, "any", false},
		{"with_elems", args{elems: []string{"elems"}}, "elems", true},
		{"with_elems_and_more_elems", args{elems: []string{"elems"}, moreElems: []string{"elems", "moreElems"}}, "elems", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder()
			if tt.args.elems != nil && len(tt.args.elems) > 0 {
				builder.Add(tt.args.elems...)
			}

			if tt.args.moreElems != nil && len(tt.args.moreElems) > 0 {
				builder.Add(tt.args.moreElems...)
			}

			list := builder.Build()
			got := list.Contains(tt.token)

			assert.Equal(t, tt.want, got)
		})
	}
}
