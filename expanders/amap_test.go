package expanders

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAmap_ShouldReturnNewAmapInstance(t *testing.T) {
	assert.Fail(t, "not yet implemented")
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

func TestSingleWordExpansion_OnAmapWithNoMatches_ShouldReturnEmptyLongFormCandidates(t *testing.T) {
	variableDeclarations := []string{"cpol Carpool"}
	methodName := "buildCarpool"
	var emptyMethodBody, emptyMethodComments, emptyPackageComments string

	amap := NewAmap()
	pattern := (&patternBuilder{}).kind("prefix").shortForm("cp").build()
	got := amap.singleWordExpansion(pattern, variableDeclarations, methodName, emptyMethodBody,
		emptyMethodComments, emptyPackageComments)

	assert.Empty(t, got)
}

func TestSingleWordExpansion_OnAmapWithPrefixPatternButManyWovels_ShouldReturnEmptyLongForms(t *testing.T) {
	possibleButSkippedMatch := []string{"iooboooo ioob"}
	var emptyMethodName, emptyMethodBody, emptyMethodComments, emptyPackageComments string

	amap := NewAmap()
	pattern := (&patternBuilder{}).kind("prefix").shortForm("ioob").build()
	got := amap.singleWordExpansion(pattern, possibleButSkippedMatch, emptyMethodName, emptyMethodBody,
		emptyMethodComments, emptyPackageComments)

	assert.Empty(t, got)
}

func TestSingleWordExpansion_OnAmapWithPrefixPattern_ShouldReturnMatchingLongFormCandidates(t *testing.T) {
	cases := []struct {
		name       string
		shortForm  string
		candidates []string
	}{
		{"many_vowels", "cp", []string{}},
		{"not_many_vowels_and_only_match_on_variable_declarations", "carp", []string{"carpool"}},
		{"not_many_vowels_and_only_match_on_method_name", "bu", []string{"buildCarpool"}},
		{"not_many_vowels_and_short_form_size_not_two_one_match_method_body", "syn", []string{"syntax"}},
		{"not_many_vowels_and_short_form_size_not_two_one_match_method_comments", "abs", []string{"abstract"}},
		{"not_many_vowels_and_short_form_size_two_no_match_method_body", "sy", []string{}},
		{"not_many_vowels_and_short_form_size_two_no_match_method_comments", "ab", []string{}},
		{"not_many_vowels_and_short_form_size_higher_than_one_one_match_package_comments", "wal", []string{"walker"}},
	}

	amap := NewAmap()
	variableDeclarations := []string{"carpool carp"}
	methodName := "buildCarpool"
	methodBodyText := "syntax analizer"
	methodComments := "abstract watcher"
	packageComments := "walker"

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			pattern := (&patternBuilder{}).kind("prefix").shortForm(fixture.shortForm).build()

			got := amap.singleWordExpansion(pattern, variableDeclarations, methodName, methodBodyText,
				methodComments, packageComments)

			assert.ElementsMatch(t, fixture.candidates, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}
