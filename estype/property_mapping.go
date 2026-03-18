package estype

// MappingProperty is implemented by all Elasticsearch field property types.
// Use this interface as the type for [MappingField.Property] to enforce that
// only recognized property values can be assigned at compile time.
//
// Implementations include typed property structs such as [TextProperty],
// [KeywordProperty], [DateProperty], [ObjectProperty], [NestedProperty],
// [DenseVectorProperty], and others, as well as the [FieldType] string type
// for cases where a plain ES type name string suffices.
type MappingProperty interface {
	// ESTypeName returns the Elasticsearch type name for this property
	// (e.g. "text", "keyword", "integer").
	ESTypeName() string
}

// FieldType is an Elasticsearch field type name as a typed string.
// It satisfies the [MappingProperty] interface and is the correct way to
// specify a plain ES type name in a [MappingField]:
//
// estype.MappingField{Path: "price", Property: estype.FieldType("integer")}
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
	// SearchQuoteAnalyzer is the analyzer used for quoted phrases at query time.
	SearchQuoteAnalyzer *string
	// Fielddata enables in-memory fielddata for sorting, aggregations, and scripting.
	Fielddata *bool
	// Index controls whether the field is indexed.
	Index *bool
	// Store controls whether the field value is stored separately.
	Store *bool
	// Norms controls whether norms are enabled for scoring.
	Norms *bool
	// Similarity is the similarity algorithm to use.
	Similarity *string
	// IndexPhrases controls whether two-term word combinations are indexed.
	IndexPhrases *bool
	// PositionIncrementGap is the number of fake term positions between indexed values.
	PositionIncrementGap *int
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
// estype.WithField("keyword", estype.NewKeywordProperty())
func WithField(name string, property MappingProperty) TextPropertyOption {
	return func(p *TextProperty) {
		if p.Fields == nil {
			p.Fields = make(map[string]MappingProperty)
		}
		p.Fields[name] = property
	}
}

// WithTextSearchQuoteAnalyzer sets the search quote analyzer for the text field.
func WithTextSearchQuoteAnalyzer(v string) TextPropertyOption {
	return func(p *TextProperty) { p.SearchQuoteAnalyzer = &v }
}

// WithTextFielddata sets whether fielddata is enabled.
func WithTextFielddata(v bool) TextPropertyOption {
	return func(p *TextProperty) { p.Fielddata = &v }
}

// WithTextIndex sets whether the field is indexed.
func WithTextIndex(v bool) TextPropertyOption {
	return func(p *TextProperty) { p.Index = &v }
}

// WithTextStore sets whether the field value is stored.
func WithTextStore(v bool) TextPropertyOption {
	return func(p *TextProperty) { p.Store = &v }
}

// WithTextNorms sets whether norms are enabled.
func WithTextNorms(v bool) TextPropertyOption {
	return func(p *TextProperty) { p.Norms = &v }
}

// WithTextSimilarity sets the similarity algorithm.
func WithTextSimilarity(v string) TextPropertyOption {
	return func(p *TextProperty) { p.Similarity = &v }
}

// WithTextIndexPhrases sets whether two-term word combinations are indexed.
func WithTextIndexPhrases(v bool) TextPropertyOption {
	return func(p *TextProperty) { p.IndexPhrases = &v }
}

// WithTextPositionIncrementGap sets the number of fake term positions between indexed values.
func WithTextPositionIncrementGap(v int) TextPropertyOption {
	return func(p *TextProperty) { p.PositionIncrementGap = &v }
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
	IgnoreAbove *int
	// DocValues controls whether doc values are enabled.
	DocValues *bool
	// Index controls whether the field is indexed.
	Index *bool
	// Store controls whether the field value is stored separately.
	Store *bool
	// NullValue is the value substituted for explicit null values.
	NullValue *string
	// Normalizer is the normalizer applied to the keyword field.
	Normalizer *string
	// Norms controls whether norms are enabled for scoring.
	Norms *bool
	// Similarity is the similarity algorithm to use.
	Similarity *string
	// EagerGlobalOrdinals controls whether global ordinals are loaded eagerly.
	EagerGlobalOrdinals *bool
	// SplitQueriesOnWhitespace controls whether queries are split on whitespace.
	SplitQueriesOnWhitespace *bool
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

// WithKeywordDocValues sets whether doc values are enabled.
func WithKeywordDocValues(v bool) KeywordPropertyOption {
	return func(p *KeywordProperty) { p.DocValues = &v }
}

// WithKeywordIndex sets whether the field is indexed.
func WithKeywordIndex(v bool) KeywordPropertyOption {
	return func(p *KeywordProperty) { p.Index = &v }
}

// WithKeywordStore sets whether the field value is stored.
func WithKeywordStore(v bool) KeywordPropertyOption {
	return func(p *KeywordProperty) { p.Store = &v }
}

// WithKeywordNullValue sets the null value for the field.
func WithKeywordNullValue(v string) KeywordPropertyOption {
	return func(p *KeywordProperty) { p.NullValue = &v }
}

// WithKeywordNormalizer sets the normalizer for the keyword field.
func WithKeywordNormalizer(v string) KeywordPropertyOption {
	return func(p *KeywordProperty) { p.Normalizer = &v }
}

// WithKeywordNorms sets whether norms are enabled.
func WithKeywordNorms(v bool) KeywordPropertyOption {
	return func(p *KeywordProperty) { p.Norms = &v }
}

// WithKeywordSimilarity sets the similarity algorithm.
func WithKeywordSimilarity(v string) KeywordPropertyOption {
	return func(p *KeywordProperty) { p.Similarity = &v }
}

// WithKeywordEagerGlobalOrdinals sets whether to eagerly load global ordinals.
func WithKeywordEagerGlobalOrdinals(v bool) KeywordPropertyOption {
	return func(p *KeywordProperty) { p.EagerGlobalOrdinals = &v }
}

// WithKeywordSplitQueriesOnWhitespace sets whether to split queries on whitespace.
func WithKeywordSplitQueriesOnWhitespace(v bool) KeywordPropertyOption {
	return func(p *KeywordProperty) { p.SplitQueriesOnWhitespace = &v }
}

// ---------------------------------------------------------------------------
// Constant Keyword
// ---------------------------------------------------------------------------

// ConstantKeywordPropertyOption is a functional option for configuring a [ConstantKeywordProperty].
type ConstantKeywordPropertyOption func(*ConstantKeywordProperty)

// ConstantKeywordProperty represents an Elasticsearch "constant_keyword" field mapping.
// All documents in the index have the same value for this field.
// Use [NewConstantKeywordProperty] to construct one with functional options.
type ConstantKeywordProperty struct {
	// Value is the constant value for this field across all documents.
	Value *string
}

// ESTypeName returns the Elasticsearch type name for a constant_keyword property.
func (ConstantKeywordProperty) ESTypeName() string { return "constant_keyword" }

// NewConstantKeywordProperty creates a new [ConstantKeywordProperty] with the given options applied.
func NewConstantKeywordProperty(opts ...ConstantKeywordPropertyOption) ConstantKeywordProperty {
	var p ConstantKeywordProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithConstantKeywordValue sets the constant value for this field.
func WithConstantKeywordValue(v string) ConstantKeywordPropertyOption {
	return func(p *ConstantKeywordProperty) { p.Value = &v }
}

// ---------------------------------------------------------------------------
// Counted Keyword
// ---------------------------------------------------------------------------

// CountedKeywordPropertyOption is a functional option for configuring a [CountedKeywordProperty].
type CountedKeywordPropertyOption func(*CountedKeywordProperty)

// CountedKeywordProperty represents an Elasticsearch "counted_keyword" field mapping.
// Use [NewCountedKeywordProperty] to construct one with functional options.
type CountedKeywordProperty struct {
	// Index controls whether the field is indexed.
	Index *bool
}

// ESTypeName returns the Elasticsearch type name for a counted_keyword property.
func (CountedKeywordProperty) ESTypeName() string { return "counted_keyword" }

// NewCountedKeywordProperty creates a new [CountedKeywordProperty] with the given options applied.
func NewCountedKeywordProperty(opts ...CountedKeywordPropertyOption) CountedKeywordProperty {
	var p CountedKeywordProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithCountedKeywordIndex sets whether the field is indexed.
func WithCountedKeywordIndex(v bool) CountedKeywordPropertyOption {
	return func(p *CountedKeywordProperty) { p.Index = &v }
}

// ---------------------------------------------------------------------------
// Wildcard
// ---------------------------------------------------------------------------

// WildcardPropertyOption is a functional option for configuring a [WildcardProperty].
type WildcardPropertyOption func(*WildcardProperty)

// WildcardProperty represents an Elasticsearch "wildcard" field mapping.
// Use [NewWildcardProperty] to construct one with functional options.
type WildcardProperty struct {
	// IgnoreAbove is the maximum string length that will be indexed.
	IgnoreAbove *int
	// DocValues controls whether doc values are enabled.
	DocValues *bool
	// NullValue is the value substituted for explicit null values.
	NullValue *string
	// Store controls whether the field value is stored separately.
	Store *bool
}

// ESTypeName returns the Elasticsearch type name for a wildcard property.
func (WildcardProperty) ESTypeName() string { return "wildcard" }

// NewWildcardProperty creates a new [WildcardProperty] with the given options applied.
func NewWildcardProperty(opts ...WildcardPropertyOption) WildcardProperty {
	var p WildcardProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithWildcardIgnoreAbove sets the maximum string length that will be indexed.
func WithWildcardIgnoreAbove(v int) WildcardPropertyOption {
	return func(p *WildcardProperty) { p.IgnoreAbove = &v }
}

// WithWildcardDocValues sets whether doc values are enabled.
func WithWildcardDocValues(v bool) WildcardPropertyOption {
	return func(p *WildcardProperty) { p.DocValues = &v }
}

// WithWildcardNullValue sets the null value for the field.
func WithWildcardNullValue(v string) WildcardPropertyOption {
	return func(p *WildcardProperty) { p.NullValue = &v }
}

// WithWildcardStore sets whether the field value is stored.
func WithWildcardStore(v bool) WildcardPropertyOption {
	return func(p *WildcardProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Match Only Text
// ---------------------------------------------------------------------------

// MatchOnlyTextPropertyOption is a functional option for configuring a [MatchOnlyTextProperty].
type MatchOnlyTextPropertyOption func(*MatchOnlyTextProperty)

// MatchOnlyTextProperty represents an Elasticsearch "match_only_text" field mapping.
// A space-optimized variant of text that disables scoring and positional queries.
// Use [NewMatchOnlyTextProperty] to construct one with functional options.
type MatchOnlyTextProperty struct {
	// Fields holds named multi-field sub-properties.
	Fields map[string]MappingProperty
}

// ESTypeName returns the Elasticsearch type name for a match_only_text property.
func (MatchOnlyTextProperty) ESTypeName() string { return "match_only_text" }

// NewMatchOnlyTextProperty creates a new [MatchOnlyTextProperty] with the given options applied.
func NewMatchOnlyTextProperty(opts ...MatchOnlyTextPropertyOption) MatchOnlyTextProperty {
	var p MatchOnlyTextProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithMatchOnlyTextField adds a named multi-field sub-property.
func WithMatchOnlyTextField(name string, property MappingProperty) MatchOnlyTextPropertyOption {
	return func(p *MatchOnlyTextProperty) {
		if p.Fields == nil {
			p.Fields = make(map[string]MappingProperty)
		}
		p.Fields[name] = property
	}
}

// ---------------------------------------------------------------------------
// Completion
// ---------------------------------------------------------------------------

// CompletionPropertyOption is a functional option for configuring a [CompletionProperty].
type CompletionPropertyOption func(*CompletionProperty)

// CompletionProperty represents an Elasticsearch "completion" field mapping
// used for auto-complete suggestions.
// Use [NewCompletionProperty] to construct one with functional options.
type CompletionProperty struct {
	// Analyzer is the analyzer used for this completion field.
	Analyzer *Analyzer
	// SearchAnalyzer is the analyzer used at query time.
	SearchAnalyzer *Analyzer
	// MaxInputLength is the maximum length of a single input.
	MaxInputLength *int
	// PreservePositionIncrements controls whether position increments are preserved.
	PreservePositionIncrements *bool
	// PreserveSeparators controls whether separators are preserved.
	PreserveSeparators *bool
}

// ESTypeName returns the Elasticsearch type name for a completion property.
func (CompletionProperty) ESTypeName() string { return "completion" }

// NewCompletionProperty creates a new [CompletionProperty] with the given options applied.
func NewCompletionProperty(opts ...CompletionPropertyOption) CompletionProperty {
	var p CompletionProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithCompletionAnalyzer sets the analyzer for the completion field.
func WithCompletionAnalyzer(a Analyzer) CompletionPropertyOption {
	return func(p *CompletionProperty) { p.Analyzer = &a }
}

// WithCompletionSearchAnalyzer sets the query-time analyzer for the completion field.
func WithCompletionSearchAnalyzer(a Analyzer) CompletionPropertyOption {
	return func(p *CompletionProperty) { p.SearchAnalyzer = &a }
}

// WithCompletionMaxInputLength sets the maximum length of a single input.
func WithCompletionMaxInputLength(v int) CompletionPropertyOption {
	return func(p *CompletionProperty) { p.MaxInputLength = &v }
}

// WithCompletionPreservePositionIncrements sets whether to preserve position increments.
func WithCompletionPreservePositionIncrements(v bool) CompletionPropertyOption {
	return func(p *CompletionProperty) { p.PreservePositionIncrements = &v }
}

// WithCompletionPreserveSeparators sets whether to preserve separators.
func WithCompletionPreserveSeparators(v bool) CompletionPropertyOption {
	return func(p *CompletionProperty) { p.PreserveSeparators = &v }
}

// ---------------------------------------------------------------------------
// Search As You Type
// ---------------------------------------------------------------------------

// SearchAsYouTypePropertyOption is a functional option for configuring a [SearchAsYouTypeProperty].
type SearchAsYouTypePropertyOption func(*SearchAsYouTypeProperty)

// SearchAsYouTypeProperty represents an Elasticsearch "search_as_you_type" field mapping.
// Use [NewSearchAsYouTypeProperty] to construct one with functional options.
type SearchAsYouTypeProperty struct {
	// Analyzer is the analyzer used for this field.
	Analyzer *Analyzer
	// SearchAnalyzer is the analyzer used at query time.
	SearchAnalyzer *Analyzer
	// MaxShingleSize is the maximum number of terms in the shingle sub-field (2-4).
	MaxShingleSize *int
	// SearchQuoteAnalyzer is the analyzer used for quoted phrases at query time.
	SearchQuoteAnalyzer *string
	// Index controls whether the field is indexed.
	Index *bool
	// Store controls whether the field value is stored separately.
	Store *bool
	// Norms controls whether norms are enabled for scoring.
	Norms *bool
	// Similarity is the similarity algorithm to use.
	Similarity *string
}

// ESTypeName returns the Elasticsearch type name for a search_as_you_type property.
func (SearchAsYouTypeProperty) ESTypeName() string { return "search_as_you_type" }

// NewSearchAsYouTypeProperty creates a new [SearchAsYouTypeProperty] with the given options applied.
func NewSearchAsYouTypeProperty(opts ...SearchAsYouTypePropertyOption) SearchAsYouTypeProperty {
	var p SearchAsYouTypeProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithSearchAsYouTypeAnalyzer sets the analyzer.
func WithSearchAsYouTypeAnalyzer(a Analyzer) SearchAsYouTypePropertyOption {
	return func(p *SearchAsYouTypeProperty) { p.Analyzer = &a }
}

// WithSearchAsYouTypeSearchAnalyzer sets the query-time analyzer.
func WithSearchAsYouTypeSearchAnalyzer(a Analyzer) SearchAsYouTypePropertyOption {
	return func(p *SearchAsYouTypeProperty) { p.SearchAnalyzer = &a }
}

// WithSearchAsYouTypeMaxShingleSize sets the maximum shingle size (2-4).
func WithSearchAsYouTypeMaxShingleSize(v int) SearchAsYouTypePropertyOption {
	return func(p *SearchAsYouTypeProperty) { p.MaxShingleSize = &v }
}

// WithSearchAsYouTypeSearchQuoteAnalyzer sets the search quote analyzer.
func WithSearchAsYouTypeSearchQuoteAnalyzer(v string) SearchAsYouTypePropertyOption {
	return func(p *SearchAsYouTypeProperty) { p.SearchQuoteAnalyzer = &v }
}

// WithSearchAsYouTypeIndex sets whether the field is indexed.
func WithSearchAsYouTypeIndex(v bool) SearchAsYouTypePropertyOption {
	return func(p *SearchAsYouTypeProperty) { p.Index = &v }
}

// WithSearchAsYouTypeStore sets whether the field value is stored.
func WithSearchAsYouTypeStore(v bool) SearchAsYouTypePropertyOption {
	return func(p *SearchAsYouTypeProperty) { p.Store = &v }
}

// WithSearchAsYouTypeNorms sets whether norms are enabled.
func WithSearchAsYouTypeNorms(v bool) SearchAsYouTypePropertyOption {
	return func(p *SearchAsYouTypeProperty) { p.Norms = &v }
}

// WithSearchAsYouTypeSimilarity sets the similarity algorithm.
func WithSearchAsYouTypeSimilarity(v string) SearchAsYouTypePropertyOption {
	return func(p *SearchAsYouTypeProperty) { p.Similarity = &v }
}

// ---------------------------------------------------------------------------
// Boolean
// ---------------------------------------------------------------------------

// BooleanPropertyOption is a functional option for configuring a [BooleanProperty].
type BooleanPropertyOption func(*BooleanProperty)

// BooleanProperty represents an Elasticsearch "boolean" field mapping.
// Use [NewBooleanProperty] to construct one with functional options.
type BooleanProperty struct {
	// NullValue is the value substituted for any explicit null values.
	NullValue *bool
	// DocValues controls whether doc values are enabled.
	DocValues *bool
	// Index controls whether the field is indexed.
	Index *bool
	// Store controls whether the field value is stored separately.
	Store *bool
}

// ESTypeName returns the Elasticsearch type name for a boolean property.
func (BooleanProperty) ESTypeName() string { return "boolean" }

// NewBooleanProperty creates a new [BooleanProperty] with the given options applied.
func NewBooleanProperty(opts ...BooleanPropertyOption) BooleanProperty {
	var p BooleanProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithBooleanNullValue sets the null value for the boolean field.
func WithBooleanNullValue(v bool) BooleanPropertyOption {
	return func(p *BooleanProperty) { p.NullValue = &v }
}

// WithBooleanDocValues sets whether doc values are enabled.
func WithBooleanDocValues(v bool) BooleanPropertyOption {
	return func(p *BooleanProperty) { p.DocValues = &v }
}

// WithBooleanIndex sets whether the field is indexed.
func WithBooleanIndex(v bool) BooleanPropertyOption {
	return func(p *BooleanProperty) { p.Index = &v }
}

// WithBooleanStore sets whether the field value is stored.
func WithBooleanStore(v bool) BooleanPropertyOption {
	return func(p *BooleanProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Numeric properties
// ---------------------------------------------------------------------------

// IntegerNumberPropertyOption is a functional option for configuring an [IntegerNumberProperty].
type IntegerNumberPropertyOption func(*IntegerNumberProperty)

// IntegerNumberProperty represents an Elasticsearch "integer" field mapping.
// Use [NewIntegerNumberProperty] to construct one with functional options.
type IntegerNumberProperty struct {
	Coerce          *bool
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *int
}

// ESTypeName returns the Elasticsearch type name for an integer property.
func (IntegerNumberProperty) ESTypeName() string { return "integer" }

// NewIntegerNumberProperty creates a new [IntegerNumberProperty] with the given options applied.
func NewIntegerNumberProperty(opts ...IntegerNumberPropertyOption) IntegerNumberProperty {
	var p IntegerNumberProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithIntegerNumberCoerce sets whether to coerce values to the correct type.
func WithIntegerNumberCoerce(v bool) IntegerNumberPropertyOption {
	return func(p *IntegerNumberProperty) { p.Coerce = &v }
}

// WithIntegerNumberDocValues sets whether doc values are enabled.
func WithIntegerNumberDocValues(v bool) IntegerNumberPropertyOption {
	return func(p *IntegerNumberProperty) { p.DocValues = &v }
}

// WithIntegerNumberIgnoreMalformed sets whether to ignore malformed values.
func WithIntegerNumberIgnoreMalformed(v bool) IntegerNumberPropertyOption {
	return func(p *IntegerNumberProperty) { p.IgnoreMalformed = &v }
}

// WithIntegerNumberIndex sets whether the field is indexed.
func WithIntegerNumberIndex(v bool) IntegerNumberPropertyOption {
	return func(p *IntegerNumberProperty) { p.Index = &v }
}

// WithIntegerNumberStore sets whether the field value is stored.
func WithIntegerNumberStore(v bool) IntegerNumberPropertyOption {
	return func(p *IntegerNumberProperty) { p.Store = &v }
}

// WithIntegerNumberNullValue sets the null value.
func WithIntegerNumberNullValue(v int) IntegerNumberPropertyOption {
	return func(p *IntegerNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------

// LongNumberPropertyOption is a functional option for configuring a [LongNumberProperty].
type LongNumberPropertyOption func(*LongNumberProperty)

// LongNumberProperty represents an Elasticsearch "long" field mapping.
// Use [NewLongNumberProperty] to construct one with functional options.
type LongNumberProperty struct {
	Coerce          *bool
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *int64
}

// ESTypeName returns the Elasticsearch type name for a long property.
func (LongNumberProperty) ESTypeName() string { return "long" }

// NewLongNumberProperty creates a new [LongNumberProperty] with the given options applied.
func NewLongNumberProperty(opts ...LongNumberPropertyOption) LongNumberProperty {
	var p LongNumberProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithLongNumberCoerce sets whether to coerce values to the correct type.
func WithLongNumberCoerce(v bool) LongNumberPropertyOption {
	return func(p *LongNumberProperty) { p.Coerce = &v }
}

// WithLongNumberDocValues sets whether doc values are enabled.
func WithLongNumberDocValues(v bool) LongNumberPropertyOption {
	return func(p *LongNumberProperty) { p.DocValues = &v }
}

// WithLongNumberIgnoreMalformed sets whether to ignore malformed values.
func WithLongNumberIgnoreMalformed(v bool) LongNumberPropertyOption {
	return func(p *LongNumberProperty) { p.IgnoreMalformed = &v }
}

// WithLongNumberIndex sets whether the field is indexed.
func WithLongNumberIndex(v bool) LongNumberPropertyOption {
	return func(p *LongNumberProperty) { p.Index = &v }
}

// WithLongNumberStore sets whether the field value is stored.
func WithLongNumberStore(v bool) LongNumberPropertyOption {
	return func(p *LongNumberProperty) { p.Store = &v }
}

// WithLongNumberNullValue sets the null value.
func WithLongNumberNullValue(v int64) LongNumberPropertyOption {
	return func(p *LongNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------

// ShortNumberPropertyOption is a functional option for configuring a [ShortNumberProperty].
type ShortNumberPropertyOption func(*ShortNumberProperty)

// ShortNumberProperty represents an Elasticsearch "short" field mapping.
// Use [NewShortNumberProperty] to construct one with functional options.
type ShortNumberProperty struct {
	Coerce          *bool
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *int
}

// ESTypeName returns the Elasticsearch type name for a short property.
func (ShortNumberProperty) ESTypeName() string { return "short" }

// NewShortNumberProperty creates a new [ShortNumberProperty] with the given options applied.
func NewShortNumberProperty(opts ...ShortNumberPropertyOption) ShortNumberProperty {
	var p ShortNumberProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithShortNumberCoerce sets whether to coerce values to the correct type.
func WithShortNumberCoerce(v bool) ShortNumberPropertyOption {
	return func(p *ShortNumberProperty) { p.Coerce = &v }
}

// WithShortNumberDocValues sets whether doc values are enabled.
func WithShortNumberDocValues(v bool) ShortNumberPropertyOption {
	return func(p *ShortNumberProperty) { p.DocValues = &v }
}

// WithShortNumberIgnoreMalformed sets whether to ignore malformed values.
func WithShortNumberIgnoreMalformed(v bool) ShortNumberPropertyOption {
	return func(p *ShortNumberProperty) { p.IgnoreMalformed = &v }
}

// WithShortNumberIndex sets whether the field is indexed.
func WithShortNumberIndex(v bool) ShortNumberPropertyOption {
	return func(p *ShortNumberProperty) { p.Index = &v }
}

// WithShortNumberStore sets whether the field value is stored.
func WithShortNumberStore(v bool) ShortNumberPropertyOption {
	return func(p *ShortNumberProperty) { p.Store = &v }
}

// WithShortNumberNullValue sets the null value.
func WithShortNumberNullValue(v int) ShortNumberPropertyOption {
	return func(p *ShortNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------

// ByteNumberPropertyOption is a functional option for configuring a [ByteNumberProperty].
type ByteNumberPropertyOption func(*ByteNumberProperty)

// ByteNumberProperty represents an Elasticsearch "byte" field mapping.
// Use [NewByteNumberProperty] to construct one with functional options.
type ByteNumberProperty struct {
	Coerce          *bool
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *byte
}

// ESTypeName returns the Elasticsearch type name for a byte property.
func (ByteNumberProperty) ESTypeName() string { return "byte" }

// NewByteNumberProperty creates a new [ByteNumberProperty] with the given options applied.
func NewByteNumberProperty(opts ...ByteNumberPropertyOption) ByteNumberProperty {
	var p ByteNumberProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithByteNumberCoerce sets whether to coerce values to the correct type.
func WithByteNumberCoerce(v bool) ByteNumberPropertyOption {
	return func(p *ByteNumberProperty) { p.Coerce = &v }
}

// WithByteNumberDocValues sets whether doc values are enabled.
func WithByteNumberDocValues(v bool) ByteNumberPropertyOption {
	return func(p *ByteNumberProperty) { p.DocValues = &v }
}

// WithByteNumberIgnoreMalformed sets whether to ignore malformed values.
func WithByteNumberIgnoreMalformed(v bool) ByteNumberPropertyOption {
	return func(p *ByteNumberProperty) { p.IgnoreMalformed = &v }
}

// WithByteNumberIndex sets whether the field is indexed.
func WithByteNumberIndex(v bool) ByteNumberPropertyOption {
	return func(p *ByteNumberProperty) { p.Index = &v }
}

// WithByteNumberStore sets whether the field value is stored.
func WithByteNumberStore(v bool) ByteNumberPropertyOption {
	return func(p *ByteNumberProperty) { p.Store = &v }
}

// WithByteNumberNullValue sets the null value.
func WithByteNumberNullValue(v byte) ByteNumberPropertyOption {
	return func(p *ByteNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------

// DoubleNumberPropertyOption is a functional option for configuring a [DoubleNumberProperty].
type DoubleNumberPropertyOption func(*DoubleNumberProperty)

// DoubleNumberProperty represents an Elasticsearch "double" field mapping.
// Use [NewDoubleNumberProperty] to construct one with functional options.
type DoubleNumberProperty struct {
	Coerce          *bool
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *float64
}

// ESTypeName returns the Elasticsearch type name for a double property.
func (DoubleNumberProperty) ESTypeName() string { return "double" }

// NewDoubleNumberProperty creates a new [DoubleNumberProperty] with the given options applied.
func NewDoubleNumberProperty(opts ...DoubleNumberPropertyOption) DoubleNumberProperty {
	var p DoubleNumberProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithDoubleNumberCoerce sets whether to coerce values to the correct type.
func WithDoubleNumberCoerce(v bool) DoubleNumberPropertyOption {
	return func(p *DoubleNumberProperty) { p.Coerce = &v }
}

// WithDoubleNumberDocValues sets whether doc values are enabled.
func WithDoubleNumberDocValues(v bool) DoubleNumberPropertyOption {
	return func(p *DoubleNumberProperty) { p.DocValues = &v }
}

// WithDoubleNumberIgnoreMalformed sets whether to ignore malformed values.
func WithDoubleNumberIgnoreMalformed(v bool) DoubleNumberPropertyOption {
	return func(p *DoubleNumberProperty) { p.IgnoreMalformed = &v }
}

// WithDoubleNumberIndex sets whether the field is indexed.
func WithDoubleNumberIndex(v bool) DoubleNumberPropertyOption {
	return func(p *DoubleNumberProperty) { p.Index = &v }
}

// WithDoubleNumberStore sets whether the field value is stored.
func WithDoubleNumberStore(v bool) DoubleNumberPropertyOption {
	return func(p *DoubleNumberProperty) { p.Store = &v }
}

// WithDoubleNumberNullValue sets the null value.
func WithDoubleNumberNullValue(v float64) DoubleNumberPropertyOption {
	return func(p *DoubleNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------

// FloatNumberPropertyOption is a functional option for configuring a [FloatNumberProperty].
type FloatNumberPropertyOption func(*FloatNumberProperty)

// FloatNumberProperty represents an Elasticsearch "float" field mapping.
// Use [NewFloatNumberProperty] to construct one with functional options.
type FloatNumberProperty struct {
	Coerce          *bool
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *float32
}

// ESTypeName returns the Elasticsearch type name for a float property.
func (FloatNumberProperty) ESTypeName() string { return "float" }

// NewFloatNumberProperty creates a new [FloatNumberProperty] with the given options applied.
func NewFloatNumberProperty(opts ...FloatNumberPropertyOption) FloatNumberProperty {
	var p FloatNumberProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithFloatNumberCoerce sets whether to coerce values to the correct type.
func WithFloatNumberCoerce(v bool) FloatNumberPropertyOption {
	return func(p *FloatNumberProperty) { p.Coerce = &v }
}

// WithFloatNumberDocValues sets whether doc values are enabled.
func WithFloatNumberDocValues(v bool) FloatNumberPropertyOption {
	return func(p *FloatNumberProperty) { p.DocValues = &v }
}

// WithFloatNumberIgnoreMalformed sets whether to ignore malformed values.
func WithFloatNumberIgnoreMalformed(v bool) FloatNumberPropertyOption {
	return func(p *FloatNumberProperty) { p.IgnoreMalformed = &v }
}

// WithFloatNumberIndex sets whether the field is indexed.
func WithFloatNumberIndex(v bool) FloatNumberPropertyOption {
	return func(p *FloatNumberProperty) { p.Index = &v }
}

// WithFloatNumberStore sets whether the field value is stored.
func WithFloatNumberStore(v bool) FloatNumberPropertyOption {
	return func(p *FloatNumberProperty) { p.Store = &v }
}

// WithFloatNumberNullValue sets the null value.
func WithFloatNumberNullValue(v float32) FloatNumberPropertyOption {
	return func(p *FloatNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------

// HalfFloatNumberPropertyOption is a functional option for configuring a [HalfFloatNumberProperty].
type HalfFloatNumberPropertyOption func(*HalfFloatNumberProperty)

// HalfFloatNumberProperty represents an Elasticsearch "half_float" field mapping.
// Use [NewHalfFloatNumberProperty] to construct one with functional options.
type HalfFloatNumberProperty struct {
	Coerce          *bool
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *float32
}

// ESTypeName returns the Elasticsearch type name for a half_float property.
func (HalfFloatNumberProperty) ESTypeName() string { return "half_float" }

// NewHalfFloatNumberProperty creates a new [HalfFloatNumberProperty] with the given options applied.
func NewHalfFloatNumberProperty(opts ...HalfFloatNumberPropertyOption) HalfFloatNumberProperty {
	var p HalfFloatNumberProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithHalfFloatNumberCoerce sets whether to coerce values to the correct type.
func WithHalfFloatNumberCoerce(v bool) HalfFloatNumberPropertyOption {
	return func(p *HalfFloatNumberProperty) { p.Coerce = &v }
}

// WithHalfFloatNumberDocValues sets whether doc values are enabled.
func WithHalfFloatNumberDocValues(v bool) HalfFloatNumberPropertyOption {
	return func(p *HalfFloatNumberProperty) { p.DocValues = &v }
}

// WithHalfFloatNumberIgnoreMalformed sets whether to ignore malformed values.
func WithHalfFloatNumberIgnoreMalformed(v bool) HalfFloatNumberPropertyOption {
	return func(p *HalfFloatNumberProperty) { p.IgnoreMalformed = &v }
}

// WithHalfFloatNumberIndex sets whether the field is indexed.
func WithHalfFloatNumberIndex(v bool) HalfFloatNumberPropertyOption {
	return func(p *HalfFloatNumberProperty) { p.Index = &v }
}

// WithHalfFloatNumberStore sets whether the field value is stored.
func WithHalfFloatNumberStore(v bool) HalfFloatNumberPropertyOption {
	return func(p *HalfFloatNumberProperty) { p.Store = &v }
}

// WithHalfFloatNumberNullValue sets the null value.
func WithHalfFloatNumberNullValue(v float32) HalfFloatNumberPropertyOption {
	return func(p *HalfFloatNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------

// UnsignedLongNumberPropertyOption is a functional option for configuring an [UnsignedLongNumberProperty].
type UnsignedLongNumberPropertyOption func(*UnsignedLongNumberProperty)

// UnsignedLongNumberProperty represents an Elasticsearch "unsigned_long" field mapping.
// Use [NewUnsignedLongNumberProperty] to construct one with functional options.
type UnsignedLongNumberProperty struct {
	Coerce          *bool
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *uint64
}

// ESTypeName returns the Elasticsearch type name for an unsigned_long property.
func (UnsignedLongNumberProperty) ESTypeName() string { return "unsigned_long" }

// NewUnsignedLongNumberProperty creates a new [UnsignedLongNumberProperty] with the given options applied.
func NewUnsignedLongNumberProperty(opts ...UnsignedLongNumberPropertyOption) UnsignedLongNumberProperty {
	var p UnsignedLongNumberProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithUnsignedLongNumberCoerce sets whether to coerce values to the correct type.
func WithUnsignedLongNumberCoerce(v bool) UnsignedLongNumberPropertyOption {
	return func(p *UnsignedLongNumberProperty) { p.Coerce = &v }
}

// WithUnsignedLongNumberDocValues sets whether doc values are enabled.
func WithUnsignedLongNumberDocValues(v bool) UnsignedLongNumberPropertyOption {
	return func(p *UnsignedLongNumberProperty) { p.DocValues = &v }
}

// WithUnsignedLongNumberIgnoreMalformed sets whether to ignore malformed values.
func WithUnsignedLongNumberIgnoreMalformed(v bool) UnsignedLongNumberPropertyOption {
	return func(p *UnsignedLongNumberProperty) { p.IgnoreMalformed = &v }
}

// WithUnsignedLongNumberIndex sets whether the field is indexed.
func WithUnsignedLongNumberIndex(v bool) UnsignedLongNumberPropertyOption {
	return func(p *UnsignedLongNumberProperty) { p.Index = &v }
}

// WithUnsignedLongNumberStore sets whether the field value is stored.
func WithUnsignedLongNumberStore(v bool) UnsignedLongNumberPropertyOption {
	return func(p *UnsignedLongNumberProperty) { p.Store = &v }
}

// WithUnsignedLongNumberNullValue sets the null value.
func WithUnsignedLongNumberNullValue(v uint64) UnsignedLongNumberPropertyOption {
	return func(p *UnsignedLongNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------

// ScaledFloatNumberPropertyOption is a functional option for configuring a [ScaledFloatNumberProperty].
type ScaledFloatNumberPropertyOption func(*ScaledFloatNumberProperty)

// ScaledFloatNumberProperty represents an Elasticsearch "scaled_float" field mapping.
// Use [NewScaledFloatNumberProperty] to construct one with functional options.
type ScaledFloatNumberProperty struct {
	// ScalingFactor is the scaling factor to use when encoding values.
	// Values will be multiplied by this factor at index time and rounded
	// to the nearest long value (e.g. 100 for two decimal places).
	ScalingFactor   *float64
	Coerce          *bool
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *float64
}

// ESTypeName returns the Elasticsearch type name for a scaled_float property.
func (ScaledFloatNumberProperty) ESTypeName() string { return "scaled_float" }

// NewScaledFloatNumberProperty creates a new [ScaledFloatNumberProperty] with the given options applied.
func NewScaledFloatNumberProperty(opts ...ScaledFloatNumberPropertyOption) ScaledFloatNumberProperty {
	var p ScaledFloatNumberProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithScalingFactor sets the scaling factor for the scaled_float field.
func WithScalingFactor(v float64) ScaledFloatNumberPropertyOption {
	return func(p *ScaledFloatNumberProperty) { p.ScalingFactor = &v }
}

// WithScaledFloatNumberCoerce sets whether to coerce values to the correct type.
func WithScaledFloatNumberCoerce(v bool) ScaledFloatNumberPropertyOption {
	return func(p *ScaledFloatNumberProperty) { p.Coerce = &v }
}

// WithScaledFloatNumberDocValues sets whether doc values are enabled.
func WithScaledFloatNumberDocValues(v bool) ScaledFloatNumberPropertyOption {
	return func(p *ScaledFloatNumberProperty) { p.DocValues = &v }
}

// WithScaledFloatNumberIgnoreMalformed sets whether to ignore malformed values.
func WithScaledFloatNumberIgnoreMalformed(v bool) ScaledFloatNumberPropertyOption {
	return func(p *ScaledFloatNumberProperty) { p.IgnoreMalformed = &v }
}

// WithScaledFloatNumberIndex sets whether the field is indexed.
func WithScaledFloatNumberIndex(v bool) ScaledFloatNumberPropertyOption {
	return func(p *ScaledFloatNumberProperty) { p.Index = &v }
}

// WithScaledFloatNumberStore sets whether the field value is stored.
func WithScaledFloatNumberStore(v bool) ScaledFloatNumberPropertyOption {
	return func(p *ScaledFloatNumberProperty) { p.Store = &v }
}

// WithScaledFloatNumberNullValue sets the null value.
func WithScaledFloatNumberNullValue(v float64) ScaledFloatNumberPropertyOption {
	return func(p *ScaledFloatNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Date
// ---------------------------------------------------------------------------

// DatePropertyOption is a functional option for configuring a [DateProperty].
type DatePropertyOption func(*DateProperty)

// DateProperty represents an Elasticsearch "date" field mapping.
// Use [NewDateProperty] to construct one with functional options.
type DateProperty struct {
	// Format is the date format(s) that can be parsed.
	// Multiple formats can be separated by "||".
	Format          *string
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	Locale          *string
}

// ESTypeName returns the Elasticsearch type name for a date property.
func (DateProperty) ESTypeName() string { return "date" }

// NewDateProperty creates a new [DateProperty] with the given options applied.
func NewDateProperty(opts ...DatePropertyOption) DateProperty {
	var p DateProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithDateFormat sets the date format(s) that can be parsed.
func WithDateFormat(v string) DatePropertyOption {
	return func(p *DateProperty) { p.Format = &v }
}

// WithDateDocValues sets whether doc values are enabled.
func WithDateDocValues(v bool) DatePropertyOption {
	return func(p *DateProperty) { p.DocValues = &v }
}

// WithDateIgnoreMalformed sets whether to ignore malformed values.
func WithDateIgnoreMalformed(v bool) DatePropertyOption {
	return func(p *DateProperty) { p.IgnoreMalformed = &v }
}

// WithDateIndex sets whether the field is indexed.
func WithDateIndex(v bool) DatePropertyOption {
	return func(p *DateProperty) { p.Index = &v }
}

// WithDateStore sets whether the field value is stored.
func WithDateStore(v bool) DatePropertyOption {
	return func(p *DateProperty) { p.Store = &v }
}

// WithDateLocale sets the locale for parsing dates.
func WithDateLocale(v string) DatePropertyOption {
	return func(p *DateProperty) { p.Locale = &v }
}

// ---------------------------------------------------------------------------

// DateNanosPropertyOption is a functional option for configuring a [DateNanosProperty].
type DateNanosPropertyOption func(*DateNanosProperty)

// DateNanosProperty represents an Elasticsearch "date_nanos" field mapping
// with nanosecond resolution.
// Use [NewDateNanosProperty] to construct one with functional options.
type DateNanosProperty struct {
	// Format is the date format(s) that can be parsed.
	Format          *string
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
}

// ESTypeName returns the Elasticsearch type name for a date_nanos property.
func (DateNanosProperty) ESTypeName() string { return "date_nanos" }

// NewDateNanosProperty creates a new [DateNanosProperty] with the given options applied.
func NewDateNanosProperty(opts ...DateNanosPropertyOption) DateNanosProperty {
	var p DateNanosProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithDateNanosFormat sets the date format(s) that can be parsed.
func WithDateNanosFormat(v string) DateNanosPropertyOption {
	return func(p *DateNanosProperty) { p.Format = &v }
}

// WithDateNanosDocValues sets whether doc values are enabled.
func WithDateNanosDocValues(v bool) DateNanosPropertyOption {
	return func(p *DateNanosProperty) { p.DocValues = &v }
}

// WithDateNanosIgnoreMalformed sets whether to ignore malformed values.
func WithDateNanosIgnoreMalformed(v bool) DateNanosPropertyOption {
	return func(p *DateNanosProperty) { p.IgnoreMalformed = &v }
}

// WithDateNanosIndex sets whether the field is indexed.
func WithDateNanosIndex(v bool) DateNanosPropertyOption {
	return func(p *DateNanosProperty) { p.Index = &v }
}

// WithDateNanosStore sets whether the field value is stored.
func WithDateNanosStore(v bool) DateNanosPropertyOption {
	return func(p *DateNanosProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Geo
// ---------------------------------------------------------------------------

// GeoPointPropertyOption is a functional option for configuring a [GeoPointProperty].
type GeoPointPropertyOption func(*GeoPointProperty)

// GeoPointProperty represents an Elasticsearch "geo_point" field mapping.
// Use [NewGeoPointProperty] to construct one with functional options.
type GeoPointProperty struct {
	IgnoreMalformed *bool
	IgnoreZValue    *bool
	DocValues       *bool
	Index           *bool
	Store           *bool
}

// ESTypeName returns the Elasticsearch type name for a geo_point property.
func (GeoPointProperty) ESTypeName() string { return "geo_point" }

// NewGeoPointProperty creates a new [GeoPointProperty] with the given options applied.
func NewGeoPointProperty(opts ...GeoPointPropertyOption) GeoPointProperty {
	var p GeoPointProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithGeoPointIgnoreMalformed sets whether to ignore malformed values.
func WithGeoPointIgnoreMalformed(v bool) GeoPointPropertyOption {
	return func(p *GeoPointProperty) { p.IgnoreMalformed = &v }
}

// WithGeoPointIgnoreZValue sets whether to ignore z-values.
func WithGeoPointIgnoreZValue(v bool) GeoPointPropertyOption {
	return func(p *GeoPointProperty) { p.IgnoreZValue = &v }
}

// WithGeoPointDocValues sets whether doc values are enabled.
func WithGeoPointDocValues(v bool) GeoPointPropertyOption {
	return func(p *GeoPointProperty) { p.DocValues = &v }
}

// WithGeoPointIndex sets whether the field is indexed.
func WithGeoPointIndex(v bool) GeoPointPropertyOption {
	return func(p *GeoPointProperty) { p.Index = &v }
}

// WithGeoPointStore sets whether the field value is stored.
func WithGeoPointStore(v bool) GeoPointPropertyOption {
	return func(p *GeoPointProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// GeoShapePropertyOption is a functional option for configuring a [GeoShapeProperty].
type GeoShapePropertyOption func(*GeoShapeProperty)

// GeoShapeProperty represents an Elasticsearch "geo_shape" field mapping.
// Use [NewGeoShapeProperty] to construct one with functional options.
type GeoShapeProperty struct {
	Coerce          *bool
	IgnoreMalformed *bool
	IgnoreZValue    *bool
	DocValues       *bool
	Index           *bool
	Store           *bool
}

// ESTypeName returns the Elasticsearch type name for a geo_shape property.
func (GeoShapeProperty) ESTypeName() string { return "geo_shape" }

// NewGeoShapeProperty creates a new [GeoShapeProperty] with the given options applied.
func NewGeoShapeProperty(opts ...GeoShapePropertyOption) GeoShapeProperty {
	var p GeoShapeProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithGeoShapeCoerce sets whether to coerce values.
func WithGeoShapeCoerce(v bool) GeoShapePropertyOption {
	return func(p *GeoShapeProperty) { p.Coerce = &v }
}

// WithGeoShapeIgnoreMalformed sets whether to ignore malformed values.
func WithGeoShapeIgnoreMalformed(v bool) GeoShapePropertyOption {
	return func(p *GeoShapeProperty) { p.IgnoreMalformed = &v }
}

// WithGeoShapeIgnoreZValue sets whether to ignore z-values.
func WithGeoShapeIgnoreZValue(v bool) GeoShapePropertyOption {
	return func(p *GeoShapeProperty) { p.IgnoreZValue = &v }
}

// WithGeoShapeDocValues sets whether doc values are enabled.
func WithGeoShapeDocValues(v bool) GeoShapePropertyOption {
	return func(p *GeoShapeProperty) { p.DocValues = &v }
}

// WithGeoShapeIndex sets whether the field is indexed.
func WithGeoShapeIndex(v bool) GeoShapePropertyOption {
	return func(p *GeoShapeProperty) { p.Index = &v }
}

// WithGeoShapeStore sets whether the field value is stored.
func WithGeoShapeStore(v bool) GeoShapePropertyOption {
	return func(p *GeoShapeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// ShapePropertyOption is a functional option for configuring a [ShapeProperty].
type ShapePropertyOption func(*ShapeProperty)

// ShapeProperty represents an Elasticsearch "shape" field mapping
// for arbitrary cartesian geometries.
// Use [NewShapeProperty] to construct one with functional options.
type ShapeProperty struct {
	Coerce          *bool
	IgnoreMalformed *bool
	IgnoreZValue    *bool
	DocValues       *bool
	Store           *bool
}

// ESTypeName returns the Elasticsearch type name for a shape property.
func (ShapeProperty) ESTypeName() string { return "shape" }

// NewShapeProperty creates a new [ShapeProperty] with the given options applied.
func NewShapeProperty(opts ...ShapePropertyOption) ShapeProperty {
	var p ShapeProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithShapeCoerce sets whether to coerce values.
func WithShapeCoerce(v bool) ShapePropertyOption {
	return func(p *ShapeProperty) { p.Coerce = &v }
}

// WithShapeIgnoreMalformed sets whether to ignore malformed values.
func WithShapeIgnoreMalformed(v bool) ShapePropertyOption {
	return func(p *ShapeProperty) { p.IgnoreMalformed = &v }
}

// WithShapeIgnoreZValue sets whether to ignore z-values.
func WithShapeIgnoreZValue(v bool) ShapePropertyOption {
	return func(p *ShapeProperty) { p.IgnoreZValue = &v }
}

// WithShapeDocValues sets whether doc values are enabled.
func WithShapeDocValues(v bool) ShapePropertyOption {
	return func(p *ShapeProperty) { p.DocValues = &v }
}

// WithShapeStore sets whether the field value is stored.
func WithShapeStore(v bool) ShapePropertyOption {
	return func(p *ShapeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// PointPropertyOption is a functional option for configuring a [PointProperty].
type PointPropertyOption func(*PointProperty)

// PointProperty represents an Elasticsearch "point" field mapping
// for arbitrary cartesian points.
// Use [NewPointProperty] to construct one with functional options.
type PointProperty struct {
	IgnoreMalformed *bool
	IgnoreZValue    *bool
	DocValues       *bool
	Store           *bool
	NullValue       *string
}

// ESTypeName returns the Elasticsearch type name for a point property.
func (PointProperty) ESTypeName() string { return "point" }

// NewPointProperty creates a new [PointProperty] with the given options applied.
func NewPointProperty(opts ...PointPropertyOption) PointProperty {
	var p PointProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithPointIgnoreMalformed sets whether to ignore malformed values.
func WithPointIgnoreMalformed(v bool) PointPropertyOption {
	return func(p *PointProperty) { p.IgnoreMalformed = &v }
}

// WithPointIgnoreZValue sets whether to ignore z-values.
func WithPointIgnoreZValue(v bool) PointPropertyOption {
	return func(p *PointProperty) { p.IgnoreZValue = &v }
}

// WithPointDocValues sets whether doc values are enabled.
func WithPointDocValues(v bool) PointPropertyOption {
	return func(p *PointProperty) { p.DocValues = &v }
}

// WithPointStore sets whether the field value is stored.
func WithPointStore(v bool) PointPropertyOption {
	return func(p *PointProperty) { p.Store = &v }
}

// WithPointNullValue sets the null value.
func WithPointNullValue(v string) PointPropertyOption {
	return func(p *PointProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Range
// ---------------------------------------------------------------------------

// IntegerRangePropertyOption is a functional option for configuring an [IntegerRangeProperty].
type IntegerRangePropertyOption func(*IntegerRangeProperty)

// IntegerRangeProperty represents an Elasticsearch "integer_range" field mapping.
// Use [NewIntegerRangeProperty] to construct one with functional options.
type IntegerRangeProperty struct {
	Coerce    *bool
	DocValues *bool
	Index     *bool
	Store     *bool
}

// ESTypeName returns the Elasticsearch type name for an integer_range property.
func (IntegerRangeProperty) ESTypeName() string { return "integer_range" }

// NewIntegerRangeProperty creates a new [IntegerRangeProperty] with the given options applied.
func NewIntegerRangeProperty(opts ...IntegerRangePropertyOption) IntegerRangeProperty {
	var p IntegerRangeProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithIntegerRangeCoerce sets whether to coerce values.
func WithIntegerRangeCoerce(v bool) IntegerRangePropertyOption {
	return func(p *IntegerRangeProperty) { p.Coerce = &v }
}

// WithIntegerRangeDocValues sets whether doc values are enabled.
func WithIntegerRangeDocValues(v bool) IntegerRangePropertyOption {
	return func(p *IntegerRangeProperty) { p.DocValues = &v }
}

// WithIntegerRangeIndex sets whether the field is indexed.
func WithIntegerRangeIndex(v bool) IntegerRangePropertyOption {
	return func(p *IntegerRangeProperty) { p.Index = &v }
}

// WithIntegerRangeStore sets whether the field value is stored.
func WithIntegerRangeStore(v bool) IntegerRangePropertyOption {
	return func(p *IntegerRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// LongRangePropertyOption is a functional option for configuring a [LongRangeProperty].
type LongRangePropertyOption func(*LongRangeProperty)

// LongRangeProperty represents an Elasticsearch "long_range" field mapping.
// Use [NewLongRangeProperty] to construct one with functional options.
type LongRangeProperty struct {
	Coerce    *bool
	DocValues *bool
	Index     *bool
	Store     *bool
}

// ESTypeName returns the Elasticsearch type name for a long_range property.
func (LongRangeProperty) ESTypeName() string { return "long_range" }

// NewLongRangeProperty creates a new [LongRangeProperty] with the given options applied.
func NewLongRangeProperty(opts ...LongRangePropertyOption) LongRangeProperty {
	var p LongRangeProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithLongRangeCoerce sets whether to coerce values.
func WithLongRangeCoerce(v bool) LongRangePropertyOption {
	return func(p *LongRangeProperty) { p.Coerce = &v }
}

// WithLongRangeDocValues sets whether doc values are enabled.
func WithLongRangeDocValues(v bool) LongRangePropertyOption {
	return func(p *LongRangeProperty) { p.DocValues = &v }
}

// WithLongRangeIndex sets whether the field is indexed.
func WithLongRangeIndex(v bool) LongRangePropertyOption {
	return func(p *LongRangeProperty) { p.Index = &v }
}

// WithLongRangeStore sets whether the field value is stored.
func WithLongRangeStore(v bool) LongRangePropertyOption {
	return func(p *LongRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// FloatRangePropertyOption is a functional option for configuring a [FloatRangeProperty].
type FloatRangePropertyOption func(*FloatRangeProperty)

// FloatRangeProperty represents an Elasticsearch "float_range" field mapping.
// Use [NewFloatRangeProperty] to construct one with functional options.
type FloatRangeProperty struct {
	Coerce    *bool
	DocValues *bool
	Index     *bool
	Store     *bool
}

// ESTypeName returns the Elasticsearch type name for a float_range property.
func (FloatRangeProperty) ESTypeName() string { return "float_range" }

// NewFloatRangeProperty creates a new [FloatRangeProperty] with the given options applied.
func NewFloatRangeProperty(opts ...FloatRangePropertyOption) FloatRangeProperty {
	var p FloatRangeProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithFloatRangeCoerce sets whether to coerce values.
func WithFloatRangeCoerce(v bool) FloatRangePropertyOption {
	return func(p *FloatRangeProperty) { p.Coerce = &v }
}

// WithFloatRangeDocValues sets whether doc values are enabled.
func WithFloatRangeDocValues(v bool) FloatRangePropertyOption {
	return func(p *FloatRangeProperty) { p.DocValues = &v }
}

// WithFloatRangeIndex sets whether the field is indexed.
func WithFloatRangeIndex(v bool) FloatRangePropertyOption {
	return func(p *FloatRangeProperty) { p.Index = &v }
}

// WithFloatRangeStore sets whether the field value is stored.
func WithFloatRangeStore(v bool) FloatRangePropertyOption {
	return func(p *FloatRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// DoubleRangePropertyOption is a functional option for configuring a [DoubleRangeProperty].
type DoubleRangePropertyOption func(*DoubleRangeProperty)

// DoubleRangeProperty represents an Elasticsearch "double_range" field mapping.
// Use [NewDoubleRangeProperty] to construct one with functional options.
type DoubleRangeProperty struct {
	Coerce    *bool
	DocValues *bool
	Index     *bool
	Store     *bool
}

// ESTypeName returns the Elasticsearch type name for a double_range property.
func (DoubleRangeProperty) ESTypeName() string { return "double_range" }

// NewDoubleRangeProperty creates a new [DoubleRangeProperty] with the given options applied.
func NewDoubleRangeProperty(opts ...DoubleRangePropertyOption) DoubleRangeProperty {
	var p DoubleRangeProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithDoubleRangeCoerce sets whether to coerce values.
func WithDoubleRangeCoerce(v bool) DoubleRangePropertyOption {
	return func(p *DoubleRangeProperty) { p.Coerce = &v }
}

// WithDoubleRangeDocValues sets whether doc values are enabled.
func WithDoubleRangeDocValues(v bool) DoubleRangePropertyOption {
	return func(p *DoubleRangeProperty) { p.DocValues = &v }
}

// WithDoubleRangeIndex sets whether the field is indexed.
func WithDoubleRangeIndex(v bool) DoubleRangePropertyOption {
	return func(p *DoubleRangeProperty) { p.Index = &v }
}

// WithDoubleRangeStore sets whether the field value is stored.
func WithDoubleRangeStore(v bool) DoubleRangePropertyOption {
	return func(p *DoubleRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// DateRangePropertyOption is a functional option for configuring a [DateRangeProperty].
type DateRangePropertyOption func(*DateRangeProperty)

// DateRangeProperty represents an Elasticsearch "date_range" field mapping.
// Use [NewDateRangeProperty] to construct one with functional options.
type DateRangeProperty struct {
	// Format is the date format(s) that can be parsed.
	Format    *string
	Coerce    *bool
	DocValues *bool
	Index     *bool
	Store     *bool
}

// ESTypeName returns the Elasticsearch type name for a date_range property.
func (DateRangeProperty) ESTypeName() string { return "date_range" }

// NewDateRangeProperty creates a new [DateRangeProperty] with the given options applied.
func NewDateRangeProperty(opts ...DateRangePropertyOption) DateRangeProperty {
	var p DateRangeProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithDateRangeFormat sets the date format(s) that can be parsed.
func WithDateRangeFormat(v string) DateRangePropertyOption {
	return func(p *DateRangeProperty) { p.Format = &v }
}

// WithDateRangeCoerce sets whether to coerce values.
func WithDateRangeCoerce(v bool) DateRangePropertyOption {
	return func(p *DateRangeProperty) { p.Coerce = &v }
}

// WithDateRangeDocValues sets whether doc values are enabled.
func WithDateRangeDocValues(v bool) DateRangePropertyOption {
	return func(p *DateRangeProperty) { p.DocValues = &v }
}

// WithDateRangeIndex sets whether the field is indexed.
func WithDateRangeIndex(v bool) DateRangePropertyOption {
	return func(p *DateRangeProperty) { p.Index = &v }
}

// WithDateRangeStore sets whether the field value is stored.
func WithDateRangeStore(v bool) DateRangePropertyOption {
	return func(p *DateRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// IpRangePropertyOption is a functional option for configuring an [IpRangeProperty].
type IpRangePropertyOption func(*IpRangeProperty)

// IpRangeProperty represents an Elasticsearch "ip_range" field mapping.
// Use [NewIpRangeProperty] to construct one with functional options.
type IpRangeProperty struct {
	Coerce    *bool
	DocValues *bool
	Index     *bool
	Store     *bool
}

// ESTypeName returns the Elasticsearch type name for an ip_range property.
func (IpRangeProperty) ESTypeName() string { return "ip_range" }

// NewIpRangeProperty creates a new [IpRangeProperty] with the given options applied.
func NewIpRangeProperty(opts ...IpRangePropertyOption) IpRangeProperty {
	var p IpRangeProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithIpRangeCoerce sets whether to coerce values.
func WithIpRangeCoerce(v bool) IpRangePropertyOption {
	return func(p *IpRangeProperty) { p.Coerce = &v }
}

// WithIpRangeDocValues sets whether doc values are enabled.
func WithIpRangeDocValues(v bool) IpRangePropertyOption {
	return func(p *IpRangeProperty) { p.DocValues = &v }
}

// WithIpRangeIndex sets whether the field is indexed.
func WithIpRangeIndex(v bool) IpRangePropertyOption {
	return func(p *IpRangeProperty) { p.Index = &v }
}

// WithIpRangeStore sets whether the field value is stored.
func WithIpRangeStore(v bool) IpRangePropertyOption {
	return func(p *IpRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Object / Nested
// ---------------------------------------------------------------------------

// ObjectPropertyOption is a functional option for configuring an [ObjectProperty].
type ObjectPropertyOption func(*ObjectProperty)

// ObjectProperty represents an Elasticsearch "object" field mapping.
// Use [NewObjectProperty] to construct one with functional options.
type ObjectProperty struct {
	// Enabled controls whether the JSON object is parsed and indexed (true)
	// or stored as-is without indexing (false).
	Enabled *bool
	// Properties holds the child field mappings.
	Properties map[string]MappingProperty
	// Store controls whether the field value is stored separately.
	Store *bool
}

// ESTypeName returns the Elasticsearch type name for an object property.
func (ObjectProperty) ESTypeName() string { return "object" }

// NewObjectProperty creates a new [ObjectProperty] with the given options applied.
func NewObjectProperty(opts ...ObjectPropertyOption) ObjectProperty {
	var p ObjectProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithObjectEnabled sets whether the object is indexed.
func WithObjectEnabled(v bool) ObjectPropertyOption {
	return func(p *ObjectProperty) { p.Enabled = &v }
}

// WithObjectProperty adds a named child field mapping.
func WithObjectProperty(name string, property MappingProperty) ObjectPropertyOption {
	return func(p *ObjectProperty) {
		if p.Properties == nil {
			p.Properties = make(map[string]MappingProperty)
		}
		p.Properties[name] = property
	}
}

// WithObjectStore sets whether the field value is stored.
func WithObjectStore(v bool) ObjectPropertyOption {
	return func(p *ObjectProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// NestedPropertyOption is a functional option for configuring a [NestedProperty].
type NestedPropertyOption func(*NestedProperty)

// NestedProperty represents an Elasticsearch "nested" field mapping.
// Each nested object is indexed as a separate hidden document, allowing
// independent querying of nested objects.
// Use [NewNestedProperty] to construct one with functional options.
type NestedProperty struct {
	// Properties holds the child field mappings.
	Properties map[string]MappingProperty
	// Enabled controls whether the nested object is parsed and indexed.
	Enabled *bool
	// IncludeInParent controls whether nested fields are added to the parent document.
	IncludeInParent *bool
	// IncludeInRoot controls whether nested fields are added to the root document.
	IncludeInRoot *bool
	// Store controls whether the field value is stored separately.
	Store *bool
}

// ESTypeName returns the Elasticsearch type name for a nested property.
func (NestedProperty) ESTypeName() string { return "nested" }

// NewNestedProperty creates a new [NestedProperty] with the given options applied.
func NewNestedProperty(opts ...NestedPropertyOption) NestedProperty {
	var p NestedProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithNestedProperty adds a named child field mapping.
func WithNestedProperty(name string, property MappingProperty) NestedPropertyOption {
	return func(p *NestedProperty) {
		if p.Properties == nil {
			p.Properties = make(map[string]MappingProperty)
		}
		p.Properties[name] = property
	}
}

// WithNestedEnabled sets whether the nested mapping is enabled.
func WithNestedEnabled(v bool) NestedPropertyOption {
	return func(p *NestedProperty) { p.Enabled = &v }
}

// WithNestedIncludeInParent sets whether to include in parent.
func WithNestedIncludeInParent(v bool) NestedPropertyOption {
	return func(p *NestedProperty) { p.IncludeInParent = &v }
}

// WithNestedIncludeInRoot sets whether to include in root.
func WithNestedIncludeInRoot(v bool) NestedPropertyOption {
	return func(p *NestedProperty) { p.IncludeInRoot = &v }
}

// WithNestedStore sets whether the field value is stored.
func WithNestedStore(v bool) NestedPropertyOption {
	return func(p *NestedProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------

// FlattenedPropertyOption is a functional option for configuring a [FlattenedProperty].
type FlattenedPropertyOption func(*FlattenedProperty)

// FlattenedProperty represents an Elasticsearch "flattened" field mapping.
// An entire object is mapped as a single field, useful for objects with
// a large or unknown number of unique keys.
// Use [NewFlattenedProperty] to construct one with functional options.
type FlattenedProperty struct {
	// DepthLimit is the maximum depth of nested inner objects.
	DepthLimit *int
	// IgnoreAbove is the maximum string length for leaf values.
	IgnoreAbove *int
	// DocValues controls whether doc values are enabled.
	DocValues *bool
	// Index controls whether the field is indexed.
	Index *bool
	// NullValue is the value substituted for explicit null values.
	NullValue *string
	// EagerGlobalOrdinals controls whether global ordinals are loaded eagerly.
	EagerGlobalOrdinals *bool
	// Similarity is the similarity algorithm to use.
	Similarity *string
	// SplitQueriesOnWhitespace controls whether queries are split on whitespace.
	SplitQueriesOnWhitespace *bool
}

// ESTypeName returns the Elasticsearch type name for a flattened property.
func (FlattenedProperty) ESTypeName() string { return "flattened" }

// NewFlattenedProperty creates a new [FlattenedProperty] with the given options applied.
func NewFlattenedProperty(opts ...FlattenedPropertyOption) FlattenedProperty {
	var p FlattenedProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithFlattenedDepthLimit sets the maximum depth of nested inner objects.
func WithFlattenedDepthLimit(v int) FlattenedPropertyOption {
	return func(p *FlattenedProperty) { p.DepthLimit = &v }
}

// WithFlattenedIgnoreAbove sets the maximum string length for leaf values.
func WithFlattenedIgnoreAbove(v int) FlattenedPropertyOption {
	return func(p *FlattenedProperty) { p.IgnoreAbove = &v }
}

// WithFlattenedDocValues sets whether doc values are enabled.
func WithFlattenedDocValues(v bool) FlattenedPropertyOption {
	return func(p *FlattenedProperty) { p.DocValues = &v }
}

// WithFlattenedIndex sets whether the field is indexed.
func WithFlattenedIndex(v bool) FlattenedPropertyOption {
	return func(p *FlattenedProperty) { p.Index = &v }
}

// WithFlattenedNullValue sets the null value.
func WithFlattenedNullValue(v string) FlattenedPropertyOption {
	return func(p *FlattenedProperty) { p.NullValue = &v }
}

// WithFlattenedEagerGlobalOrdinals sets whether to eagerly load global ordinals.
func WithFlattenedEagerGlobalOrdinals(v bool) FlattenedPropertyOption {
	return func(p *FlattenedProperty) { p.EagerGlobalOrdinals = &v }
}

// WithFlattenedSimilarity sets the similarity algorithm.
func WithFlattenedSimilarity(v string) FlattenedPropertyOption {
	return func(p *FlattenedProperty) { p.Similarity = &v }
}

// WithFlattenedSplitQueriesOnWhitespace sets whether to split queries on whitespace.
func WithFlattenedSplitQueriesOnWhitespace(v bool) FlattenedPropertyOption {
	return func(p *FlattenedProperty) { p.SplitQueriesOnWhitespace = &v }
}

// ---------------------------------------------------------------------------

// JoinPropertyOption is a functional option for configuring a [JoinProperty].
type JoinPropertyOption func(*JoinProperty)

// JoinProperty represents an Elasticsearch "join" field mapping that defines
// parent/child relationships within a single index.
// Use [NewJoinProperty] to construct one with functional options.
type JoinProperty struct {
	// Relations maps parent names to their child names.
	Relations map[string][]string
	// EagerGlobalOrdinals controls whether global ordinals are loaded eagerly.
	EagerGlobalOrdinals *bool
}

// ESTypeName returns the Elasticsearch type name for a join property.
func (JoinProperty) ESTypeName() string { return "join" }

// NewJoinProperty creates a new [JoinProperty] with the given options applied.
func NewJoinProperty(opts ...JoinPropertyOption) JoinProperty {
	var p JoinProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithJoinRelation adds a parent-children relationship.
func WithJoinRelation(parent string, children ...string) JoinPropertyOption {
	return func(p *JoinProperty) {
		if p.Relations == nil {
			p.Relations = make(map[string][]string)
		}
		p.Relations[parent] = children
	}
}

// WithJoinEagerGlobalOrdinals sets whether to eagerly load global ordinals.
func WithJoinEagerGlobalOrdinals(v bool) JoinPropertyOption {
	return func(p *JoinProperty) { p.EagerGlobalOrdinals = &v }
}

// ---------------------------------------------------------------------------

// PassthroughObjectPropertyOption is a functional option for configuring a [PassthroughObjectProperty].
type PassthroughObjectPropertyOption func(*PassthroughObjectProperty)

// PassthroughObjectProperty represents an Elasticsearch "passthrough" object field mapping.
// Use [NewPassthroughObjectProperty] to construct one with functional options.
type PassthroughObjectProperty struct {
	// Properties holds the child field mappings.
	Properties map[string]MappingProperty
	// Enabled controls whether the object is parsed and indexed.
	Enabled *bool
	// Priority controls the priority of the passthrough object.
	Priority *int
	// Store controls whether the field value is stored separately.
	Store *bool
	// TimeSeriesDimension controls whether this is a time series dimension.
	TimeSeriesDimension *bool
}

// ESTypeName returns the Elasticsearch type name for a passthrough property.
func (PassthroughObjectProperty) ESTypeName() string { return "passthrough" }

// NewPassthroughObjectProperty creates a new [PassthroughObjectProperty] with the given options applied.
func NewPassthroughObjectProperty(opts ...PassthroughObjectPropertyOption) PassthroughObjectProperty {
	var p PassthroughObjectProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithPassthroughObjectProperty adds a named child field mapping.
func WithPassthroughObjectProperty(name string, property MappingProperty) PassthroughObjectPropertyOption {
	return func(p *PassthroughObjectProperty) {
		if p.Properties == nil {
			p.Properties = make(map[string]MappingProperty)
		}
		p.Properties[name] = property
	}
}

// WithPassthroughObjectEnabled sets whether the object is enabled.
func WithPassthroughObjectEnabled(v bool) PassthroughObjectPropertyOption {
	return func(p *PassthroughObjectProperty) { p.Enabled = &v }
}

// WithPassthroughObjectPriority sets the priority.
func WithPassthroughObjectPriority(v int) PassthroughObjectPropertyOption {
	return func(p *PassthroughObjectProperty) { p.Priority = &v }
}

// WithPassthroughObjectStore sets whether the field value is stored.
func WithPassthroughObjectStore(v bool) PassthroughObjectPropertyOption {
	return func(p *PassthroughObjectProperty) { p.Store = &v }
}

// WithPassthroughObjectTimeSeriesDimension sets whether this is a time series dimension.
func WithPassthroughObjectTimeSeriesDimension(v bool) PassthroughObjectPropertyOption {
	return func(p *PassthroughObjectProperty) { p.TimeSeriesDimension = &v }
}

// ---------------------------------------------------------------------------
// IP
// ---------------------------------------------------------------------------

// IpPropertyOption is a functional option for configuring an [IpProperty].
type IpPropertyOption func(*IpProperty)

// IpProperty represents an Elasticsearch "ip" field mapping for IPv4 and IPv6 addresses.
// Use [NewIpProperty] to construct one with functional options.
type IpProperty struct {
	DocValues       *bool
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	NullValue       *string
}

// ESTypeName returns the Elasticsearch type name for an ip property.
func (IpProperty) ESTypeName() string { return "ip" }

// NewIpProperty creates a new [IpProperty] with the given options applied.
func NewIpProperty(opts ...IpPropertyOption) IpProperty {
	var p IpProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithIpDocValues sets whether doc values are enabled.
func WithIpDocValues(v bool) IpPropertyOption {
	return func(p *IpProperty) { p.DocValues = &v }
}

// WithIpIgnoreMalformed sets whether to ignore malformed values.
func WithIpIgnoreMalformed(v bool) IpPropertyOption {
	return func(p *IpProperty) { p.IgnoreMalformed = &v }
}

// WithIpIndex sets whether the field is indexed.
func WithIpIndex(v bool) IpPropertyOption {
	return func(p *IpProperty) { p.Index = &v }
}

// WithIpStore sets whether the field value is stored.
func WithIpStore(v bool) IpPropertyOption {
	return func(p *IpProperty) { p.Store = &v }
}

// WithIpNullValue sets the null value.
func WithIpNullValue(v string) IpPropertyOption {
	return func(p *IpProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Binary
// ---------------------------------------------------------------------------

// BinaryPropertyOption is a functional option for configuring a [BinaryProperty].
type BinaryPropertyOption func(*BinaryProperty)

// BinaryProperty represents an Elasticsearch "binary" field mapping.
// Values are stored as Base64 encoded strings and are not searchable by default.
// Use [NewBinaryProperty] to construct one with functional options.
type BinaryProperty struct {
	DocValues *bool
	Store     *bool
}

// ESTypeName returns the Elasticsearch type name for a binary property.
func (BinaryProperty) ESTypeName() string { return "binary" }

// NewBinaryProperty creates a new [BinaryProperty] with the given options applied.
func NewBinaryProperty(opts ...BinaryPropertyOption) BinaryProperty {
	var p BinaryProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithBinaryDocValues sets whether doc values are enabled.
func WithBinaryDocValues(v bool) BinaryPropertyOption {
	return func(p *BinaryProperty) { p.DocValues = &v }
}

// WithBinaryStore sets whether the field value is stored.
func WithBinaryStore(v bool) BinaryPropertyOption {
	return func(p *BinaryProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Token Count
// ---------------------------------------------------------------------------

// TokenCountPropertyOption is a functional option for configuring a [TokenCountProperty].
type TokenCountPropertyOption func(*TokenCountProperty)

// TokenCountProperty represents an Elasticsearch "token_count" field mapping.
// It counts the number of tokens produced by an analyzer for a string value.
// Use [NewTokenCountProperty] to construct one with functional options.
type TokenCountProperty struct {
	// Analyzer is the analyzer used to count tokens.
	Analyzer                 *Analyzer
	DocValues                *bool
	Index                    *bool
	Store                    *bool
	EnablePositionIncrements *bool
}

// ESTypeName returns the Elasticsearch type name for a token_count property.
func (TokenCountProperty) ESTypeName() string { return "token_count" }

// NewTokenCountProperty creates a new [TokenCountProperty] with the given options applied.
func NewTokenCountProperty(opts ...TokenCountPropertyOption) TokenCountProperty {
	var p TokenCountProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithTokenCountAnalyzer sets the analyzer used to count tokens.
func WithTokenCountAnalyzer(a Analyzer) TokenCountPropertyOption {
	return func(p *TokenCountProperty) { p.Analyzer = &a }
}

// WithTokenCountDocValues sets whether doc values are enabled.
func WithTokenCountDocValues(v bool) TokenCountPropertyOption {
	return func(p *TokenCountProperty) { p.DocValues = &v }
}

// WithTokenCountIndex sets whether the field is indexed.
func WithTokenCountIndex(v bool) TokenCountPropertyOption {
	return func(p *TokenCountProperty) { p.Index = &v }
}

// WithTokenCountStore sets whether the field value is stored.
func WithTokenCountStore(v bool) TokenCountPropertyOption {
	return func(p *TokenCountProperty) { p.Store = &v }
}

// WithTokenCountEnablePositionIncrements sets whether to count position increments.
func WithTokenCountEnablePositionIncrements(v bool) TokenCountPropertyOption {
	return func(p *TokenCountProperty) { p.EnablePositionIncrements = &v }
}

// ---------------------------------------------------------------------------
// Percolator
// ---------------------------------------------------------------------------

// PercolatorProperty represents an Elasticsearch "percolator" field mapping
// that stores query DSL for the percolate query.
type PercolatorProperty struct{}

// ESTypeName returns the Elasticsearch type name for a percolator property.
func (PercolatorProperty) ESTypeName() string { return "percolator" }

// NewPercolatorProperty creates a new [PercolatorProperty].
func NewPercolatorProperty() PercolatorProperty {
	return PercolatorProperty{}
}

// ---------------------------------------------------------------------------
// Field Alias
// ---------------------------------------------------------------------------

// FieldAliasPropertyOption is a functional option for configuring a [FieldAliasProperty].
type FieldAliasPropertyOption func(*FieldAliasProperty)

// FieldAliasProperty represents an Elasticsearch "alias" field mapping
// that provides an alternate name for a field.
// Use [NewFieldAliasProperty] to construct one with functional options.
type FieldAliasProperty struct {
	// Path is the target field that this alias points to.
	Path *string
}

// ESTypeName returns the Elasticsearch type name for an alias property.
func (FieldAliasProperty) ESTypeName() string { return "alias" }

// NewFieldAliasProperty creates a new [FieldAliasProperty] with the given options applied.
func NewFieldAliasProperty(opts ...FieldAliasPropertyOption) FieldAliasProperty {
	var p FieldAliasProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithFieldAliasPath sets the target field path that this alias points to.
func WithFieldAliasPath(v string) FieldAliasPropertyOption {
	return func(p *FieldAliasProperty) { p.Path = &v }
}

// ---------------------------------------------------------------------------
// Histogram
// ---------------------------------------------------------------------------

// HistogramPropertyOption is a functional option for configuring a [HistogramProperty].
type HistogramPropertyOption func(*HistogramProperty)

// HistogramProperty represents an Elasticsearch "histogram" field mapping
// for pre-aggregated numerical data.
// Use [NewHistogramProperty] to construct one with functional options.
type HistogramProperty struct {
	IgnoreMalformed *bool
}

// ESTypeName returns the Elasticsearch type name for a histogram property.
func (HistogramProperty) ESTypeName() string { return "histogram" }

// NewHistogramProperty creates a new [HistogramProperty] with the given options applied.
func NewHistogramProperty(opts ...HistogramPropertyOption) HistogramProperty {
	var p HistogramProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithHistogramIgnoreMalformed sets whether to ignore malformed values.
func WithHistogramIgnoreMalformed(v bool) HistogramPropertyOption {
	return func(p *HistogramProperty) { p.IgnoreMalformed = &v }
}

// ---------------------------------------------------------------------------
// Version
// ---------------------------------------------------------------------------

// VersionPropertyOption is a functional option for configuring a [VersionProperty].
type VersionPropertyOption func(*VersionProperty)

// VersionProperty represents an Elasticsearch "version" field mapping
// for software version values following semver rules.
// Use [NewVersionProperty] to construct one with functional options.
type VersionProperty struct {
	DocValues *bool
	Store     *bool
}

// ESTypeName returns the Elasticsearch type name for a version property.
func (VersionProperty) ESTypeName() string { return "version" }

// NewVersionProperty creates a new [VersionProperty] with the given options applied.
func NewVersionProperty(opts ...VersionPropertyOption) VersionProperty {
	var p VersionProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithVersionDocValues sets whether doc values are enabled.
func WithVersionDocValues(v bool) VersionPropertyOption {
	return func(p *VersionProperty) { p.DocValues = &v }
}

// WithVersionStore sets whether the field value is stored.
func WithVersionStore(v bool) VersionPropertyOption {
	return func(p *VersionProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Dense Vector
// ---------------------------------------------------------------------------

// DenseVectorPropertyOption is a functional option for configuring a [DenseVectorProperty].
type DenseVectorPropertyOption func(*DenseVectorProperty)

// DenseVectorProperty represents an Elasticsearch "dense_vector" field mapping
// for storing dense vectors of float values for kNN search.
// Use [NewDenseVectorProperty] to construct one with functional options.
type DenseVectorProperty struct {
	// Dims is the number of dimensions in the vector.
	Dims *int
	// Similarity is the similarity metric used for kNN search
	// (e.g. "l2_norm", "dot_product", "cosine").
	Similarity *string
	// Index controls whether the field is indexed for kNN search.
	Index *bool
}

// ESTypeName returns the Elasticsearch type name for a dense_vector property.
func (DenseVectorProperty) ESTypeName() string { return "dense_vector" }

// NewDenseVectorProperty creates a new [DenseVectorProperty] with the given options applied.
func NewDenseVectorProperty(opts ...DenseVectorPropertyOption) DenseVectorProperty {
	var p DenseVectorProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithDenseVectorDims sets the number of dimensions.
func WithDenseVectorDims(v int) DenseVectorPropertyOption {
	return func(p *DenseVectorProperty) { p.Dims = &v }
}

// WithDenseVectorSimilarity sets the similarity metric (e.g. "l2_norm", "dot_product", "cosine").
func WithDenseVectorSimilarity(v string) DenseVectorPropertyOption {
	return func(p *DenseVectorProperty) { p.Similarity = &v }
}

// WithDenseVectorIndex sets whether the field is indexed for kNN search.
func WithDenseVectorIndex(v bool) DenseVectorPropertyOption {
	return func(p *DenseVectorProperty) { p.Index = &v }
}

// ---------------------------------------------------------------------------
// Sparse Vector
// ---------------------------------------------------------------------------

// SparseVectorPropertyOption is a functional option for configuring a [SparseVectorProperty].
type SparseVectorPropertyOption func(*SparseVectorProperty)

// SparseVectorProperty represents an Elasticsearch "sparse_vector" field mapping.
// Use [NewSparseVectorProperty] to construct one with functional options.
type SparseVectorProperty struct {
	Store *bool
}

// ESTypeName returns the Elasticsearch type name for a sparse_vector property.
func (SparseVectorProperty) ESTypeName() string { return "sparse_vector" }

// NewSparseVectorProperty creates a new [SparseVectorProperty] with the given options applied.
func NewSparseVectorProperty(opts ...SparseVectorPropertyOption) SparseVectorProperty {
	var p SparseVectorProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithSparseVectorStore sets whether the field value is stored.
func WithSparseVectorStore(v bool) SparseVectorPropertyOption {
	return func(p *SparseVectorProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Rank Feature / Rank Features
// ---------------------------------------------------------------------------

// RankFeaturePropertyOption is a functional option for configuring a [RankFeatureProperty].
type RankFeaturePropertyOption func(*RankFeatureProperty)

// RankFeatureProperty represents an Elasticsearch "rank_feature" field mapping
// for numeric feature values that boost relevance scoring.
// Use [NewRankFeatureProperty] to construct one with functional options.
type RankFeatureProperty struct {
	// PositiveScoreImpact indicates whether the feature positively
	// correlates with relevance score. Defaults to true.
	PositiveScoreImpact *bool
}

// ESTypeName returns the Elasticsearch type name for a rank_feature property.
func (RankFeatureProperty) ESTypeName() string { return "rank_feature" }

// NewRankFeatureProperty creates a new [RankFeatureProperty] with the given options applied.
func NewRankFeatureProperty(opts ...RankFeaturePropertyOption) RankFeatureProperty {
	var p RankFeatureProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithRankFeaturePositiveScoreImpact sets whether the feature positively
// correlates with the relevance score.
func WithRankFeaturePositiveScoreImpact(v bool) RankFeaturePropertyOption {
	return func(p *RankFeatureProperty) { p.PositiveScoreImpact = &v }
}

// ---------------------------------------------------------------------------

// RankFeaturesPropertyOption is a functional option for configuring a [RankFeaturesProperty].
type RankFeaturesPropertyOption func(*RankFeaturesProperty)

// RankFeaturesProperty represents an Elasticsearch "rank_features" field mapping
// for multiple named rank features in a single field.
// Use [NewRankFeaturesProperty] to construct one with functional options.
type RankFeaturesProperty struct {
	PositiveScoreImpact *bool
}

// ESTypeName returns the Elasticsearch type name for a rank_features property.
func (RankFeaturesProperty) ESTypeName() string { return "rank_features" }

// NewRankFeaturesProperty creates a new [RankFeaturesProperty] with the given options applied.
func NewRankFeaturesProperty(opts ...RankFeaturesPropertyOption) RankFeaturesProperty {
	var p RankFeaturesProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithRankFeaturesPositiveScoreImpact sets whether a positive score has a positive impact.
func WithRankFeaturesPositiveScoreImpact(v bool) RankFeaturesPropertyOption {
	return func(p *RankFeaturesProperty) { p.PositiveScoreImpact = &v }
}

// ---------------------------------------------------------------------------

// RankVectorPropertyOption is a functional option for configuring a [RankVectorProperty].
type RankVectorPropertyOption func(*RankVectorProperty)

// RankVectorProperty represents an Elasticsearch "rank_vectors" field mapping.
// Use [NewRankVectorProperty] to construct one with functional options.
type RankVectorProperty struct {
	Dims *int
}

// ESTypeName returns the Elasticsearch type name for a rank_vectors property.
func (RankVectorProperty) ESTypeName() string { return "rank_vectors" }

// NewRankVectorProperty creates a new [RankVectorProperty] with the given options applied.
func NewRankVectorProperty(opts ...RankVectorPropertyOption) RankVectorProperty {
	var p RankVectorProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithRankVectorDims sets the number of dimensions.
func WithRankVectorDims(v int) RankVectorPropertyOption {
	return func(p *RankVectorProperty) { p.Dims = &v }
}

// ---------------------------------------------------------------------------
// Semantic Text
// ---------------------------------------------------------------------------

// SemanticTextPropertyOption is a functional option for configuring a [SemanticTextProperty].
type SemanticTextPropertyOption func(*SemanticTextProperty)

// SemanticTextProperty represents an Elasticsearch "semantic_text" field mapping
// for text fields that use inference to generate embeddings.
// Use [NewSemanticTextProperty] to construct one with functional options.
type SemanticTextProperty struct {
	// InferenceId is the identifier of the inference endpoint.
	InferenceId *string
	// SearchInferenceId is the identifier of the inference endpoint used at search time.
	SearchInferenceId *string
}

// ESTypeName returns the Elasticsearch type name for a semantic_text property.
func (SemanticTextProperty) ESTypeName() string { return "semantic_text" }

// NewSemanticTextProperty creates a new [SemanticTextProperty] with the given options applied.
func NewSemanticTextProperty(opts ...SemanticTextPropertyOption) SemanticTextProperty {
	var p SemanticTextProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithSemanticTextInferenceId sets the inference endpoint identifier.
func WithSemanticTextInferenceId(v string) SemanticTextPropertyOption {
	return func(p *SemanticTextProperty) { p.InferenceId = &v }
}

// WithSemanticTextSearchInferenceId sets the search inference endpoint identifier.
func WithSemanticTextSearchInferenceId(v string) SemanticTextPropertyOption {
	return func(p *SemanticTextProperty) { p.SearchInferenceId = &v }
}

// ---------------------------------------------------------------------------
// Aggregate Metric Double
// ---------------------------------------------------------------------------

// AggregateMetricDoublePropertyOption is a functional option for configuring an [AggregateMetricDoubleProperty].
type AggregateMetricDoublePropertyOption func(*AggregateMetricDoubleProperty)

// AggregateMetricDoubleProperty represents an Elasticsearch "aggregate_metric_double"
// field mapping for pre-aggregated metric data.
// Use [NewAggregateMetricDoubleProperty] to construct one with functional options.
type AggregateMetricDoubleProperty struct {
	// DefaultMetric is the default metric used in queries, scoring, and sorting.
	DefaultMetric *string
	// Metrics is the list of sub-metric fields (e.g. "min", "max", "sum", "value_count").
	Metrics []string
	// IgnoreMalformed controls whether malformed values are ignored.
	IgnoreMalformed *bool
}

// ESTypeName returns the Elasticsearch type name for an aggregate_metric_double property.
func (AggregateMetricDoubleProperty) ESTypeName() string { return "aggregate_metric_double" }

// NewAggregateMetricDoubleProperty creates a new [AggregateMetricDoubleProperty] with the given options applied.
func NewAggregateMetricDoubleProperty(opts ...AggregateMetricDoublePropertyOption) AggregateMetricDoubleProperty {
	var p AggregateMetricDoubleProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithAggregateMetricDoubleDefaultMetric sets the default metric.
func WithAggregateMetricDoubleDefaultMetric(v string) AggregateMetricDoublePropertyOption {
	return func(p *AggregateMetricDoubleProperty) { p.DefaultMetric = &v }
}

// WithAggregateMetricDoubleMetrics sets the list of sub-metric fields.
func WithAggregateMetricDoubleMetrics(v ...string) AggregateMetricDoublePropertyOption {
	return func(p *AggregateMetricDoubleProperty) { p.Metrics = v }
}

// WithAggregateMetricDoubleIgnoreMalformed sets whether to ignore malformed values.
func WithAggregateMetricDoubleIgnoreMalformed(v bool) AggregateMetricDoublePropertyOption {
	return func(p *AggregateMetricDoubleProperty) { p.IgnoreMalformed = &v }
}

// ---------------------------------------------------------------------------
// Murmur3 Hash
// ---------------------------------------------------------------------------

// Murmur3HashPropertyOption is a functional option for configuring a [Murmur3HashProperty].
type Murmur3HashPropertyOption func(*Murmur3HashProperty)

// Murmur3HashProperty represents an Elasticsearch "murmur3" field mapping
// that computes and stores a murmur3 hash of the field value.
// Use [NewMurmur3HashProperty] to construct one with functional options.
type Murmur3HashProperty struct {
	DocValues *bool
	Store     *bool
}

// ESTypeName returns the Elasticsearch type name for a murmur3 property.
func (Murmur3HashProperty) ESTypeName() string { return "murmur3" }

// NewMurmur3HashProperty creates a new [Murmur3HashProperty] with the given options applied.
func NewMurmur3HashProperty(opts ...Murmur3HashPropertyOption) Murmur3HashProperty {
	var p Murmur3HashProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithMurmur3HashDocValues sets whether doc values are enabled.
func WithMurmur3HashDocValues(v bool) Murmur3HashPropertyOption {
	return func(p *Murmur3HashProperty) { p.DocValues = &v }
}

// WithMurmur3HashStore sets whether the field value is stored.
func WithMurmur3HashStore(v bool) Murmur3HashPropertyOption {
	return func(p *Murmur3HashProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// ICU Collation Keyword
// ---------------------------------------------------------------------------

// IcuCollationPropertyOption is a functional option for configuring an [IcuCollationProperty].
type IcuCollationPropertyOption func(*IcuCollationProperty)

// IcuCollationProperty represents an Elasticsearch "icu_collation_keyword" field mapping
// provided by the ICU analysis plugin. The field indexes text as a keyword
// using ICU collation rules for locale-sensitive sorting and comparison.
// Use [NewIcuCollationProperty] to construct one with functional options.
type IcuCollationProperty struct {
	Language               *string
	Country                *string
	DocValues              *bool
	Index                  *bool
	Store                  *bool
	NullValue              *string
	Norms                  *bool
	Rules                  *string
	Variant                *string
	CaseLevel              *bool
	Numeric                *bool
	HiraganaQuaternaryMode *bool
	VariableTop            *string
}

// ESTypeName returns the Elasticsearch type name for an icu_collation_keyword property.
func (IcuCollationProperty) ESTypeName() string { return "icu_collation_keyword" }

// NewIcuCollationProperty creates a new [IcuCollationProperty] with the given options applied.
func NewIcuCollationProperty(opts ...IcuCollationPropertyOption) IcuCollationProperty {
	var p IcuCollationProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithIcuCollationLanguage sets the language.
func WithIcuCollationLanguage(v string) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.Language = &v }
}

// WithIcuCollationCountry sets the country.
func WithIcuCollationCountry(v string) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.Country = &v }
}

// WithIcuCollationDocValues sets whether doc values are enabled.
func WithIcuCollationDocValues(v bool) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.DocValues = &v }
}

// WithIcuCollationIndex sets whether the field is indexed.
func WithIcuCollationIndex(v bool) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.Index = &v }
}

// WithIcuCollationStore sets whether the field value is stored.
func WithIcuCollationStore(v bool) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.Store = &v }
}

// WithIcuCollationNullValue sets the null value.
func WithIcuCollationNullValue(v string) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.NullValue = &v }
}

// WithIcuCollationNorms sets whether norms are enabled.
func WithIcuCollationNorms(v bool) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.Norms = &v }
}

// WithIcuCollationRules sets the collation rules.
func WithIcuCollationRules(v string) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.Rules = &v }
}

// WithIcuCollationVariant sets the collation variant.
func WithIcuCollationVariant(v string) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.Variant = &v }
}

// WithIcuCollationCaseLevel sets whether case-level sorting is enabled.
func WithIcuCollationCaseLevel(v bool) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.CaseLevel = &v }
}

// WithIcuCollationNumeric sets whether numeric sorting is enabled.
func WithIcuCollationNumeric(v bool) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.Numeric = &v }
}

// WithIcuCollationHiraganaQuaternaryMode sets whether hiragana quaternary mode is enabled.
func WithIcuCollationHiraganaQuaternaryMode(v bool) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.HiraganaQuaternaryMode = &v }
}

// WithIcuCollationVariableTop sets the variable top.
func WithIcuCollationVariableTop(v string) IcuCollationPropertyOption {
	return func(p *IcuCollationProperty) { p.VariableTop = &v }
}

// ---------------------------------------------------------------------------
// Dynamic
// ---------------------------------------------------------------------------

// DynamicPropertyOption is a functional option for configuring a [DynamicProperty].
type DynamicPropertyOption func(*DynamicProperty)

// DynamicProperty represents an Elasticsearch "{dynamic_type}" field mapping
// used in dynamic templates.
// Use [NewDynamicProperty] to construct one with functional options.
type DynamicProperty struct {
	Analyzer        *Analyzer
	SearchAnalyzer  *Analyzer
	Coerce          *bool
	DocValues       *bool
	Enabled         *bool
	Format          *string
	IgnoreMalformed *bool
	Index           *bool
	Store           *bool
	Norms           *bool
	Locale          *string
}

// ESTypeName returns the Elasticsearch type name for a dynamic property.
func (DynamicProperty) ESTypeName() string { return "{dynamic_type}" }

// NewDynamicProperty creates a new [DynamicProperty] with the given options applied.
func NewDynamicProperty(opts ...DynamicPropertyOption) DynamicProperty {
	var p DynamicProperty
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

// WithDynamicAnalyzer sets the analyzer.
func WithDynamicAnalyzer(a Analyzer) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.Analyzer = &a }
}

// WithDynamicSearchAnalyzer sets the search analyzer.
func WithDynamicSearchAnalyzer(a Analyzer) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.SearchAnalyzer = &a }
}

// WithDynamicCoerce sets whether to coerce values.
func WithDynamicCoerce(v bool) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.Coerce = &v }
}

// WithDynamicDocValues sets whether doc values are enabled.
func WithDynamicDocValues(v bool) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.DocValues = &v }
}

// WithDynamicEnabled sets whether the field is enabled.
func WithDynamicEnabled(v bool) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.Enabled = &v }
}

// WithDynamicFormat sets the format.
func WithDynamicFormat(v string) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.Format = &v }
}

// WithDynamicIgnoreMalformed sets whether to ignore malformed values.
func WithDynamicIgnoreMalformed(v bool) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.IgnoreMalformed = &v }
}

// WithDynamicIndex sets whether the field is indexed.
func WithDynamicIndex(v bool) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.Index = &v }
}

// WithDynamicStore sets whether the field value is stored.
func WithDynamicStore(v bool) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.Store = &v }
}

// WithDynamicNorms sets whether norms are enabled.
func WithDynamicNorms(v bool) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.Norms = &v }
}

// WithDynamicLocale sets the locale.
func WithDynamicLocale(v string) DynamicPropertyOption {
	return func(p *DynamicProperty) { p.Locale = &v }
}
