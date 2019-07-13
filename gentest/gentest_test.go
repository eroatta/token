package gentest

import (
	"fmt"
	"testing"

	"github.com/eroatta/token-splitex/lists"

	"math"

	"github.com/stretchr/testify/assert"
)

func TestSplit_OnGenTest_ShouldReturnValidSplits(t *testing.T) {
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
		{5, "type_notype", []string{"type", "no", "type"}},
	}

	dicc := lists.NewBuilder().
		Add("car", "get", "string", "no", "not", "notary", "type", "typo").Build()

	simCalculatorMock := simCalculatorMock{
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
	genTest := NewGenTest(simCalculatorMock)
	expansionsSet := NewExpansions(dicc)

	for _, c := range cases {
		got := Split(c.token, genTest, context, expansionsSet)

		assert.Equal(t, c.expected, got, "elements should match in number and order for identifier number")
	}
}

func TestGeneratePotentialSplits_ShouldReturnEveryPossibleCombination(t *testing.T) {
	cases := []struct {
		token    string
		expected []string
	}{
		{"car", []string{"car", "c_ar", "c_a_r", "ca_r"}},
		{"bond", []string{"bond", "b_ond", "b_o_nd", "b_on_d", "bo_nd", "bo_n_d", "bon_d"}},
	}

	for _, c := range cases {
		var got []string
		for _, potentialSplit := range generatePotentialSplits(c.token) {
			got = append(got, potentialSplit.split)
		}

		assert.ElementsMatch(t, c.expected, got, "elements should match")
	}
}

func TestFindExpansions_OnGenTestWithCustomList_ShouldReturnAllMatches(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expansions []string
	}{
		{"input_st", "st", []string{"string", "steer", "set"}},
		{"input_rlen", "rlen", []string{}},
		{"input_str", "str", []string{"steer", "string"}},
		{"input_len", "len", []string{"lender", "length"}},
		{"empty_input", "", []string{}},
		{"blankspace_input", " ", []string{}},
	}

	dicc := lists.NewBuilder().
		Add("car", "string", "steer", "set", "riflemen", "lender", "bar", "length", "kamikaze").Build()
	expansionsSet := NewExpansions(dicc)

	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			got := findExpansions(fixture.input, expansionsSet)

			assert.ElementsMatch(t, fixture.expansions, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestSimilarityScore_OnEqualWords_ShouldReturnZero(t *testing.T) {
	genTest := NewGenTest(nil)

	got := genTest.similarityScore("car", "car")

	assert.Equal(t, float64(0), got)
}

func TestSimilarityScore_OnWordsWithHighProb_ShouldReturnValue(t *testing.T) {
	simCalculatorMock := simCalculatorMock{
		"car-wheel": 0.8211,
	}
	genTest := NewGenTest(simCalculatorMock)

	got := genTest.similarityScore("car", "wheel")

	expected := math.Log(0.8211)
	assert.Equal(t, expected, got)
}

func TestSimilarityScore_OnDifferentWordsWithZeroProb_ShouldReturnCustomMinimalValue(t *testing.T) {
	genTest := NewGenTest(simCalculatorMock{})

	got := genTest.similarityScore("disco", "egypt")

	expected := math.Log(closeToZeroProbability)
	assert.Equal(t, expected, got)
}

func TestScore_OnOneWordSplitAndNoContext_ShouldReturnZero(t *testing.T) {
	split := potentialSplit{
		split: "bar",
		softwords: []softword{
			{"bar", []expansion{{"bar", 1.2345}}},
		},
	}

	emptyContext := lists.NewBuilder().Build()
	genTest := NewGenTest(nil)
	got := genTest.score(split, emptyContext.Elements())

	assert.Equal(t, 0.0, got)
}

func TestScore_OnTwoWordsSplitAndNoContext_ShouldReturnScore(t *testing.T) {
	split := potentialSplit{
		split: "str_len",
		softwords: []softword{
			{"str", []expansion{{"string", 2.3432}}},
			{"len", []expansion{{"length", 2.0011}}},
		},
	}

	simCalculatorMock := simCalculatorMock{
		"length-string": 0.9123,
	}
	emptyContext := lists.NewBuilder().Build()
	genTest := NewGenTest(simCalculatorMock)

	got := genTest.score(split, emptyContext.Elements())

	expected := (math.Log(0.9123) + math.Log(0.9123)) / (2.0 * (2.0 + 0.0))
	assert.Equal(t, expected, got)
}

func TestScore_OnTwoWordsSplitAndContext_ShouldReturnScore(t *testing.T) {
	split := potentialSplit{
		split: "str_len",
		softwords: []softword{
			{"str", []expansion{{"string", 2.3432}}},
			{"len", []expansion{{"length", 2.0011}}},
		},
	}

	simCalculatorMock := simCalculatorMock{
		"length-string":        0.9123,
		"concatenation-string": 0.8912,
	}
	context := lists.NewBuilder().Add("concatenation").Build()
	genTest := NewGenTest(simCalculatorMock)

	got := genTest.score(split, context.Elements())

	expected := (math.Log(0.9123) + math.Log(0.8912) + math.Log(0.9123) + math.Log(closeToZeroProbability)) / (2.0 * (2.0 + 1.0))
	assert.Equal(t, expected, got)
}

// mocks
type simCalculatorMock map[string]float64

func (s simCalculatorMock) Sim(word string, another string) float64 {
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
