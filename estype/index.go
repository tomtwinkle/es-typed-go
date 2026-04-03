package estype

import "fmt"

// Index represents an Elasticsearch index name as a distinct type to prevent misuse.
type Index string

// String returns the string representation of the Index.
func (i Index) String() string {
	return string(i)
}

// Ptr returns a pointer to the string representation of the Index.
// This is useful when passing typed Index constants to ES query types that accept *string.
func (i Index) Ptr() *string {
	s := string(i)
	return &s
}

// ParseESIndex parses an index name string into an Index type.
// Returns an error if the name is empty.
func ParseESIndex(name string) (Index, error) {
	if name == "" {
		return "", fmt.Errorf("index name must not be empty")
	}
	return Index(name), nil
}

// IndexProvider is implemented by types that declare a canonical Elasticsearch index name.
// The estyped generator reads this method when running in struct mode with the -group flag
// to include a typed Index field in the generated model accessor.
//
// Example usage in a definition file:
//
//	func (Product) Index() estype.Index { return "product-000001" }
type IndexProvider interface {
	Index() Index
}
