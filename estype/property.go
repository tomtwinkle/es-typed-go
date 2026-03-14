package estype

// Deprecated: RawKeyword is deprecated. Use WithTextRawKeyword functional option instead.
// RawKeyword configures a keyword sub-field for a text property.
type RawKeyword struct {
	// IgnoreAbove specifies the maximum string length for the keyword sub-field.
	// Strings longer than this value are not indexed.
	IgnoreAbove int
}
