package estype

// MappingProperty is implemented by all Elasticsearch field property types.
// Use this interface as the type for [MappingField.Property] to enforce that
// only recognized property values can be assigned at compile time.
//
// The two built-in implementations are typed property structs such as
// [TextProperty] and [KeywordProperty], and the [FieldType] string type for
// cases where a plain ES type name string suffices (e.g. "integer", "date").
type MappingProperty interface {
	// ESTypeName returns the Elasticsearch type name for this property
	// (e.g. "text", "keyword", "integer").
	ESTypeName() string
}

// FieldType is an Elasticsearch field type name as a typed string.
// It satisfies the [MappingProperty] interface and is the correct way to
// specify a plain ES type name in a [MappingField]:
//
//	estype.MappingField{Path: "price", Property: estype.FieldType("integer")}
//
// Common type names: "text", "keyword", "integer", "long", "float", "double",
// "boolean", "date", "object", "nested", "geo_point", "dense_vector", etc.
type FieldType string

// ESTypeName returns the underlying Elasticsearch type name string.
func (f FieldType) ESTypeName() string { return string(f) }

// Analyzer is a named Elasticsearch analyzer.
// Use a typed Analyzer value instead of a plain string to avoid typos in
// analyzer names when defining field mappings via [ESMappingProvider].
type Analyzer string

// String returns the string representation of the Analyzer.
func (a Analyzer) String() string { return string(a) }

// ---------------------------------------------------------------------------
// Text
// ---------------------------------------------------------------------------

// TextPropertyOption is a functional option for configuring a [TextProperty].
type TextPropertyOption func(*TextProperty)

// TextProperty represents an Elasticsearch "text" field mapping.
// Use [NewTextProperty] to construct one with functional options.
type TextProperty struct {
	// SearchAnalyzer is the analyzer used at query time.
	SearchAnalyzer *Analyzer
	// IndexAnalyzer is the analyzer used at index time.
	IndexAnalyzer *Analyzer
	// Fields holds named multi-field sub-properties (e.g. a keyword sub-field).
	Fields map[string]MappingProperty
}

// ESTypeName returns the Elasticsearch type name for a text property.
func (TextProperty) ESTypeName() string { return "text" }

// NewTextProperty creates a new [TextProperty] with the given options applied.
func NewTextProperty(opts ...TextPropertyOption) TextProperty {
	var p TextProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithSearchAnalyzer sets the analyzer used at query time.
func WithSearchAnalyzer(a Analyzer) TextPropertyOption {
	return func(p *TextProperty) { p.SearchAnalyzer = &a }
}

// WithIndexAnalyzer sets the analyzer used at index (analysis) time.
func WithIndexAnalyzer(a Analyzer) TextPropertyOption {
	return func(p *TextProperty) { p.IndexAnalyzer = &a }
}

// WithField adds a named multi-field sub-property.
// For example, adding a keyword sub-field enables exact-match queries on
// a text field:
//
//	estype.WithField("keyword", estype.NewKeywordProperty())
func WithField(name string, property MappingProperty) TextPropertyOption {
	return func(p *TextProperty) {
		if p.Fields == nil {
			p.Fields = make(map[string]MappingProperty)
		}
		p.Fields[name] = property
	}
}

// ---------------------------------------------------------------------------
// Keyword
// ---------------------------------------------------------------------------

// KeywordPropertyOption is a functional option for configuring a [KeywordProperty].
type KeywordPropertyOption func(*KeywordProperty)

// KeywordProperty represents an Elasticsearch "keyword" field mapping.
// Use [NewKeywordProperty] to construct one with functional options.
type KeywordProperty struct {
	// IgnoreAbove is the maximum string length that will be indexed.
	// Strings longer than this value are not indexed or stored.
	IgnoreAbove *int
}

// ESTypeName returns the Elasticsearch type name for a keyword property.
func (KeywordProperty) ESTypeName() string { return "keyword" }

// NewKeywordProperty creates a new [KeywordProperty] with the given options applied.
func NewKeywordProperty(opts ...KeywordPropertyOption) KeywordProperty {
	var p KeywordProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithIgnoreAbove sets the maximum string length that will be indexed.
// The default value is 256 when no argument is provided.
// Strings longer than this value are not indexed or stored.
func WithIgnoreAbove(v ...int) KeywordPropertyOption {
	n := 256
	if len(v) > 0 {
		n = v[0]
	}
	return func(p *KeywordProperty) { p.IgnoreAbove = &n }
}
