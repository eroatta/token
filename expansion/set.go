package expansion

import (
	"strings"

	"github.com/eroatta/token-splitex/lists"
)

// Set represents a set expansions stored in convenient format.
type Set interface {
	// Array returns a string array representation of the given set.
	Array() []string
	// String returns a string representation of the given set.
	String() string
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
	return &set{
		words:         list,
		wordsAsString: strings.Join(list.Elements(), " "),
	}
}
