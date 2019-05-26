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
		{"special_x_case", "xt", "^e?xt[a-z]+"},
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
	assert.FailNow(t, "not yet implemented or properly tested")
}

func TestBuild_OnPatternBuilderForAcronym_ShouldReturnPatternWithRegex(t *testing.T) {
	assert.FailNow(t, "not yet implemented or properly tested")
}

func TestBuild_OnPatternBuilderForWordCombination_ShouldReturnPatternWithRegex(t *testing.T) {
	assert.FailNow(t, "not yet implemented or properly tested")
}
