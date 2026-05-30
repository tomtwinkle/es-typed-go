package estype

import (
	"fmt"
	"strings"
)

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

// ParseESAlias parses one or more alias name strings into a single Alias.
// Multiple names are joined with "," which Elasticsearch interprets as a
// combined target — e.g. ParseESAlias("orders", "archive") produces the
// alias "orders,archive" that queries both indices in one request.
// Returns an error if no names are provided or any name is empty.
func ParseESAlias(names ...string) (Alias, error) {
	if len(names) == 0 {
		return "", fmt.Errorf("at least one alias name must be provided")
	}
	for _, name := range names {
		if name == "" {
			return "", fmt.Errorf("alias name must not be empty")
		}
	}
	return Alias(strings.Join(names, ",")), nil
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
