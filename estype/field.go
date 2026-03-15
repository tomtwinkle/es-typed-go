package estype

// Field represents an Elasticsearch field name as a distinct type to prevent typos.
// Use the code generator (cmd/estyped) to produce typed Field constants from an
// Elasticsearch index mapping, similar to how sqlc generates Go code from SQL schemas.
type Field string

// String returns the string representation of the Field.
func (f Field) String() string {
	return string(f)
}

// FieldNames converts a list of Field values to a []string slice.
// This is useful when passing typed Field constants to ES query types that
// accept []string, such as MultiMatchQuery.Fields.
func FieldNames(fields ...Field) []string {
	result := make([]string, len(fields))
	for i, f := range fields {
		result[i] = string(f)
	}
	return result
}
