package query

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

// SearchParams holds all parameters for a search request built by a SearchBuilder.
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

// ISearchBuilder is a self-referential type constraint for search builders.
// The type parameter B must itself satisfy ISearchBuilder[B], so every chain
// method returns the concrete builder type rather than a base type.
// This pattern is enabled by Go 1.26's lifted restriction on self-referential
// generic type parameter lists.
type ISearchBuilder[B ISearchBuilder[B]] interface {
	// Query sets the full query directly. This replaces any bool query
	// clauses set via Where, Must, Should, or MustNot.
	Query(q types.Query) B

	// Where adds filter clauses to the bool query.
	// Multiple calls accumulate filter clauses.
	Where(queries ...types.Query) B

	// Must adds must clauses to the bool query.
	// Multiple calls accumulate must clauses.
	Must(queries ...types.Query) B

	// Should adds should clauses to the bool query.
	// Multiple calls accumulate should clauses.
	Should(queries ...types.Query) B

	// MustNot adds must_not clauses to the bool query.
	// Multiple calls accumulate must_not clauses.
	MustNot(queries ...types.Query) B

	// MinimumShouldMatch sets the minimum_should_match value on the bool query.
	MinimumShouldMatch(v types.MinimumShouldMatch) B

	// Sort sets the sort specification. Each call replaces the previous sort.
	Sort(sorts ...types.SortCombinations) B

	// Limit sets the maximum number of results to return.
	Limit(size int) B

	// Offset sets the number of results to skip.
	Offset(from int) B

	// Aggregation sets the aggregations. Each call replaces the previous aggregations.
	Aggregation(aggs map[string]types.Aggregations) B

	// Highlight sets the highlight configuration.
	Highlight(h *types.Highlight) B

	// Collapse sets the field collapse configuration.
	Collapse(c *types.FieldCollapse) B

	// ScriptFields sets the script fields.
	ScriptFields(sf map[string]types.ScriptField) B

	// Build returns the constructed SearchParams.
	Build() SearchParams
}

// Compile-time assertion that the concrete type satisfies the interface.
var _ ISearchBuilder[*SearchBuilder] = (*SearchBuilder)(nil)

// SearchBuilder constructs search parameters using method chaining.
// It provides an ActiveRecord-style fluent API: Where/Must/Should/MustNot
// for query clauses, Sort for ordering, Limit/Offset for pagination,
// and Build to produce the final SearchParams.
type SearchBuilder struct {
	query              *types.Query
	filters            []types.Query
	musts              []types.Query
	shoulds            []types.Query
	mustNots           []types.Query
	minimumShouldMatch types.MinimumShouldMatch
	sorts              []types.SortCombinations
	aggs               map[string]types.Aggregations
	highlight          *types.Highlight
	collapse           *types.FieldCollapse
	scriptFields       map[string]types.ScriptField
	size               int
	from               int
}

// NewSearch returns a new empty SearchBuilder.
func NewSearch() *SearchBuilder {
	return &SearchBuilder{}
}

// Query sets the full query directly. This replaces any bool query
// clauses set via Where, Must, Should, or MustNot.
func (b *SearchBuilder) Query(q types.Query) *SearchBuilder {
	b.query = &q
	return b
}

// Where adds filter clauses to the bool query.
func (b *SearchBuilder) Where(queries ...types.Query) *SearchBuilder {
	b.filters = append(b.filters, queries...)
	return b
}

// Must adds must clauses to the bool query.
func (b *SearchBuilder) Must(queries ...types.Query) *SearchBuilder {
	b.musts = append(b.musts, queries...)
	return b
}

// Should adds should clauses to the bool query.
func (b *SearchBuilder) Should(queries ...types.Query) *SearchBuilder {
	b.shoulds = append(b.shoulds, queries...)
	return b
}

// MustNot adds must_not clauses to the bool query.
func (b *SearchBuilder) MustNot(queries ...types.Query) *SearchBuilder {
	b.mustNots = append(b.mustNots, queries...)
	return b
}

// MinimumShouldMatch sets the minimum_should_match value on the bool query.
func (b *SearchBuilder) MinimumShouldMatch(v types.MinimumShouldMatch) *SearchBuilder {
	b.minimumShouldMatch = v
	return b
}

// Sort sets the sort specification.
func (b *SearchBuilder) Sort(sorts ...types.SortCombinations) *SearchBuilder {
	b.sorts = sorts
	return b
}

// Limit sets the maximum number of results to return.
func (b *SearchBuilder) Limit(size int) *SearchBuilder {
	b.size = size
	return b
}

// Offset sets the number of results to skip.
func (b *SearchBuilder) Offset(from int) *SearchBuilder {
	b.from = from
	return b
}

// Aggregation sets the aggregations.
func (b *SearchBuilder) Aggregation(aggs map[string]types.Aggregations) *SearchBuilder {
	b.aggs = aggs
	return b
}

// Highlight sets the highlight configuration.
func (b *SearchBuilder) Highlight(h *types.Highlight) *SearchBuilder {
	b.highlight = h
	return b
}

// Collapse sets the field collapse configuration.
func (b *SearchBuilder) Collapse(c *types.FieldCollapse) *SearchBuilder {
	b.collapse = c
	return b
}

// ScriptFields sets the script fields.
func (b *SearchBuilder) ScriptFields(sf map[string]types.ScriptField) *SearchBuilder {
	b.scriptFields = sf
	return b
}

// Build returns the constructed SearchParams.
// If a full query was set via Query(), it takes precedence.
// Otherwise, bool query clauses from Where, Must, Should, and MustNot are combined.
func (b *SearchBuilder) Build() SearchParams {
	var q types.Query
	if b.query != nil {
		q = *b.query
	} else if len(b.filters) > 0 || len(b.musts) > 0 || len(b.shoulds) > 0 || len(b.mustNots) > 0 {
		bq := &types.BoolQuery{}
		if len(b.filters) > 0 {
			bq.Filter = b.filters
		}
		if len(b.musts) > 0 {
			bq.Must = b.musts
		}
		if len(b.shoulds) > 0 {
			bq.Should = b.shoulds
		}
		if len(b.mustNots) > 0 {
			bq.MustNot = b.mustNots
		}
		if b.minimumShouldMatch != nil {
			bq.MinimumShouldMatch = b.minimumShouldMatch
		}
		q.Bool = bq
	}

	return SearchParams{
		Query:        q,
		Sort:         b.sorts,
		Aggregations: b.aggs,
		Highlight:    b.highlight,
		Collapse:     b.collapse,
		ScriptFields: b.scriptFields,
		Size:         b.size,
		From:         b.from,
	}
}
