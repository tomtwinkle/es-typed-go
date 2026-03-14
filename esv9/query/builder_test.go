package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv9/query"
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
	q := query.New().Term(estype.Field("status"), types.TermQuery{Value: val}).Build()
	assert.Assert(t, q.Term != nil)
	assert.Equal(t, val, q.Term["status"].Value)
}

func TestBuilder_Match(t *testing.T) {
	t.Parallel()
	q := query.New().Match(estype.Field("title"), types.MatchQuery{Query: "hello"}).Build()
	assert.Assert(t, q.Match != nil)
	assert.Equal(t, "hello", q.Match["title"].Query)
}

func TestBuilder_Range(t *testing.T) {
	t.Parallel()
	gt := "2023-01-01"
	rq := types.NewDateRangeQuery()
	rq.Gte = &gt
	q := query.New().Range(estype.Field("date"), rq).Build()
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["date"] != nil)
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
