package splitters

// Conserv represents a conservative splitter.
type Conserv struct {
	Splitter
}

// Split on Conserv receives a token and returns an array of hard-words.
func (c Conserv) Split(token string) ([]string, error) {
	return nil, nil
}
