# 検索ガイド

このドキュメントでは、`es-typed-go` が提供する検索 API、それぞれを使うべき場面、そして低レベルな Elasticsearch request との関係を説明します。

トップレベルの `README*` は簡潔に保ち、検索 API の詳細はこのドキュメントにまとめます。

## 概要

`es-typed-go` には、主に次の 3 つの検索入口があります。

- `Search[T](...)` — hit metadata、total hits、aggregations、raw response を含む高レベル typed search
- `SearchDocuments[T](...)` — decoded された document source だけを返す高レベル helper
- `SearchRaw(...)` — 高度な Elasticsearch request shape のための低レベル escape hatch

通常のアプリケーションコードでは、次の方針を推奨します。

- hit metadata や aggregation も必要なら `Search[T](...)`
- decoded した `_source` だけ必要なら `SearchDocuments[T](...)`
- 件数だけ必要なら `Count(...)`
- 高レベル helper で表現できない request shape が必要なときだけ `SearchRaw(...)`

## 検索 helper の使い分け

### `Search[T](...)`

次のような情報が必要な場合は `Search[T](...)` を使います。

- decoded document
- total hit count
- `_id`、`_index`、`_score` などの hit metadata
- typed aggregation access
- raw Elasticsearch response

例:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Offset(0).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
if err != nil {
    return err
}

fmt.Println(resp.Total)
for _, hit := range resp.Hits {
    fmt.Println(hit.ID, hit.Index, hit.Source)
}
```

### `SearchDocuments[T](...)`

decoded された `_source` だけが必要で、hit metadata や aggregation が不要なら `SearchDocuments[T](...)` を使います。

例:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
if err != nil {
    return err
}

for _, doc := range docs {
    fmt.Println(doc)
}
```

### `SearchRaw(...)`

高レベル helper では扱っていない request shape が必要な場合は `SearchRaw(...)` を使います。

代表例:

- `search_after`
- point-in-time
- custom `_source` filtering
- typed client の request struct を直接組み立てたい場合
- まだライブラリが wrapper を提供していない Elasticsearch 機能

例:

```go
req := search.NewRequest()
req.Query = &types.Query{
    MatchAll: &types.MatchAllQuery{},
}

rawResp, err := client.SearchRaw(ctx, alias, req)
if err != nil {
    return err
}

fmt.Println(rawResp.Hits.Total)
```

## どの API を選ぶべきか

### typed hit と metadata が欲しい

`Search[T](...)` を使います。

通常のアプリケーション検索では、これが最も自然なデフォルトです。

### decoded document だけ欲しい

`SearchDocuments[T](...)` を使います。

これは `Search[T](...)` の convenience helper です。

### 1 件だけ欲しい

`Search[T](...)` または `SearchDocuments[T](...)` に `Limit(1)` あるいは `Size: 1` を指定します。

例:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(1).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
if err != nil {
    return err
}

if len(resp.Hits) == 0 {
    fmt.Println("not found")
    return nil
}

fmt.Println(resp.Hits[0].Source)
```

decoded document だけでよい場合:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(1).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
if err != nil {
    return err
}

if len(docs) == 0 {
    fmt.Println("not found")
    return nil
}

fmt.Println(docs[0])
```

## 件数だけ知りたいなら search を使わない

一致件数だけが必要なら、`Count(...)` を優先してください。

例:

```go
res, err := client.Count(ctx, alias)
if err != nil {
    return err
}

fmt.Println(res.Count)
```

これを推奨する理由:

- 意図が明確
- API が単純
- 不要な hit を返させなくてよい

### `Limit(0)` はどう扱うか

`Limit(0)` は「hit を返さない」という意味です。

次の用途では still useful です。

- aggregation だけ欲しい
- search response の total hits を見たい
- document payload なしで search shape の request を使いたい

ただし、目的が純粋に「何件一致したか」を知ることなら、通常は `Count(...)` の方が適切です。

## query builder の params をそのまま渡す

`query.NewSearch().Build()` は `query.SearchParams` を返します。

この値は、そのまま高レベル helper に渡せます。

package-level の `SearchParams` に手でコピーする必要はありません。

例:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
if err != nil {
    return err
}
```

これは次の両方で使えます。

- `esv8`
- `esv9`

## package-level `SearchParams` と builder `query.SearchParams`

よく使う search parameter type は 2 種類あります。

- `esv8.SearchParams` / `esv9.SearchParams`
- `esv8/query.SearchParams` / `esv9/query.SearchParams`

どちらも typed Elasticsearch search request に変換できる限り、高レベル helper で利用できます。

### builder params を使う場面

次のような場合は builder params を優先します。

- fluent に query を組み立てたい
- type-safe な query helper を使いたい
- query、sort、aggregation、pagination を 1 箇所にまとめたい

例:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Sort(
        query.NewSort().
            ScoreDesc().
            Build()...,
    ).
    Limit(20).
    Offset(0).
    Build()
```

### package-level `SearchParams` を使う場面

次のような場合に向いています。

- request の各要素をすでに別々に持っている
- fluent builder を使わない
- request field を直接組み立てたい

例:

```go
params := esv8.SearchParams{
    Query: types.Query{
        MatchAll: &types.MatchAllQuery{},
    },
    Size: 10,
    From: 0,
}
```

## 検索レスポンスの構造

### `SearchHit[T]`

各 hit には次が含まれます。

- `ID`
- `Index`
- `Score`
- `Source`
- `Raw`

typed document を扱いつつ、低レベル metadata にもアクセスできます。

### `SearchResponse[T]`

高レベル response には次が含まれます。

- `Total`
- `Hits`
- `Aggregations`
- `Raw`

typed response を得ながら escape hatch も保持したいアプリケーションコードに適しています。

## Aggregation を扱う

typed document と一緒に aggregation も扱いたいなら `Search[T](...)` を使います。

例:

```go
avgPriceAgg := query.AvgAgg("avg_price", ProductFields.Price)
byCategoryAgg := query.StringTermsAgg(
    "by_category",
    ProductFields.Category,
    query.WithSubAggs(avgPriceAgg),
)

params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Aggregation(query.Aggs(byCategoryAgg).Build()).
    Limit(10).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
if err != nil {
    return err
}

terms := resp.Aggregations.MustStringTerms(byCategoryAgg)
for _, bucket := range terms.Buckets() {
    avg := bucket.Aggregations().MustAvg(avgPriceAgg)
    fmt.Println(bucket.Key(), avg.Value())
}
```

aggregation が不要なら、`SearchDocuments[T](...)` の方が適していることが多いです。

## ページネーション

search params を組み立てるときは `Limit(...)` と `Offset(...)` を使います。

例:

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(20).
    Offset(40).
    Build()
```

ここで重要なのは:

- `resp.Total` は一致した総件数
- `len(resp.Hits)` は現在のページで返ってきた件数

たとえば:

- `Total = 125`
- `len(resp.Hits) = 20`

なら、

- 125 件一致した
- そのうち 20 件がこのページに含まれている

という意味です。

## v8 / v9 parity

検索 helper の設計は、次の間で揃える方針です。

- `esv8`
- `esv9`

一方を理解すれば、もう一方もほぼ同じ感覚で扱えるはずです。

典型的な形:

```go
// v8
v8Resp, err := esv8.Search[Product](ctx, v8Client, alias, params)

// v9
v9Resp, err := esv9.Search[Product](ctx, v9Client, alias, params)
```

document-only helper も同様です。

```go
// v8
v8Docs, err := esv8.SearchDocuments[Product](ctx, v8Client, alias, params)

// v9
v9Docs, err := esv9.SearchDocuments[Product](ctx, v9Client, alias, params)
```

## 推奨パターン

### 通常のアプリケーション検索

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
```

### document だけ取得したい

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
```

### 1 件だけ取得したい

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(1).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
```

### 件数だけ知りたい

```go
res, err := client.Count(ctx, alias)
```

### 高度な request shape が必要

```go
req := search.NewRequest()
// advanced Elasticsearch request details を直接設定する

rawResp, err := client.SearchRaw(ctx, alias, req)
```

## 避けたいパターン

### 理由なく builder params を手で展開する

特別な変換が必要な場合を除き、次のような手動展開は避けます。

```go
resp, err := esv8.Search[Product](ctx, client, alias, esv8.SearchParams{
    Query:        params.Query,
    Sort:         params.Sort,
    Aggregations: params.Aggregations,
    Highlight:    params.Highlight,
    Collapse:     params.Collapse,
    ScriptFields: params.ScriptFields,
    Size:         params.Size,
    From:         params.From,
})
```

通常は次で十分です。

```go
resp, err := esv8.Search[Product](ctx, client, alias, params)
```

### 件数目的なのに search を使う

total だけ見たいのに search helper を使うより、`Count(...)` を優先してください。

## 関連ドキュメント

- [../README.md](../README.md) — 簡潔なトップレベル概要
- [../examples/quickstart/README.md](../examples/quickstart/README.md) — 実行可能な quickstart 解説
- [property-reference.md](property-reference.md) — property builder のリファレンス
- [contributing.md](contributing.md) — 開発フローとリポジトリ規約