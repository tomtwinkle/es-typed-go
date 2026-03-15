package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv8/query"
)

func TestNewSearch_Empty(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().Build()
	assert.DeepEqual(t, types.Query{}, params.Query)
	assert.Assert(t, len(params.Sort) == 0)
	assert.Assert(t, params.Aggregations == nil)
	assert.Assert(t, params.Highlight == nil)
	assert.Assert(t, params.Collapse == nil)
	assert.Assert(t, params.ScriptFields == nil)
	assert.Equal(t, 0, params.Size)
	assert.Equal(t, 0, params.From)
}

func TestSearchBuilder_Where(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		Where(query.TermValue(estype.Field("status"), "active")).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Filter) == 1)
	assert.Assert(t, params.Query.Bool.Filter[0].Term != nil)
	assert.Equal(t, "active", params.Query.Bool.Filter[0].Term["status"].Value)
}

func TestSearchBuilder_Where_Multiple(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		Where(query.TermValue(estype.Field("status"), "active")).
		Where(query.TermValue(estype.Field("type"), "document")).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Filter) == 2)
}

func TestSearchBuilder_Must(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		Must(query.TermValue(estype.Field("status"), "active")).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Must) == 1)
}

func TestSearchBuilder_Should(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		Should(
			query.MatchPhrase(estype.Field("title"), "keyword"),
			query.MatchPhrase(estype.Field("name"), "keyword"),
		).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Should) == 2)
}

func TestSearchBuilder_MustNot(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		MustNot(query.ExistsField(estype.Field("tags"))).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.MustNot) == 1)
}

func TestSearchBuilder_MinimumShouldMatch(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		Should(
			query.MatchPhrase(estype.Field("title"), "keyword"),
			query.MatchPhrase(estype.Field("name"), "keyword"),
		).
		MinimumShouldMatch(1).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Should) == 2)
	assert.Equal(t, 1, params.Query.Bool.MinimumShouldMatch)
}

func TestSearchBuilder_QueryOverridesBoolClauses(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		Where(query.TermValue(estype.Field("status"), "active")).
		Query(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Build()

	// Query() takes precedence over Where()
	assert.Assert(t, params.Query.MatchAll != nil)
	assert.Assert(t, params.Query.Bool == nil)
}

func TestSearchBuilder_Sort(t *testing.T) {
	t.Parallel()
	sorts := query.NewSort().
		Field(estype.Field("date"), sortorder.Desc).
		Build()

	params := query.NewSearch().
		Where(query.TermValue(estype.Field("status"), "active")).
		Sort(sorts...).
		Build()

	assert.Assert(t, len(params.Sort) == 1)
}

func TestSearchBuilder_LimitAndOffset(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		Limit(20).
		Offset(40).
		Build()

	assert.Equal(t, 20, params.Size)
	assert.Equal(t, 40, params.From)
}

func TestSearchBuilder_Aggregation(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().
		Terms("by_category", estype.Field("category")).
		Build()

	params := query.NewSearch().
		Aggregation(aggs).
		Build()

	assert.Assert(t, params.Aggregations != nil)
	_, ok := params.Aggregations["by_category"]
	assert.Assert(t, ok)
}

func TestSearchBuilder_Highlight(t *testing.T) {
	t.Parallel()
	h := &types.Highlight{
		Fields: map[string]types.HighlightField{
			"title": {},
		},
	}
	params := query.NewSearch().
		Highlight(h).
		Build()

	assert.Assert(t, params.Highlight != nil)
	_, ok := params.Highlight.Fields["title"]
	assert.Assert(t, ok)
}

func TestSearchBuilder_Collapse(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		Collapse(&types.FieldCollapse{Field: "category"}).
		Build()

	assert.Assert(t, params.Collapse != nil)
	assert.Equal(t, "category", params.Collapse.Field)
}

func TestSearchBuilder_ScriptFields(t *testing.T) {
	t.Parallel()
	source := "doc['price'].value * 2"
	params := query.NewSearch().
		ScriptFields(map[string]types.ScriptField{
			"value": {Script: types.Script{Source: &source}},
		}).
		Build()

	assert.Assert(t, params.ScriptFields != nil)
	sf, ok := params.ScriptFields["value"]
	assert.Assert(t, ok)
	assert.Equal(t, source, *sf.Script.Source)
}

func TestSearchBuilder_CombinedBoolClauses(t *testing.T) {
	t.Parallel()
	params := query.NewSearch().
		Where(query.TermValue(estype.Field("status"), "active")).
		Must(query.MatchPhrase(estype.Field("title"), "keyword")).
		Should(query.TermValue(estype.Field("category"), "a")).
		MustNot(query.ExistsField(estype.Field("tags"))).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Filter) == 1)
	assert.Assert(t, len(params.Query.Bool.Must) == 1)
	assert.Assert(t, len(params.Query.Bool.Should) == 1)
	assert.Assert(t, len(params.Query.Bool.MustNot) == 1)
}

func TestSearchBuilder_FullChaining(t *testing.T) {
	t.Parallel()
	// Demonstrates the ActiveRecord-style fluent API.
	aggs := query.NewAggregations().
		DateHistogram("by_month", estype.Field("date"), calendarinterval.Month).
		Build()

	sorts := query.NewSort().
		Field(estype.Field("date"), sortorder.Desc).
		Field(estype.Field("id"), sortorder.Asc).
		Build()

	params := query.NewSearch().
		Where(
			query.TermValue(estype.Field("status"), "active"),
			query.TermsValues(estype.Field("category"), query.FieldValues("a", "b")...),
		).
		MustNot(query.NotExists(estype.Field("tags"))).
		Sort(sorts...).
		Limit(10).
		Offset(0).
		Aggregation(aggs).
		Highlight(&types.Highlight{
			Fields: map[string]types.HighlightField{"title": {}},
		}).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Filter) == 2)
	assert.Assert(t, len(params.Query.Bool.MustNot) == 1)
	assert.Assert(t, len(params.Sort) == 2)
	assert.Equal(t, 10, params.Size)
	assert.Equal(t, 0, params.From)
	assert.Assert(t, params.Aggregations != nil)
	assert.Assert(t, params.Highlight != nil)
}

func TestSearchBuilder_NoBoolClausesNoQuery(t *testing.T) {
	t.Parallel()
	// When no Where/Must/Should/MustNot and no Query() is set,
	// the resulting query should be a zero-value Query.
	params := query.NewSearch().
		Limit(10).
		Build()

	assert.DeepEqual(t, types.Query{}, params.Query)
	assert.Equal(t, 10, params.Size)
}
