// Package query provides bool query building helpers for Elasticsearch v8.
//
// Deprecated: import github.com/tomtwinkle/es-typed-go/query instead.
package query

import sharedquery "github.com/tomtwinkle/es-typed-go/query"

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.BoolQueryBuilder.
type BoolQueryBuilder = sharedquery.BoolQueryBuilder

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.NewBoolQuery.
func NewBoolQuery() *sharedquery.BoolQueryBuilder { return sharedquery.NewBoolQuery() }
