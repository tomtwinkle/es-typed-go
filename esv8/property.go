package esv8

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// ---------------------------------------------------------------------------
// Boolean
// ---------------------------------------------------------------------------

// BooleanPropertyOption is a functional option for configuring BooleanProperty.
type BooleanPropertyOption func(*types.BooleanProperty)

// NewBooleanProperty creates a new boolean property mapping.
func NewBooleanProperty(opts ...BooleanPropertyOption) *types.BooleanProperty {
	prop := types.NewBooleanProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithBooleanDocValues sets whether doc values are enabled.
func WithBooleanDocValues(v bool) BooleanPropertyOption {
	return func(p *types.BooleanProperty) { p.DocValues = &v }
}

// WithBooleanIndex sets whether the field is indexed.
func WithBooleanIndex(v bool) BooleanPropertyOption {
	return func(p *types.BooleanProperty) { p.Index = &v }
}

// WithBooleanStore sets whether the field value is stored.
func WithBooleanStore(v bool) BooleanPropertyOption {
	return func(p *types.BooleanProperty) { p.Store = &v }
}

// WithBooleanNullValue sets the null value for the field.
func WithBooleanNullValue(v bool) BooleanPropertyOption {
	return func(p *types.BooleanProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Keyword
// ---------------------------------------------------------------------------

// KeywordPropertyOption is a functional option for configuring KeywordProperty.
type KeywordPropertyOption func(*types.KeywordProperty)

// NewKeywordProperty creates a new keyword property mapping.
func NewKeywordProperty(opts ...KeywordPropertyOption) *types.KeywordProperty {
	prop := types.NewKeywordProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithKeywordIgnoreAbove sets the maximum string length for the keyword field.
// Strings longer than this value are not indexed.
func WithKeywordIgnoreAbove(v int) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.IgnoreAbove = &v }
}

// WithKeywordDocValues sets whether doc values are enabled.
func WithKeywordDocValues(v bool) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.DocValues = &v }
}

// WithKeywordIndex sets whether the field is indexed.
func WithKeywordIndex(v bool) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.Index = &v }
}

// WithKeywordStore sets whether the field value is stored.
func WithKeywordStore(v bool) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.Store = &v }
}

// WithKeywordNullValue sets the null value for the field.
func WithKeywordNullValue(v string) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.NullValue = &v }
}

// WithKeywordNormalizer sets the normalizer for the keyword field.
func WithKeywordNormalizer(v string) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.Normalizer = &v }
}

// WithKeywordNorms sets whether norms are enabled.
func WithKeywordNorms(v bool) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.Norms = &v }
}

// WithKeywordSimilarity sets the similarity algorithm.
func WithKeywordSimilarity(v string) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.Similarity = &v }
}

// WithKeywordEagerGlobalOrdinals sets whether to eagerly load global ordinals.
func WithKeywordEagerGlobalOrdinals(v bool) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.EagerGlobalOrdinals = &v }
}

// WithKeywordSplitQueriesOnWhitespace sets whether to split queries on whitespace.
func WithKeywordSplitQueriesOnWhitespace(v bool) KeywordPropertyOption {
	return func(p *types.KeywordProperty) { p.SplitQueriesOnWhitespace = &v }
}

// ---------------------------------------------------------------------------
// Constant Keyword
// ---------------------------------------------------------------------------

// ConstantKeywordPropertyOption is a functional option for configuring ConstantKeywordProperty.
type ConstantKeywordPropertyOption func(*types.ConstantKeywordProperty)

// NewConstantKeywordProperty creates a new constant_keyword property mapping.
func NewConstantKeywordProperty(opts ...ConstantKeywordPropertyOption) *types.ConstantKeywordProperty {
	prop := types.NewConstantKeywordProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// ---------------------------------------------------------------------------
// Counted Keyword
// ---------------------------------------------------------------------------

// CountedKeywordPropertyOption is a functional option for configuring CountedKeywordProperty.
type CountedKeywordPropertyOption func(*types.CountedKeywordProperty)

// NewCountedKeywordProperty creates a new counted_keyword property mapping.
func NewCountedKeywordProperty(opts ...CountedKeywordPropertyOption) *types.CountedKeywordProperty {
	prop := types.NewCountedKeywordProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithCountedKeywordIndex sets whether the field is indexed.
func WithCountedKeywordIndex(v bool) CountedKeywordPropertyOption {
	return func(p *types.CountedKeywordProperty) { p.Index = &v }
}

// ---------------------------------------------------------------------------
// Wildcard
// ---------------------------------------------------------------------------

// WildcardPropertyOption is a functional option for configuring WildcardProperty.
type WildcardPropertyOption func(*types.WildcardProperty)

// NewWildcardProperty creates a new wildcard property mapping.
func NewWildcardProperty(opts ...WildcardPropertyOption) *types.WildcardProperty {
	prop := types.NewWildcardProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithWildcardIgnoreAbove sets the maximum string length.
func WithWildcardIgnoreAbove(v int) WildcardPropertyOption {
	return func(p *types.WildcardProperty) { p.IgnoreAbove = &v }
}

// WithWildcardDocValues sets whether doc values are enabled.
func WithWildcardDocValues(v bool) WildcardPropertyOption {
	return func(p *types.WildcardProperty) { p.DocValues = &v }
}

// WithWildcardNullValue sets the null value for the field.
func WithWildcardNullValue(v string) WildcardPropertyOption {
	return func(p *types.WildcardProperty) { p.NullValue = &v }
}

// WithWildcardStore sets whether the field value is stored.
func WithWildcardStore(v bool) WildcardPropertyOption {
	return func(p *types.WildcardProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Text
// ---------------------------------------------------------------------------

// TextPropertyOption is a functional option for configuring TextProperty.
type TextPropertyOption func(*types.TextProperty)

// NewTextProperty creates a new text property mapping.
func NewTextProperty(opts ...TextPropertyOption) *types.TextProperty {
	prop := types.NewTextProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithTextAnalyzer sets the analyzer for the text field.
func WithTextAnalyzer(v string) TextPropertyOption {
	return func(p *types.TextProperty) { p.Analyzer = &v }
}

// WithTextSearchAnalyzer sets the search analyzer for the text field.
func WithTextSearchAnalyzer(v string) TextPropertyOption {
	return func(p *types.TextProperty) { p.SearchAnalyzer = &v }
}

// WithTextSearchQuoteAnalyzer sets the search quote analyzer for the text field.
func WithTextSearchQuoteAnalyzer(v string) TextPropertyOption {
	return func(p *types.TextProperty) { p.SearchQuoteAnalyzer = &v }
}

// WithTextFielddata sets whether fielddata is enabled.
func WithTextFielddata(v bool) TextPropertyOption {
	return func(p *types.TextProperty) { p.Fielddata = &v }
}

// WithTextIndex sets whether the field is indexed.
func WithTextIndex(v bool) TextPropertyOption {
	return func(p *types.TextProperty) { p.Index = &v }
}

// WithTextStore sets whether the field value is stored.
func WithTextStore(v bool) TextPropertyOption {
	return func(p *types.TextProperty) { p.Store = &v }
}

// WithTextNorms sets whether norms are enabled.
func WithTextNorms(v bool) TextPropertyOption {
	return func(p *types.TextProperty) { p.Norms = &v }
}

// WithTextSimilarity sets the similarity algorithm.
func WithTextSimilarity(v string) TextPropertyOption {
	return func(p *types.TextProperty) { p.Similarity = &v }
}

// WithTextIndexPhrases sets whether two-term word combinations are indexed.
func WithTextIndexPhrases(v bool) TextPropertyOption {
	return func(p *types.TextProperty) { p.IndexPhrases = &v }
}

// WithTextPositionIncrementGap sets the number of fake term positions between indexed values.
func WithTextPositionIncrementGap(v int) TextPropertyOption {
	return func(p *types.TextProperty) { p.PositionIncrementGap = &v }
}

// WithTextRawKeyword adds a "keyword" multi-field to the text property
// for exact-match queries alongside full-text search.
// ignoreAbove specifies the maximum string length for the keyword sub-field.
func WithTextRawKeyword(ignoreAbove int) TextPropertyOption {
	return func(p *types.TextProperty) {
		if p.Fields == nil {
			p.Fields = make(map[string]types.Property)
		}
		p.Fields["keyword"] = NewKeywordProperty(WithKeywordIgnoreAbove(ignoreAbove))
	}
}

// WithTextFields sets the multi-fields for the text property.
func WithTextFields(fields map[string]types.Property) TextPropertyOption {
	return func(p *types.TextProperty) { p.Fields = fields }
}

// ---------------------------------------------------------------------------
// Match Only Text
// ---------------------------------------------------------------------------

// MatchOnlyTextPropertyOption is a functional option for configuring MatchOnlyTextProperty.
type MatchOnlyTextPropertyOption func(*types.MatchOnlyTextProperty)

// NewMatchOnlyTextProperty creates a new match_only_text property mapping.
func NewMatchOnlyTextProperty(opts ...MatchOnlyTextPropertyOption) *types.MatchOnlyTextProperty {
	prop := types.NewMatchOnlyTextProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// ---------------------------------------------------------------------------
// Completion
// ---------------------------------------------------------------------------

// CompletionPropertyOption is a functional option for configuring CompletionProperty.
type CompletionPropertyOption func(*types.CompletionProperty)

// NewCompletionProperty creates a new completion property mapping.
func NewCompletionProperty(opts ...CompletionPropertyOption) *types.CompletionProperty {
	prop := types.NewCompletionProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithCompletionAnalyzer sets the analyzer for the completion field.
func WithCompletionAnalyzer(v string) CompletionPropertyOption {
	return func(p *types.CompletionProperty) { p.Analyzer = &v }
}

// WithCompletionSearchAnalyzer sets the search analyzer.
func WithCompletionSearchAnalyzer(v string) CompletionPropertyOption {
	return func(p *types.CompletionProperty) { p.SearchAnalyzer = &v }
}

// WithCompletionMaxInputLength sets the maximum length of a single input.
func WithCompletionMaxInputLength(v int) CompletionPropertyOption {
	return func(p *types.CompletionProperty) { p.MaxInputLength = &v }
}

// WithCompletionPreservePositionIncrements sets whether to preserve position increments.
func WithCompletionPreservePositionIncrements(v bool) CompletionPropertyOption {
	return func(p *types.CompletionProperty) { p.PreservePositionIncrements = &v }
}

// WithCompletionPreserveSeparators sets whether to preserve separators.
func WithCompletionPreserveSeparators(v bool) CompletionPropertyOption {
	return func(p *types.CompletionProperty) { p.PreserveSeparators = &v }
}

// ---------------------------------------------------------------------------
// Search As You Type
// ---------------------------------------------------------------------------

// SearchAsYouTypePropertyOption is a functional option for configuring SearchAsYouTypeProperty.
type SearchAsYouTypePropertyOption func(*types.SearchAsYouTypeProperty)

// NewSearchAsYouTypeProperty creates a new search_as_you_type property mapping.
func NewSearchAsYouTypeProperty(opts ...SearchAsYouTypePropertyOption) *types.SearchAsYouTypeProperty {
	prop := types.NewSearchAsYouTypeProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithSearchAsYouTypeAnalyzer sets the analyzer.
func WithSearchAsYouTypeAnalyzer(v string) SearchAsYouTypePropertyOption {
	return func(p *types.SearchAsYouTypeProperty) { p.Analyzer = &v }
}

// WithSearchAsYouTypeSearchAnalyzer sets the search analyzer.
func WithSearchAsYouTypeSearchAnalyzer(v string) SearchAsYouTypePropertyOption {
	return func(p *types.SearchAsYouTypeProperty) { p.SearchAnalyzer = &v }
}

// WithSearchAsYouTypeSearchQuoteAnalyzer sets the search quote analyzer.
func WithSearchAsYouTypeSearchQuoteAnalyzer(v string) SearchAsYouTypePropertyOption {
	return func(p *types.SearchAsYouTypeProperty) { p.SearchQuoteAnalyzer = &v }
}

// WithSearchAsYouTypeMaxShingleSize sets the max shingle size (2-4).
func WithSearchAsYouTypeMaxShingleSize(v int) SearchAsYouTypePropertyOption {
	return func(p *types.SearchAsYouTypeProperty) { p.MaxShingleSize = &v }
}

// WithSearchAsYouTypeIndex sets whether the field is indexed.
func WithSearchAsYouTypeIndex(v bool) SearchAsYouTypePropertyOption {
	return func(p *types.SearchAsYouTypeProperty) { p.Index = &v }
}

// WithSearchAsYouTypeStore sets whether the field value is stored.
func WithSearchAsYouTypeStore(v bool) SearchAsYouTypePropertyOption {
	return func(p *types.SearchAsYouTypeProperty) { p.Store = &v }
}

// WithSearchAsYouTypeNorms sets whether norms are enabled.
func WithSearchAsYouTypeNorms(v bool) SearchAsYouTypePropertyOption {
	return func(p *types.SearchAsYouTypeProperty) { p.Norms = &v }
}

// WithSearchAsYouTypeSimilarity sets the similarity algorithm.
func WithSearchAsYouTypeSimilarity(v string) SearchAsYouTypePropertyOption {
	return func(p *types.SearchAsYouTypeProperty) { p.Similarity = &v }
}

// ---------------------------------------------------------------------------
// Integer Number
// ---------------------------------------------------------------------------

// IntegerNumberPropertyOption is a functional option for configuring IntegerNumberProperty.
type IntegerNumberPropertyOption func(*types.IntegerNumberProperty)

// NewIntegerNumberProperty creates a new integer number property mapping.
func NewIntegerNumberProperty(opts ...IntegerNumberPropertyOption) *types.IntegerNumberProperty {
	prop := types.NewIntegerNumberProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithIntegerNumberCoerce sets whether to coerce values to the correct type.
func WithIntegerNumberCoerce(v bool) IntegerNumberPropertyOption {
	return func(p *types.IntegerNumberProperty) { p.Coerce = &v }
}

// WithIntegerNumberDocValues sets whether doc values are enabled.
func WithIntegerNumberDocValues(v bool) IntegerNumberPropertyOption {
	return func(p *types.IntegerNumberProperty) { p.DocValues = &v }
}

// WithIntegerNumberIgnoreMalformed sets whether to ignore malformed values.
func WithIntegerNumberIgnoreMalformed(v bool) IntegerNumberPropertyOption {
	return func(p *types.IntegerNumberProperty) { p.IgnoreMalformed = &v }
}

// WithIntegerNumberIndex sets whether the field is indexed.
func WithIntegerNumberIndex(v bool) IntegerNumberPropertyOption {
	return func(p *types.IntegerNumberProperty) { p.Index = &v }
}

// WithIntegerNumberStore sets whether the field value is stored.
func WithIntegerNumberStore(v bool) IntegerNumberPropertyOption {
	return func(p *types.IntegerNumberProperty) { p.Store = &v }
}

// WithIntegerNumberNullValue sets the null value.
func WithIntegerNumberNullValue(v int) IntegerNumberPropertyOption {
	return func(p *types.IntegerNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Long Number
// ---------------------------------------------------------------------------

// LongNumberPropertyOption is a functional option for configuring LongNumberProperty.
type LongNumberPropertyOption func(*types.LongNumberProperty)

// NewLongNumberProperty creates a new long number property mapping.
func NewLongNumberProperty(opts ...LongNumberPropertyOption) *types.LongNumberProperty {
	prop := types.NewLongNumberProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithLongNumberCoerce sets whether to coerce values to the correct type.
func WithLongNumberCoerce(v bool) LongNumberPropertyOption {
	return func(p *types.LongNumberProperty) { p.Coerce = &v }
}

// WithLongNumberDocValues sets whether doc values are enabled.
func WithLongNumberDocValues(v bool) LongNumberPropertyOption {
	return func(p *types.LongNumberProperty) { p.DocValues = &v }
}

// WithLongNumberIgnoreMalformed sets whether to ignore malformed values.
func WithLongNumberIgnoreMalformed(v bool) LongNumberPropertyOption {
	return func(p *types.LongNumberProperty) { p.IgnoreMalformed = &v }
}

// WithLongNumberIndex sets whether the field is indexed.
func WithLongNumberIndex(v bool) LongNumberPropertyOption {
	return func(p *types.LongNumberProperty) { p.Index = &v }
}

// WithLongNumberStore sets whether the field value is stored.
func WithLongNumberStore(v bool) LongNumberPropertyOption {
	return func(p *types.LongNumberProperty) { p.Store = &v }
}

// WithLongNumberNullValue sets the null value.
func WithLongNumberNullValue(v int64) LongNumberPropertyOption {
	return func(p *types.LongNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Short Number
// ---------------------------------------------------------------------------

// ShortNumberPropertyOption is a functional option for configuring ShortNumberProperty.
type ShortNumberPropertyOption func(*types.ShortNumberProperty)

// NewShortNumberProperty creates a new short number property mapping.
func NewShortNumberProperty(opts ...ShortNumberPropertyOption) *types.ShortNumberProperty {
	prop := types.NewShortNumberProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithShortNumberCoerce sets whether to coerce values to the correct type.
func WithShortNumberCoerce(v bool) ShortNumberPropertyOption {
	return func(p *types.ShortNumberProperty) { p.Coerce = &v }
}

// WithShortNumberDocValues sets whether doc values are enabled.
func WithShortNumberDocValues(v bool) ShortNumberPropertyOption {
	return func(p *types.ShortNumberProperty) { p.DocValues = &v }
}

// WithShortNumberIgnoreMalformed sets whether to ignore malformed values.
func WithShortNumberIgnoreMalformed(v bool) ShortNumberPropertyOption {
	return func(p *types.ShortNumberProperty) { p.IgnoreMalformed = &v }
}

// WithShortNumberIndex sets whether the field is indexed.
func WithShortNumberIndex(v bool) ShortNumberPropertyOption {
	return func(p *types.ShortNumberProperty) { p.Index = &v }
}

// WithShortNumberStore sets whether the field value is stored.
func WithShortNumberStore(v bool) ShortNumberPropertyOption {
	return func(p *types.ShortNumberProperty) { p.Store = &v }
}

// WithShortNumberNullValue sets the null value.
func WithShortNumberNullValue(v int) ShortNumberPropertyOption {
	return func(p *types.ShortNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Byte Number
// ---------------------------------------------------------------------------

// ByteNumberPropertyOption is a functional option for configuring ByteNumberProperty.
type ByteNumberPropertyOption func(*types.ByteNumberProperty)

// NewByteNumberProperty creates a new byte number property mapping.
func NewByteNumberProperty(opts ...ByteNumberPropertyOption) *types.ByteNumberProperty {
	prop := types.NewByteNumberProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithByteNumberCoerce sets whether to coerce values to the correct type.
func WithByteNumberCoerce(v bool) ByteNumberPropertyOption {
	return func(p *types.ByteNumberProperty) { p.Coerce = &v }
}

// WithByteNumberDocValues sets whether doc values are enabled.
func WithByteNumberDocValues(v bool) ByteNumberPropertyOption {
	return func(p *types.ByteNumberProperty) { p.DocValues = &v }
}

// WithByteNumberIgnoreMalformed sets whether to ignore malformed values.
func WithByteNumberIgnoreMalformed(v bool) ByteNumberPropertyOption {
	return func(p *types.ByteNumberProperty) { p.IgnoreMalformed = &v }
}

// WithByteNumberIndex sets whether the field is indexed.
func WithByteNumberIndex(v bool) ByteNumberPropertyOption {
	return func(p *types.ByteNumberProperty) { p.Index = &v }
}

// WithByteNumberStore sets whether the field value is stored.
func WithByteNumberStore(v bool) ByteNumberPropertyOption {
	return func(p *types.ByteNumberProperty) { p.Store = &v }
}

// WithByteNumberNullValue sets the null value.
func WithByteNumberNullValue(v byte) ByteNumberPropertyOption {
	return func(p *types.ByteNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Double Number
// ---------------------------------------------------------------------------

// DoubleNumberPropertyOption is a functional option for configuring DoubleNumberProperty.
type DoubleNumberPropertyOption func(*types.DoubleNumberProperty)

// NewDoubleNumberProperty creates a new double number property mapping.
func NewDoubleNumberProperty(opts ...DoubleNumberPropertyOption) *types.DoubleNumberProperty {
	prop := types.NewDoubleNumberProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithDoubleNumberCoerce sets whether to coerce values to the correct type.
func WithDoubleNumberCoerce(v bool) DoubleNumberPropertyOption {
	return func(p *types.DoubleNumberProperty) { p.Coerce = &v }
}

// WithDoubleNumberDocValues sets whether doc values are enabled.
func WithDoubleNumberDocValues(v bool) DoubleNumberPropertyOption {
	return func(p *types.DoubleNumberProperty) { p.DocValues = &v }
}

// WithDoubleNumberIgnoreMalformed sets whether to ignore malformed values.
func WithDoubleNumberIgnoreMalformed(v bool) DoubleNumberPropertyOption {
	return func(p *types.DoubleNumberProperty) { p.IgnoreMalformed = &v }
}

// WithDoubleNumberIndex sets whether the field is indexed.
func WithDoubleNumberIndex(v bool) DoubleNumberPropertyOption {
	return func(p *types.DoubleNumberProperty) { p.Index = &v }
}

// WithDoubleNumberStore sets whether the field value is stored.
func WithDoubleNumberStore(v bool) DoubleNumberPropertyOption {
	return func(p *types.DoubleNumberProperty) { p.Store = &v }
}

// WithDoubleNumberNullValue sets the null value.
func WithDoubleNumberNullValue(v float64) DoubleNumberPropertyOption {
	return func(p *types.DoubleNumberProperty) {
		fv := types.Float64(v)
		p.NullValue = &fv
	}
}

// ---------------------------------------------------------------------------
// Float Number
// ---------------------------------------------------------------------------

// FloatNumberPropertyOption is a functional option for configuring FloatNumberProperty.
type FloatNumberPropertyOption func(*types.FloatNumberProperty)

// NewFloatNumberProperty creates a new float number property mapping.
func NewFloatNumberProperty(opts ...FloatNumberPropertyOption) *types.FloatNumberProperty {
	prop := types.NewFloatNumberProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithFloatNumberCoerce sets whether to coerce values to the correct type.
func WithFloatNumberCoerce(v bool) FloatNumberPropertyOption {
	return func(p *types.FloatNumberProperty) { p.Coerce = &v }
}

// WithFloatNumberDocValues sets whether doc values are enabled.
func WithFloatNumberDocValues(v bool) FloatNumberPropertyOption {
	return func(p *types.FloatNumberProperty) { p.DocValues = &v }
}

// WithFloatNumberIgnoreMalformed sets whether to ignore malformed values.
func WithFloatNumberIgnoreMalformed(v bool) FloatNumberPropertyOption {
	return func(p *types.FloatNumberProperty) { p.IgnoreMalformed = &v }
}

// WithFloatNumberIndex sets whether the field is indexed.
func WithFloatNumberIndex(v bool) FloatNumberPropertyOption {
	return func(p *types.FloatNumberProperty) { p.Index = &v }
}

// WithFloatNumberStore sets whether the field value is stored.
func WithFloatNumberStore(v bool) FloatNumberPropertyOption {
	return func(p *types.FloatNumberProperty) { p.Store = &v }
}

// WithFloatNumberNullValue sets the null value.
func WithFloatNumberNullValue(v float32) FloatNumberPropertyOption {
	return func(p *types.FloatNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Half Float Number
// ---------------------------------------------------------------------------

// HalfFloatNumberPropertyOption is a functional option for configuring HalfFloatNumberProperty.
type HalfFloatNumberPropertyOption func(*types.HalfFloatNumberProperty)

// NewHalfFloatNumberProperty creates a new half_float number property mapping.
func NewHalfFloatNumberProperty(opts ...HalfFloatNumberPropertyOption) *types.HalfFloatNumberProperty {
	prop := types.NewHalfFloatNumberProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithHalfFloatNumberCoerce sets whether to coerce values to the correct type.
func WithHalfFloatNumberCoerce(v bool) HalfFloatNumberPropertyOption {
	return func(p *types.HalfFloatNumberProperty) { p.Coerce = &v }
}

// WithHalfFloatNumberDocValues sets whether doc values are enabled.
func WithHalfFloatNumberDocValues(v bool) HalfFloatNumberPropertyOption {
	return func(p *types.HalfFloatNumberProperty) { p.DocValues = &v }
}

// WithHalfFloatNumberIgnoreMalformed sets whether to ignore malformed values.
func WithHalfFloatNumberIgnoreMalformed(v bool) HalfFloatNumberPropertyOption {
	return func(p *types.HalfFloatNumberProperty) { p.IgnoreMalformed = &v }
}

// WithHalfFloatNumberIndex sets whether the field is indexed.
func WithHalfFloatNumberIndex(v bool) HalfFloatNumberPropertyOption {
	return func(p *types.HalfFloatNumberProperty) { p.Index = &v }
}

// WithHalfFloatNumberStore sets whether the field value is stored.
func WithHalfFloatNumberStore(v bool) HalfFloatNumberPropertyOption {
	return func(p *types.HalfFloatNumberProperty) { p.Store = &v }
}

// WithHalfFloatNumberNullValue sets the null value.
func WithHalfFloatNumberNullValue(v float32) HalfFloatNumberPropertyOption {
	return func(p *types.HalfFloatNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Unsigned Long Number
// ---------------------------------------------------------------------------

// UnsignedLongNumberPropertyOption is a functional option for configuring UnsignedLongNumberProperty.
type UnsignedLongNumberPropertyOption func(*types.UnsignedLongNumberProperty)

// NewUnsignedLongNumberProperty creates a new unsigned_long number property mapping.
func NewUnsignedLongNumberProperty(opts ...UnsignedLongNumberPropertyOption) *types.UnsignedLongNumberProperty {
	prop := types.NewUnsignedLongNumberProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithUnsignedLongNumberCoerce sets whether to coerce values to the correct type.
func WithUnsignedLongNumberCoerce(v bool) UnsignedLongNumberPropertyOption {
	return func(p *types.UnsignedLongNumberProperty) { p.Coerce = &v }
}

// WithUnsignedLongNumberDocValues sets whether doc values are enabled.
func WithUnsignedLongNumberDocValues(v bool) UnsignedLongNumberPropertyOption {
	return func(p *types.UnsignedLongNumberProperty) { p.DocValues = &v }
}

// WithUnsignedLongNumberIgnoreMalformed sets whether to ignore malformed values.
func WithUnsignedLongNumberIgnoreMalformed(v bool) UnsignedLongNumberPropertyOption {
	return func(p *types.UnsignedLongNumberProperty) { p.IgnoreMalformed = &v }
}

// WithUnsignedLongNumberIndex sets whether the field is indexed.
func WithUnsignedLongNumberIndex(v bool) UnsignedLongNumberPropertyOption {
	return func(p *types.UnsignedLongNumberProperty) { p.Index = &v }
}

// WithUnsignedLongNumberStore sets whether the field value is stored.
func WithUnsignedLongNumberStore(v bool) UnsignedLongNumberPropertyOption {
	return func(p *types.UnsignedLongNumberProperty) { p.Store = &v }
}

// WithUnsignedLongNumberNullValue sets the null value.
func WithUnsignedLongNumberNullValue(v uint64) UnsignedLongNumberPropertyOption {
	return func(p *types.UnsignedLongNumberProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Scaled Float Number
// ---------------------------------------------------------------------------

// ScaledFloatNumberPropertyOption is a functional option for configuring ScaledFloatNumberProperty.
type ScaledFloatNumberPropertyOption func(*types.ScaledFloatNumberProperty)

// NewScaledFloatNumberProperty creates a new scaled_float number property mapping.
func NewScaledFloatNumberProperty(opts ...ScaledFloatNumberPropertyOption) *types.ScaledFloatNumberProperty {
	prop := types.NewScaledFloatNumberProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithScaledFloatNumberScalingFactor sets the scaling factor.
func WithScaledFloatNumberScalingFactor(v float64) ScaledFloatNumberPropertyOption {
	return func(p *types.ScaledFloatNumberProperty) {
		fv := types.Float64(v)
		p.ScalingFactor = &fv
	}
}

// WithScaledFloatNumberCoerce sets whether to coerce values to the correct type.
func WithScaledFloatNumberCoerce(v bool) ScaledFloatNumberPropertyOption {
	return func(p *types.ScaledFloatNumberProperty) { p.Coerce = &v }
}

// WithScaledFloatNumberDocValues sets whether doc values are enabled.
func WithScaledFloatNumberDocValues(v bool) ScaledFloatNumberPropertyOption {
	return func(p *types.ScaledFloatNumberProperty) { p.DocValues = &v }
}

// WithScaledFloatNumberIgnoreMalformed sets whether to ignore malformed values.
func WithScaledFloatNumberIgnoreMalformed(v bool) ScaledFloatNumberPropertyOption {
	return func(p *types.ScaledFloatNumberProperty) { p.IgnoreMalformed = &v }
}

// WithScaledFloatNumberIndex sets whether the field is indexed.
func WithScaledFloatNumberIndex(v bool) ScaledFloatNumberPropertyOption {
	return func(p *types.ScaledFloatNumberProperty) { p.Index = &v }
}

// WithScaledFloatNumberStore sets whether the field value is stored.
func WithScaledFloatNumberStore(v bool) ScaledFloatNumberPropertyOption {
	return func(p *types.ScaledFloatNumberProperty) { p.Store = &v }
}

// WithScaledFloatNumberNullValue sets the null value.
func WithScaledFloatNumberNullValue(v float64) ScaledFloatNumberPropertyOption {
	return func(p *types.ScaledFloatNumberProperty) {
		fv := types.Float64(v)
		p.NullValue = &fv
	}
}

// ---------------------------------------------------------------------------
// Date
// ---------------------------------------------------------------------------

// DatePropertyOption is a functional option for configuring DateProperty.
type DatePropertyOption func(*types.DateProperty)

// NewDateProperty creates a new date property mapping.
//
// https://www.elastic.co/docs/reference/elasticsearch/mapping-reference/mapping-date-format
func NewDateProperty(opts ...DatePropertyOption) *types.DateProperty {
	prop := types.NewDateProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithDateFormat sets the date format(s).
// Multiple formats are joined with "||".
func WithDateFormat(formats ...estype.DateFormat) DatePropertyOption {
	return func(p *types.DateProperty) {
		if len(formats) > 0 {
			format := estype.JoinDateFormats(formats...)
			p.Format = &format
		}
	}
}

// WithDateDocValues sets whether doc values are enabled.
func WithDateDocValues(v bool) DatePropertyOption {
	return func(p *types.DateProperty) { p.DocValues = &v }
}

// WithDateIgnoreMalformed sets whether to ignore malformed values.
func WithDateIgnoreMalformed(v bool) DatePropertyOption {
	return func(p *types.DateProperty) { p.IgnoreMalformed = &v }
}

// WithDateIndex sets whether the field is indexed.
func WithDateIndex(v bool) DatePropertyOption {
	return func(p *types.DateProperty) { p.Index = &v }
}

// WithDateStore sets whether the field value is stored.
func WithDateStore(v bool) DatePropertyOption {
	return func(p *types.DateProperty) { p.Store = &v }
}

// WithDateLocale sets the locale for parsing dates.
func WithDateLocale(v string) DatePropertyOption {
	return func(p *types.DateProperty) { p.Locale = &v }
}

// ---------------------------------------------------------------------------
// Date Nanos
// ---------------------------------------------------------------------------

// DateNanosPropertyOption is a functional option for configuring DateNanosProperty.
type DateNanosPropertyOption func(*types.DateNanosProperty)

// NewDateNanosProperty creates a new date_nanos property mapping.
// Date nanos limits its range of dates from roughly 1970 to 2262.
//
// https://www.elastic.co/docs/reference/elasticsearch/mapping-reference/date_nanos
func NewDateNanosProperty(opts ...DateNanosPropertyOption) *types.DateNanosProperty {
	prop := types.NewDateNanosProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithDateNanosFormat sets the date format(s).
func WithDateNanosFormat(formats ...estype.DateFormat) DateNanosPropertyOption {
	return func(p *types.DateNanosProperty) {
		if len(formats) > 0 {
			format := estype.JoinDateFormats(formats...)
			p.Format = &format
		}
	}
}

// WithDateNanosDocValues sets whether doc values are enabled.
func WithDateNanosDocValues(v bool) DateNanosPropertyOption {
	return func(p *types.DateNanosProperty) { p.DocValues = &v }
}

// WithDateNanosIgnoreMalformed sets whether to ignore malformed values.
func WithDateNanosIgnoreMalformed(v bool) DateNanosPropertyOption {
	return func(p *types.DateNanosProperty) { p.IgnoreMalformed = &v }
}

// WithDateNanosIndex sets whether the field is indexed.
func WithDateNanosIndex(v bool) DateNanosPropertyOption {
	return func(p *types.DateNanosProperty) { p.Index = &v }
}

// WithDateNanosStore sets whether the field value is stored.
func WithDateNanosStore(v bool) DateNanosPropertyOption {
	return func(p *types.DateNanosProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Geo Point
// ---------------------------------------------------------------------------

// GeoPointPropertyOption is a functional option for configuring GeoPointProperty.
type GeoPointPropertyOption func(*types.GeoPointProperty)

// NewGeoPointProperty creates a new geo_point property mapping.
func NewGeoPointProperty(opts ...GeoPointPropertyOption) *types.GeoPointProperty {
	prop := types.NewGeoPointProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithGeoPointIgnoreMalformed sets whether to ignore malformed values.
func WithGeoPointIgnoreMalformed(v bool) GeoPointPropertyOption {
	return func(p *types.GeoPointProperty) { p.IgnoreMalformed = &v }
}

// WithGeoPointIgnoreZValue sets whether to ignore z-values.
func WithGeoPointIgnoreZValue(v bool) GeoPointPropertyOption {
	return func(p *types.GeoPointProperty) { p.IgnoreZValue = &v }
}

// WithGeoPointDocValues sets whether doc values are enabled.
func WithGeoPointDocValues(v bool) GeoPointPropertyOption {
	return func(p *types.GeoPointProperty) { p.DocValues = &v }
}

// WithGeoPointIndex sets whether the field is indexed.
func WithGeoPointIndex(v bool) GeoPointPropertyOption {
	return func(p *types.GeoPointProperty) { p.Index = &v }
}

// WithGeoPointStore sets whether the field value is stored.
func WithGeoPointStore(v bool) GeoPointPropertyOption {
	return func(p *types.GeoPointProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Geo Shape
// ---------------------------------------------------------------------------

// GeoShapePropertyOption is a functional option for configuring GeoShapeProperty.
type GeoShapePropertyOption func(*types.GeoShapeProperty)

// NewGeoShapeProperty creates a new geo_shape property mapping.
func NewGeoShapeProperty(opts ...GeoShapePropertyOption) *types.GeoShapeProperty {
	prop := types.NewGeoShapeProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithGeoShapeCoerce sets whether to coerce values.
func WithGeoShapeCoerce(v bool) GeoShapePropertyOption {
	return func(p *types.GeoShapeProperty) { p.Coerce = &v }
}

// WithGeoShapeIgnoreMalformed sets whether to ignore malformed values.
func WithGeoShapeIgnoreMalformed(v bool) GeoShapePropertyOption {
	return func(p *types.GeoShapeProperty) { p.IgnoreMalformed = &v }
}

// WithGeoShapeIgnoreZValue sets whether to ignore z-values.
func WithGeoShapeIgnoreZValue(v bool) GeoShapePropertyOption {
	return func(p *types.GeoShapeProperty) { p.IgnoreZValue = &v }
}

// WithGeoShapeDocValues sets whether doc values are enabled.
func WithGeoShapeDocValues(v bool) GeoShapePropertyOption {
	return func(p *types.GeoShapeProperty) { p.DocValues = &v }
}

// WithGeoShapeIndex sets whether the field is indexed.
func WithGeoShapeIndex(v bool) GeoShapePropertyOption {
	return func(p *types.GeoShapeProperty) { p.Index = &v }
}

// WithGeoShapeStore sets whether the field value is stored.
func WithGeoShapeStore(v bool) GeoShapePropertyOption {
	return func(p *types.GeoShapeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Shape
// ---------------------------------------------------------------------------

// ShapePropertyOption is a functional option for configuring ShapeProperty.
type ShapePropertyOption func(*types.ShapeProperty)

// NewShapeProperty creates a new shape property mapping.
func NewShapeProperty(opts ...ShapePropertyOption) *types.ShapeProperty {
	prop := types.NewShapeProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithShapeCoerce sets whether to coerce values.
func WithShapeCoerce(v bool) ShapePropertyOption {
	return func(p *types.ShapeProperty) { p.Coerce = &v }
}

// WithShapeIgnoreMalformed sets whether to ignore malformed values.
func WithShapeIgnoreMalformed(v bool) ShapePropertyOption {
	return func(p *types.ShapeProperty) { p.IgnoreMalformed = &v }
}

// WithShapeIgnoreZValue sets whether to ignore z-values.
func WithShapeIgnoreZValue(v bool) ShapePropertyOption {
	return func(p *types.ShapeProperty) { p.IgnoreZValue = &v }
}

// WithShapeDocValues sets whether doc values are enabled.
func WithShapeDocValues(v bool) ShapePropertyOption {
	return func(p *types.ShapeProperty) { p.DocValues = &v }
}

// WithShapeStore sets whether the field value is stored.
func WithShapeStore(v bool) ShapePropertyOption {
	return func(p *types.ShapeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Point
// ---------------------------------------------------------------------------

// PointPropertyOption is a functional option for configuring PointProperty.
type PointPropertyOption func(*types.PointProperty)

// NewPointProperty creates a new point property mapping.
func NewPointProperty(opts ...PointPropertyOption) *types.PointProperty {
	prop := types.NewPointProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithPointIgnoreMalformed sets whether to ignore malformed values.
func WithPointIgnoreMalformed(v bool) PointPropertyOption {
	return func(p *types.PointProperty) { p.IgnoreMalformed = &v }
}

// WithPointIgnoreZValue sets whether to ignore z-values.
func WithPointIgnoreZValue(v bool) PointPropertyOption {
	return func(p *types.PointProperty) { p.IgnoreZValue = &v }
}

// WithPointDocValues sets whether doc values are enabled.
func WithPointDocValues(v bool) PointPropertyOption {
	return func(p *types.PointProperty) { p.DocValues = &v }
}

// WithPointStore sets whether the field value is stored.
func WithPointStore(v bool) PointPropertyOption {
	return func(p *types.PointProperty) { p.Store = &v }
}

// WithPointNullValue sets the null value.
func WithPointNullValue(v string) PointPropertyOption {
	return func(p *types.PointProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Integer Range
// ---------------------------------------------------------------------------

// IntegerRangePropertyOption is a functional option for configuring IntegerRangeProperty.
type IntegerRangePropertyOption func(*types.IntegerRangeProperty)

// NewIntegerRangeProperty creates a new integer_range property mapping.
func NewIntegerRangeProperty(opts ...IntegerRangePropertyOption) *types.IntegerRangeProperty {
	prop := types.NewIntegerRangeProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithIntegerRangeCoerce sets whether to coerce values.
func WithIntegerRangeCoerce(v bool) IntegerRangePropertyOption {
	return func(p *types.IntegerRangeProperty) { p.Coerce = &v }
}

// WithIntegerRangeDocValues sets whether doc values are enabled.
func WithIntegerRangeDocValues(v bool) IntegerRangePropertyOption {
	return func(p *types.IntegerRangeProperty) { p.DocValues = &v }
}

// WithIntegerRangeIndex sets whether the field is indexed.
func WithIntegerRangeIndex(v bool) IntegerRangePropertyOption {
	return func(p *types.IntegerRangeProperty) { p.Index = &v }
}

// WithIntegerRangeStore sets whether the field value is stored.
func WithIntegerRangeStore(v bool) IntegerRangePropertyOption {
	return func(p *types.IntegerRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Long Range
// ---------------------------------------------------------------------------

// LongRangePropertyOption is a functional option for configuring LongRangeProperty.
type LongRangePropertyOption func(*types.LongRangeProperty)

// NewLongRangeProperty creates a new long_range property mapping.
func NewLongRangeProperty(opts ...LongRangePropertyOption) *types.LongRangeProperty {
	prop := types.NewLongRangeProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithLongRangeCoerce sets whether to coerce values.
func WithLongRangeCoerce(v bool) LongRangePropertyOption {
	return func(p *types.LongRangeProperty) { p.Coerce = &v }
}

// WithLongRangeDocValues sets whether doc values are enabled.
func WithLongRangeDocValues(v bool) LongRangePropertyOption {
	return func(p *types.LongRangeProperty) { p.DocValues = &v }
}

// WithLongRangeIndex sets whether the field is indexed.
func WithLongRangeIndex(v bool) LongRangePropertyOption {
	return func(p *types.LongRangeProperty) { p.Index = &v }
}

// WithLongRangeStore sets whether the field value is stored.
func WithLongRangeStore(v bool) LongRangePropertyOption {
	return func(p *types.LongRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Float Range
// ---------------------------------------------------------------------------

// FloatRangePropertyOption is a functional option for configuring FloatRangeProperty.
type FloatRangePropertyOption func(*types.FloatRangeProperty)

// NewFloatRangeProperty creates a new float_range property mapping.
func NewFloatRangeProperty(opts ...FloatRangePropertyOption) *types.FloatRangeProperty {
	prop := types.NewFloatRangeProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithFloatRangeCoerce sets whether to coerce values.
func WithFloatRangeCoerce(v bool) FloatRangePropertyOption {
	return func(p *types.FloatRangeProperty) { p.Coerce = &v }
}

// WithFloatRangeDocValues sets whether doc values are enabled.
func WithFloatRangeDocValues(v bool) FloatRangePropertyOption {
	return func(p *types.FloatRangeProperty) { p.DocValues = &v }
}

// WithFloatRangeIndex sets whether the field is indexed.
func WithFloatRangeIndex(v bool) FloatRangePropertyOption {
	return func(p *types.FloatRangeProperty) { p.Index = &v }
}

// WithFloatRangeStore sets whether the field value is stored.
func WithFloatRangeStore(v bool) FloatRangePropertyOption {
	return func(p *types.FloatRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Double Range
// ---------------------------------------------------------------------------

// DoubleRangePropertyOption is a functional option for configuring DoubleRangeProperty.
type DoubleRangePropertyOption func(*types.DoubleRangeProperty)

// NewDoubleRangeProperty creates a new double_range property mapping.
func NewDoubleRangeProperty(opts ...DoubleRangePropertyOption) *types.DoubleRangeProperty {
	prop := types.NewDoubleRangeProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithDoubleRangeCoerce sets whether to coerce values.
func WithDoubleRangeCoerce(v bool) DoubleRangePropertyOption {
	return func(p *types.DoubleRangeProperty) { p.Coerce = &v }
}

// WithDoubleRangeDocValues sets whether doc values are enabled.
func WithDoubleRangeDocValues(v bool) DoubleRangePropertyOption {
	return func(p *types.DoubleRangeProperty) { p.DocValues = &v }
}

// WithDoubleRangeIndex sets whether the field is indexed.
func WithDoubleRangeIndex(v bool) DoubleRangePropertyOption {
	return func(p *types.DoubleRangeProperty) { p.Index = &v }
}

// WithDoubleRangeStore sets whether the field value is stored.
func WithDoubleRangeStore(v bool) DoubleRangePropertyOption {
	return func(p *types.DoubleRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Date Range
// ---------------------------------------------------------------------------

// DateRangePropertyOption is a functional option for configuring DateRangeProperty.
type DateRangePropertyOption func(*types.DateRangeProperty)

// NewDateRangeProperty creates a new date_range property mapping.
func NewDateRangeProperty(opts ...DateRangePropertyOption) *types.DateRangeProperty {
	prop := types.NewDateRangeProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithDateRangeFormat sets the date format.
func WithDateRangeFormat(formats ...estype.DateFormat) DateRangePropertyOption {
	return func(p *types.DateRangeProperty) {
		if len(formats) > 0 {
			format := estype.JoinDateFormats(formats...)
			p.Format = &format
		}
	}
}

// WithDateRangeCoerce sets whether to coerce values.
func WithDateRangeCoerce(v bool) DateRangePropertyOption {
	return func(p *types.DateRangeProperty) { p.Coerce = &v }
}

// WithDateRangeDocValues sets whether doc values are enabled.
func WithDateRangeDocValues(v bool) DateRangePropertyOption {
	return func(p *types.DateRangeProperty) { p.DocValues = &v }
}

// WithDateRangeIndex sets whether the field is indexed.
func WithDateRangeIndex(v bool) DateRangePropertyOption {
	return func(p *types.DateRangeProperty) { p.Index = &v }
}

// WithDateRangeStore sets whether the field value is stored.
func WithDateRangeStore(v bool) DateRangePropertyOption {
	return func(p *types.DateRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// IP Range
// ---------------------------------------------------------------------------

// IpRangePropertyOption is a functional option for configuring IpRangeProperty.
type IpRangePropertyOption func(*types.IpRangeProperty)

// NewIpRangeProperty creates a new ip_range property mapping.
func NewIpRangeProperty(opts ...IpRangePropertyOption) *types.IpRangeProperty {
	prop := types.NewIpRangeProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithIpRangeCoerce sets whether to coerce values.
func WithIpRangeCoerce(v bool) IpRangePropertyOption {
	return func(p *types.IpRangeProperty) { p.Coerce = &v }
}

// WithIpRangeDocValues sets whether doc values are enabled.
func WithIpRangeDocValues(v bool) IpRangePropertyOption {
	return func(p *types.IpRangeProperty) { p.DocValues = &v }
}

// WithIpRangeIndex sets whether the field is indexed.
func WithIpRangeIndex(v bool) IpRangePropertyOption {
	return func(p *types.IpRangeProperty) { p.Index = &v }
}

// WithIpRangeStore sets whether the field value is stored.
func WithIpRangeStore(v bool) IpRangePropertyOption {
	return func(p *types.IpRangeProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Object
// ---------------------------------------------------------------------------

// ObjectPropertyOption is a functional option for configuring ObjectProperty.
type ObjectPropertyOption func(*types.ObjectProperty)

// NewObjectProperty creates a new object property mapping.
func NewObjectProperty(opts ...ObjectPropertyOption) *types.ObjectProperty {
	prop := types.NewObjectProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithObjectProperties sets the properties of the object.
func WithObjectProperties(v map[string]types.Property) ObjectPropertyOption {
	return func(p *types.ObjectProperty) { p.Properties = v }
}

// WithObjectEnabled sets whether the object is enabled.
func WithObjectEnabled(v bool) ObjectPropertyOption {
	return func(p *types.ObjectProperty) { p.Enabled = &v }
}

// WithObjectStore sets whether the field value is stored.
func WithObjectStore(v bool) ObjectPropertyOption {
	return func(p *types.ObjectProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Nested
// ---------------------------------------------------------------------------

// NestedPropertyOption is a functional option for configuring NestedProperty.
type NestedPropertyOption func(*types.NestedProperty)

// NewNestedProperty creates a new nested property mapping.
func NewNestedProperty(opts ...NestedPropertyOption) *types.NestedProperty {
	prop := types.NewNestedProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithNestedProperties sets the properties of the nested mapping.
func WithNestedProperties(v map[string]types.Property) NestedPropertyOption {
	return func(p *types.NestedProperty) { p.Properties = v }
}

// WithNestedEnabled sets whether the nested mapping is enabled.
func WithNestedEnabled(v bool) NestedPropertyOption {
	return func(p *types.NestedProperty) { p.Enabled = &v }
}

// WithNestedIncludeInParent sets whether to include in parent.
func WithNestedIncludeInParent(v bool) NestedPropertyOption {
	return func(p *types.NestedProperty) { p.IncludeInParent = &v }
}

// WithNestedIncludeInRoot sets whether to include in root.
func WithNestedIncludeInRoot(v bool) NestedPropertyOption {
	return func(p *types.NestedProperty) { p.IncludeInRoot = &v }
}

// WithNestedStore sets whether the field value is stored.
func WithNestedStore(v bool) NestedPropertyOption {
	return func(p *types.NestedProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Flattened
// ---------------------------------------------------------------------------

// FlattenedPropertyOption is a functional option for configuring FlattenedProperty.
type FlattenedPropertyOption func(*types.FlattenedProperty)

// NewFlattenedProperty creates a new flattened property mapping.
func NewFlattenedProperty(opts ...FlattenedPropertyOption) *types.FlattenedProperty {
	prop := types.NewFlattenedProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithFlattenedDepthLimit sets the maximum depth.
func WithFlattenedDepthLimit(v int) FlattenedPropertyOption {
	return func(p *types.FlattenedProperty) { p.DepthLimit = &v }
}

// WithFlattenedDocValues sets whether doc values are enabled.
func WithFlattenedDocValues(v bool) FlattenedPropertyOption {
	return func(p *types.FlattenedProperty) { p.DocValues = &v }
}

// WithFlattenedIndex sets whether the field is indexed.
func WithFlattenedIndex(v bool) FlattenedPropertyOption {
	return func(p *types.FlattenedProperty) { p.Index = &v }
}

// WithFlattenedIgnoreAbove sets the maximum string length.
func WithFlattenedIgnoreAbove(v int) FlattenedPropertyOption {
	return func(p *types.FlattenedProperty) { p.IgnoreAbove = &v }
}

// WithFlattenedNullValue sets the null value.
func WithFlattenedNullValue(v string) FlattenedPropertyOption {
	return func(p *types.FlattenedProperty) { p.NullValue = &v }
}

// WithFlattenedEagerGlobalOrdinals sets whether to eagerly load global ordinals.
func WithFlattenedEagerGlobalOrdinals(v bool) FlattenedPropertyOption {
	return func(p *types.FlattenedProperty) { p.EagerGlobalOrdinals = &v }
}

// WithFlattenedSimilarity sets the similarity algorithm.
func WithFlattenedSimilarity(v string) FlattenedPropertyOption {
	return func(p *types.FlattenedProperty) { p.Similarity = &v }
}

// WithFlattenedSplitQueriesOnWhitespace sets whether to split queries on whitespace.
func WithFlattenedSplitQueriesOnWhitespace(v bool) FlattenedPropertyOption {
	return func(p *types.FlattenedProperty) { p.SplitQueriesOnWhitespace = &v }
}

// ---------------------------------------------------------------------------
// Join
// ---------------------------------------------------------------------------

// JoinPropertyOption is a functional option for configuring JoinProperty.
type JoinPropertyOption func(*types.JoinProperty)

// NewJoinProperty creates a new join property mapping.
func NewJoinProperty(opts ...JoinPropertyOption) *types.JoinProperty {
	prop := types.NewJoinProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithJoinRelations sets the parent/child relations.
func WithJoinRelations(v map[string][]string) JoinPropertyOption {
	return func(p *types.JoinProperty) { p.Relations = v }
}

// WithJoinEagerGlobalOrdinals sets whether to eagerly load global ordinals.
func WithJoinEagerGlobalOrdinals(v bool) JoinPropertyOption {
	return func(p *types.JoinProperty) { p.EagerGlobalOrdinals = &v }
}

// ---------------------------------------------------------------------------
// Passthrough Object
// ---------------------------------------------------------------------------

// PassthroughObjectPropertyOption is a functional option for configuring PassthroughObjectProperty.
type PassthroughObjectPropertyOption func(*types.PassthroughObjectProperty)

// NewPassthroughObjectProperty creates a new passthrough object property mapping.
func NewPassthroughObjectProperty(opts ...PassthroughObjectPropertyOption) *types.PassthroughObjectProperty {
	prop := types.NewPassthroughObjectProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithPassthroughObjectProperties sets the properties.
func WithPassthroughObjectProperties(v map[string]types.Property) PassthroughObjectPropertyOption {
	return func(p *types.PassthroughObjectProperty) { p.Properties = v }
}

// WithPassthroughObjectEnabled sets whether the object is enabled.
func WithPassthroughObjectEnabled(v bool) PassthroughObjectPropertyOption {
	return func(p *types.PassthroughObjectProperty) { p.Enabled = &v }
}

// WithPassthroughObjectPriority sets the priority.
func WithPassthroughObjectPriority(v int) PassthroughObjectPropertyOption {
	return func(p *types.PassthroughObjectProperty) { p.Priority = &v }
}

// WithPassthroughObjectStore sets whether the field value is stored.
func WithPassthroughObjectStore(v bool) PassthroughObjectPropertyOption {
	return func(p *types.PassthroughObjectProperty) { p.Store = &v }
}

// WithPassthroughObjectTimeSeriesDimension sets whether this is a time series dimension.
func WithPassthroughObjectTimeSeriesDimension(v bool) PassthroughObjectPropertyOption {
	return func(p *types.PassthroughObjectProperty) { p.TimeSeriesDimension = &v }
}

// ---------------------------------------------------------------------------
// IP
// ---------------------------------------------------------------------------

// IpPropertyOption is a functional option for configuring IpProperty.
type IpPropertyOption func(*types.IpProperty)

// NewIpProperty creates a new ip property mapping.
func NewIpProperty(opts ...IpPropertyOption) *types.IpProperty {
	prop := types.NewIpProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithIpDocValues sets whether doc values are enabled.
func WithIpDocValues(v bool) IpPropertyOption {
	return func(p *types.IpProperty) { p.DocValues = &v }
}

// WithIpIgnoreMalformed sets whether to ignore malformed values.
func WithIpIgnoreMalformed(v bool) IpPropertyOption {
	return func(p *types.IpProperty) { p.IgnoreMalformed = &v }
}

// WithIpIndex sets whether the field is indexed.
func WithIpIndex(v bool) IpPropertyOption {
	return func(p *types.IpProperty) { p.Index = &v }
}

// WithIpStore sets whether the field value is stored.
func WithIpStore(v bool) IpPropertyOption {
	return func(p *types.IpProperty) { p.Store = &v }
}

// WithIpNullValue sets the null value.
func WithIpNullValue(v string) IpPropertyOption {
	return func(p *types.IpProperty) { p.NullValue = &v }
}

// ---------------------------------------------------------------------------
// Binary
// ---------------------------------------------------------------------------

// BinaryPropertyOption is a functional option for configuring BinaryProperty.
type BinaryPropertyOption func(*types.BinaryProperty)

// NewBinaryProperty creates a new binary property mapping.
func NewBinaryProperty(opts ...BinaryPropertyOption) *types.BinaryProperty {
	prop := types.NewBinaryProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithBinaryDocValues sets whether doc values are enabled.
func WithBinaryDocValues(v bool) BinaryPropertyOption {
	return func(p *types.BinaryProperty) { p.DocValues = &v }
}

// WithBinaryStore sets whether the field value is stored.
func WithBinaryStore(v bool) BinaryPropertyOption {
	return func(p *types.BinaryProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Token Count
// ---------------------------------------------------------------------------

// TokenCountPropertyOption is a functional option for configuring TokenCountProperty.
type TokenCountPropertyOption func(*types.TokenCountProperty)

// NewTokenCountProperty creates a new token_count property mapping.
func NewTokenCountProperty(opts ...TokenCountPropertyOption) *types.TokenCountProperty {
	prop := types.NewTokenCountProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithTokenCountAnalyzer sets the analyzer.
func WithTokenCountAnalyzer(v string) TokenCountPropertyOption {
	return func(p *types.TokenCountProperty) { p.Analyzer = &v }
}

// WithTokenCountDocValues sets whether doc values are enabled.
func WithTokenCountDocValues(v bool) TokenCountPropertyOption {
	return func(p *types.TokenCountProperty) { p.DocValues = &v }
}

// WithTokenCountIndex sets whether the field is indexed.
func WithTokenCountIndex(v bool) TokenCountPropertyOption {
	return func(p *types.TokenCountProperty) { p.Index = &v }
}

// WithTokenCountStore sets whether the field value is stored.
func WithTokenCountStore(v bool) TokenCountPropertyOption {
	return func(p *types.TokenCountProperty) { p.Store = &v }
}

// WithTokenCountEnablePositionIncrements sets whether to count position increments.
func WithTokenCountEnablePositionIncrements(v bool) TokenCountPropertyOption {
	return func(p *types.TokenCountProperty) { p.EnablePositionIncrements = &v }
}

// ---------------------------------------------------------------------------
// Percolator
// ---------------------------------------------------------------------------

// PercolatorPropertyOption is a functional option for configuring PercolatorProperty.
type PercolatorPropertyOption func(*types.PercolatorProperty)

// NewPercolatorProperty creates a new percolator property mapping.
func NewPercolatorProperty(opts ...PercolatorPropertyOption) *types.PercolatorProperty {
	prop := types.NewPercolatorProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// ---------------------------------------------------------------------------
// Field Alias
// ---------------------------------------------------------------------------

// FieldAliasPropertyOption is a functional option for configuring FieldAliasProperty.
type FieldAliasPropertyOption func(*types.FieldAliasProperty)

// NewFieldAliasProperty creates a new alias property mapping.
func NewFieldAliasProperty(opts ...FieldAliasPropertyOption) *types.FieldAliasProperty {
	prop := types.NewFieldAliasProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithFieldAliasPath sets the path to the target field.
func WithFieldAliasPath(v string) FieldAliasPropertyOption {
	return func(p *types.FieldAliasProperty) { p.Path = &v }
}

// ---------------------------------------------------------------------------
// Histogram
// ---------------------------------------------------------------------------

// HistogramPropertyOption is a functional option for configuring HistogramProperty.
type HistogramPropertyOption func(*types.HistogramProperty)

// NewHistogramProperty creates a new histogram property mapping.
func NewHistogramProperty(opts ...HistogramPropertyOption) *types.HistogramProperty {
	prop := types.NewHistogramProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithHistogramIgnoreMalformed sets whether to ignore malformed values.
func WithHistogramIgnoreMalformed(v bool) HistogramPropertyOption {
	return func(p *types.HistogramProperty) { p.IgnoreMalformed = &v }
}

// ---------------------------------------------------------------------------
// Version
// ---------------------------------------------------------------------------

// VersionPropertyOption is a functional option for configuring VersionProperty.
type VersionPropertyOption func(*types.VersionProperty)

// NewVersionProperty creates a new version property mapping.
func NewVersionProperty(opts ...VersionPropertyOption) *types.VersionProperty {
	prop := types.NewVersionProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithVersionDocValues sets whether doc values are enabled.
func WithVersionDocValues(v bool) VersionPropertyOption {
	return func(p *types.VersionProperty) { p.DocValues = &v }
}

// WithVersionStore sets whether the field value is stored.
func WithVersionStore(v bool) VersionPropertyOption {
	return func(p *types.VersionProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Dense Vector
// ---------------------------------------------------------------------------

// DenseVectorPropertyOption is a functional option for configuring DenseVectorProperty.
type DenseVectorPropertyOption func(*types.DenseVectorProperty)

// NewDenseVectorProperty creates a new dense_vector property mapping.
func NewDenseVectorProperty(opts ...DenseVectorPropertyOption) *types.DenseVectorProperty {
	prop := types.NewDenseVectorProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithDenseVectorDims sets the number of dimensions.
func WithDenseVectorDims(v int) DenseVectorPropertyOption {
	return func(p *types.DenseVectorProperty) { p.Dims = &v }
}

// WithDenseVectorIndex sets whether the field is indexed.
func WithDenseVectorIndex(v bool) DenseVectorPropertyOption {
	return func(p *types.DenseVectorProperty) { p.Index = &v }
}

// ---------------------------------------------------------------------------
// Sparse Vector
// ---------------------------------------------------------------------------

// SparseVectorPropertyOption is a functional option for configuring SparseVectorProperty.
type SparseVectorPropertyOption func(*types.SparseVectorProperty)

// NewSparseVectorProperty creates a new sparse_vector property mapping.
func NewSparseVectorProperty(opts ...SparseVectorPropertyOption) *types.SparseVectorProperty {
	prop := types.NewSparseVectorProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithSparseVectorStore sets whether the field value is stored.
func WithSparseVectorStore(v bool) SparseVectorPropertyOption {
	return func(p *types.SparseVectorProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// Rank Feature
// ---------------------------------------------------------------------------

// RankFeaturePropertyOption is a functional option for configuring RankFeatureProperty.
type RankFeaturePropertyOption func(*types.RankFeatureProperty)

// NewRankFeatureProperty creates a new rank_feature property mapping.
func NewRankFeatureProperty(opts ...RankFeaturePropertyOption) *types.RankFeatureProperty {
	prop := types.NewRankFeatureProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithRankFeaturePositiveScoreImpact sets whether a positive score has a positive impact.
func WithRankFeaturePositiveScoreImpact(v bool) RankFeaturePropertyOption {
	return func(p *types.RankFeatureProperty) { p.PositiveScoreImpact = &v }
}

// ---------------------------------------------------------------------------
// Rank Features
// ---------------------------------------------------------------------------

// RankFeaturesPropertyOption is a functional option for configuring RankFeaturesProperty.
type RankFeaturesPropertyOption func(*types.RankFeaturesProperty)

// NewRankFeaturesProperty creates a new rank_features property mapping.
func NewRankFeaturesProperty(opts ...RankFeaturesPropertyOption) *types.RankFeaturesProperty {
	prop := types.NewRankFeaturesProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithRankFeaturesPositiveScoreImpact sets whether a positive score has a positive impact.
func WithRankFeaturesPositiveScoreImpact(v bool) RankFeaturesPropertyOption {
	return func(p *types.RankFeaturesProperty) { p.PositiveScoreImpact = &v }
}

// ---------------------------------------------------------------------------
// Rank Vector
// ---------------------------------------------------------------------------

// RankVectorPropertyOption is a functional option for configuring RankVectorProperty.
type RankVectorPropertyOption func(*types.RankVectorProperty)

// NewRankVectorProperty creates a new rank_vector property mapping.
func NewRankVectorProperty(opts ...RankVectorPropertyOption) *types.RankVectorProperty {
	prop := types.NewRankVectorProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithRankVectorDims sets the number of dimensions.
func WithRankVectorDims(v int) RankVectorPropertyOption {
	return func(p *types.RankVectorProperty) { p.Dims = &v }
}

// ---------------------------------------------------------------------------
// Semantic Text
// ---------------------------------------------------------------------------

// SemanticTextPropertyOption is a functional option for configuring SemanticTextProperty.
type SemanticTextPropertyOption func(*types.SemanticTextProperty)

// NewSemanticTextProperty creates a new semantic_text property mapping.
func NewSemanticTextProperty(opts ...SemanticTextPropertyOption) *types.SemanticTextProperty {
	prop := types.NewSemanticTextProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithSemanticTextInferenceId sets the inference ID.
func WithSemanticTextInferenceId(v string) SemanticTextPropertyOption {
	return func(p *types.SemanticTextProperty) { p.InferenceId = &v }
}

// WithSemanticTextSearchInferenceId sets the search inference ID.
func WithSemanticTextSearchInferenceId(v string) SemanticTextPropertyOption {
	return func(p *types.SemanticTextProperty) { p.SearchInferenceId = &v }
}

// ---------------------------------------------------------------------------
// Aggregate Metric Double
// ---------------------------------------------------------------------------

// AggregateMetricDoublePropertyOption is a functional option for configuring AggregateMetricDoubleProperty.
type AggregateMetricDoublePropertyOption func(*types.AggregateMetricDoubleProperty)

// NewAggregateMetricDoubleProperty creates a new aggregate_metric_double property mapping.
func NewAggregateMetricDoubleProperty(opts ...AggregateMetricDoublePropertyOption) *types.AggregateMetricDoubleProperty {
	prop := types.NewAggregateMetricDoubleProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithAggregateMetricDoubleDefaultMetric sets the default metric.
func WithAggregateMetricDoubleDefaultMetric(v string) AggregateMetricDoublePropertyOption {
	return func(p *types.AggregateMetricDoubleProperty) { p.DefaultMetric = v }
}

// WithAggregateMetricDoubleMetrics sets the list of metrics.
func WithAggregateMetricDoubleMetrics(v []string) AggregateMetricDoublePropertyOption {
	return func(p *types.AggregateMetricDoubleProperty) { p.Metrics = v }
}

// WithAggregateMetricDoubleIgnoreMalformed sets whether to ignore malformed values.
func WithAggregateMetricDoubleIgnoreMalformed(v bool) AggregateMetricDoublePropertyOption {
	return func(p *types.AggregateMetricDoubleProperty) { p.IgnoreMalformed = &v }
}

// ---------------------------------------------------------------------------
// Murmur3 Hash
// ---------------------------------------------------------------------------

// Murmur3HashPropertyOption is a functional option for configuring Murmur3HashProperty.
type Murmur3HashPropertyOption func(*types.Murmur3HashProperty)

// NewMurmur3HashProperty creates a new murmur3 hash property mapping.
func NewMurmur3HashProperty(opts ...Murmur3HashPropertyOption) *types.Murmur3HashProperty {
	prop := types.NewMurmur3HashProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithMurmur3HashDocValues sets whether doc values are enabled.
func WithMurmur3HashDocValues(v bool) Murmur3HashPropertyOption {
	return func(p *types.Murmur3HashProperty) { p.DocValues = &v }
}

// WithMurmur3HashStore sets whether the field value is stored.
func WithMurmur3HashStore(v bool) Murmur3HashPropertyOption {
	return func(p *types.Murmur3HashProperty) { p.Store = &v }
}

// ---------------------------------------------------------------------------
// ICU Collation
// ---------------------------------------------------------------------------

// IcuCollationPropertyOption is a functional option for configuring IcuCollationProperty.
type IcuCollationPropertyOption func(*types.IcuCollationProperty)

// NewIcuCollationProperty creates a new icu_collation_keyword property mapping.
func NewIcuCollationProperty(opts ...IcuCollationPropertyOption) *types.IcuCollationProperty {
	prop := types.NewIcuCollationProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithIcuCollationLanguage sets the language.
func WithIcuCollationLanguage(v string) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.Language = &v }
}

// WithIcuCollationCountry sets the country.
func WithIcuCollationCountry(v string) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.Country = &v }
}

// WithIcuCollationDocValues sets whether doc values are enabled.
func WithIcuCollationDocValues(v bool) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.DocValues = &v }
}

// WithIcuCollationIndex sets whether the field is indexed.
func WithIcuCollationIndex(v bool) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.Index = &v }
}

// WithIcuCollationStore sets whether the field value is stored.
func WithIcuCollationStore(v bool) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.Store = &v }
}

// WithIcuCollationNullValue sets the null value.
func WithIcuCollationNullValue(v string) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.NullValue = &v }
}

// WithIcuCollationNorms sets whether norms are enabled.
func WithIcuCollationNorms(v bool) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.Norms = &v }
}

// WithIcuCollationRules sets the collation rules.
func WithIcuCollationRules(v string) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.Rules = &v }
}

// WithIcuCollationVariant sets the collation variant.
func WithIcuCollationVariant(v string) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.Variant = &v }
}

// WithIcuCollationCaseLevel sets whether case-level sorting is enabled.
func WithIcuCollationCaseLevel(v bool) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.CaseLevel = &v }
}

// WithIcuCollationNumeric sets whether numeric sorting is enabled.
func WithIcuCollationNumeric(v bool) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.Numeric = &v }
}

// WithIcuCollationHiraganaQuaternaryMode sets whether hiragana quaternary mode is enabled.
func WithIcuCollationHiraganaQuaternaryMode(v bool) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.HiraganaQuaternaryMode = &v }
}

// WithIcuCollationVariableTop sets the variable top.
func WithIcuCollationVariableTop(v string) IcuCollationPropertyOption {
	return func(p *types.IcuCollationProperty) { p.VariableTop = &v }
}

// ---------------------------------------------------------------------------
// Dynamic
// ---------------------------------------------------------------------------

// DynamicPropertyOption is a functional option for configuring DynamicProperty.
type DynamicPropertyOption func(*types.DynamicProperty)

// NewDynamicProperty creates a new dynamic property mapping.
func NewDynamicProperty(opts ...DynamicPropertyOption) *types.DynamicProperty {
	prop := types.NewDynamicProperty()
	for _, opt := range opts {
		opt(prop)
	}
	return prop
}

// WithDynamicAnalyzer sets the analyzer.
func WithDynamicAnalyzer(v string) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.Analyzer = &v }
}

// WithDynamicSearchAnalyzer sets the search analyzer.
func WithDynamicSearchAnalyzer(v string) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.SearchAnalyzer = &v }
}

// WithDynamicCoerce sets whether to coerce values.
func WithDynamicCoerce(v bool) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.Coerce = &v }
}

// WithDynamicDocValues sets whether doc values are enabled.
func WithDynamicDocValues(v bool) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.DocValues = &v }
}

// WithDynamicEnabled sets whether the field is enabled.
func WithDynamicEnabled(v bool) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.Enabled = &v }
}

// WithDynamicFormat sets the format.
func WithDynamicFormat(v string) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.Format = &v }
}

// WithDynamicIgnoreMalformed sets whether to ignore malformed values.
func WithDynamicIgnoreMalformed(v bool) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.IgnoreMalformed = &v }
}

// WithDynamicIndex sets whether the field is indexed.
func WithDynamicIndex(v bool) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.Index = &v }
}

// WithDynamicStore sets whether the field value is stored.
func WithDynamicStore(v bool) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.Store = &v }
}

// WithDynamicNorms sets whether norms are enabled.
func WithDynamicNorms(v bool) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.Norms = &v }
}

// WithDynamicLocale sets the locale.
func WithDynamicLocale(v string) DynamicPropertyOption {
	return func(p *types.DynamicProperty) { p.Locale = &v }
}
