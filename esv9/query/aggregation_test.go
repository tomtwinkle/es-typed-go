package query_test

import (
	"fmt"
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
	assert.ErrorContains(t, err, `aggregation "avg_price" (avg) not found`)
}

func TestAggResults_Raw(t *testing.T) {
	t.Parallel()

	raw := map[string]types.Aggregate{
		"avg_price": &types.AvgAggregate{},
	}

	results := query.NewAggResults(raw)
	assert.DeepEqual(t, raw, results.Raw())
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

func TestMustAgg_GenericHelper_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "avg_price" (avg) not found`)
		}()

		_ = query.MustAgg(query.NewAggResults(nil), def)
		return false
	}())
}

func TestMustAvg_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "avg_price" (avg) not found`)
		}()

		_ = query.NewAggResults(nil).MustAvg(def)
		return false
	}())
}

func TestMustSum_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.SumAgg("sum_value", FieldValue)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "sum_value" (sum) not found`)
		}()

		_ = query.NewAggResults(nil).MustSum(def)
		return false
	}())
}

func TestMustMin_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.MinAgg("min_price", FieldPrice)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "min_price" (min) not found`)
		}()

		_ = query.NewAggResults(nil).MustMin(def)
		return false
	}())
}

func TestMustMax_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.MaxAgg("max_price", FieldPrice)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "max_price" (max) not found`)
		}()

		_ = query.NewAggResults(nil).MustMax(def)
		return false
	}())
}

func TestMustStats_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.StatsAgg("price_stats", FieldPrice)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "price_stats" (stats) not found`)
		}()

		_ = query.NewAggResults(nil).MustStats(def)
		return false
	}())
}

func TestMustValueCount_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.ValueCountAgg("order_count", FieldId)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "order_count" (value_count) not found`)
		}()

		_ = query.NewAggResults(nil).MustValueCount(def)
		return false
	}())
}

func TestMustCardinality_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.CardinalityAgg("unique_users", FieldId)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "unique_users" (cardinality) not found`)
		}()

		_ = query.NewAggResults(nil).MustCardinality(def)
		return false
	}())
}

func TestMustStringTerms_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("by_category", FieldCategory)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "by_category" (string_terms) not found`)
		}()

		_ = query.NewAggResults(nil).MustStringTerms(def)
		return false
	}())
}

func TestMustDateHistogram_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "by_month" (date_histogram) not found`)
		}()

		_ = query.NewAggResults(nil).MustDateHistogram(def)
		return false
	}())
}

func TestMustHistogram_PanicsOnMissingAggregation(t *testing.T) {
	t.Parallel()

	def := query.HistogramAgg("price_ranges", FieldPrice, 100.0)

	assert.Assert(t, func() (panicked bool) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			panicked = true
			err, ok := r.(error)
			assert.Assert(t, ok)
			assert.ErrorContains(t, err, `aggregation "price_ranges" (histogram) not found`)
		}()

		_ = query.NewAggResults(nil).MustHistogram(def)
		return false
	}())
}

func TestAggResults_UnexpectedType(t *testing.T) {
	t.Parallel()

	def := query.AvgAgg("avg_price", FieldPrice)
	value := 1.0
	raw := map[string]types.Aggregate{
		"avg_price": &types.SumAggregate{Value: (*types.Float64)(&value)},
	}

	_, err := query.NewAggResults(raw).GetAvg(def)
	assert.ErrorContains(t, err, `aggregation "avg_price" is not of expected type avg`)
}

func TestAggResults_StringTerms_UnexpectedBucketsType(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("by_category", FieldCategory)
	raw := map[string]types.Aggregate{
		"by_category": &types.StringTermsAggregate{
			Buckets: map[string]types.StringTermsBucket{},
		},
	}

	_, err := query.NewAggResults(raw).GetStringTerms(def)
	assert.ErrorContains(t, err, `aggregation "by_category" returned unexpected buckets type for string terms`)
}

func TestAggResults_StringTerms_NonStringBucketKey(t *testing.T) {
	t.Parallel()

	def := query.StringTermsAgg("by_category", FieldCategory)
	raw := map[string]types.Aggregate{
		"by_category": &types.StringTermsAggregate{
			Buckets: []types.StringTermsBucket{
				{Key: 123, DocCount: 1},
			},
		},
	}

	_, err := query.NewAggResults(raw).GetStringTerms(def)
	assert.ErrorContains(t, err, `aggregation "by_category" returned non-string bucket key of type int`)
}

func TestAggResults_StringTerms_BucketNilAggregations(t *testing.T) {
	t.Parallel()

	avgDef := query.AvgAgg("avg_price", FieldPrice)
	def := query.StringTermsAgg("by_category", FieldCategory, query.TermsAggSubAggs(avgDef))
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

	_, err = res.Buckets()[0].Aggregations().GetAvg(avgDef)
	assert.ErrorContains(t, err, `aggregation "avg_price" (avg) not found`)
}

func TestAggResults_DateHistogram_UnexpectedBucketsType(t *testing.T) {
	t.Parallel()

	def := query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month)
	raw := map[string]types.Aggregate{
		"by_month": &types.DateHistogramAggregate{
			Buckets: map[string]types.DateHistogramBucket{},
		},
	}

	_, err := query.NewAggResults(raw).GetDateHistogram(def)
	assert.ErrorContains(t, err, `aggregation "by_month" returned unexpected buckets type for date histogram`)
}

func TestAggResults_DateHistogram_NilKeyAsString(t *testing.T) {
	t.Parallel()

	def := query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month)
	raw := map[string]types.Aggregate{
		"by_month": &types.DateHistogramAggregate{
			Buckets: []types.DateHistogramBucket{
				{Key: 1740787200000, KeyAsString: nil, DocCount: 4, Aggregations: nil},
			},
		},
	}

	res, err := query.NewAggResults(raw).GetDateHistogram(def)
	assert.NilError(t, err)
	assert.Equal(t, 1, len(res.Buckets()))
	assert.Equal(t, "", res.Buckets()[0].KeyAsString())
	assert.Assert(t, res.Buckets()[0].Aggregations().Raw() == nil)
}

func TestAggResults_Histogram_UnexpectedBucketsType(t *testing.T) {
	t.Parallel()

	def := query.HistogramAgg("price_ranges", FieldPrice, 100.0)
	raw := map[string]types.Aggregate{
		"price_ranges": &types.HistogramAggregate{
			Buckets: map[string]types.HistogramBucket{},
		},
	}

	_, err := query.NewAggResults(raw).GetHistogram(def)
	assert.ErrorContains(t, err, `aggregation "price_ranges" returned unexpected buckets type for histogram`)
}

func TestTermsAggOptionAliases(t *testing.T) {
	t.Parallel()

	defWithAlias := query.StringTermsAgg(
		"by_category",
		FieldCategory,
		query.WithTermsSize(10),
		query.WithSubAggs(query.AvgAgg("avg_price", FieldPrice)),
	)
	defWithCanonical := query.StringTermsAgg(
		"by_category",
		FieldCategory,
		query.TermsAggSize(10),
		query.TermsAggSubAggs(query.AvgAgg("avg_price", FieldPrice)),
	)

	aliasAggs := query.Aggs(defWithAlias).Build()
	canonicalAggs := query.Aggs(defWithCanonical).Build()

	assert.DeepEqual(t, aliasAggs, canonicalAggs)
}

func TestAggResults_PointerConversions_NilAndZeroValues(t *testing.T) {
	t.Parallel()

	var nilFloat *types.Float64
	zero := 0.0

	avgDef := query.AvgAgg("avg_price", FieldPrice)
	sumDef := query.SumAgg("sum_value", FieldValue)
	minDef := query.MinAgg("min_price", FieldPrice)
	maxDef := query.MaxAgg("max_price", FieldPrice)
	valueCountDef := query.ValueCountAgg("order_count", FieldId)
	cardinalityDef := query.CardinalityAgg("unique_users", FieldId)
	statsDef := query.StatsAgg("price_stats", FieldPrice)

	raw := map[string]types.Aggregate{
		"avg_price":    &types.AvgAggregate{Value: nilFloat},
		"sum_value":    &types.SumAggregate{Value: (*types.Float64)(&zero)},
		"min_price":    &types.MinAggregate{Value: nilFloat},
		"max_price":    &types.MaxAggregate{Value: nilFloat},
		"order_count":  &types.ValueCountAggregate{Value: nilFloat},
		"unique_users": &types.CardinalityAggregate{Value: 0},
		"price_stats": &types.StatsAggregate{
			Count: 0,
			Min:   nilFloat,
			Max:   nilFloat,
			Avg:   nilFloat,
			Sum:   0,
		},
	}

	results := query.NewAggResults(raw)

	avgRes, err := results.GetAvg(avgDef)
	assert.NilError(t, err)
	assert.Assert(t, avgRes.Value() == nil)

	sumRes, err := results.GetSum(sumDef)
	assert.NilError(t, err)
	assert.Assert(t, sumRes.Value() != nil)
	assert.Equal(t, 0.0, *sumRes.Value())

	minRes, err := results.GetMin(minDef)
	assert.NilError(t, err)
	assert.Assert(t, minRes.Value() == nil)

	maxRes, err := results.GetMax(maxDef)
	assert.NilError(t, err)
	assert.Assert(t, maxRes.Value() == nil)

	valueCountRes, err := results.GetValueCount(valueCountDef)
	assert.NilError(t, err)
	assert.Assert(t, valueCountRes.Value() == nil)

	cardinalityRes, err := results.GetCardinality(cardinalityDef)
	assert.NilError(t, err)
	assert.Assert(t, cardinalityRes.Value() != nil)
	assert.Equal(t, int64(0), *cardinalityRes.Value())

	statsRes, err := results.GetStats(statsDef)
	assert.NilError(t, err)
	assert.Equal(t, int64(0), statsRes.Count())
	assert.Assert(t, statsRes.Min() == nil)
	assert.Assert(t, statsRes.Max() == nil)
	assert.Assert(t, statsRes.Avg() == nil)
	assert.Assert(t, statsRes.Sum() != nil)
	assert.Equal(t, 0.0, *statsRes.Sum())
}

func TestMustHelpers_PanicValuesAreErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		run  func()
		want string
	}{
		{
			name: "must avg",
			run: func() {
				_ = query.NewAggResults(nil).MustAvg(query.AvgAgg("avg_price", FieldPrice))
			},
			want: `aggregation "avg_price" (avg) not found`,
		},
		{
			name: "must sum",
			run: func() {
				_ = query.NewAggResults(nil).MustSum(query.SumAgg("sum_value", FieldValue))
			},
			want: `aggregation "sum_value" (sum) not found`,
		},
		{
			name: "must min",
			run: func() {
				_ = query.NewAggResults(nil).MustMin(query.MinAgg("min_price", FieldPrice))
			},
			want: `aggregation "min_price" (min) not found`,
		},
		{
			name: "must max",
			run: func() {
				_ = query.NewAggResults(nil).MustMax(query.MaxAgg("max_price", FieldPrice))
			},
			want: `aggregation "max_price" (max) not found`,
		},
		{
			name: "must stats",
			run: func() {
				_ = query.NewAggResults(nil).MustStats(query.StatsAgg("price_stats", FieldPrice))
			},
			want: `aggregation "price_stats" (stats) not found`,
		},
		{
			name: "must value count",
			run: func() {
				_ = query.NewAggResults(nil).MustValueCount(query.ValueCountAgg("order_count", FieldId))
			},
			want: `aggregation "order_count" (value_count) not found`,
		},
		{
			name: "must cardinality",
			run: func() {
				_ = query.NewAggResults(nil).MustCardinality(query.CardinalityAgg("unique_users", FieldId))
			},
			want: `aggregation "unique_users" (cardinality) not found`,
		},
		{
			name: "must string terms",
			run: func() {
				_ = query.NewAggResults(nil).MustStringTerms(query.StringTermsAgg("by_category", FieldCategory))
			},
			want: `aggregation "by_category" (string_terms) not found`,
		},
		{
			name: "must date histogram",
			run: func() {
				_ = query.NewAggResults(nil).MustDateHistogram(query.DateHistogramAgg("by_month", FieldDate, calendarinterval.Month))
			},
			want: `aggregation "by_month" (date_histogram) not found`,
		},
		{
			name: "must histogram",
			run: func() {
				_ = query.NewAggResults(nil).MustHistogram(query.HistogramAgg("price_ranges", FieldPrice, 100.0))
			},
			want: `aggregation "price_ranges" (histogram) not found`,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				r := recover()
				assert.Assert(t, r != nil)

				err, ok := r.(error)
				assert.Assert(t, ok, fmt.Sprintf("panic value should be error, got %T", r))
				assert.ErrorContains(t, err, tt.want)
			}()

			tt.run()
		})
	}
}
