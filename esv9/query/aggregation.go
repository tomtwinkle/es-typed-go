package query

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/calendarinterval"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// GetNested returns a typed nested aggregation result.
func (r AggResults) GetNested(def NestedAggregation) (NestedResult, error) {
	return def.parse(r.raw)
}

// MustNested returns a typed nested aggregation result or panics.
func (r AggResults) MustNested(def NestedAggregation) NestedResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetFilter returns a typed filter aggregation result.
func (r AggResults) GetFilter(def FilterAggregation) (FilterResult, error) {
	return def.parse(r.raw)
}

// MustFilter returns a typed filter aggregation result or panics.
func (r AggResults) MustFilter(def FilterAggregation) FilterResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetMultiTerms returns a typed multi_terms aggregation result.
func (r AggResults) GetMultiTerms(def MultiTermsAggregation) (MultiTermsResult, error) {
	return def.parse(r.raw)
}

// MustMultiTerms returns a typed multi_terms aggregation result or panics.
func (r AggResults) MustMultiTerms(def MultiTermsAggregation) MultiTermsResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// AggDefinition is a typed aggregation definition that can build an Elasticsearch
// aggregation request and decode the corresponding response value.
type AggDefinition[R any] interface {
	Name() string
	build() types.Aggregations
	parse(map[string]types.Aggregate) (R, error)
}

// AggSet is a collection of typed aggregation definitions.
type AggSet struct {
	defs []aggDefinitionAny
}

type aggDefinitionAny interface {
	Name() string
	build() types.Aggregations
}

// Aggs creates an aggregation set from typed aggregation definitions.
func Aggs(defs ...aggDefinitionAny) AggSet {
	copied := make([]aggDefinitionAny, len(defs))
	copy(copied, defs)
	return AggSet{defs: copied}
}

// Build converts the aggregation set into the Elasticsearch request shape.
func (s AggSet) Build() map[string]types.Aggregations {
	if len(s.defs) == 0 {
		return nil
	}
	out := make(map[string]types.Aggregations, len(s.defs))
	for _, def := range s.defs {
		out[def.Name()] = def.build()
	}
	return out
}

// AggResults wraps raw aggregation results and provides typed accessors.
type AggResults struct {
	raw map[string]types.Aggregate
}

// NewAggResults wraps raw aggregation results.
func NewAggResults(raw map[string]types.Aggregate) AggResults {
	return AggResults{raw: raw}
}

// Raw returns the raw Elasticsearch aggregation map.
func (r AggResults) Raw() map[string]types.Aggregate {
	return r.raw
}

// Get returns a typed aggregation result.
func GetAgg[R any](r AggResults, def AggDefinition[R]) (R, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics if the aggregation is
// missing or has an unexpected type.
func MustAgg[R any](r AggResults, def AggDefinition[R]) R {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetAvg(def AvgAggregation) (AvgResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustAvg(def AvgAggregation) AvgResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetSum(def SumAggregation) (SumResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustSum(def SumAggregation) SumResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetMin(def MinAggregation) (MinResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustMin(def MinAggregation) MinResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetMax(def MaxAggregation) (MaxResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustMax(def MaxAggregation) MaxResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetStats(def StatsAggregation) (StatsResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustStats(def StatsAggregation) StatsResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetValueCount(def ValueCountAggregation) (ValueCountResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustValueCount(def ValueCountAggregation) ValueCountResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetCardinality(def CardinalityAggregation) (CardinalityResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustCardinality(def CardinalityAggregation) CardinalityResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetStringTerms(def TermsAggregation) (StringTermsResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustStringTerms(def TermsAggregation) StringTermsResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetDateHistogram(def DateHistogramAggregation) (DateHistogramResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustDateHistogram(def DateHistogramAggregation) DateHistogramResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// Get returns a typed aggregation result.
func (r AggResults) GetHistogram(def HistogramAggregation) (HistogramResult, error) {
	return def.parse(r.raw)
}

// Must returns a typed aggregation result or panics.
func (r AggResults) MustHistogram(def HistogramAggregation) HistogramResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// AvgResult wraps an avg aggregation result.
type AvgResult struct {
	value *float64
}

// Value returns the avg value.
func (r AvgResult) Value() *float64 { return r.value }

// SumResult wraps a sum aggregation result.
type SumResult struct {
	value *float64
}

// Value returns the sum value.
func (r SumResult) Value() *float64 { return r.value }

// MinResult wraps a min aggregation result.
type MinResult struct {
	value *float64
}

// Value returns the min value.
func (r MinResult) Value() *float64 { return r.value }

// MaxResult wraps a max aggregation result.
type MaxResult struct {
	value *float64
}

// Value returns the max value.
func (r MaxResult) Value() *float64 { return r.value }

// ValueCountResult wraps a value_count aggregation result.
type ValueCountResult struct {
	value *int64
}

// Value returns the value_count value.
func (r ValueCountResult) Value() *int64 { return r.value }

// CardinalityResult wraps a cardinality aggregation result.
type CardinalityResult struct {
	value *int64
}

// Value returns the cardinality value.
func (r CardinalityResult) Value() *int64 { return r.value }

// StatsResult wraps a stats aggregation result.
type StatsResult struct {
	count int64
	min   *float64
	max   *float64
	avg   *float64
	sum   *float64
}

// Count returns the count value.
func (r StatsResult) Count() int64 { return r.count }

// Min returns the min value.
func (r StatsResult) Min() *float64 { return r.min }

// Max returns the max value.
func (r StatsResult) Max() *float64 { return r.max }

// Avg returns the avg value.
func (r StatsResult) Avg() *float64 { return r.avg }

// Sum returns the sum value.
func (r StatsResult) Sum() *float64 { return r.sum }

// StringTermsResult wraps a string terms aggregation result.
type StringTermsResult struct {
	buckets []StringTermsBucket
}

// Buckets returns all buckets.
func (r StringTermsResult) Buckets() []StringTermsBucket { return r.buckets }

// StringTermsBucket wraps a string terms bucket.
type StringTermsBucket struct {
	key          string
	docCount     int64
	aggregations AggResults
}

// Key returns the bucket key.
func (b StringTermsBucket) Key() string { return b.key }

// DocCount returns the document count.
func (b StringTermsBucket) DocCount() int64 { return b.docCount }

// Aggregations returns sub-aggregation results for the bucket.
func (b StringTermsBucket) Aggregations() AggResults { return b.aggregations }

// DateHistogramResult wraps a date histogram aggregation result.
type DateHistogramResult struct {
	buckets []DateHistogramBucket
}

// Buckets returns all buckets.
func (r DateHistogramResult) Buckets() []DateHistogramBucket { return r.buckets }

// DateHistogramBucket wraps a date histogram bucket.
type DateHistogramBucket struct {
	key          int64
	keyAsString  string
	docCount     int64
	aggregations AggResults
}

// Key returns the bucket key.
func (b DateHistogramBucket) Key() int64 { return b.key }

// KeyAsString returns the formatted key if present.
func (b DateHistogramBucket) KeyAsString() string { return b.keyAsString }

// DocCount returns the document count.
func (b DateHistogramBucket) DocCount() int64 { return b.docCount }

// Aggregations returns sub-aggregation results for the bucket.
func (b DateHistogramBucket) Aggregations() AggResults { return b.aggregations }

// HistogramResult wraps a histogram aggregation result.
type HistogramResult struct {
	buckets []HistogramBucket
}

// Buckets returns all buckets.
func (r HistogramResult) Buckets() []HistogramBucket { return r.buckets }

// HistogramBucket wraps a histogram bucket.
type HistogramBucket struct {
	key          float64
	docCount     int64
	aggregations AggResults
}

// Key returns the bucket key.
func (b HistogramBucket) Key() float64 { return b.key }

// DocCount returns the document count.
func (b HistogramBucket) DocCount() int64 { return b.docCount }

// Aggregations returns sub-aggregation results for the bucket.
func (b HistogramBucket) Aggregations() AggResults { return b.aggregations }

// TermsAggOption configures a string terms aggregation.
type TermsAggOption func(*TermsAggregation)

// DateHistogramAggOption configures a date histogram aggregation.
type DateHistogramAggOption func(*DateHistogramAggregation)

// HistogramAggOption configures a histogram aggregation.
type HistogramAggOption func(*HistogramAggregation)

// TermsAggregation defines a typed string terms aggregation.
type TermsAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a TermsAggregation) Name() string { return a.name }

func (a TermsAggregation) build() types.Aggregations { return a.agg }

// StringTermsAgg creates a typed string terms aggregation definition.
func StringTermsAgg(name string, field estype.Field, opts ...TermsAggOption) TermsAggregation {
	agg := types.NewTermsAggregation()
	f := string(field)
	agg.Field = &f
	def := TermsAggregation{
		name: name,
		agg:  types.Aggregations{Terms: agg},
	}
	for _, opt := range opts {
		opt(&def)
	}
	return def
}

// TermsAggSize sets the terms size.
func TermsAggSize(size int) TermsAggOption {
	return func(a *TermsAggregation) {
		if a.agg.Terms == nil {
			a.agg.Terms = types.NewTermsAggregation()
		}
		a.agg.Terms.Size = &size
	}
}

// WithTermsSize sets the terms size.
func WithTermsSize(size int) TermsAggOption {
	return TermsAggSize(size)
}

// TermsAggSubAggs sets sub-aggregations.
func TermsAggSubAggs(defs ...aggDefinitionAny) TermsAggOption {
	return func(a *TermsAggregation) {
		a.agg.Aggregations = Aggs(defs...).Build()
	}
}

// WithSubAggs sets sub-aggregations.
func WithSubAggs(defs ...aggDefinitionAny) TermsAggOption {
	return TermsAggSubAggs(defs...)
}

func (a TermsAggregation) parse(raw map[string]types.Aggregate) (StringTermsResult, error) {
	agg, err := getAggregate[*types.StringTermsAggregate](raw, a.name, "string_terms")
	if err != nil {
		return StringTermsResult{}, err
	}
	buckets, ok := agg.Buckets.([]types.StringTermsBucket)
	if !ok {
		return StringTermsResult{}, fmt.Errorf("aggregation %q returned unexpected buckets type for string terms", a.name)
	}
	out := make([]StringTermsBucket, 0, len(buckets))
	for _, bucket := range buckets {
		key, ok := bucket.Key.(string)
		if !ok {
			return StringTermsResult{}, fmt.Errorf("aggregation %q returned non-string bucket key of type %T", a.name, bucket.Key)
		}
		out = append(out, StringTermsBucket{
			key:          key,
			docCount:     bucket.DocCount,
			aggregations: NewAggResults(bucket.Aggregations),
		})
	}
	return StringTermsResult{buckets: out}, nil
}

// DateHistogramAggregation defines a typed date histogram aggregation.
type DateHistogramAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a DateHistogramAggregation) Name() string { return a.name }

func (a DateHistogramAggregation) build() types.Aggregations { return a.agg }

// DateHistogramAgg creates a typed date histogram aggregation definition.
func DateHistogramAgg(name string, field estype.Field, interval calendarinterval.CalendarInterval, opts ...DateHistogramAggOption) DateHistogramAggregation {
	agg := types.NewDateHistogramAggregation()
	f := string(field)
	agg.Field = &f
	agg.CalendarInterval = &interval
	def := DateHistogramAggregation{
		name: name,
		agg:  types.Aggregations{DateHistogram: agg},
	}
	for _, opt := range opts {
		opt(&def)
	}
	return def
}

// DateHistogramAggFormat sets the date histogram format.
func DateHistogramAggFormat(format string) DateHistogramAggOption {
	return func(a *DateHistogramAggregation) {
		if a.agg.DateHistogram == nil {
			a.agg.DateHistogram = types.NewDateHistogramAggregation()
		}
		a.agg.DateHistogram.Format = &format
	}
}

// DateHistogramAggSubAggs sets sub-aggregations.
func DateHistogramAggSubAggs(defs ...aggDefinitionAny) DateHistogramAggOption {
	return func(a *DateHistogramAggregation) {
		a.agg.Aggregations = Aggs(defs...).Build()
	}
}

func (a DateHistogramAggregation) parse(raw map[string]types.Aggregate) (DateHistogramResult, error) {
	agg, err := getAggregate[*types.DateHistogramAggregate](raw, a.name, "date_histogram")
	if err != nil {
		return DateHistogramResult{}, err
	}
	buckets, ok := agg.Buckets.([]types.DateHistogramBucket)
	if !ok {
		return DateHistogramResult{}, fmt.Errorf("aggregation %q returned unexpected buckets type for date histogram", a.name)
	}
	out := make([]DateHistogramBucket, 0, len(buckets))
	for _, bucket := range buckets {
		keyAsString := ""
		if bucket.KeyAsString != nil {
			keyAsString = *bucket.KeyAsString
		}
		out = append(out, DateHistogramBucket{
			key:          bucket.Key,
			keyAsString:  keyAsString,
			docCount:     bucket.DocCount,
			aggregations: NewAggResults(bucket.Aggregations),
		})
	}
	return DateHistogramResult{buckets: out}, nil
}

// HistogramAggregation defines a typed histogram aggregation.
type HistogramAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a HistogramAggregation) Name() string { return a.name }

func (a HistogramAggregation) build() types.Aggregations { return a.agg }

// HistogramAgg creates a typed histogram aggregation definition.
func HistogramAgg(name string, field estype.Field, interval float64, opts ...HistogramAggOption) HistogramAggregation {
	agg := types.NewHistogramAggregation()
	f := string(field)
	agg.Field = &f
	iv := types.Float64(interval)
	agg.Interval = &iv
	def := HistogramAggregation{
		name: name,
		agg:  types.Aggregations{Histogram: agg},
	}
	for _, opt := range opts {
		opt(&def)
	}
	return def
}

// HistogramAggSubAggs sets sub-aggregations.
func HistogramAggSubAggs(defs ...aggDefinitionAny) HistogramAggOption {
	return func(a *HistogramAggregation) {
		a.agg.Aggregations = Aggs(defs...).Build()
	}
}

func (a HistogramAggregation) parse(raw map[string]types.Aggregate) (HistogramResult, error) {
	agg, err := getAggregate[*types.HistogramAggregate](raw, a.name, "histogram")
	if err != nil {
		return HistogramResult{}, err
	}
	buckets, ok := agg.Buckets.([]types.HistogramBucket)
	if !ok {
		return HistogramResult{}, fmt.Errorf("aggregation %q returned unexpected buckets type for histogram", a.name)
	}
	out := make([]HistogramBucket, 0, len(buckets))
	for _, bucket := range buckets {
		out = append(out, HistogramBucket{
			key:          float64(bucket.Key),
			docCount:     bucket.DocCount,
			aggregations: NewAggResults(bucket.Aggregations),
		})
	}
	return HistogramResult{buckets: out}, nil
}

// AvgAggregation defines a typed avg aggregation.
type AvgAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a AvgAggregation) Name() string { return a.name }

func (a AvgAggregation) build() types.Aggregations { return a.agg }

// AvgAgg creates a typed avg aggregation definition.
func AvgAgg(name string, field estype.Field) AvgAggregation {
	agg := types.NewAverageAggregation()
	f := string(field)
	agg.Field = &f
	return AvgAggregation{name: name, agg: types.Aggregations{Avg: agg}}
}

func (a AvgAggregation) parse(raw map[string]types.Aggregate) (AvgResult, error) {
	agg, err := getAggregate[*types.AvgAggregate](raw, a.name, "avg")
	if err != nil {
		return AvgResult{}, err
	}
	return AvgResult{value: float64PtrFromTypes(agg.Value)}, nil
}

// SumAggregation defines a typed sum aggregation.
type SumAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a SumAggregation) Name() string { return a.name }

func (a SumAggregation) build() types.Aggregations { return a.agg }

// SumAgg creates a typed sum aggregation definition.
func SumAgg(name string, field estype.Field) SumAggregation {
	agg := types.NewSumAggregation()
	f := string(field)
	agg.Field = &f
	return SumAggregation{name: name, agg: types.Aggregations{Sum: agg}}
}

func (a SumAggregation) parse(raw map[string]types.Aggregate) (SumResult, error) {
	agg, err := getAggregate[*types.SumAggregate](raw, a.name, "sum")
	if err != nil {
		return SumResult{}, err
	}
	return SumResult{value: float64PtrFromTypes(agg.Value)}, nil
}

// MinAggregation defines a typed min aggregation.
type MinAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a MinAggregation) Name() string { return a.name }

func (a MinAggregation) build() types.Aggregations { return a.agg }

// MinAgg creates a typed min aggregation definition.
func MinAgg(name string, field estype.Field) MinAggregation {
	agg := types.NewMinAggregation()
	f := string(field)
	agg.Field = &f
	return MinAggregation{name: name, agg: types.Aggregations{Min: agg}}
}

func (a MinAggregation) parse(raw map[string]types.Aggregate) (MinResult, error) {
	agg, err := getAggregate[*types.MinAggregate](raw, a.name, "min")
	if err != nil {
		return MinResult{}, err
	}
	return MinResult{value: float64PtrFromTypes(agg.Value)}, nil
}

// MaxAggregation defines a typed max aggregation.
type MaxAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a MaxAggregation) Name() string { return a.name }

func (a MaxAggregation) build() types.Aggregations { return a.agg }

// MaxAgg creates a typed max aggregation definition.
func MaxAgg(name string, field estype.Field) MaxAggregation {
	agg := types.NewMaxAggregation()
	f := string(field)
	agg.Field = &f
	return MaxAggregation{name: name, agg: types.Aggregations{Max: agg}}
}

func (a MaxAggregation) parse(raw map[string]types.Aggregate) (MaxResult, error) {
	agg, err := getAggregate[*types.MaxAggregate](raw, a.name, "max")
	if err != nil {
		return MaxResult{}, err
	}
	return MaxResult{value: float64PtrFromTypes(agg.Value)}, nil
}

// StatsAggregation defines a typed stats aggregation.
type StatsAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a StatsAggregation) Name() string { return a.name }

func (a StatsAggregation) build() types.Aggregations { return a.agg }

// StatsAgg creates a typed stats aggregation definition.
func StatsAgg(name string, field estype.Field) StatsAggregation {
	agg := types.NewStatsAggregation()
	f := string(field)
	agg.Field = &f
	return StatsAggregation{name: name, agg: types.Aggregations{Stats: agg}}
}

func (a StatsAggregation) parse(raw map[string]types.Aggregate) (StatsResult, error) {
	agg, err := getAggregate[*types.StatsAggregate](raw, a.name, "stats")
	if err != nil {
		return StatsResult{}, err
	}
	return StatsResult{
		count: agg.Count,
		min:   float64PtrFromTypes(agg.Min),
		max:   float64PtrFromTypes(agg.Max),
		avg:   float64PtrFromTypes(agg.Avg),
		sum:   float64PtrFromFloat64(agg.Sum),
	}, nil
}

// ValueCountAggregation defines a typed value_count aggregation.
type ValueCountAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a ValueCountAggregation) Name() string { return a.name }

func (a ValueCountAggregation) build() types.Aggregations { return a.agg }

// ValueCountAgg creates a typed value_count aggregation definition.
func ValueCountAgg(name string, field estype.Field) ValueCountAggregation {
	agg := types.NewValueCountAggregation()
	f := string(field)
	agg.Field = &f
	return ValueCountAggregation{name: name, agg: types.Aggregations{ValueCount: agg}}
}

func (a ValueCountAggregation) parse(raw map[string]types.Aggregate) (ValueCountResult, error) {
	agg, err := getAggregate[*types.ValueCountAggregate](raw, a.name, "value_count")
	if err != nil {
		return ValueCountResult{}, err
	}
	return ValueCountResult{value: int64PtrFromTypes(agg.Value)}, nil
}

// CardinalityAggregation defines a typed cardinality aggregation.
type CardinalityAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a CardinalityAggregation) Name() string { return a.name }

func (a CardinalityAggregation) build() types.Aggregations { return a.agg }

// CardinalityAgg creates a typed cardinality aggregation definition.
func CardinalityAgg(name string, field estype.Field) CardinalityAggregation {
	agg := types.NewCardinalityAggregation()
	f := string(field)
	agg.Field = &f
	return CardinalityAggregation{name: name, agg: types.Aggregations{Cardinality: agg}}
}

func (a CardinalityAggregation) parse(raw map[string]types.Aggregate) (CardinalityResult, error) {
	agg, err := getAggregate[*types.CardinalityAggregate](raw, a.name, "cardinality")
	if err != nil {
		return CardinalityResult{}, err
	}
	return CardinalityResult{value: int64PtrFromInt64(agg.Value)}, nil
}

// NestedResult wraps a nested aggregation result.
type NestedResult struct {
	docCount     int64
	aggregations AggResults
}

// DocCount returns the document count.
func (r NestedResult) DocCount() int64 { return r.docCount }

// Aggregations returns sub-aggregation results.
func (r NestedResult) Aggregations() AggResults { return r.aggregations }

// FilterResult wraps a filter aggregation result.
type FilterResult struct {
	docCount     int64
	aggregations AggResults
}

// DocCount returns the document count.
func (r FilterResult) DocCount() int64 { return r.docCount }

// Aggregations returns sub-aggregation results.
func (r FilterResult) Aggregations() AggResults { return r.aggregations }

// MultiTermsBucket wraps a multi_terms bucket.
type MultiTermsBucket struct {
	keys         []string
	docCount     int64
	aggregations AggResults
}

// Keys returns the composite bucket keys as strings.
func (b MultiTermsBucket) Keys() []string { return b.keys }

// DocCount returns the bucket document count.
func (b MultiTermsBucket) DocCount() int64 { return b.docCount }

// Aggregations returns nested aggregation results for this bucket.
func (b MultiTermsBucket) Aggregations() AggResults { return b.aggregations }

// MultiTermsResult wraps a multi_terms aggregation result.
type MultiTermsResult struct {
	buckets []MultiTermsBucket
}

// Buckets returns the typed multi_terms buckets.
func (r MultiTermsResult) Buckets() []MultiTermsBucket { return r.buckets }

// NestedAggOption configures a nested aggregation.
type NestedAggOption func(*NestedAggregation)

// FilterAggOption configures a filter aggregation.
type FilterAggOption func(*FilterAggregation)

// MultiTermsAggOption configures a multi_terms aggregation.
type MultiTermsAggOption func(*MultiTermsAggregation)

// MultiTermLookup represents a single field entry in a multi_terms aggregation.
// Set Missing to a non-nil value to specify a replacement for documents that
// do not have a value for the field.
type MultiTermLookup struct {
	Field   estype.Field
	Missing any // optional; nil means no missing value
}

// NestedAggregation defines a typed nested aggregation.
type NestedAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a NestedAggregation) Name() string { return a.name }

func (a NestedAggregation) build() types.Aggregations { return a.agg }

// NestedAggSubAggs sets sub-aggregations for a nested aggregation.
func NestedAggSubAggs(defs ...aggDefinitionAny) NestedAggOption {
	return func(a *NestedAggregation) {
		a.agg.Aggregations = Aggs(defs...).Build()
	}
}

// NestedAgg creates a typed nested aggregation definition.
func NestedAgg(name string, path estype.Field, opts ...NestedAggOption) NestedAggregation {
	nestedAgg := types.NewNestedAggregation()
	p := string(path)
	nestedAgg.Path = &p
	def := NestedAggregation{
		name: name,
		agg:  types.Aggregations{Nested: nestedAgg},
	}
	for _, opt := range opts {
		opt(&def)
	}
	return def
}

func (a NestedAggregation) parse(raw map[string]types.Aggregate) (NestedResult, error) {
	agg, err := getAggregate[*types.NestedAggregate](raw, a.name, "nested")
	if err != nil {
		return NestedResult{}, err
	}
	return NestedResult{
		docCount:     agg.DocCount,
		aggregations: NewAggResults(agg.Aggregations),
	}, nil
}

// FilterAggregation defines a typed filter aggregation.
type FilterAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a FilterAggregation) Name() string { return a.name }

func (a FilterAggregation) build() types.Aggregations { return a.agg }

// FilterAggSubAggs sets sub-aggregations for a filter aggregation.
func FilterAggSubAggs(defs ...aggDefinitionAny) FilterAggOption {
	return func(a *FilterAggregation) {
		a.agg.Aggregations = Aggs(defs...).Build()
	}
}

// FilterAgg creates a typed filter aggregation definition.
func FilterAgg(name string, filter types.Query, opts ...FilterAggOption) FilterAggregation {
	q := filter
	def := FilterAggregation{
		name: name,
		agg:  types.Aggregations{Filter: &q},
	}
	for _, opt := range opts {
		opt(&def)
	}
	return def
}

func (a FilterAggregation) parse(raw map[string]types.Aggregate) (FilterResult, error) {
	agg, err := getAggregate[*types.FilterAggregate](raw, a.name, "filter")
	if err != nil {
		return FilterResult{}, err
	}
	return FilterResult{
		docCount:     agg.DocCount,
		aggregations: NewAggResults(agg.Aggregations),
	}, nil
}

// MultiTermsAggregation defines a typed multi_terms aggregation.
type MultiTermsAggregation struct {
	name string
	agg  types.Aggregations
}

// Name returns the aggregation name.
func (a MultiTermsAggregation) Name() string { return a.name }

func (a MultiTermsAggregation) build() types.Aggregations { return a.agg }

// MultiTermsAggSize sets the number of buckets to return.
func MultiTermsAggSize(size int) MultiTermsAggOption {
	return func(a *MultiTermsAggregation) {
		if a.agg.MultiTerms == nil {
			a.agg.MultiTerms = types.NewMultiTermsAggregation()
		}
		a.agg.MultiTerms.Size = &size
	}
}

// MultiTermsAggSubAggs sets sub-aggregations for a multi_terms aggregation.
func MultiTermsAggSubAggs(defs ...aggDefinitionAny) MultiTermsAggOption {
	return func(a *MultiTermsAggregation) {
		a.agg.Aggregations = Aggs(defs...).Build()
	}
}

// MultiTermsAgg creates a typed multi_terms aggregation definition.
func MultiTermsAgg(name string, fields []MultiTermLookup, opts ...MultiTermsAggOption) MultiTermsAggregation {
	multiAgg := types.NewMultiTermsAggregation()
	terms := make([]types.MultiTermLookup, 0, len(fields))
	for _, f := range fields {
		lookup := types.MultiTermLookup{Field: string(f.Field)}
		if f.Missing != nil {
			lookup.Missing = f.Missing
		}
		terms = append(terms, lookup)
	}
	multiAgg.Terms = terms
	def := MultiTermsAggregation{
		name: name,
		agg:  types.Aggregations{MultiTerms: multiAgg},
	}
	for _, opt := range opts {
		opt(&def)
	}
	return def
}

func (a MultiTermsAggregation) parse(raw map[string]types.Aggregate) (MultiTermsResult, error) {
	agg, err := getAggregate[*types.MultiTermsAggregate](raw, a.name, "multi_terms")
	if err != nil {
		return MultiTermsResult{}, err
	}

	buckets, ok := agg.Buckets.([]types.MultiTermsBucket)
	if !ok {
		return MultiTermsResult{}, fmt.Errorf("aggregation %q has unexpected buckets type %T", a.name, agg.Buckets)
	}

	out := make([]MultiTermsBucket, 0, len(buckets))
	for _, bucket := range buckets {
		keys := make([]string, 0, len(bucket.Key))
		for _, k := range bucket.Key {
			keys = append(keys, fmt.Sprintf("%v", k))
		}
		out = append(out, MultiTermsBucket{
			keys:         keys,
			docCount:     bucket.DocCount,
			aggregations: NewAggResults(bucket.Aggregations),
		})
	}
	return MultiTermsResult{buckets: out}, nil
}

func getAggregate[T any](raw map[string]types.Aggregate, name string, kind string) (T, error) {
	var zero T
	if len(raw) == 0 {
		return zero, fmt.Errorf("aggregation %q (%s) not found", name, kind)
	}
	v, ok := raw[name]
	if !ok {
		return zero, fmt.Errorf("aggregation %q (%s) not found", name, kind)
	}
	typed, ok := v.(T)
	if !ok {
		return zero, fmt.Errorf("aggregation %q is not of expected type %s", name, kind)
	}
	return typed, nil
}

func float64PtrFromTypes(v *types.Float64) *float64 {
	if v == nil {
		return nil
	}
	out := float64(*v)
	return &out
}

func float64PtrFromFloat64(v types.Float64) *float64 {
	out := float64(v)
	return &out
}

func int64PtrFromTypes(v *types.Float64) *int64 {
	if v == nil {
		return nil
	}
	out := int64(*v)
	return &out
}

func int64PtrFromInt64(v int64) *int64 {
	out := v
	return &out
}
