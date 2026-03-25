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
- builder が返す search params をそのまま高レベル検索 helper に渡せる設計

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

実行可能なエンドツーエンド例は以下を参照してください。

- [`examples/quickstart/main.go`](examples/quickstart/main.go)
- [`examples/quickstart/README.md`](examples/quickstart/README.md)

## ドキュメント

詳細な説明は `docs/` 以下にあります。

### 利用者向け
- [Search Guide](docs/search-guide.md)
- [Property Reference](docs/property-reference.md)
- [ドキュメント一覧](docs/README.md)

### コントリビューター向け
- [Contributing Guide](docs/contributing.md)

## バージョンサポート

- `esv8` は Elasticsearch v8 向け
- `esv9` は Elasticsearch v9 向け

両パッケージは、できるだけ近い API 形状を維持する方針です。

## ライセンス

[MIT](LICENSE)