package expanders

// Amap represents an Automatically Mining Abbreviations in Programs expander.
type Amap struct {
}

// NewAmap creates an AMAP expander.
func NewAmap() *Amap {
	return &Amap{}
}

// Expand on AMAP receives a token and returns and array of possible expansions.
//
// The AMAP expansion algorithm handles single-word and multi-word abbreviations.
// For each type of abbreviation AMAP creates and applies a pattern to look for possible
// expansions. AMAP is capable of select the more appropiate expansions based on available
// information on the given context.
func (a Amap) Expand(token string) []string {
	return []string{}
}
