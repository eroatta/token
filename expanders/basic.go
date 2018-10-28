package expanders

import (
	"github.com/deckarep/golang-set"
)

// Basic represents a Basic expansion algorithm, proposed by Lawrie, Feild and Binkley.
type Basic struct {
	stopList *mapset.Set
	dicc     *mapset.Set
}

// NewBasic creates a new Basic expander with the given stop lists.
func NewBasic(stopList *mapset.Set, dicc *mapset.Set) *Basic {
	return &Basic{
		stopList: stopList,
		dicc:     dicc,
	}
}

// Expand on Basic receives a token and returns an array of possible expansions.
func (b *Basic) Expand(token string) ([]string, error) {
	// the first list includes words contained in the comments that appear before and within the function
	// the first list also includes dictionary hard words found in the identifiers of the function
	srcWords := mapset.NewSet()

	// the phrase list is obtained by running the comments and multiword-identifiers through a phrase finder
	srcPhrases := make(map[string]string)

	stoplist := *b.stopList
	if stoplist.Contains(token) {
		return token, nil
	}

	if srcPhrases[token] != "" {
		return srcPhrases[token], nil
	}

	if srcWords.Contains(token) {
		return token, nil
	}

	var expansions []string

	// build the search regex

	it := b.dicc.Iterator()
	for elem := range it.C {

	}

	return expansions, nil
}
