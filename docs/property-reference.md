# Property Reference

[English](property-reference.md) | [日本語](property-reference.ja.md)

This document lists all property constructors and their functional options provided by `esv8` and `esv9` packages. All constructors follow the functional-option pattern:

```go
prop := esv8.NewTextProperty(
    esv8.WithTextAnalyzer("standard"),
    esv8.WithTextStore(true),
)
```

Unless noted otherwise, every property listed here is available in both `esv8` and `esv9` with identical signatures.

---

## Table of Contents

- [Text Family](#text-family)
- [Numeric](#numeric)
- [Date and Boolean](#date-and-boolean)
- [Geographic](#geographic)
- [Range](#range)
- [Object and Nested](#object-and-nested)
- [Join](#join)
- [Network](#network)
- [Vector](#vector)
- [Ranking](#ranking)
- [Special](#special)
- [Plugin-Dependent](#plugin-dependent)

---

## Text Family

### NewTextProperty

Full-text field with analyzer support.

```go
esv8.NewTextProperty(opts ...TextPropertyOption) *types.TextProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithTextAnalyzer` | `v string` | Sets the analyzer for indexing |
| `WithTextSearchAnalyzer` | `v string` | Sets the analyzer for search queries |
| `WithTextSearchQuoteAnalyzer` | `v string` | Sets the analyzer for phrase queries |
| `WithTextFielddata` | `v bool` | Enables fielddata for sorting/aggregations on text |
| `WithTextIndex` | `v bool` | Whether the field is searchable |
| `WithTextStore` | `v bool` | Whether the field value is stored separately |
| `WithTextNorms` | `v bool` | Whether norms are stored for scoring |
| `WithTextSimilarity` | `v string` | Sets the similarity algorithm (e.g. `"BM25"`) |
| `WithTextIndexPhrases` | `v bool` | Whether two-term word combinations are indexed |
| `WithTextPositionIncrementGap` | `v int` | Number of fake term positions between array values |
| `WithTextRawKeyword` | `ignoreAbove int` | Adds a `.keyword` sub-field with the given `ignore_above` |
| `WithTextFields` | `fields map[string]types.Property` | Sets custom multi-fields |

### NewKeywordProperty

Exact-value string field for filtering, sorting, and aggregations.

```go
esv8.NewKeywordProperty(opts ...KeywordPropertyOption) *types.KeywordProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithKeywordIgnoreAbove` | `v int` | Maximum string length; longer values are not indexed |
| `WithKeywordDocValues` | `v bool` | Whether doc values are enabled |
| `WithKeywordIndex` | `v bool` | Whether the field is searchable |
| `WithKeywordStore` | `v bool` | Whether the field value is stored separately |
| `WithKeywordNullValue` | `v string` | Value substituted for `null` |
| `WithKeywordNormalizer` | `v string` | Normalizer applied before indexing |
| `WithKeywordNorms` | `v bool` | Whether norms are stored for scoring |
| `WithKeywordSimilarity` | `v string` | Similarity algorithm |
| `WithKeywordEagerGlobalOrdinals` | `v bool` | Whether to eagerly load global ordinals |
| `WithKeywordSplitQueriesOnWhitespace` | `v bool` | Whether to split queries on whitespace |

### NewConstantKeywordProperty

Keyword field where all documents share the same value.

```go
esv8.NewConstantKeywordProperty(opts ...ConstantKeywordPropertyOption) *types.ConstantKeywordProperty
```

No options. Value is set on the first indexed document.

### NewCountedKeywordProperty

Keyword field that also tracks the count of each term.

```go
esv8.NewCountedKeywordProperty(opts ...CountedKeywordPropertyOption) *types.CountedKeywordProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithCountedKeywordIndex` | `v bool` | Whether the field is searchable |

### NewWildcardProperty

Keyword-like field optimized for wildcard and regexp queries.

```go
esv8.NewWildcardProperty(opts ...WildcardPropertyOption) *types.WildcardProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithWildcardIgnoreAbove` | `v int` | Maximum string length |
| `WithWildcardNullValue` | `v string` | Value substituted for `null` |

### NewMatchOnlyTextProperty

Text field optimized for match queries that does not store positions or norms.

```go
esv8.NewMatchOnlyTextProperty(opts ...MatchOnlyTextPropertyOption) *types.MatchOnlyTextProperty
```

No options.

### NewCompletionProperty

Autocomplete suggestion field.

```go
esv8.NewCompletionProperty(opts ...CompletionPropertyOption) *types.CompletionProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithCompletionAnalyzer` | `v string` | Analyzer for indexing |
| `WithCompletionSearchAnalyzer` | `v string` | Analyzer for search |
| `WithCompletionMaxInputLength` | `v int` | Maximum length of a single input |
| `WithCompletionPreservePositionIncrements` | `v bool` | Whether to preserve position increments |
| `WithCompletionPreserveSeparators` | `v bool` | Whether to preserve separators |

### NewSearchAsYouTypeProperty

Field type for search-as-you-type autocomplete.

```go
esv8.NewSearchAsYouTypeProperty(opts ...SearchAsYouTypePropertyOption) *types.SearchAsYouTypeProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithSearchAsYouTypeAnalyzer` | `v string` | Analyzer for indexing |
| `WithSearchAsYouTypeSearchAnalyzer` | `v string` | Analyzer for search |
| `WithSearchAsYouTypeSearchQuoteAnalyzer` | `v string` | Analyzer for phrase queries |
| `WithSearchAsYouTypeMaxShingleSize` | `v int` | Maximum shingle size (2-4) |
| `WithSearchAsYouTypeIndex` | `v bool` | Whether the field is searchable |
| `WithSearchAsYouTypeStore` | `v bool` | Whether the field value is stored separately |
| `WithSearchAsYouTypeNorms` | `v bool` | Whether norms are stored |
| `WithSearchAsYouTypeSimilarity` | `v string` | Similarity algorithm |

---

## Numeric

All numeric property types share a common set of options (coerce, doc_values, ignore_malformed, index, store, null_value). Exceptions are noted per type.

### NewIntegerNumberProperty

32-bit signed integer.

```go
esv8.NewIntegerNumberProperty(opts ...IntegerNumberPropertyOption) *types.IntegerNumberProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithIntegerNumberCoerce` | `v bool` | Whether to coerce values to the correct type |
| `WithIntegerNumberDocValues` | `v bool` | Whether doc values are enabled |
| `WithIntegerNumberIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithIntegerNumberIndex` | `v bool` | Whether the field is searchable |
| `WithIntegerNumberStore` | `v bool` | Whether the field value is stored separately |
| `WithIntegerNumberNullValue` | `v int` | Value substituted for `null` |

### NewLongNumberProperty

64-bit signed integer.

```go
esv8.NewLongNumberProperty(opts ...LongNumberPropertyOption) *types.LongNumberProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithLongNumberCoerce` | `v bool` | Whether to coerce values to the correct type |
| `WithLongNumberDocValues` | `v bool` | Whether doc values are enabled |
| `WithLongNumberIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithLongNumberIndex` | `v bool` | Whether the field is searchable |
| `WithLongNumberStore` | `v bool` | Whether the field value is stored separately |
| `WithLongNumberNullValue` | `v int64` | Value substituted for `null` |

### NewShortNumberProperty

16-bit signed integer.

```go
esv8.NewShortNumberProperty(opts ...ShortNumberPropertyOption) *types.ShortNumberProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithShortNumberCoerce` | `v bool` | Whether to coerce values to the correct type |
| `WithShortNumberDocValues` | `v bool` | Whether doc values are enabled |
| `WithShortNumberIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithShortNumberIndex` | `v bool` | Whether the field is searchable |
| `WithShortNumberStore` | `v bool` | Whether the field value is stored separately |
| `WithShortNumberNullValue` | `v int` | Value substituted for `null` |

### NewByteNumberProperty

8-bit signed integer.

```go
esv8.NewByteNumberProperty(opts ...ByteNumberPropertyOption) *types.ByteNumberProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithByteNumberCoerce` | `v bool` | Whether to coerce values to the correct type |
| `WithByteNumberDocValues` | `v bool` | Whether doc values are enabled |
| `WithByteNumberIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithByteNumberIndex` | `v bool` | Whether the field is searchable |
| `WithByteNumberStore` | `v bool` | Whether the field value is stored separately |
| `WithByteNumberNullValue` | `v byte` | Value substituted for `null` |

### NewDoubleNumberProperty

64-bit IEEE 754 floating point.

```go
esv8.NewDoubleNumberProperty(opts ...DoubleNumberPropertyOption) *types.DoubleNumberProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithDoubleNumberCoerce` | `v bool` | Whether to coerce values to the correct type |
| `WithDoubleNumberDocValues` | `v bool` | Whether doc values are enabled |
| `WithDoubleNumberIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithDoubleNumberIndex` | `v bool` | Whether the field is searchable |
| `WithDoubleNumberStore` | `v bool` | Whether the field value is stored separately |
| `WithDoubleNumberNullValue` | `v float64` | Value substituted for `null` |

### NewFloatNumberProperty

32-bit IEEE 754 floating point.

```go
esv8.NewFloatNumberProperty(opts ...FloatNumberPropertyOption) *types.FloatNumberProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithFloatNumberCoerce` | `v bool` | Whether to coerce values to the correct type |
| `WithFloatNumberDocValues` | `v bool` | Whether doc values are enabled |
| `WithFloatNumberIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithFloatNumberIndex` | `v bool` | Whether the field is searchable |
| `WithFloatNumberStore` | `v bool` | Whether the field value is stored separately |
| `WithFloatNumberNullValue` | `v float32` | Value substituted for `null` |

### NewHalfFloatNumberProperty

16-bit IEEE 754 floating point.

```go
esv8.NewHalfFloatNumberProperty(opts ...HalfFloatNumberPropertyOption) *types.HalfFloatNumberProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithHalfFloatNumberCoerce` | `v bool` | Whether to coerce values to the correct type |
| `WithHalfFloatNumberDocValues` | `v bool` | Whether doc values are enabled |
| `WithHalfFloatNumberIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithHalfFloatNumberIndex` | `v bool` | Whether the field is searchable |
| `WithHalfFloatNumberStore` | `v bool` | Whether the field value is stored separately |
| `WithHalfFloatNumberNullValue` | `v float32` | Value substituted for `null` |

### NewUnsignedLongNumberProperty

Unsigned 64-bit integer (0 to 2^64-1).

```go
esv8.NewUnsignedLongNumberProperty(opts ...UnsignedLongNumberPropertyOption) *types.UnsignedLongNumberProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithUnsignedLongNumberDocValues` | `v bool` | Whether doc values are enabled |
| `WithUnsignedLongNumberIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithUnsignedLongNumberIndex` | `v bool` | Whether the field is searchable |
| `WithUnsignedLongNumberStore` | `v bool` | Whether the field value is stored separately |
| `WithUnsignedLongNumberNullValue` | `v uint64` | Value substituted for `null` |

> **Note:** `unsigned_long` does not support the `coerce` parameter.

### NewScaledFloatNumberProperty

Floating-point number stored as a scaled `long` for compact storage.

```go
esv8.NewScaledFloatNumberProperty(opts ...ScaledFloatNumberPropertyOption) *types.ScaledFloatNumberProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithScaledFloatNumberScalingFactor` | `v float64` | Scaling factor (required) |
| `WithScaledFloatNumberCoerce` | `v bool` | Whether to coerce values to the correct type |
| `WithScaledFloatNumberDocValues` | `v bool` | Whether doc values are enabled |
| `WithScaledFloatNumberIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithScaledFloatNumberIndex` | `v bool` | Whether the field is searchable |
| `WithScaledFloatNumberStore` | `v bool` | Whether the field value is stored separately |
| `WithScaledFloatNumberNullValue` | `v float64` | Value substituted for `null` |

---

## Date and Boolean

### NewDateProperty

Date/time field.

```go
esv8.NewDateProperty(opts ...DatePropertyOption) *types.DateProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithDateFormat` | `formats ...estype.DateFormat` | Accepted date formats (joined with `\|\|`) |
| `WithDateDocValues` | `v bool` | Whether doc values are enabled |
| `WithDateIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithDateIndex` | `v bool` | Whether the field is searchable |
| `WithDateStore` | `v bool` | Whether the field value is stored separately |
| `WithDateLocale` | `v string` | Locale for parsing dates (e.g. `"en"`) |

### NewDateNanosProperty

Date/time field with nanosecond resolution.

```go
esv8.NewDateNanosProperty(opts ...DateNanosPropertyOption) *types.DateNanosProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithDateNanosFormat` | `formats ...estype.DateFormat` | Accepted date formats |
| `WithDateNanosDocValues` | `v bool` | Whether doc values are enabled |
| `WithDateNanosIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithDateNanosIndex` | `v bool` | Whether the field is searchable |
| `WithDateNanosStore` | `v bool` | Whether the field value is stored separately |

### NewBooleanProperty

Boolean (`true`/`false`) field.

```go
esv8.NewBooleanProperty(opts ...BooleanPropertyOption) *types.BooleanProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithBooleanDocValues` | `v bool` | Whether doc values are enabled |
| `WithBooleanIndex` | `v bool` | Whether the field is searchable |
| `WithBooleanStore` | `v bool` | Whether the field value is stored separately |
| `WithBooleanNullValue` | `v bool` | Value substituted for `null` |

---

## Geographic

### NewGeoPointProperty

Latitude/longitude point.

```go
esv8.NewGeoPointProperty(opts ...GeoPointPropertyOption) *types.GeoPointProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithGeoPointIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithGeoPointIgnoreZValue` | `v bool` | Whether to ignore the Z value |
| `WithGeoPointDocValues` | `v bool` | Whether doc values are enabled |
| `WithGeoPointIndex` | `v bool` | Whether the field is searchable |
| `WithGeoPointStore` | `v bool` | Whether the field value is stored separately |

### NewGeoShapeProperty

Arbitrary GeoJSON geometry.

```go
esv8.NewGeoShapeProperty(opts ...GeoShapePropertyOption) *types.GeoShapeProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithGeoShapeCoerce` | `v bool` | Whether to coerce unclosed polygons |
| `WithGeoShapeIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithGeoShapeIgnoreZValue` | `v bool` | Whether to ignore the Z value |
| `WithGeoShapeDocValues` | `v bool` | Whether doc values are enabled |
| `WithGeoShapeIndex` | `v bool` | Whether the field is searchable |
| `WithGeoShapeStore` | `v bool` | Whether the field value is stored separately |

### NewShapeProperty

Arbitrary Cartesian geometry (non-geographic).

```go
esv8.NewShapeProperty(opts ...ShapePropertyOption) *types.ShapeProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithShapeCoerce` | `v bool` | Whether to coerce unclosed polygons |
| `WithShapeIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithShapeIgnoreZValue` | `v bool` | Whether to ignore the Z value |
| `WithShapeDocValues` | `v bool` | Whether doc values are enabled |

### NewPointProperty

Cartesian (x, y) point.

```go
esv8.NewPointProperty(opts ...PointPropertyOption) *types.PointProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithPointIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithPointIgnoreZValue` | `v bool` | Whether to ignore the Z value |
| `WithPointDocValues` | `v bool` | Whether doc values are enabled |
| `WithPointStore` | `v bool` | Whether the field value is stored separately |
| `WithPointNullValue` | `v string` | Value substituted for `null` (WKT point) |

---

## Range

All range property types share a common set of options: coerce, doc_values, index, store. The date range type additionally supports a format option.

### NewIntegerRangeProperty

```go
esv8.NewIntegerRangeProperty(opts ...IntegerRangePropertyOption) *types.IntegerRangeProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithIntegerRangeCoerce` | `v bool` | Whether to coerce values |
| `WithIntegerRangeDocValues` | `v bool` | Whether doc values are enabled |
| `WithIntegerRangeIndex` | `v bool` | Whether the field is searchable |
| `WithIntegerRangeStore` | `v bool` | Whether the field value is stored separately |

### NewLongRangeProperty

```go
esv8.NewLongRangeProperty(opts ...LongRangePropertyOption) *types.LongRangeProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithLongRangeCoerce` | `v bool` | Whether to coerce values |
| `WithLongRangeDocValues` | `v bool` | Whether doc values are enabled |
| `WithLongRangeIndex` | `v bool` | Whether the field is searchable |
| `WithLongRangeStore` | `v bool` | Whether the field value is stored separately |

### NewFloatRangeProperty

```go
esv8.NewFloatRangeProperty(opts ...FloatRangePropertyOption) *types.FloatRangeProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithFloatRangeCoerce` | `v bool` | Whether to coerce values |
| `WithFloatRangeDocValues` | `v bool` | Whether doc values are enabled |
| `WithFloatRangeIndex` | `v bool` | Whether the field is searchable |
| `WithFloatRangeStore` | `v bool` | Whether the field value is stored separately |

### NewDoubleRangeProperty

```go
esv8.NewDoubleRangeProperty(opts ...DoubleRangePropertyOption) *types.DoubleRangeProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithDoubleRangeCoerce` | `v bool` | Whether to coerce values |
| `WithDoubleRangeDocValues` | `v bool` | Whether doc values are enabled |
| `WithDoubleRangeIndex` | `v bool` | Whether the field is searchable |
| `WithDoubleRangeStore` | `v bool` | Whether the field value is stored separately |

### NewDateRangeProperty

```go
esv8.NewDateRangeProperty(opts ...DateRangePropertyOption) *types.DateRangeProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithDateRangeFormat` | `formats ...estype.DateFormat` | Accepted date formats |
| `WithDateRangeCoerce` | `v bool` | Whether to coerce values |
| `WithDateRangeDocValues` | `v bool` | Whether doc values are enabled |
| `WithDateRangeIndex` | `v bool` | Whether the field is searchable |
| `WithDateRangeStore` | `v bool` | Whether the field value is stored separately |

### NewIpRangeProperty

```go
esv8.NewIpRangeProperty(opts ...IpRangePropertyOption) *types.IpRangeProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithIpRangeCoerce` | `v bool` | Whether to coerce values |
| `WithIpRangeDocValues` | `v bool` | Whether doc values are enabled |
| `WithIpRangeIndex` | `v bool` | Whether the field is searchable |
| `WithIpRangeStore` | `v bool` | Whether the field value is stored separately |

---

## Object and Nested

### NewObjectProperty

JSON object (flat structure, no independent querying of child fields).

```go
esv8.NewObjectProperty(opts ...ObjectPropertyOption) *types.ObjectProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithObjectProperties` | `v map[string]types.Property` | Child field mappings |
| `WithObjectEnabled` | `v bool` | Whether the object is enabled for indexing |

### NewNestedProperty

JSON object where child fields can be queried independently.

```go
esv8.NewNestedProperty(opts ...NestedPropertyOption) *types.NestedProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithNestedProperties` | `v map[string]types.Property` | Child field mappings |
| `WithNestedEnabled` | `v bool` | Whether the nested object is enabled |
| `WithNestedIncludeInParent` | `v bool` | Whether to include nested fields in parent |
| `WithNestedIncludeInRoot` | `v bool` | Whether to include nested fields in root |

### NewFlattenedProperty

Entire JSON object mapped as a single field.

```go
esv8.NewFlattenedProperty(opts ...FlattenedPropertyOption) *types.FlattenedProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithFlattenedDepthLimit` | `v int` | Maximum nesting depth |
| `WithFlattenedDocValues` | `v bool` | Whether doc values are enabled |
| `WithFlattenedIndex` | `v bool` | Whether the field is searchable |
| `WithFlattenedIgnoreAbove` | `v int` | Maximum string length for leaf values |
| `WithFlattenedNullValue` | `v string` | Value substituted for `null` |
| `WithFlattenedEagerGlobalOrdinals` | `v bool` | Whether to eagerly load global ordinals |
| `WithFlattenedSimilarity` | `v string` | Similarity algorithm |
| `WithFlattenedSplitQueriesOnWhitespace` | `v bool` | Whether to split queries on whitespace |

### NewPassthroughObjectProperty

Object where fields are mapped to the parent level (used in data streams).

```go
esv8.NewPassthroughObjectProperty(opts ...PassthroughObjectPropertyOption) *types.PassthroughObjectProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithPassthroughObjectProperties` | `v map[string]types.Property` | Child field mappings |
| `WithPassthroughObjectEnabled` | `v bool` | Whether the object is enabled |
| `WithPassthroughObjectPriority` | `v int` | Priority for conflicting field names |
| `WithPassthroughObjectTimeSeriesDimension` | `v bool` | Whether this is a time series dimension |

---

## Join

### NewJoinProperty

Parent-child relationship field.

```go
esv8.NewJoinProperty(opts ...JoinPropertyOption) *types.JoinProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithJoinRelations` | `v map[string][]string` | Parent-child relation definitions |
| `WithJoinEagerGlobalOrdinals` | `v bool` | Whether to eagerly load global ordinals |

---

## Network

### NewIpProperty

IPv4 or IPv6 address.

```go
esv8.NewIpProperty(opts ...IpPropertyOption) *types.IpProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithIpDocValues` | `v bool` | Whether doc values are enabled |
| `WithIpIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithIpIndex` | `v bool` | Whether the field is searchable |
| `WithIpStore` | `v bool` | Whether the field value is stored separately |
| `WithIpNullValue` | `v string` | Value substituted for `null` |

---

## Vector

### NewDenseVectorProperty

Dense floating-point vector for kNN search.

```go
esv8.NewDenseVectorProperty(opts ...DenseVectorPropertyOption) *types.DenseVectorProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithDenseVectorDims` | `v int` | Number of dimensions |
| `WithDenseVectorIndex` | `v bool` | Whether to index for kNN search |

### NewSparseVectorProperty

Sparse floating-point vector for term-based ranking.

```go
esv8.NewSparseVectorProperty(opts ...SparseVectorPropertyOption) *types.SparseVectorProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithSparseVectorStore` | `v bool` | Whether the field value is stored separately |

### NewRankVectorProperty

Fixed-length float vector for rank-based scoring.

```go
esv8.NewRankVectorProperty(opts ...RankVectorPropertyOption) *types.RankVectorProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithRankVectorDims` | `v int` | Number of dimensions |

---

## Ranking

### NewRankFeatureProperty

Numeric feature used to boost relevance scoring.

```go
esv8.NewRankFeatureProperty(opts ...RankFeaturePropertyOption) *types.RankFeatureProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithRankFeaturePositiveScoreImpact` | `v bool` | Whether higher values boost relevance |

### NewRankFeaturesProperty

Map of named numeric features for boosting relevance scoring.

```go
esv8.NewRankFeaturesProperty(opts ...RankFeaturesPropertyOption) *types.RankFeaturesProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithRankFeaturesPositiveScoreImpact` | `v bool` | Whether higher values boost relevance |

---

## Special

### NewBinaryProperty

Base64-encoded binary data.

```go
esv8.NewBinaryProperty(opts ...BinaryPropertyOption) *types.BinaryProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithBinaryDocValues` | `v bool` | Whether doc values are enabled |
| `WithBinaryStore` | `v bool` | Whether the field value is stored separately |

### NewTokenCountProperty

Integer count of analyzed tokens.

```go
esv8.NewTokenCountProperty(opts ...TokenCountPropertyOption) *types.TokenCountProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithTokenCountAnalyzer` | `v string` | Analyzer for counting tokens |
| `WithTokenCountDocValues` | `v bool` | Whether doc values are enabled |
| `WithTokenCountIndex` | `v bool` | Whether the field is searchable |
| `WithTokenCountStore` | `v bool` | Whether the field value is stored separately |
| `WithTokenCountEnablePositionIncrements` | `v bool` | Whether to count position increments |

### NewPercolatorProperty

Stores a query for use with the percolate query.

```go
esv8.NewPercolatorProperty(opts ...PercolatorPropertyOption) *types.PercolatorProperty
```

No options.

### NewFieldAliasProperty

Alternative name for an existing field.

```go
esv8.NewFieldAliasProperty(opts ...FieldAliasPropertyOption) *types.FieldAliasProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithFieldAliasPath` | `v string` | Path to the target field |

### NewHistogramProperty

Pre-aggregated histogram data.

```go
esv8.NewHistogramProperty(opts ...HistogramPropertyOption) *types.HistogramProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithHistogramIgnoreMalformed` | `v bool` | Whether to ignore malformed values |

### NewExponentialHistogramProperty (v9 only)

Pre-aggregated exponential histogram data.

```go
esv9.NewExponentialHistogramProperty(opts ...ExponentialHistogramPropertyOption) *types.ExponentialHistogramProperty
```

No options. Available only in `esv9`.

### NewVersionProperty

Software version string with semantic version ordering.

```go
esv8.NewVersionProperty(opts ...VersionPropertyOption) *types.VersionProperty
```

No options.

### NewAggregateMetricDoubleProperty

Pre-aggregated metric values.

```go
esv8.NewAggregateMetricDoubleProperty(opts ...AggregateMetricDoublePropertyOption) *types.AggregateMetricDoubleProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithAggregateMetricDoubleDefaultMetric` | `v string` | Default metric for queries |
| `WithAggregateMetricDoubleMetrics` | `v []string` | List of stored metrics (e.g. `["min","max","sum","value_count"]`) |
| `WithAggregateMetricDoubleIgnoreMalformed` | `v bool` | Whether to ignore malformed values |

### NewSemanticTextProperty

Field for ML-powered semantic search using inference endpoints.

```go
esv8.NewSemanticTextProperty(opts ...SemanticTextPropertyOption) *types.SemanticTextProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithSemanticTextInferenceId` | `v string` | Inference endpoint ID for indexing |
| `WithSemanticTextSearchInferenceId` | `v string` | Inference endpoint ID for search |

### NewDynamicProperty

Template-based property for dynamic field mappings.

```go
esv8.NewDynamicProperty(opts ...DynamicPropertyOption) *types.DynamicProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithDynamicAnalyzer` | `v string` | Analyzer for text-like types |
| `WithDynamicSearchAnalyzer` | `v string` | Search analyzer |
| `WithDynamicCoerce` | `v bool` | Whether to coerce values |
| `WithDynamicDocValues` | `v bool` | Whether doc values are enabled |
| `WithDynamicEnabled` | `v bool` | Whether the field is enabled |
| `WithDynamicFormat` | `v string` | Date format string |
| `WithDynamicIgnoreMalformed` | `v bool` | Whether to ignore malformed values |
| `WithDynamicIndex` | `v bool` | Whether the field is searchable |
| `WithDynamicStore` | `v bool` | Whether the field value is stored separately |
| `WithDynamicNorms` | `v bool` | Whether norms are stored |
| `WithDynamicLocale` | `v string` | Locale for date parsing |

---

## Plugin-Dependent

These property types require specific Elasticsearch plugins to be installed.

### NewMurmur3HashProperty

Stores a murmur3 hash of the field value. Requires the `mapper-murmur3` plugin.

```go
esv8.NewMurmur3HashProperty(opts ...Murmur3HashPropertyOption) *types.Murmur3HashProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithMurmur3HashDocValues` | `v bool` | Whether doc values are enabled |
| `WithMurmur3HashStore` | `v bool` | Whether the field value is stored separately |

### NewIcuCollationProperty

Keyword field with ICU collation-based sorting. Requires the `analysis-icu` plugin.

```go
esv8.NewIcuCollationProperty(opts ...IcuCollationPropertyOption) *types.IcuCollationProperty
```

| Option | Parameter | Description |
|--------|-----------|-------------|
| `WithIcuCollationLanguage` | `v string` | Language code (e.g. `"en"`) |
| `WithIcuCollationCountry` | `v string` | Country code (e.g. `"US"`) |
| `WithIcuCollationDocValues` | `v bool` | Whether doc values are enabled |
| `WithIcuCollationIndex` | `v bool` | Whether the field is searchable |
| `WithIcuCollationStore` | `v bool` | Whether the field value is stored separately |
| `WithIcuCollationNullValue` | `v string` | Value substituted for `null` |
| `WithIcuCollationNorms` | `v bool` | Whether norms are stored |
| `WithIcuCollationRules` | `v string` | ICU collation rules string |
| `WithIcuCollationVariant` | `v string` | Collation variant |
| `WithIcuCollationCaseLevel` | `v bool` | Whether case-level comparison is enabled |
| `WithIcuCollationNumeric` | `v bool` | Whether numeric collation is enabled |
| `WithIcuCollationHiraganaQuaternaryMode` | `v bool` | Whether Hiragana quaternary mode is enabled |
| `WithIcuCollationVariableTop` | `v string` | Variable top setting for collation |

---

## Version Differences (v8 vs v9)

| Property | v8 | v9 |
|----------|----|----|
| All listed above (except noted) | Yes | Yes |
| `NewExponentialHistogramProperty` | -- | Yes |

All option function signatures are identical across `esv8` and `esv9`. Code written for `esv8` works unchanged in `esv9`, with the addition of `NewExponentialHistogramProperty` available only in `esv9`.
