package estype

// RawKeyword configures a keyword sub-field for a text property.
// When passed to NewTextProperty, a "keyword" multi-field is added to the
// text property with the specified IgnoreAbove setting, enabling exact-match
// queries alongside full-text search.
type RawKeyword struct {
	// IgnoreAbove specifies the maximum string length for the keyword sub-field.
	// Strings longer than this value are not indexed.
	IgnoreAbove int
}
