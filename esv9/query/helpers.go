package query

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

// TermValue creates a Query with a single TermQuery for the given field.
func TermValue(field string, value any) types.Query {
	return types.Query{
		Term: map[string]types.TermQuery{
			field: {Value: value},
		},
	}
}

// TermsValues creates a Query with a TermsQuery for the given field and values.
func TermsValues(field string, values ...types.FieldValue) types.Query {
	return types.Query{
		Terms: &types.TermsQuery{
			TermsQuery: map[string]types.TermsQueryField{
				field: values,
			},
		},
	}
}

// MatchPhrase creates a Query with a MatchPhraseQuery for the given field.
func MatchPhrase(field, query string) types.Query {
	return types.Query{
		MatchPhrase: map[string]types.MatchPhraseQuery{
			field: {Query: query},
		},
	}
}

// ExistsField creates a Query that checks for field existence.
func ExistsField(field string) types.Query {
	return types.Query{
		Exists: &types.ExistsQuery{Field: field},
	}
}

// NotExists creates a Query that matches documents where the field does not exist.
func NotExists(field string) types.Query {
	return types.Query{
		Bool: &types.BoolQuery{
			MustNot: []types.Query{
				{Exists: &types.ExistsQuery{Field: field}},
			},
		},
	}
}

// NestedFilter creates a nested query wrapping filter clauses.
func NestedFilter(path string, queries ...types.Query) types.Query {
	return types.Query{
		Nested: &types.NestedQuery{
			Path: path,
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
func DateRangeQuery(field string, gte, lte string) types.Query {
	rq := types.NewDateRangeQuery()
	if gte != "" {
		rq.Gte = &gte
	}
	if lte != "" {
		rq.Lte = &lte
	}
	return types.Query{
		Range: map[string]types.RangeQuery{
			field: rq,
		},
	}
}

// NumberRangeQuery creates a number range query for the given field.
// Pass nil for gte or lte to omit that bound.
func NumberRangeQuery(field string, gte, lte *types.Float64) types.Query {
	rq := types.NewNumberRangeQuery()
	if gte != nil {
		rq.Gte = gte
	}
	if lte != nil {
		rq.Lte = lte
	}
	return types.Query{
		Range: map[string]types.RangeQuery{
			field: rq,
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

// FieldValues converts a slice of values to []types.FieldValue for use in TermsQuery.
func FieldValues[T any](values ...T) []types.FieldValue {
	result := make([]types.FieldValue, len(values))
	for i, v := range values {
		result[i] = v
	}
	return result
}
