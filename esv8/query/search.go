// Package query provides search building helpers for Elasticsearch v8.
//
// Deprecated: import github.com/tomtwinkle/es-typed-go/query instead.
package query

import sharedquery "github.com/tomtwinkle/es-typed-go/query"

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.SearchParams.
type SearchParams = sharedquery.SearchParams

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.SearchBuilder.
type SearchBuilder = sharedquery.SearchBuilder

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.NewSearch.
func NewSearch() *sharedquery.SearchBuilder { return sharedquery.NewSearch() }
