package gentest

import (
	"fmt"
	"testing"

	"github.com/eroatta/token/expansion"
	"github.com/eroatta/token/lists"

	"math"

	"github.com/stretchr/testify/assert"
)

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
		{"with_markers_and_lowercase", "type_notype", []string{"type", "no", "type"}},
	}

	dict := lists.NewBuilder().
		Add("car", "get", "string", "no", "not", "notary", "type", "typo").Build()

	similarityCalculatorMock := similarityCalculatorMock{
		"car-gps":     0.9999,
		"gps-status":  0.9999,
		"state-tire":  0.9001,
		"ast-visitor": 1.0,
		"no-type":     0.8564,
		"no-typo":     0.0001,
	}
	context := lists.NewBuilder().
		Add("none", "no", "never", "nine", "type", "typeset").
		Add("typhoon", "tire", "car", "boo", "foo").Build()
	expansionsSet := expansion.NewSetBuilder().AddList(dict).Build()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.token, similarityCalculatorMock, context, expansionsSet)

			assert.Equal(t, tt.want, got, "elements should match in number and order")
		})
	}
}

func TestGeneratePotentialSplits_ShouldReturnEveryPossibleCombination(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  []string
	}{
		{"car_token", "car", []string{"car", "c_ar", "c_a_r", "ca_r"}},
		{"bond_token", "bond", []string{"bond", "b_ond", "b_o_nd", "b_on_d", "bo_nd", "bo_n_d", "bon_d"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []string
			for _, potentialSplit := range generatePotentialSplits(tt.token) {
				got = append(got, potentialSplit.split)
			}

			assert.ElementsMatch(t, tt.want, got, "elements should match")
		})
	}
}

func TestFindExpansions_ShouldReturnAllMatches(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"input_st", "st", []string{"string", "steer", "set"}},
		{"input_rlen", "rlen", []string{}},
		{"input_str", "str", []string{"steer", "string"}},
		{"input_len", "len", []string{"lender", "length"}},
		{"empty_input", "", []string{}},
		{"blankspace_input", " ", []string{}},
	}

	dict := lists.NewBuilder().
		Add("car", "string", "steer", "set", "riflemen", "lender", "bar", "length", "kamikaze").Build()
	expansionsSet := expansion.NewSetBuilder().AddList(dict).Build()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findExpansions(tt.input, expansionsSet)

			assert.ElementsMatch(t, tt.want, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestSimilarityScore_ShouldReturnSimilarity(t *testing.T) {
	tests := []struct {
		name          string
		simCalculator similarityCalculatorMock
		word1         string
		word2         string
		want          float64
	}{
		{"equal_words", nil, "car", "car", float64(0)},
		{"high_probability", similarityCalculatorMock{"car-wheel": 0.8211}, "car", "wheel", math.Log(0.8211)},
		{"no_probability", similarityCalculatorMock{}, "disco", "egypt", math.Log(closeToZeroProbability)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := similarityScore(tt.simCalculator, tt.word1, tt.word2)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestScore_ShouldReturnSimilarity(t *testing.T) {
	tests := []struct {
		name          string
		split         potentialSplit
		simCalculator similarityCalculatorMock
		context       []string
		want          float64
	}{
		{"one_word_and_no_context",
			potentialSplit{
				split: "bar",
				softwords: []softword{
					{"bar", []possibleExpansion{{"bar", 1.2345}}},
				}},
			similarityCalculatorMock{},
			[]string{},
			0.0},
		{"two_words_and_no_context", potentialSplit{
			split: "str_len",
			softwords: []softword{
				{"str", []possibleExpansion{{"string", 2.3432}}},
				{"len", []possibleExpansion{{"length", 2.0011}}},
			}},
			similarityCalculatorMock{"length-string": 0.9123},
			[]string{},
			(math.Log(0.9123) + math.Log(0.9123)) / (2.0 * (2.0 + 0.0))},
		{"two_words_and_context",
			potentialSplit{
				split: "str_len",
				softwords: []softword{
					{"str", []possibleExpansion{{"string", 2.3432}}},
					{"len", []possibleExpansion{{"length", 2.0011}}},
				},
			}, similarityCalculatorMock{
				"length-string":        0.9123,
				"concatenation-string": 0.8912},
			[]string{"concatenation"},
			(math.Log(0.9123) + math.Log(0.8912) + math.Log(0.9123) + math.Log(closeToZeroProbability)) / (2.0 * (2.0 + 1.0))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simFunc := func(w1 string, w2 string) float64 {
				return similarityScore(tt.simCalculator, w1, w2)
			}
			got := score(simFunc, tt.split, tt.context)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExpand_ShouldReturnValidExpansions(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  []string
	}{
		{"no_split", "cr", []string{"car"}},
		{"by_lower_to_upper_case", "getString", []string{"get", "String"}},
		{"by_upper_to_lower_case", "GPSsttus", []string{"GPS", "status"}},
		{"with_upper_case_and_softword_starting_with_upper_case", "ASTVisitor", []string{"AST", "Visitor"}},
		{"lowercase_softword", "notype", []string{"no", "type"}},
		{"with_markers_and_lowercase", "type_notype", []string{"type", "no", "type"}},
	}

	dict := lists.NewBuilder().
		Add("car", "get", "string", "no", "not", "notary", "type", "typo", "status").Build()

	similarityCalculatorMock := similarityCalculatorMock{
		"car-gps":     0.9999,
		"gps-status":  0.9999,
		"state-tire":  0.9001,
		"ast-visitor": 1.0,
		"no-type":     0.8564,
		"no-typo":     0.0001,
	}
	context := lists.NewBuilder().
		Add("none", "no", "never", "nine", "type", "typeset").
		Add("typhoon", "tire", "car", "boo", "foo", "state", "states").Build()
	expansionsSet := expansion.NewSetBuilder().AddList(dict).Build()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Expand(tt.token, similarityCalculatorMock, context, expansionsSet)

			assert.Equal(t, tt.want, got, "elements should match in number and order")
		})
	}
}

// mocks
type similarityCalculatorMock map[string]float64

func (s similarityCalculatorMock) Similarity(word string, another string) float64 {
	var key string
	if word < another {
		key = word + "-" + another
	} else {
		key = another + "-" + word
	}

	return s[key]
}

// end of mocks

func BenchmarkGenerate(b *testing.B) {
	benchmarks := []struct {
		name  string
		token string
	}{
		{"Short token", "car"},
		{"Medium token", "numsize"},
		{"Long token", "allocatedsize"},
		{"Longest token", "veryverylongtokennameforsplitting"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				generatePotentialSplits(bm.token)
			}
		})
	}
}
