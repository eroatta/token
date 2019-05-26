package expanders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild_OnPatternBuilderForPrefix_ShouldReturnPatternWithRegex(t *testing.T) {
	cases := []struct {
		name          string
		shortForm     string
		expectedRegex string
	}{
		{"regular_short_form", "arg", "^arg[a-z]+"},
		{"starting_with_x_short_form", "xt", "^e?xt[a-z]+"},
		{"empty_short_form", "", ""},
	}

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			builder := &patternBuilder{}
			pattern := builder.kind("prefix").shortForm(fixture.shortForm).build()

			assert.Equal(t, pattern.group, singleWordGroup)
			assert.Equal(t, pattern.kind, prefixType)
			assert.Equal(t, pattern.shortForm, fixture.shortForm)
			assert.Equal(t, pattern.regex, fixture.expectedRegex)
		})
	}
}

func TestBuild_OnPatternBuilderForDroppedLetters_ShouldReturnPatternWithRegex(t *testing.T) {
	cases := []struct {
		name          string
		shortForm     string
		expectedRegex string
	}{
		{"regular_short_form", "arg", "^a[a-z]*r[a-z]*g[a-z]*"},
		{"starting_with_x_short_form", "xt", "^e?x[a-z]*t[a-z]*"},
		{"empty_short_form", "", ""},
	}

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			builder := &patternBuilder{}
			pattern := builder.kind("dropped-letters").shortForm(fixture.shortForm).build()

			assert.Equal(t, pattern.group, singleWordGroup)
			assert.Equal(t, pattern.kind, droppedLettersType)
			assert.Equal(t, pattern.shortForm, fixture.shortForm)
			assert.Equal(t, pattern.regex, fixture.expectedRegex)
		})
	}
}

func TestBuild_OnPatternBuilderForAcronym_ShouldReturnPatternWithRegex(t *testing.T) {
	cases := []struct {
		name          string
		shortForm     string
		expectedRegex string
	}{
		{"regular_short_form", "json", "(j[a-z]+[ ]s[a-z]+[ ]o[a-z]+[ ]n[a-z]+)"},
		{"starting_with_x_short_form", "xml", "(e?x[a-z]+[ ]m[a-z]+[ ]l[a-z]+)"},
		{"empty_short_form", "", ""},
	}

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			builder := &patternBuilder{}
			pattern := builder.kind("acronym").shortForm(fixture.shortForm).build()

			assert.Equal(t, pattern.group, multiWordGroup)
			assert.Equal(t, pattern.kind, acronymType)
			assert.Equal(t, pattern.shortForm, fixture.shortForm)
			assert.Equal(t, pattern.regex, fixture.expectedRegex)
		})
	}
}

func TestBuild_OnPatternBuilderForWordCombination_ShouldReturnPatternWithRegex(t *testing.T) {
	cases := []struct {
		name          string
		shortForm     string
		expectedRegex string
	}{
		{"regular_short_form", "pdef", "^p[a-z]*?[ ]*?d[a-z]*?[ ]*?e[a-z]*?[ ]*?f[a-z]*?[ ]*?$"},
		{"starting_with_x_short_form", "xp", "^e?x[a-z]*?[ ]*?p[a-z]*?[ ]*?$"},
		{"empty_short_form", "", ""},
	}

	for _, fixture := range cases {
		t.Run(fixture.name, func(t *testing.T) {
			builder := &patternBuilder{}
			pattern := builder.kind("word-combination").shortForm(fixture.shortForm).build()

			assert.Equal(t, pattern.group, multiWordGroup)
			assert.Equal(t, pattern.kind, wordCombinationType)
			assert.Equal(t, pattern.shortForm, fixture.shortForm)
			assert.Equal(t, pattern.regex, fixture.expectedRegex)
		})
	}
}
