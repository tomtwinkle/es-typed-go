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
- **只需更改 import 路径即可在 Elasticsearch v8 和 v9 之间切换**
- 两个版本共享的 version-agnostic `query/` 包

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
	"github.com/tomtwinkle/es-typed-go/query"
	"github.com/tomtwinkle/es-typed-go/examples/quickstart/esmodel"
)

type Product struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func main() {
	ctx := context.Background()

	client, err := esv8.NewClient(...)
	if err != nil {
		panic(err)
	}

	// esmodel.Product.Alias / .Index 来自生成的模型文件。
	// esmodel.Product.Fields.Status 是类型化字段名 (estype.Field)。
	params := query.NewSearch().
		Where(query.TermValue(esmodel.Product.Fields.Status, "active")).
		Limit(10).
		Build()

	_, _ = esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)
}
```

## 从 v8 切换到 v9

由于查询构建使用共享的 `query/` 包，切换 Elasticsearch 版本只需**更改客户端的 import**：

```go
// v8
import (
    es8 "github.com/elastic/go-elasticsearch/v8"
    "github.com/tomtwinkle/es-typed-go/esv8"
    "github.com/tomtwinkle/es-typed-go/query"
)
client, _ := esv8.NewClient(es8.Config{Addresses: []string{"http://localhost:19200"}})
resp, err := esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)

// v9 — 只更改上面两行
import (
    es9 "github.com/elastic/go-elasticsearch/v9"
    "github.com/tomtwinkle/es-typed-go/esv9"
    "github.com/tomtwinkle/es-typed-go/query"
)
client, _ := esv9.NewClient(es9.Config{Addresses: []string{"http://localhost:19201"}})
resp, err := esv9.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

查询构建、字段名、聚合、排序定义均无需更改。

## 生成模型格式

使用 `-struct` 和 `-group` 标志运行 `estyped`，可生成一个统一的模型访问器，将类型化字段名、别名和索引名集中在一个变量中：

```go
// esdefinition/product.go
func (Product) Alias() estype.Alias { return "product" }
func (Product) Index() estype.Index { return "product-000001" }

//go:generate go tool estyped -struct Product -package esmodel -out ../esmodel/product_gen.go -group Product
```

生成的访问器：

```go
// esmodel/product_gen.go（生成文件 — 请勿直接编辑）
var Product = struct {
    Fields struct {
        Status   estype.Field
        Category estype.Field
        // ...
    }
    Alias estype.Alias
    Index estype.Index
}{
    Fields: struct{ ... }{Status: "status", Category: "category", ...},
    Alias: "product",
    Index: "product-000001",
}
```

使用方式：

```go
query.TermValue(esmodel.Product.Fields.Status, "active")
esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

可运行的端到端示例请参见：

- [`examples/quickstart/main.go`](examples/quickstart/main.go) — Elasticsearch v8
- [`examples/quickstart_v9/main.go`](examples/quickstart_v9/main.go) — Elasticsearch v9

## 文档

更详细的说明已移到 `docs/` 目录。

### 面向使用者
- [Search Guide](docs/search-guide.md)
- [Property Reference](docs/property-reference.md)
- [迁移指南 (v2)](docs/migration-v2.md)
- [文档索引](docs/README.md)

### 面向贡献者
- [Contributing Guide](docs/contributing.md)

## 版本支持

- `esv8` 面向 Elasticsearch v8
- `esv9` 面向 Elasticsearch v9

两个包共享顶层 `query/` 包进行查询构建，并保持相同的 API 签名。唯一的区别是底层使用的 Elasticsearch 客户端版本不同。

## 许可证

[MIT](LICENSE)
