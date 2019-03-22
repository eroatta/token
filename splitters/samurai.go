package splitters

import (
	"math"
	"regexp"

	"github.com/eroatta/token-splitex/lists"

	"github.com/deckarep/golang-set"
)

var defaultLocalFreqTable FrequencyTable
var defaultGlobalFreqTable FrequencyTable

// Samurai represents the Samurai splitting algorithm, proposed by Hill et all.
type Samurai struct {
	localFreqTable  *FrequencyTable
	globalFreqTable *FrequencyTable
	prefixes        mapset.Set
	suffixes        mapset.Set
}

// NewSamurai creates a new Samurai splitter with the provided frequency tables. If no frequency
// tables are provided, the default tables are used.
func NewSamurai(localFreqTable *FrequencyTable, globalFreqTable *FrequencyTable, prefixes mapset.Set, suffixes mapset.Set) *Samurai {
	local := &defaultLocalFreqTable
	if localFreqTable != nil {
		local = localFreqTable
	}

	global := &defaultGlobalFreqTable
	if globalFreqTable != nil {
		global = globalFreqTable
	}

	commonPrefixes := lists.Prefixes
	if prefixes != nil {
		commonPrefixes = prefixes
	}

	commonSuffixes := lists.Suffixes
	if suffixes != nil {
		commonSuffixes = suffixes
	}

	return &Samurai{
		localFreqTable:  local,
		globalFreqTable: global,
		prefixes:        commonPrefixes,
		suffixes:        commonSuffixes,
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

	splitToken := []string{token}
	n := len(token)

	for i := 0; i < n; i++ {
		left := token[0:i]
		scoreLeft := s.score(left)
		shouldSplitLeft := math.Sqrt(scoreLeft) > math.Max(s.score(token), baseScore)

		right := token[i:n]
		scoreRight := s.score(right)
		shouldSplitRight := math.Sqrt(scoreRight) > math.Max(s.score(token), baseScore)

		isPreffixOrSuffix := s.isPrefix(left) || s.isSuffix(right)
		if !isPreffixOrSuffix && shouldSplitLeft && shouldSplitRight {
			if (scoreLeft + scoreRight) > maxScore {
				maxScore = scoreLeft + scoreRight
				splitToken = []string{left, right}
			}
		} else if !isPreffixOrSuffix && shouldSplitLeft {
			temp := s.sameCaseSplit(right, baseScore)
			if len(temp) > 1 {
				splitToken = []string{left}
				splitToken = append(splitToken, temp...)
			}
		}
	}

	return splitToken
}

// score calculates the score for a string based on how frequently a word
// appears in the program under analysis and in a more global scope of a large set of programs.
func (s *Samurai) score(word string) float64 {
	freqS := s.localFreqTable.Frequency(word)
	globalFreqS := s.globalFreqTable.Frequency(word)
	allStrsFreqP := float64(s.localFreqTable.TotalOccurrences())

	// Freq(s,p) + (globalFreq(s) / log_10 (AllStrsFreq(p))
	return freqS + globalFreqS/math.Log10(allStrsFreqP)
}

// isPrefix checks if the current token is found on a list of common prefixes.
func (s *Samurai) isPrefix(token string) bool {
	return s.prefixes.Contains(token)
}

// isSuffix checks if the current token is found on a list of common suffixes.
func (s *Samurai) isSuffix(token string) bool {
	return s.suffixes.Contains(token)
}
