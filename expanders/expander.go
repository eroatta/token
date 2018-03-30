package expanders

// Expander defines the required behaviour for any token expander.
type Expander interface {
	Expand(token string) ([]string, error)
}
