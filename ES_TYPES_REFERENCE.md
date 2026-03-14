# Elasticsearch Typed API v8 - Go Type Definitions

## Overview
This document provides a comprehensive reference for the Elasticsearch go-elasticsearch v8 library's type definitions, specifically for building type-safe query builders.

**Location**: `/home/runner/go/pkg/mod/github.com/elastic/go-elasticsearch/v8@v8.19.3/typedapi/types/`

---

## 1. Query Type

**File**: `query.go`

The `Query` struct is a comprehensive union type that represents all possible Elasticsearch query types. It uses a map-based approach to handle multiple query types within a single structure.

### Query Struct Definition
```go
type Query struct {
    AdditionalQueryProperty map[string]json.RawMessage `json:"-"`
    
    // Core query types (partial list):
    Bool              *BoolQuery
    Boosting          *BoostingQuery
    CombinedFields    *CombinedFieldsQuery
    Common            map[string]CommonTermsQuery
    ConstantScore     *ConstantScoreQuery
    DisMax            *DisMaxQuery
    DistanceFeature   DistanceFeatureQuery
    Exists            *ExistsQuery
    FunctionScore     *FunctionScoreQuery
    Fuzzy             map[string]FuzzyQuery
    GeoBoundingBox    *GeoBoundingBoxQuery
    GeoDistance       *GeoDistanceQuery
    GeoGrid           map[string]GeoGridQuery
    GeoPolygon        *GeoPolygonQuery
    GeoShape          *GeoShapeQuery
    HasChild          *HasChildQuery
    HasParent         *HasParentQuery
    Ids               *IdsQuery
    Intervals         map[string]IntervalsQuery
    Knn               *KnnQuery
    Match             map[string]MatchQuery
    MatchAll          *MatchAllQuery
    MatchBoolPrefix   map[string]MatchBoolPrefixQuery
    MatchNone         *MatchNoneQuery
    MatchPhrase       map[string]MatchPhraseQuery
    MatchPhrasePrefix map[string]MatchPhrasePrefixQuery
    MoreLikeThis      *MoreLikeThisQuery
    MultiMatch        *MultiMatchQuery
    Nested            *NestedQuery
    ParentId          *ParentIdQuery
    Percolate         *PercolateQuery
    Pinned            *PinnedQuery
    Prefix            map[string]PrefixQuery
    QueryString       *QueryStringQuery
    Range             map[string]RangeQuery
    RankFeature       *RankFeatureQuery
    Regexp            map[string]RegexpQuery
    Rule              *RuleQuery
    Script            *ScriptQuery
    ScriptScore       *ScriptScoreQuery
    Semantic          *SemanticQuery
    Shape             *ShapeQuery
    SimpleQueryString *SimpleQueryStringQuery
    SpanContaining    *SpanContainingQuery
    SpanFieldMasking  *SpanFieldMaskingQuery
    SpanFirst         *SpanFirstQuery
    SpanMulti         *SpanMultiTermQuery
    SpanNear          *SpanNearQuery
    SpanNot           *SpanNotQuery
    SpanOr            *SpanOrQuery
    SpanTerm          map[string]SpanTermQuery
    SpanWithin        *SpanWithinQuery
    SparseVector      *SparseVectorQuery
    Term              map[string]TermQuery
    Terms             *TermsQuery
    TermsSet          map[string]TermsSetQuery
    TextExpansion     map[string]TextExpansionQuery
    Type              *TypeQuery
    WeightedTokens    map[string]WeightedTokensQuery
    Wildcard          map[string]WildcardQuery
    Wrapper           *WrapperQuery
}
```

### Key Features:
- **Union Type Pattern**: Uses `json:"-"` for maps to handle field-per-query-type pattern
- **Flexible Nesting**: Supports nested queries and complex boolean combinations
- **Additional Properties**: `AdditionalQueryProperty` map allows for future/custom query types

---

## 2. TermQuery Type

**File**: `termquery.go`

Matches documents with an exact term in a provided field.

### TermQuery Struct Definition
```go
type TermQuery struct {
    Boost          *float32   `json:"boost,omitempty"`
    CaseInsensitive *bool      `json:"case_insensitive,omitempty"`
    QueryName_     *string    `json:"_name,omitempty"`
    Value          FieldValue `json:"value"`  // Required field
}
```

### Fields:
- **Value** (FieldValue): The term value to search for - **REQUIRED**
- **Boost** (*float32): Relevance score multiplier (default 1.0)
- **CaseInsensitive** (*bool): Enable ASCII case-insensitive matching
- **QueryName_** (*string): Optional query name for debugging

---

## 3. TermsQuery Type

**File**: `termsquery.go`

Matches documents with one or more exact terms in a provided field.

### TermsQuery Struct Definition
```go
type TermsQuery struct {
    Boost      *float32                   `json:"boost,omitempty"`
    QueryName_ *string                    `json:"_name,omitempty"`
    TermsQuery map[string]TermsQueryField `json:"-"`  // Field-based query map
}
```

### Fields:
- **TermsQuery** (map[string]TermsQueryField): Maps field names to their term values
- **Boost** (*float32): Relevance score multiplier
- **QueryName_** (*string): Optional query name

### Supporting Type: TermsQueryField
```go
// Union type holding:
// - []FieldValue (array of values)
// - TermsLookup (lookup from another document)
type TermsQueryField any
```

---

## 4. FieldValue Type

**File**: `fieldvalue.go`

A union type representing any value that can be stored in a field.

```go
// Supports:
// - int64
// - Float64
// - string
// - bool
// - nil
// - json.RawMessage
type FieldValue any
```

**Usage**: Used across term queries, aggregations, and other contexts requiring flexible value types.

---

## 5. Sort Types

### 5a. Sort Type (Alias)
**File**: `sort.go`

```go
// Type alias for a slice of sort options
type Sort []SortCombinations
```

### 5b. SortCombinations (Union Type)
**File**: `sortcombinations.go`

```go
// Can be either:
// - string (simple field name like "field_name" or "_score", "_doc")
// - SortOptions (complex sort configuration)
type SortCombinations any
```

### 5c. SortOptions Struct
**File**: `sortoptions.go`

Comprehensive sort configuration supporting multiple sort criteria.

```go
type SortOptions struct {
    Doc_         *ScoreSort           `json:"_doc,omitempty"`
    GeoDistance_ *GeoDistanceSort     `json:"_geo_distance,omitempty"`
    Score_       *ScoreSort           `json:"_score,omitempty"`
    Script_      *ScriptSort          `json:"_script,omitempty"`
    SortOptions  map[string]FieldSort `json:"-"`  // Field-based sorts
}
```

### 5d. FieldSort Struct
**File**: `fieldsort.go`

Detailed sort configuration for a specific field.

```go
type FieldSort struct {
    Format       *string                                    `json:"format,omitempty"`
    Missing      Missing                                    `json:"missing,omitempty"`
    Mode         *sortmode.SortMode                         `json:"mode,omitempty"`
    Nested       *NestedSortValue                           `json:"nested,omitempty"`
    NumericType  *fieldsortnumerictype.FieldSortNumericType `json:"numeric_type,omitempty"`
    Order        *sortorder.SortOrder                       `json:"order,omitempty"`
    UnmappedType *fieldtype.FieldType                       `json:"unmapped_type,omitempty"`
}
```

### Fields:
- **Order** (*sortorder.SortOrder): `Asc` or `Desc`
- **Mode** (*sortmode.SortMode): How to handle multi-valued fields
- **Format** (*string): Date format for date fields
- **Missing** (Missing): Value to use for missing fields
- **Nested** (*NestedSortValue): For sorting nested documents
- **NumericType** (*fieldsortnumerictype.FieldSortNumericType): Numeric type hint
- **UnmappedType** (*fieldtype.FieldType): Type for unmapped fields

---

## 6. Range Query Types

### 6a. RangeQuery (Union Type)
**File**: `rangequery.go`

```go
// Union holding one of:
// - UntypedRangeQuery
// - DateRangeQuery
// - NumberRangeQuery
// - TermRangeQuery
type RangeQuery any
```

### 6b. DateRangeQuery
**File**: `daterangequery.go`

For date field range queries.

```go
type DateRangeQuery struct {
    Boost      *float32                       `json:"boost,omitempty"`
    Format     *string                        `json:"format,omitempty"`
    From       *string                        `json:"from,omitempty"`
    Gt         *string                        `json:"gt,omitempty"`  // Greater than
    Gte        *string                        `json:"gte,omitempty"` // Greater than or equal
    Lt         *string                        `json:"lt,omitempty"`  // Less than
    Lte        *string                        `json:"lte,omitempty"` // Less than or equal
    QueryName_ *string                        `json:"_name,omitempty"`
    Relation   *rangerelation.RangeRelation   `json:"relation,omitempty"`
    TimeZone   *string                        `json:"time_zone,omitempty"`
    To         *string                        `json:"to,omitempty"`
}
```

### 6c. NumberRangeQuery
**File**: `numberrangequery.go`

For numeric field range queries.

```go
type NumberRangeQuery struct {
    Boost      *float32                       `json:"boost,omitempty"`
    From       *Float64                       `json:"from,omitempty"`
    Gt         *Float64                       `json:"gt,omitempty"`
    Gte        *Float64                       `json:"gte,omitempty"`
    Lt         *Float64                       `json:"lt,omitempty"`
    Lte        *Float64                       `json:"lte,omitempty"`
    QueryName_ *string                        `json:"_name,omitempty"`
    Relation   *rangerelation.RangeRelation   `json:"relation,omitempty"`
    To         *Float64                       `json:"to,omitempty"`
}
```

### Fields (both types):
- **Gt/Gte/Lt/Lte**: Boundary operators
- **From/To**: Legacy boundary fields
- **Boost**: Query relevance multiplier
- **Relation**: How to match range field values (`intersects`, `contains`, `within`)

---

## 7. MatchPhraseQuery Type

**File**: `matchphrasequery.go`

Finds documents containing an exact phrase (ordered terms).

```go
type MatchPhraseQuery struct {
    Analyzer       *string                          `json:"analyzer,omitempty"`
    Boost          *float32                         `json:"boost,omitempty"`
    Query          string                           `json:"query"`  // Required
    QueryName_     *string                          `json:"_name,omitempty"`
    Slop           *int                             `json:"slop,omitempty"`
    ZeroTermsQuery *zerotermsquery.ZeroTermsQuery   `json:"zero_terms_query,omitempty"`
}
```

### Fields:
- **Query** (string): The phrase text - **REQUIRED**
- **Analyzer** (*string): Custom analyzer for tokenization
- **Slop** (*int): Max intervening unmatched positions
- **ZeroTermsQuery** (*ZeroTermsQuery): Behavior when analyzer removes all tokens

---

## 8. NestedQuery Type

**File**: `nestedquery.go`

Searches nested documents and returns the root parent.

```go
type NestedQuery struct {
    Boost          *float32                             `json:"boost,omitempty"`
    IgnoreUnmapped *bool                               `json:"ignore_unmapped,omitempty"`
    InnerHits      *InnerHits                          `json:"inner_hits,omitempty"`
    Path           string                              `json:"path"`  // Required
    Query          Query                               `json:"query"` // Required
    QueryName_     *string                             `json:"_name,omitempty"`
    ScoreMode      *childscoremode.ChildScoreMode      `json:"score_mode,omitempty"`
}
```

### Fields:
- **Path** (string): Path to nested object - **REQUIRED**
- **Query** (Query): Query to run on nested objects - **REQUIRED**
- **ScoreMode** (*ChildScoreMode): How to score parent documents
- **InnerHits** (*InnerHits): Include matching nested docs in results
- **IgnoreUnmapped** (*bool): Ignore unmapped paths

### Supporting Type: NestedSortValue
**File**: `nestedsortvalue.go`

```go
type NestedSortValue struct {
    Filter      *Query            `json:"filter,omitempty"`
    MaxChildren *int              `json:"max_children,omitempty"`
    Nested      *NestedSortValue  `json:"nested,omitempty"`  // Recursive nesting
    Path        string            `json:"path"`
}
```

---

## 9. ExistsQuery Type

**File**: `existsquery.go`

Matches documents where a field has any indexed value.

```go
type ExistsQuery struct {
    Boost      *float32  `json:"boost,omitempty"`
    Field      string    `json:"field"`  // Required
    QueryName_ *string   `json:"_name,omitempty"`
}
```

### Fields:
- **Field** (string): Field name to check - **REQUIRED**
- **Boost** (*float32): Relevance multiplier

---

## 10. BoolQuery Type

**File**: `boolquery.go`

Combines multiple queries using boolean logic.

```go
type BoolQuery struct {
    Boost              *float32          `json:"boost,omitempty"`
    Filter             []Query           `json:"filter,omitempty"`
    MinimumShouldMatch MinimumShouldMatch `json:"minimum_should_match,omitempty"`
    Must               []Query           `json:"must,omitempty"`
    MustNot            []Query           `json:"must_not,omitempty"`
    QueryName_         *string           `json:"_name,omitempty"`
    Should             []Query           `json:"should,omitempty"`
}
```

### Fields:
- **Must** ([]Query): All clauses must match (AND logic)
- **Should** ([]Query): At least one clause should match (OR logic)
- **MustNot** ([]Query): No clauses must match (NOT logic)
- **Filter** ([]Query): Must match, but doesn't affect scoring
- **MinimumShouldMatch** (MinimumShouldMatch): Min clauses from `should` to match
- **Boost** (*float32): Score multiplier

### Supporting Type: MinimumShouldMatch
**File**: `minimumshouldmatch.go`

```go
// Union type that can be:
// - int (exact count)
// - string (percentage like "75%", count with + or -)
type MinimumShouldMatch any
```

---

## 11. SortMode Enum

**File**: `/enums/sortmode/sortmode.go`

Controls how multi-valued fields are sorted.

```go
type SortMode struct {
    Name string
}

var (
    Min    = SortMode{"min"}     // Use minimum value
    Max    = SortMode{"max"}     // Use maximum value
    Sum    = SortMode{"sum"}     // Sum all values
    Avg    = SortMode{"avg"}     // Average of values
    Median = SortMode{"median"}  // Median value
)

// Usage example:
// sortMode := sortmode.Min
// sortMode.String() // Returns "min"
```

---

## 12. SortOrder Enum

**File**: `/enums/sortorder/sortorder.go`

Direction for sorting results.

```go
type SortOrder struct {
    Name string
}

var (
    Asc  = SortOrder{"asc"}
    Desc = SortOrder{"desc"}
)

// Usage example:
// order := sortorder.Desc
// order.String() // Returns "desc"
```

---

## Architecture Patterns

### 1. Union Types
Go doesn't have native union types, so the library uses `interface{}` (typed as `any`):
```go
type RangeQuery any  // Can hold DateRangeQuery, NumberRangeQuery, etc.
type FieldValue any  // Can hold string, int64, float, bool, etc.
```

### 2. Field-Per-Type Pattern
For queries that can apply to any field:
```go
type Query struct {
    Term  map[string]TermQuery `json:"term,omitempty"`
    Range map[string]RangeQuery `json:"range,omitempty"`
}
```

### 3. Custom JSON Marshaling
Most types implement custom `UnmarshalJSON()` to handle:
- Type flexibility (parsing `"asc"` into `SortOrder`)
- Multi-format support (can be string or object)
- Recursive structures

---

## Design Considerations for Query Builder

### Key Insights:
1. **Query is a mega-struct**: Contains fields for all 50+ query types - use conditional initialization
2. **Maps for field-level queries**: Use `map[string]SpecificQuery` for term, range, etc.
3. **Pointers for optional fields**: Most fields are pointers, make them nil for default behavior
4. **Custom marshaling**: Types handle JSON flexibly (e.g., `"asc"` or `{order: "asc"}`)
5. **Enums are not strict**: SortMode/SortOrder have fallback to custom values

### Builder Pattern Recommendations:
```go
// Initialize with NewQuery() to get empty maps pre-allocated
q := types.NewQuery()

// Build conditionally
if needTerm {
    q.Term = make(map[string]types.TermQuery)
    q.Term["status"] = types.TermQuery{Value: "active"}
}

// For BoolQuery composition
bq := types.NewBoolQuery()
bq.Must = append(bq.Must, *q)
```

---

## References
- Module: `github.com/elastic/go-elasticsearch/v8 v8.19.3`
- Specification: https://github.com/elastic/elasticsearch-specification
- Generated from Elasticsearch Specification v8

