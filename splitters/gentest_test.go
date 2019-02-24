package splitters

import (
	"fmt"
	"testing"

	"math"

	"github.com/stretchr/testify/mock"

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

	gentest := NewGenTest()
	// TODO: review gentest creation
	testDicc := map[string]interface{}{
		"car":    true,
		"get":    true,
		"string": true,
		"no":     true,
		"type":   true,
	}
	gentest.dicctionary = testDicc
	gentest.simCalculator = &simCalculatorMock{}
	gentest.list = "no not notary type typo"
	gentest.context = []string{"none", "no", "never", "nine", "type", "typeset", "typhoon", "tire", "car", "boo", "foo"}

	for _, c := range cases {
		got, err := gentest.Split(c.token)
		if err != nil {
			assert.Fail(t, "we shouldn't get any errors at this point", err)
		}

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
		{"input_rlen", "rlen", []string{"riflemen"}},
		{"input_str", "str", []string{"steer", "string"}},
		{"input_len", "len", []string{"lender", "length"}},
		{"empty_input", "", []string{}},
		{"blankspace_input", " ", []string{}},
	}

	gentest := NewGenTest()
	gentest.list = "car string steer set riflemen lender bar length kamikaze"
	for _, fixture := range tests {
		t.Run(fixture.name, func(t *testing.T) {
			got := gentest.findExpansions(fixture.input)

			assert.ElementsMatch(t, fixture.expansions, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestSimilarityScore_OnEqualWords_ShouldReturnZero(t *testing.T) {
	genTest := NewGenTest()

	got := genTest.similarityScore("car", "car")

	assert.Equal(t, float64(0), got)
}

// TODO: change test name
func TestSimilarityScore_OnDifferentWordsWithSomeProb_ShouldReturnValue(t *testing.T) {
	simCalculatorMock := new(simCalculatorMock)
	// TODO: review why it's not working
	simCalculatorMock.On("Sim", "car", "wheel").Return(1.2345)

	genTest := NewGenTest()
	genTest.simCalculator = simCalculatorMock

	got := genTest.similarityScore("car", "wheel")

	expected := math.Log(1.2345)
	assert.Equal(t, expected, got)
}

func TestSimilarityScore_OnDifferentWordsWithZeroProb_ShouldReturnCustomMinimalValue(t *testing.T) {
	genTest := NewGenTest()
	genTest.simCalculator = &simCalculatorMock{}

	got := genTest.similarityScore("disco", "egypt")

	expected := math.Log(0.000000000998587)
	assert.Equal(t, expected, got)
}

// mocks
type simCalculatorMock struct {
	mock.Mock
}

func (s *simCalculatorMock) Sim(word string, another string) float64 {
	if word == "gps" {
		if another == "status" || another == "car" {
			return 0.9999
		}
	}

	if word == "state" {
		if another == "tire" {
			return 0.9001
		}
	}

	if word == "ast" && another == "visitor" {
		return 1
	}

	if word == "no" && another == "type" {
		return 0.8564
	}

	if word == "no" && another == "typo" {
		return 0.0001
	}

	if word == "not" {
		return 0.0021
	}

	if word == "notary" && another == "p" {
		return 0
	}

	if word == "notary" && another == "pe" {
		return 0
	}

	if word == "pe" && another == "notary" {
		return 0
	}

	if word == "notary" || another == "notary" {
		return 0.0023
	}

	return 0
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
