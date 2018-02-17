package splitters

var defaultDictionary map[string]interface{}
var defaultKnownAbbreviations map[string]interface{}
var defaultStopList map[string]interface{}

// Greedy represents the Greedy splitting algorithm, proposed by Field, Binkley and Lawrie.
type Greedy struct {
	Splitter

	dictionary         *map[string]interface{}
	knownAbbreviations *map[string]interface{}
	stopList           *map[string]interface{}
}

// NewGreedy creates a new Greedy splitter, using the specified dictionary.
// If the provided dictionary is null, a default dictionary is used.
func NewGreedy(dicc *map[string]interface{}, knownAbbrs *map[string]interface{}, stop *map[string]interface{}) *Greedy {
	return &Greedy{
		dictionary:         dicc,
		knownAbbreviations: knownAbbrs,
		stopList:           stop,
	}
}

// Split on Greedy receives a token and returns an array of hard/soft words,
// split by the Greedy algorithm proposed by TODO.
func (g *Greedy) Split(token string) ([]string, error) {
	preprocessedToken := addMarkersOnDigits(token)
	preprocessedToken = addMarkersOnLowerToUpperCase(preprocessedToken)

	splitToken := make([]string, 0, 10)
	for _, s := range splitOnMarkers(preprocessedToken) {
		if g.inAnyList(s) != true {
			preffixSplittings := splitOnMarkers(g.findPreffix(s, ""))
			suffixSplittings := splitOnMarkers(g.findSuffix(s, ""))
			chosenSplittings := g.compare(preffixSplittings, suffixSplittings)

			splitToken = append(splitToken, chosenSplittings...)
		} else {
			splitToken = append(splitToken, s)
		}
	}

	return splitToken, nil
}

func (g *Greedy) inAnyList(token string) bool {
	return (*g.dictionary)[token] != nil ||
		(*g.knownAbbreviations)[token] != nil ||
		(*g.stopList)[token] != nil
}

func (g *Greedy) findPreffix(token string, splitToken string) string {
	if len(token) == 0 {
		return ""
	}

	if g.inAnyList(token) {
		return token + "_" + g.findPreffix(token, "")
	}

	sToken := string(token[0]) + splitToken
	s := token[:len(token)-1]

	return g.findPreffix(s, sToken)
}

func (g *Greedy) findSuffix(token string, splitToken string) string {
	if len(token) == 0 {
		return ""
	}

	if g.inAnyList(token) {
		return g.findSuffix(splitToken, "") + "_" + token
	}

	sToken := splitToken + string(token[0])
	s := token[1:len(token)]

	return g.findSuffix(s, sToken)
}

// Compare calculates the ratio between found words on any list and
// the total number of splittings. The string with the highest ratio
// is chosen as the proper splitting.
func (g *Greedy) compare(preffixSplittings []string, suffixSplittings []string) []string {
	if g.calculateWordsInListsRatio(preffixSplittings) > g.calculateWordsInListsRatio(suffixSplittings) {
		return preffixSplittings
	}

	return suffixSplittings
}

// calculateWordsInListsRatio calculates the ratio between the total words
// passed as parameter vs. the total of those words found on any list.
func (g *Greedy) calculateWordsInListsRatio(words []string) float64 {
	found := 0
	for _, word := range words {
		if g.inAnyList(word) {
			found++
		}
	}

	ratio := 0.0
	if len(words) > 0 {
		ratio = float64(found) / float64(len(words))
	}

	return ratio
}
