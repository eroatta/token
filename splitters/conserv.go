package splitters

// Conserv represents a conservative splitter.
type Conserv struct {
	Splitter
}

// Split on Conserv receives a token and returns an array of hard/soft words,
// split by:
// * Underscores
// * Numbers
// * CamelCase.
func (c Conserv) Split(token string) ([]string, error) {
	processedToken := addMarkersOnDigits(token)
	processedToken = addMarkersOnLowerToUpperCase(processedToken)
	processedToken = addMarkersOnUpperToLowerCase(processedToken)

	return splitOnMarkers(processedToken), nil
}
