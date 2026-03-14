package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv9/query"
)

func TestTermValue(t *testing.T) {
	t.Parallel()
	q := query.TermValue(estype.Field("status"), "active")
	assert.Assert(t, q.Term != nil)
	assert.Equal(t, "active", q.Term["status"].Value)
}

func TestTermsValues(t *testing.T) {
	t.Parallel()
	q := query.TermsValues(estype.Field("category"),
		types.FieldValue("val1"),
		types.FieldValue("val2"),
		types.FieldValue("val3"),
	)
	assert.Assert(t, q.Terms != nil)
	assert.Assert(t, q.Terms.TermsQuery != nil)
	vals, ok := q.Terms.TermsQuery["category"].([]types.FieldValue)
	assert.Assert(t, ok)
	assert.Assert(t, len(vals) == 3)
}

func TestMatchPhrase(t *testing.T) {
	t.Parallel()
	q := query.MatchPhrase(estype.Field("title.ngram"), "test keyword")
	assert.Assert(t, q.MatchPhrase != nil)
	assert.Equal(t, "test keyword", q.MatchPhrase["title.ngram"].Query)
}

func TestExistsField(t *testing.T) {
	t.Parallel()
	q := query.ExistsField(estype.Field("tags"))
	assert.Assert(t, q.Exists != nil)
	assert.Equal(t, "tags", q.Exists.Field)
}

func TestNotExists(t *testing.T) {
	t.Parallel()
	q := query.NotExists(estype.Field("tags"))
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.MustNot) == 1)
	assert.Assert(t, q.Bool.MustNot[0].Exists != nil)
	assert.Equal(t, "tags", q.Bool.MustNot[0].Exists.Field)
}

func TestNestedFilter(t *testing.T) {
	t.Parallel()
	inner := query.TermValue(estype.Field("items.color"), "red")
	q := query.NestedFilter(estype.Field("items"), inner)
	assert.Assert(t, q.Nested != nil)
	assert.Equal(t, "items", q.Nested.Path)
	assert.Assert(t, q.Nested.Query.Bool != nil)
	assert.Assert(t, len(q.Nested.Query.Bool.Filter) == 1)
}

func TestNestedFilter_MultipleQueries(t *testing.T) {
	t.Parallel()
	q1 := query.TermValue(estype.Field("items.color"), "red")
	q2 := query.TermValue(estype.Field("items.status"), "active")
	q := query.NestedFilter(estype.Field("items"), q1, q2)
	assert.Assert(t, q.Nested != nil)
	assert.Assert(t, len(q.Nested.Query.Bool.Filter) == 2)
}

func TestDateRangeQuery(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery(estype.Field("date"), "2024-01-01", "2024-12-31")
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["date"] != nil)
}

func TestDateRangeQuery_GteOnly(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery(estype.Field("date"), "2024-01-01", "")
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["date"] != nil)
}

func TestDateRangeQuery_LteOnly(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery(estype.Field("date"), "", "2024-12-31")
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["date"] != nil)
}

func TestNumberRangeQuery(t *testing.T) {
	t.Parallel()
	gte := types.Float64(10.0)
	lte := types.Float64(100.0)
	q := query.NumberRangeQuery(estype.Field("price"), &gte, &lte)
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["price"] != nil)
}

func TestNumberRangeQuery_GteOnly(t *testing.T) {
	t.Parallel()
	gte := types.Float64(10.0)
	q := query.NumberRangeQuery(estype.Field("price"), &gte, nil)
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["price"] != nil)
}

func TestNumberRangeQuery_LteOnly(t *testing.T) {
	t.Parallel()
	lte := types.Float64(100.0)
	q := query.NumberRangeQuery(estype.Field("price"), nil, &lte)
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["price"] != nil)
}

func TestBoolMust(t *testing.T) {
	t.Parallel()
	q := query.BoolMust(
		query.TermValue(estype.Field("status"), "active"),
		query.TermValue(estype.Field("enabled"), true),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Must) == 2)
}

func TestBoolShould(t *testing.T) {
	t.Parallel()
	q := query.BoolShould(
		query.MatchPhrase(estype.Field("title.ngram"), "keyword"),
		query.MatchPhrase(estype.Field("title.raw"), "keyword"),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Should) == 2)
}

func TestBoolFilter(t *testing.T) {
	t.Parallel()
	q := query.BoolFilter(
		query.TermValue(estype.Field("type"), "document"),
		query.TermValue(estype.Field("status"), "active"),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Filter) == 2)
}

func TestBoolMustNot(t *testing.T) {
	t.Parallel()
	q := query.BoolMustNot(
		query.ExistsField(estype.Field("date")),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.MustNot) == 1)
}

func TestFieldValues_Strings(t *testing.T) {
	t.Parallel()
	vals := query.FieldValues("a", "b", "c")
	assert.Assert(t, len(vals) == 3)
	assert.Equal(t, "a", vals[0])
	assert.Equal(t, "b", vals[1])
	assert.Equal(t, "c", vals[2])
}

func TestFieldValues_Int32(t *testing.T) {
	t.Parallel()
	vals := query.FieldValues(int32(1), int32(2), int32(3))
	assert.Assert(t, len(vals) == 3)
	assert.Equal(t, int32(1), vals[0])
}

func TestFieldValues_Empty(t *testing.T) {
	t.Parallel()
	vals := query.FieldValues[string]()
	assert.Assert(t, len(vals) == 0)
}

func TestFieldValues_WithTermsValues(t *testing.T) {
	t.Parallel()
	q := query.TermsValues(estype.Field("tags"), query.FieldValues(int32(1), int32(2), int32(3))...)
	assert.Assert(t, q.Terms != nil)
	vals, ok := q.Terms.TermsQuery["tags"].([]types.FieldValue)
	assert.Assert(t, ok)
	assert.Assert(t, len(vals) == 3)
}

func TestComplexQueryCombination(t *testing.T) {
	t.Parallel()
	// Demonstrates building a complex query using the helper functions
	// instead of verbose struct construction.
	filters := []types.Query{
		query.TermValue(estype.Field("type"), "document"),
		query.TermsValues(estype.Field("category"), query.FieldValues("cat1", "cat2")...),
		query.TermValue(estype.Field("status"), "active"),
		query.MatchPhrase(estype.Field("title.ngram"), "search keyword"),
		query.NotExists(estype.Field("tags")),
		query.NestedFilter(estype.Field("items"),
			query.TermsValues(estype.Field("items.ids"), query.FieldValues(int32(1), int32(2))...),
		),
	}

	q := query.New().Bool(
		query.NewBoolQuery().Filter(filters...).Build(),
	).Build()

	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Filter) == 6)
}
