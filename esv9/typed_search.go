package esv9

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/query"
)

// SearchParams defines the high-level search input surface.
// It mirrors the commonly used parts of search.Request while staying focused on
// the library's typed search workflow.
type SearchParams struct {
	Query        types.Query
	Sort         []types.SortCombinations
	Aggregations map[string]types.Aggregations
	Highlight    *types.Highlight
	Collapse     *types.FieldCollapse
	ScriptFields map[string]types.ScriptField
	Size         int
	From         int
}

// ToRequest converts SearchParams into a typed Elasticsearch search.Request.
// Deprecated: Use ToV9Request instead for clarity.
func (p SearchParams) ToRequest() *search.Request {
	return p.ToV9Request()
}

// ToV9Request converts SearchParams into a typed Elasticsearch v9 search.Request.
func (p SearchParams) ToV9Request() *search.Request {
	req := search.NewRequest()

	if !query.IsZeroQuery(p.Query) {
		req.Query = &p.Query
	}

	if len(p.Sort) > 0 {
		req.Sort = p.Sort
	}
	if len(p.Aggregations) > 0 {
		req.Aggregations = p.Aggregations
	}
	if p.Highlight != nil {
		req.Highlight = p.Highlight
	}
	if p.Collapse != nil {
		req.Collapse = p.Collapse
	}
	if len(p.ScriptFields) > 0 {
		req.ScriptFields = p.ScriptFields
	}

	if p.Size > 0 {
		size := p.Size
		req.Size = &size
	}

	if p.From > 0 {
		from := p.From
		req.From = &from
	}

	timeout := "10s"
	req.Timeout = &timeout
	req.Source_ = true

	return req
}

// SearchHit is a typed view of a single search hit.
type SearchHit[T any] struct {
	ID     string
	Index  string
	Score  *float64
	Source T
	Raw    types.Hit
}

// SearchResponse is the high-level search response.
// It exposes decoded hits, typed aggregation accessors, and the raw response.
type SearchResponse[T any] struct {
	Total        int64
	Hits         []SearchHit[T]
	Aggregations query.AggResults
	Raw          *search.Response
}

// SearchClient is the minimal capability required by Search.
// Any client that can execute SearchRaw can participate in the high-level search flow.
type SearchClient interface {
	SearchRaw(ctx context.Context, aliasName estype.Alias, req *search.Request) (*search.Response, error)
}

// SearchRequest is the minimal input accepted by Search[T] and related helpers.
// esv9.SearchParams, esv9/query.SearchParams, and query.SearchParams all satisfy
// this interface.
type SearchRequest interface {
	ToV9Request() *search.Request
}

// Search executes a high-level search against an alias and decodes each hit's
// _source into T.
// Prefer this over SearchRaw for normal application code.
func Search[T any](
	ctx context.Context,
	client SearchClient,
	aliasName estype.Alias,
	params SearchRequest,
) (*SearchResponse[T], error) {
	rawResp, err := client.SearchRaw(ctx, aliasName, params.ToV9Request())
	if err != nil {
		return nil, err
	}

	resp := &SearchResponse[T]{
		Raw:          rawResp,
		Aggregations: aggResultsFromV9(rawResp.Aggregations),
	}

	if rawResp.Hits.Total != nil {
		resp.Total = rawResp.Hits.Total.Value
	}

	if len(rawResp.Hits.Hits) == 0 {
		return resp, nil
	}

	resp.Hits = make([]SearchHit[T], 0, len(rawResp.Hits.Hits))
	for _, hit := range rawResp.Hits.Hits {
		var src T
		if len(hit.Source_) > 0 {
			if err := json.Unmarshal(hit.Source_, &src); err != nil {
				var hitID string
				if hit.Id_ != nil {
					hitID = *hit.Id_
				}
				return nil, fmt.Errorf("decode search hit %q: %w", hitID, err)
			}
		}

		var score *float64
		if hit.Score_ != nil {
			v := float64(*hit.Score_)
			score = &v
		}

		var id string
		if hit.Id_ != nil {
			id = *hit.Id_
		}

		resp.Hits = append(resp.Hits, SearchHit[T]{
			ID:     id,
			Index:  hit.Index_,
			Score:  score,
			Source: src,
			Raw:    hit,
		})
	}

	return resp, nil
}

// SearchDocuments executes Search and returns only the decoded document sources.
func SearchDocuments[T any](
	ctx context.Context,
	client SearchClient,
	aliasName estype.Alias,
	params SearchRequest,
) ([]T, error) {
	resp, err := Search[T](ctx, client, aliasName, params)
	if err != nil {
		return nil, err
	}

	docs := make([]T, 0, len(resp.Hits))
	for _, hit := range resp.Hits {
		docs = append(docs, hit.Source)
	}
	return docs, nil
}

// aggResultsFromV9 converts v9 aggregation results to query.AggResults via JSON
// round-trip, since the v8 and v9 aggregation JSON formats are identical.
func aggResultsFromV9(raw map[string]types.Aggregate) query.AggResults {
	if len(raw) == 0 {
		return query.NewAggResultsFromJSON(nil)
	}
	b, _ := json.Marshal(raw)
	return query.NewAggResultsFromJSON(b)
}
