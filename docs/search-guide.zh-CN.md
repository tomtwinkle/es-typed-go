# 搜索指南

## 应使用哪个 API？

| 使用场景 | API |
|---|---|
| 需要 typed hits、metadata、aggregations | `Search[T](ctx, client, alias, params)` |
| 只需要解码后的文档 `_source` | `SearchDocuments[T](ctx, client, alias, params)` |
| 只需要匹配文档数量 | `Count(ctx, alias)` |
| 高层 helper 无法覆盖的请求格式 | `SearchRaw(ctx, alias, req)` |

## 构建搜索参数

使用顶层 `query` 包构建所有查询。该包与版本无关，`esv8` 和 `esv9` 均可使用。

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

将 `params` 直接传给 `Search[T]`：

```go
// v8
resp, err := esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)

// v9 — 完全相同，只有 package 不同
resp, err := esv9.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

## 在 v8 与 v9 之间切换

由于查询构建使用共享的 `query/` 包，切换 Elasticsearch 版本只需修改客户端 import。完整变更列表请参阅 [migration-v2.md](migration-v2.md)。

## 获取 Aggregation 结果

`SearchResponse.Aggregations` 类型为 `query.AggResults`。使用 `GetXxx` / `MustXxx` 方法获取类型化结果：

```go
avgDef := query.AvgAgg("avg_price", esmodel.Product.Fields.Price)
termsDef := query.StringTermsAgg("by_category", esmodel.Product.Fields.Category,
    query.WithSubAggs(avgDef))

// ... 执行搜索 ...

terms := resp.Aggregations.MustStringTerms(termsDef)
for _, bucket := range terms.Buckets() {
    avg, _ := bucket.Aggregations().GetAvg(avgDef)
    // avg.Value() 为 *float64
}
```

## 排序方向

使用 `query.SortAsc` / `query.SortDesc`，无需导入版本特定的 `sortorder` 包：

```go
query.NewSort().Field(esmodel.Product.Fields.Date, query.SortDesc)
```

## 查询 helper 参考

### DateRangeQuery

`DateRangeQuery` 接受 functional options，可以任意组合四种比较运算符：

```go
// Gte + Lte（闭区间）
query.DateRangeQuery(field, query.DateRangeGte("2024-01-01"), query.DateRangeLte("2024-12-31"))

// Gt + Lt（开区间）
query.DateRangeQuery(field, query.DateRangeGt("2024-01-01"), query.DateRangeLt("2025-01-01"))

// 单侧边界
query.DateRangeQuery(field, query.DateRangeGte("2024-01-01"))
```

可用选项：`DateRangeGt`、`DateRangeGte`、`DateRangeLt`、`DateRangeLte`。

### MultiTermsAgg 的字段级 Missing

使用 `query.MultiTermLookup` 可对每个字段单独配置。`Missing` 用于为没有该字段的文档提供替代值：

```go
query.MultiTermsAgg("by_date_tz", []query.MultiTermLookup{
    {Field: esmodel.Item.Fields.BusinessDate},
    {Field: esmodel.Item.Fields.Timezone, Missing: "UTC"},
})
```

### Field.Ptr() — typed field 转 *string

当 raw go-elasticsearch 类型需要 `*string` 时（如 `NestedAggregation.Path`、`SumAggregation.Field`），可使用 `Ptr()` 代替临时变量：

```go
// Before
path := string(esmodel.Item.Fields.Items)
types.NestedAggregation{Path: &path}

// After
types.NestedAggregation{Path: esmodel.Item.Fields.Items.Ptr()}
```

`Ptr()` 同样适用于 `estype.Alias` 和 `estype.Index`。

## 补充说明

- `Limit(0)` 不返回 hits（适用于仅需 aggregation 或 count 的搜索）。
- 获取单条结果使用 `Limit(1)`。
- `SearchRaw` 是高层 helper 无法满足时的逃生通道，接受任意 `*search.Request`。
- `esv8` 与 `esv9` 的搜索 helper 用法完全一致。

## 相关文档

- [../README.zh-CN.md](../README.zh-CN.md)
- [migration-v2.md](migration-v2.md) — v2 架构与迁移步骤
- [property-reference.md](property-reference.md)
- [contributing.md](contributing.md)
