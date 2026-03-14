package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
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
		Field("created_at", sortorder.Desc).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs, ok := so.SortOptions["created_at"]
	assert.Assert(t, ok)
	assert.Equal(t, sortorder.Desc, *fs.Order)
}

func TestSortBuilder_FieldWithMissing_Last(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		FieldWithMissing("name.keyword", sortorder.Asc, query.MissingLast).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs := so.SortOptions["name.keyword"]
	assert.Equal(t, sortorder.Asc, *fs.Order)
	assert.Equal(t, "_last", fs.Missing.(string))
}

func TestSortBuilder_FieldWithMissing_First(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		FieldWithMissing("date", sortorder.Desc, query.MissingFirst).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs := so.SortOptions["date"]
	assert.Equal(t, sortorder.Desc, *fs.Order)
	assert.Equal(t, "_first", fs.Missing.(string))
}

func TestSortBuilder_FieldNested(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		FieldNested("items.date", sortorder.Asc, "items", sortmode.Min).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs := so.SortOptions["items.date"]
	assert.Equal(t, sortorder.Asc, *fs.Order)
	assert.Equal(t, sortmode.Min, *fs.Mode)
	assert.Assert(t, fs.Nested != nil)
	assert.Equal(t, "items", fs.Nested.Path)
}

func TestSortBuilder_FieldCustom(t *testing.T) {
	t.Parallel()
	mode := sortmode.Max
	order := sortorder.Desc
	sorts := query.NewSort().
		FieldCustom("price", types.FieldSort{
			Order: &order,
			Mode:  &mode,
		}).
		Build()

	assert.Assert(t, len(sorts) == 1)
	so, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs := so.SortOptions["price"]
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
		FieldNested("items.date", sortorder.Asc, "items", sortmode.Min).
		Field("id", sortorder.Asc).
		Build()

	assert.Assert(t, len(sorts) == 2)
	so1, ok := sorts[0].(types.SortOptions)
	assert.Assert(t, ok)
	fs1 := so1.SortOptions["items.date"]
	assert.Equal(t, sortorder.Asc, *fs1.Order)
	assert.Assert(t, fs1.Nested != nil)
	so2, ok := sorts[1].(types.SortOptions)
	assert.Assert(t, ok)
	fs2 := so2.SortOptions["id"]
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
				FieldWithMissing("date", tt.order, tt.missing).
				Field("id", tt.order).
				Build()

			assert.Assert(t, len(sorts) == 2)
			so, ok := sorts[0].(types.SortOptions)
			assert.Assert(t, ok)
			fs := so.SortOptions["date"]
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
