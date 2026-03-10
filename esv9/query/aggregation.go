package query

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/calendarinterval"
)

// AggregationBuilder constructs Elasticsearch v9 aggregations using method chaining.
// Method chaining is intentionally allowed only in this package (query).
type AggregationBuilder struct {
	aggs map[string]types.Aggregations
}

// NewAggregations returns a new empty AggregationBuilder.
func NewAggregations() *AggregationBuilder {
	return &AggregationBuilder{
		aggs: make(map[string]types.Aggregations),
	}
}

// Terms adds a terms bucket aggregation on the given field.
func (b *AggregationBuilder) Terms(name, field string) *AggregationBuilder {
	agg := types.NewTermsAggregation()
	agg.Field = &field
	b.aggs[name] = types.Aggregations{Terms: agg}
	return b
}

// TermsWithSize adds a terms bucket aggregation with a custom size.
func (b *AggregationBuilder) TermsWithSize(name, field string, size int) *AggregationBuilder {
	agg := types.NewTermsAggregation()
	agg.Field = &field
	agg.Size = &size
	b.aggs[name] = types.Aggregations{Terms: agg}
	return b
}

// DateHistogram adds a date_histogram bucket aggregation using a calendar interval.
func (b *AggregationBuilder) DateHistogram(name, field string, interval calendarinterval.CalendarInterval) *AggregationBuilder {
	agg := types.NewDateHistogramAggregation()
	agg.Field = &field
	agg.CalendarInterval = &interval
	b.aggs[name] = types.Aggregations{DateHistogram: agg}
	return b
}

// DateHistogramWithFormat adds a date_histogram aggregation with a date format.
func (b *AggregationBuilder) DateHistogramWithFormat(name, field, format string, interval calendarinterval.CalendarInterval) *AggregationBuilder {
	agg := types.NewDateHistogramAggregation()
	agg.Field = &field
	agg.CalendarInterval = &interval
	agg.Format = &format
	b.aggs[name] = types.Aggregations{DateHistogram: agg}
	return b
}

// Histogram adds a histogram bucket aggregation with a numeric interval.
func (b *AggregationBuilder) Histogram(name, field string, interval float64) *AggregationBuilder {
	agg := types.NewHistogramAggregation()
	agg.Field = &field
	iv := types.Float64(interval)
	agg.Interval = &iv
	b.aggs[name] = types.Aggregations{Histogram: agg}
	return b
}

// Avg adds an avg metric aggregation.
func (b *AggregationBuilder) Avg(name, field string) *AggregationBuilder {
	agg := types.NewAverageAggregation()
	agg.Field = &field
	b.aggs[name] = types.Aggregations{Avg: agg}
	return b
}

// Max adds a max metric aggregation.
func (b *AggregationBuilder) Max(name, field string) *AggregationBuilder {
	agg := types.NewMaxAggregation()
	agg.Field = &field
	b.aggs[name] = types.Aggregations{Max: agg}
	return b
}

// Min adds a min metric aggregation.
func (b *AggregationBuilder) Min(name, field string) *AggregationBuilder {
	agg := types.NewMinAggregation()
	agg.Field = &field
	b.aggs[name] = types.Aggregations{Min: agg}
	return b
}

// Sum adds a sum metric aggregation.
func (b *AggregationBuilder) Sum(name, field string) *AggregationBuilder {
	agg := types.NewSumAggregation()
	agg.Field = &field
	b.aggs[name] = types.Aggregations{Sum: agg}
	return b
}

// ValueCount adds a value_count metric aggregation.
func (b *AggregationBuilder) ValueCount(name, field string) *AggregationBuilder {
	agg := types.NewValueCountAggregation()
	agg.Field = &field
	b.aggs[name] = types.Aggregations{ValueCount: agg}
	return b
}

// Cardinality adds a cardinality metric aggregation (approximate distinct count).
func (b *AggregationBuilder) Cardinality(name, field string) *AggregationBuilder {
	agg := types.NewCardinalityAggregation()
	agg.Field = &field
	b.aggs[name] = types.Aggregations{Cardinality: agg}
	return b
}

// Stats adds a stats metric aggregation (count, min, max, avg, sum in one).
func (b *AggregationBuilder) Stats(name, field string) *AggregationBuilder {
	agg := types.NewStatsAggregation()
	agg.Field = &field
	b.aggs[name] = types.Aggregations{Stats: agg}
	return b
}

// Nested adds a nested bucket aggregation with sub-aggregations.
func (b *AggregationBuilder) Nested(name, path string, sub *AggregationBuilder) *AggregationBuilder {
	agg := types.NewNestedAggregation()
	agg.Path = &path
	b.aggs[name] = types.Aggregations{
		Nested:       agg,
		Aggregations: sub.Build(),
	}
	return b
}

// Filter adds a single-bucket filter aggregation with sub-aggregations.
func (b *AggregationBuilder) Filter(name string, filter types.Query, sub *AggregationBuilder) *AggregationBuilder {
	entry := types.Aggregations{
		Filter:       &filter,
		Aggregations: sub.Build(),
	}
	b.aggs[name] = entry
	return b
}

// SubAggregations adds sub-aggregations to an existing named aggregation.
// Call this after adding the parent aggregation (Terms, DateHistogram, etc.).
func (b *AggregationBuilder) SubAggregations(name string, sub *AggregationBuilder) *AggregationBuilder {
	if entry, ok := b.aggs[name]; ok {
		entry.Aggregations = sub.Build()
		b.aggs[name] = entry
	}
	return b
}

// Build returns the aggregations map ready for use in a search request.
func (b *AggregationBuilder) Build() map[string]types.Aggregations {
	return b.aggs
}
