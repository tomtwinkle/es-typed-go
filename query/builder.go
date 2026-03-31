package query

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// IBoolQueryBuilder is a self-referential type constraint for bool query builders.
// The type parameter B must itself satisfy IBoolQueryBuilder[B], so every chain
// method returns the concrete builder type rather than a base type.
type IBoolQueryBuilder[B IBoolQueryBuilder[B]] interface {
	Must(queries ...types.Query) B
	MustNot(queries ...types.Query) B
	Should(queries ...types.Query) B
	Filter(queries ...types.Query) B
	MinimumShouldMatch(v types.MinimumShouldMatch) B
	Build() *types.BoolQuery
}

// Compile-time assertion that the concrete type satisfies the interface.
var _ IBoolQueryBuilder[*BoolQueryBuilder] = (*BoolQueryBuilder)(nil)

// BoolQueryBuilder constructs a types.BoolQuery using method chaining.
type BoolQueryBuilder struct {
	q types.BoolQuery
}

// NewBoolQuery returns a new BoolQueryBuilder.
func NewBoolQuery() *BoolQueryBuilder {
	return &BoolQueryBuilder{}
}

// Must adds must clauses to the bool query.
func (b *BoolQueryBuilder) Must(queries ...types.Query) *BoolQueryBuilder {
	b.q.Must = append(b.q.Must, queries...)
	return b
}

// MustNot adds must_not clauses to the bool query.
func (b *BoolQueryBuilder) MustNot(queries ...types.Query) *BoolQueryBuilder {
	b.q.MustNot = append(b.q.MustNot, queries...)
	return b
}

// Should adds should clauses to the bool query.
func (b *BoolQueryBuilder) Should(queries ...types.Query) *BoolQueryBuilder {
	b.q.Should = append(b.q.Should, queries...)
	return b
}

// Filter adds filter clauses to the bool query.
func (b *BoolQueryBuilder) Filter(queries ...types.Query) *BoolQueryBuilder {
	b.q.Filter = append(b.q.Filter, queries...)
	return b
}

// MinimumShouldMatch sets the minimum_should_match value.
func (b *BoolQueryBuilder) MinimumShouldMatch(v types.MinimumShouldMatch) *BoolQueryBuilder {
	b.q.MinimumShouldMatch = v
	return b
}

// Build returns the constructed types.BoolQuery.
func (b *BoolQueryBuilder) Build() *types.BoolQuery {
	return &b.q
}
