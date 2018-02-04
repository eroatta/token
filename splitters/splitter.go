package splitters

// Splitter defines the required behavior for any token splitter.
type Splitter interface {
	Split(token string) ([]string, error)
}
