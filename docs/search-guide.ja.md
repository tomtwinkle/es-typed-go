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
        query.DateRangeQuery(esmodel.Product.Fields.Date, query.DateRangeGte("2024-01-01"), query.DateRangeLte("2024-12-31")),
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

### Nested aggregation

ネストされたオブジェクト内のドキュメントを集計し、そのスコープでサブ集計を実行します:

```go
avgDef := query.AvgAgg("avg_price", esmodel.Order.Fields.Items.Price)
nestedDef := query.NestedAgg("items", esmodel.Order.Fields.Items,
    query.NestedAggSubAggs(avgDef))

// ... 検索実行 ...

nested := resp.Aggregations.MustNested(nestedDef)
// nested.DocCount()  — ネストドキュメントの総数
avg, _ := nested.Aggregations().GetAvg(avgDef)
```

### Filter aggregation

クエリにマッチするドキュメントに集計を絞り込みます:

```go
avgDef := query.AvgAgg("avg_price", esmodel.Product.Fields.Price)
filterDef := query.FilterAgg("active_products",
    query.TermValue(esmodel.Product.Fields.Status, "active"),
    query.FilterAggSubAggs(avgDef),
)

// ... 検索実行 ...

filtered := resp.Aggregations.MustFilter(filterDef)
// filtered.DocCount()
avg, _ := filtered.Aggregations().GetAvg(avgDef)
```

### Multi-terms aggregation

複数フィールドでドキュメントをまとめてグループ化します:

```go
multiDef := query.MultiTermsAgg("by_category_status",
    []estype.Field{esmodel.Product.Fields.Category, esmodel.Product.Fields.Status},
    query.MultiTermsAggSize(10),
)

// ... 検索実行 ...

multi := resp.Aggregations.MustMultiTerms(multiDef)
for _, bucket := range multi.Buckets() {
    // bucket.Keys()     — 複合キーの値 []string
    // bucket.DocCount() — このキーの組み合わせのドキュメント数
    fmt.Println(bucket.Keys(), bucket.DocCount())
}
```

## ソート方向

`query.SortAsc` / `query.SortDesc` を使います。バージョン固有の `sortorder` パッケージのインポートは不要です:

```go
query.NewSort().Field(esmodel.Product.Fields.Date, query.SortDesc)
```

## クエリ helper リファレンス

### DateRangeQuery

`DateRangeQuery` は functional options を受け取るため、4種類の比較演算子を任意に組み合わせて使えます:

```go
// Gte + Lte（閉区間）
query.DateRangeQuery(field, query.DateRangeGte("2024-01-01"), query.DateRangeLte("2024-12-31"))

// Gt + Lt（開区間）
query.DateRangeQuery(field, query.DateRangeGt("2024-01-01"), query.DateRangeLt("2025-01-01"))

// 片側のみ
query.DateRangeQuery(field, query.DateRangeGte("2024-01-01"))
```

使用可能なオプション: `DateRangeGt`, `DateRangeGte`, `DateRangeLt`, `DateRangeLte`。

### MultiTermsAgg のフィールドごとの Missing

`query.MultiTermLookup` を使うとフィールドごとに設定を行えます。`Missing` には、そのフィールドを持たないドキュメントへの代替値を指定します:

```go
query.MultiTermsAgg("by_date_tz", []query.MultiTermLookup{
    {Field: esmodel.Item.Fields.BusinessDate},
    {Field: esmodel.Item.Fields.Timezone, Missing: "UTC"},
})
```

### Field.Ptr() / Field.String() — typed field の変換

raw go-elasticsearch 型が `*string` を要求する場合（例: `NestedAggregation.Path`、`SumAggregation.Field`）、一時変数の代わりに `Ptr()` を使えます:

```go
// Before
path := string(esmodel.Item.Fields.Items)
types.NestedAggregation{Path: &path}

// After
types.NestedAggregation{Path: esmodel.Item.Fields.Items.Ptr()}
```

フィールドが `string`（`*string` ではない）の場合、例えば `types.NestedSortValue.Path` では `String()` を使います（`Ptr()` は使用できません）:

```go
// NestedSortValue.Path は string 型 — Ptr() は使用不可
types.NestedSortValue{Path: esmodel.Item.Fields.Items.String()}
```

`Ptr()` と `String()` は `estype.Alias` と `estype.Index` でも使用できます。

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
