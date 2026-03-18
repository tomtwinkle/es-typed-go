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

// CountedKeywordProperty represents an Elasticsearch "counted_keyword" field mapping.
// Use [NewCountedKeywordProperty] to construct one.
type CountedKeywordProperty struct{}

// ESTypeName returns the Elasticsearch type name for a counted_keyword property.
func (CountedKeywordProperty) ESTypeName() string { return "counted_keyword" }

// NewCountedKeywordProperty creates a new [CountedKeywordProperty].
func NewCountedKeywordProperty() CountedKeywordProperty {
	return CountedKeywordProperty{}
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

// ---------------------------------------------------------------------------
// Numeric properties
// ---------------------------------------------------------------------------

// IntegerNumberProperty represents an Elasticsearch "integer" field mapping.
type IntegerNumberProperty struct{}

// ESTypeName returns the Elasticsearch type name for an integer property.
func (IntegerNumberProperty) ESTypeName() string { return "integer" }

// NewIntegerNumberProperty creates a new [IntegerNumberProperty].
func NewIntegerNumberProperty() IntegerNumberProperty {
	return IntegerNumberProperty{}
}

// LongNumberProperty represents an Elasticsearch "long" field mapping.
type LongNumberProperty struct{}

// ESTypeName returns the Elasticsearch type name for a long property.
func (LongNumberProperty) ESTypeName() string { return "long" }

// NewLongNumberProperty creates a new [LongNumberProperty].
func NewLongNumberProperty() LongNumberProperty {
	return LongNumberProperty{}
}

// ShortNumberProperty represents an Elasticsearch "short" field mapping.
type ShortNumberProperty struct{}

// ESTypeName returns the Elasticsearch type name for a short property.
func (ShortNumberProperty) ESTypeName() string { return "short" }

// NewShortNumberProperty creates a new [ShortNumberProperty].
func NewShortNumberProperty() ShortNumberProperty {
	return ShortNumberProperty{}
}

// ByteNumberProperty represents an Elasticsearch "byte" field mapping.
type ByteNumberProperty struct{}

// ESTypeName returns the Elasticsearch type name for a byte property.
func (ByteNumberProperty) ESTypeName() string { return "byte" }

// NewByteNumberProperty creates a new [ByteNumberProperty].
func NewByteNumberProperty() ByteNumberProperty {
	return ByteNumberProperty{}
}

// DoubleNumberProperty represents an Elasticsearch "double" field mapping.
type DoubleNumberProperty struct{}

// ESTypeName returns the Elasticsearch type name for a double property.
func (DoubleNumberProperty) ESTypeName() string { return "double" }

// NewDoubleNumberProperty creates a new [DoubleNumberProperty].
func NewDoubleNumberProperty() DoubleNumberProperty {
	return DoubleNumberProperty{}
}

// FloatNumberProperty represents an Elasticsearch "float" field mapping.
type FloatNumberProperty struct{}

// ESTypeName returns the Elasticsearch type name for a float property.
func (FloatNumberProperty) ESTypeName() string { return "float" }

// NewFloatNumberProperty creates a new [FloatNumberProperty].
func NewFloatNumberProperty() FloatNumberProperty {
	return FloatNumberProperty{}
}

// HalfFloatNumberProperty represents an Elasticsearch "half_float" field mapping.
type HalfFloatNumberProperty struct{}

// ESTypeName returns the Elasticsearch type name for a half_float property.
func (HalfFloatNumberProperty) ESTypeName() string { return "half_float" }

// NewHalfFloatNumberProperty creates a new [HalfFloatNumberProperty].
func NewHalfFloatNumberProperty() HalfFloatNumberProperty {
	return HalfFloatNumberProperty{}
}

// UnsignedLongNumberProperty represents an Elasticsearch "unsigned_long" field mapping.
type UnsignedLongNumberProperty struct{}

// ESTypeName returns the Elasticsearch type name for an unsigned_long property.
func (UnsignedLongNumberProperty) ESTypeName() string { return "unsigned_long" }

// NewUnsignedLongNumberProperty creates a new [UnsignedLongNumberProperty].
func NewUnsignedLongNumberProperty() UnsignedLongNumberProperty {
	return UnsignedLongNumberProperty{}
}

// ScaledFloatNumberPropertyOption is a functional option for configuring a [ScaledFloatNumberProperty].
type ScaledFloatNumberPropertyOption func(*ScaledFloatNumberProperty)

// ScaledFloatNumberProperty represents an Elasticsearch "scaled_float" field mapping.
// Use [NewScaledFloatNumberProperty] to construct one with functional options.
type ScaledFloatNumberProperty struct {
	// ScalingFactor is the scaling factor to use when encoding values.
	// Values will be multiplied by this factor at index time and rounded
	// to the nearest long value (e.g. 100 for two decimal places).
	ScalingFactor *float64
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
	Format *string
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

// DateNanosPropertyOption is a functional option for configuring a [DateNanosProperty].
type DateNanosPropertyOption func(*DateNanosProperty)

// DateNanosProperty represents an Elasticsearch "date_nanos" field mapping
// with nanosecond resolution.
// Use [NewDateNanosProperty] to construct one with functional options.
type DateNanosProperty struct {
	// Format is the date format(s) that can be parsed.
	Format *string
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

// ---------------------------------------------------------------------------
// Geo
// ---------------------------------------------------------------------------

// GeoPointProperty represents an Elasticsearch "geo_point" field mapping.
type GeoPointProperty struct{}

// ESTypeName returns the Elasticsearch type name for a geo_point property.
func (GeoPointProperty) ESTypeName() string { return "geo_point" }

// NewGeoPointProperty creates a new [GeoPointProperty].
func NewGeoPointProperty() GeoPointProperty {
	return GeoPointProperty{}
}

// GeoShapeProperty represents an Elasticsearch "geo_shape" field mapping.
type GeoShapeProperty struct{}

// ESTypeName returns the Elasticsearch type name for a geo_shape property.
func (GeoShapeProperty) ESTypeName() string { return "geo_shape" }

// NewGeoShapeProperty creates a new [GeoShapeProperty].
func NewGeoShapeProperty() GeoShapeProperty {
	return GeoShapeProperty{}
}

// ShapeProperty represents an Elasticsearch "shape" field mapping
// for arbitrary cartesian geometries.
type ShapeProperty struct{}

// ESTypeName returns the Elasticsearch type name for a shape property.
func (ShapeProperty) ESTypeName() string { return "shape" }

// NewShapeProperty creates a new [ShapeProperty].
func NewShapeProperty() ShapeProperty {
	return ShapeProperty{}
}

// PointProperty represents an Elasticsearch "point" field mapping
// for arbitrary cartesian points.
type PointProperty struct{}

// ESTypeName returns the Elasticsearch type name for a point property.
func (PointProperty) ESTypeName() string { return "point" }

// NewPointProperty creates a new [PointProperty].
func NewPointProperty() PointProperty {
	return PointProperty{}
}

// ---------------------------------------------------------------------------
// Range
// ---------------------------------------------------------------------------

// IntegerRangeProperty represents an Elasticsearch "integer_range" field mapping.
type IntegerRangeProperty struct{}

// ESTypeName returns the Elasticsearch type name for an integer_range property.
func (IntegerRangeProperty) ESTypeName() string { return "integer_range" }

// NewIntegerRangeProperty creates a new [IntegerRangeProperty].
func NewIntegerRangeProperty() IntegerRangeProperty {
	return IntegerRangeProperty{}
}

// LongRangeProperty represents an Elasticsearch "long_range" field mapping.
type LongRangeProperty struct{}

// ESTypeName returns the Elasticsearch type name for a long_range property.
func (LongRangeProperty) ESTypeName() string { return "long_range" }

// NewLongRangeProperty creates a new [LongRangeProperty].
func NewLongRangeProperty() LongRangeProperty {
	return LongRangeProperty{}
}

// FloatRangeProperty represents an Elasticsearch "float_range" field mapping.
type FloatRangeProperty struct{}

// ESTypeName returns the Elasticsearch type name for a float_range property.
func (FloatRangeProperty) ESTypeName() string { return "float_range" }

// NewFloatRangeProperty creates a new [FloatRangeProperty].
func NewFloatRangeProperty() FloatRangeProperty {
	return FloatRangeProperty{}
}

// DoubleRangeProperty represents an Elasticsearch "double_range" field mapping.
type DoubleRangeProperty struct{}

// ESTypeName returns the Elasticsearch type name for a double_range property.
func (DoubleRangeProperty) ESTypeName() string { return "double_range" }

// NewDoubleRangeProperty creates a new [DoubleRangeProperty].
func NewDoubleRangeProperty() DoubleRangeProperty {
	return DoubleRangeProperty{}
}

// DateRangePropertyOption is a functional option for configuring a [DateRangeProperty].
type DateRangePropertyOption func(*DateRangeProperty)

// DateRangeProperty represents an Elasticsearch "date_range" field mapping.
// Use [NewDateRangeProperty] to construct one with functional options.
type DateRangeProperty struct {
	// Format is the date format(s) that can be parsed.
	Format *string
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

// IpRangeProperty represents an Elasticsearch "ip_range" field mapping.
type IpRangeProperty struct{}

// ESTypeName returns the Elasticsearch type name for an ip_range property.
func (IpRangeProperty) ESTypeName() string { return "ip_range" }

// NewIpRangeProperty creates a new [IpRangeProperty].
func NewIpRangeProperty() IpRangeProperty {
	return IpRangeProperty{}
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

// NestedPropertyOption is a functional option for configuring a [NestedProperty].
type NestedPropertyOption func(*NestedProperty)

// NestedProperty represents an Elasticsearch "nested" field mapping.
// Each nested object is indexed as a separate hidden document, allowing
// independent querying of nested objects.
// Use [NewNestedProperty] to construct one with functional options.
type NestedProperty struct {
	// Properties holds the child field mappings.
	Properties map[string]MappingProperty
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

// JoinPropertyOption is a functional option for configuring a [JoinProperty].
type JoinPropertyOption func(*JoinProperty)

// JoinProperty represents an Elasticsearch "join" field mapping that defines
// parent/child relationships within a single index.
// Use [NewJoinProperty] to construct one with functional options.
type JoinProperty struct {
	// Relations maps parent names to their child names.
	Relations map[string][]string
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

// PassthroughObjectPropertyOption is a functional option for configuring a [PassthroughObjectProperty].
type PassthroughObjectPropertyOption func(*PassthroughObjectProperty)

// PassthroughObjectProperty represents an Elasticsearch "passthrough" object field mapping.
// Use [NewPassthroughObjectProperty] to construct one with functional options.
type PassthroughObjectProperty struct {
	// Properties holds the child field mappings.
	Properties map[string]MappingProperty
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

// ---------------------------------------------------------------------------
// IP
// ---------------------------------------------------------------------------

// IpProperty represents an Elasticsearch "ip" field mapping for IPv4 and IPv6 addresses.
type IpProperty struct{}

// ESTypeName returns the Elasticsearch type name for an ip property.
func (IpProperty) ESTypeName() string { return "ip" }

// NewIpProperty creates a new [IpProperty].
func NewIpProperty() IpProperty {
	return IpProperty{}
}

// ---------------------------------------------------------------------------
// Binary
// ---------------------------------------------------------------------------

// BinaryProperty represents an Elasticsearch "binary" field mapping.
// Values are stored as Base64 encoded strings and are not searchable by default.
type BinaryProperty struct{}

// ESTypeName returns the Elasticsearch type name for a binary property.
func (BinaryProperty) ESTypeName() string { return "binary" }

// NewBinaryProperty creates a new [BinaryProperty].
func NewBinaryProperty() BinaryProperty {
	return BinaryProperty{}
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
	Analyzer *Analyzer
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

// HistogramProperty represents an Elasticsearch "histogram" field mapping
// for pre-aggregated numerical data.
type HistogramProperty struct{}

// ESTypeName returns the Elasticsearch type name for a histogram property.
func (HistogramProperty) ESTypeName() string { return "histogram" }

// NewHistogramProperty creates a new [HistogramProperty].
func NewHistogramProperty() HistogramProperty {
	return HistogramProperty{}
}

// ---------------------------------------------------------------------------
// Version
// ---------------------------------------------------------------------------

// VersionProperty represents an Elasticsearch "version" field mapping
// for software version values following semver rules.
type VersionProperty struct{}

// ESTypeName returns the Elasticsearch type name for a version property.
func (VersionProperty) ESTypeName() string { return "version" }

// NewVersionProperty creates a new [VersionProperty].
func NewVersionProperty() VersionProperty {
	return VersionProperty{}
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

// ---------------------------------------------------------------------------
// Sparse Vector
// ---------------------------------------------------------------------------

// SparseVectorProperty represents an Elasticsearch "sparse_vector" field mapping.
type SparseVectorProperty struct{}

// ESTypeName returns the Elasticsearch type name for a sparse_vector property.
func (SparseVectorProperty) ESTypeName() string { return "sparse_vector" }

// NewSparseVectorProperty creates a new [SparseVectorProperty].
func NewSparseVectorProperty() SparseVectorProperty {
	return SparseVectorProperty{}
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

// RankFeaturesProperty represents an Elasticsearch "rank_features" field mapping
// for multiple named rank features in a single field.
type RankFeaturesProperty struct{}

// ESTypeName returns the Elasticsearch type name for a rank_features property.
func (RankFeaturesProperty) ESTypeName() string { return "rank_features" }

// NewRankFeaturesProperty creates a new [RankFeaturesProperty].
func NewRankFeaturesProperty() RankFeaturesProperty {
	return RankFeaturesProperty{}
}

// RankVectorProperty represents an Elasticsearch "rank_vectors" field mapping.
type RankVectorProperty struct{}

// ESTypeName returns the Elasticsearch type name for a rank_vectors property.
func (RankVectorProperty) ESTypeName() string { return "rank_vectors" }

// NewRankVectorProperty creates a new [RankVectorProperty].
func NewRankVectorProperty() RankVectorProperty {
	return RankVectorProperty{}
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

// ---------------------------------------------------------------------------
// Murmur3 Hash
// ---------------------------------------------------------------------------

// Murmur3HashProperty represents an Elasticsearch "murmur3" field mapping
// that computes and stores a murmur3 hash of the field value.
type Murmur3HashProperty struct{}

// ESTypeName returns the Elasticsearch type name for a murmur3 property.
func (Murmur3HashProperty) ESTypeName() string { return "murmur3" }

// NewMurmur3HashProperty creates a new [Murmur3HashProperty].
func NewMurmur3HashProperty() Murmur3HashProperty {
	return Murmur3HashProperty{}
}

// ---------------------------------------------------------------------------
// ICU Collation Keyword
// ---------------------------------------------------------------------------

// IcuCollationProperty represents an Elasticsearch "icu_collation_keyword" field mapping
// provided by the ICU analysis plugin. The field indexes text as a keyword
// using ICU collation rules for locale-sensitive sorting and comparison.
type IcuCollationProperty struct{}

// ESTypeName returns the Elasticsearch type name for an icu_collation_keyword property.
func (IcuCollationProperty) ESTypeName() string { return "icu_collation_keyword" }

// NewIcuCollationProperty creates a new [IcuCollationProperty].
func NewIcuCollationProperty() IcuCollationProperty {
	return IcuCollationProperty{}
}

// ---------------------------------------------------------------------------
// Dynamic
// ---------------------------------------------------------------------------

// DynamicProperty represents an Elasticsearch "{dynamic_type}" field mapping
// used in dynamic templates.
type DynamicProperty struct{}

// ESTypeName returns the Elasticsearch type name for a dynamic property.
func (DynamicProperty) ESTypeName() string { return "{dynamic_type}" }

// NewDynamicProperty creates a new [DynamicProperty].
func NewDynamicProperty() DynamicProperty {
	return DynamicProperty{}
}
