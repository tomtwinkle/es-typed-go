# Search Guide

This guide explains the search APIs provided by `es-typed-go`, when to use each one, and how they relate to the lower-level Elasticsearch request types.

The top-level `README*` files should stay brief. This document is the detailed reference for search usage.

## Overview

`es-typed-go` provides three main search entry points:

- `Search[T](...)` — high-level typed search that returns hit metadata, total hits, aggregations, and the raw response
- `SearchDocuments[T](...)` — high-level typed search that returns only decoded document sources
- `SearchRaw(...)` — low-level escape hatch for advanced Elasticsearch request shapes

In normal application code:

- prefer `Search[T](...)` when you need hit metadata or aggregations
- prefer `SearchDocuments[T](...)` when you only need decoded `_source` values
- prefer `Count(...)` when you only need the number of matching documents
- use `SearchRaw(...)` only when you need Elasticsearch request features that are not modeled by the high-level helpers

## Search helper summary

### `Search[T](...)`

Use `Search[T](...)` when you want:

- decoded documents
- total hit count
- hit metadata such as `_id`, `_index`, and `_score`
- typed aggregation access
- access to the raw Elasticsearch response

Example:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Offset(0).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
if err != nil {
    return err
}

fmt.Println(resp.Total)
for _, hit := range resp.Hits {
    fmt.Println(hit.ID, hit.Index, hit.Source)
}
```

### `SearchDocuments[T](...)`

Use `SearchDocuments[T](...)` when you only want decoded `_source` values and do not care about hit metadata or aggregations.

Example:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
if err != nil {
    return err
}

for _, doc := range docs {
    fmt.Println(doc)
}
```

### `SearchRaw(...)`

Use `SearchRaw(...)` when you need advanced request shapes that are not covered by the high-level helper surface.

Typical examples:

- `search_after`
- point-in-time
- custom `_source` filtering
- request forms that need direct typed client request access
- Elasticsearch features that the library has not wrapped yet

Example:

```go
req := search.NewRequest()
req.Query = &types.Query{
    MatchAll: &types.MatchAllQuery{},
}

rawResp, err := client.SearchRaw(ctx, alias, req)
if err != nil {
    return err
}

fmt.Println(rawResp.Hits.Total)
```

## Which API should I use?

### I want typed hits plus metadata

Use `Search[T](...)`.

This is the best default for application searches.

### I only want decoded documents

Use `SearchDocuments[T](...)`.

This is a convenience helper over `Search[T](...)`.

### I only want one document

Use `Search[T](...)` or `SearchDocuments[T](...)` with `Limit(1)` or `Size: 1`.

Example:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(1).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
if err != nil {
    return err
}

if len(resp.Hits) == 0 {
    fmt.Println("not found")
    return nil
}

fmt.Println(resp.Hits[0].Source)
```

If you only want decoded documents:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(1).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
if err != nil {
    return err
}

if len(docs) == 0 {
    fmt.Println("not found")
    return nil
}

fmt.Println(docs[0])
```

## Do not use search just to get counts

If you only need the number of matching documents, prefer `Count(...)`.

Example:

```go
res, err := client.Count(ctx, alias)
if err != nil {
    return err
}

fmt.Println(res.Count)
```

Why this is preferred:

- the intent is clearer
- the API surface is simpler
- you are not asking Elasticsearch to return hits you do not need

### What about `Limit(0)`?

`Limit(0)` means “return no hits”.

That can still be useful when you want:

- aggregations only
- total hits from a search request
- a search-shaped request without document payloads

But if your goal is strictly “how many documents match?”, `Count(...)` is usually the better API.

## Directly passing query builder params

`query.NewSearch().Build()` returns `query.SearchParams`.

You can pass that value directly into the high-level helpers.

You do not need to manually copy fields into package-level `SearchParams`.

Example:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
if err != nil {
    return err
}
```

This works in both:

- `esv8`
- `esv9`

## Package-level `SearchParams` vs builder `query.SearchParams`

There are two common search parameter types:

- `esv8.SearchParams` / `esv9.SearchParams`
- `esv8/query.SearchParams` / `esv9/query.SearchParams`

Both can be used with the high-level helpers as long as they can convert themselves into a typed Elasticsearch search request.

### When to use builder params

Prefer builder params when:

- you are constructing queries fluently
- you want type-safe query helpers
- you want to compose query, sort, aggregations, and pagination in one place

Example:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Sort(
        query.NewSort().
            ScoreDesc().
            Build()...,
    ).
    Limit(20).
    Offset(0).
    Build()
```

### When to use package-level `SearchParams`

Use package-level `SearchParams` when:

- you already have the pieces separately
- you are not using the fluent builder
- you want to assemble the request fields directly

Example:

```go
params := esv8.SearchParams{
    Query: types.Query{
        MatchAll: &types.MatchAllQuery{},
    },
    Size: 10,
    From: 0,
}
```

## Search response structure

### `SearchHit[T]`

Each hit contains:

- `ID`
- `Index`
- `Score`
- `Source`
- `Raw`

This lets you work with a typed document while still having access to low-level hit metadata.

### `SearchResponse[T]`

A high-level response contains:

- `Total`
- `Hits`
- `Aggregations`
- `Raw`

This makes it suitable for application code that wants a typed response while still preserving escape hatches.

## Aggregations

If you need aggregation results alongside typed documents, use `Search[T](...)`.

Example:

```go
avgPriceAgg := query.AvgAgg("avg_price", ProductFields.Price)
byCategoryAgg := query.StringTermsAgg(
    "by_category",
    ProductFields.Category,
    query.WithSubAggs(avgPriceAgg),
)

params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Aggregation(query.Aggs(byCategoryAgg).Build()).
    Limit(10).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
if err != nil {
    return err
}

terms := resp.Aggregations.MustStringTerms(byCategoryAgg)
for _, bucket := range terms.Buckets() {
    avg := bucket.Aggregations().MustAvg(avgPriceAgg)
    fmt.Println(bucket.Key(), avg.Value())
}
```

If you do not need aggregation data, `SearchDocuments[T](...)` is often a better fit.

## Pagination

Use `Limit(...)` and `Offset(...)` when constructing search params.

Example:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(20).
    Offset(40).
    Build()
```

Important distinction:

- `resp.Total` is the total number of matching documents
- `len(resp.Hits)` is only the number of hits returned in the current page

For example:

- `Total = 125`
- `len(resp.Hits) = 20`

means:

- 125 documents matched
- 20 were returned in this page

## v8 and v9 parity

The search helper design is intended to stay aligned between:

- `esv8`
- `esv9`

In general, if you learn one, the other should feel the same.

Typical usage shape:

```go
// v8
v8Resp, err := esv8.Search[Product](ctx, v8Client, alias, params)

// v9
v9Resp, err := esv9.Search[Product](ctx, v9Client, alias, params)
```

Likewise for document-only usage:

```go
// v8
v8Docs, err := esv8.SearchDocuments[Product](ctx, v8Client, alias, params)

// v9
v9Docs, err := esv9.SearchDocuments[Product](ctx, v9Client, alias, params)
```

## Recommended patterns

### Normal application search

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
```

### Document-only retrieval

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
```

### Single-result retrieval

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(1).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
```

### Count-only retrieval

```go
res, err := client.Count(ctx, alias)
```

### Advanced request shape

```go
req := search.NewRequest()
// configure advanced Elasticsearch request details directly

rawResp, err := client.SearchRaw(ctx, alias, req)
```

## Anti-patterns

### Manually expanding builder params without a reason

Avoid this unless you specifically need to transform fields:

```go
resp, err := esv8.Search[Product](ctx, client, alias, esv8.SearchParams{
    Query:        params.Query,
    Sort:         params.Sort,
    Aggregations: params.Aggregations,
    Highlight:    params.Highlight,
    Collapse:     params.Collapse,
    ScriptFields: params.ScriptFields,
    Size:         params.Size,
    From:         params.From,
})
```

Prefer:

```go
resp, err := esv8.Search[Product](ctx, client, alias, params)
```

### Using search when count is the real goal

Avoid using search helpers only to inspect totals when `Count(...)` is sufficient.

## Related documents

- [../README.md](../README.md) — concise top-level overview
- [../examples/quickstart/README.md](../examples/quickstart/README.md) — runnable quickstart walkthrough
- [property-reference.md](property-reference.md) — property builder reference
- [contributing.md](contributing.md) — contributor workflow and repository rules