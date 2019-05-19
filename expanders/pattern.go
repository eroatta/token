package expanders

const (
	prefixType          = "prefix"
	droppedLettersType  = "dropped-letters"
	acronymType         = "acronym"
	wordCombinationType = "word-combination"
)

type pattern struct {
	typ   string
	regex string
}

func buildPrefixPattern(shortForm string) pattern {
	return pattern{
		typ: prefixType,
	}
}

func buildDroppedLettersPattern(shortForm string) pattern {
	return pattern{
		typ: droppedLettersType,
	}
}

func buildAcronymPattern(shortForm string) pattern {
	return pattern{
		typ: acronymType,
	}
}

func buildWordCombinationPattern(shortForm string) pattern {
	return pattern{
		typ: wordCombinationType,
	}
}
