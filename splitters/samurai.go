package splitters

type frequencyTable map[string]float32

var defaultLocalFreqTable frequencyTable
var defaultGlobalFreqTable frequencyTable

// Samurai represents the Samurai splitting algorithm, proposed by Hill et all.
type Samurai struct {
	localFreqTable  *frequencyTable
	globalFreqTable *frequencyTable
}

// NewSamurai creates a new Samurai splitter with the provided frequency tables. If no frequency
// tables are provided, the default tables are used.
func NewSamurai(localFreqTable *frequencyTable, globalFreqTable *frequencyTable) *Samurai {
	var local *frequencyTable
	if localFreqTable != nil {
		local = localFreqTable
	} else {
		local = &defaultLocalFreqTable
	}

	var global *frequencyTable
	if globalFreqTable != nil {
		global = globalFreqTable
	} else {
		global = &defaultGlobalFreqTable
	}

	return &Samurai{
		localFreqTable:  local,
		globalFreqTable: global,
	}
}

// Split on Samurai receives a token and returns an array of hard/soft words,
// split by the Samurai algorithm proposed by Hill et all.
func (s *Samurai) Split(token string) ([]string, error) {
	return []string{token}, nil
}
