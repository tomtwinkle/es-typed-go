# es-typed-go

[![Test](https://github.com/tomtwinkle/es-typed-go/actions/workflows/test.yml/badge.svg)](https://github.com/tomtwinkle/es-typed-go/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tomtwinkle/es-typed-go.svg)](https://pkg.go.dev/github.com/tomtwinkle/es-typed-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**English** | [日本語](README.ja.md) | [中文](README.zh-CN.md)

A type-safe Go wrapper for [go-elasticsearch](https://github.com/elastic/go-elasticsearch) (v8 and v9) that prevents field-name typos and index/alias confusion at compile time.

## Why this library?

The official Elasticsearch Go typed client is powerful, but in practice it still exposes many request shapes that are difficult to use safely and ergonomically in normal application code.

`es-typed-go` improves that experience with:

- distinct types for field names, index names, and alias names
- typed field accessor generation from mappings or Go structs
- fluent builders for queries, sorting, and aggregations
- functional-option builders for Elasticsearch property definitions
- high-level typed search helpers that decode `_source` into Go structs

## Features

- compile-time separation of `Field`, `Index`, and `Alias`
- code generation with `estyped`
- typed query / sort / aggregation builders
- Elasticsearch property builders
- high-level search helpers for v8 and v9
- **single import path to switch between Elasticsearch v8 and v9**
- version-agnostic `query/` package shared by both versions

## Installation

Install the library:

```bash
go get github.com/tomtwinkle/es-typed-go
```

Install the code generator as a Go tool:

```bash
go get -tool github.com/tomtwinkle/es-typed-go/cmd/estyped
```

Run it with:

```bash
go tool estyped
```

## Quick example

```go
package main

import (
	"context"

	"github.com/tomtwinkle/es-typed-go/esv8"
	"github.com/tomtwinkle/es-typed-go/query"
	"github.com/tomtwinkle/es-typed-go/examples/quickstart/esmodel"
)

type Product struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func main() {
	ctx := context.Background()

	client, err := esv8.NewClient(...)
	if err != nil {
		panic(err)
	}

	// esmodel.Product.Alias and .Index come from the generated model file.
	// esmodel.Product.Fields.Status is a typed field name (estype.Field).
	params := query.NewSearch().
		Where(query.TermValue(esmodel.Product.Fields.Status, "active")).
		Limit(10).
		Build()

	_, _ = esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)
}
```

## Switching from v8 to v9

Because query building uses the shared `query/` package, switching Elasticsearch versions requires changing **only the client import**:

```go
// v8
import (
    es8 "github.com/elastic/go-elasticsearch/v8"
    "github.com/tomtwinkle/es-typed-go/esv8"
    "github.com/tomtwinkle/es-typed-go/query"
)
client, _ := esv8.NewClient(es8.Config{Addresses: []string{"http://localhost:19200"}})
resp, err := esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)

// v9 — only the two lines above change
import (
    es9 "github.com/elastic/go-elasticsearch/v9"
    "github.com/tomtwinkle/es-typed-go/esv9"
    "github.com/tomtwinkle/es-typed-go/query"
)
client, _ := esv9.NewClient(es9.Config{Addresses: []string{"http://localhost:19201"}})
resp, err := esv9.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

All query building, field names, aggregations, and sort definitions are unchanged.

## Generated model format

Run `estyped` with `-struct` and `-group` to generate a unified model accessor that includes typed field names, the canonical alias, and the canonical index name in one place:

```go
// esdefinition/product.go
func (Product) Alias() estype.Alias { return "product" }
func (Product) Index() estype.Index { return "product-000001" }

//go:generate go tool estyped -struct Product -package esmodel -out ../esmodel/product_gen.go -group Product
```

The generated accessor:

```go
// esmodel/product_gen.go (generated — do not edit)
var Product = struct {
    Fields struct {
        Status   estype.Field
        Category estype.Field
        // ...
    }
    Alias estype.Alias
    Index estype.Index
}{
    Fields: struct{ ... }{Status: "status", Category: "category", ...},
    Alias: "product",
    Index: "product-000001",
}
```

Usage:

```go
query.TermValue(esmodel.Product.Fields.Status, "active")
esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

For runnable end-to-end examples, see:

- [`examples/quickstart/main.go`](examples/quickstart/main.go) — Elasticsearch v8
- [`examples/quickstart_v9/main.go`](examples/quickstart_v9/main.go) — Elasticsearch v9

## Documentation

Detailed documentation is under `docs/`.

### User guides
- [Search Guide](docs/search-guide.md)
- [Property Reference](docs/property-reference.md)
- [Migration Guide (v2)](docs/migration-v2.md)
- [Documentation index](docs/README.md)

### Contributor guides
- [Contributing Guide](docs/contributing.md)

## Version support

- `esv8` targets Elasticsearch v8
- `esv9` targets Elasticsearch v9

Both packages share the top-level `query/` package for query building and expose identical API signatures. The only difference between them is the underlying Elasticsearch client version.

## License

[MIT](LICENSE)
