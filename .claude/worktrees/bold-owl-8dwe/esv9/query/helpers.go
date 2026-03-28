package query

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// TermValue creates a Query with a single TermQuery for the given field.
func TermValue(field estype.Field, value any) types.Query {
	return types.Query{
		Term: map[string]types.TermQuery{
			string(field): {Value: value},
		},
	}
}

// TermsValues creates a Query with a TermsQuery for the given field and values.
func TermsValues(field estype.Field, values ...types.FieldValue) types.Query {
	return types.Query{
		Terms: &types.TermsQuery{
			TermsQuery: map[string]types.TermsQueryField{
				string(field): values,
			},
		},
	}
}

// MatchPhrase creates a Query with a MatchPhraseQuery for the given field.
func MatchPhrase(field estype.Field, query string) types.Query {
	return types.Query{
		MatchPhrase: map[string]types.MatchPhraseQuery{
			string(field): {Query: query},
		},
	}
}

// ExistsField creates a Query that checks for field existence.
func ExistsField(field estype.Field) types.Query {
	return types.Query{
		Exists: &types.ExistsQuery{Field: string(field)},
	}
}

// NotExists creates a Query that matches documents where the field does not exist.
func NotExists(field estype.Field) types.Query {
	return types.Query{
		Bool: &types.BoolQuery{
			MustNot: []types.Query{
				{Exists: &types.ExistsQuery{Field: string(field)}},
			},
		},
	}
}

// NestedFilter creates a nested query wrapping filter clauses.
func NestedFilter(path estype.Field, queries ...types.Query) types.Query {
	return types.Query{
		Nested: &types.NestedQuery{
			Path: string(path),
			Query: types.Query{
				Bool: &types.BoolQuery{
					Filter: queries,
				},
			},
		},
	}
}

// DateRangeQuery creates a date range query for the given field.
// Both gte and lte are optional — pass empty string to omit.
func DateRangeQuery(field estype.Field, gte, lte string) types.Query {
	rq := types.NewDateRangeQuery()
	if gte != "" {
		rq.Gte = &gte
	}
	if lte != "" {
		rq.Lte = &lte
	}
	return types.Query{
		Range: map[string]types.RangeQuery{
			string(field): rq,
		},
	}
}

// NumberRangeQuery creates a number range query for the given field.
// Pass nil for gte or lte to omit that bound.
func NumberRangeQuery(field estype.Field, gte, lte *types.Float64) types.Query {
	rq := types.NewNumberRangeQuery()
	if gte != nil {
		rq.Gte = gte
	}
	if lte != nil {
		rq.Lte = lte
	}
	return types.Query{
		Range: map[string]types.RangeQuery{
			string(field): rq,
		},
	}
}

// BoolMust creates a bool query with must clauses.
func BoolMust(queries ...types.Query) types.Query {
	return types.Query{
		Bool: &types.BoolQuery{
			Must: queries,
		},
	}
}

// BoolShould creates a bool query with should clauses.
func BoolShould(queries ...types.Query) types.Query {
	return types.Query{
		Bool: &types.BoolQuery{
			Should: queries,
		},
	}
}

// BoolFilter creates a bool query with filter clauses.
func BoolFilter(queries ...types.Query) types.Query {
	return types.Query{
		Bool: &types.BoolQuery{
			Filter: queries,
		},
	}
}

// BoolMustNot creates a bool query with must_not clauses.
func BoolMustNot(queries ...types.Query) types.Query {
	return types.Query{
		Bool: &types.BoolQuery{
			MustNot: queries,
		},
	}
}

// MatchValue creates a Query with a MatchQuery for the given field.
func MatchValue(field estype.Field, query string) types.Query {
	return types.Query{
		Match: map[string]types.MatchQuery{
			string(field): {Query: query},
		},
	}
}

// MatchAll creates a Query that matches all documents.
func MatchAll() types.Query {
	return types.Query{MatchAll: &types.MatchAllQuery{}}
}

// MatchNone creates a Query that matches no documents.
func MatchNone() types.Query {
	return types.Query{MatchNone: &types.MatchNoneQuery{}}
}

// IdsQuery creates a Query that matches documents with the given IDs.
func IdsQuery(ids ...string) types.Query {
	return types.Query{Ids: &types.IdsQuery{Values: ids}}
}

// PrefixValue creates a Query with a PrefixQuery for the given field.
func PrefixValue(field estype.Field, value string) types.Query {
	return types.Query{
		Prefix: map[string]types.PrefixQuery{
			string(field): {Value: value},
		},
	}
}

// WildcardValue creates a Query with a WildcardQuery for the given field.
func WildcardValue(field estype.Field, value string) types.Query {
	v := value
	return types.Query{
		Wildcard: map[string]types.WildcardQuery{
			string(field): {Value: &v},
		},
	}
}

// MultiMatchQuery creates a Query with a MultiMatchQuery for the given query string and fields.
func MultiMatchQuery(query string, fields ...estype.Field) types.Query {
	return types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Query:  query,
			Fields: estype.FieldNames(fields...),
		},
	}
}

// FunctionScoreQuery creates a Query with a FunctionScoreQuery.
func FunctionScoreQuery(fsq *types.FunctionScoreQuery) types.Query {
	return types.Query{FunctionScore: fsq}
}

// BoolQuery creates a Query wrapping a BoolQuery.
func BoolQuery(bq *types.BoolQuery) types.Query {
	return types.Query{Bool: bq}
}

// FieldValues converts a slice of values to []types.FieldValue for use in TermsQuery.
func FieldValues[T any](values ...T) []types.FieldValue {
	result := make([]types.FieldValue, len(values))
	for i, v := range values {
		result[i] = v
	}
	return result
}
