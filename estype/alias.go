package estype

import "fmt"

// Alias represents an Elasticsearch alias name as a distinct type to prevent misuse.
type Alias string

// String returns the string representation of the Alias.
func (a Alias) String() string {
	return string(a)
}

// ParseESAlias parses an alias name string into an Alias type.
// Returns an error if the name is empty.
func ParseESAlias(name string) (Alias, error) {
	if name == "" {
		return "", fmt.Errorf("alias name must not be empty")
	}
	return Alias(name), nil
}
