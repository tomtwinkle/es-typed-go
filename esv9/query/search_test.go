package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/calendarinterval"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortorder"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/esv9/query"
)

func TestNewSearch_Empty(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().Build()

	assert.DeepEqual(t, types.Query{}, params.Query)
	assert.Assert(t, len(params.Sort) == 0)
	assert.Assert(t, params.Aggregations.Build() == nil)
	assert.Assert(t, params.Highlight == nil)
	assert.Assert(t, params.Collapse == nil)
	assert.Assert(t, params.ScriptFields == nil)
	assert.Assert(t, params.TrackTotalHits == nil)
	assert.Equal(t, 0, params.Size)
	assert.Equal(t, 0, params.From)
}

func TestSearchBuilder_Where(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		Where(query.TermValue(FieldStatus, "active")).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Filter) == 1)
	assert.Assert(t, params.Query.Bool.Filter[0].Term != nil)
	assert.Equal(t, "active", params.Query.Bool.Filter[0].Term[string(FieldStatus)].Value)
}

func TestSearchBuilder_Where_Multiple(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		Where(query.TermValue(FieldStatus, "active")).
		Where(query.TermValue(FieldType, "document")).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Filter) == 2)
}

func TestSearchBuilder_Must(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		Must(query.TermValue(FieldStatus, "active")).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Must) == 1)
}

func TestSearchBuilder_Should(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		Should(
			query.MatchPhrase(FieldTitle, "keyword"),
			query.MatchPhrase(FieldName, "keyword"),
		).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Should) == 2)
}

func TestSearchBuilder_MustNot(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		MustNot(query.ExistsField(FieldTags)).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.MustNot) == 1)
}

func TestSearchBuilder_MinimumShouldMatch(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		Should(
			query.MatchPhrase(FieldTitle, "keyword"),
			query.MatchPhrase(FieldName, "keyword"),
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
		Where(query.TermValue(FieldStatus, "active")).
		Query(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Build()

	assert.Assert(t, params.Query.MatchAll != nil)
	assert.Assert(t, params.Query.Bool == nil)
}

func TestSearchBuilder_Sort(t *testing.T) {
	t.Parallel()

	sorts := query.NewSort().
		Field(FieldDate, sortorder.Desc).
		Build()

	params := query.NewSearch().
		Where(query.TermValue(FieldStatus, "active")).
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

	termsDef := query.StringTermsAgg("by_category", FieldCategory)
	aggs := query.Aggs(termsDef)

	params := query.NewSearch().
		Aggregation(aggs).
		Build()

	built := params.Aggregations.Build()
	assert.Assert(t, built != nil)

	agg, ok := built["by_category"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Terms != nil)
	assert.Equal(t, string(FieldCategory), *agg.Terms.Field)
}

func TestSearchBuilder_Highlight(t *testing.T) {
	t.Parallel()

	h := &types.Highlight{
		Fields: []map[string]types.HighlightField{
			{FieldTitle.String(): {}},
		},
	}

	params := query.NewSearch().
		Highlight(h).
		Build()

	assert.Assert(t, params.Highlight != nil)
	assert.Assert(t, len(params.Highlight.Fields) == 1)
}

func TestSearchBuilder_Collapse(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		Collapse(&types.FieldCollapse{Field: FieldCategory.String()}).
		Build()

	assert.Assert(t, params.Collapse != nil)
	assert.Equal(t, FieldCategory.String(), params.Collapse.Field)
}

func TestSearchBuilder_ScriptFields(t *testing.T) {
	t.Parallel()

	source := "doc['price'].value * 2"
	params := query.NewSearch().
		ScriptFields(map[string]types.ScriptField{
			FieldValue.String(): {Script: types.Script{Source: source}},
		}).
		Build()

	assert.Assert(t, params.ScriptFields != nil)
	sf, ok := params.ScriptFields[FieldValue.String()]
	assert.Assert(t, ok)
	assert.Equal(t, source, sf.Script.Source)
}

func TestSearchBuilder_TrackTotalHits(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		TrackTotalHits(true).
		Build()

	assert.Equal(t, true, params.TrackTotalHits)
}

func TestSearchBuilder_CombinedBoolClauses(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		Where(query.TermValue(FieldStatus, "active")).
		Must(query.MatchPhrase(FieldTitle, "keyword")).
		Should(query.TermValue(FieldCategory, "a")).
		MustNot(query.ExistsField(FieldTags)).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Filter) == 1)
	assert.Assert(t, len(params.Query.Bool.Must) == 1)
	assert.Assert(t, len(params.Query.Bool.Should) == 1)
	assert.Assert(t, len(params.Query.Bool.MustNot) == 1)
}

func TestSearchBuilder_FullChaining(t *testing.T) {
	t.Parallel()

	dateHistDef := query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month)
	aggs := query.Aggs(dateHistDef)

	sorts := query.NewSort().
		Field(FieldDate, sortorder.Desc).
		Field(FieldId, sortorder.Asc).
		Build()

	params := query.NewSearch().
		Where(
			query.TermValue(FieldStatus, "active"),
			query.TermsValues(FieldCategory, query.FieldValues("a", "b")...),
		).
		MustNot(query.NotExists(FieldTags)).
		Sort(sorts...).
		Limit(10).
		Offset(0).
		Aggregation(aggs).
		Highlight(&types.Highlight{
			Fields: []map[string]types.HighlightField{
				{FieldTitle.String(): {}},
			},
		}).
		Build()

	assert.Assert(t, params.Query.Bool != nil)
	assert.Assert(t, len(params.Query.Bool.Filter) == 2)
	assert.Assert(t, len(params.Query.Bool.MustNot) == 1)
	assert.Assert(t, len(params.Sort) == 2)
	assert.Equal(t, 10, params.Size)
	assert.Equal(t, 0, params.From)
	assert.Assert(t, params.Aggregations.Build() != nil)
	assert.Assert(t, params.Highlight != nil)
}

func TestSearchBuilder_NoBoolClausesNoQuery(t *testing.T) {
	t.Parallel()

	params := query.NewSearch().
		Limit(10).
		Build()

	assert.DeepEqual(t, types.Query{}, params.Query)
	assert.Equal(t, 10, params.Size)
}
