package estype

import "fmt"

// Index represents an Elasticsearch index name as a distinct type to prevent misuse.
type Index string

// String returns the string representation of the Index.
func (i Index) String() string {
	return string(i)
}

// ParseESIndex parses an index name string into an Index type.
// Returns an error if the name is empty.
func ParseESIndex(name string) (Index, error) {
	if name == "" {
		return "", fmt.Errorf("index name must not be empty")
	}
	return Index(name), nil
}
