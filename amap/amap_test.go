package amap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTokenScope_ShouldReturnNewTokenScope(t *testing.T) {
	scope := NewTokenScope([]string{"vardecl"}, "methodName", "bodyText",
		[]string{"method comments"}, []string{"package comments"})

	assert.ElementsMatch(t, []string{"vardecl"}, scope.variableDeclarations)
	assert.Equal(t, "methodName", scope.methodName)
	assert.Equal(t, "bodyText", scope.methodBodyText)
	assert.ElementsMatch(t, []string{"method comments"}, scope.methodComments)
	assert.ElementsMatch(t, []string{"package comments"}, scope.packageComments)
}

func TestExpand_OnAmap_ShouldReturnExpansion(t *testing.T) {
	cases := []struct {
		name     string
		token    string
		expected []string
	}{
		{"long_form_found_using_prefix_pattern", "exP", []string{"expansion"}},
		{"long_form_found_using_dropped_letter_pattern", "expnsn", []string{"expansion"}},
		{"long_form_found_using_acronym_pattern", "GUI", []string{"graphical user interface"}},
		{"long_form_found_using_word_combination_pattern", "cdrdr", []string{"card reader"}},
		{"long_form_found_handling_multiple_matches", "int", []string{"interface"}},
		{"short_from_skipped_during_validation", "ex", []string{}},
	}

	methodBodyText := "expansion interface interfacing interface interfaces"
	methodComments := []string{"providing graphical user interface for Linux, setting a card reader implementation"}
	packageComments := []string{"provides a card reader implementation"}
	scope := NewTokenScope([]string{}, "", methodBodyText, methodComments, packageComments)

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			got := Expand(fixture.token, scope, []string{})

			assert.ElementsMatch(t, fixture.expected, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestSingleWordExpansion_OnAmapWithNoMatches_ShouldReturnEmptyLongForms(t *testing.T) {
	pattern := (&patternBuilder{}).kind("prefix").shortForm("cp").build()

	variableDeclarations := []string{"cpol Carpool"}
	emptyComments := []string{}
	scope := NewTokenScope(variableDeclarations, "buildCarpool", "", emptyComments, emptyComments)

	got := searchSingleWordExpansion(pattern, scope)

	assert.Empty(t, got)
}

func TestSingleWordExpansion_OnAmapWithPrefixPatternButTooManyWovels_ShouldReturnEmptyLongForms(t *testing.T) {
	pattern := (&patternBuilder{}).kind("prefix").shortForm("ioob").build()

	possibleButSkippedMatch := []string{"iooboooo ioob"}
	emptyComments := []string{}
	scope := NewTokenScope(possibleButSkippedMatch, "", "", emptyComments, emptyComments)

	got := searchSingleWordExpansion(pattern, scope)

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

	variableDeclarations := []string{"carpool carp"}
	methodName := "buildCarpool"
	methodBodyText := "syntax analizer"
	methodComments := []string{"abstract watcher"}
	packageComments := []string{"walker"}
	scope := NewTokenScope(variableDeclarations, methodName, methodBodyText, methodComments, packageComments)

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			pattern := (&patternBuilder{}).kind("prefix").shortForm(fixture.shortForm).build()

			got := searchSingleWordExpansion(pattern, scope)

			assert.ElementsMatch(t, fixture.candidates, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestSingleWordExpansion_OnAmapWithNoPrefixPatternButTooManyVowels_ShouldReturnEmptyLongForms(t *testing.T) {
	pattern := (&patternBuilder{}).kind("dropped-letters").shortForm("cool").build()

	possibleButSkippedMatch := []string{"contextpool coo"}
	emptyComments := []string{}
	scope := NewTokenScope(possibleButSkippedMatch, "", "", emptyComments, emptyComments)

	got := searchSingleWordExpansion(pattern, scope)

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

	variableDeclarations := []string{"carpool cpl"}
	methodName := "buildcarpool"
	methodBodyText := "syntax analizer"
	methodComments := []string{"abstract watcher"}
	packageComments := []string{"walker"}
	scope := NewTokenScope(variableDeclarations, methodName, methodBodyText, methodComments, packageComments)

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			pattern := (&patternBuilder{}).kind("dropped-letters").shortForm(fixture.shortForm).build()

			got := searchSingleWordExpansion(pattern, scope)

			assert.ElementsMatch(t, fixture.expected, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestSingleWordExpansion_OnAmapWithPrefixPatternButNoSingleMatch_ShouldReturnMatchingLongForms(t *testing.T) {
	pattern := (&patternBuilder{}).kind("prefix").shortForm("abs").build()

	methodComments := []string{"abstraction for abstract syntax tree"}
	packageComments := []string{"absolute"}
	scope := NewTokenScope([]string{}, "", "", methodComments, packageComments)

	got := searchSingleWordExpansion(pattern, scope)

	assert.ElementsMatch(t, []string{"abstraction", "abstract", "absolute"}, got, fmt.Sprintf("found elements: %v", got))
}

func TestMultiWordExpansion_OnAmapWithNoMatches_ShouldReturnEmptyLongForms(t *testing.T) {
	pattern := (&patternBuilder{}).kind("acronym").shortForm("json").build()

	variableDeclarations := []string{"cpol Carpool"}
	emptyComments := []string{}
	scope := NewTokenScope(variableDeclarations, "buildCarpool", "", emptyComments, emptyComments)

	got := searchMultiWordExpansion(pattern, scope)

	assert.Empty(t, got)
}

func TestMultiWordExpansion_OnAmapWithNotAcronymPatternButTooShortShortForm_ShouldReturnEmptyLongForms(t *testing.T) {
	pattern := (&patternBuilder{}).kind("word-combination").shortForm("jsn").build()

	var emptyVariableDecls, emptyComments []string
	possibleButSkippedMethodComments := []string{"javascript notation"}
	scope := NewTokenScope(emptyVariableDecls, "", "", possibleButSkippedMethodComments, emptyComments)

	got := searchMultiWordExpansion(pattern, scope)

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
		{"match_acronym_on_method_name", "acronym", "jpb", []string{"json parser builder"}},
		{"match_acronym_on_method_body_text", "acronym", "ff", []string{"factory function"}},
		{"match_acronym_on_method_comments", "acronym", "xml", []string{"extensible markup language"}},
		{"match_acronym_on_package_comments", "acronym", "ftp", []string{"file transfer protocol"}},
		{"match_word-combination_on_variable_decl", "word-combination", "jsonpf", []string{"json parser factory"}},
		{"match_word-combination_on_method_name", "word-combination", "jpbder", []string{"json parser builder"}},
		{"match_word-combination_on_method_body_text", "word-combination", "facfunc", []string{"factory function"}},
		{"match_word-combination_on_method_comments", "word-combination", "xmlang", []string{"extensible markup language"}},
		{"no_match_word-combination_skipped_package_comments", "word-combination", "filetp", []string{}},
	}

	variableDeclarations := []string{"json parser factory jpf", "json parser factory jsonpf"}
	methodName := "json parser builder"
	methodBodyText := "factory function"
	methodComments := []string{"extensible markup language"}
	packageComments := []string{"file transfer protocol enables file transferring between"}
	scope := NewTokenScope(variableDeclarations, methodName, methodBodyText, methodComments, packageComments)

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			pattern := (&patternBuilder{}).kind(fixture.patternType).shortForm(fixture.shortForm).build()

			got := searchMultiWordExpansion(pattern, scope)

			assert.ElementsMatch(t, fixture.expected, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}

func TestMultiWordExpansion_OnAmapWithAcronymPatternButNoSingleMatch_ShouldReturnMatchingLongForms(t *testing.T) {
	pattern := (&patternBuilder{}).kind("acronym").shortForm("js").build()

	methodComments := []string{"java script define json source"}
	packageComments := []string{"java script"}
	scope := NewTokenScope([]string{}, "", "", methodComments, packageComments)

	got := searchMultiWordExpansion(pattern, scope)

	assert.ElementsMatch(t, []string{"java script", "json source", "java script"}, got, fmt.Sprintf("found elements: %v", got))
}

func TestMultiWordExpansion_OnAmapWithWordCombinationPatternButNoSingleMatch_ShouldReturnMatchingLongForms(t *testing.T) {
	pattern := (&patternBuilder{}).kind("word-combination").shortForm("jspt").build()

	methodBodyText := "java script is not a java scripting language"
	methodComments := []string{"java script define json source"}
	scope := NewTokenScope([]string{}, "", methodBodyText, methodComments, []string{})

	got := searchMultiWordExpansion(pattern, scope)

	assert.ElementsMatch(t, []string{"java script", "java scripting", "java script"}, got, fmt.Sprintf("found elements: %v", got))
}

func TestMostFrequentWord_WhenNoWords_ShouldReturnEmptyWord(t *testing.T) {
	words := []string{}

	mfw := mostFrequentWord(words)

	assert.Equal(t, "", mfw)
}

func TestMostFrequentWord_WhenTwoRepeatedWords_ShouldReturnWord(t *testing.T) {
	words := []string{"valid", "valid"}

	mfw := mostFrequentWord(words)

	assert.Equal(t, "valid", mfw)
}

func TestMostFrequentWord_WhenTwoRepeatedWordsAndAnotherWord_ShouldReturnRepeatedWord(t *testing.T) {
	words := []string{"valid", "validate", "valid"}

	mfw := mostFrequentWord(words)

	assert.Equal(t, "valid", mfw)
}

func TestMostFrequentWord_WhenSeveralRepeatedWords_ShouldReturnMostFrequentWord(t *testing.T) {
	words := []string{"valid", "validate", "valid", "valid", "validation", "validation"}

	mfw := mostFrequentWord(words)

	assert.Equal(t, "valid", mfw)
}

func TestMostFrequentWord_WhenSeveralRepeatedWordsWithNoWinningWord_ShouldReturnEmptyWord(t *testing.T) {
	words := []string{"valid", "validate", "validation", "valid", "valid", "validation", "validation"}

	mfw := mostFrequentWord(words)

	assert.Equal(t, "", mfw)
}

func TestStemmedWords_WhenEmptyArray_ShouldReturnEmptyArray(t *testing.T) {
	words := []string{}

	stemmedWords := stemmedWords(words)

	assert.Empty(t, stemmedWords)
}

func TestStemmedWords_WhenThreeWords_ShouldReturnThreeStemmedWords(t *testing.T) {
	words := []string{"default", "defaults", "running"}

	stemmedWords := stemmedWords(words)

	assert.Equal(t, 3, len(stemmedWords))
	assert.ElementsMatch(t, []string{"default", "default", "run"}, stemmedWords)
}

func TestStemmedWord_WhenCombinedWords_ShouldReturnEachWordStemmed(t *testing.T) {
	words := []string{"graphical user interface"}

	stemmedWords := stemmedWords(words)

	assert.Equal(t, 1, len(stemmedWords))
	assert.ElementsMatch(t, []string{"graphic user interfac"}, stemmedWords)
}

func TestMostFrequentExpansion_WhenNoMatches_ShouldReturnEmptyWord(t *testing.T) {
	pttrn := (&patternBuilder{}).kind("prefix").shortForm("val").build()
	text := []string{"no matching expansion"}

	mfe := mostFrequentExpansion(pttrn, text)

	assert.Equal(t, "", mfe)
}

func TestMostFrequentExpansion_WhenRelativeFrequencyLessThanZeroFive_ShouldReturnEmptyWord(t *testing.T) {
	pttrn := (&patternBuilder{}).kind("prefix").shortForm("val").build()
	text := []string{
		"big value",
		"medium value",
		"small value",
		"very small value",
		"check validation",
		"review validation",
		"check validity",
		"review validity",
	}

	mfe := mostFrequentExpansion(pttrn, text)

	assert.Equal(t, "", mfe)
}

func TestMostFrequentExpansion_WhenLessThanThreeMatches_ShouldReturnEmptyWord(t *testing.T) {
	pttrn := (&patternBuilder{}).kind("prefix").shortForm("val").build()
	text := []string{
		"big value",
		"another value",
	}

	mfe := mostFrequentExpansion(pttrn, text)

	assert.Equal(t, "", mfe)
}

func TestMostFrequentExpansion_WhenExpansionFound_ShouldReturnExpansion(t *testing.T) {
	pttrn := (&patternBuilder{}).kind("prefix").shortForm("val").build()
	text := []string{
		"big value",
		"medium value",
		"small value",
		"very small value",
		"tiny value",
		"check validation",
		"review validation",
		"check validity",
	}

	mfe := mostFrequentExpansion(pttrn, text)

	assert.Equal(t, "value", mfe)
}
