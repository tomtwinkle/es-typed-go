package estype

// WriteIndex specifies whether an index should be designated as the write target for an alias.
type WriteIndex int8

const (
	// WriteIndexUnset leaves the is_write_index field unset (uses Elasticsearch default behaviour).
	WriteIndexUnset WriteIndex = iota
	// WriteIndexEnabled designates the index as the write target for the alias (is_write_index = true).
	WriteIndexEnabled
	// WriteIndexDisabled explicitly marks the index as not the write target (is_write_index = false).
	WriteIndexDisabled
)

// BoolPtr converts a WriteIndex to a *bool suitable for Elasticsearch requests.
// WriteIndexUnset returns nil.
func (w WriteIndex) BoolPtr() *bool {
	switch w {
	case WriteIndexEnabled:
		v := true
		return &v
	case WriteIndexDisabled:
		v := false
		return &v
	default:
		return nil
	}
}
