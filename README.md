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
- direct use of builder search params with typed search helpers
- dual-version support for Elasticsearch v8 and v9

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
	"github.com/tomtwinkle/es-typed-go/esv8/query"
	"github.com/tomtwinkle/es-typed-go/estype"
)

type Product struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

var FieldStatus estype.Field = "status"

func main() {
	ctx := context.Background()

	client, err := esv8.NewClient(...)
	if err != nil {
		panic(err)
	}

	alias := estype.Alias("products")

	params := query.NewSearch().
		Where(query.TermValue(FieldStatus, "active")).
		Limit(10).
		Build()

	_, _ = esv8.Search[Product](ctx, client, alias, params)
}
```

For a runnable end-to-end example, see:

- [`examples/quickstart/main.go`](examples/quickstart/main.go)
- [`examples/quickstart/README.md`](examples/quickstart/README.md)

## Documentation

Detailed documentation has been moved under `docs/`.

### User guides
- [Search Guide](docs/search-guide.md)
- [Property Reference](docs/property-reference.md)
- [Documentation index](docs/README.md)

### Contributor guides
- [Contributing Guide](docs/contributing.md)

## Version support

- `esv8` targets Elasticsearch v8
- `esv9` targets Elasticsearch v9

The two packages are intended to stay closely aligned.

## License

[MIT](LICENSE)