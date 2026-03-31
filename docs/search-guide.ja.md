# 検索ガイド

## どの API を使うか

| ユースケース | API |
|---|---|
| typed hits・metadata・aggregations が必要 | `Search[T](ctx, client, alias, params)` |
| decoded document `_source` のみ必要 | `SearchDocuments[T](ctx, client, alias, params)` |
| 件数確認のみ | `Count(ctx, alias)` |
| 高レベル helper で表現できない request shape | `SearchRaw(ctx, alias, req)` |

## 検索パラメータの構築

クエリ構築にはトップレベルの `query` パッケージを使います。このパッケージは version-agnostic で、`esv8` と `esv9` 両方で共通して使えます。

```go
import "github.com/tomtwinkle/es-typed-go/query"

params := query.NewSearch().
    Where(query.TermValue(esmodel.Product.Fields.Status, "active")).
    Where(
        query.TermValue(esmodel.Product.Fields.Category, "electronics"),
        query.DateRangeQuery(esmodel.Product.Fields.Date, "2024-01-01", "2024-12-31"),
    ).
    Sort(
        query.NewSort().
            Field(esmodel.Product.Fields.Date, query.SortDesc).
            ScoreDesc().
            Build()...,
    ).
    Aggregation(query.Aggs(
        query.AvgAgg("avg_price", esmodel.Product.Fields.Price),
    ).Build()).
    Limit(20).
    Offset(0).
    Build()
```

`params` は `Search[T]` にそのまま渡せます:

```go
// v8
resp, err := esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)

// v9 — 全く同じ、パッケージだけが異なる
resp, err := esv9.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

## v8 と v9 の切り替え

クエリ構築が共通の `query/` パッケージを使うため、バージョン切り替えはクライアントの import 変更のみで済みます。変更点の全体像は [migration-v2.md](migration-v2.md) を参照してください。

## Aggregation 結果の取得

`SearchResponse.Aggregations` は `query.AggResults` 型です。`GetXxx` / `MustXxx` メソッドで型付き結果を取得できます:

```go
avgDef := query.AvgAgg("avg_price", esmodel.Product.Fields.Price)
termsDef := query.StringTermsAgg("by_category", esmodel.Product.Fields.Category,
    query.WithSubAggs(avgDef))

// ... 検索実行 ...

terms := resp.Aggregations.MustStringTerms(termsDef)
for _, bucket := range terms.Buckets() {
    avg, _ := bucket.Aggregations().GetAvg(avgDef)
    // avg.Value() は *float64
}
```

## ソート方向

`query.SortAsc` / `query.SortDesc` を使います。バージョン固有の `sortorder` パッケージのインポートは不要です:

```go
query.NewSort().Field(esmodel.Product.Fields.Date, query.SortDesc)
```

## 補足

- `Limit(0)` は hits を返しません（aggregation のみ・件数確認のみの検索に便利）。
- 1件取得には `Limit(1)` を使います。
- `SearchRaw` は高レベル helper でカバーできない request shape のための escape hatch です。
- `esv8` と `esv9` で検索 helper の使い方は同じです。

## 関連ドキュメント

- [../README.ja.md](../README.ja.md)
- [migration-v2.md](migration-v2.md) — v2 アーキテクチャと移行手順
- [property-reference.ja.md](property-reference.ja.md)
- [contributing.md](contributing.md)
