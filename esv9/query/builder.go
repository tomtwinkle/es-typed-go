// Package query provides a fluent QueryBuilder for constructing Elasticsearch v9 queries.
// Method chaining is intentionally limited to this package only.
package query

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// IQueryBuilder is a self-referential type constraint for query builders.
// The type parameter B must itself satisfy IQueryBuilder[B], so every chain
// method returns the concrete builder type rather than a base type.
// This pattern is enabled by Go 1.26's lifted restriction on self-referential
// generic type parameter lists.
type IQueryBuilder[B IQueryBuilder[B]] interface {
	Bool(bq *types.BoolQuery) B
	Match(field estype.Field, mq types.MatchQuery) B
	Term(field estype.Field, tq types.TermQuery) B
	Terms(tq *types.TermsQuery) B
	Range(field estype.Field, rq types.RangeQuery) B
	Exists(eq *types.ExistsQuery) B
	MatchAll(maq *types.MatchAllQuery) B
	MatchNone(mnq *types.MatchNoneQuery) B
	Ids(iq *types.IdsQuery) B
	Prefix(field estype.Field, pq types.PrefixQuery) B
	Wildcard(field estype.Field, wq types.WildcardQuery) B
	MultiMatch(mmq *types.MultiMatchQuery) B
	FunctionScore(fsq *types.FunctionScoreQuery) B
	Build() types.Query
}

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

// Compile-time assertions that the concrete types satisfy the interfaces.
var _ IQueryBuilder[*Builder] = (*Builder)(nil)
var _ IBoolQueryBuilder[*BoolQueryBuilder] = (*BoolQueryBuilder)(nil)

// Builder constructs an Elasticsearch types.Query using method chaining.
// This is the only place in the codebase where method chaining is used.
type Builder struct {
	query types.Query
}

// New returns a new empty Builder.
func New() *Builder {
	return &Builder{}
}

// Bool sets a BoolQuery on the query.
func (b *Builder) Bool(bq *types.BoolQuery) *Builder {
	b.query.Bool = bq
	return b
}

// Match sets a MatchQuery for the given field.
func (b *Builder) Match(field estype.Field, mq types.MatchQuery) *Builder {
	if b.query.Match == nil {
		b.query.Match = make(map[string]types.MatchQuery)
	}
	b.query.Match[string(field)] = mq
	return b
}

// Term sets a TermQuery for the given field.
func (b *Builder) Term(field estype.Field, tq types.TermQuery) *Builder {
	if b.query.Term == nil {
		b.query.Term = make(map[string]types.TermQuery)
	}
	b.query.Term[string(field)] = tq
	return b
}

// Terms sets a TermsQuery.
func (b *Builder) Terms(tq *types.TermsQuery) *Builder {
	b.query.Terms = tq
	return b
}

// Range sets a RangeQuery for the given field.
func (b *Builder) Range(field estype.Field, rq types.RangeQuery) *Builder {
	if b.query.Range == nil {
		b.query.Range = make(map[string]types.RangeQuery)
	}
	b.query.Range[string(field)] = rq
	return b
}

// Exists sets an ExistsQuery.
func (b *Builder) Exists(eq *types.ExistsQuery) *Builder {
	b.query.Exists = eq
	return b
}

// MatchAll sets a MatchAllQuery.
func (b *Builder) MatchAll(maq *types.MatchAllQuery) *Builder {
	b.query.MatchAll = maq
	return b
}

// MatchNone sets a MatchNoneQuery.
func (b *Builder) MatchNone(mnq *types.MatchNoneQuery) *Builder {
	b.query.MatchNone = mnq
	return b
}

// Ids sets an IdsQuery.
func (b *Builder) Ids(iq *types.IdsQuery) *Builder {
	b.query.Ids = iq
	return b
}

// Prefix sets a PrefixQuery for the given field.
func (b *Builder) Prefix(field estype.Field, pq types.PrefixQuery) *Builder {
	if b.query.Prefix == nil {
		b.query.Prefix = make(map[string]types.PrefixQuery)
	}
	b.query.Prefix[string(field)] = pq
	return b
}

// Wildcard sets a WildcardQuery for the given field.
func (b *Builder) Wildcard(field estype.Field, wq types.WildcardQuery) *Builder {
	if b.query.Wildcard == nil {
		b.query.Wildcard = make(map[string]types.WildcardQuery)
	}
	b.query.Wildcard[string(field)] = wq
	return b
}

// MultiMatch sets a MultiMatchQuery.
func (b *Builder) MultiMatch(mmq *types.MultiMatchQuery) *Builder {
	b.query.MultiMatch = mmq
	return b
}

// FunctionScore sets a FunctionScoreQuery.
func (b *Builder) FunctionScore(fsq *types.FunctionScoreQuery) *Builder {
	b.query.FunctionScore = fsq
	return b
}

// Build returns the constructed types.Query.
func (b *Builder) Build() types.Query {
	return b.query
}

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
