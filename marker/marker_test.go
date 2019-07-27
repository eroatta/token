package marker

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnDigits_ShouldAddMarkersBeforeAndAfterNumbers(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{"before_numbers", "brooklyn99", "brooklyn_99"},
		{"after_numbers", "99brooklyn", "99_brooklyn"},
		{"before_and_after", "leto2nd", "leto_2_nd"},
		{"no_markers_for_all_numbers", "99", "99"},
		{"no_markers_for_empty_string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OnDigits(tt.token)

			assert.Equal(t, tt.want, got, fmt.Sprintf("got: %v", got))
		})
	}
}

func TestOnLowerToUpperCase_ShouldAddMarkersOnChangesFromLowerToUpperCase(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{"one_variation", "squarePants", "square_Pants"},
		{"multiple_variations", "squarePantsBob", "square_Pants_Bob"},
		{"one_variation_multiple_uppercase_letters", "responseHTTP", "response_HTTP"},
		{"no_marker_at_the_beginning", "HTTPresponse", "HTTPresponse"},
		{"no_markers_for_empty_string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OnLowerToUpperCase(tt.token)

			assert.Equal(t, tt.want, got, fmt.Sprintf("got: %v", got))
		})
	}
}

func TestOnUpperToLowerCase_ShouldAddMarkersOnChangesFromUpperToLowerCase(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string
	}{
		{"one_variation", "SQUAREPants", "SQUARE_Pants"},
		{"multiple_variations", "spongeBOBSquarePANTSBob", "spongeBOB_SquarePANTS_Bob"},
		{"one_variation_at_beginning", "HTTP_response", "HTTP_response"},
		{"no_marker_when_only_one_uppercase_letter", "httpResponse", "httpResponse"},
		{"no_marker_at_the_end", "responseHTTP", "responseHTTP"},
		{"no_markers_for_empty_string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := OnUpperToLowerCase(tt.token)

			assert.Equal(t, tt.want, got, fmt.Sprintf("got: %v", got))
		})
	}
}

func TestSplitBy_ShouldSplitMarkedTokenInArray(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  []string
	}{
		{"marked_token", "http_response", []string{"http", "response"}},
		{"marked_token_with_leading_and_trailing_blank_chars", " http_response ", []string{" http", "response "}},
		{"no_markers", "word", []string{"word"}},
		{"empty_token", "", []string{""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitBy(tt.token)

			assert.ElementsMatch(t, tt.want, got, fmt.Sprintf("found elements: %v", got))
		})
	}
}
