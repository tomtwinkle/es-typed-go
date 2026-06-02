package esv8

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/query"
)

// SearchParams defines the high-level search input for Search[T].
// It mirrors the commonly used parts of search.Request while keeping the
// library's typed-search workflow focused on the most common application needs.
type SearchParams struct {
	Query          types.Query
	Sort           []types.SortCombinations
	Aggregations   map[string]types.Aggregations
	Highlight      *types.Highlight
	Collapse       *types.FieldCollapse
	ScriptFields   map[string]types.ScriptField
	Size           int
	HasSize        bool
	From           int
	HasFrom        bool
	SearchAfter    []types.FieldValue
	TrackTotalHits types.TrackHits
	Source         types.SourceConfig
	Timeout        *string
}

// ToRequest converts SearchParams into a typed Elasticsearch search.Request.
// It is primarily used internally by Search[T] and SearchDocuments[T] to call
// SearchRaw.
func (p SearchParams) ToRequest() *search.Request {
	return p.ToV8Request()
}

// ToV8Request converts SearchParams into a typed Elasticsearch v8 search.Request.
func (p SearchParams) ToV8Request() *search.Request {
	return query.SearchParams{
		Query:          p.Query,
		Sort:           p.Sort,
		Aggregations:   p.Aggregations,
		Highlight:      p.Highlight,
		Collapse:       p.Collapse,
		ScriptFields:   p.ScriptFields,
		Size:           p.Size,
		HasSize:        p.HasSize,
		From:           p.From,
		HasFrom:        p.HasFrom,
		SearchAfter:    p.SearchAfter,
		TrackTotalHits: p.TrackTotalHits,
		Source:         p.Source,
		Timeout:        p.Timeout,
	}.ToV8Request()
}

// SearchHit is a typed view of a single search hit.
type SearchHit[T any] struct {
	ID     string
	Index  string
	Score  *float64
	Sort   []types.FieldValue
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

// SearchRequest is the minimal input accepted by Search[T] and related helpers.
// esv8.SearchParams, esv8/query.SearchParams, and query.SearchParams all satisfy
// this interface.
type SearchRequest interface {
	ToRequest() *search.Request
}

// Search executes the preferred high-level search flow against an alias.
// It accepts any value that can convert itself into a search.Request,
// executes SearchRaw, decodes each hit's _source into T, and exposes typed
// aggregation accessors through SearchResponse[T].
func Search[T any](
	ctx context.Context,
	client SearchClient,
	aliasName estype.Alias,
	params SearchRequest,
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
			Sort:   append([]types.FieldValue(nil), hit.Sort...),
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
