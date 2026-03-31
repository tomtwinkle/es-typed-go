# Migration Guide: v2 Architecture

This document describes the changes introduced in the v2 architecture and explains how to update existing code.

## Goals

- Switching between Elasticsearch v8 and v9 requires changing **only one import path** (`esv8` → `esv9`).
- Query building, field access, aggregations, and sort definitions are written once against the top-level `query` package and work with both versions.
- Generated model accessors (`esmodel`) now include `Fields`, `Alias`, and `Index` in one place.

## Summary of Changes

| Area | Before | After |
|---|---|---|
| Query package | `esv8/query` or `esv9/query` | `query` (top-level, version-agnostic) |
| Sort direction | `sortorder.Desc` (version-specific) | `query.SortDesc` / `query.SortAsc` |
| `esmodel` field access | `esmodel.Product.Status` | `esmodel.Product.Fields.Status` |
| `esmodel` alias access | `estype.Alias("product")` (manual) | `esmodel.Product.Alias` |
| `esmodel` index access | `estype.Index("product-000001")` (manual) | `esmodel.Product.Index` |
| `SearchRequest` interface (v8) | `ToRequest() *v8search.Request` | `ToV8Request() *v8search.Request` |
| `SearchRequest` interface (v9) | `ToRequest() *v9search.Request` | `ToV9Request() *v9search.Request` |

---

## Step-by-Step Migration

### 1. Replace the query package import

Replace version-specific query imports with the top-level `query` package:

```go
// Before
import "github.com/tomtwinkle/es-typed-go/esv8/query"

// After
import "github.com/tomtwinkle/es-typed-go/query"
```

All function names (`TermValue`, `MatchAll`, `BoolQuery`, `NewSearch`, `AvgAgg`, `StringTermsAgg`, etc.) are identical. No call-site changes are needed beyond the import path.

### 2. Replace sort direction imports

The `query` package exports `query.SortAsc` and `query.SortDesc`, eliminating the need to import a version-specific `sortorder` package:

```go
// Before
import "github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"

query.NewSort().Field(esmodel.Product.Fields.Date, sortorder.Desc)

// After (no sortorder import needed)
query.NewSort().Field(esmodel.Product.Fields.Date, query.SortDesc)
```

### 3. Update esmodel field access

Regenerate your model files with `go generate` after adding `Alias()` and `Index()` methods to your definition struct.

**Add methods to your definition struct:**

```go
// esdefinition/product.go
func (Product) Alias() estype.Alias { return "product" }
func (Product) Index() estype.Index { return "product-000001" }
```

**Update the go:generate directive** (use `-group` with `-struct`):

```
//go:generate go tool estyped -struct Product -package esmodel -out ../esmodel/product_gen.go -group Product
```

**Run code generation:**

```bash
go generate ./...
```

**Update field access:**

```go
// Before
query.TermValue(esmodel.Product.Status, "active")
estype.Alias("product")

// After
query.TermValue(esmodel.Product.Fields.Status, "active")
esmodel.Product.Alias
```

### 4. Switch from v8 to v9 (the main goal)

After completing steps 1–3, switching from Elasticsearch v8 to v9 requires only two import changes:

```go
// Before (v8)
import (
    es8 "github.com/elastic/go-elasticsearch/v8"
    "github.com/tomtwinkle/es-typed-go/esv8"
    "github.com/tomtwinkle/es-typed-go/query"
    "github.com/tomtwinkle/es-typed-go/examples/quickstart/esmodel"
)

client, _ := esv8.NewClientWithLogger(es8.Config{Addresses: []string{"http://localhost:19200"}}, slog.Default())
resp, err := esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)

// After (v9) — only the two lines above change
import (
    es9 "github.com/elastic/go-elasticsearch/v9"
    "github.com/tomtwinkle/es-typed-go/esv9"
    "github.com/tomtwinkle/es-typed-go/query"
    "github.com/tomtwinkle/es-typed-go/examples/quickstart/esmodel"
)

client, _ := esv9.NewClientWithLogger(es9.Config{Addresses: []string{"http://localhost:19201"}}, slog.Default())
resp, err := esv9.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

Query building, aggregation access, field names, sort definitions, and model accessors are unchanged.

---

## Breaking Changes

### `SearchRequest` interface method renamed

If you have **custom types** that implement the `SearchRequest` interface, update the method name:

```go
// esv8: rename ToRequest → ToV8Request
func (p MyParams) ToV8Request() *v8search.Request { ... }

// esv9: rename ToRequest → ToV9Request
func (p MyParams) ToV9Request() *v9search.Request { ... }
```

`esv8.SearchParams`, `esv9.SearchParams`, `esv8/query.SearchParams`, `esv9/query.SearchParams`, and `query.SearchParams` all implement both `ToV8Request()` and `ToV9Request()` — no change needed for standard usage.

### Generated esmodel format changed

The generated model variable structure changed from a flat struct to a nested `Fields` sub-struct. Regenerate with `go generate ./...` and update field accesses from `esmodel.Product.Status` to `esmodel.Product.Fields.Status`.

---

## Backward Compatibility

The following remain available with deprecation notices to ease gradual migration:

- `esv8/query` and `esv9/query` packages are kept as thin wrappers re-exporting from the top-level `query` package.
- `esv8.SearchParams.ToRequest()` and `esv9.SearchParams.ToRequest()` are kept as deprecated aliases.
- The old constant-mode and group-mode estyped output (without `Alias`/`Index`) remains available if the definition struct does not implement `AliasProvider` or `IndexProvider`.

---

## Future

The top-level `query` package currently uses `go-elasticsearch/v8` types as its internal representation and performs a JSON round-trip when producing v9 requests. A future version may introduce a fully version-independent intermediate representation to remove this dependency, but the public API surface will remain stable.
