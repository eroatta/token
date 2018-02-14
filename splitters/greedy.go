package splitters

var dictionary map[string]interface{}
var knownAbbreviations map[string]interface{}
var stopList map[string]interface{}

// Greedy represents the Greedy splitting algorithm, proposed by Field, Binkley and Lawrie.
type Greedy struct {
	Splitter
}

// Split on Greedy receives a token and returns an array of hard/soft words,
// split by the Greedy algorithm proposed by TODO.
func (g Greedy) Split(token string) ([]string, error) {
	preprocessedToken := addMarkersOnDigits(token)
	preprocessedToken = addMarkersOnLowerToUpperCase(preprocessedToken)

	splitToken := make([]string, 0, 10)
	for _, s := range splitOnMarkers(preprocessedToken) {
		if inAnyList(s) != true {
			preffixSplittings := splitOnMarkers(findPreffix(s, ""))
			suffixSplittings := splitOnMarkers(findSuffix(s, ""))
			chosenSplittings := compare(preffixSplittings, suffixSplittings)

			splitToken = append(splitToken, chosenSplittings...)
		} else {
			splitToken = append(splitToken, s)
		}
	}

	return splitToken, nil
}

func inAnyList(token string) bool {
	return dictionary[token] != nil || knownAbbreviations[token] != nil || stopList[token] != nil
}

func findPreffix(token string, splitToken string) string {
	if len(token) == 0 {
		return ""
	}

	if inAnyList(token) == true {
		return token + "_" + findPreffix(token, "")
	}

	sToken := string(token[0]) + splitToken
	s := token[:len(token)-1]

	return findPreffix(s, sToken)
}

func findSuffix(token string, splitToken string) string {
	return ""
}

// Compare calculates the ratio between found words on any list and
// the total number of splittings. The string with the highest ratio
// is chosen as the proper splitting.
func compare(preffixSplittings []string, suffixSplittings []string) []string {
	var splittings []string
	if calculateWordsInListsRatio(preffixSplittings) > calculateWordsInListsRatio(suffixSplittings) {
		splittings = preffixSplittings
	} else {
		splittings = suffixSplittings
	}

	return splittings
}

// calculateWordsInListsRatio calculates the ratio between the total words
// passed as parameter vs. the total of those words found on any list.
func calculateWordsInListsRatio(words []string) float64 {
	found := 0
	for _, word := range words {
		if inAnyList(word) {
			found++
		}
	}

	ratio := 0.0
	if len(words) > 0 {
		ratio = float64(found) / float64(len(words))
	}

	return ratio
}
