package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/query"
)

func TestTermValue(t *testing.T) {
	t.Parallel()
	q := query.TermValue(FieldStatus, "active")
	assert.Assert(t, q.Term != nil)
	assert.Equal(t, "active", q.Term[string(FieldStatus)].Value)
}

func TestTermsValues(t *testing.T) {
	t.Parallel()
	q := query.TermsValues(FieldCategory,
		types.FieldValue("val1"),
		types.FieldValue("val2"),
		types.FieldValue("val3"),
	)
	assert.Assert(t, q.Terms != nil)
	assert.Assert(t, q.Terms.TermsQuery != nil)
	rawVals, ok := q.Terms.TermsQuery[string(FieldCategory)]
	assert.Assert(t, ok)
	vals, ok := rawVals.([]types.FieldValue)
	assert.Assert(t, ok)
	assert.Assert(t, len(vals) == 3)
}

func TestMatchPhrase(t *testing.T) {
	t.Parallel()
	q := query.MatchPhrase(FieldTitleNgram, "test keyword")
	assert.Assert(t, q.MatchPhrase != nil)
	assert.Equal(t, "test keyword", q.MatchPhrase[string(FieldTitleNgram)].Query)
}

func TestExistsField(t *testing.T) {
	t.Parallel()
	q := query.ExistsField(FieldTags)
	assert.Assert(t, q.Exists != nil)
	assert.Equal(t, string(FieldTags), q.Exists.Field)
}

func TestNotExists(t *testing.T) {
	t.Parallel()
	q := query.NotExists(FieldTags)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.MustNot) == 1)
	assert.Assert(t, q.Bool.MustNot[0].Exists != nil)
	assert.Equal(t, string(FieldTags), q.Bool.MustNot[0].Exists.Field)
}

func TestNestedFilter(t *testing.T) {
	t.Parallel()
	inner := query.TermValue(FieldItemsColor, "red")
	q := query.NestedFilter(FieldItems, inner)
	assert.Assert(t, q.Nested != nil)
	assert.Equal(t, string(FieldItems), q.Nested.Path)
	assert.Assert(t, q.Nested.Query.Bool != nil)
	assert.Assert(t, len(q.Nested.Query.Bool.Filter) == 1)
}

func TestNestedFilter_MultipleQueries(t *testing.T) {
	t.Parallel()
	q1 := query.TermValue(FieldItemsColor, "red")
	q2 := query.TermValue(FieldItemsStatus, "active")
	q := query.NestedFilter(FieldItems, q1, q2)
	assert.Assert(t, q.Nested != nil)
	assert.Assert(t, len(q.Nested.Query.Bool.Filter) == 2)
}

func TestDateRangeQuery(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery(FieldDate, query.DateRangeGte("2024-01-01"), query.DateRangeLte("2024-12-31"))
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range[string(FieldDate)] != nil)
}

func TestDateRangeQuery_GteOnly(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery(FieldDate, query.DateRangeGte("2024-01-01"))
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range[string(FieldDate)] != nil)
}

func TestDateRangeQuery_LteOnly(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery(FieldDate, query.DateRangeLte("2024-12-31"))
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range[string(FieldDate)] != nil)
}

func TestDateRangeQuery_GtLt(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery(FieldDate, query.DateRangeGt("2024-01-01"), query.DateRangeLt("2025-01-01"))
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range[string(FieldDate)] != nil)
}

func TestDateRangeQuery_NoOpts(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery(FieldDate)
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range[string(FieldDate)] != nil)
}

func TestNumberRangeQuery(t *testing.T) {
	t.Parallel()
	gte := types.Float64(10.0)
	lte := types.Float64(100.0)
	q := query.NumberRangeQuery(FieldPrice, &gte, &lte)
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range[string(FieldPrice)] != nil)
}

func TestNumberRangeQuery_GteOnly(t *testing.T) {
	t.Parallel()
	gte := types.Float64(10.0)
	q := query.NumberRangeQuery(FieldPrice, &gte, nil)
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range[string(FieldPrice)] != nil)
}

func TestNumberRangeQuery_LteOnly(t *testing.T) {
	t.Parallel()
	lte := types.Float64(100.0)
	q := query.NumberRangeQuery(FieldPrice, nil, &lte)
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range[string(FieldPrice)] != nil)
}

func TestBoolMust(t *testing.T) {
	t.Parallel()
	q := query.BoolMust(
		query.TermValue(FieldStatus, "active"),
		query.TermValue(FieldEnabled, true),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Must) == 2)
}

func TestBoolShould(t *testing.T) {
	t.Parallel()
	q := query.BoolShould(
		query.MatchPhrase(FieldTitleNgram, "keyword"),
		query.MatchPhrase(FieldTitleRaw, "keyword"),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Should) == 2)
}

func TestBoolFilter(t *testing.T) {
	t.Parallel()
	q := query.BoolFilter(
		query.TermValue(FieldType, "document"),
		query.TermValue(FieldStatus, "active"),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Filter) == 2)
}

func TestBoolMustNot(t *testing.T) {
	t.Parallel()
	q := query.BoolMustNot(
		query.ExistsField(FieldDate),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.MustNot) == 1)
}

func TestFieldValues_Strings(t *testing.T) {
	t.Parallel()
	vals := query.FieldValues("a", "b", "c")
	assert.Assert(t, len(vals) == 3)
	assert.Equal(t, "a", vals[0])
	assert.Equal(t, "b", vals[1])
	assert.Equal(t, "c", vals[2])
}

func TestFieldValues_Int32(t *testing.T) {
	t.Parallel()
	vals := query.FieldValues(int32(1), int32(2), int32(3))
	assert.Assert(t, len(vals) == 3)
	assert.Equal(t, int32(1), vals[0])
}

func TestFieldValues_Empty(t *testing.T) {
	t.Parallel()
	vals := query.FieldValues[string]()
	assert.Assert(t, len(vals) == 0)
}

func TestFieldValues_WithTermsValues(t *testing.T) {
	t.Parallel()
	// Test that FieldValues output is compatible with TermsValues.
	q := query.TermsValues(FieldTags, query.FieldValues(int32(1), int32(2), int32(3))...)
	assert.Assert(t, q.Terms != nil)
	rawVals, ok := q.Terms.TermsQuery[string(FieldTags)]
	assert.Assert(t, ok)
	vals, ok := rawVals.([]types.FieldValue)
	assert.Assert(t, ok)
	assert.Assert(t, len(vals) == 3)
}

func TestComplexQueryCombination(t *testing.T) {
	t.Parallel()
	// Demonstrates building a complex query using the helper functions
	// instead of verbose struct construction.
	filters := []types.Query{
		query.TermValue(FieldType, "document"),
		query.TermsValues(FieldCategory, query.FieldValues("cat1", "cat2")...),
		query.TermValue(FieldStatus, "active"),
		query.MatchPhrase(FieldTitleNgram, "search keyword"),
		query.NotExists(FieldTags),
		query.NestedFilter(FieldItems,
			query.TermsValues(FieldItemsIds, query.FieldValues(int32(1), int32(2))...),
		),
	}

	q := query.BoolQuery(
		query.NewBoolQuery().Filter(filters...).Build(),
	)

	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Filter) == 6)
}

func TestMatchValue(t *testing.T) {
	t.Parallel()
	q := query.MatchValue(FieldTitle, "hello")
	assert.Assert(t, q.Match != nil)
	assert.Equal(t, "hello", q.Match[string(FieldTitle)].Query)
}

func TestMatchAll(t *testing.T) {
	t.Parallel()
	q := query.MatchAll()
	assert.Assert(t, q.MatchAll != nil)
}

func TestMatchNone(t *testing.T) {
	t.Parallel()
	q := query.MatchNone()
	assert.Assert(t, q.MatchNone != nil)
}

func TestIdsQuery(t *testing.T) {
	t.Parallel()
	q := query.IdsQuery("id1", "id2", "id3")
	assert.Assert(t, q.Ids != nil)
	assert.Assert(t, len(q.Ids.Values) == 3)
	assert.Equal(t, "id1", q.Ids.Values[0])
}

func TestPrefixValue(t *testing.T) {
	t.Parallel()
	q := query.PrefixValue(FieldName, "pre")
	assert.Assert(t, q.Prefix != nil)
	assert.Equal(t, "pre", q.Prefix[string(FieldName)].Value)
}

func TestWildcardValue(t *testing.T) {
	t.Parallel()
	q := query.WildcardValue(FieldName, "val*")
	assert.Assert(t, q.Wildcard != nil)
	assert.Equal(t, "val*", *q.Wildcard[string(FieldName)].Value)
}

func TestMultiMatchQuery(t *testing.T) {
	t.Parallel()
	q := query.MultiMatchQuery("search text", FieldTitle, FieldName)
	assert.Assert(t, q.MultiMatch != nil)
	assert.Equal(t, "search text", q.MultiMatch.Query)
	assert.Assert(t, len(q.MultiMatch.Fields) == 2)
}

func TestFunctionScoreQuery(t *testing.T) {
	t.Parallel()
	q := query.FunctionScoreQuery(&types.FunctionScoreQuery{
		Query: &types.Query{MatchAll: &types.MatchAllQuery{}},
	})
	assert.Assert(t, q.FunctionScore != nil)
	assert.Assert(t, q.FunctionScore.Query != nil)
	assert.Assert(t, q.FunctionScore.Query.MatchAll != nil)
}

func TestBoolQuery(t *testing.T) {
	t.Parallel()
	bq := query.NewBoolQuery().
		Must(query.TermValue(FieldStatus, "active")).
		Build()
	q := query.BoolQuery(bq)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Must) == 1)
}
