package esv8

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv8/query"
)

// SearchParams defines the high-level search input for Search[T].
// It mirrors the commonly used parts of search.Request while keeping the
// library's typed-search workflow focused on the most common application needs.
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
// It is primarily used internally by Search[T], SearchDocuments[T], and
// SearchOne[T] to call SearchRaw.
func (p SearchParams) ToRequest() *search.Request {
	req := search.NewRequest()

	req.Query = &p.Query

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

// SearchResponse is the high-level response returned by Search[T].
// It contains decoded hits, typed aggregation access, and the underlying raw
// Elasticsearch response for escape-hatch scenarios.
type SearchResponse[T any] struct {
	Total        int64
	Hits         []SearchHit[T]
	Aggregations query.AggResults
	Raw          *search.Response
}

// SearchClient is the minimal capability required by Search[T].
// Any client that can execute SearchRaw can participate in the high-level
// typed-search helper flow.
type SearchClient interface {
	SearchRaw(ctx context.Context, aliasName estype.Alias, req *search.Request) (*search.Response, error)
}

// Search executes the preferred high-level search flow against an alias.
// It converts SearchParams into a search.Request, executes SearchRaw, decodes
// each hit's _source into T, and exposes typed aggregation accessors through
// SearchResponse[T].
func Search[T any](
	ctx context.Context,
	client SearchClient,
	aliasName estype.Alias,
	params SearchParams,
) (*SearchResponse[T], error) {
	rawResp, err := client.SearchRaw(ctx, aliasName, params.ToRequest())
	if err != nil {
		return nil, err
	}

	resp := &SearchResponse[T]{
		Raw:          rawResp,
		Aggregations: query.NewAggResults(rawResp.Aggregations),
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

// SearchDocuments executes Search[T] and returns only the decoded sources.
// Use it when hit metadata and aggregations are not needed.
func SearchDocuments[T any](
	ctx context.Context,
	client SearchClient,
	aliasName estype.Alias,
	params SearchParams,
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

// SearchOne executes Search[T] with Size forced to 1 and returns the first
// decoded source plus a boolean indicating whether a hit was found.
func SearchOne[T any](
	ctx context.Context,
	client SearchClient,
	aliasName estype.Alias,
	params SearchParams,
) (T, bool, error) {
	params.Size = 1

	resp, err := Search[T](ctx, client, aliasName, params)
	if err != nil {
		var zero T
		return zero, false, err
	}
	if len(resp.Hits) == 0 {
		var zero T
		return zero, false, nil
	}

	return resp.Hits[0].Source, true, nil
}
