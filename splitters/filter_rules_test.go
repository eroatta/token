package splitters

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

func TestIsTruncation_OnAbbrEqualThanWord_ShouldReturnTrue(t *testing.T) {
	got := isTruncation("car", "car")

	assert.True(t, got)
}
