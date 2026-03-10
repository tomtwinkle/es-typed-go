package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomtwinkle/es-typed-go/query"
)

func TestNewAggregations_Empty(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Build()
	assert.Empty(t, aggs)
}

func TestAggregationBuilder_Terms(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Terms("by_category", "category").Build()
	require.Contains(t, aggs, "by_category")
	require.NotNil(t, aggs["by_category"].Terms)
	assert.Equal(t, "category", *aggs["by_category"].Terms.Field)
}

func TestAggregationBuilder_TermsWithSize(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().TermsWithSize("top10", "brand", 10).Build()
	require.Contains(t, aggs, "top10")
	require.NotNil(t, aggs["top10"].Terms)
	assert.Equal(t, "brand", *aggs["top10"].Terms.Field)
	assert.Equal(t, 10, *aggs["top10"].Terms.Size)
}

func TestAggregationBuilder_DateHistogram(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().
		DateHistogram("by_month", "created_at", calendarinterval.Month).
		Build()
	require.Contains(t, aggs, "by_month")
	require.NotNil(t, aggs["by_month"].DateHistogram)
	assert.Equal(t, "created_at", *aggs["by_month"].DateHistogram.Field)
	assert.Equal(t, calendarinterval.Month, *aggs["by_month"].DateHistogram.CalendarInterval)
}

func TestAggregationBuilder_DateHistogramWithFormat(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().
		DateHistogramWithFormat("by_year", "date", "yyyy", calendarinterval.Year).
		Build()
	require.Contains(t, aggs, "by_year")
	require.NotNil(t, aggs["by_year"].DateHistogram)
	assert.Equal(t, "yyyy", *aggs["by_year"].DateHistogram.Format)
}

func TestAggregationBuilder_Histogram(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Histogram("price_ranges", "price", 100.0).Build()
	require.Contains(t, aggs, "price_ranges")
	require.NotNil(t, aggs["price_ranges"].Histogram)
	assert.Equal(t, "price", *aggs["price_ranges"].Histogram.Field)
	assert.InDelta(t, 100.0, float64(*aggs["price_ranges"].Histogram.Interval), 0.001)
}

func TestAggregationBuilder_Avg(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Avg("avg_price", "price").Build()
	require.Contains(t, aggs, "avg_price")
	require.NotNil(t, aggs["avg_price"].Avg)
	assert.Equal(t, "price", *aggs["avg_price"].Avg.Field)
}

func TestAggregationBuilder_Max(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Max("max_price", "price").Build()
	require.Contains(t, aggs, "max_price")
	require.NotNil(t, aggs["max_price"].Max)
}

func TestAggregationBuilder_Min(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Min("min_price", "price").Build()
	require.Contains(t, aggs, "min_price")
	require.NotNil(t, aggs["min_price"].Min)
}

func TestAggregationBuilder_Sum(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Sum("total_sales", "amount").Build()
	require.Contains(t, aggs, "total_sales")
	require.NotNil(t, aggs["total_sales"].Sum)
	assert.Equal(t, "amount", *aggs["total_sales"].Sum.Field)
}

func TestAggregationBuilder_ValueCount(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().ValueCount("order_count", "order_id").Build()
	require.Contains(t, aggs, "order_count")
	require.NotNil(t, aggs["order_count"].ValueCount)
}

func TestAggregationBuilder_Cardinality(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Cardinality("unique_users", "user_id").Build()
	require.Contains(t, aggs, "unique_users")
	require.NotNil(t, aggs["unique_users"].Cardinality)
	assert.Equal(t, "user_id", *aggs["unique_users"].Cardinality.Field)
}

func TestAggregationBuilder_Stats(t *testing.T) {
	t.Parallel()
	aggs := query.NewAggregations().Stats("price_stats", "price").Build()
	require.Contains(t, aggs, "price_stats")
	require.NotNil(t, aggs["price_stats"].Stats)
	assert.Equal(t, "price", *aggs["price_stats"].Stats.Field)
}

func TestAggregationBuilder_Nested(t *testing.T) {
	t.Parallel()
	sub := query.NewAggregations().Avg("avg_price", "products.price")
	aggs := query.NewAggregations().Nested("products_agg", "products", sub).Build()
	require.Contains(t, aggs, "products_agg")
	require.NotNil(t, aggs["products_agg"].Nested)
	assert.Equal(t, "products", *aggs["products_agg"].Nested.Path)
	require.Contains(t, aggs["products_agg"].Aggregations, "avg_price")
}

func TestAggregationBuilder_Filter(t *testing.T) {
	t.Parallel()
	sub := query.NewAggregations().Avg("avg_price", "price")
	filter := types.Query{
		Term: map[string]types.TermQuery{
			"status": {Value: "active"},
		},
	}
	aggs := query.NewAggregations().Filter("active_products", filter, sub).Build()
	require.Contains(t, aggs, "active_products")
	require.NotNil(t, aggs["active_products"].Filter)
	require.Contains(t, aggs["active_products"].Aggregations, "avg_price")
}

func TestAggregationBuilder_SubAggregations(t *testing.T) {
	t.Parallel()
	sub := query.NewAggregations().Avg("avg_price", "price").Sum("total_sales", "price")
	aggs := query.NewAggregations().
		Terms("by_category", "category").
		SubAggregations("by_category", sub).
		Build()
	require.Contains(t, aggs, "by_category")
	require.Contains(t, aggs["by_category"].Aggregations, "avg_price")
	require.Contains(t, aggs["by_category"].Aggregations, "total_sales")
}

func TestAggregationBuilder_SubAggregations_NonExistentParent(t *testing.T) {
	t.Parallel()
	// SubAggregations on a name that was never added should be a no-op.
	sub := query.NewAggregations().Avg("avg_price", "price")
	aggs := query.NewAggregations().
		SubAggregations("does_not_exist", sub).
		Build()
	assert.Empty(t, aggs)
}

func TestAggregationBuilder_Chaining(t *testing.T) {
	t.Parallel()
	// Verify that multiple aggregations can be added in one chain.
	aggs := query.NewAggregations().
		Terms("by_category", "category").
		Avg("avg_price", "price").
		Max("max_price", "price").
		Min("min_price", "price").
		Sum("total_revenue", "revenue").
		Build()
	assert.Len(t, aggs, 5)
	assert.NotNil(t, aggs["by_category"].Terms)
	assert.NotNil(t, aggs["avg_price"].Avg)
	assert.NotNil(t, aggs["max_price"].Max)
	assert.NotNil(t, aggs["min_price"].Min)
	assert.NotNil(t, aggs["total_revenue"].Sum)
}
