package esv8

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	searchapi "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv8/query"
)

func TestSearchParams_ToRequest_DefaultSizeUsesElasticsearchDefault(t *testing.T) {
	t.Parallel()

	req := (SearchParams{}).ToRequest()

	assert.Assert(t, req != nil)
	assert.Assert(t, req.Size == nil)
	assert.Equal(t, req.Source_, true)
	assert.Assert(t, req.Timeout != nil)
	assert.Equal(t, *req.Timeout, "10s")
}

func TestSearchParams_ToRequest_ExplicitSizeIsApplied(t *testing.T) {
	t.Parallel()

	req := (SearchParams{
		Size: 25,
	}).ToRequest()

	assert.Assert(t, req != nil)
	assert.Assert(t, req.Size != nil)
	assert.Equal(t, *req.Size, 25)
	assert.Equal(t, req.Source_, true)
	assert.Assert(t, req.Timeout != nil)
	assert.Equal(t, *req.Timeout, "10s")
}

func TestSearchHelpers_WithDirectParams(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		params func() SearchRequest
		size   int
	}{
		{
			name: "package_search_params",
			params: func() SearchRequest {
				return SearchParams{
					Query: query.TermValue(estype.Field("status"), "active"),
					Size:  2,
				}
			},
			size: 2,
		},
		{
			name: "query_builder_search_params",
			params: func() SearchRequest {
				return query.NewSearch().
					Where(query.TermValue(estype.Field("status"), "active")).
					Limit(2).
					Build()
			},
			size: 2,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			t.Run("search", func(t *testing.T) {
				t.Parallel()

				client := searchClientFunc(func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
					assert.Equal(t, estype.Alias("products"), aliasName)
					assert.Assert(t, req != nil)
					assert.Assert(t, req.Query != nil)

					if req.Query.Bool != nil {
						assert.Assert(t, len(req.Query.Bool.Filter) == 1)
					} else {
						assert.Assert(t, req.Query.Term != nil)
						assert.Equal(t, "active", req.Query.Term["status"].Value)
					}

					assert.Assert(t, req.Size != nil)
					assert.Equal(t, tc.size, *req.Size)

					src, err := json.Marshal(searchTestDoc{
						ID:    "doc-1",
						Name:  "Widget",
						Price: 42,
					})
					assert.NilError(t, err)

					id := "doc-1"
					score := types.Float64(1.5)

					return &searchapi.Response{
						Hits: types.HitsMetadata{
							Total: &types.TotalHits{Value: 1},
							Hits: []types.Hit{
								{
									Id_:     &id,
									Index_:  "products-000001",
									Score_:  &score,
									Source_: src,
								},
							},
						},
					}, nil
				})

				resp, err := Search[searchTestDoc](context.Background(), client, estype.Alias("products"), tc.params())
				assert.NilError(t, err)
				assert.Assert(t, resp != nil)
				assert.Equal(t, int64(1), resp.Total)
				assert.Equal(t, 1, len(resp.Hits))
				assert.Equal(t, "doc-1", resp.Hits[0].ID)
				assert.Equal(t, "Widget", resp.Hits[0].Source.Name)
				assert.Assert(t, resp.Hits[0].Score != nil)
				assert.Equal(t, 1.5, *resp.Hits[0].Score)
			})

			t.Run("search_documents", func(t *testing.T) {
				t.Parallel()

				client := searchClientFunc(func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
					assert.Assert(t, req != nil)
					assert.Assert(t, req.Size != nil)
					assert.Equal(t, tc.size, *req.Size)

					src1, err := json.Marshal(searchTestDoc{ID: "doc-1", Name: "Alpha", Price: 10})
					assert.NilError(t, err)
					src2, err := json.Marshal(searchTestDoc{ID: "doc-2", Name: "Beta", Price: 20})
					assert.NilError(t, err)

					id1 := "doc-1"
					id2 := "doc-2"

					return &searchapi.Response{
						Hits: types.HitsMetadata{
							Total: &types.TotalHits{Value: 2},
							Hits: []types.Hit{
								{Id_: &id1, Index_: "products-000001", Source_: src1},
								{Id_: &id2, Index_: "products-000001", Source_: src2},
							},
						},
					}, nil
				})

				docs, err := SearchDocuments[searchTestDoc](context.Background(), client, estype.Alias("products"), tc.params())
				assert.NilError(t, err)
				assert.Equal(t, 2, len(docs))
				assert.Equal(t, "Alpha", docs[0].Name)
				assert.Equal(t, "Beta", docs[1].Name)
			})
		})
	}
}

func TestSearchParams_ToRequest_AllOptionalFields(t *testing.T) {
	t.Parallel()

	highlight := &types.Highlight{
		Fields: map[string]types.HighlightField{
			"title": {},
		},
	}
	collapse := &types.FieldCollapse{
		Field: "category",
	}
	scriptFields := map[string]types.ScriptField{
		"score_x2": {},
	}
	aggs := map[string]types.Aggregations{
		"avg_price": {
			Avg: types.NewAverageAggregation(),
		},
	}
	sorts := []types.SortCombinations{
		types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"date": {},
			},
		},
	}

	req := (SearchParams{
		Query:        query.TermValue(estype.Field("status"), "active"),
		Sort:         sorts,
		Aggregations: aggs,
		Highlight:    highlight,
		Collapse:     collapse,
		ScriptFields: scriptFields,
		Size:         3,
		From:         7,
	}).ToRequest()

	assert.Assert(t, req != nil)
	assert.Assert(t, req.Query != nil)
	assert.DeepEqual(t, sorts, req.Sort)
	assert.DeepEqual(t, aggs, req.Aggregations)
	assert.Equal(t, highlight, req.Highlight)
	assert.Equal(t, collapse, req.Collapse)
	assert.DeepEqual(t, scriptFields, req.ScriptFields)
	assert.Assert(t, req.Size != nil)
	assert.Equal(t, 3, *req.Size)
	assert.Assert(t, req.From != nil)
	assert.Equal(t, 7, *req.From)
	assert.Equal(t, true, req.Source_)
	assert.Assert(t, req.Timeout != nil)
	assert.Equal(t, "10s", *req.Timeout)
}

func TestSearch_EmptyHitsAndNilTotal(t *testing.T) {
	t.Parallel()

	client := searchClientFunc(func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
		assert.Equal(t, estype.Alias("products"), aliasName)
		assert.Assert(t, req != nil)

		return &searchapi.Response{
			Aggregations: map[string]types.Aggregate{},
			Hits: types.HitsMetadata{
				Hits: nil,
			},
		}, nil
	})

	resp, err := Search[searchTestDoc](context.Background(), client, estype.Alias("products"), SearchParams{})
	assert.NilError(t, err)
	assert.Assert(t, resp != nil)
	assert.Equal(t, int64(0), resp.Total)
	assert.Equal(t, 0, len(resp.Hits))
	assert.Assert(t, resp.Raw != nil)
}

func TestSearch_DecodeErrorIncludesHitID(t *testing.T) {
	t.Parallel()

	client := searchClientFunc(func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
		id := "broken-doc"
		return &searchapi.Response{
			Hits: types.HitsMetadata{
				Hits: []types.Hit{
					{
						Id_:     &id,
						Index_:  "products-000001",
						Source_: []byte("{"),
					},
				},
			},
		}, nil
	})

	resp, err := Search[searchTestDoc](context.Background(), client, estype.Alias("products"), SearchParams{})
	assert.Assert(t, resp == nil)
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), `decode search hit "broken-doc"`))
}

func TestSearch_DecodeErrorWithoutHitID(t *testing.T) {
	t.Parallel()

	client := searchClientFunc(func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
		return &searchapi.Response{
			Hits: types.HitsMetadata{
				Hits: []types.Hit{
					{
						Index_:  "products-000001",
						Source_: []byte("{"),
					},
				},
			},
		}, nil
	})

	resp, err := Search[searchTestDoc](context.Background(), client, estype.Alias("products"), SearchParams{})
	assert.Assert(t, resp == nil)
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), `decode search hit ""`))
}

func TestSearch_AllowsHitWithoutSourceIDOrScore(t *testing.T) {
	t.Parallel()

	client := searchClientFunc(func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
		return &searchapi.Response{
			Hits: types.HitsMetadata{
				Total: &types.TotalHits{Value: 1},
				Hits: []types.Hit{
					{
						Index_: "products-000001",
					},
				},
			},
		}, nil
	})

	resp, err := Search[searchTestDoc](context.Background(), client, estype.Alias("products"), SearchParams{})
	assert.NilError(t, err)
	assert.Assert(t, resp != nil)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, 1, len(resp.Hits))
	assert.Equal(t, "", resp.Hits[0].ID)
	assert.Equal(t, "products-000001", resp.Hits[0].Index)
	assert.Assert(t, resp.Hits[0].Score == nil)
	assert.DeepEqual(t, searchTestDoc{}, resp.Hits[0].Source)
}

func TestSearch_PropagatesSearchRawError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("search failed")
	client := searchClientFunc(func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
		return nil, wantErr
	})

	resp, err := Search[searchTestDoc](context.Background(), client, estype.Alias("products"), SearchParams{})
	assert.Assert(t, resp == nil)
	assert.Equal(t, wantErr, err)
}

func TestSearchDocuments_PropagatesSearchError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("search documents failed")
	client := searchClientFunc(func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
		return nil, wantErr
	})

	docs, err := SearchDocuments[searchTestDoc](context.Background(), client, estype.Alias("products"), SearchParams{})
	assert.Assert(t, docs == nil)
	assert.Equal(t, wantErr, err)
}

func TestSearchDocuments_EmptyResult(t *testing.T) {
	t.Parallel()

	client := searchClientFunc(func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
		return &searchapi.Response{
			Hits: types.HitsMetadata{},
		}, nil
	})

	docs, err := SearchDocuments[searchTestDoc](context.Background(), client, estype.Alias("products"), SearchParams{})
	assert.NilError(t, err)
	assert.Assert(t, docs != nil)
	assert.Equal(t, 0, len(docs))
}

type searchTestDoc struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type searchClientFunc func(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error)

func (f searchClientFunc) SearchRaw(ctx context.Context, aliasName estype.Alias, req *searchapi.Request) (*searchapi.Response, error) {
	return f(ctx, aliasName, req)
}
