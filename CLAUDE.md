# es-typed-go — CLAUDE.md

A type-safe Go wrapper for [go-elasticsearch](https://github.com/elastic/go-elasticsearch) (v8 and v9) that prevents field-name typos and index/alias confusion at compile time.

---

## Build and Test Commands

```bash
# Build all packages
go build ./...

# Run all unit tests
go test ./... -short -count=1

# Run unit tests with verbose output
go test -v -count=1 -timeout=60s ./...

# Run integration tests against a running Elasticsearch instance
go test -tags=integration -v ./estype/... ./esv8/...   # v8
go test -tags=integration -v ./esv9/...                # v9

# Regenerate generated files (property wrappers, API coverage tests)
go generate ./...

# Generate field constants from a mapping file
go tool estyped -mapping mapping.json -out model.go -package model
# Struct mode (grouped access via model.Sample.FieldName)
go tool estyped -mapping mapping.json -out model.go -package model -name Sample
# Generate field constants from a Go struct with JSON tags
go tool estyped -struct MyType -out model.go -package model
```

---

## Repository Structure

```
estype/           Core shared types
  field.go        estype.Field — distinct string type for ES field names
  index.go        estype.Index — distinct string type for index names
  alias.go        estype.Alias — distinct string type for alias names
  date_format.go  All built-in ES date format constants + JoinDateFormats()
  mapping.go      ParseMapping() — parses ES mapping JSON into []MappingField
  property_mapping.go  MappingProperty interface and typed mapping property definitions

esv8/             Elasticsearch v8 wrapper
  esclient.go     ESClient interface (curated operations, uses estype.Index/Alias)
  client.go       Concrete implementation
  typed_search.go High-level Search[T], SearchDocuments, SearchOne, SearchParams
  property.go     Functional-option constructors for all 52 ES property types
  esclient_spec.go  ESClientSpec interface (full ES API spec — generated)
  client_spec.go    Full spec implementation (generated)
  query/          Query/sort/aggregation builders for v8

esv9/             Elasticsearch v9 wrapper (mirrors esv8/ exactly)
  typed_search.go High-level Search[T], SearchDocuments, SearchOne, SearchParams
  query/          Query/sort/aggregation builders for v9

cmd/estyped/      CLI code generator
  main.go         Reads ES mapping JSON → outputs typed estype.Field constants

esv8/generator/   Generates api_coverage_test.go and esclient_spec.go for v8
esv9/generator/   Same for v9
```

---

## Key Design Principles

### Type Safety
`estype.Field`, `estype.Index`, and `estype.Alias` are distinct string types. This means passing an index name where a field name is expected is a compile error — similar to how sqlc generates typed access from SQL schemas.

```go
// Good
func Search(ctx context.Context, index estype.Index) { ... }

// Won't compile — passing Field where Index is expected
client.Search(ctx, estype.Field("my_field"))
```

### Functional Options for Property Builders
All property constructors follow the functional-options pattern:

```go
prop := esv8.NewTextProperty(
    func(p *types.TextProperty) { p.Analyzer = ptr("kuromoji") },
)
```

### Fluent Builders in query Sub-Packages
Method chaining is intentionally confined to `esv8/query` and `esv9/query`:

```go
q := query.BoolQuery(
    query.NewBoolQuery().
        Must(query.TermValue(FieldStatus, "active")).
        Filter(query.TermsValues(FieldCategory, "a", "b")).
        Build(),
)

sort := query.NewSort().
    Field(FieldDate, sortorder.Desc).
    ScoreDesc().
    Build()
```

### High-Level Search vs SearchRaw
- `Search[T](ctx, client, alias, params)` — preferred high-level search API for application code. It decodes `_source` into caller-provided structs and exposes typed aggregation accessors through `SearchResponse[T]`.
- `SearchDocuments[T](...)` — convenience helper when only decoded documents are needed.
- `SearchOne[T](...)` — convenience helper for fetching the first decoded hit.
- `SearchRaw(ctx, alias, req)` — lower-level escape hatch for advanced Elasticsearch request shapes not yet modeled by `SearchParams`.

### ESClient vs ESClientSpec
- `ESClient` — curated set of common operations; preferred for application code. It keeps `SearchRaw` on the interface, while high-level typed search is provided by top-level helpers such as `Search[T](...)`.
- `ESClientSpec` — complete ES API coverage (generated from spec); use when `ESClient` lacks the needed endpoint.

---

## Coding Conventions

| Concern | Convention |
|---------|-----------|
| Field names | Always `estype.Field`, never bare `string` |
| Error wrapping | `fmt.Errorf("context: %w", err)` |
| Logging | `log/slog` only |
| Assertions in tests | `gotest.tools/v3/assert` (`assert.NilError`, `assert.Equal`, …) |
| Parallel tests | `t.Parallel()` at the top of every test and sub-test |
| Integration tests | `//go:build integration` build tag |
| Field names in tests | Generic names only: `status`, `title`, `category`, `tags`, `items`, `date`, `price`, `id`, `name`, `type`, `enabled`, `value` |
| Generated files | Header `// Code generated by estyped; DO NOT EDIT.` |
| v8 ↔ v9 parity | Changes to `esv8/` typically need the same change in `esv9/` |
| Search API guidance | Prefer `Search[T](...)` for normal application searches; reserve `SearchRaw(...)` for advanced escape-hatch scenarios |
| SearchBuilder output | `query.NewSearch().Build()` returns `query.SearchParams`; map it into `esv8.SearchParams` / `esv9.SearchParams` before calling `Search[T](...)` |
| Documentation style | No emoji; no excessive bold for inline emphasis (bold is reserved for list-item headings and table headers) |

---

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/elastic/go-elasticsearch/v8` | v8.19.3 | ES v8 typed client |
| `github.com/elastic/go-elasticsearch/v9` | v9.3.1 | ES v9 typed client |
| `github.com/google/uuid` | v1.6.0 | UUID generation |
| `gotest.tools/v3` | v3.5.2 | Test assertions |
| Go | 1.26 | Minimum required version |

---

## Common Patterns

### Generating field constants from an ES mapping

```bash
# mapping.json can be the full Get Mapping API response or just {"properties":{...}}
go tool estyped -mapping mapping.json -out esmodel/fields.go -package esmodel
```

### Generating field constants from a Go struct with JSON tags

Place the `//go:generate` directive in the file that defines the struct:

```go
//go:generate go tool estyped -struct Product -out product_fields.go

type Product struct {
    Status string `json:"status"`
    Items  []Item `json:"items"`
}

type Item struct {
    Name string `json:"name"`
}
```

The `-package` flag defaults to `$GOPACKAGE` when using `-struct`, so it can be omitted
inside a `go generate` run.

```go
// Generated output (constant mode)
const FieldStatus estype.Field = "status"
const FieldItemsColor estype.Field = "items.color"
```

### Adding a new property type wrapper

1. Add a `NewXxxProperty(opts ...XxxPropertyOption)` function to `esv8/property.go` and `esv9/property.go`.
2. Follow the functional-options pattern already used by `NewTextProperty`, `NewKeywordProperty`, etc.
3. Add a unit test in `esv8/property_test.go` (or `esv9/`).

### Adding a new query helper

1. Add the helper to `esv8/query/helpers.go` and `esv9/query/helpers.go`.
2. Write a test in the corresponding `helpers_test.go` using only generic field names.
