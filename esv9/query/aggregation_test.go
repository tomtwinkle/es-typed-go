package query_test

import (
	"math"
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/calendarinterval"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/esv9/query"
)

func TestAggs_Empty(t *testing.T) {
	t.Parallel()

	aggs := query.Aggs().Build()
	assert.Assert(t, aggs == nil)
}

func TestStringTermsAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("by_category", FieldCategory)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["by_category"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Terms != nil)
	assert.Equal(t, string(FieldCategory), *agg.Terms.Field)
}

func TestStringTermsAgg_WithSize(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("top10", FieldCategory, query.TermsAggSize(10))
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["top10"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Terms != nil)
	assert.Equal(t, string(FieldCategory), *agg.Terms.Field)
	assert.Equal(t, 10, *agg.Terms.Size)
}

func TestDateHistogramAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["by_month"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.DateHistogram != nil)
	assert.Equal(t, string(FieldDate), *agg.DateHistogram.Field)
	assert.Equal(t, calendarinterval.Month, *agg.DateHistogram.CalendarInterval)
}

func TestDateHistogramAgg_WithFormat(t *testing.T) {
	t.Parallel()

	def := query.DateHistogramAgg(
		"by_year",
		FieldDate,
		calendarinterval.Year,
		query.DateHistogramAggFormat("yyyy"),
	)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["by_year"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.DateHistogram != nil)
	assert.Equal(t, "yyyy", *agg.DateHistogram.Format)
}

func TestHistogramAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.HistogramAgg("price_ranges", FieldPrice, 100.0)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["price_ranges"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Histogram != nil)
	assert.Equal(t, string(FieldPrice), *agg.Histogram.Field)
	assert.Assert(t, math.Abs(float64(*agg.Histogram.Interval)-100.0) < 0.001)
}

func TestAvgAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["avg_price"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Avg != nil)
	assert.Equal(t, string(FieldPrice), *agg.Avg.Field)
}

func TestMaxAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.MaxAgg("max_price", FieldPrice)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["max_price"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Max != nil)
	assert.Equal(t, string(FieldPrice), *agg.Max.Field)
}

func TestMinAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.MinAgg("min_price", FieldPrice)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["min_price"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Min != nil)
	assert.Equal(t, string(FieldPrice), *agg.Min.Field)
}

func TestSumAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.SumAgg("sum_value", FieldValue)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["sum_value"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Sum != nil)
	assert.Equal(t, string(FieldValue), *agg.Sum.Field)
}

func TestValueCountAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.ValueCountAgg("order_count", FieldId)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["order_count"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.ValueCount != nil)
	assert.Equal(t, string(FieldId), *agg.ValueCount.Field)
}

func TestCardinalityAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.CardinalityAgg("unique_users", FieldId)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["unique_users"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Cardinality != nil)
	assert.Equal(t, string(FieldId), *agg.Cardinality.Field)
}

func TestStatsAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.StatsAgg("price_stats", FieldPrice)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["price_stats"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Stats != nil)
	assert.Equal(t, string(FieldPrice), *agg.Stats.Field)
}

func TestStringTermsAgg_WithSubAggs(t *testing.T) {
	t.Parallel()

	avgDef := query.AvgAgg("avg_price", FieldPrice)
	sumDef := query.SumAgg("sum_value", FieldPrice)
	def := query.StringTermsAgg(
		"by_category",
		FieldCategory,
		query.TermsAggSubAggs(avgDef, sumDef),
	)

	aggs := query.Aggs(def).Build()

	agg, ok := aggs["by_category"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Terms != nil)
	assert.Assert(t, agg.Aggregations != nil)

	_, ok = agg.Aggregations["avg_price"]
	assert.Assert(t, ok)
	_, ok = agg.Aggregations["sum_value"]
	assert.Assert(t, ok)
}

func TestDateHistogramAgg_WithSubAggs(t *testing.T) {
	t.Parallel()

	avgDef := query.AvgAgg("avg_price", FieldPrice)
	def := query.DateHistogramAgg(
		"by_month",
		FieldDate,
		calendarinterval.Month,
		query.DateHistogramAggSubAggs(avgDef),
	)

	aggs := query.Aggs(def).Build()

	agg, ok := aggs["by_month"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.DateHistogram != nil)
	assert.Assert(t, agg.Aggregations != nil)

	_, ok = agg.Aggregations["avg_price"]
	assert.Assert(t, ok)
}

func TestHistogramAgg_WithSubAggs(t *testing.T) {
	t.Parallel()

	sumDef := query.SumAgg("sum_value", FieldValue)
	def := query.HistogramAgg(
		"price_ranges",
		FieldPrice,
		50.0,
		query.HistogramAggSubAggs(sumDef),
	)

	aggs := query.Aggs(def).Build()

	agg, ok := aggs["price_ranges"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Histogram != nil)
	assert.Assert(t, agg.Aggregations != nil)

	_, ok = agg.Aggregations["sum_value"]
	assert.Assert(t, ok)
}

func TestAggResults_Avg(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)
	value := 42.5
	raw := map[string]types.Aggregate{
		"avg_price": &types.AvgAggregate{Value: (*types.Float64)(&value)},
	}

	res, err := query.NewAggResults(raw).GetAvg(def)
	assert.NilError(t, err)

	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, value, *res.Value())
}

func TestAggResults_Sum(t *testing.T) {
	t.Parallel()

	def := query.SumAgg("sum_value", FieldValue)
	value := 100.0
	raw := map[string]types.Aggregate{
		"sum_value": &types.SumAggregate{Value: (*types.Float64)(&value)},
	}

	res, err := query.NewAggResults(raw).GetSum(def)
	assert.NilError(t, err)

	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, value, *res.Value())
}

func TestAggResults_Cardinality(t *testing.T) {
	t.Parallel()

	def := query.CardinalityAgg("unique_users", FieldId)
	value := int64(7)
	raw := map[string]types.Aggregate{
		"unique_users": &types.CardinalityAggregate{Value: value},
	}

	res, err := query.NewAggResults(raw).GetCardinality(def)
	assert.NilError(t, err)

	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, value, *res.Value())
}

func TestAggResults_Stats(t *testing.T) {
	t.Parallel()

	def := query.StatsAgg("price_stats", FieldPrice)
	min := 10.0
	max := 50.0
	avg := 30.0
	sum := 90.0
	raw := map[string]types.Aggregate{
		"price_stats": &types.StatsAggregate{
			Count: 3,
			Min:   (*types.Float64)(&min),
			Max:   (*types.Float64)(&max),
			Avg:   (*types.Float64)(&avg),
			Sum:   types.Float64(sum),
		},
	}

	res, err := query.NewAggResults(raw).GetStats(def)
	assert.NilError(t, err)

	assert.Equal(t, int64(3), res.Count())
	assert.Assert(t, res.Min() != nil)
	assert.Assert(t, res.Max() != nil)
	assert.Assert(t, res.Avg() != nil)
	assert.Assert(t, res.Sum() != nil)
	assert.Equal(t, min, *res.Min())
	assert.Equal(t, max, *res.Max())
	assert.Equal(t, avg, *res.Avg())
	assert.Equal(t, sum, *res.Sum())
}

func TestAggResults_StringTerms(t *testing.T) {
	t.Parallel()

	avgDef := query.AvgAgg("avg_price", FieldPrice)
	def := query.StringTermsAgg("by_category", FieldCategory, query.TermsAggSubAggs(avgDef))

	avgValue := 25.0
	raw := map[string]types.Aggregate{
		"by_category": &types.StringTermsAggregate{
			Buckets: []types.StringTermsBucket{
				{
					Key:      "a",
					DocCount: 2,
					Aggregations: map[string]types.Aggregate{
						"avg_price": &types.AvgAggregate{Value: (*types.Float64)(&avgValue)},
					},
				},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetStringTerms(def)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(res.Buckets()))
	assert.Equal(t, "a", res.Buckets()[0].Key())
	assert.Equal(t, int64(2), res.Buckets()[0].DocCount())

	avgRes, err := res.Buckets()[0].Aggregations().GetAvg(avgDef)
	assert.NilError(t, err)
	assert.Assert(t, avgRes.Value() != nil)
	assert.Equal(t, avgValue, *avgRes.Value())
}

func TestAggResults_DateHistogram(t *testing.T) {
	t.Parallel()

	sumDef := query.SumAgg("sum_value", FieldValue)
	def := query.DateHistogramAgg(
		"by_month",
		FieldDate,
		calendarinterval.Month,
		query.DateHistogramAggSubAggs(sumDef),
	)

	sumValue := 120.0
	keyAsString := "2026-03-01T00:00:00.000Z"
	raw := map[string]types.Aggregate{
		"by_month": &types.DateHistogramAggregate{
			Buckets: []types.DateHistogramBucket{
				{
					Key:         1740787200000,
					KeyAsString: &keyAsString,
					DocCount:    4,
					Aggregations: map[string]types.Aggregate{
						"sum_value": &types.SumAggregate{Value: (*types.Float64)(&sumValue)},
					},
				},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetDateHistogram(def)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(res.Buckets()))
	assert.Equal(t, int64(1740787200000), res.Buckets()[0].Key())
	assert.Equal(t, keyAsString, res.Buckets()[0].KeyAsString())
	assert.Equal(t, int64(4), res.Buckets()[0].DocCount())
}

func TestAggResults_Histogram(t *testing.T) {
	t.Parallel()

	def := query.HistogramAgg("price_ranges", FieldPrice, 100.0)
	raw := map[string]types.Aggregate{
		"price_ranges": &types.HistogramAggregate{
			Buckets: []types.HistogramBucket{
				{
					Key:      100.0,
					DocCount: 5,
				},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetHistogram(def)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(res.Buckets()))
	assert.Equal(t, 100.0, res.Buckets()[0].Key())
	assert.Equal(t, int64(5), res.Buckets()[0].DocCount())
}

func TestAggResults_MissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)
	_, err := query.NewAggResults(nil).GetAvg(def)
	assert.Assert(t, err != nil)
}
