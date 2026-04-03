package estype

import "fmt"

// Alias represents an Elasticsearch alias name as a distinct type to prevent misuse.
type Alias string

// String returns the string representation of the Alias.
func (a Alias) String() string {
	return string(a)
}

// Ptr returns a pointer to the string representation of the Alias.
// This is useful when passing typed Alias constants to ES query types that accept *string.
func (a Alias) Ptr() *string {
	s := string(a)
	return &s
}

// ParseESAlias parses an alias name string into an Alias type.
// Returns an error if the name is empty.
func ParseESAlias(name string) (Alias, error) {
	if name == "" {
		return "", fmt.Errorf("alias name must not be empty")
	}
	return Alias(name), nil
}

// AliasProvider is implemented by types that declare a canonical Elasticsearch alias name.
// The estyped generator reads this method when running in struct mode with the -group flag
// to include a typed Alias field in the generated model accessor.
//
// Example usage in a definition file:
//
//	func (Product) Alias() Alias { return "product" }
type AliasProvider interface {
	Alias() Alias
}
