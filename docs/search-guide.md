# Search Guide

## Which API should I use?

| Use case | API |
|---|---|
| Typed hits, metadata, aggregations | `Search[T](ctx, client, alias, params)` |
| Only decoded document `_source` | `SearchDocuments[T](ctx, client, alias, params)` |
| Count matching documents only | `Count(ctx, alias)` |
| Advanced request shapes not covered by helpers | `SearchRaw(ctx, alias, req)` |

## Building search parameters

Use the top-level `query` package for all query building. It is version-agnostic — the same code works with both `esv8` and `esv9`.

```go
import "github.com/tomtwinkle/es-typed-go/query"

params := query.NewSearch().
    Where(query.TermValue(esmodel.Product.Fields.Status, "active")).
    Where(
        query.TermValue(esmodel.Product.Fields.Category, "electronics"),
        query.DateRangeQuery(esmodel.Product.Fields.Date, query.DateRangeGte("2024-01-01"), query.DateRangeLte("2024-12-31")),
    ).
    Sort(
        query.NewSort().
            Field(esmodel.Product.Fields.Date, query.SortDesc).
            ScoreDesc().
            Build()...,
    ).
    Aggregation(query.Aggs(
        query.AvgAgg("avg_price", esmodel.Product.Fields.Price),
    ).Build()).
    Limit(20).
    Offset(0).
    Build()
```

Pass `params` directly to `Search[T]`:

```go
// v8
resp, err := esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)

// v9 — identical, only the package changes
resp, err := esv9.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

## Switching between v8 and v9

Because all query building uses the shared `query/` package, switching Elasticsearch versions only requires changing the client import and instantiation. See [migration-v2.md](migration-v2.md) for the full list of changes.

## Aggregation results

`SearchResponse.Aggregations` is of type `query.AggResults`. Use `GetXxx` / `MustXxx` methods to retrieve typed results:

```go
avgDef := query.AvgAgg("avg_price", esmodel.Product.Fields.Price)
termsDef := query.StringTermsAgg("by_category", esmodel.Product.Fields.Category,
    query.WithSubAggs(avgDef))

// ... run search ...

terms := resp.Aggregations.MustStringTerms(termsDef)
for _, bucket := range terms.Buckets() {
    avg, _ := bucket.Aggregations().GetAvg(avgDef)
    // avg.Value() is *float64
}
```

### Nested aggregation

Count documents inside a nested object and run sub-aggregations within that scope:

```go
avgDef := query.AvgAgg("avg_price", esmodel.Order.Fields.Items.Price)
nestedDef := query.NestedAgg("items", esmodel.Order.Fields.Items,
    query.NestedAggSubAggs(avgDef))

// ... run search ...

nested := resp.Aggregations.MustNested(nestedDef)
// nested.DocCount()  — total nested document count
avg, _ := nested.Aggregations().GetAvg(avgDef)
```

### Filter aggregation

Restrict an aggregation to documents matching a query:

```go
avgDef := query.AvgAgg("avg_price", esmodel.Product.Fields.Price)
filterDef := query.FilterAgg("active_products",
    query.TermValue(esmodel.Product.Fields.Status, "active"),
    query.FilterAggSubAggs(avgDef),
)

// ... run search ...

filtered := resp.Aggregations.MustFilter(filterDef)
// filtered.DocCount()
avg, _ := filtered.Aggregations().GetAvg(avgDef)
```

### Multi-terms aggregation

Group documents by multiple fields simultaneously:

```go
multiDef := query.MultiTermsAgg("by_category_status",
    []estype.Field{esmodel.Product.Fields.Category, esmodel.Product.Fields.Status},
    query.MultiTermsAggSize(10),
)

// ... run search ...

multi := resp.Aggregations.MustMultiTerms(multiDef)
for _, bucket := range multi.Buckets() {
    // bucket.Keys()     — []string composite key values
    // bucket.DocCount() — document count for this key combination
    fmt.Println(bucket.Keys(), bucket.DocCount())
}
```

## Sort directions

Use `query.SortAsc` and `query.SortDesc` — no need to import a version-specific `sortorder` package:

```go
query.NewSort().Field(esmodel.Product.Fields.Date, query.SortDesc)
```

## Query helpers reference

### DateRangeQuery

`DateRangeQuery` accepts functional options so you can use any combination of the four comparison operators:

```go
// Gte + Lte (closed range)
query.DateRangeQuery(field, query.DateRangeGte("2024-01-01"), query.DateRangeLte("2024-12-31"))

// Gt + Lt (open range)
query.DateRangeQuery(field, query.DateRangeGt("2024-01-01"), query.DateRangeLt("2025-01-01"))

// One-sided bound
query.DateRangeQuery(field, query.DateRangeGte("2024-01-01"))
```

Available options: `DateRangeGt`, `DateRangeGte`, `DateRangeLt`, `DateRangeLte`.

### MultiTermsAgg with per-field Missing

Use `query.MultiTermLookup` to configure each field individually. Set `Missing` to substitute a value for documents that do not have that field:

```go
query.MultiTermsAgg("by_date_tz", []query.MultiTermLookup{
    {Field: esmodel.Item.Fields.BusinessDate},
    {Field: esmodel.Item.Fields.Timezone, Missing: "UTC"},
})
```

### Field.Ptr() and Field.String() — typed field conversion

When a raw go-elasticsearch type requires `*string` (e.g. `NestedAggregation.Path`, `SumAggregation.Field`), use `Ptr()` instead of a temporary variable:

```go
// Before
path := string(esmodel.Item.Fields.Items)
types.NestedAggregation{Path: &path}

// After
types.NestedAggregation{Path: esmodel.Item.Fields.Items.Ptr()}
```

When a field is `string` (not `*string`), e.g. `types.NestedSortValue.Path`, use `String()` instead:

```go
// NestedSortValue.Path is string, not *string — Ptr() cannot be used here
types.NestedSortValue{Path: esmodel.Item.Fields.Items.String()}
```

`Ptr()` and `String()` are also available on `estype.Alias` and `estype.Index`.

## Notes

- `Limit(0)` returns no hits (useful for aggregation-only or count-oriented searches).
- For a single result, use `Limit(1)`.
- `SearchRaw` accepts any `*search.Request` and is the escape hatch for request shapes not covered by the high-level helpers.
- The same API selection guidance applies to both `esv8` and `esv9`.

## Related documents

- [../README.md](../README.md) — concise top-level overview
- [migration-v2.md](migration-v2.md) — v2 architecture and migration steps
- [property-reference.md](property-reference.md) — property builder reference
- [contributing.md](contributing.md) — contributor setup and validation steps
