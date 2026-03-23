# Search Guide

本文档说明 `es-typed-go` 提供的搜索 API、各自适用的场景，以及它们与底层 Elasticsearch 请求类型之间的关系。

顶层 `README*` 文件应保持简洁；搜索相关的详细说明统一放在这里。

## 概览

`es-typed-go` 主要提供三种搜索入口：

- `Search[T](...)` — 高层 typed search，返回命中元信息、总命中数、聚合结果以及原始响应
- `SearchDocuments[T](...)` — 高层 typed search，只返回解码后的文档 `_source`
- `SearchRaw(...)` — 面向高级 Elasticsearch 请求形状的低层 escape hatch

在正常的应用代码中：

- 需要命中元信息或聚合时，优先使用 `Search[T](...)`
- 只需要解码后的 `_source` 时，优先使用 `SearchDocuments[T](...)`
- 只需要匹配文档数量时，优先使用 `Count(...)`
- 只有在高层 helper 尚未覆盖的 Elasticsearch 能力上，才使用 `SearchRaw(...)`

## 搜索 helper 总览

### `Search[T](...)`

当你需要以下能力时，使用 `Search[T](...)`：

- 解码后的文档
- 总命中数
- `_id`、`_index`、`_score` 等命中元信息
- 类型化聚合访问
- 原始 Elasticsearch 响应

示例：

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

如果你只想拿到解码后的 `_source`，而不关心命中元信息或聚合结果，就使用 `SearchDocuments[T](...)`。

示例：

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

当你需要高层 helper 没有建模的高级请求形状时，使用 `SearchRaw(...)`。

典型场景包括：

- `search_after`
- point-in-time
- 自定义 `_source` 过滤
- 需要直接访问 typed client request 的场景
- 库中尚未封装的 Elasticsearch 搜索能力

示例：

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

## 我该使用哪个 API？

### 我需要 typed hits 和元信息

使用 `Search[T](...)`。

这是应用搜索场景下最推荐的默认选择。

### 我只需要解码后的文档

使用 `SearchDocuments[T](...)`。

它本质上是 `Search[T](...)` 的文档提取版 convenience helper。

### 我只想拿一条文档

使用 `Search[T](...)` 或 `SearchDocuments[T](...)`，并设置 `Limit(1)` 或 `Size: 1`。

示例：

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

如果你只需要解码后的文档：

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

## 不要为了 count 而滥用 search

如果你只需要匹配文档数量，优先使用 `Count(...)`。

示例：

```go
res, err := client.Count(ctx, alias)
if err != nil {
    return err
}

fmt.Println(res.Count)
```

为什么更推荐这样做：

- 意图更清晰
- API 更简单
- 不会让 Elasticsearch 返回你根本不需要的 hits

### 那 `Limit(0)` 呢？

`Limit(0)` 的含义是“不返回 hits”。

它仍然有用，例如你想要：

- 只做 aggregations
- 从搜索请求里拿 total hits
- 发起一个 search 形状的请求，但不需要文档内容

但是如果你的目标只是“有多少条文档匹配”，通常 `Count(...)` 更合适。

## 直接传递 query builder 生成的 params

`query.NewSearch().Build()` 返回的是 `query.SearchParams`。

你可以把它直接传给高层 helper。

不需要手动把各个字段再拷贝到 package-level 的 `SearchParams` 里。

示例：

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

这在以下两个包中都成立：

- `esv8`
- `esv9`

## package-level `SearchParams` 与 builder `query.SearchParams`

常见有两类搜索参数类型：

- `esv8.SearchParams` / `esv9.SearchParams`
- `esv8/query.SearchParams` / `esv9/query.SearchParams`

只要它们能把自己转换成 typed Elasticsearch search request，就都可以用于高层 helper。

### 什么时候用 builder params

以下场景优先使用 builder params：

- 你希望 fluently 组织查询逻辑
- 你希望使用类型安全的 query helper
- 你希望在一个地方组合 query、sort、aggregation、pagination

示例：

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

### 什么时候用 package-level `SearchParams`

以下场景适合使用 package-level `SearchParams`：

- 你已经单独拥有各个请求片段
- 你不打算使用 fluent builder
- 你希望直接拼装 request 的各个字段

示例：

```go
params := esv8.SearchParams{
    Query: types.Query{
        MatchAll: &types.MatchAllQuery{},
    },
    Size: 10,
    From: 0,
}
```

## 搜索响应结构

### `SearchHit[T]`

每个命中包含：

- `ID`
- `Index`
- `Score`
- `Source`
- `Raw`

这样你既可以直接处理 typed 文档，也可以保留底层命中元信息。

### `SearchResponse[T]`

高层响应包含：

- `Total`
- `Hits`
- `Aggregations`
- `Raw`

这使它非常适合应用代码：默认是 typed 的，但依然保留 escape hatch。

## 聚合

如果你需要 typed 文档和聚合结果一起返回，请使用 `Search[T](...)`。

示例：

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

如果你根本不需要聚合，通常 `SearchDocuments[T](...)` 更合适。

## 分页

使用 `Limit(...)` 和 `Offset(...)` 来构建分页参数。

示例：

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(20).
    Offset(40).
    Build()
```

要注意这两个值的区别：

- `resp.Total` 表示总匹配文档数
- `len(resp.Hits)` 只表示当前页返回的命中数

例如：

- `Total = 125`
- `len(resp.Hits) = 20`

表示：

- 一共有 125 条文档匹配
- 当前这一页只返回了 20 条

## v8 与 v9 的一致性

搜索 helper 的设计目标是在以下两个包之间保持对齐：

- `esv8`
- `esv9`

一般来说，学会一个版本后，另一个版本的用法应该非常接近。

典型用法：

```go
// v8
v8Resp, err := esv8.Search[Product](ctx, v8Client, alias, params)

// v9
v9Resp, err := esv9.Search[Product](ctx, v9Client, alias, params)
```

文档数组场景也是一样：

```go
// v8
v8Docs, err := esv8.SearchDocuments[Product](ctx, v8Client, alias, params)

// v9
v9Docs, err := esv9.SearchDocuments[Product](ctx, v9Client, alias, params)
```

## 推荐模式

### 常规应用搜索

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

resp, err := esv8.Search[Product](ctx, client, alias, params)
```

### 只取文档

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(10).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
```

### 只取一条

```go
params := query.NewSearch().
    Where(query.TermValue(ProductFields.Status, "active")).
    Limit(1).
    Build()

docs, err := esv8.SearchDocuments[Product](ctx, client, alias, params)
```

### 只取 count

```go
res, err := client.Count(ctx, alias)
```

### 高级请求形状

```go
req := search.NewRequest()
// 直接配置高级 Elasticsearch request 字段

rawResp, err := client.SearchRaw(ctx, alias, req)
```

## 反模式

### 无理由地手动展开 builder params

除非你确实需要转换字段，否则不要这么写：

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

更推荐这样写：

```go
resp, err := esv8.Search[Product](ctx, client, alias, params)
```

### 实际只想 count，却去调用 search

如果你最终只是在看总命中数，而并不需要 hits 或 aggregations，就不要用 search helper，优先使用 `Count(...)`。

## 相关文档

- [../README.md](../README.md) — 简洁的顶层概览
- [../examples/quickstart/README.md](../examples/quickstart/README.md) — 可运行 quickstart 说明
- [property-reference.md](property-reference.md) — property builder 参考
- [contributing.md](contributing.md) — 贡献流程与仓库规则