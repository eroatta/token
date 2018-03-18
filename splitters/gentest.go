package splitters

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
	return []string{}, nil
}
