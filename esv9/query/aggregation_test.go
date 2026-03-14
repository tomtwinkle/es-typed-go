package query_test

import (
	"math"
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/calendarinterval"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv9/query"
)

func TestNewAggregations_Empty(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Build()
	assert.Assert(t, len(aggs) == 0)
}

func TestAggregationBuilder_Terms(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Terms("by_category", estype.Field("category")).Build()
	_, ok := aggs["by_category"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["by_category"].Terms != nil)
	assert.Equal(t, "category", *aggs["by_category"].Terms.Field)
}

func TestAggregationBuilder_TermsWithSize(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().TermsWithSize("top10", estype.Field("category"), 10).Build()
	_, ok := aggs["top10"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["top10"].Terms != nil)
	assert.Equal(t, "category", *aggs["top10"].Terms.Field)
	assert.Equal(t, 10, *aggs["top10"].Terms.Size)
}

func TestAggregationBuilder_DateHistogram(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().
		DateHistogram("by_month", estype.Field("date"), calendarinterval.Month).
		Build()
	_, ok := aggs["by_month"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["by_month"].DateHistogram != nil)
	assert.Equal(t, "date", *aggs["by_month"].DateHistogram.Field)
	assert.Equal(t, calendarinterval.Month, *aggs["by_month"].DateHistogram.CalendarInterval)
}

func TestAggregationBuilder_DateHistogramWithFormat(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().
		DateHistogramWithFormat("by_year", estype.Field("date"), "yyyy", calendarinterval.Year).
		Build()
	_, ok := aggs["by_year"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["by_year"].DateHistogram != nil)
	assert.Equal(t, "yyyy", *aggs["by_year"].DateHistogram.Format)
}

func TestAggregationBuilder_Histogram(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Histogram("price_ranges", estype.Field("price"), 100.0).Build()
	_, ok := aggs["price_ranges"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["price_ranges"].Histogram != nil)
	assert.Equal(t, "price", *aggs["price_ranges"].Histogram.Field)
	assert.Assert(t, math.Abs(float64(*aggs["price_ranges"].Histogram.Interval)-100.0) < 0.001)
}

func TestAggregationBuilder_Avg(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Avg("avg_price", estype.Field("price")).Build()
	_, ok := aggs["avg_price"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["avg_price"].Avg != nil)
	assert.Equal(t, "price", *aggs["avg_price"].Avg.Field)
}

func TestAggregationBuilder_Max(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Max("max_price", estype.Field("price")).Build()
	_, ok := aggs["max_price"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["max_price"].Max != nil)
}

func TestAggregationBuilder_Min(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Min("min_price", estype.Field("price")).Build()
	_, ok := aggs["min_price"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["min_price"].Min != nil)
}

func TestAggregationBuilder_Sum(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Sum("total_sales", estype.Field("value")).Build()
	_, ok := aggs["total_sales"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["total_sales"].Sum != nil)
	assert.Equal(t, "value", *aggs["total_sales"].Sum.Field)
}

func TestAggregationBuilder_ValueCount(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().ValueCount("order_count", estype.Field("id")).Build()
	_, ok := aggs["order_count"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["order_count"].ValueCount != nil)
}

func TestAggregationBuilder_Cardinality(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Cardinality("unique_users", estype.Field("id")).Build()
	_, ok := aggs["unique_users"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["unique_users"].Cardinality != nil)
	assert.Equal(t, "id", *aggs["unique_users"].Cardinality.Field)
}

func TestAggregationBuilder_Stats(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Stats("price_stats", estype.Field("price")).Build()
	_, ok := aggs["price_stats"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["price_stats"].Stats != nil)
	assert.Equal(t, "price", *aggs["price_stats"].Stats.Field)
}

func TestAggregationBuilder_Nested(t *testing.T) {
	t.Parallel()
	sub := query.NewAggregations().Avg("avg_price", estype.Field("items.price"))
	aggs := query.NewAggregations().Nested("items_agg", estype.Field("items"), sub).Build()
	_, ok := aggs["items_agg"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["items_agg"].Nested != nil)
	assert.Equal(t, "items", *aggs["items_agg"].Nested.Path)
	_, ok = aggs["items_agg"].Aggregations["avg_price"]
	assert.Assert(t, ok)
}

func TestAggregationBuilder_Filter(t *testing.T) {
	t.Parallel()
	sub := query.NewAggregations().Avg("avg_price", estype.Field("price"))
	filter := types.Query{
		Term: map[string]types.TermQuery{
			"status": {Value: "active"},
		},
	}
	aggs := query.NewAggregations().Filter("active_items", filter, sub).Build()
	_, ok := aggs["active_items"]
	assert.Assert(t, ok)
	assert.Assert(t, aggs["active_items"].Filter != nil)
	_, ok = aggs["active_items"].Aggregations["avg_price"]
	assert.Assert(t, ok)
}

func TestAggregationBuilder_SubAggregations(t *testing.T) {
	t.Parallel()
	sub := query.NewAggregations().Avg("avg_price", estype.Field("price")).Sum("total_sales", estype.Field("price"))
	aggs := query.NewAggregations().
		Terms("by_category", estype.Field("category")).
		SubAggregations("by_category", sub).
		Build()
	_, ok := aggs["by_category"]
	assert.Assert(t, ok)
	_, ok = aggs["by_category"].Aggregations["avg_price"]
	assert.Assert(t, ok)
	_, ok = aggs["by_category"].Aggregations["total_sales"]
	assert.Assert(t, ok)
}

func TestAggregationBuilder_SubAggregations_NonExistentParent(t *testing.T) {
	t.Parallel()
	// SubAggregations on a name that was never added should be a no-op.
	sub := query.NewAggregations().Avg("avg_price", estype.Field("price"))
	aggs := query.NewAggregations().
		SubAggregations("does_not_exist", sub).
		Build()
	assert.Assert(t, len(aggs) == 0)
}

func TestAggregationBuilder_Chaining(t *testing.T) {
	t.Parallel()
	// Verify that multiple aggregations can be added in one chain.
	aggs := query.NewAggregations().
		Terms("by_category", estype.Field("category")).
		Avg("avg_price", estype.Field("price")).
		Max("max_price", estype.Field("price")).
		Min("min_price", estype.Field("price")).
		Sum("total_revenue", estype.Field("value")).
		Build()
	assert.Assert(t, len(aggs) == 5)
	assert.Assert(t, aggs["by_category"].Terms != nil)
	assert.Assert(t, aggs["avg_price"].Avg != nil)
	assert.Assert(t, aggs["max_price"].Max != nil)
	assert.Assert(t, aggs["min_price"].Min != nil)
	assert.Assert(t, aggs["total_revenue"].Sum != nil)
}
