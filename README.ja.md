# es-typed-go

[![Test](https://github.com/tomtwinkle/es-typed-go/actions/workflows/test.yml/badge.svg)](https://github.com/tomtwinkle/es-typed-go/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tomtwinkle/es-typed-go.svg)](https://pkg.go.dev/github.com/tomtwinkle/es-typed-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[English](README.md) | **日本語** | [中文](README.zh-CN.md)

[go-elasticsearch](https://github.com/elastic/go-elasticsearch) (v8 / v9) 向けの型安全な Go ラッパーです。フィールド名のタイポやインデックス / エイリアスの取り違えをコンパイル時に防止します。

## このライブラリの目的

Elasticsearch 公式の Go typed client は強力ですが、通常のアプリケーションコードでは次のような課題があります。

- request shape が広く、使い方が分かりにくい
- 文字列ベースの指定ミスをコンパイル時に防ぎにくい
- 検索クエリや mapping 定義の記述が冗長になりやすい

`es-typed-go` はその改善のために、次を提供します。

- `Field` / `Index` / `Alias` の distinct type
- mapping や Go struct からのフィールド定数生成
- query / sort / aggregation の fluent builder
- Elasticsearch property 定義用の functional-option builder
- `_source` を Go struct にデコードする高レベル検索 helper

## 主な機能

- `estype.Field` / `estype.Index` / `estype.Alias` の型分離
- `estyped` によるコード生成
- 型安全な query / sort / aggregation builder
- Elasticsearch property builder
- `esv8` / `esv9` 両対応
- **import パスの変更のみで Elasticsearch v8 / v9 を切り替え可能**
- 両バージョンで共有される version-agnostic な `query/` パッケージ

## インストール

ライブラリ本体をインストール:

```bash
go get github.com/tomtwinkle/es-typed-go
```

コード生成 CLI を Go tool として追加:

```bash
go get -tool github.com/tomtwinkle/es-typed-go/cmd/estyped
```

実行:

```bash
go tool estyped
```

## 最小例

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

	// esmodel.Product.Alias / .Index は生成ファイルから取得。
	// esmodel.Product.Fields.Status は型付きフィールド名 (estype.Field)。
	params := query.NewSearch().
		Where(query.TermValue(esmodel.Product.Fields.Status, "active")).
		Limit(10).
		Build()

	_, _ = esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)
}
```

## v8 から v9 への切り替え

クエリ構築には共通の `query/` パッケージを使うため、Elasticsearch のバージョン切り替えは **クライアントの import 変更のみ** で完了します:

```go
// v8
import (
    es8 "github.com/elastic/go-elasticsearch/v8"
    "github.com/tomtwinkle/es-typed-go/esv8"
    "github.com/tomtwinkle/es-typed-go/query"
)
client, _ := esv8.NewClient(es8.Config{Addresses: []string{"http://localhost:19200"}})
resp, err := esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)

// v9 — 上の2行だけ変更する
import (
    es9 "github.com/elastic/go-elasticsearch/v9"
    "github.com/tomtwinkle/es-typed-go/esv9"
    "github.com/tomtwinkle/es-typed-go/query"
)
client, _ := esv9.NewClient(es9.Config{Addresses: []string{"http://localhost:19201"}})
resp, err := esv9.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

クエリ構築・フィールド名・アグリゲーション・ソート定義はすべて変更不要です。

## 生成モデル形式

`estyped` を `-struct` と `-group` フラグで実行すると、型付きフィールド名・エイリアス・インデックス名を1つの変数にまとめた統合アクセサを生成できます:

```go
// esdefinition/product.go
func (Product) Alias() estype.Alias { return "product" }
func (Product) Index() estype.Index { return "product-000001" }

//go:generate go tool estyped -struct Product -package esmodel -out ../esmodel/product_gen.go -group Product
```

生成されるアクセサ:

```go
// esmodel/product_gen.go (生成ファイル — 直接編集不可)
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

使用例:

```go
query.TermValue(esmodel.Product.Fields.Status, "active")
esv8.Search[Product](ctx, client, esmodel.Product.Alias, params)
```

実行可能なエンドツーエンド例は以下を参照してください:

- [`examples/quickstart/main.go`](examples/quickstart/main.go) — Elasticsearch v8
- [`examples/quickstart_v9/main.go`](examples/quickstart_v9/main.go) — Elasticsearch v9

## ドキュメント

詳細な説明は `docs/` 以下にあります。

### 利用者向け
- [Search Guide](docs/search-guide.md)
- [Property Reference](docs/property-reference.md)
- [マイグレーションガイド (v2)](docs/migration-v2.md)
- [ドキュメント一覧](docs/README.md)

### コントリビューター向け
- [Contributing Guide](docs/contributing.md)

## バージョンサポート

- `esv8` は Elasticsearch v8 向け
- `esv9` は Elasticsearch v9 向け

両パッケージはトップレベルの `query/` パッケージを共有し、同一の API シグネチャを維持します。違いは内部で使用する Elasticsearch クライアントのバージョンのみです。

## ライセンス

[MIT](LICENSE)
