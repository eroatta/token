package expanders

import "strings"

const (
	singleWordGroup = "single-word"
	multiWordGroup  = "multi-word"
)

const (
	prefixType          = "prefix"
	droppedLettersType  = "dropped-letters"
	acronymType         = "acronym"
	wordCombinationType = "word-combination"
)

type pattern struct {
	group     string
	kind      string
	shortForm string
	regex     string
}

type patternBuilder struct {
	pattern
}

func (pb *patternBuilder) kind(kind string) *patternBuilder {
	pb.pattern.kind = kind

	switch kind {
	case prefixType, droppedLettersType:
		pb.pattern.group = singleWordGroup
	case acronymType, wordCombinationType:
		pb.pattern.group = multiWordGroup
	}

	return pb
}

func (pb *patternBuilder) shortForm(sf string) *patternBuilder {
	pb.pattern.shortForm = sf
	return pb
}

func (pb *patternBuilder) build() pattern {
	var regex string
	switch pb.pattern.kind {
	case prefixType:
		regex = buildPrefixRegex(pb.pattern.shortForm)
	case droppedLettersType:
		regex = buildDroppedLettersRegex(pb.pattern.shortForm)
	case acronymType:
		regex = buildAcronymRegex(pb.pattern.shortForm)
	case wordCombinationType:
		regex = buildWordCombinationRegex(pb.pattern.shortForm)
	}
	pb.pattern.regex = regex

	return pb.pattern
}

func buildPrefixRegex(input string) string {
	if len(input) == 0 {
		return ""
	}

	var builder strings.Builder
	// TODO: review regexp starting char
	builder.WriteRune('^')
	if input[0] == 'x' {
		builder.WriteString("e?")
	}
	builder.WriteString(input)
	builder.WriteString("[a-z]+")

	return builder.String()
}

func buildDroppedLettersRegex(input string) string {
	if len(input) == 0 {
		return ""
	}

	var builder strings.Builder
	// TODO: review regexp starting char
	builder.WriteRune('^')
	if input[0] == 'x' {
		builder.WriteString("e?")
	}

	for _, letter := range input {
		builder.WriteRune(letter)
		builder.WriteString("[a-z]*")
	}

	return builder.String()
}

func buildAcronymRegex(input string) string {
	if len(input) == 0 {
		return ""
	}

	var builder strings.Builder
	// TODO: review regexp starting char
	builder.WriteRune('(')
	if input[0] == 'x' {
		builder.WriteString("e?")
	}

	for i := 0; i < len(input); i++ {
		builder.WriteByte(input[i])
		builder.WriteString("[a-z]+")
		if i < len(input)-1 {
			builder.WriteString("[ ]")
		}
	}
	builder.WriteRune(')')

	return builder.String()
}

func buildWordCombinationRegex(input string) string {
	return ""
}
