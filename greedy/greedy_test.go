package greedy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewListBuilder_ShouldReturnListBuilder(t *testing.T) {
	got := NewListBuilder()

	assert.NotNil(t, got, "ListBuilder shouldn't be nil")
}

func TestBuild_OnListBuilder_ShouldIncludeProvidedWords(t *testing.T) {
	type args struct {
		dicc       []string
		knownAbbrs []string
		stopList   []string
	}
	tests := []struct {
		name  string
		args  args
		token string
		want  bool
	}{
		{"empty_list", args{}, "any", false},
		{"with_dicctionary", args{dicc: []string{"dicc"}}, "dicc", true},
		{"with_known_abbrvs", args{knownAbbrs: []string{"known"}}, "known", true},
		{"with_stoplist", args{stopList: []string{"case"}}, "case", true},
		{"merged_list", args{dicc: []string{"merged"}, knownAbbrs: []string{"merged"}, stopList: []string{"merged"}}, "merged", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewListBuilder()
			if tt.args.dicc != nil && len(tt.args.dicc) > 0 {
				builder.Dicctionary(tt.args.dicc)
			}

			if tt.args.knownAbbrs != nil && len(tt.args.knownAbbrs) > 0 {
				builder.KnownAbbreviations(tt.args.knownAbbrs)
			}

			if tt.args.stopList != nil && len(tt.args.stopList) > 0 {
				builder.StopList(tt.args.stopList)
			}

			list := builder.Build()
			got := list.Contains(tt.token)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSplit_ShouldReturnValidSplits(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  []string
	}{
		{"no_split", "car", []string{"car"}},
		{"by_lower_to_upper_case", "getString", []string{"get", "String"}},
		{"by_upper_to_lower_case", "GPSstate", []string{"GPS", "state"}},
		{"with_upper_case_and_softword_starting_with_upper_case", "ASTVisitor", []string{"AST", "Visitor"}},
		{"lowercase_softword", "notype", []string{"no", "type"}},
	}

	dicc := []string{"get", "string", "gps", "state", "ast", "visitor", "no", "type"}
	list := NewListBuilder().Dicctionary(dicc).Build()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.token, list)

			assert.Equal(t, tt.want, got, "elements should match in number and order")
		})
	}
}

func BenchmarkGreedySplitting(b *testing.B) {
	dicc := []string{"get", "string", "gps", "state", "ast", "visitor", "no", "type"}
	list := NewListBuilder().Dicctionary(dicc).Build()
	for i := 0; i < b.N; i++ {
		Split("GPSstate", list)
	}
}
