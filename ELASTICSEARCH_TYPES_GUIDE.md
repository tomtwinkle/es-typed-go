# Elasticsearch Go v8 Typed API - Complete Type Guide

This guide provides a comprehensive reference for the Elasticsearch `go-elasticsearch/v8` library's type definitions, essential for building type-safe query builders and search applications.

## 📚 Documentation Files

This repository contains three complementary documentation files:

### 1. **QUICK_TYPE_REFERENCE.md** ⚡
**Start here for quick lookups**
- Quick reference table of all 16 key types
- Common operations with code examples
- Field explanations for each type
- Common mistakes and solutions
- Best practices and performance tips
- **Best for**: Developers building queries quickly

### 2. **TYPE_RELATIONSHIPS.md** 🗺️
**Visual understanding of type structure**
- ASCII diagram showing type hierarchy
- Query and sort structure breakdown
- Initialization patterns with examples
- Common query patterns (term, range, bool, nested, sort)
- Type system strengths and limitations
- **Best for**: Understanding architecture and designing query builders

### 3. **ES_TYPES_REFERENCE.md** 📖
**Deep dive into each type**
- Detailed struct definitions with all fields
- Field-by-field explanations
- Design patterns and architecture
- Supporting type definitions
- References and specification links
- **Best for**: Complete understanding and documentation

## 🎯 Quick Start: The 11 Core Types

```
Query          → Main query container (50+ query types)
├── TermQuery           → Exact term match
├── TermsQuery          → Multiple term match
├── RangeQuery          → Range matching (union)
│   ├── DateRangeQuery
│   └── NumberRangeQuery
├── MatchPhraseQuery    → Phrase matching
├── BoolQuery           → Boolean logic
├── NestedQuery         → Nested documents
├── ExistsQuery         → Field existence
└── [40+ more...]

Sort           → Sorting results
├── SortCombinations    → Union (string | SortOptions)
├── SortOptions         → Complex sort
└── FieldSort           → Per-field sort config
```

## 🔍 Type Locations

All types are in: `/home/runner/go/pkg/mod/github.com/elastic/go-elasticsearch/v8@v8.19.3/typedapi/types/`

**Key files:**
```
query.go                    671 lines - Query union type
boolquery.go                182 lines - BoolQuery (Must/Should/MustNot)
nestedquery.go              147 lines - NestedQuery
termquery.go                130 lines - TermQuery
termsquery.go               150 lines - TermsQuery
rangequery.go               31 lines  - RangeQuery union
daterangequery.go           163 lines - DateRangeQuery
numberrangequery.go         192 lines - NumberRangeQuery
matchphrasequery.go         167 lines - MatchPhraseQuery
existsquery.go              106 lines - ExistsQuery
sort.go                     27 lines  - Sort alias
sortoptions.go              138 lines - SortOptions
fieldsort.go                123 lines - FieldSort
enums/sortorder/            58 lines  - Asc | Desc
enums/sortmode/             70 lines  - Min | Max | Sum | Avg | Median
```

## 💡 Key Design Insights

### 1. Mega-Struct Pattern
`Query` contains 50+ optional fields for different query types:
```go
type Query struct {
    Term   map[string]TermQuery     `json:"term,omitempty"`
    Bool   *BoolQuery               `json:"bool,omitempty"`
    Nested *NestedQuery             `json:"nested,omitempty"`
    Range  map[string]RangeQuery    `json:"range,omitempty"`
    // ... 50+ more fields
}
```

### 2. Maps for Field-Level Queries
Queries that can target any field use maps:
```go
q.Term = map[string]TermQuery{
    "status": {Value: "active"},
    "verified": {Value: true},
}
// Can query multiple fields in one map
```

### 3. Union Types via `any`
Go lacks native unions, so the library uses `interface{}` (aliased `any`):
```go
type RangeQuery any  // Can hold DateRangeQuery OR NumberRangeQuery
type FieldValue any  // Can hold string OR int64 OR bool
```

### 4. Slices for Boolean Composition
Compound queries use slices to combine multiple conditions:
```go
type BoolQuery struct {
    Must    []Query  // ALL must match
    Should  []Query  // ANY may match
    MustNot []Query  // NONE must match
    Filter  []Query  // ALL must match (no scoring)
}
```

### 5. Custom JSON Marshaling
Types implement custom `UnmarshalJSON()` to handle flexibility:
- Parse `"asc"` string into `SortOrder` enum
- Accept both string and object forms
- Support recursive structures

## 📋 Type Inventory

### Query Types (50+)
- **Simple**: TermQuery, TermsQuery, ExistsQuery, PrefixQuery, WildcardQuery
- **Range**: DateRangeQuery, NumberRangeQuery, TermRangeQuery, UntypedRangeQuery
- **Text**: MatchQuery, MatchPhraseQuery, MatchPhrasePrefixQuery, MultiMatchQuery
- **Pattern**: FuzzyQuery, RegexpQuery, WildcardQuery
- **Compound**: BoolQuery, ConstantScoreQuery, DisMaxQuery, FunctionScoreQuery
- **Joining**: NestedQuery, HasChildQuery, HasParentQuery
- **Geo**: GeoDistanceQuery, GeoBoundingBoxQuery, GeoShapeQuery, GeoPolygonQuery
- **Span**: SpanTermQuery, SpanNearQuery, SpanFirstQuery, etc.
- **Special**: MatchAllQuery, MatchNoneQuery, ScriptQuery, IdQuery

### Support Types
- **FieldValue** (union): int64, Float64, string, bool, nil, json.RawMessage
- **TermsQueryField** (union): []FieldValue, TermsLookup
- **RangeQuery** (union): DateRangeQuery, NumberRangeQuery, TermRangeQuery, UntypedRangeQuery
- **SortCombinations** (union): string, SortOptions
- **MinimumShouldMatch** (union): int, string (percentage)

### Enums
- **SortOrder**: Asc, Desc
- **SortMode**: Min, Max, Sum, Avg, Median
- **ChildScoreMode**: None, Sum, Min, Max, Avg
- **ZeroTermsQuery**: All, None
- **RangeRelation**: Intersects, Contains, Within

## ✅ Best Practices

### 1. Always Use Constructor Functions
```go
✅ q := types.NewQuery()          // Maps pre-allocated
✅ bq := types.NewBoolQuery()     // Slices ready
❌ q := types.Query{}             // Nil maps, will panic on assignment
```

### 2. Use Pointers for Optional Fields
```go
✅ Order: &sortorder.Desc         // Pointer to enum
✅ Boost: &f                      // Pointer to float
❌ Order: sortorder.Desc          // Not a pointer!
```

### 3. Build Incrementally
```go
q := types.NewQuery()
if needTerm {
    q.Term = make(map[string]types.TermQuery)
    q.Term["status"] = types.TermQuery{Value: "active"}
}
```

### 4. Validate Required Fields
```go
// TermQuery requires Value
tq := types.TermQuery{
    Value: fieldValue,  // Must provide this
}

// NestedQuery requires Path and Query
nq := types.NestedQuery{
    Path: "comments",   // Required
    Query: query,       // Required
}
```

### 5. Leverage JSON Marshaling
```go
// Custom marshaling handles:
// - "asc" → SortOrder{Name: "asc"}
// - Omitting nil pointers
// - Handling map keys as field names

data, _ := json.Marshal(query)  // Works seamlessly
var q types.Query
json.Unmarshal(data, &q)        // Reconstructs correctly
```

## 🔧 Common Patterns

### Build a Term Query
```go
q := types.NewQuery()
q.Term = map[string]types.TermQuery{
    "status": {Value: "active"},
}
```

### Build a Complex Bool Query
```go
bq := types.NewBoolQuery()
bq.Must = []types.Query{
    {Term: map[string]types.TermQuery{"status": {Value: "active"}}},
}
bq.Filter = []types.Query{
    {Range: map[string]types.RangeQuery{
        "age": types.NumberRangeQuery{Gte: &Float64(18)},
    }},
}
q := types.NewQuery()
q.Bool = bq
```

### Build a Nested Query
```go
nq := types.NewNestedQuery()
nq.Path = "comments"
nq.Query = types.Query{
    Match: map[string]types.MatchQuery{
        "comments.text": {Query: "elasticsearch"},
    },
}
q := types.NewQuery()
q.Nested = nq
```

### Build Sorting
```go
sorts := types.Sort{
    "_score",  // Simple string sort
    types.SortOptions{
        SortOptions: map[string]types.FieldSort{
            "timestamp": {
                Order: &sortorder.Desc,
                Mode:  &sortmode.Max,
            },
        },
    },
}
```

## 🎓 Learning Path

1. **Start**: Read `QUICK_TYPE_REFERENCE.md` for overview
2. **Understand**: Check `TYPE_RELATIONSHIPS.md` for architecture
3. **Deep Dive**: Reference `ES_TYPES_REFERENCE.md` for details
4. **Practice**: Build queries incrementally using these patterns
5. **Optimize**: Refer to performance tips section

## 📚 External Resources

- **Specification**: https://github.com/elastic/elasticsearch-specification
- **go-elasticsearch**: https://github.com/elastic/go-elasticsearch
- **Elasticsearch DSL**: https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html

## 🤝 Query Builder Tips

### For Your Type-Safe Query Builder:

1. **Wrapper Types**: Create builder interfaces that wrap Query types
2. **Validation**: Add compile-time checks for required fields
3. **Chaining**: Implement method chaining for fluent API
4. **Defaults**: Provide sensible defaults (e.g., boost=1.0)
5. **Type Safety**: Use generics/interfaces to prevent wrong field usage

### Example Structure:
```go
type QueryBuilder struct {
    query *types.Query
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{query: types.NewQuery()}
}

func (qb *QueryBuilder) Term(field string, value interface{}) *QueryBuilder {
    if qb.query.Term == nil {
        qb.query.Term = make(map[string]types.TermQuery)
    }
    qb.query.Term[field] = types.TermQuery{Value: value}
    return qb  // Chain!
}

func (qb *QueryBuilder) Build() *types.Query {
    return qb.query
}
```

## 📝 Module Information

- **Package**: `github.com/elastic/go-elasticsearch/v8`
- **Version**: v8.19.3
- **Generated**: From Elasticsearch specification
- **Go Version**: 1.16+
- **License**: SSPL (Elastic)

---

**Last Updated**: 2024-03-14
**Generated From**: go-elasticsearch v8.19.3 typedapi/types
