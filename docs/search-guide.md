# Search Guide

This guide only keeps the search API selection notes for `es-typed-go`.

## Which API should I use?

- `Search[T](...)`
  - Use when you need typed hits, hit metadata, total hits, or aggregations.
  - This is the normal default for application searches.

- `SearchDocuments[T](...)`
  - Use when you only need decoded document `_source` values.
  - Prefer this when hit metadata and aggregations are unnecessary.

- `Count(...)`
  - Use when you only need the number of matching documents.
  - Prefer this over search helpers for count-only use cases.

- `SearchRaw(...)`
  - Use only for advanced Elasticsearch request shapes that are not modeled by the high-level helpers.

## Notes

- For a single result, use `Search[T](...)` or `SearchDocuments[T](...)` with `Limit(1)` or `Size: 1`.
- `Limit(0)` means no hits are returned. This can still be useful for aggregation-only or total-hit-oriented searches.
- `query.NewSearch().Build()` can be passed directly to the high-level search helpers.
- The same API selection guidance applies to both `esv8` and `esv9`.

## Related documents

- [../README.md](../README.md) — concise top-level overview
- [property-reference.md](property-reference.md) — property builder reference
- [contributing.md](contributing.md) — contributor setup and validation steps