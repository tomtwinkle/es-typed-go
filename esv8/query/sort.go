// Package query provides sort building helpers for Elasticsearch v8.
//
// Deprecated: import github.com/tomtwinkle/es-typed-go/query instead.
package query

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/distanceunit"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/geodistancetype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortmode"
	sharedquery "github.com/tomtwinkle/es-typed-go/query"
)

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MissingFirst.
const MissingFirst = sharedquery.MissingFirst

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MissingLast.
const MissingLast = sharedquery.MissingLast

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.SortBuilder.
type SortBuilder = sharedquery.SortBuilder

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.GeoDistanceSortOption.
type GeoDistanceSortOption = sharedquery.GeoDistanceSortOption

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.ScriptSortOption.
type ScriptSortOption = sharedquery.ScriptSortOption

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.NewSort.
func NewSort() *sharedquery.SortBuilder { return sharedquery.NewSort() }

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.WithGeoDistanceUnit.
func WithGeoDistanceUnit(unit distanceunit.DistanceUnit) sharedquery.GeoDistanceSortOption {
	return sharedquery.WithGeoDistanceUnit(unit)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.WithGeoDistanceType.
func WithGeoDistanceType(dt geodistancetype.GeoDistanceType) sharedquery.GeoDistanceSortOption {
	return sharedquery.WithGeoDistanceType(dt)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.WithGeoDistanceMode.
func WithGeoDistanceMode(mode sortmode.SortMode) sharedquery.GeoDistanceSortOption {
	return sharedquery.WithGeoDistanceMode(mode)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.WithGeoDistanceIgnoreUnmapped.
func WithGeoDistanceIgnoreUnmapped(ignore bool) sharedquery.GeoDistanceSortOption {
	return sharedquery.WithGeoDistanceIgnoreUnmapped(ignore)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.WithGeoDistanceNested.
func WithGeoDistanceNested(nested *types.NestedSortValue) sharedquery.GeoDistanceSortOption {
	return sharedquery.WithGeoDistanceNested(nested)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.WithScriptSortMode.
func WithScriptSortMode(mode sortmode.SortMode) sharedquery.ScriptSortOption {
	return sharedquery.WithScriptSortMode(mode)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.WithScriptSortNested.
func WithScriptSortNested(nested *types.NestedSortValue) sharedquery.ScriptSortOption {
	return sharedquery.WithScriptSortNested(nested)
}
