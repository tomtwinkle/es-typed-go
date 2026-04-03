package query_test

import (
	"math"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/query"
)

func assertPanicsWithErrorContains(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		assert.Assert(t, recovered != nil)

		err, ok := recovered.(error)
		assert.Assert(t, ok)
		assert.ErrorContains(t, err, want)
	}()

	fn()
}

func TestAggs_Empty(t *testing.T) {
	t.Parallel()

	aggs := query.Aggs().Build()
	assert.Assert(t, aggs == nil)
}

func TestAggResults_NewAggResultsRaw(t *testing.T) {
	t.Parallel()

	raw := map[string]types.Aggregate{
		"avg_price": &types.AvgAggregate{},
	}

	results := query.NewAggResults(raw)
	assert.DeepEqual(t, results.Raw(), raw)
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

func TestTermsAgg_AliasBuild(t *testing.T) {
	t.Parallel()

	def := query.TermsAgg("top5", FieldCategory, query.WithTermsSize(5))
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["top5"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Terms != nil)
	assert.Equal(t, string(FieldCategory), *agg.Terms.Field)
	assert.Equal(t, 5, *agg.Terms.Size)
}

func TestTermsAgg_AliasWithSubAggs(t *testing.T) {
	t.Parallel()

	avgDef := query.AvgAgg("avg_price", FieldPrice)
	def := query.TermsAgg("by_category", FieldCategory, query.WithSubAggs(avgDef))

	aggs := query.Aggs(def).Build()
	agg, ok := aggs["by_category"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Aggregations != nil)

	_, ok = agg.Aggregations["avg_price"]
	assert.Assert(t, ok)
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

func TestAggResults_AvgNilValue(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)
	raw := map[string]types.Aggregate{
		"avg_price": &types.AvgAggregate{Value: nil},
	}

	res, err := query.NewAggResults(raw).GetAvg(def)
	assert.NilError(t, err)
	assert.Assert(t, res.Value() == nil)
}

func TestAggResults_Sum(t *testing.T) {
	t.Parallel()

	def := query.SumAgg("sum_value", FieldValue)
	value := 123.4
	raw := map[string]types.Aggregate{
		"sum_value": &types.SumAggregate{Value: (*types.Float64)(&value)},
	}

	res, err := query.NewAggResults(raw).GetSum(def)
	assert.NilError(t, err)

	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, value, *res.Value())
}

func TestAggResults_SumNilValue(t *testing.T) {
	t.Parallel()

	def := query.SumAgg("sum_value", FieldValue)
	raw := map[string]types.Aggregate{
		"sum_value": &types.SumAggregate{Value: nil},
	}

	res, err := query.NewAggResults(raw).GetSum(def)
	assert.NilError(t, err)
	assert.Assert(t, res.Value() == nil)
}

func TestAggResults_Min(t *testing.T) {
	t.Parallel()

	def := query.MinAgg("min_price", FieldPrice)
	value := 10.5
	raw := map[string]types.Aggregate{
		"min_price": &types.MinAggregate{Value: (*types.Float64)(&value)},
	}

	res, err := query.NewAggResults(raw).GetMin(def)
	assert.NilError(t, err)

	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, value, *res.Value())
}

func TestAggResults_MinNilValue(t *testing.T) {
	t.Parallel()

	def := query.MinAgg("min_price", FieldPrice)
	raw := map[string]types.Aggregate{
		"min_price": &types.MinAggregate{Value: nil},
	}

	res, err := query.NewAggResults(raw).GetMin(def)
	assert.NilError(t, err)
	assert.Assert(t, res.Value() == nil)
}

func TestAggResults_Max(t *testing.T) {
	t.Parallel()

	def := query.MaxAgg("max_price", FieldPrice)
	value := 99.9
	raw := map[string]types.Aggregate{
		"max_price": &types.MaxAggregate{Value: (*types.Float64)(&value)},
	}

	res, err := query.NewAggResults(raw).GetMax(def)
	assert.NilError(t, err)

	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, value, *res.Value())
}

func TestAggResults_MaxNilValue(t *testing.T) {
	t.Parallel()

	def := query.MaxAgg("max_price", FieldPrice)
	raw := map[string]types.Aggregate{
		"max_price": &types.MaxAggregate{Value: nil},
	}

	res, err := query.NewAggResults(raw).GetMax(def)
	assert.NilError(t, err)
	assert.Assert(t, res.Value() == nil)
}

func TestAggResults_ValueCount(t *testing.T) {
	t.Parallel()

	def := query.ValueCountAgg("order_count", FieldId)
	value := 7
	raw := map[string]types.Aggregate{
		"order_count": &types.ValueCountAggregate{Value: (*types.Float64)(&[]float64{float64(value)}[0])},
	}

	res, err := query.NewAggResults(raw).GetValueCount(def)
	assert.NilError(t, err)

	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, int64(7), *res.Value())
}

func TestAggResults_ValueCountNilValue(t *testing.T) {
	t.Parallel()

	def := query.ValueCountAgg("order_count", FieldId)
	raw := map[string]types.Aggregate{
		"order_count": &types.ValueCountAggregate{Value: nil},
	}

	res, err := query.NewAggResults(raw).GetValueCount(def)
	assert.NilError(t, err)
	assert.Assert(t, res.Value() == nil)
}

func TestAggResults_Cardinality(t *testing.T) {
	t.Parallel()

	def := query.CardinalityAgg("unique_users", FieldId)
	value := 5
	raw := map[string]types.Aggregate{
		"unique_users": &types.CardinalityAggregate{Value: int64(value)},
	}

	res, err := query.NewAggResults(raw).GetCardinality(def)
	assert.NilError(t, err)

	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, int64(5), *res.Value())
}

func TestAggResults_Stats(t *testing.T) {
	t.Parallel()

	def := query.StatsAgg("price_stats", FieldPrice)
	min := 10.0
	max := 30.0
	avg := 20.0
	sum := 60.0
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

	def := query.StringTermsAgg("by_category", FieldCategory)
	raw := map[string]types.Aggregate{
		"by_category": &types.StringTermsAggregate{
			Buckets: []types.StringTermsBucket{
				{Key: "electronics", DocCount: 2},
				{Key: "clothing", DocCount: 3},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetStringTerms(def)
	assert.NilError(t, err)

	buckets := res.Buckets()
	assert.Equal(t, 2, len(buckets))
	assert.Equal(t, "electronics", buckets[0].Key())
	assert.Equal(t, int64(2), buckets[0].DocCount())
	assert.Equal(t, "clothing", buckets[1].Key())
	assert.Equal(t, int64(3), buckets[1].DocCount())
}

func TestAggResults_StringTermsNilBucketAggs(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("by_category", FieldCategory)
	raw := map[string]types.Aggregate{
		"by_category": &types.StringTermsAggregate{
			Buckets: []types.StringTermsBucket{
				{Key: "electronics", DocCount: 2, Aggregations: nil},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetStringTerms(def)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(res.Buckets()))
	assert.Assert(t, res.Buckets()[0].Aggregations().Raw() == nil)
}

func TestAggResults_StringTermsBucketKeyTypeError(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("by_category", FieldCategory)
	raw := map[string]types.Aggregate{
		"by_category": &types.StringTermsAggregate{
			Buckets: []types.StringTermsBucket{
				{Key: 123, DocCount: 2},
			},
		},
	}

	_, err := query.NewAggResults(raw).GetStringTerms(def)
	assert.ErrorContains(t, err, `aggregation "by_category" has unexpected bucket key type int`)
}

func TestAggResults_StringTermsBucketsTypeError(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("by_category", FieldCategory)
	raw := map[string]types.Aggregate{
		"by_category": &types.StringTermsAggregate{
			Buckets: map[string]types.StringTermsBucket{},
		},
	}

	_, err := query.NewAggResults(raw).GetStringTerms(def)
	assert.ErrorContains(t, err, `aggregation "by_category" has unexpected buckets type`)
}

func TestAggResults_DateHistogram(t *testing.T) {
	t.Parallel()

	def := query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month)
	raw := map[string]types.Aggregate{
		"by_month": &types.DateHistogramAggregate{
			Buckets: []types.DateHistogramBucket{
				{Key: 1704067200000, KeyAsString: func() *string { s := "2024-01"; return &s }(), DocCount: 2},
				{Key: 1706745600000, KeyAsString: func() *string { s := "2024-02"; return &s }(), DocCount: 1},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetDateHistogram(def)
	assert.NilError(t, err)

	buckets := res.Buckets()
	assert.Equal(t, 2, len(buckets))
	assert.Equal(t, int64(1704067200000), buckets[0].Key())
	assert.Equal(t, "2024-01", buckets[0].KeyAsString())
	assert.Equal(t, int64(2), buckets[0].DocCount())
}

func TestAggResults_DateHistogramNilKeyAsString(t *testing.T) {
	t.Parallel()

	def := query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month)
	raw := map[string]types.Aggregate{
		"by_month": &types.DateHistogramAggregate{
			Buckets: []types.DateHistogramBucket{
				{Key: 1704067200000, KeyAsString: nil, DocCount: 2, Aggregations: nil},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetDateHistogram(def)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(res.Buckets()))
	assert.Equal(t, "", res.Buckets()[0].KeyAsString())
	assert.Assert(t, res.Buckets()[0].Aggregations().Raw() == nil)
}

func TestAggResults_DateHistogramBucketsTypeError(t *testing.T) {
	t.Parallel()

	def := query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month)
	raw := map[string]types.Aggregate{
		"by_month": &types.DateHistogramAggregate{
			Buckets: map[string]types.DateHistogramBucket{},
		},
	}

	_, err := query.NewAggResults(raw).GetDateHistogram(def)
	assert.ErrorContains(t, err, `aggregation "by_month" has unexpected buckets type`)
}

func TestAggResults_Histogram(t *testing.T) {
	t.Parallel()

	def := query.HistogramAgg("price_ranges", FieldPrice, 50.0)
	k1 := types.Float64(0)
	k2 := types.Float64(50)
	raw := map[string]types.Aggregate{
		"price_ranges": &types.HistogramAggregate{
			Buckets: []types.HistogramBucket{
				{Key: k1, DocCount: 4},
				{Key: k2, DocCount: 2},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetHistogram(def)
	assert.NilError(t, err)

	buckets := res.Buckets()
	assert.Equal(t, 2, len(buckets))
	assert.Equal(t, 0.0, buckets[0].Key())
	assert.Equal(t, int64(4), buckets[0].DocCount())
	assert.Equal(t, 50.0, buckets[1].Key())
	assert.Equal(t, int64(2), buckets[1].DocCount())
}

func TestAggResults_HistogramNilBucketAggs(t *testing.T) {
	t.Parallel()

	def := query.HistogramAgg("price_ranges", FieldPrice, 50.0)
	k1 := types.Float64(0)
	raw := map[string]types.Aggregate{
		"price_ranges": &types.HistogramAggregate{
			Buckets: []types.HistogramBucket{
				{Key: k1, DocCount: 4, Aggregations: nil},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetHistogram(def)
	assert.NilError(t, err)

	assert.Equal(t, 1, len(res.Buckets()))
	assert.Assert(t, res.Buckets()[0].Aggregations().Raw() == nil)
}

func TestAggResults_HistogramBucketsTypeError(t *testing.T) {
	t.Parallel()

	def := query.HistogramAgg("price_ranges", FieldPrice, 50.0)
	raw := map[string]types.Aggregate{
		"price_ranges": &types.HistogramAggregate{
			Buckets: map[string]types.HistogramBucket{},
		},
	}

	_, err := query.NewAggResults(raw).GetHistogram(def)
	assert.ErrorContains(t, err, `aggregation "price_ranges" has unexpected buckets type`)
}

func TestGetAgg_GenericHelper(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)
	value := 12.5
	raw := map[string]types.Aggregate{
		"avg_price": &types.AvgAggregate{Value: (*types.Float64)(&value)},
	}

	res, err := query.GetAgg(query.NewAggResults(raw), def)
	assert.NilError(t, err)
	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, value, *res.Value())
}

func TestMustAgg_GenericHelper(t *testing.T) {
	t.Parallel()

	def := query.SumAgg("sum_value", FieldValue)
	value := 55.0
	raw := map[string]types.Aggregate{
		"sum_value": &types.SumAggregate{Value: (*types.Float64)(&value)},
	}

	res := query.MustAgg(query.NewAggResults(raw), def)
	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, value, *res.Value())
}

func TestMustAgg_GenericHelperPanics(t *testing.T) {
	t.Parallel()

	def := query.SumAgg("sum_value", FieldValue)

	assertPanicsWithErrorContains(t, `aggregation "sum_value" not found`, func() {
		_ = query.MustAgg(query.NewAggResults(nil), def)
	})
}

func TestMustAvg(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)
	value := 33.3
	raw := map[string]types.Aggregate{
		"avg_price": &types.AvgAggregate{Value: (*types.Float64)(&value)},
	}

	res := query.NewAggResults(raw).MustAvg(def)
	assert.Assert(t, res.Value() != nil)
	assert.Equal(t, value, *res.Value())
}

func TestMustAvgPanics(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)

	assertPanicsWithErrorContains(t, `aggregation "avg_price" not found`, func() {
		_ = query.NewAggResults(nil).MustAvg(def)
	})
}

func TestMustSumPanics(t *testing.T) {
	t.Parallel()

	def := query.SumAgg("sum_value", FieldValue)

	assertPanicsWithErrorContains(t, `aggregation "sum_value" not found`, func() {
		_ = query.NewAggResults(nil).MustSum(def)
	})
}

func TestMustMinPanics(t *testing.T) {
	t.Parallel()

	def := query.MinAgg("min_price", FieldPrice)

	assertPanicsWithErrorContains(t, `aggregation "min_price" not found`, func() {
		_ = query.NewAggResults(nil).MustMin(def)
	})
}

func TestMustMaxPanics(t *testing.T) {
	t.Parallel()

	def := query.MaxAgg("max_price", FieldPrice)

	assertPanicsWithErrorContains(t, `aggregation "max_price" not found`, func() {
		_ = query.NewAggResults(nil).MustMax(def)
	})
}

func TestMustStatsPanics(t *testing.T) {
	t.Parallel()

	def := query.StatsAgg("price_stats", FieldPrice)

	assertPanicsWithErrorContains(t, `aggregation "price_stats" not found`, func() {
		_ = query.NewAggResults(nil).MustStats(def)
	})
}

func TestMustValueCountPanics(t *testing.T) {
	t.Parallel()

	def := query.ValueCountAgg("order_count", FieldId)

	assertPanicsWithErrorContains(t, `aggregation "order_count" not found`, func() {
		_ = query.NewAggResults(nil).MustValueCount(def)
	})
}

func TestMustCardinalityPanics(t *testing.T) {
	t.Parallel()

	def := query.CardinalityAgg("unique_users", FieldId)

	assertPanicsWithErrorContains(t, `aggregation "unique_users" not found`, func() {
		_ = query.NewAggResults(nil).MustCardinality(def)
	})
}

func TestMustStringTerms(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("by_category", FieldCategory)
	raw := map[string]types.Aggregate{
		"by_category": &types.StringTermsAggregate{
			Buckets: []types.StringTermsBucket{
				{Key: "electronics", DocCount: 2},
			},
		},
	}

	res := query.NewAggResults(raw).MustStringTerms(def)
	assert.Equal(t, 1, len(res.Buckets()))
	assert.Equal(t, "electronics", res.Buckets()[0].Key())
}

func TestMustStringTermsPanics(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("by_category", FieldCategory)

	assertPanicsWithErrorContains(t, `aggregation "by_category" not found`, func() {
		_ = query.NewAggResults(nil).MustStringTerms(def)
	})
}

func TestMustDateHistogramPanics(t *testing.T) {
	t.Parallel()

	def := query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month)

	assertPanicsWithErrorContains(t, `aggregation "by_month" not found`, func() {
		_ = query.NewAggResults(nil).MustDateHistogram(def)
	})
}

func TestMustHistogramPanics(t *testing.T) {
	t.Parallel()

	def := query.HistogramAgg("price_ranges", FieldPrice, 50.0)

	assertPanicsWithErrorContains(t, `aggregation "price_ranges" not found`, func() {
		_ = query.NewAggResults(nil).MustHistogram(def)
	})
}

func TestNestedAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.NestedAgg("items_agg", FieldItems)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["items_agg"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Nested != nil)
	assert.Equal(t, string(FieldItems), *agg.Nested.Path)
}

func TestNestedAgg_WithSubAggs(t *testing.T) {
	t.Parallel()

	sumDef := query.SumAgg("total_value", FieldValue)
	def := query.NestedAgg("items_agg", FieldItems, query.NestedAggSubAggs(sumDef))
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["items_agg"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Nested != nil)
	assert.Assert(t, agg.Aggregations != nil)
	_, ok = agg.Aggregations["total_value"]
	assert.Assert(t, ok)
}

func TestAggResults_Nested(t *testing.T) {
	t.Parallel()

	sumDef := query.SumAgg("total_value", FieldValue)
	def := query.NestedAgg("items_agg", FieldItems, query.NestedAggSubAggs(sumDef))
	v := 99.9
	raw := map[string]types.Aggregate{
		"items_agg": &types.NestedAggregate{
			DocCount: 5,
			Aggregations: map[string]types.Aggregate{
				"total_value": &types.SumAggregate{Value: (*types.Float64)(&v)},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetNested(def)
	assert.NilError(t, err)
	assert.Equal(t, int64(5), res.DocCount())

	sumRes, err := res.Aggregations().GetSum(sumDef)
	assert.NilError(t, err)
	assert.Assert(t, sumRes.Value() != nil)
	assert.Equal(t, v, *sumRes.Value())
}

func TestAggResults_NestedNilSubAggs(t *testing.T) {
	t.Parallel()

	def := query.NestedAgg("items_agg", FieldItems)
	raw := map[string]types.Aggregate{
		"items_agg": &types.NestedAggregate{DocCount: 3},
	}

	res, err := query.NewAggResults(raw).GetNested(def)
	assert.NilError(t, err)
	assert.Equal(t, int64(3), res.DocCount())
	assert.Assert(t, res.Aggregations().Raw() == nil)
}

func TestMustNested(t *testing.T) {
	t.Parallel()

	def := query.NestedAgg("items_agg", FieldItems)
	raw := map[string]types.Aggregate{
		"items_agg": &types.NestedAggregate{DocCount: 2},
	}

	res := query.NewAggResults(raw).MustNested(def)
	assert.Equal(t, int64(2), res.DocCount())
}

func TestMustNestedPanics(t *testing.T) {
	t.Parallel()

	def := query.NestedAgg("items_agg", FieldItems)
	assertPanicsWithErrorContains(t, `aggregation "items_agg" not found`, func() {
		_ = query.NewAggResults(nil).MustNested(def)
	})
}

func TestFilterAgg_Build(t *testing.T) {
	t.Parallel()

	filter := types.Query{Term: map[string]types.TermQuery{
		string(FieldStatus): {Value: "active"},
	}}
	def := query.FilterAgg("active_items", filter)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["active_items"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Filter != nil)
}

func TestFilterAgg_WithSubAggs(t *testing.T) {
	t.Parallel()

	filter := types.Query{Term: map[string]types.TermQuery{
		string(FieldStatus): {Value: "active"},
	}}
	sumDef := query.SumAgg("total_value", FieldValue)
	def := query.FilterAgg("active_items", filter, query.FilterAggSubAggs(sumDef))
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["active_items"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Filter != nil)
	assert.Assert(t, agg.Aggregations != nil)
	_, ok = agg.Aggregations["total_value"]
	assert.Assert(t, ok)
}

func TestAggResults_Filter(t *testing.T) {
	t.Parallel()

	filter := types.Query{Term: map[string]types.TermQuery{
		string(FieldStatus): {Value: "active"},
	}}
	sumDef := query.SumAgg("total_value", FieldValue)
	def := query.FilterAgg("active_items", filter, query.FilterAggSubAggs(sumDef))
	v := 42.0
	raw := map[string]types.Aggregate{
		"active_items": &types.FilterAggregate{
			DocCount: 7,
			Aggregations: map[string]types.Aggregate{
				"total_value": &types.SumAggregate{Value: (*types.Float64)(&v)},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetFilter(def)
	assert.NilError(t, err)
	assert.Equal(t, int64(7), res.DocCount())

	sumRes, err := res.Aggregations().GetSum(sumDef)
	assert.NilError(t, err)
	assert.Assert(t, sumRes.Value() != nil)
	assert.Equal(t, v, *sumRes.Value())
}

func TestMustFilterPanics(t *testing.T) {
	t.Parallel()

	filter := types.Query{}
	def := query.FilterAgg("active_items", filter)
	assertPanicsWithErrorContains(t, `aggregation "active_items" not found`, func() {
		_ = query.NewAggResults(nil).MustFilter(def)
	})
}

func TestMultiTermsAgg_Build(t *testing.T) {
	t.Parallel()

	def := query.MultiTermsAgg("by_category_status", []query.MultiTermLookup{{Field: FieldCategory}, {Field: FieldStatus}})
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["by_category_status"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.MultiTerms != nil)
	assert.Equal(t, 2, len(agg.MultiTerms.Terms))
	assert.Equal(t, string(FieldCategory), agg.MultiTerms.Terms[0].Field)
	assert.Equal(t, string(FieldStatus), agg.MultiTerms.Terms[1].Field)
}

func TestMultiTermsAgg_WithSize(t *testing.T) {
	t.Parallel()

	def := query.MultiTermsAgg(
		"by_category_status",
		[]query.MultiTermLookup{{Field: FieldCategory}, {Field: FieldStatus}},
		query.MultiTermsAggSize(100),
	)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["by_category_status"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.MultiTerms != nil)
	assert.Equal(t, 100, *agg.MultiTerms.Size)
}

func TestMultiTermsAgg_WithMissing(t *testing.T) {
	t.Parallel()

	def := query.MultiTermsAgg("by_category_status", []query.MultiTermLookup{
		{Field: FieldCategory},
		{Field: FieldStatus, Missing: "unknown"},
	})
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["by_category_status"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.MultiTerms != nil)
	assert.Equal(t, 2, len(agg.MultiTerms.Terms))
	assert.Equal(t, string(FieldCategory), agg.MultiTerms.Terms[0].Field)
	assert.Assert(t, agg.MultiTerms.Terms[0].Missing == nil)
	assert.Equal(t, string(FieldStatus), agg.MultiTerms.Terms[1].Field)
	assert.Equal(t, "unknown", agg.MultiTerms.Terms[1].Missing)
}

func TestMultiTermsAgg_WithSubAggs(t *testing.T) {
	t.Parallel()

	sumDef := query.SumAgg("total_value", FieldValue)
	def := query.MultiTermsAgg(
		"by_category_status",
		[]query.MultiTermLookup{{Field: FieldCategory}, {Field: FieldStatus}},
		query.MultiTermsAggSubAggs(sumDef),
	)
	aggs := query.Aggs(def).Build()

	agg, ok := aggs["by_category_status"]
	assert.Assert(t, ok)
	assert.Assert(t, agg.Aggregations != nil)
	_, ok = agg.Aggregations["total_value"]
	assert.Assert(t, ok)
}

func TestAggResults_MultiTerms(t *testing.T) {
	t.Parallel()

	def := query.MultiTermsAgg("by_category_status", []query.MultiTermLookup{{Field: FieldCategory}, {Field: FieldStatus}})
	raw := map[string]types.Aggregate{
		"by_category_status": &types.MultiTermsAggregate{
			Buckets: []types.MultiTermsBucket{
				{Key: []types.FieldValue{"electronics", "active"}, DocCount: 10},
				{Key: []types.FieldValue{"clothing", "inactive"}, DocCount: 3},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetMultiTerms(def)
	assert.NilError(t, err)

	buckets := res.Buckets()
	assert.Equal(t, 2, len(buckets))
	assert.Equal(t, 2, len(buckets[0].Keys()))
	assert.Equal(t, "electronics", buckets[0].Keys()[0])
	assert.Equal(t, "active", buckets[0].Keys()[1])
	assert.Equal(t, int64(10), buckets[0].DocCount())
	assert.Equal(t, "clothing", buckets[1].Keys()[0])
	assert.Equal(t, int64(3), buckets[1].DocCount())
}

func TestAggResults_MultiTermsWithSubAggs(t *testing.T) {
	t.Parallel()

	sumDef := query.SumAgg("total_value", FieldValue)
	def := query.MultiTermsAgg("by_category_status", []query.MultiTermLookup{{Field: FieldCategory}, {Field: FieldStatus}},
		query.MultiTermsAggSubAggs(sumDef),
	)
	v := 55.5
	raw := map[string]types.Aggregate{
		"by_category_status": &types.MultiTermsAggregate{
			Buckets: []types.MultiTermsBucket{
				{
					Key:      []types.FieldValue{"electronics", "active"},
					DocCount: 10,
					Aggregations: map[string]types.Aggregate{
						"total_value": &types.SumAggregate{Value: (*types.Float64)(&v)},
					},
				},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetMultiTerms(def)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(res.Buckets()))

	sumRes, err := res.Buckets()[0].Aggregations().GetSum(sumDef)
	assert.NilError(t, err)
	assert.Assert(t, sumRes.Value() != nil)
	assert.Equal(t, v, *sumRes.Value())
}

func TestAggResults_MultiTermsBucketsTypeError(t *testing.T) {
	t.Parallel()

	def := query.MultiTermsAgg("by_category_status", []query.MultiTermLookup{{Field: FieldCategory}, {Field: FieldStatus}})
	raw := map[string]types.Aggregate{
		"by_category_status": &types.MultiTermsAggregate{
			Buckets: map[string]types.MultiTermsBucket{},
		},
	}

	_, err := query.NewAggResults(raw).GetMultiTerms(def)
	assert.ErrorContains(t, err, `aggregation "by_category_status" has unexpected buckets type`)
}

func TestMustMultiTermsPanics(t *testing.T) {
	t.Parallel()

	def := query.MultiTermsAgg("by_category_status", []query.MultiTermLookup{{Field: FieldCategory}, {Field: FieldStatus}})
	assertPanicsWithErrorContains(t, `aggregation "by_category_status" not found`, func() {
		_ = query.NewAggResults(nil).MustMultiTerms(def)
	})
}

func TestAggResults_FilterNilSubAggs(t *testing.T) {
	t.Parallel()

	filter := types.Query{}
	def := query.FilterAgg("active_items", filter)
	raw := map[string]types.Aggregate{
		"active_items": &types.FilterAggregate{DocCount: 3},
	}

	res, err := query.NewAggResults(raw).GetFilter(def)
	assert.NilError(t, err)
	assert.Equal(t, int64(3), res.DocCount())
	assert.Assert(t, res.Aggregations().Raw() == nil)
}

func TestMustFilter(t *testing.T) {
	t.Parallel()

	filter := types.Query{}
	def := query.FilterAgg("active_items", filter)
	raw := map[string]types.Aggregate{
		"active_items": &types.FilterAggregate{DocCount: 7},
	}

	res := query.NewAggResults(raw).MustFilter(def)
	assert.Equal(t, int64(7), res.DocCount())
}

func TestAggResults_MultiTermsEmptyBuckets(t *testing.T) {
	t.Parallel()

	def := query.MultiTermsAgg("by_category_status", []query.MultiTermLookup{{Field: FieldCategory}, {Field: FieldStatus}})
	raw := map[string]types.Aggregate{
		"by_category_status": &types.MultiTermsAggregate{
			Buckets: []types.MultiTermsBucket{},
		},
	}

	res, err := query.NewAggResults(raw).GetMultiTerms(def)
	assert.NilError(t, err)
	assert.Equal(t, 0, len(res.Buckets()))
}

func TestGetAgg_Nested(t *testing.T) {
	t.Parallel()

	def := query.NestedAgg("items_agg", FieldItems)
	raw := map[string]types.Aggregate{
		"items_agg": &types.NestedAggregate{DocCount: 3},
	}

	res, err := query.GetAgg(query.NewAggResults(raw), def)
	assert.NilError(t, err)
	assert.Equal(t, int64(3), res.DocCount())
}

func TestGetAgg_Filter(t *testing.T) {
	t.Parallel()

	filter := types.Query{}
	def := query.FilterAgg("active_items", filter)
	raw := map[string]types.Aggregate{
		"active_items": &types.FilterAggregate{DocCount: 5},
	}

	res, err := query.GetAgg(query.NewAggResults(raw), def)
	assert.NilError(t, err)
	assert.Equal(t, int64(5), res.DocCount())
}

func TestGetAgg_MultiTerms(t *testing.T) {
	t.Parallel()

	def := query.MultiTermsAgg("by_category_status", []query.MultiTermLookup{{Field: FieldCategory}, {Field: FieldStatus}})
	raw := map[string]types.Aggregate{
		"by_category_status": &types.MultiTermsAggregate{
			Buckets: []types.MultiTermsBucket{
				{Key: []types.FieldValue{"electronics", "active"}, DocCount: 4},
			},
		},
	}

	res, err := query.GetAgg(query.NewAggResults(raw), def)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(res.Buckets()))
	assert.Equal(t, "electronics", res.Buckets()[0].Keys()[0])
	assert.Equal(t, "active", res.Buckets()[0].Keys()[1])
}

func TestAggResults_MissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("missing_avg", FieldPrice)
	_, err := query.NewAggResults(map[string]types.Aggregate{}).GetAvg(def)
	assert.ErrorContains(t, err, `aggregation "missing_avg" not found`)
}

func TestAggResults_NilRawMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)
	_, err := query.NewAggResults(nil).GetAvg(def)
	assert.ErrorContains(t, err, `aggregation "avg_price" not found`)
}

func TestAggResults_UnexpectedType(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)
	value := 1.0
	raw := map[string]types.Aggregate{
		"avg_price": &types.SumAggregate{Value: (*types.Float64)(&value)},
	}

	_, err := query.NewAggResults(raw).GetAvg(def)
	assert.ErrorContains(t, err, `aggregation "avg_price" has unexpected type`)
}
