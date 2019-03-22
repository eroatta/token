package lists

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/packr/v2"

	"github.com/deckarep/golang-set"
)

// Dicctionary is a list of words extracted from the aspell GNU tool.
var Dicctionary map[string]interface{}

// KnownAbbreviations is a list of strings that are known and common abbreviations on the language.
var KnownAbbreviations map[string]interface{}

// StopList is a list of reserved words, data types and Go library names.
var StopList map[string]interface{}

// Prefixes is a list of common prefixes.
var Prefixes mapset.Set

// Suffixes is a list of common suffixes.
var Suffixes mapset.Set

func init() {
	box := packr.New("lists", ".")

	Dicctionary = buildMapFromFile(box, "dicctionary.txt")
	KnownAbbreviations = buildMapFromFile(box, "known_abbreviations.txt")
	StopList = buildMapFromFile(box, "stoplist.txt")
	Prefixes = buildSetFromFile(box, "prefixes.txt")
	Suffixes = buildSetFromFile(box, "suffixes.txt")
}

func readLinesFromFile(box *packr.Box, filename string) []string {
	content, err := box.Find(filename)
	if err != nil {
		panic(fmt.Sprintf("Unable to open file %s: %v", filename, err))
	}

	return strings.Split(string(content), "\n")
}

// buildMapFromFile reads a file, line by line, and builds a set of strings from it.
func buildMapFromFile(box *packr.Box, filename string) map[string]interface{} {
	set := make(map[string]interface{}, 100)
	for _, line := range readLinesFromFile(box, filename) {
		set[strings.TrimSpace(line)] = true
	}

	return set
}

// loadFileIntoSet reads a file, line by line, and builds a mapset.Set of strings from it.
func buildSetFromFile(box *packr.Box, filename string) mapset.Set {
	set := mapset.NewSet()
	for _, line := range readLinesFromFile(box, filename) {
		set.Add(strings.TrimSpace(line))
	}

	return set
}
