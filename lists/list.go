package lists

// List declares the contract for a list.
type List interface {
	// Contains checks if a word is contained on the list.
	Contains(string) bool
}
