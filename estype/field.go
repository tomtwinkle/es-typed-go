package estype

// Field represents an Elasticsearch field name as a distinct type to prevent typos.
// Use the code generator (cmd/estyped) to produce typed Field constants from an
// Elasticsearch index mapping, similar to how sqlc generates Go code from SQL schemas.
type Field string

// String returns the string representation of the Field.
func (f Field) String() string {
	return string(f)
}
