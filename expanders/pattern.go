package expanders

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
	}
	pb.pattern.regex = regex

	return pb.pattern
}

func buildPrefixRegex(input string) string {
	return ""
}

func buildDroppedLettersPattern(shortForm string) pattern {
	return pattern{
		kind: droppedLettersType,
	}
}

func buildAcronymPattern(shortForm string) pattern {
	return pattern{
		kind: acronymType,
	}
}

func buildWordCombinationPattern(shortForm string) pattern {
	return pattern{
		kind: wordCombinationType,
	}
}
