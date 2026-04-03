# examples/quickstart

`examples/quickstart` is a runnable companion to the top-level README.

It shows how to:

1. define a document struct
2. generate typed field constants with `estyped`
3. seed example documents into Elasticsearch
4. build type-safe queries and aggregations with `github.com/tomtwinkle/es-typed-go/query`
5. run a high-level typed search with `esv8.Search[T](...)`
6. demonstrate pagination by returning fewer hits than the total match count

## Files

- `main.go` — example application
- `product.go` — document model and `//go:generate` directive
- `esmodel/product_gen.go` — generated field constants (created by `go generate`)

## Generate field constants

From the repository root:

```bash
go generate ./examples/quickstart/...
```

This generates:

```text
examples/quickstart/esmodel/product_gen.go
```

The example uses the generated constants via the `esmodel` package, such as:

- `esmodel.FieldStatus`
- `esmodel.FieldCategory`
- `esmodel.FieldDate`
- `esmodel.FieldPrice`

## Run the example

If you have Elasticsearch running locally, you can run:

```bash
go run ./examples/quickstart
```

By default, the example is intended to connect to the local Elasticsearch instance described by the repository `compose.yaml`.

When the example starts, it:

1. recreates the demo index and alias
2. seeds a small set of example documents
3. runs a filtered search with sorting and aggregations
4. prints both the total number of matching documents and the number of hits returned in the current page

## Seeded data and query behavior

The seeded dataset intentionally includes both matching and non-matching documents.

The example query filters for:

- `status = active`
- `category = electronics`
- `date` between `2024-06-01` and `2024-12-31`

This means the seeded records are designed to show clear filtering behavior:

- `product-1` matches
- `product-4` matches
- `product-2` is excluded by category
- `product-3` is excluded by status
- `product-5` is excluded by category
- `product-6` is excluded by date

## Pagination behavior

The quickstart intentionally uses a page size of `1`.

That means the output is expected to show a difference between:

- total matched documents
- hits returned in the current page

For example, the console output may look like:

```text
Total hits: 2 (page size=1)
Raw total hits: value=2 relation=eq
Hits (1):
```

This is intentional: `Total hits` is the full match count, while `Hits (1)` is only the current page.

## Notes

- The generated file name is intentionally `product_gen.go` so it is obvious that the file is generated.
- This example is meant to mirror the README style: it uses generated field constants instead of handwritten `const Field...` declarations.
- The example recreates its demo index on each run so stale mappings do not affect the result.
- The example also prints the raw Elasticsearch total-hit metadata to make the pagination behavior easier to understand.
- If you modify the `Product` struct or its mapping metadata, run `go generate` again.