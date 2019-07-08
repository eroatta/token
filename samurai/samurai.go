package samurai

import (
	"math"
	"regexp"

	"github.com/eroatta/token-splitex/marker"

	"github.com/eroatta/token-splitex/lists"
)

var cutLocationRegex = regexp.MustCompile("[A-Z][a-z]")

// TokenContext holds the local and global frequency tables in the context of a given token.
type TokenContext struct {
	localFT  *FrequencyTable
	globalFT *FrequencyTable
}

// Score calculates the score for a string based on how frequently a word
// appears in the program under analysis and in a more global scope of a large set of programs.
func (ctx TokenContext) Score(word string) float64 {
	freqS := ctx.localFT.Frequency(word)
	globalFreqS := ctx.globalFT.Frequency(word)
	allStrsFreqP := float64(ctx.localFT.TotalOccurrences())

	// Freq(s,p) + (globalFreq(s) / log_10 (AllStrsFreq(p))
	return freqS + globalFreqS/math.Log10(allStrsFreqP)
}

// Split on Samurai receives a token and returns an array of hard/soft words,
// split by the Samurai algorithm proposed by Hill et all.
func Split(token string, tCtx TokenContext, prefixes lists.List, suffixes lists.List) []string {
	preprocessedToken := marker.OnDigits(token)
	preprocessedToken = marker.OnLowerToUpperCase(token)

	var processedToken string
	for _, word := range marker.SplitBy(preprocessedToken) {
		cutLocation := cutLocationRegex.FindStringIndex(word)
		if len(word) > 1 && cutLocation != nil {
			n := len(word) - 1
			i := cutLocation[0]

			var camelScore float64
			if i > 0 {
				camelScore = tCtx.Score(word[i:n])
			} else {
				camelScore = tCtx.Score(word[0:n])
			}

			altCamelScore := tCtx.Score(word[i+1 : n])
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
	for _, word := range marker.SplitBy(preprocessedToken) {
		sameCaseSplitting := sameCaseSplit(word, tCtx, prefixes, suffixes, tCtx.Score(word))
		splitToken = append(splitToken, sameCaseSplitting...)
	}

	return splitToken
}

func sameCaseSplit(token string, tCtx TokenContext, prefixes lists.List, suffixes lists.List, baseScore float64) []string {
	maxScore := -1.0

	splitToken := []string{token}
	n := len(token)

	for i := 0; i < n; i++ {
		left := token[0:i]
		scoreLeft := tCtx.Score(left)
		shouldSplitLeft := math.Sqrt(scoreLeft) > math.Max(tCtx.Score(token), baseScore)

		right := token[i:n]
		scoreRight := tCtx.Score(right)
		shouldSplitRight := math.Sqrt(scoreRight) > math.Max(tCtx.Score(token), baseScore)

		isPreffixOrSuffix := prefixes.Contains(left) || suffixes.Contains(right)
		if !isPreffixOrSuffix && shouldSplitLeft && shouldSplitRight {
			if (scoreLeft + scoreRight) > maxScore {
				maxScore = scoreLeft + scoreRight
				splitToken = []string{left, right}
			}
		} else if !isPreffixOrSuffix && shouldSplitLeft {
			temp := sameCaseSplit(right, tCtx, prefixes, suffixes, baseScore)
			if len(temp) > 1 {
				splitToken = []string{left}
				splitToken = append(splitToken, temp...)
			}
		}
	}

	return splitToken
}
