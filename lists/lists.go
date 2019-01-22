package lists

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Dicctionary is a list of words extracted from the aspell GNU tool.
var Dicctionary map[string]interface{}

// KnownAbbreviations is a list of strings that are known and common abbreviations on the language.
var KnownAbbreviations map[string]interface{}

// StopList is a list of reserved words, data types and Go library names.
var StopList map[string]interface{}

func init() {
	Dicctionary = loadFile("dicctionary.txt")
	KnownAbbreviations = loadFile("known_abbreviations.txt")
	StopList = loadFile("stoplist.txt")
}

// loadFile reads a file, line by line, and builds a set of strings from it.
func loadFile(filename string) map[string]interface{} {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("Unable to open file %s: %v", filename, err))
	}
	defer file.Close()

	set := make(map[string]interface{}, 100)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		set[strings.TrimSpace(scanner.Text())] = true
	}

	return set
}
