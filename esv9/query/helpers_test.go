package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/esv9/query"
)

func TestTermValue(t *testing.T) {
	t.Parallel()
	q := query.TermValue("customer_id", "abc123")
	assert.Assert(t, q.Term != nil)
	assert.Equal(t, "abc123", q.Term["customer_id"].Value)
}

func TestTermsValues(t *testing.T) {
	t.Parallel()
	q := query.TermsValues("integration_id",
		types.FieldValue("id1"),
		types.FieldValue("id2"),
		types.FieldValue("id3"),
	)
	assert.Assert(t, q.Terms != nil)
	assert.Assert(t, q.Terms.TermsQuery != nil)
	vals, ok := q.Terms.TermsQuery["integration_id"].([]types.FieldValue)
	assert.Assert(t, ok)
	assert.Assert(t, len(vals) == 3)
}

func TestMatchPhrase(t *testing.T) {
	t.Parallel()
	q := query.MatchPhrase("name.ngram_case_insensitive", "test keyword")
	assert.Assert(t, q.MatchPhrase != nil)
	assert.Equal(t, "test keyword", q.MatchPhrase["name.ngram_case_insensitive"].Query)
}

func TestExistsField(t *testing.T) {
	t.Parallel()
	q := query.ExistsField("construction_status_ids")
	assert.Assert(t, q.Exists != nil)
	assert.Equal(t, "construction_status_ids", q.Exists.Field)
}

func TestNotExists(t *testing.T) {
	t.Parallel()
	q := query.NotExists("construction_status_ids")
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.MustNot) == 1)
	assert.Assert(t, q.Bool.MustNot[0].Exists != nil)
	assert.Equal(t, "construction_status_ids", q.Bool.MustNot[0].Exists.Field)
}

func TestNestedFilter(t *testing.T) {
	t.Parallel()
	inner := query.TermValue("construction_items.color_id", "red")
	q := query.NestedFilter("construction_items", inner)
	assert.Assert(t, q.Nested != nil)
	assert.Equal(t, "construction_items", q.Nested.Path)
	assert.Assert(t, q.Nested.Query.Bool != nil)
	assert.Assert(t, len(q.Nested.Query.Bool.Filter) == 1)
}

func TestNestedFilter_MultipleQueries(t *testing.T) {
	t.Parallel()
	q1 := query.TermValue("construction_items.color_id", "red")
	q2 := query.TermValue("construction_items.status", "active")
	q := query.NestedFilter("construction_items", q1, q2)
	assert.Assert(t, q.Nested != nil)
	assert.Assert(t, len(q.Nested.Query.Bool.Filter) == 2)
}

func TestDateRangeQuery(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery("created_at", "2024-01-01", "2024-12-31")
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["created_at"] != nil)
}

func TestDateRangeQuery_GteOnly(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery("created_at", "2024-01-01", "")
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["created_at"] != nil)
}

func TestDateRangeQuery_LteOnly(t *testing.T) {
	t.Parallel()
	q := query.DateRangeQuery("created_at", "", "2024-12-31")
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["created_at"] != nil)
}

func TestNumberRangeQuery(t *testing.T) {
	t.Parallel()
	gte := types.Float64(10.0)
	lte := types.Float64(100.0)
	q := query.NumberRangeQuery("price", &gte, &lte)
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["price"] != nil)
}

func TestNumberRangeQuery_GteOnly(t *testing.T) {
	t.Parallel()
	gte := types.Float64(10.0)
	q := query.NumberRangeQuery("price", &gte, nil)
	assert.Assert(t, q.Range != nil)
	assert.Assert(t, q.Range["price"] != nil)
}

func TestBoolMust(t *testing.T) {
	t.Parallel()
	q := query.BoolMust(
		query.TermValue("status", "active"),
		query.TermValue("draft", false),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Must) == 2)
}

func TestBoolShould(t *testing.T) {
	t.Parallel()
	q := query.BoolShould(
		query.MatchPhrase("name.ngram_case_insensitive", "keyword"),
		query.MatchPhrase("name.ngram_exact", "keyword"),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Should) == 2)
}

func TestBoolFilter(t *testing.T) {
	t.Parallel()
	q := query.BoolFilter(
		query.TermValue("client_id", "c1"),
		query.TermValue("draft", false),
	)
	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Filter) == 2)
}

func TestBoolMustNot(t *testing.T) {
	t.Parallel()
	q := query.BoolMustNot(
		query.ExistsField("deleted_at"),
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
	q := query.TermsValues("tag_ids", query.FieldValues(int32(1), int32(2), int32(3))...)
	assert.Assert(t, q.Terms != nil)
	vals, ok := q.Terms.TermsQuery["tag_ids"].([]types.FieldValue)
	assert.Assert(t, ok)
	assert.Assert(t, len(vals) == 3)
}

func TestComplexQueryCombination(t *testing.T) {
	t.Parallel()
	filters := []types.Query{
		query.TermValue("client_id", "client1"),
		query.TermsValues("integration_id", query.FieldValues("id1", "id2")...),
		query.TermValue("draft", false),
		query.MatchPhrase("name.ngram_case_insensitive", "search keyword"),
		query.NotExists("construction_status_ids"),
		query.NestedFilter("construction_items",
			query.TermsValues("construction_items.crafts_user_ids", query.FieldValues(int32(1), int32(2))...),
		),
	}

	q := query.New().Bool(
		query.NewBoolQuery().Filter(filters...).Build(),
	).Build()

	assert.Assert(t, q.Bool != nil)
	assert.Assert(t, len(q.Bool.Filter) == 6)
}
