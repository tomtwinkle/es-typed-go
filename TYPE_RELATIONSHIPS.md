# Elasticsearch Types - Relationship Diagram

## Type Hierarchy & Relationships

```
Query (Main Container)
├── Leaf Query Types (Match documents)
│   ├── TermQuery
│   │   └── Value: FieldValue → (int64 | Float64 | string | bool | nil | json.RawMessage)
│   │
│   ├── TermsQuery
│   │   └── TermsQuery[fieldName]: TermsQueryField → ([]FieldValue | TermsLookup)
│   │
│   ├── MatchPhraseQuery
│   │   ├── Query: string (required)
│   │   ├── Analyzer: *string
│   │   ├── Slop: *int
│   │   └── ZeroTermsQuery: *ZeroTermsQuery
│   │
│   ├── ExistsQuery
│   │   └── Field: string (required)
│   │
│   └── RangeQuery (Union Type)
│       ├── DateRangeQuery
│       │   ├── Gt/Gte/Lt/Lte: *string (date format)
│       │   ├── Format: *string
│       │   ├── TimeZone: *string
│       │   └── Relation: *RangeRelation
│       │
│       ├── NumberRangeQuery
│       │   ├── Gt/Gte/Lt/Lte: *Float64
│       │   └── Relation: *RangeRelation
│       │
│       └── UntypedRangeQuery + TermRangeQuery
│
├── Compound Query Types
│   ├── BoolQuery
│   │   ├── Must: []Query (AND)
│   │   ├── Should: []Query (OR)
│   │   ├── MustNot: []Query (NOT)
│   │   ├── Filter: []Query (AND, no scoring)
│   │   └── MinimumShouldMatch: (int | string)
│   │
│   └── NestedQuery
│       ├── Path: string (required)
│       ├── Query: Query (required)
│       ├── ScoreMode: *ChildScoreMode
│       └── InnerHits: *InnerHits
│           └── Nested sorting support via NestedSortValue
│
└── [50+ more query types...]
    ├── Fuzzy, Wildcard, Prefix, Regexp (Pattern Matching)
    ├── Match, MultiMatch (Full-text Search)
    ├── Geo* queries (Geographic Queries)
    ├── Span* queries (Proximity Queries)
    └── Semantic, TextExpansion, etc.


Sort Structure
├── Sort (Type Alias)
│   └── []SortCombinations
│
├── SortCombinations (Union Type)
│   ├── string (simple: "field_name", "_score", "_doc")
│   └── SortOptions
│
└── SortOptions (Complex sort)
    ├── Score_: *ScoreSort
    ├── Doc_: *ScoreSort
    ├── GeoDistance_: *GeoDistanceSort
    ├── Script_: *ScriptSort
    └── SortOptions[fieldName]: FieldSort
        ├── Order: *SortOrder
        │   └── Enum: Asc | Desc
        ├── Mode: *SortMode
        │   └── Enum: Min | Max | Sum | Avg | Median
        ├── Format: *string
        ├── Missing: Missing (what value for missing fields)
        ├── NumericType: *FieldSortNumericType
        ├── UnmappedType: *FieldType
        └── Nested: *NestedSortValue
            ├── Path: string
            ├── Filter: *Query
            ├── MaxChildren: *int
            └── Nested: *NestedSortValue (recursive)


Key Type Aliases
├── FieldValue any
│   └── Union of: int64 | Float64 | string | bool | nil | json.RawMessage
│
├── TermsQueryField any
│   └── Union of: []FieldValue | TermsLookup
│
├── RangeQuery any
│   └── Union of: UntypedRangeQuery | DateRangeQuery | NumberRangeQuery | TermRangeQuery
│
├── SortCombinations any
│   └── Union of: string | SortOptions
│
└── MinimumShouldMatch any
    └── Union of: int | string
```

## Initialization Patterns

### Creating a Query

```go
// 1. Initialize empty Query
q := types.NewQuery()
// Auto-initializes maps for all query types

// 2. Create leaf query
termQuery := types.TermQuery{
    Value: "active",  // FieldValue - can be string, int, bool, etc.
}

// 3. Add to Query
q.Term = make(map[string]types.TermQuery)
q.Term["status"] = termQuery

// 4. Create compound query
boolQ := types.NewBoolQuery()
boolQ.Must = append(boolQ.Must, *q)

// 5. Wrap in outer Query
outerQ := types.NewQuery()
outerQ.Bool = boolQ
```

### Creating Sorts

```go
// 1. Simple sort (just field name)
var sorts types.Sort
sorts = append(sorts, "field_name")  // string as SortCombinations

// 2. Complex sort with options
fieldSort := types.FieldSort{
    Order: &sortorder.Desc,
    Mode:  &sortmode.Min,
}

// 3. Build SortOptions
sortOpts := types.NewSortOptions()
sortOpts.SortOptions["timestamp"] = fieldSort

// 4. Add to Sort
sorts = append(sorts, sortOpts)  // SortOptions as SortCombinations
```

## Important Design Notes

### 1. Pointers vs Values
- **Pointers (`*T`)**: Optional fields (marshaled as omitempty)
- **Values (`T`)**: Required or special handling
  - Exception: `Query`, `FieldValue`, `SortCombinations` are `any` (flexibly typed)

### 2. Maps for Multi-Field Queries
```go
Query.Term      map[string]TermQuery        // Can match multiple fields
Query.Range     map[string]RangeQuery       // Multiple range filters
Query.Match     map[string]MatchQuery       // Multiple fields
SortOptions     map[string]FieldSort        // Sort by multiple fields
```

### 3. Slices for Compound Queries
```go
BoolQuery.Must     []Query   // ALL must match
BoolQuery.Should   []Query   // ANY should match
BoolQuery.MustNot  []Query   // NONE must match
BoolQuery.Filter   []Query   // ALL must match (no scoring)
```

### 4. Union Type Handling
Since Go has no native union types, the library uses `any`:
- At encode time: Type system enforces correct field
- At decode time: Custom JSON unmarshaling handles parsing
- At use time: Type assertion or interface{} pattern needed

### 5. Recursion
Some types are recursive:
```go
NestedSortValue {
    Nested *NestedSortValue  // Can nest arbitrarily deep
}

BoolQuery {
    Must []Query  // Query is recursive, can contain BoolQuery
}
```

## Common Query Patterns

### Pattern 1: Single Term Match
```go
q := types.NewQuery()
q.Term = map[string]types.TermQuery{
    "status": {Value: "active"},
}
```

### Pattern 2: Multiple Conditions (AND)
```go
bq := types.NewBoolQuery()
bq.Must = []types.Query{
    {Term: map[string]types.TermQuery{"status": {Value: "active"}}},
    {Term: map[string]types.TermQuery{"verified": {Value: true}}},
}
```

### Pattern 3: Multiple Conditions with Range
```go
bq := types.NewBoolQuery()
bq.Must = append(bq.Must, types.Query{
    Range: map[string]types.RangeQuery{
        "timestamp": types.DateRangeQuery{
            Gte: pointers.String("2024-01-01"),
        },
    },
})
```

### Pattern 4: Nested Document Search
```go
nq := types.NewNestedQuery()
nq.Path = "comments"
nq.Query = types.Query{
    Term: map[string]types.TermQuery{
        "comments.author": {Value: "john"},
    },
}
```

### Pattern 5: Complex Sorting
```go
sorts := types.Sort{
    types.SortOptions{
        SortOptions: map[string]types.FieldSort{
            "_score": {Order: &sortorder.Desc},
            "timestamp": {
                Order: &sortorder.Desc,
                Mode:  &sortmode.Max,
            },
        },
    },
}
```

---

## File Locations Reference

```
$GOMODCACHE/github.com/elastic/go-elasticsearch/v8@v8.19.3/typedapi/types/

Core Types:
├── query.go                      # Query union type
├── boolquery.go                  # BoolQuery
├── nestedquery.go                # NestedQuery
├── termquery.go                  # TermQuery
├── termsquery.go                 # TermsQuery
├── matchphrasequery.go           # MatchPhraseQuery
├── existsquery.go                # ExistsQuery
├── daterangequery.go             # DateRangeQuery
├── numberrangequery.go           # NumberRangeQuery
├── rangequery.go                 # RangeQuery (union)

Sort Types:
├── sort.go                       # Sort alias
├── sortcombinations.go           # SortCombinations union
├── sortoptions.go                # SortOptions
├── fieldsort.go                  # FieldSort
├── nestedsortvalue.go            # NestedSortValue

Union/Support Types:
├── fieldvalue.go                 # FieldValue union
├── termsqueryfield.go            # TermsQueryField union
├── minimumshouldmatch.go         # MinimumShouldMatch union

Enums (in enums/ subdirectory):
├── sortorder/sortorder.go        # Asc | Desc
├── sortmode/sortmode.go          # Min | Max | Sum | Avg | Median
└── [many other enum types]
```

---

## Type System Strengths

✅ **Type Safety**: Compile-time checks for field names, value types  
✅ **Flexibility**: Union types via `any` + custom marshaling  
✅ **Composability**: Nested structs enable complex query building  
✅ **Spec Compliance**: Generated from official Elasticsearch spec  
✅ **JSON Agnostic**: Seamless JSON encoding/decoding  

## Type System Limitations

⚠️ **No Strict Enums**: SortOrder/SortMode accept custom values  
⚠️ **No Compile-Time Field Validation**: Map keys are strings  
⚠️ **Type Assertions Needed**: When working with `any` types  
⚠️ **Verbose Construction**: Requires lots of boilerplate for complex queries  

