package expanders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild_OnPatternBuilderForPrefix_ShouldReturnPatternWithRegex(t *testing.T) {
	builder := &patternBuilder{}
	pattern := builder.kind("prefix").shortForm("arg").build()

	assert.Equal(t, pattern.group, singleWordGroup)
	assert.Equal(t, pattern.kind, prefixType)
	assert.Equal(t, pattern.shortForm, "arg")
	assert.Equal(t, pattern.regex, "[arg]*")
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
