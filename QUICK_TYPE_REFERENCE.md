# Elasticsearch Go Types - Quick Reference Guide

## At a Glance

| Type | Purpose | File | Key Fields |
|------|---------|------|-----------|
| **Query** | Main query container (union of 50+ query types) | query.go | Bool, Range, Term, Nested, etc. |
| **TermQuery** | Exact term match | termquery.go | `Value: FieldValue` (required) |
| **TermsQuery** | Match one of multiple terms | termsquery.go | `TermsQuery: map[string]TermsQueryField` |
| **RangeQuery** | Range matching (union) | rangequery.go | Union of Date/Number/Term/Untyped |
| **DateRangeQuery** | Date range | daterangequery.go | `Gte/Lte: *string`, `TimeZone` |
| **NumberRangeQuery** | Numeric range | numberrangequery.go | `Gte/Lte: *Float64` |
| **MatchPhraseQuery** | Exact phrase match | matchphrasequery.go | `Query: string` (required), `Slop` |
| **BoolQuery** | Boolean logic | boolquery.go | `Must/Should/MustNot/Filter: []Query` |
| **NestedQuery** | Nested doc search | nestedquery.go | `Path/Query: string/Query` (required) |
| **ExistsQuery** | Field existence | existsquery.go | `Field: string` (required) |
| **FieldValue** | Value union | fieldvalue.go | int64, Float64, string, bool, nil, json.RawMessage |
| **Sort** | Sort array | sort.go | `[]SortCombinations` |
| **SortCombinations** | Sort union | sortcombinations.go | string OR SortOptions |
| **SortOptions** | Complex sort | sortoptions.go | `SortOptions: map[string]FieldSort` |
| **FieldSort** | Field sort config | fieldsort.go | `Order/Mode/Format/Nested` |
| **SortOrder** | Enum | enums/sortorder/sortorder.go | `Asc`, `Desc` |
| **SortMode** | Enum | enums/sortmode/sortmode.go | `Min`, `Max`, `Sum`, `Avg`, `Median` |

## Common Operations

### Build a Term Query
```go
q := types.NewQuery()
q.Term = map[string]types.TermQuery{
    "status": {Value: "active"},
}
// JSON: {"term": {"status": {"value": "active"}}}
```

### Build a Range Query
```go
q := types.NewQuery()
q.Range = map[string]types.RangeQuery{
    "age": types.NumberRangeQuery{
        Gte: &Float64(18),
        Lte: &Float64(65),
    },
}
// JSON: {"range": {"age": {"gte": 18, "lte": 65}}}
```

### Build a Bool Query
```go
bq := types.NewBoolQuery()
bq.Must = []types.Query{
    {Term: map[string]types.TermQuery{"status": {Value: "active"}}},
}
bq.Filter = []types.Query{
    {Range: map[string]types.RangeQuery{"age": types.NumberRangeQuery{Gte: &Float64(18)}}},
}

q := types.NewQuery()
q.Bool = bq
// JSON: {"bool": {"must": [...], "filter": [...]}}
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
// JSON: {"nested": {"path": "comments", "query": {...}}}
```

### Build a Sort
```go
sorts := types.Sort{
    types.SortCombinations("_score"),  // String literal sort
    types.SortCombinations(types.SortOptions{
        SortOptions: map[string]types.FieldSort{
            "timestamp": {
                Order: &sortorder.Desc,
                Mode:  &sortmode.Max,
            },
        },
    }),
}
```

## Field Explanation

### TermQuery
- **Value**: What to search for (can be int, string, bool, etc.)
- **Boost**: Score multiplier (1.0 default)
- **CaseInsensitive**: ASCII case-insensitive match

### DateRangeQuery
- **Gte/Gt/Lte/Lt**: Boundary conditions
- **Format**: Date format string (e.g., "yyyy-MM-dd")
- **TimeZone**: Convert dates to UTC (e.g., "America/New_York")
- **Relation**: How range fields match (`intersects`, `contains`, `within`)

### NumberRangeQuery
- **Gte/Gt/Lte/Lt**: Numeric boundaries
- **Relation**: Same as DateRangeQuery

### MatchPhraseQuery
- **Query**: The phrase text (required)
- **Analyzer**: Tokenization method
- **Slop**: Max unmatched positions between terms
- **ZeroTermsQuery**: Behavior if analyzer removes all tokens

### BoolQuery
- **Must**: ALL clauses must match (AND, affects score)
- **Should**: ANY clause may match (OR)
- **MustNot**: NO clauses must match (NOT)
- **Filter**: ALL clauses must match (AND, no scoring)
- **MinimumShouldMatch**: Minimum count/percentage of `should` to match

### NestedQuery
- **Path**: Path to nested object (required)
- **Query**: Query on nested docs (required)
- **ScoreMode**: How to score parent (none, sum, min, max, avg)
- **InnerHits**: Return matching nested docs

### FieldSort
- **Order**: `Asc` or `Desc`
- **Mode**: How to handle multi-valued fields (`min`, `max`, `sum`, `avg`, `median`)
- **Format**: Date format for date fields
- **Nested**: For sorting nested docs (with `NestedSortValue`)
- **Missing**: Value for missing fields

## Pointers and Nil

### Rule
- **`*T` fields**: Optional, set to `nil` or `&value`
- **Non-pointer fields**: Required (except special `any` types)

### Example
```go
query := types.TermQuery{
    Value:          "active",           // Required
    Boost:          nil,                // Optional - omitted if nil
    CaseInsensitive: &boolTrue,         // Optional - use pointer
    QueryName_:     nil,                // Optional - omitted if nil
}
```

## NewQuery() Initialization

```go
q := types.NewQuery()
// Pre-allocated maps:
// - Term: map[string]TermQuery{}
// - Range: map[string]RangeQuery{}
// - Match: map[string]MatchQuery{}
// - Fuzzy: map[string]FuzzyQuery{}
// ... and many more
```

## NewBoolQuery() Initialization

```go
bq := types.NewBoolQuery()
// Fields are nil/empty slices:
// - Must: nil → append to it
// - Should: nil → append to it
// - MustNot: nil → append to it
// - Filter: nil → append to it
```

## NewSortOptions() Initialization

```go
so := types.NewSortOptions()
// Pre-allocated:
// - SortOptions: map[string]FieldSort{}
```

## Enum Usage

### SortOrder
```go
import "github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"

order := sortorder.Asc
order.String()  // Returns "asc"
order.MarshalText()  // Returns []byte("asc")
```

### SortMode
```go
import "github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortmode"

mode := sortmode.Min
mode.String()  // Returns "min"
```

## Union Type Handling

### FieldValue (any)
```go
var value types.FieldValue

// Can assign:
value = "active"           // string
value = int64(42)          // int64
value = 3.14              // float64 → Float64
value = true              // bool
value = nil               // nil
```

### RangeQuery (any)
```go
var rq types.RangeQuery

// Can assign:
rq = types.DateRangeQuery{Gte: pointers.String("2024-01-01")}
rq = types.NumberRangeQuery{Gte: &Float64(18)}
rq = types.TermRangeQuery{Gte: "aaa"}
rq = types.UntypedRangeQuery{Gte: pointers.String("2024-01-01")}
```

### SortCombinations (any)
```go
var sc types.SortCombinations

// Can assign:
sc = "field_name"  // string
sc = types.SortOptions{...}  // SortOptions struct
```

## Common Mistakes

❌ **Forget to initialize maps**
```go
q := types.Query{}  // ← Term, Range maps are nil!
q.Term["status"] = ...  // ← Panic: assignment to entry in nil map
```
✅ **Use NewQuery()**
```go
q := types.NewQuery()  // ← Maps pre-allocated
q.Term["status"] = ...  // ✓ Works
```

---

❌ **Forget pointers for optional fields**
```go
sort := types.FieldSort{
    Order: sortorder.Desc,  // ← Not a pointer!
}
```
✅ **Use pointer to enum**
```go
sort := types.FieldSort{
    Order: &sortorder.Desc,  // ✓ Correct
}
```

---

❌ **Mix union types**
```go
rq := types.DateRangeQuery{}
q.Range["age"] = rq  // ← OK: DateRangeQuery is valid RangeQuery
```
✅ **Both work because RangeQuery is `any`**

---

## Performance Tips

1. **Reuse maps**: Don't recreate `map[string]TermQuery` for each query
2. **Pre-allocate**: Use `make(map[string]TermQuery, capacity)` if you know size
3. **Avoid JSON round-trips**: Build types directly, not JSON then decode
4. **Lazy initialization**: Only set fields you need (use nil for rest)

## JSON Serialization Notes

- **Omitempty rule**: Nil pointers and zero values are omitted
- **Custom marshaling**: Some types have custom JSON handlers
- **Map fields**: `json:"-"` prevents double-marshaling
- **Underscore fields**: `QueryName_` becomes `_name` in JSON

