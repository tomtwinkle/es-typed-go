// Package query provides aggregation building helpers for Elasticsearch v8.
//
// Deprecated: import github.com/tomtwinkle/es-typed-go/query instead.
package query

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/tomtwinkle/es-typed-go/estype"
	sharedquery "github.com/tomtwinkle/es-typed-go/query"
)

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.AggEntry.
type AggEntry = sharedquery.AggEntry

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.AggDefinition.
type AggDefinition[R any] = sharedquery.AggDefinition[R]

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.AggSet.
type AggSet = sharedquery.AggSet

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.AggResults.
type AggResults = sharedquery.AggResults

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.AvgAggregation.
type AvgAggregation = sharedquery.AvgAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.AvgResult.
type AvgResult = sharedquery.AvgResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.SumAggregation.
type SumAggregation = sharedquery.SumAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.SumResult.
type SumResult = sharedquery.SumResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MinAggregation.
type MinAggregation = sharedquery.MinAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MinResult.
type MinResult = sharedquery.MinResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MaxAggregation.
type MaxAggregation = sharedquery.MaxAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MaxResult.
type MaxResult = sharedquery.MaxResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.StatsAggregation.
type StatsAggregation = sharedquery.StatsAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.StatsResult.
type StatsResult = sharedquery.StatsResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.ValueCountAggregation.
type ValueCountAggregation = sharedquery.ValueCountAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.ValueCountResult.
type ValueCountResult = sharedquery.ValueCountResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.CardinalityAggregation.
type CardinalityAggregation = sharedquery.CardinalityAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.CardinalityResult.
type CardinalityResult = sharedquery.CardinalityResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.TermsAggregation.
type TermsAggregation = sharedquery.TermsAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.TermsAggOption.
type TermsAggOption = sharedquery.TermsAggOption

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.StringTermsResult.
type StringTermsResult = sharedquery.StringTermsResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.StringTermsBucket.
type StringTermsBucketResult = sharedquery.StringTermsBucket

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.DateHistogramAggregation.
type DateHistogramAggregation = sharedquery.DateHistogramAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.DateHistogramAggOption.
type DateHistogramAggOption = sharedquery.DateHistogramAggOption

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.DateHistogramResult.
type DateHistogramResult = sharedquery.DateHistogramResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.DateHistogramBucket.
type DateHistogramBucketResult = sharedquery.DateHistogramBucket

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.HistogramAggregation.
type HistogramAggregation = sharedquery.HistogramAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.HistogramAggOption.
type HistogramAggOption = sharedquery.HistogramAggOption

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.HistogramResult.
type HistogramResult = sharedquery.HistogramResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.HistogramBucket.
type HistogramBucketResult = sharedquery.HistogramBucket

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.NewAggResults.
func NewAggResults(raw map[string]types.Aggregate) sharedquery.AggResults {
	return sharedquery.NewAggResults(raw)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.Aggs.
func Aggs(defs ...sharedquery.AggEntry) sharedquery.AggSet {
	return sharedquery.Aggs(defs...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.GetAgg.
func GetAgg[R any](r sharedquery.AggResults, def sharedquery.AggDefinition[R]) (R, error) {
	return sharedquery.GetAgg[R](r, def)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MustAgg.
func MustAgg[R any](r sharedquery.AggResults, def sharedquery.AggDefinition[R]) R {
	return sharedquery.MustAgg[R](r, def)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.AvgAgg.
func AvgAgg(name string, field estype.Field) sharedquery.AvgAggregation {
	return sharedquery.AvgAgg(name, field)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.SumAgg.
func SumAgg(name string, field estype.Field) sharedquery.SumAggregation {
	return sharedquery.SumAgg(name, field)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MinAgg.
func MinAgg(name string, field estype.Field) sharedquery.MinAggregation {
	return sharedquery.MinAgg(name, field)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MaxAgg.
func MaxAgg(name string, field estype.Field) sharedquery.MaxAggregation {
	return sharedquery.MaxAgg(name, field)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.StatsAgg.
func StatsAgg(name string, field estype.Field) sharedquery.StatsAggregation {
	return sharedquery.StatsAgg(name, field)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.ValueCountAgg.
func ValueCountAgg(name string, field estype.Field) sharedquery.ValueCountAggregation {
	return sharedquery.ValueCountAgg(name, field)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.CardinalityAgg.
func CardinalityAgg(name string, field estype.Field) sharedquery.CardinalityAggregation {
	return sharedquery.CardinalityAgg(name, field)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.StringTermsAgg.
func StringTermsAgg(name string, field estype.Field, opts ...sharedquery.TermsAggOption) sharedquery.TermsAggregation {
	return sharedquery.StringTermsAgg(name, field, opts...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.TermsAgg.
func TermsAgg(name string, field estype.Field, opts ...sharedquery.TermsAggOption) sharedquery.TermsAggregation {
	return sharedquery.TermsAgg(name, field, opts...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.WithTermsSize.
func WithTermsSize(size int) sharedquery.TermsAggOption {
	return sharedquery.WithTermsSize(size)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.TermsAggSize.
func TermsAggSize(size int) sharedquery.TermsAggOption {
	return sharedquery.TermsAggSize(size)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.WithSubAggs.
func WithSubAggs(defs ...sharedquery.AggEntry) sharedquery.TermsAggOption {
	return sharedquery.WithSubAggs(defs...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.TermsAggSubAggs.
func TermsAggSubAggs(defs ...sharedquery.AggEntry) sharedquery.TermsAggOption {
	return sharedquery.TermsAggSubAggs(defs...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.DateHistogramAgg.
func DateHistogramAgg(
	name string,
	field estype.Field,
	interval calendarinterval.CalendarInterval,
	opts ...sharedquery.DateHistogramAggOption,
) sharedquery.DateHistogramAggregation {
	return sharedquery.DateHistogramAgg(name, field, interval, opts...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.DateHistogramAggFormat.
func DateHistogramAggFormat(format string) sharedquery.DateHistogramAggOption {
	return sharedquery.DateHistogramAggFormat(format)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.DateHistogramAggSubAggs.
func DateHistogramAggSubAggs(defs ...sharedquery.AggEntry) sharedquery.DateHistogramAggOption {
	return sharedquery.DateHistogramAggSubAggs(defs...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.HistogramAgg.
func HistogramAgg(name string, field estype.Field, interval float64, opts ...sharedquery.HistogramAggOption) sharedquery.HistogramAggregation {
	return sharedquery.HistogramAgg(name, field, interval, opts...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.HistogramAggSubAggs.
func HistogramAggSubAggs(defs ...sharedquery.AggEntry) sharedquery.HistogramAggOption {
	return sharedquery.HistogramAggSubAggs(defs...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.NestedResult.
type NestedResult = sharedquery.NestedResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.FilterResult.
type FilterResult = sharedquery.FilterResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MultiTermsBucketResult.
type MultiTermsBucketResult = sharedquery.MultiTermsBucketResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MultiTermsResult.
type MultiTermsResult = sharedquery.MultiTermsResult

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.NestedAggregation.
type NestedAggregation = sharedquery.NestedAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.NestedAggOption.
type NestedAggOption = sharedquery.NestedAggOption

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.FilterAggregation.
type FilterAggregation = sharedquery.FilterAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.FilterAggOption.
type FilterAggOption = sharedquery.FilterAggOption

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MultiTermsAggregation.
type MultiTermsAggregation = sharedquery.MultiTermsAggregation

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MultiTermsAggOption.
type MultiTermsAggOption = sharedquery.MultiTermsAggOption

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.NestedAgg.
func NestedAgg(name string, path estype.Field, opts ...sharedquery.NestedAggOption) sharedquery.NestedAggregation {
	return sharedquery.NestedAgg(name, path, opts...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.NestedAggSubAggs.
func NestedAggSubAggs(defs ...sharedquery.AggEntry) sharedquery.NestedAggOption {
	return sharedquery.NestedAggSubAggs(defs...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.FilterAgg.
func FilterAgg(name string, filter types.Query, opts ...sharedquery.FilterAggOption) sharedquery.FilterAggregation {
	return sharedquery.FilterAgg(name, filter, opts...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.FilterAggSubAggs.
func FilterAggSubAggs(defs ...sharedquery.AggEntry) sharedquery.FilterAggOption {
	return sharedquery.FilterAggSubAggs(defs...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MultiTermsAgg.
func MultiTermsAgg(name string, fields []estype.Field, opts ...sharedquery.MultiTermsAggOption) sharedquery.MultiTermsAggregation {
	return sharedquery.MultiTermsAgg(name, fields, opts...)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MultiTermsAggSize.
func MultiTermsAggSize(size int) sharedquery.MultiTermsAggOption {
	return sharedquery.MultiTermsAggSize(size)
}

// Deprecated: use github.com/tomtwinkle/es-typed-go/query.MultiTermsAggSubAggs.
func MultiTermsAggSubAggs(defs ...sharedquery.AggEntry) sharedquery.MultiTermsAggOption {
	return sharedquery.MultiTermsAggSubAggs(defs...)
}
