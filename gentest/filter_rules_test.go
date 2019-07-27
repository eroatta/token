package gentest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTruncation_OnEmptyAbbr_ShouldReturnFalse(t *testing.T) {
	got := isTruncation("", "word")

	assert.False(t, got)
}

func TestIsTruncation_OnEmptyWord_ShouldReturnFalse(t *testing.T) {
	got := isTruncation("car", "")

	assert.False(t, got)
}

func TestIsTruncation_OnLongerAbbrThanWord_ShouldReturnFalse(t *testing.T) {
	got := isTruncation("carpool", "car")

	assert.False(t, got)
}

func TestIsTruncation_OnNonAbbrOfWord_ShouldReturnFalse(t *testing.T) {
	got := isTruncation("car", "dog")

	assert.False(t, got)
}

func TestIsTruncation_OnAbbrOfWord_ShouldReturnTrue(t *testing.T) {
	got := isTruncation("car", "carpool")

	assert.True(t, got)
}

func TestIsTruncation_OnAbbrEqualsThanWord_ShouldReturnTrue(t *testing.T) {
	got := isTruncation("car", "car")

	assert.True(t, got)
}

func TestIsTruncation_OnSinlgeLetterAbbrAndPluralWord_ShouldReturnFalse(t *testing.T) {
	got := isTruncation("s", "arguments")

	assert.False(t, got)
}

func TestIsTruncation_OnPluralAbbr_ShouldReturnTrue(t *testing.T) {
	got := isTruncation("args", "arguments")

	assert.True(t, got)
}

func TestHasRemovedChar_OnEmptyAbbr_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedChar("", "car")

	assert.False(t, got)
}

func TestHasRemovedChar_OnEmptyWord_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedChar("cr", "")

	assert.False(t, got)
}

func TestHasRemovedChar_OnAbbrWithSameSizeThanWord_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedChar("sts", "set")

	assert.False(t, got)
}

func TestHasRemovedChar_OnAbbrThanHasDifferentStartCharThanWord_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedChar("st", "ast")

	assert.False(t, got)
}

func TestHasRemovedChar_OnNonAbbrOfWord_ShouldReturnFalse(t *testing.T) {
	cases := []struct {
		abbr string
		word string
	}{
		{"st", "sub"},
		{"snit", "saint"},
	}
	for _, c := range cases {
		got := hasRemovedChar(c.abbr, c.word)

		assert.False(t, got)
	}
}

func TestHasRemovedChar_OnAbbrOfWord_ShouldReturnTrue(t *testing.T) {
	got := hasRemovedChar("st", "set")

	assert.True(t, got)
}

func TestHasRemovedVowels_OnEmptyAbbr_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedVowels("", "o")

	assert.False(t, got)
}

func TestHasRemovedVowels_OnEmptyWord_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedVowels("cr", "")

	assert.False(t, got)
}

func TestHasRemovedVowels_OnAbbrWithVowels_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedVowels("car", "car")

	assert.False(t, got)
}

func TestHasRemovedVowels_OnNonAbbrOfWord_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedVowels("st", "street")

	assert.False(t, got)
}

func TestHasRemovedVowels_OnAbbrOfWord_ShouldReturnTrue(t *testing.T) {
	got := hasRemovedVowels("st", "set")

	assert.True(t, got)
}

func TestHasRemoverCharAfterRemovedVowels_OnEmptyAbbr_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedCharAfterRemovedVowels("", "car")

	assert.False(t, got)
}

func TestHasRemoverCharAfterRemovedVowels_OnEmptyWord_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedCharAfterRemovedVowels("cr", "")

	assert.False(t, got)
}

func TestHasRemoverCharAfterRemovedVowels_OnNonAbbrOfWord_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedCharAfterRemovedVowels("str", "stood")

	assert.False(t, got)
}

func TestHasRemoverCharAfterRemovedVowels_OnAbbrOfWordButMoreCharsRemoved_ShouldReturnFalse(t *testing.T) {
	got := hasRemovedCharAfterRemovedVowels("str", "string")

	assert.False(t, got)
}

func TestHasRemoverCharAfterRemovedVowels_OnAbbrOfWord_ShouldReturnTrue(t *testing.T) {
	got := hasRemovedCharAfterRemovedVowels("strg", "string")

	assert.True(t, got)
}
