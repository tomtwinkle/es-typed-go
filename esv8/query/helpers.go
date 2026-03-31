// Package query provides query building helpers for Elasticsearch v8.
//
// Deprecated: import github.com/tomtwinkle/es-typed-go/query instead.
// All symbols in this package are type aliases or forwarding wrappers.
package query

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
	sharedquery "github.com/tomtwinkle/es-typed-go/query"
)

func TermValue(field estype.Field, value any) types.Query {
	return sharedquery.TermValue(field, value)
}
func TermsValues(field estype.Field, values ...types.FieldValue) types.Query {
	return sharedquery.TermsValues(field, values...)
}
func MatchPhrase(field estype.Field, query string) types.Query {
	return sharedquery.MatchPhrase(field, query)
}
func ExistsField(field estype.Field) types.Query {
	return sharedquery.ExistsField(field)
}
func NotExists(field estype.Field) types.Query {
	return sharedquery.NotExists(field)
}
func NestedFilter(path estype.Field, queries ...types.Query) types.Query {
	return sharedquery.NestedFilter(path, queries...)
}
func DateRangeQuery(field estype.Field, gte, lte string) types.Query {
	return sharedquery.DateRangeQuery(field, gte, lte)
}
func NumberRangeQuery(field estype.Field, gte, lte *types.Float64) types.Query {
	return sharedquery.NumberRangeQuery(field, gte, lte)
}
func BoolMust(queries ...types.Query) types.Query   { return sharedquery.BoolMust(queries...) }
func BoolShould(queries ...types.Query) types.Query { return sharedquery.BoolShould(queries...) }
func BoolFilter(queries ...types.Query) types.Query { return sharedquery.BoolFilter(queries...) }
func BoolMustNot(queries ...types.Query) types.Query {
	return sharedquery.BoolMustNot(queries...)
}
func MatchValue(field estype.Field, query string) types.Query {
	return sharedquery.MatchValue(field, query)
}
func MatchAll() types.Query  { return sharedquery.MatchAll() }
func MatchNone() types.Query { return sharedquery.MatchNone() }
func IdsQuery(ids ...string) types.Query { return sharedquery.IdsQuery(ids...) }
func PrefixValue(field estype.Field, value string) types.Query {
	return sharedquery.PrefixValue(field, value)
}
func WildcardValue(field estype.Field, value string) types.Query {
	return sharedquery.WildcardValue(field, value)
}
func MultiMatchQuery(query string, fields ...estype.Field) types.Query {
	return sharedquery.MultiMatchQuery(query, fields...)
}
func FunctionScoreQuery(fsq *types.FunctionScoreQuery) types.Query {
	return sharedquery.FunctionScoreQuery(fsq)
}
func BoolQuery(bq *types.BoolQuery) types.Query { return sharedquery.BoolQuery(bq) }
func FieldValues[T any](values ...T) []types.FieldValue {
	return sharedquery.FieldValues[T](values...)
}
