# プロパティリファレンス

[English](property-reference.md) | [日本語](property-reference.ja.md)

このドキュメントでは、`esv8` / `esv9` パッケージが提供するすべてのプロパティコンストラクタとその Functional Option を一覧にまとめています。すべてのコンストラクタは Functional-option パターンに従います。

```go
prop := esv8.NewTextProperty(
    esv8.WithTextAnalyzer("standard"),
    esv8.WithTextStore(true),
)
```

特に記載がない限り、ここに掲載するすべてのプロパティは `esv8` と `esv9` の両方で同一のシグネチャで利用できます。

---

## 目次

- [テキストファミリー](#テキストファミリー)
- [数値](#数値)
- [日付とブーリアン](#日付とブーリアン)
- [地理](#地理)
- [レンジ](#レンジ)
- [オブジェクトとネスト](#オブジェクトとネスト)
- [Join](#join)
- [ネットワーク](#ネットワーク)
- [ベクトル](#ベクトル)
- [ランキング](#ランキング)
- [特殊](#特殊)
- [プラグイン依存](#プラグイン依存)

---

## テキストファミリー

### NewTextProperty

アナライザをサポートする全文検索フィールド。

```go
esv8.NewTextProperty(opts ...TextPropertyOption) *types.TextProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithTextAnalyzer` | `v string` | インデクシング用アナライザを設定 |
| `WithTextSearchAnalyzer` | `v string` | 検索クエリ用アナライザを設定 |
| `WithTextSearchQuoteAnalyzer` | `v string` | フレーズクエリ用アナライザを設定 |
| `WithTextFielddata` | `v bool` | テキストのソート/アグリゲーション用 fielddata を有効化 |
| `WithTextIndex` | `v bool` | フィールドを検索可能にするか |
| `WithTextStore` | `v bool` | フィールド値を個別に保存するか |
| `WithTextNorms` | `v bool` | スコアリング用の norms を保存するか |
| `WithTextSimilarity` | `v string` | 類似度アルゴリズム（例：`"BM25"`） |
| `WithTextIndexPhrases` | `v bool` | 2 語の組み合わせをインデックスするか |
| `WithTextPositionIncrementGap` | `v int` | 配列値間の疑似トークン位置数 |
| `WithTextFields` | `fields map[string]types.Property` | カスタムマルチフィールドを設定 |

### NewKeywordProperty

フィルタリング、ソート、アグリゲーション用の完全一致文字列フィールド。

```go
esv8.NewKeywordProperty(opts ...KeywordPropertyOption) *types.KeywordProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithKeywordIgnoreAbove` | `v int` | 最大文字列長。超過する値はインデックスされない |
| `WithKeywordDocValues` | `v bool` | doc values を有効にするか |
| `WithKeywordIndex` | `v bool` | フィールドを検索可能にするか |
| `WithKeywordStore` | `v bool` | フィールド値を個別に保存するか |
| `WithKeywordNullValue` | `v string` | `null` の代替値 |
| `WithKeywordNormalizer` | `v string` | インデクシング前に適用するノーマライザ |
| `WithKeywordNorms` | `v bool` | スコアリング用の norms を保存するか |
| `WithKeywordSimilarity` | `v string` | 類似度アルゴリズム |
| `WithKeywordEagerGlobalOrdinals` | `v bool` | グローバルオーディナルを事前ロードするか |
| `WithKeywordSplitQueriesOnWhitespace` | `v bool` | クエリを空白で分割するか |

### NewConstantKeywordProperty

全ドキュメントが同一の値を持つキーワードフィールド。

```go
esv8.NewConstantKeywordProperty(opts ...ConstantKeywordPropertyOption) *types.ConstantKeywordProperty
```

オプションなし。値は最初にインデックスされたドキュメントで設定されます。

### NewCountedKeywordProperty

各タームのカウントも追跡するキーワードフィールド。

```go
esv8.NewCountedKeywordProperty(opts ...CountedKeywordPropertyOption) *types.CountedKeywordProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithCountedKeywordIndex` | `v bool` | フィールドを検索可能にするか |

### NewWildcardProperty

ワイルドカードクエリと正規表現クエリに最適化されたキーワードフィールド。

```go
esv8.NewWildcardProperty(opts ...WildcardPropertyOption) *types.WildcardProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithWildcardIgnoreAbove` | `v int` | 最大文字列長 |
| `WithWildcardNullValue` | `v string` | `null` の代替値 |

### NewMatchOnlyTextProperty

位置情報や norms を保存しない、マッチクエリに最適化されたテキストフィールド。

```go
esv8.NewMatchOnlyTextProperty(opts ...MatchOnlyTextPropertyOption) *types.MatchOnlyTextProperty
```

オプションなし。

### NewCompletionProperty

オートコンプリートサジェストフィールド。

```go
esv8.NewCompletionProperty(opts ...CompletionPropertyOption) *types.CompletionProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithCompletionAnalyzer` | `v string` | インデクシング用アナライザ |
| `WithCompletionSearchAnalyzer` | `v string` | 検索用アナライザ |
| `WithCompletionMaxInputLength` | `v int` | 単一入力の最大長 |
| `WithCompletionPreservePositionIncrements` | `v bool` | 位置インクリメントを保持するか |
| `WithCompletionPreserveSeparators` | `v bool` | セパレータを保持するか |

### NewSearchAsYouTypeProperty

入力中検索（search-as-you-type）向けフィールド。

```go
esv8.NewSearchAsYouTypeProperty(opts ...SearchAsYouTypePropertyOption) *types.SearchAsYouTypeProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithSearchAsYouTypeAnalyzer` | `v string` | インデクシング用アナライザ |
| `WithSearchAsYouTypeSearchAnalyzer` | `v string` | 検索用アナライザ |
| `WithSearchAsYouTypeSearchQuoteAnalyzer` | `v string` | フレーズクエリ用アナライザ |
| `WithSearchAsYouTypeMaxShingleSize` | `v int` | 最大シングルサイズ（2-4） |
| `WithSearchAsYouTypeIndex` | `v bool` | フィールドを検索可能にするか |
| `WithSearchAsYouTypeStore` | `v bool` | フィールド値を個別に保存するか |
| `WithSearchAsYouTypeNorms` | `v bool` | norms を保存するか |
| `WithSearchAsYouTypeSimilarity` | `v string` | 類似度アルゴリズム |

---

## 数値

すべての数値プロパティ型は共通のオプション（coerce、doc_values、ignore_malformed、index、store、null_value）を持ちます。型ごとの例外は個別に記載しています。

### NewIntegerNumberProperty

32 ビット符号付き整数。

```go
esv8.NewIntegerNumberProperty(opts ...IntegerNumberPropertyOption) *types.IntegerNumberProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithIntegerNumberCoerce` | `v bool` | 値を正しい型に変換するか |
| `WithIntegerNumberDocValues` | `v bool` | doc values を有効にするか |
| `WithIntegerNumberIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithIntegerNumberIndex` | `v bool` | フィールドを検索可能にするか |
| `WithIntegerNumberStore` | `v bool` | フィールド値を個別に保存するか |
| `WithIntegerNumberNullValue` | `v int` | `null` の代替値 |

### NewLongNumberProperty

64 ビット符号付き整数。

```go
esv8.NewLongNumberProperty(opts ...LongNumberPropertyOption) *types.LongNumberProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithLongNumberCoerce` | `v bool` | 値を正しい型に変換するか |
| `WithLongNumberDocValues` | `v bool` | doc values を有効にするか |
| `WithLongNumberIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithLongNumberIndex` | `v bool` | フィールドを検索可能にするか |
| `WithLongNumberStore` | `v bool` | フィールド値を個別に保存するか |
| `WithLongNumberNullValue` | `v int64` | `null` の代替値 |

### NewShortNumberProperty

16 ビット符号付き整数。

```go
esv8.NewShortNumberProperty(opts ...ShortNumberPropertyOption) *types.ShortNumberProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithShortNumberCoerce` | `v bool` | 値を正しい型に変換するか |
| `WithShortNumberDocValues` | `v bool` | doc values を有効にするか |
| `WithShortNumberIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithShortNumberIndex` | `v bool` | フィールドを検索可能にするか |
| `WithShortNumberStore` | `v bool` | フィールド値を個別に保存するか |
| `WithShortNumberNullValue` | `v int` | `null` の代替値 |

### NewByteNumberProperty

8 ビット符号付き整数。

```go
esv8.NewByteNumberProperty(opts ...ByteNumberPropertyOption) *types.ByteNumberProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithByteNumberCoerce` | `v bool` | 値を正しい型に変換するか |
| `WithByteNumberDocValues` | `v bool` | doc values を有効にするか |
| `WithByteNumberIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithByteNumberIndex` | `v bool` | フィールドを検索可能にするか |
| `WithByteNumberStore` | `v bool` | フィールド値を個別に保存するか |
| `WithByteNumberNullValue` | `v byte` | `null` の代替値 |

### NewDoubleNumberProperty

64 ビット IEEE 754 浮動小数点数。

```go
esv8.NewDoubleNumberProperty(opts ...DoubleNumberPropertyOption) *types.DoubleNumberProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithDoubleNumberCoerce` | `v bool` | 値を正しい型に変換するか |
| `WithDoubleNumberDocValues` | `v bool` | doc values を有効にするか |
| `WithDoubleNumberIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithDoubleNumberIndex` | `v bool` | フィールドを検索可能にするか |
| `WithDoubleNumberStore` | `v bool` | フィールド値を個別に保存するか |
| `WithDoubleNumberNullValue` | `v float64` | `null` の代替値 |

### NewFloatNumberProperty

32 ビット IEEE 754 浮動小数点数。

```go
esv8.NewFloatNumberProperty(opts ...FloatNumberPropertyOption) *types.FloatNumberProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithFloatNumberCoerce` | `v bool` | 値を正しい型に変換するか |
| `WithFloatNumberDocValues` | `v bool` | doc values を有効にするか |
| `WithFloatNumberIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithFloatNumberIndex` | `v bool` | フィールドを検索可能にするか |
| `WithFloatNumberStore` | `v bool` | フィールド値を個別に保存するか |
| `WithFloatNumberNullValue` | `v float32` | `null` の代替値 |

### NewHalfFloatNumberProperty

16 ビット IEEE 754 浮動小数点数。

```go
esv8.NewHalfFloatNumberProperty(opts ...HalfFloatNumberPropertyOption) *types.HalfFloatNumberProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithHalfFloatNumberCoerce` | `v bool` | 値を正しい型に変換するか |
| `WithHalfFloatNumberDocValues` | `v bool` | doc values を有効にするか |
| `WithHalfFloatNumberIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithHalfFloatNumberIndex` | `v bool` | フィールドを検索可能にするか |
| `WithHalfFloatNumberStore` | `v bool` | フィールド値を個別に保存するか |
| `WithHalfFloatNumberNullValue` | `v float32` | `null` の代替値 |

### NewUnsignedLongNumberProperty

符号なし 64 ビット整数（0 から 2^64-1）。

```go
esv8.NewUnsignedLongNumberProperty(opts ...UnsignedLongNumberPropertyOption) *types.UnsignedLongNumberProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithUnsignedLongNumberDocValues` | `v bool` | doc values を有効にするか |
| `WithUnsignedLongNumberIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithUnsignedLongNumberIndex` | `v bool` | フィールドを検索可能にするか |
| `WithUnsignedLongNumberStore` | `v bool` | フィールド値を個別に保存するか |
| `WithUnsignedLongNumberNullValue` | `v uint64` | `null` の代替値 |

> **注意:** `unsigned_long` は `coerce` パラメータをサポートしていません。

### NewScaledFloatNumberProperty

コンパクトな保存のためにスケーリングされた `long` として格納される浮動小数点数。

```go
esv8.NewScaledFloatNumberProperty(opts ...ScaledFloatNumberPropertyOption) *types.ScaledFloatNumberProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithScaledFloatNumberScalingFactor` | `v float64` | スケーリングファクター（必須） |
| `WithScaledFloatNumberCoerce` | `v bool` | 値を正しい型に変換するか |
| `WithScaledFloatNumberDocValues` | `v bool` | doc values を有効にするか |
| `WithScaledFloatNumberIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithScaledFloatNumberIndex` | `v bool` | フィールドを検索可能にするか |
| `WithScaledFloatNumberStore` | `v bool` | フィールド値を個別に保存するか |
| `WithScaledFloatNumberNullValue` | `v float64` | `null` の代替値 |

---

## 日付とブーリアン

### NewDateProperty

日付/時刻フィールド。

```go
esv8.NewDateProperty(opts ...DatePropertyOption) *types.DateProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithDateFormat` | `formats ...estype.DateFormat` | 受け入れる日付フォーマット（`\|\|` で結合） |
| `WithDateDocValues` | `v bool` | doc values を有効にするか |
| `WithDateIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithDateIndex` | `v bool` | フィールドを検索可能にするか |
| `WithDateStore` | `v bool` | フィールド値を個別に保存するか |
| `WithDateLocale` | `v string` | 日付解析のロケール（例：`"en"`） |

### NewDateNanosProperty

ナノ秒精度の日付/時刻フィールド。

```go
esv8.NewDateNanosProperty(opts ...DateNanosPropertyOption) *types.DateNanosProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithDateNanosFormat` | `formats ...estype.DateFormat` | 受け入れる日付フォーマット |
| `WithDateNanosDocValues` | `v bool` | doc values を有効にするか |
| `WithDateNanosIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithDateNanosIndex` | `v bool` | フィールドを検索可能にするか |
| `WithDateNanosStore` | `v bool` | フィールド値を個別に保存するか |

### NewBooleanProperty

ブーリアン（`true`/`false`）フィールド。

```go
esv8.NewBooleanProperty(opts ...BooleanPropertyOption) *types.BooleanProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithBooleanDocValues` | `v bool` | doc values を有効にするか |
| `WithBooleanIndex` | `v bool` | フィールドを検索可能にするか |
| `WithBooleanStore` | `v bool` | フィールド値を個別に保存するか |
| `WithBooleanNullValue` | `v bool` | `null` の代替値 |

---

## 地理

### NewGeoPointProperty

緯度/経度ポイント。

```go
esv8.NewGeoPointProperty(opts ...GeoPointPropertyOption) *types.GeoPointProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithGeoPointIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithGeoPointIgnoreZValue` | `v bool` | Z 値を無視するか |
| `WithGeoPointDocValues` | `v bool` | doc values を有効にするか |
| `WithGeoPointIndex` | `v bool` | フィールドを検索可能にするか |
| `WithGeoPointStore` | `v bool` | フィールド値を個別に保存するか |

### NewGeoShapeProperty

任意の GeoJSON ジオメトリ。

```go
esv8.NewGeoShapeProperty(opts ...GeoShapePropertyOption) *types.GeoShapeProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithGeoShapeCoerce` | `v bool` | 閉じていないポリゴンを補正するか |
| `WithGeoShapeIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithGeoShapeIgnoreZValue` | `v bool` | Z 値を無視するか |
| `WithGeoShapeDocValues` | `v bool` | doc values を有効にするか |
| `WithGeoShapeIndex` | `v bool` | フィールドを検索可能にするか |
| `WithGeoShapeStore` | `v bool` | フィールド値を個別に保存するか |

### NewShapeProperty

任意のデカルトジオメトリ（非地理）。

```go
esv8.NewShapeProperty(opts ...ShapePropertyOption) *types.ShapeProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithShapeCoerce` | `v bool` | 閉じていないポリゴンを補正するか |
| `WithShapeIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithShapeIgnoreZValue` | `v bool` | Z 値を無視するか |
| `WithShapeDocValues` | `v bool` | doc values を有効にするか |

### NewPointProperty

デカルト（x, y）ポイント。

```go
esv8.NewPointProperty(opts ...PointPropertyOption) *types.PointProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithPointIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithPointIgnoreZValue` | `v bool` | Z 値を無視するか |
| `WithPointDocValues` | `v bool` | doc values を有効にするか |
| `WithPointStore` | `v bool` | フィールド値を個別に保存するか |
| `WithPointNullValue` | `v string` | `null` の代替値（WKT ポイント） |

---

## レンジ

すべてのレンジプロパティ型は共通のオプション（coerce、doc_values、index、store）を持ちます。date_range 型は追加で format オプションをサポートします。

### NewIntegerRangeProperty

```go
esv8.NewIntegerRangeProperty(opts ...IntegerRangePropertyOption) *types.IntegerRangeProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithIntegerRangeCoerce` | `v bool` | 値を変換するか |
| `WithIntegerRangeDocValues` | `v bool` | doc values を有効にするか |
| `WithIntegerRangeIndex` | `v bool` | フィールドを検索可能にするか |
| `WithIntegerRangeStore` | `v bool` | フィールド値を個別に保存するか |

### NewLongRangeProperty

```go
esv8.NewLongRangeProperty(opts ...LongRangePropertyOption) *types.LongRangeProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithLongRangeCoerce` | `v bool` | 値を変換するか |
| `WithLongRangeDocValues` | `v bool` | doc values を有効にするか |
| `WithLongRangeIndex` | `v bool` | フィールドを検索可能にするか |
| `WithLongRangeStore` | `v bool` | フィールド値を個別に保存するか |

### NewFloatRangeProperty

```go
esv8.NewFloatRangeProperty(opts ...FloatRangePropertyOption) *types.FloatRangeProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithFloatRangeCoerce` | `v bool` | 値を変換するか |
| `WithFloatRangeDocValues` | `v bool` | doc values を有効にするか |
| `WithFloatRangeIndex` | `v bool` | フィールドを検索可能にするか |
| `WithFloatRangeStore` | `v bool` | フィールド値を個別に保存するか |

### NewDoubleRangeProperty

```go
esv8.NewDoubleRangeProperty(opts ...DoubleRangePropertyOption) *types.DoubleRangeProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithDoubleRangeCoerce` | `v bool` | 値を変換するか |
| `WithDoubleRangeDocValues` | `v bool` | doc values を有効にするか |
| `WithDoubleRangeIndex` | `v bool` | フィールドを検索可能にするか |
| `WithDoubleRangeStore` | `v bool` | フィールド値を個別に保存するか |

### NewDateRangeProperty

```go
esv8.NewDateRangeProperty(opts ...DateRangePropertyOption) *types.DateRangeProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithDateRangeFormat` | `formats ...estype.DateFormat` | 受け入れる日付フォーマット |
| `WithDateRangeCoerce` | `v bool` | 値を変換するか |
| `WithDateRangeDocValues` | `v bool` | doc values を有効にするか |
| `WithDateRangeIndex` | `v bool` | フィールドを検索可能にするか |
| `WithDateRangeStore` | `v bool` | フィールド値を個別に保存するか |

### NewIpRangeProperty

```go
esv8.NewIpRangeProperty(opts ...IpRangePropertyOption) *types.IpRangeProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithIpRangeCoerce` | `v bool` | 値を変換するか |
| `WithIpRangeDocValues` | `v bool` | doc values を有効にするか |
| `WithIpRangeIndex` | `v bool` | フィールドを検索可能にするか |
| `WithIpRangeStore` | `v bool` | フィールド値を個別に保存するか |

---

## オブジェクトとネスト

### NewObjectProperty

JSON オブジェクト（フラット構造、子フィールドの独立クエリ不可）。

```go
esv8.NewObjectProperty(opts ...ObjectPropertyOption) *types.ObjectProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithObjectProperties` | `v map[string]types.Property` | 子フィールドマッピング |
| `WithObjectEnabled` | `v bool` | オブジェクトのインデクシングを有効にするか |

### NewNestedProperty

子フィールドを独立してクエリできる JSON オブジェクト。

```go
esv8.NewNestedProperty(opts ...NestedPropertyOption) *types.NestedProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithNestedProperties` | `v map[string]types.Property` | 子フィールドマッピング |
| `WithNestedEnabled` | `v bool` | ネストオブジェクトを有効にするか |
| `WithNestedIncludeInParent` | `v bool` | ネストフィールドを親に含めるか |
| `WithNestedIncludeInRoot` | `v bool` | ネストフィールドをルートに含めるか |

### NewFlattenedProperty

JSON オブジェクト全体を単一フィールドとしてマッピング。

```go
esv8.NewFlattenedProperty(opts ...FlattenedPropertyOption) *types.FlattenedProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithFlattenedDepthLimit` | `v int` | 最大ネスト深度 |
| `WithFlattenedDocValues` | `v bool` | doc values を有効にするか |
| `WithFlattenedIndex` | `v bool` | フィールドを検索可能にするか |
| `WithFlattenedIgnoreAbove` | `v int` | リーフ値の最大文字列長 |
| `WithFlattenedNullValue` | `v string` | `null` の代替値 |
| `WithFlattenedEagerGlobalOrdinals` | `v bool` | グローバルオーディナルを事前ロードするか |
| `WithFlattenedSimilarity` | `v string` | 類似度アルゴリズム |
| `WithFlattenedSplitQueriesOnWhitespace` | `v bool` | クエリを空白で分割するか |

### NewPassthroughObjectProperty

フィールドが親レベルにマッピングされるオブジェクト（データストリーム用）。

```go
esv8.NewPassthroughObjectProperty(opts ...PassthroughObjectPropertyOption) *types.PassthroughObjectProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithPassthroughObjectProperties` | `v map[string]types.Property` | 子フィールドマッピング |
| `WithPassthroughObjectEnabled` | `v bool` | オブジェクトを有効にするか |
| `WithPassthroughObjectPriority` | `v int` | フィールド名競合時の優先度 |
| `WithPassthroughObjectTimeSeriesDimension` | `v bool` | 時系列ディメンションかどうか |

---

## Join

### NewJoinProperty

親子関係フィールド。

```go
esv8.NewJoinProperty(opts ...JoinPropertyOption) *types.JoinProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithJoinRelations` | `v map[string][]string` | 親子関係の定義 |
| `WithJoinEagerGlobalOrdinals` | `v bool` | グローバルオーディナルを事前ロードするか |

---

## ネットワーク

### NewIpProperty

IPv4 または IPv6 アドレス。

```go
esv8.NewIpProperty(opts ...IpPropertyOption) *types.IpProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithIpDocValues` | `v bool` | doc values を有効にするか |
| `WithIpIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithIpIndex` | `v bool` | フィールドを検索可能にするか |
| `WithIpStore` | `v bool` | フィールド値を個別に保存するか |
| `WithIpNullValue` | `v string` | `null` の代替値 |

---

## ベクトル

### NewDenseVectorProperty

kNN 検索用の密な浮動小数点ベクトル。

```go
esv8.NewDenseVectorProperty(opts ...DenseVectorPropertyOption) *types.DenseVectorProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithDenseVectorDims` | `v int` | 次元数 |
| `WithDenseVectorIndex` | `v bool` | kNN 検索用にインデックスするか |

### NewSparseVectorProperty

ターム型ランキング用のスパース浮動小数点ベクトル。

```go
esv8.NewSparseVectorProperty(opts ...SparseVectorPropertyOption) *types.SparseVectorProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithSparseVectorStore` | `v bool` | フィールド値を個別に保存するか |

### NewRankVectorProperty

ランクベーススコアリング用の固定長 float ベクトル。

```go
esv8.NewRankVectorProperty(opts ...RankVectorPropertyOption) *types.RankVectorProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithRankVectorDims` | `v int` | 次元数 |

---

## ランキング

### NewRankFeatureProperty

関連性スコアリングをブーストするための数値特徴量。

```go
esv8.NewRankFeatureProperty(opts ...RankFeaturePropertyOption) *types.RankFeatureProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithRankFeaturePositiveScoreImpact` | `v bool` | 高い値が関連性をブーストするか |

### NewRankFeaturesProperty

関連性スコアリングをブーストするための名前付き数値特徴量のマップ。

```go
esv8.NewRankFeaturesProperty(opts ...RankFeaturesPropertyOption) *types.RankFeaturesProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithRankFeaturesPositiveScoreImpact` | `v bool` | 高い値が関連性をブーストするか |

---

## 特殊

### NewBinaryProperty

Base64 エンコードされたバイナリデータ。

```go
esv8.NewBinaryProperty(opts ...BinaryPropertyOption) *types.BinaryProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithBinaryDocValues` | `v bool` | doc values を有効にするか |
| `WithBinaryStore` | `v bool` | フィールド値を個別に保存するか |

### NewTokenCountProperty

解析されたトークンの整数カウント。

```go
esv8.NewTokenCountProperty(opts ...TokenCountPropertyOption) *types.TokenCountProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithTokenCountAnalyzer` | `v string` | トークンカウント用アナライザ |
| `WithTokenCountDocValues` | `v bool` | doc values を有効にするか |
| `WithTokenCountIndex` | `v bool` | フィールドを検索可能にするか |
| `WithTokenCountStore` | `v bool` | フィールド値を個別に保存するか |
| `WithTokenCountEnablePositionIncrements` | `v bool` | 位置インクリメントをカウントするか |

### NewPercolatorProperty

percolate クエリで使用するクエリを格納。

```go
esv8.NewPercolatorProperty(opts ...PercolatorPropertyOption) *types.PercolatorProperty
```

オプションなし。

### NewFieldAliasProperty

既存フィールドの代替名。

```go
esv8.NewFieldAliasProperty(opts ...FieldAliasPropertyOption) *types.FieldAliasProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithFieldAliasPath` | `v string` | ターゲットフィールドへのパス |

### NewHistogramProperty

事前集計されたヒストグラムデータ。

```go
esv8.NewHistogramProperty(opts ...HistogramPropertyOption) *types.HistogramProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithHistogramIgnoreMalformed` | `v bool` | 不正な値を無視するか |

### NewExponentialHistogramProperty（v9 のみ）

事前集計された指数ヒストグラムデータ。

```go
esv9.NewExponentialHistogramProperty(opts ...ExponentialHistogramPropertyOption) *types.ExponentialHistogramProperty
```

オプションなし。`esv9` でのみ利用可能。

### NewVersionProperty

セマンティックバージョン順序を持つソフトウェアバージョン文字列。

```go
esv8.NewVersionProperty(opts ...VersionPropertyOption) *types.VersionProperty
```

オプションなし。

### NewAggregateMetricDoubleProperty

事前集計されたメトリック値。

```go
esv8.NewAggregateMetricDoubleProperty(opts ...AggregateMetricDoublePropertyOption) *types.AggregateMetricDoubleProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithAggregateMetricDoubleDefaultMetric` | `v string` | クエリ時のデフォルトメトリック |
| `WithAggregateMetricDoubleMetrics` | `v []string` | 保存するメトリックリスト（例：`["min","max","sum","value_count"]`） |
| `WithAggregateMetricDoubleIgnoreMalformed` | `v bool` | 不正な値を無視するか |

### NewSemanticTextProperty

推論エンドポイントを使用した ML ベースのセマンティック検索フィールド。

```go
esv8.NewSemanticTextProperty(opts ...SemanticTextPropertyOption) *types.SemanticTextProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithSemanticTextInferenceId` | `v string` | インデクシング用の推論エンドポイント ID |
| `WithSemanticTextSearchInferenceId` | `v string` | 検索用の推論エンドポイント ID |

### NewDynamicProperty

ダイナミックフィールドマッピング用のテンプレートベースプロパティ。

```go
esv8.NewDynamicProperty(opts ...DynamicPropertyOption) *types.DynamicProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithDynamicAnalyzer` | `v string` | テキスト系型用のアナライザ |
| `WithDynamicSearchAnalyzer` | `v string` | 検索用アナライザ |
| `WithDynamicCoerce` | `v bool` | 値を変換するか |
| `WithDynamicDocValues` | `v bool` | doc values を有効にするか |
| `WithDynamicEnabled` | `v bool` | フィールドを有効にするか |
| `WithDynamicFormat` | `v string` | 日付フォーマット文字列 |
| `WithDynamicIgnoreMalformed` | `v bool` | 不正な値を無視するか |
| `WithDynamicIndex` | `v bool` | フィールドを検索可能にするか |
| `WithDynamicStore` | `v bool` | フィールド値を個別に保存するか |
| `WithDynamicNorms` | `v bool` | norms を保存するか |
| `WithDynamicLocale` | `v string` | 日付解析のロケール |

---

## プラグイン依存

これらのプロパティ型は、特定の Elasticsearch プラグインのインストールが必要です。

### NewMurmur3HashProperty

フィールド値の murmur3 ハッシュを格納。`mapper-murmur3` プラグインが必要。

```go
esv8.NewMurmur3HashProperty(opts ...Murmur3HashPropertyOption) *types.Murmur3HashProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithMurmur3HashDocValues` | `v bool` | doc values を有効にするか |
| `WithMurmur3HashStore` | `v bool` | フィールド値を個別に保存するか |

### NewIcuCollationProperty

ICU 照合順序ベースのソートを持つキーワードフィールド。`analysis-icu` プラグインが必要。

```go
esv8.NewIcuCollationProperty(opts ...IcuCollationPropertyOption) *types.IcuCollationProperty
```

| オプション | パラメータ | 説明 |
|-----------|-----------|------|
| `WithIcuCollationLanguage` | `v string` | 言語コード（例：`"en"`） |
| `WithIcuCollationCountry` | `v string` | 国コード（例：`"US"`） |
| `WithIcuCollationDocValues` | `v bool` | doc values を有効にするか |
| `WithIcuCollationIndex` | `v bool` | フィールドを検索可能にするか |
| `WithIcuCollationStore` | `v bool` | フィールド値を個別に保存するか |
| `WithIcuCollationNullValue` | `v string` | `null` の代替値 |
| `WithIcuCollationNorms` | `v bool` | norms を保存するか |
| `WithIcuCollationRules` | `v string` | ICU 照合順序ルール文字列 |
| `WithIcuCollationVariant` | `v string` | 照合順序バリアント |
| `WithIcuCollationCaseLevel` | `v bool` | 大文字小文字レベル比較を有効にするか |
| `WithIcuCollationNumeric` | `v bool` | 数値照合順序を有効にするか |
| `WithIcuCollationHiraganaQuaternaryMode` | `v bool` | ひらがな4次モードを有効にするか |
| `WithIcuCollationVariableTop` | `v string` | 照合順序の variable top 設定 |

---

## バージョン間の差異（v8 と v9）

| プロパティ | v8 | v9 |
|-----------|----|----|
| 上記すべて（注記のあるものを除く） | あり | あり |
| `NewExponentialHistogramProperty` | -- | あり |

すべてのオプション関数のシグネチャは `esv8` と `esv9` で同一です。`esv8` 向けに書かれたコードは `esv9` でもそのまま動作し、`esv9` でのみ利用可能な `NewExponentialHistogramProperty` が追加されています。
