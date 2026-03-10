package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomtwinkle/es-typed-go/query"
)

func TestBuilder_Build_Empty(t *testing.T) {
	t.Parallel()
	q := query.New().Build()
	assert.Equal(t, types.Query{}, q)
}

func TestBuilder_MatchAll(t *testing.T) {
	t.Parallel()
	q := query.New().MatchAll(&types.MatchAllQuery{}).Build()
	require.NotNil(t, q.MatchAll)
}

func TestBuilder_Bool(t *testing.T) {
	t.Parallel()
	bq := query.NewBoolQuery().
		Must(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Build()

	q := query.New().Bool(bq).Build()
	require.NotNil(t, q.Bool)
	assert.Len(t, q.Bool.Must, 1)
}

func TestBuilder_Term(t *testing.T) {
	t.Parallel()
	val := "foo"
	q := query.New().Term("status", types.TermQuery{Value: val}).Build()
	require.NotNil(t, q.Term)
	assert.Equal(t, val, q.Term["status"].Value)
}

func TestBuilder_Match(t *testing.T) {
	t.Parallel()
	q := query.New().Match("title", types.MatchQuery{Query: "hello"}).Build()
	require.NotNil(t, q.Match)
	assert.Equal(t, "hello", q.Match["title"].Query)
}

func TestBuilder_Range(t *testing.T) {
	t.Parallel()
	gt := "2023-01-01"
	rq := types.NewDateRangeQuery()
	rq.Gte = &gt
	q := query.New().Range("created_at", rq).Build()
	require.NotNil(t, q.Range)
	assert.NotNil(t, q.Range["created_at"])
}

func TestBoolQueryBuilder_ShouldAndFilter(t *testing.T) {
	t.Parallel()
	bq := query.NewBoolQuery().
		Should(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Filter(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Build()

	assert.Len(t, bq.Should, 1)
	assert.Len(t, bq.Filter, 1)
}

func TestBoolQueryBuilder_MustNot(t *testing.T) {
	t.Parallel()
	bq := query.NewBoolQuery().
		MustNot(types.Query{MatchAll: &types.MatchAllQuery{}}).
		Build()

	assert.Len(t, bq.MustNot, 1)
}
