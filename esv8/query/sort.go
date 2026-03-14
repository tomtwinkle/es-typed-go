package query

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/distanceunit"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/geodistancetype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptsorttype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortmode"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// MissingFirst places documents with missing field values at the beginning of the sort order.
const MissingFirst = "_first"

// MissingLast places documents with missing field values at the end of the sort order.
const MissingLast = "_last"

// SortBuilder constructs Elasticsearch sort specifications using method chaining.
type SortBuilder struct {
	sorts []types.SortCombinations
}

// NewSort returns a new empty SortBuilder.
func NewSort() *SortBuilder {
	return &SortBuilder{}
}

// Field adds a sort on the given field with the specified order.
func (s *SortBuilder) Field(name estype.Field, order sortorder.SortOrder) *SortBuilder {
	s.sorts = append(s.sorts, types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			string(name): {Order: &order},
		},
	})
	return s
}

// FieldWithMissing adds a sort on the given field with a missing value position.
// Use MissingFirst or MissingLast to control where documents with missing field values appear.
func (s *SortBuilder) FieldWithMissing(name estype.Field, order sortorder.SortOrder, missing string) *SortBuilder {
	var m types.Missing = missing
	s.sorts = append(s.sorts, types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			string(name): {Order: &order, Missing: m},
		},
	})
	return s
}

// FieldNested adds a nested field sort with a sort mode.
func (s *SortBuilder) FieldNested(name estype.Field, order sortorder.SortOrder, path estype.Field, mode sortmode.SortMode) *SortBuilder {
	s.sorts = append(s.sorts, types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			string(name): {
				Order:  &order,
				Mode:   &mode,
				Nested: &types.NestedSortValue{Path: string(path)},
			},
		},
	})
	return s
}

// FieldCustom adds a sort on the given field with a fully customized FieldSort.
func (s *SortBuilder) FieldCustom(name estype.Field, fs types.FieldSort) *SortBuilder {
	s.sorts = append(s.sorts, types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			string(name): fs,
		},
	})
	return s
}

// ScoreDesc adds a sort by relevance score in descending order.
func (s *SortBuilder) ScoreDesc() *SortBuilder {
	order := sortorder.Desc
	s.sorts = append(s.sorts, types.SortOptions{
		Score_: &types.ScoreSort{Order: &order},
	})
	return s
}

// ScoreAsc adds a sort by relevance score in ascending order.
func (s *SortBuilder) ScoreAsc() *SortBuilder {
	order := sortorder.Asc
	s.sorts = append(s.sorts, types.SortOptions{
		Score_: &types.ScoreSort{Order: &order},
	})
	return s
}

// DocAsc adds a sort by index order in ascending order.
func (s *SortBuilder) DocAsc() *SortBuilder {
	order := sortorder.Asc
	s.sorts = append(s.sorts, types.SortOptions{
		Doc_: &types.ScoreSort{Order: &order},
	})
	return s
}

// DocDesc adds a sort by index order in descending order.
func (s *SortBuilder) DocDesc() *SortBuilder {
	order := sortorder.Desc
	s.sorts = append(s.sorts, types.SortOptions{
		Doc_: &types.ScoreSort{Order: &order},
	})
	return s
}

// Build returns the constructed sort slice ready for use in a search request.
func (s *SortBuilder) Build() []types.SortCombinations {
	return s.sorts
}

// ---------------------------------------------------------------------------
// Geo Distance Sort
// ---------------------------------------------------------------------------

// GeoDistanceSortOption is a functional option for configuring a GeoDistanceSort.
type GeoDistanceSortOption func(*types.GeoDistanceSort)

// GeoDistance adds a geo distance sort on the given field from a reference location.
func (s *SortBuilder) GeoDistance(name estype.Field, location types.GeoLocation, order sortorder.SortOrder, opts ...GeoDistanceSortOption) *SortBuilder {
	gs := &types.GeoDistanceSort{
		GeoDistanceSort: map[string][]types.GeoLocation{
			string(name): {location},
		},
		Order: &order,
	}
	for _, opt := range opts {
		opt(gs)
	}
	s.sorts = append(s.sorts, types.SortOptions{
		GeoDistance_: gs,
	})
	return s
}

// GeoDistanceCustom adds a geo distance sort with a fully customized GeoDistanceSort.
func (s *SortBuilder) GeoDistanceCustom(gs types.GeoDistanceSort) *SortBuilder {
	s.sorts = append(s.sorts, types.SortOptions{
		GeoDistance_: &gs,
	})
	return s
}

// WithGeoDistanceUnit sets the distance unit for the geo distance sort.
func WithGeoDistanceUnit(unit distanceunit.DistanceUnit) GeoDistanceSortOption {
	return func(gs *types.GeoDistanceSort) { gs.Unit = &unit }
}

// WithGeoDistanceType sets the distance calculation type (arc or plane).
func WithGeoDistanceType(dt geodistancetype.GeoDistanceType) GeoDistanceSortOption {
	return func(gs *types.GeoDistanceSort) { gs.DistanceType = &dt }
}

// WithGeoDistanceMode sets the sort mode for multi-valued geo fields.
func WithGeoDistanceMode(mode sortmode.SortMode) GeoDistanceSortOption {
	return func(gs *types.GeoDistanceSort) { gs.Mode = &mode }
}

// WithGeoDistanceIgnoreUnmapped sets whether to ignore unmapped fields.
func WithGeoDistanceIgnoreUnmapped(ignore bool) GeoDistanceSortOption {
	return func(gs *types.GeoDistanceSort) { gs.IgnoreUnmapped = &ignore }
}

// WithGeoDistanceNested sets the nested sort configuration.
func WithGeoDistanceNested(nested *types.NestedSortValue) GeoDistanceSortOption {
	return func(gs *types.GeoDistanceSort) { gs.Nested = nested }
}

// ---------------------------------------------------------------------------
// Script Sort
// ---------------------------------------------------------------------------

// ScriptSortOption is a functional option for configuring a ScriptSort.
type ScriptSortOption func(*types.ScriptSort)

// Script adds a script-based sort.
func (s *SortBuilder) Script(script types.Script, sortType scriptsorttype.ScriptSortType, order sortorder.SortOrder, opts ...ScriptSortOption) *SortBuilder {
	ss := &types.ScriptSort{
		Script: script,
		Type:   &sortType,
		Order:  &order,
	}
	for _, opt := range opts {
		opt(ss)
	}
	s.sorts = append(s.sorts, types.SortOptions{
		Script_: ss,
	})
	return s
}

// ScriptCustom adds a script-based sort with a fully customized ScriptSort.
func (s *SortBuilder) ScriptCustom(ss types.ScriptSort) *SortBuilder {
	s.sorts = append(s.sorts, types.SortOptions{
		Script_: &ss,
	})
	return s
}

// WithScriptSortMode sets the sort mode for multi-valued fields.
func WithScriptSortMode(mode sortmode.SortMode) ScriptSortOption {
	return func(ss *types.ScriptSort) { ss.Mode = &mode }
}

// WithScriptSortNested sets the nested sort configuration.
func WithScriptSortNested(nested *types.NestedSortValue) ScriptSortOption {
	return func(ss *types.ScriptSort) { ss.Nested = nested }
}
