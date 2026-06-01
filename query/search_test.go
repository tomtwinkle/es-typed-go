package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	v9types "github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/query"
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
	assert.Assert(t, !params.HasSize)
	assert.Assert(t, !params.HasFrom)
	assert.Assert(t, params.SearchAfter == nil)
	assert.Assert(t, params.TrackTotalHits == nil)
	assert.Assert(t, params.Source == nil)
	assert.Assert(t, params.Timeout == nil)
}

func TestSearchParams_ToRequest_EmptyQueryOmitsRequestQuery(t *testing.T) {
	t.Parallel()

	req := query.NewSearch().
		Limit(10).
		Offset(20).
		Build().
		ToRequest()

	assert.Assert(t, req != nil)
	assert.Assert(t, req.Query == nil)
	assert.Assert(t, req.Size != nil)
	assert.Equal(t, 10, *req.Size)
	assert.Assert(t, req.From != nil)
	assert.Equal(t, 20, *req.From)
	assert.Assert(t, req.Timeout != nil)
	assert.Equal(t, "10s", *req.Timeout)
	assert.Assert(t, req.Source_)
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

	// Query() takes precedence over Where()
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
	assert.Assert(t, params.HasSize)
	assert.Assert(t, params.HasFrom)
}

func TestSearchBuilder_AdvancedOptions(t *testing.T) {
	t.Parallel()

	searchAfter := []types.FieldValue{"cursor-1", 42}
	params := query.NewSearch().
		Limit(0).
		Offset(0).
		SearchAfter("stale").
		SearchAfter(searchAfter...).
		TrackTotalHits(true).
		TrackTotalHits(100).
		Source(true).
		Source(false).
		Timeout("5s").
		Timeout("30s").
		Build()

	searchAfter[0] = "mutated"

	assert.Equal(t, 0, params.Size)
	assert.Assert(t, params.HasSize)
	assert.Equal(t, 0, params.From)
	assert.Assert(t, params.HasFrom)
	assert.DeepEqual(t, []types.FieldValue{"cursor-1", 42}, params.SearchAfter)
	assert.Equal(t, 100, params.TrackTotalHits)
	assert.Equal(t, false, params.Source)
	assert.Assert(t, params.Timeout != nil)
	assert.Equal(t, "30s", *params.Timeout)
}

func TestSearchParams_ToRequest_AdvancedOptions(t *testing.T) {
	t.Parallel()

	timeout := "30s"
	source := types.SourceFilter{
		Includes: []string{"title"},
		Excludes: []string{"tags"},
	}

	req := (query.SearchParams{
		Size:           0,
		HasSize:        true,
		From:           0,
		HasFrom:        true,
		SearchAfter:    []types.FieldValue{"cursor-1", 42},
		TrackTotalHits: false,
		Source:         source,
		Timeout:        &timeout,
	}).ToRequest()

	assert.Assert(t, req != nil)
	assert.Assert(t, req.Size != nil)
	assert.Equal(t, 0, *req.Size)
	assert.Assert(t, req.From != nil)
	assert.Equal(t, 0, *req.From)
	assert.DeepEqual(t, []types.FieldValue{"cursor-1", 42}, req.SearchAfter)
	assert.Assert(t, req.TrackTotalHits != nil)
	assert.Equal(t, false, req.TrackTotalHits)
	assert.DeepEqual(t, source, req.Source_)
	assert.Assert(t, req.Timeout != nil)
	assert.Equal(t, "30s", *req.Timeout)
}

func TestSearchParams_ToV9Request_AdvancedOptions(t *testing.T) {
	t.Parallel()

	timeout := "30s"
	req := (query.SearchParams{
		Size:           0,
		HasSize:        true,
		From:           0,
		HasFrom:        true,
		SearchAfter:    []types.FieldValue{"cursor-1", "cursor-2"},
		TrackTotalHits: false,
		Source: types.SourceFilter{
			Includes: []string{"title"},
			Excludes: []string{"tags"},
		},
		Timeout: &timeout,
	}).ToV9Request()

	assert.Assert(t, req != nil)
	assert.Assert(t, req.Size != nil)
	assert.Equal(t, 0, *req.Size)
	assert.Assert(t, req.From != nil)
	assert.Equal(t, 0, *req.From)
	assert.DeepEqual(t, []v9types.FieldValue{"cursor-1", "cursor-2"}, req.SearchAfter)
	assert.Assert(t, req.TrackTotalHits != nil)
	assert.Equal(t, false, req.TrackTotalHits)
	rawSource, ok := req.Source_.(*v9types.SourceFilter)
	assert.Assert(t, ok)
	assert.DeepEqual(t, &v9types.SourceFilter{
		Includes: []string{"title"},
		Excludes: []string{"tags"},
	}, rawSource)
	assert.Assert(t, req.Timeout != nil)
	assert.Equal(t, "30s", *req.Timeout)
}

func TestIsZeroQuery(t *testing.T) {
	t.Parallel()

	assert.Assert(t, query.IsZeroQuery(types.Query{}))
	assert.Assert(t, !query.IsZeroQuery(types.Query{MatchAll: &types.MatchAllQuery{}}))
	assert.Assert(t, !query.IsZeroQuery(func() {}))
}

func TestSearchBuilder_Aggregation(t *testing.T) {
	t.Parallel()
	aggs := query.Aggs(
		query.StringTermsAgg("by_category", FieldCategory),
	).Build()

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
			FieldTitle.String(): {},
		},
	}
	params := query.NewSearch().
		Highlight(h).
		Build()

	assert.Assert(t, params.Highlight != nil)
	_, ok := params.Highlight.Fields[FieldTitle.String()]
	assert.Assert(t, ok)
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
			FieldValue.String(): {Script: types.Script{Source: &source}},
		}).
		Build()

	assert.Assert(t, params.ScriptFields != nil)
	sf, ok := params.ScriptFields[FieldValue.String()]
	assert.Assert(t, ok)
	assert.Equal(t, source, *sf.Script.Source)
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
	// Demonstrates the ActiveRecord-style fluent API.
	aggs := query.Aggs(
		query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month),
	).Build()

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
			Fields: map[string]types.HighlightField{FieldTitle.String(): {}},
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
