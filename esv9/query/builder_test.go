package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/esv9/query"
)

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

