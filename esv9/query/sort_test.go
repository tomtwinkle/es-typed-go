package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/distanceunit"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/geodistancetype"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/scriptsorttype"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortmode"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortorder"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/esv9/query"
)

func TestNewSort_Empty(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().Build()
	assert.Assert(t, len(sorts) == 0)
}

func TestSortBuilder_Field(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		Field(FieldDate, sortorder.Desc).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs, ok := so.SortOptions[string(FieldDate)]
	assert.Assert(t, ok)
	assert.Equal(t, sortorder.Desc, *fs.Order)
}

func TestSortBuilder_FieldWithMissing_Last(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		FieldWithMissing(FieldNameKeyword, sortorder.Asc, query.MissingLast).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs := so.SortOptions[string(FieldNameKeyword)]
	assert.Equal(t, sortorder.Asc, *fs.Order)
	assert.Equal(t, "_last", fs.Missing.(string))
}

func TestSortBuilder_FieldWithMissing_First(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		FieldWithMissing(FieldDate, sortorder.Desc, query.MissingFirst).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs := so.SortOptions[string(FieldDate)]
	assert.Equal(t, sortorder.Desc, *fs.Order)
	assert.Equal(t, "_first", fs.Missing.(string))
}

func TestSortBuilder_FieldNested(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		FieldNested(FieldItemsDate, sortorder.Asc, FieldItems, sortmode.Min).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs := so.SortOptions[string(FieldItemsDate)]
	assert.Equal(t, sortorder.Asc, *fs.Order)
	assert.Equal(t, sortmode.Min, *fs.Mode)
	assert.Assert(t, fs.Nested != nil)
	assert.Equal(t, string(FieldItems), fs.Nested.Path)
}

func TestSortBuilder_FieldCustom(t *testing.T) {
	t.Parallel()
	mode := sortmode.Max
	order := sortorder.Desc
	sorts := query.NewSort().
		FieldCustom(FieldPrice, types.FieldSort{
			Order: &order,
			Mode:  &mode,
		}).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs := so.SortOptions[string(FieldPrice)]
	assert.Equal(t, sortorder.Desc, *fs.Order)
	assert.Equal(t, sortmode.Max, *fs.Mode)
}

func TestSortBuilder_ScoreDesc(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().ScoreDesc().Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so.Score_ != nil)
	assert.Equal(t, sortorder.Desc, *so.Score_.Order)
}

func TestSortBuilder_ScoreAsc(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().ScoreAsc().Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so.Score_ != nil)
	assert.Equal(t, sortorder.Asc, *so.Score_.Order)
}

func TestSortBuilder_DocAsc(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().DocAsc().Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so.Doc_ != nil)
	assert.Equal(t, sortorder.Asc, *so.Doc_.Order)
}

func TestSortBuilder_DocDesc(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().DocDesc().Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so.Doc_ != nil)
	assert.Equal(t, sortorder.Desc, *so.Doc_.Order)
}

func TestSortBuilder_Chaining(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		FieldNested(FieldItemsDate, sortorder.Asc, FieldItems, sortmode.Min).
		Field(FieldId, sortorder.Asc).
		Build()

	assert.Assert(t, len(sorts) == 2)
	so1, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs1 := so1.SortOptions[string(FieldItemsDate)]
	assert.Equal(t, sortorder.Asc, *fs1.Order)
	assert.Assert(t, fs1.Nested != nil)
	so2, ok := sorts[1].(types.SortOptions)
	assert.Assert(t, ok)
	fs2 := so2.SortOptions[string(FieldId)]
	assert.Equal(t, sortorder.Asc, *fs2.Order)
}

func TestSortBuilder_MissingDirectionPattern(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		order   sortorder.SortOrder
		missing string
	}{
		{"asc_missing_last", sortorder.Asc, query.MissingLast},
		{"desc_missing_first", sortorder.Desc, query.MissingFirst},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sorts := query.NewSort().
				FieldWithMissing(FieldDate, tt.order, tt.missing).
				Field(FieldId, tt.order).
				Build()

			assert.Assert(t, len(sorts) == 2)
			so, ok := sorts[0].(types.SortOptions)
			assert.Assert(t, ok)
			fs := so.SortOptions[string(FieldDate)]
			assert.Equal(t, tt.order, *fs.Order)
			assert.Equal(t, tt.missing, fs.Missing.(string))
		})
	}
}

func TestMissingConstants(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "_first", query.MissingFirst)
	assert.Equal(t, "_last", query.MissingLast)
}

func TestSortBuilder_GeoDistance(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		GeoDistance(
			FieldLocation,
			types.LatLonGeoLocation{Lat: 40.7, Lon: -74.0},
			sortorder.Asc,
		).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so.GeoDistance_ != nil)
	assert.Equal(t, sortorder.Asc, *so.GeoDistance_.Order)
	locs, ok := so.GeoDistance_.GeoDistanceSort[string(FieldLocation)]
	assert.Assert(t, ok)
	assert.Assert(t, len(locs) == 1)
}

func TestSortBuilder_GeoDistance_WithOptions(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		GeoDistance(
			FieldLocation,
			types.LatLonGeoLocation{Lat: 40.7, Lon: -74.0},
			sortorder.Asc,
			query.WithGeoDistanceUnit(distanceunit.Kilometers),
			query.WithGeoDistanceType(geodistancetype.Arc),
			query.WithGeoDistanceMode(sortmode.Min),
			query.WithGeoDistanceIgnoreUnmapped(true),
		).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	gs := so.GeoDistance_
	assert.Assert(t, gs != nil)
	assert.Equal(t, sortorder.Asc, *gs.Order)
	assert.Equal(t, distanceunit.Kilometers, *gs.Unit)
	assert.Equal(t, geodistancetype.Arc, *gs.DistanceType)
	assert.Equal(t, sortmode.Min, *gs.Mode)
	assert.Equal(t, true, *gs.IgnoreUnmapped)
}

func TestSortBuilder_GeoDistance_WithNested(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		GeoDistance(
			FieldItemsLocation,
			types.LatLonGeoLocation{Lat: 40.7, Lon: -74.0},
			sortorder.Asc,
			query.WithGeoDistanceNested(&types.NestedSortValue{Path: string(FieldItems)}),
		).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	gs := so.GeoDistance_
	assert.Assert(t, gs != nil)
	assert.Equal(t, sortorder.Asc, *gs.Order)
	assert.Assert(t, gs.Nested != nil)
	assert.Equal(t, string(FieldItems), gs.Nested.Path)
}

func TestSortBuilder_GeoDistanceCustom(t *testing.T) {
	t.Parallel()
	order := sortorder.Desc
	unit := distanceunit.Miles
	sorts := query.NewSort().
		GeoDistanceCustom(types.GeoDistanceSort{
			GeoDistanceSort: map[string][]types.GeoLocation{
				string(FieldLocation): {types.LatLonGeoLocation{Lat: 40.7, Lon: -74.0}},
			},
			Order: &order,
			Unit:  &unit,
		}).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so.GeoDistance_ != nil)
	assert.Equal(t, sortorder.Desc, *so.GeoDistance_.Order)
	assert.Equal(t, distanceunit.Miles, *so.GeoDistance_.Unit)
}

func TestSortBuilder_Script(t *testing.T) {
	t.Parallel()
	var source types.ScriptSource = "doc['price'].value * params.factor"
	sorts := query.NewSort().
		Script(
			types.Script{Source: source},
			scriptsorttype.Number,
			sortorder.Desc,
		).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so.Script_ != nil)
	assert.Equal(t, sortorder.Desc, *so.Script_.Order)
	assert.Equal(t, scriptsorttype.Number, *so.Script_.Type)
	assert.Equal(t, "doc['price'].value * params.factor", so.Script_.Script.Source.(string))
}

func TestSortBuilder_Script_WithOptions(t *testing.T) {
	t.Parallel()
	var source types.ScriptSource = "doc['price'].value * params.factor"
	sorts := query.NewSort().
		Script(
			types.Script{Source: source},
			scriptsorttype.Number,
			sortorder.Asc,
			query.WithScriptSortMode(sortmode.Avg),
		).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	ss := so.Script_
	assert.Assert(t, ss != nil)
	assert.Equal(t, sortorder.Asc, *ss.Order)
	assert.Equal(t, scriptsorttype.Number, *ss.Type)
	assert.Equal(t, sortmode.Avg, *ss.Mode)
}

func TestSortBuilder_Script_WithNested(t *testing.T) {
	t.Parallel()
	var source types.ScriptSource = "doc['items.price'].value"
	sorts := query.NewSort().
		Script(
			types.Script{Source: source},
			scriptsorttype.Number,
			sortorder.Desc,
			query.WithScriptSortNested(&types.NestedSortValue{Path: string(FieldItems)}),
		).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	ss := so.Script_
	assert.Assert(t, ss != nil)
	assert.Equal(t, sortorder.Desc, *ss.Order)
	assert.Assert(t, ss.Nested != nil)
	assert.Equal(t, string(FieldItems), ss.Nested.Path)
}

func TestSortBuilder_ScriptCustom(t *testing.T) {
	t.Parallel()
	var source types.ScriptSource = "doc['name'].value"
	order := sortorder.Asc
	stype := scriptsorttype.String
	sorts := query.NewSort().
		ScriptCustom(types.ScriptSort{
			Script: types.Script{Source: source},
			Type:   &stype,
			Order:  &order,
		}).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so.Script_ != nil)
	assert.Equal(t, sortorder.Asc, *so.Script_.Order)
	assert.Equal(t, scriptsorttype.String, *so.Script_.Type)
}

func TestSortBuilder_Chaining_AllSortTypes(t *testing.T) {
	t.Parallel()
	var source types.ScriptSource = "doc['price'].value"
	sorts := query.NewSort().
		ScoreDesc().
		Field(FieldDate, sortorder.Desc).
		GeoDistance(FieldLocation, types.LatLonGeoLocation{Lat: 40.7, Lon: -74.0}, sortorder.Asc).
		Script(types.Script{Source: source}, scriptsorttype.Number, sortorder.Asc).
		DocAsc().
		Build()

	assert.Assert(t, len(sorts) == 5)
	// _score
	so0, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so0.Score_ != nil)
	// field
	so1, ok := sorts[1].(types.SortOptions)
	assert.Assert(t, ok)
	_, ok = so1.SortOptions[string(FieldDate)]
	assert.Assert(t, ok)
	// _geo_distance
	so2, ok := sorts[2].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so2.GeoDistance_ != nil)
	// _script
	so3, ok := sorts[3].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so3.Script_ != nil)
	// _doc
	so4, ok := sorts[4].(types.SortOptions)
	assert.Assert(t, ok)
	assert.Assert(t, so4.Doc_ != nil)
}
