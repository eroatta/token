package splitters

// Splitter defines the common behavior on any token splitter.
type Splitter interface {
	Split(token string) ([]string, error)
}
