---
name: go-lint-checks
description: Run repository lint checks, especially the custom type-assertion checker, and keep go vet clean.
---

# Go lint checks

Use this skill when you need to:

- validate repository-wide Go lint status
- enforce safe type assertions
- update CI or local commands for linting
- explain how linting is wired into this repository

## Purpose

This repository enforces an additional rule beyond standard `go vet`:

- when using a type assertion with the `comma, ok` form such as `v, ok := x.(T)`, the `ok` result must be explicitly checked before `v` is relied on
- direct chained assertions that hide the existence check, such as retrieving from a map and immediately asserting a type, should be rewritten into separate steps so the map lookup and the type assertion are both checked

Examples:

Good:

```go
raw, ok := m[key]
if !ok {
	return fmt.Errorf("missing key %q", key)
}
v, ok := raw.(MyType)
if !ok {
	return fmt.Errorf("unexpected type %T", raw)
}
```

Bad:

```go
v := m[key].(MyType)
```

Bad:

```go
v, ok := m[key].(MyType)
if !ok {
	return errSomething
}
```

The second example still misses the map lookup result check.

## Repository lint commands

### Primary lint command

Run the repository lint suite with:

```bash
go vet -vettool=$(go env GOPATH)/bin/okassertcheck ./...
```

If the custom analyzer is built from source in-module, use:

```bash
go build -o ./bin/okassertcheck ./tools/okassertcheck
go vet -vettool=$(pwd)/bin/okassertcheck ./...
```

### Recommended local sequence

Use this order when validating a change:

```bash
go test ./... -short -count=1
go build -o ./bin/okassertcheck ./tools/okassertcheck
go vet -vettool=$(pwd)/bin/okassertcheck ./...
```

### CI expectation

CI should run both:

1. unit tests
2. `go vet` with the custom analyzer enabled

## What to flag

Flag patterns like these:

- `v := x.(T)`
- `v, ok := expr.(T)` where `ok` is never checked
- `if !ok { ... }` missing after a `comma, ok` assertion
- map lookup + type assertion compressed into one expression:
  - `vals, ok := q.Terms.TermsQuery[string(FieldCategory)].([]types.FieldValue)`
  - this should be split into:
    - map lookup with `ok`
    - type assertion with `ok`

Also review tests with the same strictness as production code.

## Preferred fixes

### Split map lookup and type assertion

Instead of:

```go
vals, ok := m[key].([]T)
if !ok {
	t.Fatal("...")
}
```

write:

```go
raw, ok := m[key]
if !ok {
	t.Fatal("...")
}
vals, ok := raw.([]T)
if !ok {
	t.Fatal("...")
}
```

### Preserve intent in tests

In tests, use explicit assertions with clear failure messages when practical:

```go
rawAgg, ok := res.Aggregations["price_stats"]
assert.Assert(t, ok)
statsAgg, ok := rawAgg.(*types.StatsAggregate)
assert.Assert(t, ok, "expected *types.StatsAggregate")
```

### Do not weaken checks

Do not replace safe explicit checks with unchecked assertions just to make code shorter.

## Files commonly involved

- `tools/okassertcheck/` — custom analyzer implementation
- `.github/workflows/test.yml` — CI wiring
- `README.md`
- `README.ja.md`
- `.github/instructions/*.md`

## Maintenance guidance

When changing lint behavior:

1. update the analyzer implementation
2. update CI so the analyzer runs on pull requests
3. document the command in README if user-facing
4. update this skill if usage expectations changed

## Notes

- Standard `go vet` alone does not enforce this repository-specific rule.
- The custom analyzer is the source of truth for this check.
- Apply the same standards to `esv8` and `esv9` to preserve parity.