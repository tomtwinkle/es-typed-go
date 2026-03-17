package estype

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
	Fields map[string]any
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
func WithField(name string, property any) TextPropertyOption {
	return func(p *TextProperty) {
		if p.Fields == nil {
			p.Fields = make(map[string]any)
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
