# es-typed-go

[![Test](https://github.com/tomtwinkle/es-typed-go/actions/workflows/test.yml/badge.svg)](https://github.com/tomtwinkle/es-typed-go/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tomtwinkle/es-typed-go.svg)](https://pkg.go.dev/github.com/tomtwinkle/es-typed-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](README.md) | [日本語](README.ja.md) | **中文**

一个面向 [go-elasticsearch](https://github.com/elastic/go-elasticsearch)（v8 / v9）的类型安全 Go 封装库，在编译期防止字段名拼写错误以及索引 / 别名混淆。

## 这个库解决什么问题？

官方 Elasticsearch Go typed client 很强大，但在日常应用代码中通常会遇到这些问题：

- 请求结构过于宽泛，使用门槛较高
- 字段名、索引名、别名名容易混用
- 查询、排序、聚合和 mapping 定义写起来比较繁琐
- 很多错误只能在运行时才发现

`es-typed-go` 通过以下方式改善这些体验：

- 为字段、索引、别名提供不同的类型
- 从 mapping 或 Go struct 生成类型化字段访问器
- 提供 query / sort / aggregation 的 fluent builder
- 提供 Elasticsearch property 的 functional-option builder
- 提供将 `_source` 直接解码为 Go struct 的高层搜索 helper

## 主要特性

- `estype.Field` / `estype.Index` / `estype.Alias` 的类型分离
- 使用 `estyped` 进行代码生成
- 类型安全的 query / sort / aggregation builder
- Elasticsearch property builder
- 同时支持 `esv8` 和 `esv9`
- builder 返回的 search params 可直接传给高层 typed search helper

## 安装

安装库本体：

```bash
go get github.com/tomtwinkle/es-typed-go
```

将代码生成 CLI 作为 Go tool 安装：

```bash
go get -tool github.com/tomtwinkle/es-typed-go/cmd/estyped
```

运行：

```bash
go tool estyped
```

## 最小示例

```go
package main

import (
	"context"

	"github.com/tomtwinkle/es-typed-go/esv8"
	"github.com/tomtwinkle/es-typed-go/esv8/query"
	"github.com/tomtwinkle/es-typed-go/estype"
)

type Product struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

var FieldStatus estype.Field = "status"

func main() {
	ctx := context.Background()

	client, err := esv8.NewClient(...)
	if err != nil {
		panic(err)
	}

	alias := estype.Alias("products")

	params := query.NewSearch().
		Where(query.TermValue(FieldStatus, "active")).
		Limit(10).
		Build()

	_, _ = esv8.Search[Product](ctx, client, alias, params)
}
```

可运行的端到端示例请参见：

- [`examples/quickstart/main.go`](examples/quickstart/main.go)
- [`examples/quickstart/README.md`](examples/quickstart/README.md)

## 文档

更详细的说明已移到 `docs/` 目录。

### 面向使用者
- [Search Guide](docs/search-guide.md)
- [Property Reference](docs/property-reference.md)
- [文档索引](docs/README.md)

### 面向贡献者
- [Contributing Guide](docs/contributing.md)

## 版本支持

- `esv8` 面向 Elasticsearch v8
- `esv9` 面向 Elasticsearch v9

两个包会尽量保持接近的 API 形状。

## 许可证

[MIT](LICENSE)