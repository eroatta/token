package splitters

import (
	"math"
	"regexp"
)

type frequencyTable map[string]float32

var defaultLocalFreqTable frequencyTable
var defaultGlobalFreqTable frequencyTable

// Samurai represents the Samurai splitting algorithm, proposed by Hill et all.
type Samurai struct {
	Splitter

	localFreqTable  *frequencyTable
	globalFreqTable *frequencyTable
}

// NewSamurai creates a new Samurai splitter with the provided frequency tables. If no frequency
// tables are provided, the default tables are used.
func NewSamurai(localFreqTable *frequencyTable, globalFreqTable *frequencyTable) *Samurai {
	local := &defaultLocalFreqTable
	if localFreqTable != nil {
		local = localFreqTable
	}

	global := &defaultGlobalFreqTable
	if globalFreqTable != nil {
		global = globalFreqTable
	}

	return &Samurai{
		localFreqTable:  local,
		globalFreqTable: global,
	}
}

// Split on Samurai receives a token and returns an array of hard/soft words,
// split by the Samurai algorithm proposed by Hill et all.
func (s *Samurai) Split(token string) ([]string, error) {
	preprocessedToken := addMarkersOnDigits(token)
	preprocessedToken = addMarkersOnLowerToUpperCase(preprocessedToken)

	cutLocationRegex := regexp.MustCompile("[A-Z][a-z]")

	var processedToken string
	for _, word := range splitOnMarkers(preprocessedToken) {
		cutLocation := cutLocationRegex.FindStringIndex(word)
		if len(word) > 1 && cutLocation != nil {
			n := len(word) - 1
			i := cutLocation[0]

			var camelScore float64
			if i > 0 {
				camelScore = s.score(word[i:n])
			} else {
				camelScore = s.score(word[0:n])
			}

			altCamelScore := s.score(word[i+1 : n])
			if camelScore > math.Sqrt(altCamelScore) {
				if i > 0 {
					word = word[0:i-1] + "_" + word[i:n]
				}
			} else {
				word = word[0:i] + "_" + word[i+1:n]
			}
		}

		processedToken = processedToken + "_" + word
	}

	splitToken := make([]string, 0, 10)
	for _, word := range splitOnMarkers(preprocessedToken) {
		sameCaseSplitting := s.sameCaseSplit(word, s.score(word))
		splitToken = append(splitToken, sameCaseSplitting...)
	}

	return splitToken, nil
}

func (s *Samurai) sameCaseSplit(token string, baseScore float64) []string {
	maxScore := -1.0

	n := len(token)
	for i := 0; i < n-1; i++ {
		scoreLeft := s.score(token[0:i])
		shouldSplitLeft := math.Sqrt(scoreLeft) > math.Max(s.score(token), baseScore)

		scoreRight := s.score(token[i+1 : n])
		shouldSplitRight := math.Sqrt(scoreRight) > math.Max(s.score(token), baseScore)

		isPreffixOrSuffix := isPreffix(token[0:i]) || isSuffix(token[i+1:n])
		if !isPreffixOrSuffix && shouldSplitLeft && shouldSplitRight {
			if (scoreLeft + scoreRight) > maxScore {
				maxScore = scoreLeft + scoreRight
			}
		} else if !isPreffixOrSuffix && shouldSplitLeft {

		}
	}

	return []string{}
}

func (s *Samurai) score(word string) float64 {
	return 0.0
}

func isPreffix(word string) bool {
	return false
}

func isSuffix(word string) bool {
	return false
}
