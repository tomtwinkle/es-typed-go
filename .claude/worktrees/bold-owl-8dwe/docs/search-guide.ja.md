# 検索ガイド

このドキュメントは、検索 API の選び方だけを簡潔にまとめたメモです。

## 使い分け

- `Search[T](...)`
  - hit metadata、total hits、aggregations、raw response も必要な場合に使います
- `SearchDocuments[T](...)`
  - decoded された document source だけ必要な場合に使います
- `Count(...)`
  - 件数だけ必要な場合に使います
- `SearchRaw(...)`
  - 高レベル helper で表現できない Elasticsearch request shape が必要な場合に使います

## 基本方針

- 通常のアプリケーション検索では `Search[T](...)` を優先します
- document だけ欲しい場合は `SearchDocuments[T](...)` を使います
- 件数確認のためだけに search helper は使わず、`Count(...)` を使います
- `SearchRaw(...)` は escape hatch として必要な場合にだけ使います

## 補足

- `query.NewSearch().Build()` が返す params は、そのまま高レベル helper に渡せます
- `esv8` と `esv9` で検索 helper の使い方はほぼ同じです

## 関連ドキュメント

- [../README.ja.md](../README.ja.md)
- [property-reference.ja.md](property-reference.ja.md)
- [contributing.ja.md](contributing.ja.md)