package expanders

// Basic represents the Basic expansion algorithm, proposed by Lawrie, Feild and Binkley.
type Basic struct {
	srcWords    map[string]interface{}
	srcPhrases  map[string]string
	stopList    map[string]interface{}
	dicctionary map[string]interface{}
}

// NewBasic creates a new Basic expander with the given lists.
func NewBasic(srcWords map[string]interface{}, srcPhrases map[string]string, stopList map[string]interface{}, dicc map[string]interface{}) *Basic {
	return &Basic{
		srcWords:    srcWords,
		srcPhrases:  srcPhrases,
		stopList:    stopList,
		dicctionary: dicc,
	}
}

// Expand on Basic receives a token and returns an array of possible expansions.
func (b *Basic) Expand(token string) ([]string, error) {
	// the first list includes words contained in the comments that appear before and within the function
	// the first list also includes dictionary hard words found in the identifiers of the function
	/*srcWords := mapset.NewSet()

	// the phrase list is obtained by running the comments and multiword-identifiers through a phrase finder
	srcPhrases := make(map[string]string)

	stoplist := *b.stopList
	if stoplist.Contains(token) {
		return []string{token}, nil
	}

	if srcPhrases[token] != "" {
		return []string{srcPhrases[token]}, nil
	}

	if srcWords.Contains(token) {
		return []string{token}, nil
	}

	var expansions []string

	// build the search regex

	//it := b.dicc.Iterator()
	/*for elem := range it.C {

	}*/

	return nil, nil
}
