---
applyTo: "**/*_test.go"
---

# Testing Conventions

## Framework

- Use `gotest.tools/v3/assert` for all assertions — never `testing.T` methods directly.
- Common assertions: `assert.NilError(t, err)`, `assert.Equal(t, got, want)`, `assert.Assert(t, condition)`, `assert.ErrorContains(t, err, "substring")`.

## Test Structure

- Always call `t.Parallel()` at the start of every test and sub-test.
- Use table-driven tests (`map[string]struct{ input, want }{ "name": {} }`) with `t.Run(name, ...)` for multiple cases.
- Name unit test functions `TestFunctionName(t *testing.T)`.
- Name integration test functions `TestIntegration_Feature(t *testing.T)`.

## Build Tags

- Integration tests (those requiring a real Elasticsearch server) must begin with:
  ```go
  //go:build integration
  ```
- Unit tests must not include any build tags.

## Field Names in Tests

- Use only generic Elasticsearch field names for test fixtures: `status`, `title`, `category`, `tags`, `items`, `date`, `price`, `id`, `name`, `type`, `enabled`, `value`.
- Always wrap field name literals in `estype.Field(...)`.

## Resource Cleanup

- Register cleanup with `t.Cleanup(func() { ... })`.
- Use `t.Helper()` inside test helper functions.

## Package Naming for Tests

- External test packages (`package estype_test`) are preferred when testing the public API.
- Use the same package name (e.g., `package query`) for white-box tests that require access to unexported identifiers.

## Test Data and Code Generation Tests

- `cmd/estyped/main_test.go` uses `go/ast` + `go/parser` to validate generated source:
  - `collectConstDecls()` for constant mode, `collectStructVarFields()` for struct mode.
  - `assertImport()` to verify import paths in generated code.
