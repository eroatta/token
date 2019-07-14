package amap

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

var (
	groups = map[string]string{
		prefixType:          singleWordGroup,
		droppedLettersType:  singleWordGroup,
		acronymType:         multiWordGroup,
		wordCombinationType: multiWordGroup,
	}

	regexBuilders = map[string]func(string) string{
		prefixType:          buildPrefixRegex,
		droppedLettersType:  buildDroppedLettersRegex,
		acronymType:         buildAcronymRegex,
		wordCombinationType: buildWordCombinationRegex,
	}
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
	pb.pattern.group = groups[kind]

	return pb
}

func (pb *patternBuilder) shortForm(sf string) *patternBuilder {
	pb.pattern.shortForm = sf
	return pb
}

func (pb *patternBuilder) build() pattern {
	pb.pattern.regex = regexBuilders[pb.pattern.kind](pb.pattern.shortForm)
	return pb.pattern
}

func buildPrefixRegex(input string) string {
	if len(input) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("\\b")
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
	builder.WriteString("\\b")
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
	if len(input) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("\\b")
	if input[0] == 'x' {
		builder.WriteString("e?")
	}

	for _, letter := range input {
		builder.WriteRune(letter)
		builder.WriteString("[a-z]*?[ ]*?")
	}
	builder.WriteString("\\b")

	return builder.String()
}
