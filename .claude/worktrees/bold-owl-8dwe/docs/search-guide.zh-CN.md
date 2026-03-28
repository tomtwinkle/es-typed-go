# Search Guide

本文档仅保留搜索 API 的选择指引。

## API 选择

- `Search[T](...)`
  - 需要 typed hits、`Total`、aggregations 或 raw response 时使用
- `SearchDocuments[T](...)`
  - 只需要解码后的文档 `_source` 时使用
- `Count(...)`
  - 只需要匹配文档数量时使用
- `SearchRaw(...)`
  - 需要高层 helper 尚未覆盖的 Elasticsearch request shape 时使用

## 使用原则

- 默认优先使用 `Search[T](...)`
- 只取文档时使用 `SearchDocuments[T](...)`
- 不要为了 count 而使用 search
- 只有在高层 API 无法表达请求时才使用 `SearchRaw(...)`

## 补充说明

- `query.NewSearch().Build()` 返回的 params 可以直接传给高层搜索 helper
- `Limit(0)` 适用于不返回 hits、只关心 total 或 aggregations 的搜索请求
- `esv8` 与 `esv9` 的搜索 helper 用法应尽量保持一致

## 相关文档

- [../README.md](../README.md)
- [property-reference.md](property-reference.md)
- [contributing.md](contributing.md)