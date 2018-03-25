package splitters

import (
	"sync"
)

// GenTest represents a Generation and Tests splitting algorithm, proposed by Lawrie, Binkley and Morrell.
type GenTest struct {
}

// NewGenTest creates a new GenTest splitter.
func NewGenTest() *GenTest {
	return &GenTest{}
}

// Split on GenTest receives a token and returns an array of hard/soft words,
// split by the Generation and Test algorithm proposed by Lawrie, Binkley and Morrell.
func (g *GenTest) Split(token string) ([]string, error) {
	preprocessedToken := addMarkersOnDigits(token)
	preprocessedToken = addMarkersOnLowerToUpperCase(preprocessedToken)

	splitToken := make([]string, 0, 10)
	for _, word := range splitOnMarkers(preprocessedToken) {
		splitToken = append(splitToken, word)
	}

	return splitToken, nil
}

func (g *GenTest) generate(token string) []string {
	var selectedSplit []string

	wg := new(sync.WaitGroup)

	potentialSplits := make(chan []string)

	size := len(token)
	for i := 0; i < len(token); i++ {
		wg.Add(1)
		go func(i int, n int, token string, wg *sync.WaitGroup, splits *chan []string) {
			defer wg.Done()
			*splits <- []string{token[0:i], token[i:n]}
		}(i, size, token, wg, &potentialSplits)
	}
	wg.Wait()

	close(potentialSplits)

	return selectedSplit
}

// generateSplits creates every possible splitting for a given token.
func generateSplits(token string) []string {
	splits := []string{token}
	for i := 1; i < len(token); i++ {
		leading := token[:i]
		trailing := token[i:]

		split := leading + "_" + trailing
		splits = append(splits, split)

		for j := 1; j < len(trailing); j++ {
			split = leading + "_" + trailing[:j] + "_" + trailing[j:]
			splits = append(splits, split)
		}
	}

	return splits
}

func testSplit(split string) {

}

func similarity() {

}
