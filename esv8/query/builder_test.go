package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv8/query"
)

func TestBuilder_Build_Empty(t *testing.T) {
	t.Parallel()
	q := query.New().Build()
	assert.DeepEqual(t, types.Query{}, q)
}

func TestBuilder_MatchAll(t *testing.T) {
	t.Parallel()
	q := query.New().MatchAll(&types.MatchAllQuery{}).Build()
	assert.Assert(t, q.MatchAll != nil)
}

func TestBuilder_Bool(t *testing.T) {
	t.Parallel()
	bq := query.NewBoolQuery().
		Must(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Build()

	q := query.New().Bool(bq).Build()
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Must) == 1)
}

func TestBuilder_Term(t *testing.T) {
	t.Parallel()
	val := "foo"
	q := query.New().Term(FieldStatus, types.TermQuery{Value: val}).Build()
	assert.Assert(t, q.Term != nil)
	assert.Equal(t, val, q.Term[string(FieldStatus)].Value)
}

func TestBuilder_Match(t *testing.T) {
	t.Parallel()
	q := query.New().Match(FieldTitle, types.MatchQuery{Query: "hello"}).Build()
	assert.Assert(t, q.Match != nil)
	assert.Equal(t, "hello", q.Match[string(FieldTitle)].Query)
}

func TestBuilder_Range(t *testing.T) {
	t.Parallel()
	gt := "2023-01-01"
	rq := types.NewDateRangeQuery()
	rq.Gte = &gt
	q := query.New().Range(FieldDate, rq).Build()
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range[string(FieldDate)] != nil)
}

func TestBoolQueryBuilder_ShouldAndFilter(t *testing.T) {
	t.Parallel()
	bq := query.NewBoolQuery().
		Should(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Filter(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Build()

	assert.Assert(t, len(bq.Should) == 1)
	assert.Assert(t, len(bq.Filter) == 1)
}

func TestBoolQueryBuilder_MustNot(t *testing.T) {
	t.Parallel()
	bq := query.NewBoolQuery().
		MustNot(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Build()

	assert.Assert(t, len(bq.MustNot) == 1)
}

func TestBuilder_Terms(t *testing.T) {
	t.Parallel()
	q := query.New().Terms(&types.TermsQuery{
		TermsQuery: map[string]types.TermsQueryField{
			string(FieldTags): []types.FieldValue{"a", "b"},
		},
	}).Build()
	assert.Assert(t, q.Terms != nil)
	vals, ok := q.Terms.TermsQuery[string(FieldTags)].([]types.FieldValue)
	assert.Assert(t, ok)
	assert.Assert(t, len(vals) == 2)
}

func TestBuilder_Exists(t *testing.T) {
	t.Parallel()
	q := query.New().Exists(&types.ExistsQuery{Field: string(FieldStatus)}).Build()
	assert.Assert(t, q.Exists != nil)
	assert.Equal(t, string(FieldStatus), q.Exists.Field)
}

func TestBuilder_MatchNone(t *testing.T) {
	t.Parallel()
	q := query.New().MatchNone(&types.MatchNoneQuery{}).Build()
	assert.Assert(t, q.MatchNone != nil)
}

func TestBuilder_Ids(t *testing.T) {
	t.Parallel()
	q := query.New().Ids(&types.IdsQuery{Values: []string{"id1", "id2"}}).Build()
	assert.Assert(t, q.Ids != nil)
	assert.Assert(t, len(q.Ids.Values) == 2)
}

func TestBuilder_Prefix(t *testing.T) {
	t.Parallel()
	q := query.New().Prefix(FieldName, types.PrefixQuery{Value: "pre"}).Build()
	assert.Assert(t, q.Prefix != nil)
	assert.Equal(t, "pre", q.Prefix[string(FieldName)].Value)
}

func TestBuilder_Wildcard(t *testing.T) {
	t.Parallel()
	q := query.New().Wildcard(FieldName, types.WildcardQuery{Value: strPtr("val*")}).Build()
	assert.Assert(t, q.Wildcard != nil)
	assert.Equal(t, "val*", *q.Wildcard[string(FieldName)].Value)
}

func TestBuilder_MultiMatch(t *testing.T) {
	t.Parallel()
	q := query.New().MultiMatch(&types.MultiMatchQuery{
		Query:  "search text",
		Fields: estype.FieldNames(FieldTitle, FieldName),
	}).Build()
	assert.Assert(t, q.MultiMatch != nil)
	assert.Equal(t, "search text", q.MultiMatch.Query)
	assert.Assert(t, len(q.MultiMatch.Fields) == 2)
}

func TestBuilder_FunctionScore(t *testing.T) {
	t.Parallel()
	q := query.New().FunctionScore(&types.FunctionScoreQuery{
		Query: &types.Query{MatchAll: &types.MatchAllQuery{}},
	}).Build()
	assert.Assert(t, q.FunctionScore != nil)
	assert.Assert(t, q.FunctionScore.Query != nil)
	assert.Assert(t, q.FunctionScore.Query.MatchAll != nil)
}

func TestBoolQueryBuilder_MinimumShouldMatch(t *testing.T) {
	t.Parallel()
	bq := query.NewBoolQuery().
		Should(
			types.Query{MatchAll: &types.MatchAllQuery{}},
			types.Query{MatchAll: &types.MatchAllQuery{}},
		).
		MinimumShouldMatch(1).
		Build()

	assert.Assert(t, len(bq.Should) == 2)
	assert.Equal(t, 1, bq.MinimumShouldMatch)
}

func strPtr(s string) *string { return &s }
