// Package lists provides common lists, and also ways to create custom lists.
package lists

import (
	"strings"
)

var (
	// Dictionary is a list of words extracted from the aspell GNU tool.
	Dictionary = NewBuilder().Add(dictionary...).Build()
	// KnownAbbreviations is a list of strings that are known and common abbreviations on the language.
	KnownAbbreviations = NewBuilder().Add(knownAbbreviations...).Build()
	// Stop is a list of reserved words, data types and Go library names.
	Stop = NewBuilder().Add(stop...).Build()
	// Prefixes is a list of common prefixes.
	Prefixes = NewBuilder().Add(prefixes...).Build()
	// Suffixes is a list of common suffixes.
	Suffixes = NewBuilder().Add(suffixes...).Build()
)

// List declares the contract for a list.
type List interface {
	// Contains checks if a word is contained on the list.
	Contains(string) bool
	// Size returns the number elements on the list.
	Size() int
	// Elements returns an array with all the elements on the list.
	Elements() []string
}

type list struct {
	elements map[string]bool
}

func (l list) Contains(element string) bool {
	return l.elements[strings.ToLower(element)]
}

func (l list) Size() int {
	return len(l.elements)
}

func (l list) Elements() []string {
	keys := make([]string, len(l.elements))
	i := 0
	for k := range l.elements {
		keys[i] = k
		i++
	}

	return keys
}

// NewBuilder creates a new ListBuilder.
func NewBuilder() ListBuilder {
	return &listBuilder{
		elements: make(map[string]bool),
	}
}

// ListBuilder builds a list based on the added elements.
type ListBuilder interface {
	// Add adds one or more elements to list.
	Add(...string) ListBuilder
	// Build creates the list with the given elements.
	Build() List
}

type listBuilder struct {
	elements map[string]bool
}

func (lb *listBuilder) Add(elements ...string) ListBuilder {
	for _, e := range elements {
		lb.elements[strings.ToLower(e)] = true
	}

	return lb
}

func (lb *listBuilder) Build() List {
	return list{elements: lb.elements}
}
