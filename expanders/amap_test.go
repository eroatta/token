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
		// TODO: review
		{"no_expansion", "noExpansion", []string{}},
	}

	amap := NewAmap()
	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			got := amap.Expand(fixture.token)

			assert.ElementsMatch(t, fixture.expected, got, fmt.Sprintf("found elements: %v", got))
			assert.Fail(t, "not yet implemented")
		})
	}
}

func TestSingleWordExpansion_OnAmapWithNoMatches_ShouldReturnEmptyLongForms(t *testing.T) {
	variableDeclarations := []string{"cpol Carpool"}
	methodName := "buildCarpool"
	var emptyMethodBody, emptyMethodComments, emptyPackageComments string

	amap := NewAmap()
	pattern := (&patternBuilder{}).kind("prefix").shortForm("cp").build()
	got := amap.singleWordExpansion(pattern, variableDeclarations, methodName, emptyMethodBody,
		emptyMethodComments, emptyPackageComments)

	assert.Empty(t, got)
}

func TestSingleWordExpansion_OnAmapWithPrefixPatternButTooManyWovels_ShouldReturnEmptyLongForms(t *testing.T) {
	possibleButSkippedMatch := []string{"iooboooo ioob"}
	var emptyMethodName, emptyMethodBody, emptyMethodComments, emptyPackageComments string

	amap := NewAmap()
	pattern := (&patternBuilder{}).kind("prefix").shortForm("ioob").build()
	got := amap.singleWordExpansion(pattern, possibleButSkippedMatch, emptyMethodName, emptyMethodBody,
		emptyMethodComments, emptyPackageComments)

	assert.Empty(t, got)
}

func TestSingleWordExpansion_OnAmapWithPrefixPatternAndNotManyVowels_ShouldReturnMatchingLongForms(t *testing.T) {
	cases := []struct {
		name       string
		shortForm  string
		candidates []string
	}{
		{"only_match_on_variable_declarations", "carp", []string{"carpool"}},
		{"only_match_on_method_name", "bu", []string{"buildCarpool"}},
		{"short_form_size_not_two_one_match_method_body", "syn", []string{"syntax"}},
		{"short_form_size_not_two_one_match_method_comments", "abs", []string{"abstract"}},
		{"short_form_size_two_no_match_method_body", "sy", []string{}},
		{"short_form_size_two_no_match_method_comments", "ab", []string{}},
		{"short_form_size_higher_than_one_one_match_package_comments", "wal", []string{"walker"}},
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

func TestSingleWordExpansion_OnAmapWithNoPrefixPatternButTooManyVowels_ShouldReturnEmptyLongForms(t *testing.T) {
	possibleButSkippedMatch := []string{"contextpool coo"}
	var emptyMethodName, emptyMethodBody, emptyMethodComments, emptyPackageComments string

	amap := NewAmap()
	pattern := (&patternBuilder{}).kind("dropped-letters").shortForm("cool").build()
	got := amap.singleWordExpansion(pattern, possibleButSkippedMatch, emptyMethodName, emptyMethodBody,
		emptyMethodComments, emptyPackageComments)

	assert.Empty(t, got)
}

func TestSingleWordExpansion_OnAmapWithNotPrefixPatternAndNotManyVowels_ShouldReturnMatchingLongForms(t *testing.T) {
	cases := []struct {
		name      string
		shortForm string
		expected  []string
	}{
		{"only_match_on_variable_declarations", "cpl", []string{"carpool"}},
		{"only_match_on_method_name", "bcpl", []string{"buildcarpool"}},
		{"short_form_size_not_two_one_match_method_body", "syn", []string{"syntax"}},
		{"short_form_size_not_two_one_match_method_comments", "absat", []string{"abstract"}},
		{"short_form_size_two_no_match_method_body", "sy", []string{}},
		{"short_form_size_two_no_match_method_comments", "ab", []string{}},
		{"short_form_size_higher_than_one_but_skipped_package_comments", "wlkr", []string{}},
	}

	amap := NewAmap()
	variableDeclarations := []string{"carpool cpl"}
	methodName := "buildcarpool"
	methodBodyText := "syntax analizer"
	methodComments := "abstract watcher"
	packageComments := "walker"

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			pattern := (&patternBuilder{}).kind("dropped-letters").shortForm(fixture.shortForm).build()

			got := amap.singleWordExpansion(pattern, variableDeclarations, methodName, methodBodyText,
				methodComments, packageComments)

			assert.ElementsMatch(t, fixture.expected, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestSingleWordExpansion_OnAmapWithPrefixPatternButNoSingleMatch_ShouldReturnMatchingLongForms(t *testing.T) {
	var emptyVariableDecls []string
	var emptyMethodName, emptyMethodBody string
	methodComments := "abstraction for abstract syntax tree"
	packageComments := "absolute"

	amap := NewAmap()
	pattern := (&patternBuilder{}).kind("prefix").shortForm("abs").build()
	got := amap.singleWordExpansion(pattern, emptyVariableDecls, emptyMethodName, emptyMethodBody,
		methodComments, packageComments)

	assert.ElementsMatch(t, []string{"abstraction", "abstract", "absolute"}, got, fmt.Sprintf("found elements: %v", got))
}

func TestMultiWordExpansion_OnAmapWithNoMatches_ShouldReturnEmptyLongForms(t *testing.T) {
	variableDeclarations := []string{"cpol Carpool"}
	methodName := "buildCarpool"
	var emptyMethodBody, emptyMethodComments, emptyPackageComments string

	amap := NewAmap()
	pattern := (&patternBuilder{}).kind("acronym").shortForm("json").build()
	got := amap.multiWordExpansion(pattern, variableDeclarations, methodName, emptyMethodBody,
		emptyMethodComments, emptyPackageComments)

	assert.Empty(t, got)
}

func TestMultiWordExpansion_OnAmapWithNotAcronymPatternButTooShortShortForm_ShouldReturnEmptyLongForms(t *testing.T) {
	var emptyVariableDecls []string
	var emptyMethodName, emptyMethodBody, emptyPackageComments string
	possibleButSkippedMethodComments := "javascript notation"

	amap := NewAmap()
	pattern := (&patternBuilder{}).kind("word-combination").shortForm("jsn").build()
	got := amap.multiWordExpansion(pattern, emptyVariableDecls, emptyMethodName, emptyMethodBody,
		possibleButSkippedMethodComments, emptyPackageComments)

	assert.Empty(t, got)
}

func TestMultiWordExpansion_OnAmap_ShouldReturnMatchingLongForms(t *testing.T) {
	cases := []struct {
		name        string
		patternType string
		shortForm   string
		expected    []string
	}{
		{"match_acronym_on_variable_decl", "acronym", "jpf", []string{"json parser factory"}},
		// {"match_acronym_on_method_name", "acronym", "jpb", []string{"jsonParserBuilder"}}, this will never match with the pattern
		{"match_acronym_on_method_body_text", "acronym", "ff", []string{"factory function"}},
		{"match_acronym_on_method_comments", "acronym", "xml", []string{"extensible markup language"}},
		{"match_acronym_on_package_comments", "acronym", "ftp", []string{"file transfer protocol"}},
		//{"match_word-combination_on_variable_decl", "word-combination", "jsonpf", []string{"json parser factory"}},
	}

	variableDeclarations := []string{"json parser factory jpf"}
	methodName := "jsonParserBuilder"
	methodBodyText := "factory function"
	methodComments := "extensible markup language"
	packageComments := "file transfer protocol enables file transfering between"

	amap := NewAmap()
	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			pattern := (&patternBuilder{}).kind(fixture.patternType).shortForm(fixture.shortForm).build()

			got := amap.multiWordExpansion(pattern, variableDeclarations, methodName, methodBodyText,
				methodComments, packageComments)

			assert.ElementsMatch(t, fixture.expected, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestMultiWordExpansion_OnAmapWithPrefixAcronymButNoSingleMatch_ShouldReturnMatchingLongForms(t *testing.T) {
	var emptyVariableDecls []string
	var emptyMethodName, emptyMethodBody string
	methodComments := "java script define json source"
	packageComments := "java script"

	amap := NewAmap()
	pattern := (&patternBuilder{}).kind("acronym").shortForm("js").build()
	got := amap.multiWordExpansion(pattern, emptyVariableDecls, emptyMethodName, emptyMethodBody,
		methodComments, packageComments)

	assert.ElementsMatch(t, []string{"java script", "json source", "java script"}, got, fmt.Sprintf("found elements: %v", got))
}
