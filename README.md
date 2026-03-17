# es-typed-go

[![Test](https://github.com/tomtwinkle/es-typed-go/actions/workflows/test.yml/badge.svg)](https://github.com/tomtwinkle/es-typed-go/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tomtwinkle/es-typed-go.svg)](https://pkg.go.dev/github.com/tomtwinkle/es-typed-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**English** | [日本語](README.ja.md) | [中文](README.zh-CN.md)

A type-safe Go wrapper for [go-elasticsearch](https://github.com/elastic/go-elasticsearch) (v8 and v9) that prevents field-name typos and index/alias confusion at compile time.

## Motivation

### Problems with the existing elasticsearch-go Typed Client

The official [go-elasticsearch](https://github.com/elastic/go-elasticsearch) library provides a "typed" client, but despite the name, it is auto-generated from a TypeScript specification and has several practical problems:

1. **Unintuitive API** — Because the typed client is mechanically generated from TypeScript definitions, function signatures are often confusing and hard to understand without constantly consulting documentation.

2. **Pervasive use of `any` types** — When building search queries, most parameters (such as `FieldValue`, `SortCombinations`, `Missing`, and `TermsQueryField`) are typed as `any`. This means you can pass completely invalid parameters without any compiler error — you only discover the mistake when Elasticsearch rejects the query at runtime.

3. **Trial-and-error workflow** — Because of the lack of compile-time safety, developers are forced into a trial-and-error cycle: write code → send query → read the error → fix → repeat. In many cases, writing raw JSON is actually easier and faster than using the typed client.

### What es-typed-go solves

This library introduces type safety through:

- **Distinct types** for field names (`estype.Field`), index names (`estype.Index`), and alias names (`estype.Alias`) — passing an index where a field is expected is a compile error
- **Code generation** from Elasticsearch mappings — the `estyped` CLI reads your mapping JSON and generates typed field constants, similar to how [sqlc](https://sqlc.dev/) generates typed Go code from SQL schemas
- **Fluent, type-safe builders** for queries, sorts, and aggregations that accept `estype.Field` instead of bare strings
- **Functional-option constructors** for all 52+ Elasticsearch property types, making index mapping definitions safe and readable

Invalid usage is caught by the compiler, not by Elasticsearch at runtime.

## Features

- **Compile-time safety** — Distinct types (`Field`, `Index`, `Alias`) prevent mix-ups; `MappingProperty` interface eliminates `any` from field mapping definitions
- **Code generation** — Generate typed field constants from Elasticsearch mappings
- **SearchBuilder** — High-level ActiveRecord-style builder combining query, sort, aggregations, and pagination into a single `SearchParams`
- **Fluent query builders** — Type-safe Bool, Term, Match, Range, Nested, Prefix, Wildcard, MultiMatch, FunctionScore queries and more
- **Aggregation builders** — Terms, DateHistogram, Histogram, Avg, Max, Min, Sum, ValueCount, Cardinality, Stats, Nested, Filter
- **Sort builders** — Field, Score, Doc, GeoDistance, Script sorting with functional options
- **Property builders** — Functional-option constructors for all ES property types
- **Date format constants** — 80+ built-in Elasticsearch date format constants
- **Dual version support** — Elasticsearch v8 and v9 with identical APIs

## Installation

```bash
go get github.com/tomtwinkle/es-typed-go
```

This installs the core `estype` package and both `esv8` and `esv9` wrappers.

To use the `estyped` code-generation CLI, install it as a Go tool:

```bash
go get -tool github.com/tomtwinkle/es-typed-go/cmd/estyped
```

This adds an entry to your `go.mod` and lets you invoke the tool with `go tool estyped`. You can also install it globally so that `estyped` is available directly on your PATH:

```bash
go install github.com/tomtwinkle/es-typed-go/cmd/estyped@latest
```

## Quick Start

### 1. Define your mapping and generate field constants

Create your Elasticsearch mapping file:

```json
{
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "status": { "type": "keyword" },
      "price": { "type": "integer" },
      "category": { "type": "keyword" },
      "tags": { "type": "keyword" },
      "date": { "type": "date" },
      "items": {
        "type": "nested",
        "properties": {
          "name": { "type": "text" },
          "value": { "type": "integer" }
        }
      }
    }
  }
}
```

Generate typed field constants from an Elasticsearch mapping file:

```bash
go tool estyped \
  -mapping mapping.json \
  -out esmodel/fields.go \
  -package esmodel
```

Alternatively, if you already have a Go struct with JSON tags that represents your document model, you can generate directly from it by adding a `//go:generate` directive to the struct file:

```go
//go:generate go tool estyped -struct Product -out product_fields.go

type Product struct {
    Status   string   `json:"status"`
    Title    string   `json:"title"`
    Category string   `json:"category"`
    Items    []Item   `json:"items"`
}

type Item struct {
    Name  string `json:"name"`
    Price int    `json:"price"`
}
```

To give the generator accurate ES type names for each field, implement `estype.ESMappingProvider` on the struct. Without it, all field types fall back to `"unknown"` in generated comments. See [ESMappingProvider](#esmappingprovider) for details.

When using `-struct`, the `-package` flag defaults to `$GOPACKAGE` (set automatically by `go generate`). Run `go generate ./...` to regenerate. If you installed `estyped` globally with `go install`, you can also use the shorter form:

```go
//go:generate estyped -struct Product -out product_fields.go
```

This generates:

```go
// Code generated by estyped; DO NOT EDIT.
package esmodel

import "github.com/tomtwinkle/es-typed-go/estype"

// FieldCategory is the "category" field (type: keyword).
const FieldCategory estype.Field = "category"

// FieldDate is the "date" field (type: date).
const FieldDate estype.Field = "date"

// FieldItems is the "items" field (type: nested).
const FieldItems estype.Field = "items"

// FieldItemsName is the "items.name" field (type: text).
const FieldItemsName estype.Field = "items.name"

// FieldItemsValue is the "items.value" field (type: integer).
const FieldItemsValue estype.Field = "items.value"

// FieldPrice is the "price" field (type: integer).
const FieldPrice estype.Field = "price"

// FieldStatus is the "status" field (type: keyword).
const FieldStatus estype.Field = "status"

// FieldTags is the "tags" field (type: keyword).
const FieldTags estype.Field = "tags"

// FieldTitle is the "title" field (type: text).
const FieldTitle estype.Field = "title"

// FieldTitleKeyword is the "title.keyword" field (type: keyword).
const FieldTitleKeyword estype.Field = "title.keyword"
```

You can also use struct mode for grouped access:

```bash
go tool estyped \
  -mapping mapping.json \
  -out esmodel/fields.go \
  -package esmodel \
  -name Product
```

This generates:

```go
// Code generated by estyped; DO NOT EDIT.
package esmodel

import "github.com/tomtwinkle/es-typed-go/estype"

// Product provides typed field names for the Elasticsearch index mapping.
var Product = struct {
	Category    estype.Field
	Date        estype.Field
	Items       estype.Field
	Items_Name  estype.Field
	Items_Value estype.Field
	Price       estype.Field
	Status      estype.Field
	Tags        estype.Field
	Title       estype.Field
	Title_Keyword estype.Field
}{
	Category:      "category",
	Date:          "date",
	Items:         "items",
	Items_Name:    "items.name",
	Items_Value:   "items.value",
	Price:         "price",
	Status:        "status",
	Tags:          "tags",
	Title:         "title",
	Title_Keyword: "title.keyword",
}

// Usage: esmodel.Product.Status, esmodel.Product.Items_Name, etc.
```

### 2. Build type-safe queries

```go
package main

import (
	"github.com/tomtwinkle/es-typed-go/esv8/query"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// Using generated field constants
const (
	FieldStatus   estype.Field = "status"
	FieldCategory estype.Field = "category"
	FieldDate     estype.Field = "date"
	FieldPrice    estype.Field = "price"
)

func buildQuery() {
	q := query.BoolQuery(query.NewBoolQuery().
		Must(
			query.TermValue(FieldStatus, "active"),
		).
		Filter(
			query.TermsValues(FieldCategory, "electronics", "books"),
			query.DateRangeQuery(FieldDate, "2024-01-01", "2024-12-31"),
		).
		Build(),
	)

	_ = q // Use with ESClient.Search() or ESClient.SearchWithRequest()
}
```

### 3. Create an Elasticsearch client

```go
package main

import (
	"log/slog"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/tomtwinkle/es-typed-go/esv8"
)

func main() {
	client, err := esv8.NewClientWithLogger(
		es8.Config{
			Addresses: []string{"http://localhost:9200"},
		},
		slog.Default(),
	)
	if err != nil {
		panic(err)
	}
	_ = client // Use for search, indexing, etc.
}
```

## Usage Guide

### Core Types

The foundation of es-typed-go is three distinct string types that prevent mix-ups at compile time:

```go
import "github.com/tomtwinkle/es-typed-go/estype"

// These are distinct types — you cannot accidentally pass one where another is expected
var field estype.Field = "status"       // Elasticsearch field name
var index estype.Index = "my-index"     // Elasticsearch index name
var alias estype.Alias = "my-alias"     // Elasticsearch alias name

// OK — correct usage
client.Search(ctx, alias, query, ...)
client.DeleteIndex(ctx, index)

// Compile error — passing Field where Alias is expected
client.Search(ctx, field, query, ...)

// Compile error — passing Index where Alias is expected
client.Search(ctx, index, query, ...)
```

### Query Builders

#### Query Helper Functions

Convenience functions for constructing `types.Query` values:

```go
import "github.com/tomtwinkle/es-typed-go/esv8/query"

// Term queries
query.TermValue(FieldStatus, "active")
query.TermsValues(FieldCategory, "electronics", "books")

// Match queries
query.MatchValue(FieldTitle, "search keyword")
query.MatchPhrase(FieldTitle, "exact phrase")
query.MultiMatchQuery("search text", FieldTitle, FieldName)

// Match all / match none
query.MatchAll()
query.MatchNone()

// IDs query
query.IdsQuery("id1", "id2", "id3")

// Prefix / wildcard queries
query.PrefixValue(FieldTitle, "go")
query.WildcardValue(FieldTitle, "go*")

// Field existence
query.ExistsField(FieldStatus)
query.NotExists(FieldPrice)

// Nested queries
query.NestedFilter(FieldItems,
	query.TermValue(estype.Field("items.name"), "widget"),
)

// Range queries
query.DateRangeQuery(FieldDate, "2024-01-01", "2024-12-31")
gte, lte := types.Float64(100), types.Float64(500)
query.NumberRangeQuery(FieldPrice, &gte, &lte)

// Bool shorthand helpers
query.BoolMust(q1, q2)
query.BoolFilter(q1, q2)
query.BoolShould(q1, q2)
query.BoolMustNot(q1)

// Wrap a BoolQuery in a Query
query.BoolQuery(query.NewBoolQuery().
	Must(
		query.TermValue(FieldStatus, "active"),
		query.MatchPhrase(FieldTitle, "search keyword"),
	).
	Filter(
		query.DateRangeQuery(FieldDate, "2024-01-01", "2024-12-31"),
	).
	Should(
		query.TermValue(FieldCategory, "premium"),
	).
	MustNot(
		query.ExistsField(FieldPrice),
	).
	MinimumShouldMatch(1).
	Build(),
)

// Function score query
query.FunctionScoreQuery(&types.FunctionScoreQuery{
	Query: &types.Query{MatchAll: &types.MatchAllQuery{}},
})

// Convert a typed slice to []types.FieldValue for TermsQuery
ids := query.FieldValues("id1", "id2", "id3")
query.TermsValues(FieldStatus, ids...)
```

### Sort Builder

```go
import (
	"github.com/tomtwinkle/es-typed-go/esv8/query"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortmode"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptsorttype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/distanceunit"
)

sorts := query.NewSort().
	Field(FieldDate, sortorder.Desc).                      // Sort by date descending
	Field(FieldPrice, sortorder.Asc).                      // Then by price ascending
	FieldWithMissing(FieldCategory, sortorder.Asc,
		query.MissingLast).                                // Missing values last
	FieldNested(FieldItems, sortorder.Asc,
		FieldItems, sortmode.Min).                         // Nested field sort
	ScoreDesc().                                           // By relevance score descending
	ScoreAsc().                                            // By relevance score ascending
	DocAsc().                                              // By index order ascending
	DocDesc().                                             // By index order descending
	GeoDistance(FieldLocation, types.GeoLocation{...},
		sortorder.Asc,
		query.WithGeoDistanceUnit(distanceunit.Kilometers),
		query.WithGeoDistanceIgnoreUnmapped(true)).        // Geo distance sort with options
	Script(script, scriptsorttype.Number,
		sortorder.Asc).                                    // Script-based sort
	Build()
```

### Aggregation Builder

```go
import "github.com/tomtwinkle/es-typed-go/esv8/query"

aggs := query.NewAggregations().
	Terms("by_category", FieldCategory).                                         // Bucket by category
	TermsWithSize("top_tags", FieldTags, 20).                                    // Top 20 tags
	DateHistogram("over_time", FieldDate, calendarinterval.Month).               // Monthly histogram
	DateHistogramWithFormat("over_time_fmt", FieldDate, "yyyy-MM",
		calendarinterval.Month).                                                 // With date format
	Histogram("price_dist", FieldPrice, 50.0).                                   // Numeric histogram
	Avg("avg_price", FieldPrice).                                                // Average price
	Max("max_price", FieldPrice).                                                // Maximum price
	Min("min_price", FieldPrice).                                                // Minimum price
	Sum("total_price", FieldPrice).                                              // Sum
	ValueCount("count_status", FieldStatus).                                     // Value count
	Cardinality("unique_categories", FieldCategory).                             // Distinct count
	Stats("price_stats", FieldPrice).                                            // Count/min/max/avg/sum
	Nested("nested_items", FieldItems, query.NewAggregations().                  // Nested aggregation
		Terms("item_names", estype.Field("items.name")),
	).
	Filter("active_only", query.TermValue(FieldStatus, "active"),
		query.NewAggregations().Avg("avg_price", FieldPrice)).                   // Filter aggregation
	SubAggregations("by_category", query.NewAggregations().                      // Sub-aggregation
		Avg("avg_price", FieldPrice),
	).
	Build()
```

### SearchBuilder

`query.NewSearch()` provides an ActiveRecord-style builder that combines query, sort, aggregations, and pagination into a single `SearchParams` value. Use it instead of assembling search parameters manually.

```go
import (
	"github.com/tomtwinkle/es-typed-go/esv8/query"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

params := query.NewSearch().
	Where(
		query.TermValue(FieldStatus, "active"),              // filter clauses
		query.DateRangeQuery(FieldDate, "2024-01-01", ""),
	).
	Must(
		query.MatchPhrase(FieldTitle, "search keyword"),     // must clauses
	).
	Should(
		query.TermValue(FieldCategory, "premium"),           // should clauses
	).
	MustNot(
		query.ExistsField(FieldPrice),                       // must_not clauses
	).
	Sort(query.NewSort().
		Field(FieldDate, sortorder.Desc).
		ScoreDesc().
		Build()...,
	).
	Aggregation(query.NewAggregations().
		Terms("by_category", FieldCategory).
		Build(),
	).
	Limit(10).
	Offset(0).
	Build()

// params.Query, params.Sort, params.Aggregations, params.Size, params.From, etc.
resp, err := client.Search(ctx, alias, params.Query, params.Size, params.From,
	params.Sort, params.Aggregations, params.Highlight, params.Collapse, params.ScriptFields)
```

You can also set the query directly when you have already built a `types.Query`:

```go
params := query.NewSearch().
	Query(query.BoolQuery(query.NewBoolQuery().
		Must(query.TermValue(FieldStatus, "active")).
		Build()),
	).
	Limit(20).
	Build()
```

### Property Builders (Index Mappings)

Define index mappings with type-safe functional-option constructors:

```go
import (
	"github.com/tomtwinkle/es-typed-go/esv8"
	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

mappings := &types.TypeMapping{
	Properties: map[string]types.Property{
		"title": esv8.NewTextProperty(
			esv8.WithTextAnalyzer("standard"),
			esv8.WithTextRawKeyword(256),  // Adds a .keyword sub-field
		),
		"status": esv8.NewKeywordProperty(),
		"price": esv8.NewIntegerNumberProperty(
			esv8.WithIntegerNumberCoerce(true),
		),
		"date": esv8.NewDateProperty(
			esv8.WithDateFormat(
				estype.DateFormatStrictDateTime,
				estype.DateFormatEpochMillis,
			),
		),
		"enabled": esv8.NewBooleanProperty(),
		"location": esv8.NewGeoPointProperty(),
		"tags": esv8.NewKeywordProperty(
			esv8.WithKeywordIgnoreAbove(256),
		),
		"items": esv8.NewNestedProperty(
			func(p *types.NestedProperty) {
				p.Properties = map[string]types.Property{
					"name":  esv8.NewTextProperty(),
					"value": esv8.NewIntegerNumberProperty(),
				}
			},
		),
	},
}
```

### ESMappingProvider

When using `estyped -struct`, implement `estype.ESMappingProvider` on your struct to tell the generator the Elasticsearch type name of each field. Without this method, every field type falls back to `"unknown"` in generated comments.

`MappingField.Property` accepts any value that implements `estype.MappingProperty` (the `ESTypeName() string` interface):

| Property value | When to use |
|---|---|
| `estype.FieldType("integer")` | Plain ES type name: `keyword`, `text`, `integer`, `long`, `float`, `double`, `boolean`, `date`, `object`, `nested`, `geo_point`, `dense_vector`, … |
| `estype.NewTextProperty(...)` | Text field with analyzer or multi-field sub-properties |
| `estype.NewKeywordProperty(...)` | Keyword field with `ignore_above` or similar options |

`Path` is a dot-separated JSON field path and is typed as `string`. Because paths originate from JSON keys, using `string` avoids the need for an explicit type conversion when building a mapping programmatically.

```go
//go:generate go tool estyped -struct Product -out product_fields.go

type Product struct {
	Status string   `json:"status"`
	Title  string   `json:"title"`
	Price  int      `json:"price"`
	Tags   []string `json:"tags"`
	Items  []Item   `json:"items"`
}

type Item struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (Product) ESMapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status",      Property: estype.NewKeywordProperty()},
			{Path: "title",       Property: estype.NewTextProperty(
				estype.WithField("keyword", estype.NewKeywordProperty(estype.WithIgnoreAbove())),
				estype.WithSearchAnalyzer(estype.Analyzer("my_search_analyzer")),
				estype.WithIndexAnalyzer(estype.Analyzer("my_index_analyzer")),
			)},
			{Path: "price",       Property: estype.FieldType("integer")},
			{Path: "tags",        Property: estype.FieldType("keyword")},
			{Path: "items",       Property: estype.FieldType("nested")},
			{Path: "items.name",  Property: estype.NewTextProperty()},
			{Path: "items.value", Property: estype.FieldType("integer")},
		},
	}
}
```

### Date Format Constants

80+ built-in Elasticsearch date format constants:

```go
import "github.com/tomtwinkle/es-typed-go/estype"

// Use predefined constants
estype.DateFormatEpochMillis       // "epoch_millis"
estype.DateFormatStrictDate        // "strict_date" (yyyy-MM-dd)
estype.DateFormatStrictDateTime    // "strict_date_time" (yyyy-MM-dd'T'HH:mm:ss.SSSZ)
estype.DateFormatBasicDate         // "basic_date" (yyyyMMdd)

// Combine multiple formats with JoinDateFormats
format := estype.JoinDateFormats(
	estype.DateFormatEpochMillis,
	estype.DateFormatStrictDate,
)
// Result: "epoch_millis||strict_date"
```

### ESClient vs ESClientSpec

Two client interfaces are available:

| Interface | Use Case |
|-----------|----------|
| `ESClient` | Curated set of common operations (search, index, alias management). Preferred for application code. |
| `ESClientSpec` | Full Elasticsearch API coverage (auto-generated from spec). Use when `ESClient` lacks a needed endpoint. |

```go
// ESClient — curated, easy to use
client, _ := esv8.NewClientWithLogger(config, logger)

// ESClientSpec — full API coverage
specClient, _ := esv8.NewSpecClient(config)
```

### ESClient Operations

`ESClient` groups its methods into four categories:

**Cluster**
- `Info(ctx)` — cluster information

**Index management**
- `CreateIndex(ctx, index, settings, mappings)` — create an index
- `DeleteIndex(ctx, index)` — delete an index
- `IndexExists(ctx, index) bool` — check existence
- `IndexRefresh(ctx, index)` — force a refresh
- `IndexDocumentCount(ctx, index)` — document count

**Alias management**
- `CreateAlias(ctx, index, alias, isWriteIndex)` — create an alias
- `UpdateAliases(ctx, actions)` — atomic add/remove alias actions
- `AliasExists(ctx, alias) bool` — check existence
- `AliasRefresh(ctx, alias)` — force a refresh
- `GetIndicesForAlias(ctx, alias) []Index` — list backing indices
- `GetRefreshInterval(ctx, alias)` — read refresh interval
- `UpdateRefreshInterval(ctx, alias, interval)` — update refresh interval

**Documents**
- `CreateDocument(ctx, alias, id, doc)` — index a document (waits for refresh)
- `GetDocument(ctx, alias, id)` — retrieve a document by ID
- `DeleteDocument(ctx, index, id)` — delete a document
- `UpdateDocument(ctx, index, id, req)` — partial update

**Search**
- `Search(ctx, alias, query, limit, offset, sort, aggs, highlight, collapse, scriptFields)` — execute a search
- `SearchWithRequest(ctx, alias, req)` — execute a raw `search.Request`

**Reindex**
- `Reindex(ctx, srcIndex, dstIndex, waitForCompletion)` — full reindex
- `DeltaReindex(ctx, srcIndex, dstIndex, since, timestampField, waitForCompletion)` — incremental reindex
- `WaitForTaskCompletion(ctx, taskID, timeout)` — poll until a task finishes

### Complete Search Example

```go
package main

import (
	"context"
	"fmt"
	"log/slog"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv8"
	"github.com/tomtwinkle/es-typed-go/esv8/query"
)

// Generated field constants (from estyped CLI)
const (
	FieldStatus   estype.Field = "status"
	FieldCategory estype.Field = "category"
	FieldDate     estype.Field = "date"
	FieldPrice    estype.Field = "price"
)

func main() {
	// Create client
	client, err := esv8.NewClientWithLogger(
		es8.Config{Addresses: []string{"http://localhost:9200"}},
		slog.Default(),
	)
	if err != nil {
		panic(err)
	}

	// Build search parameters with SearchBuilder
	params := query.NewSearch().
		Where(query.TermValue(FieldStatus, "active")).
		Where(
			query.TermsValues(FieldCategory, "electronics", "books"),
			query.DateRangeQuery(FieldDate, "2024-01-01", "2024-12-31"),
		).
		Sort(query.NewSort().
			Field(FieldDate, sortorder.Desc).
			ScoreDesc().
			Build()...,
		).
		Aggregation(query.NewAggregations().
			Terms("by_category", FieldCategory).
			Avg("avg_price", FieldPrice).
			Build(),
		).
		Limit(10).
		Offset(0).
		Build()

	// Execute search
	ctx := context.Background()
	alias := estype.Alias("my-alias")
	resp, err := client.Search(ctx, alias, params.Query, params.Size, params.From,
		params.Sort, params.Aggregations, params.Highlight, params.Collapse, params.ScriptFields)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Total hits: %d\n", resp.Hits.Total.Value)
}
```

## Elasticsearch v8 / v9 Support

es-typed-go supports both Elasticsearch v8 and v9 with identical APIs. Simply swap the import path:

```go
// For Elasticsearch v8
import (
	"github.com/tomtwinkle/es-typed-go/esv8"
	"github.com/tomtwinkle/es-typed-go/esv8/query"
)

// For Elasticsearch v9
import (
	"github.com/tomtwinkle/es-typed-go/esv9"
	"github.com/tomtwinkle/es-typed-go/esv9/query"
)
```

All builders, helpers, and property constructors have the same signatures across both versions. Changes to `esv8/` are always mirrored in `esv9/`.

## Repository Structure

```
estype/            Core shared types (Field, Index, Alias, DateFormat, mapping parser)
esv8/              Elasticsearch v8 wrapper
  query/           Query, sort, and aggregation builders for v8
  generator/       Code generator for v8 API coverage tests
esv9/              Elasticsearch v9 wrapper (mirrors esv8)
  query/           Query, sort, and aggregation builders for v9
  generator/       Code generator for v9 API coverage tests
cmd/estyped/       CLI tool: generates typed Field constants from ES mappings
```

## Requirements

- Go 1.26 or later
- Elasticsearch v8.x or v9.x

## License

[MIT](LICENSE)
