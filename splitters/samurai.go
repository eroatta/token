package splitters

import (
	"math"
	"regexp"
)

type frequencyTable map[string]float64

func (t frequencyTable) getFrequency(word string) float64 {
	return map[string]float64(t)[word]
}

type set map[string]bool

func (s set) found(word string) bool {
	return s[word]
}

var defaultLocalFreqTable frequencyTable
var defaultGlobalFreqTable frequencyTable
var defaultPrefixes set
var defaultSuffixes set

// Samurai represents the Samurai splitting algorithm, proposed by Hill et all.
type Samurai struct {
	Splitter

	localFreqTable  *frequencyTable
	globalFreqTable *frequencyTable
	allStringsFreq  float64
	prefixes        *set
	suffixes        *set
}

// NewSamurai creates a new Samurai splitter with the provided frequency tables. If no frequency
// tables are provided, the default tables are used.
func NewSamurai(localFreqTable *frequencyTable, globalFreqTable *frequencyTable, prefixes *set, suffixes *set) *Samurai {
	local := &defaultLocalFreqTable
	if localFreqTable != nil {
		local = localFreqTable
	}

	global := &defaultGlobalFreqTable
	if globalFreqTable != nil {
		global = globalFreqTable
	}

	commonPrefixes := &defaultPrefixes
	if prefixes != nil {
		commonPrefixes = prefixes
	}

	commonSuffixes := &defaultSuffixes
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
	n := len(token) - 1

	for i := 0; i < n; i++ {
		scoreLeft := s.score(token[0:i])
		shouldSplitLeft := math.Sqrt(scoreLeft) > math.Max(s.score(token), baseScore)

		scoreRight := s.score(token[i+1 : n])
		shouldSplitRight := math.Sqrt(scoreRight) > math.Max(s.score(token), baseScore)

		isPreffixOrSuffix := s.isPrefix(token[0:i]) || s.isSuffix(token[i+1:n])
		if !isPreffixOrSuffix && shouldSplitLeft && shouldSplitRight {
			if (scoreLeft + scoreRight) > maxScore {
				maxScore = scoreLeft + scoreRight
				splitToken = []string{token[0:i], token[i+1 : n]}
			}
		} else if !isPreffixOrSuffix && shouldSplitLeft {
			temp := s.sameCaseSplit(token[i+1:n], baseScore)
			if len(temp) > 1 {
				splitToken = []string{token[0:i]}
				splitToken = append(splitToken, temp...)
			}
		}
	}

	return splitToken
}

// score calculates the a score for a string based on how frequently a word
// appears in the program under analysis and in a more global scope of a large
// set of programs.
func (s *Samurai) score(word string) float64 {
	// Freq(s,p) + (globalFreq(s) / log_10 (AllStrsFreq(p))
	return s.localFreqTable.getFrequency(word) + s.globalFreqTable.getFrequency(word)/math.Log10(s.allStringsFreq)
}

// isPrefix checks if the current token is found on a list of common prefixes.
func (s *Samurai) isPrefix(token string) bool {
	return s.prefixes.found(token)
}

// isSuffix checks if the current token is found on a list of common suffixes.
func (s *Samurai) isSuffix(token string) bool {
	return s.suffixes.found(token)
}
