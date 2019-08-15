package expansion

import (
	"sort"
	"strings"

	"github.com/eroatta/token/lists"
)

// Set represents a set expansions stored in convenient format.
type Set interface {
	// Array returns a string array representation of the given set.
	Array() []string
	// String returns a string representation of the given set.
	String() string
	// Contains checks if a word is contained on the set.
	Contains(string) bool
}

// Set represents a set expansions stored in convenient format.
type set struct {
	words         lists.List
	wordsAsString string
}

func (s set) Array() []string {
	return s.words.Elements()
}

func (s set) String() string {
	return s.wordsAsString
}

func (s set) Contains(word string) bool {
	return s.words.Contains(word)
}

// NewSetBuilder creates a new SetBuilder.
func NewSetBuilder() SetBuilder {
	return &setBuilder{
		wb: lists.NewBuilder(),
	}
}

// SetBuilder builds an expansion set based on the added lists or strings.
type SetBuilder interface {
	AddList(lists.List) SetBuilder
	AddStrings(...string) SetBuilder
	Build() Set
}

type setBuilder struct {
	wb lists.ListBuilder
}

func (sb *setBuilder) AddList(list lists.List) SetBuilder {
	if list != nil {
		sb.wb.Add(list.Elements()...)
	}
	return sb
}

func (sb *setBuilder) AddStrings(strs ...string) SetBuilder {
	sb.wb.Add(strs...)
	return sb
}

func (sb *setBuilder) Build() Set {
	list := sb.wb.Build()
	elements := list.Elements()
	sort.Strings(elements)

	return &set{
		words:         list,
		wordsAsString: strings.Join(elements, " "),
	}
}
