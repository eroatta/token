package expansion

import (
	"testing"

	"github.com/eroatta/token/lists"
	"github.com/stretchr/testify/assert"
)

func TestNewSetBuilder_ShouldReturnSetBuilder(t *testing.T) {
	got := NewSetBuilder()

	assert.NotNil(t, got, "SetBuilder shouldn't be nil")
	assert.IsType(t, new(setBuilder), got)
}

func TestArray_OnSet_ShouldRetrieveArray(t *testing.T) {
	type args struct {
		list lists.List
		strs []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty_set", args{}, []string{}},
		{"from_list", args{list: lists.NewBuilder().Add("list").Build()}, []string{"list"}},
		{"from_strings", args{strs: []string{"string"}}, []string{"string"}},
		{"combined_from_list_strings",
			args{list: lists.NewBuilder().Add("list").Build(), strs: []string{"string"}},
			[]string{"list", "string"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSetBuilder().AddList(tt.args.list).AddStrings(tt.args.strs...).Build()
			got := set.Array()

			assert.ElementsMatch(t, tt.want, got, "elements should match")
		})
	}
}

func TestString_OnSet_ShouldRetrieveArray(t *testing.T) {
	type args struct {
		list lists.List
		strs []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty_set", args{}, ""},
		{"from_list", args{list: lists.NewBuilder().Add("list").Build()}, "list"},
		{"from_strings", args{strs: []string{"string"}}, "string"},
		{"combined_from_list_strings",
			args{list: lists.NewBuilder().Add("list").Build(), strs: []string{"string"}},
			"list string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSetBuilder().AddList(tt.args.list).AddStrings(tt.args.strs...).Build()
			got := set.String()

			assert.Equal(t, tt.want, got, "strings should match")
		})
	}
}

func TestContains_OnSet_ShouldRetrieveIfContainedOrNot(t *testing.T) {
	type args struct {
		list lists.List
		strs []string
	}
	tests := []struct {
		name string
		args args
		word string
		want bool
	}{
		{"empty_set", args{}, "list", false},
		{"from_list", args{list: lists.NewBuilder().Add("list").Build()}, "list", true},
		{"from_strings", args{strs: []string{"string"}}, "string", true},
		{"combined_from_list_strings",
			args{list: lists.NewBuilder().Add("list").Build(), strs: []string{"string"}},
			"set", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := NewSetBuilder().AddList(tt.args.list).AddStrings(tt.args.strs...).Build()
			got := set.Contains(tt.word)
			assert.Equal(t, tt.want, got)
		})
	}
}
