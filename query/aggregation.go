package query

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// AggDefinition is a typed aggregation definition that can build an Elasticsearch
// aggregation request and decode the corresponding response value.
type AggDefinition[R any] interface {
	Name() string
	build() types.Aggregations
	parse(map[string]types.Aggregate) (R, error)
}

// AggSet is a collection of typed aggregation definitions.
type AggSet struct {
	defs []AggEntry
}

// AggEntry is the base interface for all aggregation definitions. External packages cannot implement this interface.
type AggEntry interface {
	Name() string
	build() types.Aggregations
}

// Aggs creates an aggregation set from typed aggregation definitions.
func Aggs(defs ...AggEntry) AggSet {
	copied := make([]AggEntry, len(defs))
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

// NewAggResultsFromJSON parses JSON-encoded aggregation results.
// Use this when the aggregation map comes from a different types package (e.g., v9).
func NewAggResultsFromJSON(b []byte) AggResults {
	if len(b) == 0 {
		return AggResults{}
	}
	var aggs map[string]types.Aggregate
	_ = json.Unmarshal(b, &aggs)
	return AggResults{raw: aggs}
}

// Raw returns the raw Elasticsearch aggregation map.
func (r AggResults) Raw() map[string]types.Aggregate {
	return r.raw
}

// GetAgg returns a typed aggregation result.
func GetAgg[R any](r AggResults, def AggDefinition[R]) (R, error) {
	return def.parse(r.raw)
}

// MustAgg returns a typed aggregation result or panics if the aggregation is
// missing or has an unexpected type.
func MustAgg[R any](r AggResults, def AggDefinition[R]) R {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetAvg returns a typed avg aggregation result.
func (r AggResults) GetAvg(def AvgAggregation) (AvgResult, error) {
	return def.parse(r.raw)
}

// MustAvg returns a typed avg aggregation result or panics.
func (r AggResults) MustAvg(def AvgAggregation) AvgResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetSum returns a typed sum aggregation result.
func (r AggResults) GetSum(def SumAggregation) (SumResult, error) {
	return def.parse(r.raw)
}

// MustSum returns a typed sum aggregation result or panics.
func (r AggResults) MustSum(def SumAggregation) SumResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetMin returns a typed min aggregation result.
func (r AggResults) GetMin(def MinAggregation) (MinResult, error) {
	return def.parse(r.raw)
}

// MustMin returns a typed min aggregation result or panics.
func (r AggResults) MustMin(def MinAggregation) MinResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetMax returns a typed max aggregation result.
func (r AggResults) GetMax(def MaxAggregation) (MaxResult, error) {
	return def.parse(r.raw)
}

// MustMax returns a typed max aggregation result or panics.
func (r AggResults) MustMax(def MaxAggregation) MaxResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetStats returns a typed stats aggregation result.
func (r AggResults) GetStats(def StatsAggregation) (StatsResult, error) {
	return def.parse(r.raw)
}

// MustStats returns a typed stats aggregation result or panics.
func (r AggResults) MustStats(def StatsAggregation) StatsResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetValueCount returns a typed value_count aggregation result.
func (r AggResults) GetValueCount(def ValueCountAggregation) (ValueCountResult, error) {
	return def.parse(r.raw)
}

// MustValueCount returns a typed value_count aggregation result or panics.
func (r AggResults) MustValueCount(def ValueCountAggregation) ValueCountResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetCardinality returns a typed cardinality aggregation result.
func (r AggResults) GetCardinality(def CardinalityAggregation) (CardinalityResult, error) {
	return def.parse(r.raw)
}

// MustCardinality returns a typed cardinality aggregation result or panics.
func (r AggResults) MustCardinality(def CardinalityAggregation) CardinalityResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetStringTerms returns a typed string terms aggregation result.
func (r AggResults) GetStringTerms(def TermsAggregation) (StringTermsResult, error) {
	return def.parse(r.raw)
}

// MustStringTerms returns a typed string terms aggregation result or panics.
func (r AggResults) MustStringTerms(def TermsAggregation) StringTermsResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetDateHistogram returns a typed date histogram aggregation result.
func (r AggResults) GetDateHistogram(def DateHistogramAggregation) (DateHistogramResult, error) {
	return def.parse(r.raw)
}

// MustDateHistogram returns a typed date histogram aggregation result or panics.
func (r AggResults) MustDateHistogram(def DateHistogramAggregation) DateHistogramResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

// GetHistogram returns a typed histogram aggregation result.
func (r AggResults) GetHistogram(def HistogramAggregation) (HistogramResult, error) {
	return def.parse(r.raw)
}

// MustHistogram returns a typed histogram aggregation result or panics.
func (r AggResults) MustHistogram(def HistogramAggregation) HistogramResult {
	v, err := def.parse(r.raw)
	if err != nil {
		panic(err)
	}
	return v
}

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

// Count returns the number of values.
func (r StatsResult) Count() int64 { return r.count }

// Min returns the minimum value.
func (r StatsResult) Min() *float64 { return r.min }

// Max returns the maximum value.
func (r StatsResult) Max() *float64 { return r.max }

// Avg returns the average value.
func (r StatsResult) Avg() *float64 { return r.avg }

// Sum returns the sum value.
func (r StatsResult) Sum() *float64 { return r.sum }

// StringTermsBucket is a typed string terms bucket.
type StringTermsBucket struct {
	key          string
	docCount     int64
	aggregations AggResults
}

// Key returns the string bucket key.
func (b StringTermsBucket) Key() string { return b.key }

// DocCount returns the bucket document count.
func (b StringTermsBucket) DocCount() int64 { return b.docCount }

// Aggregations returns nested aggregation results for this bucket.
func (b StringTermsBucket) Aggregations() AggResults { return b.aggregations }

// StringTermsResult wraps a string terms aggregation result.
type StringTermsResult struct {
	buckets []StringTermsBucket
}

// Buckets returns the typed term buckets.
func (r StringTermsResult) Buckets() []StringTermsBucket { return r.buckets }

// DateHistogramBucket is a typed date histogram bucket.
type DateHistogramBucket struct {
	key          int64
	keyAsString  string
	docCount     int64
	aggregations AggResults
}

// Key returns the numeric bucket key.
func (b DateHistogramBucket) Key() int64 { return b.key }

// KeyAsString returns the formatted bucket key.
func (b DateHistogramBucket) KeyAsString() string { return b.keyAsString }

// DocCount returns the bucket document count.
func (b DateHistogramBucket) DocCount() int64 { return b.docCount }

// Aggregations returns nested aggregation results for this bucket.
func (b DateHistogramBucket) Aggregations() AggResults { return b.aggregations }

// DateHistogramResult wraps a date histogram aggregation result.
type DateHistogramResult struct {
	buckets []DateHistogramBucket
}

// Buckets returns the typed date histogram buckets.
func (r DateHistogramResult) Buckets() []DateHistogramBucket { return r.buckets }

// HistogramBucket is a typed histogram bucket.
type HistogramBucket struct {
	key          float64
	docCount     int64
	aggregations AggResults
}

// Key returns the numeric bucket key.
func (b HistogramBucket) Key() float64 { return b.key }

// DocCount returns the bucket document count.
func (b HistogramBucket) DocCount() int64 { return b.docCount }

// Aggregations returns nested aggregation results for this bucket.
func (b HistogramBucket) Aggregations() AggResults { return b.aggregations }

// HistogramResult wraps a histogram aggregation result.
type HistogramResult struct {
	buckets []HistogramBucket
}

// Buckets returns the typed histogram buckets.
func (r HistogramResult) Buckets() []HistogramBucket { return r.buckets }

type baseAggDefinition struct {
	name string
}

func (d baseAggDefinition) Name() string { return d.name }

// AvgAggregation is a typed avg aggregation definition.
type AvgAggregation struct {
	baseAggDefinition
	field estype.Field
}

// AvgAgg creates an avg aggregation definition.
func AvgAgg(name string, field estype.Field) AvgAggregation {
	return AvgAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
	}
}

func (a AvgAggregation) build() types.Aggregations {
	agg := types.NewAverageAggregation()
	field := string(a.field)
	agg.Field = &field
	return types.Aggregations{Avg: agg}
}

func (a AvgAggregation) parse(raw map[string]types.Aggregate) (AvgResult, error) {
	agg, err := requireAggregate[*types.AvgAggregate](raw, a.name)
	if err != nil {
		return AvgResult{}, err
	}
	return AvgResult{value: float64PtrFromFloat64(agg.Value)}, nil
}

// SumAggregation is a typed sum aggregation definition.
type SumAggregation struct {
	baseAggDefinition
	field estype.Field
}

// SumAgg creates a sum aggregation definition.
func SumAgg(name string, field estype.Field) SumAggregation {
	return SumAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
	}
}

func (a SumAggregation) build() types.Aggregations {
	agg := types.NewSumAggregation()
	field := string(a.field)
	agg.Field = &field
	return types.Aggregations{Sum: agg}
}

func (a SumAggregation) parse(raw map[string]types.Aggregate) (SumResult, error) {
	agg, err := requireAggregate[*types.SumAggregate](raw, a.name)
	if err != nil {
		return SumResult{}, err
	}
	return SumResult{value: float64PtrFromFloat64(agg.Value)}, nil
}

// MinAggregation is a typed min aggregation definition.
type MinAggregation struct {
	baseAggDefinition
	field estype.Field
}

// MinAgg creates a min aggregation definition.
func MinAgg(name string, field estype.Field) MinAggregation {
	return MinAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
	}
}

func (a MinAggregation) build() types.Aggregations {
	agg := types.NewMinAggregation()
	field := string(a.field)
	agg.Field = &field
	return types.Aggregations{Min: agg}
}

func (a MinAggregation) parse(raw map[string]types.Aggregate) (MinResult, error) {
	agg, err := requireAggregate[*types.MinAggregate](raw, a.name)
	if err != nil {
		return MinResult{}, err
	}
	return MinResult{value: float64PtrFromFloat64(agg.Value)}, nil
}

// MaxAggregation is a typed max aggregation definition.
type MaxAggregation struct {
	baseAggDefinition
	field estype.Field
}

// MaxAgg creates a max aggregation definition.
func MaxAgg(name string, field estype.Field) MaxAggregation {
	return MaxAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
	}
}

func (a MaxAggregation) build() types.Aggregations {
	agg := types.NewMaxAggregation()
	field := string(a.field)
	agg.Field = &field
	return types.Aggregations{Max: agg}
}

func (a MaxAggregation) parse(raw map[string]types.Aggregate) (MaxResult, error) {
	agg, err := requireAggregate[*types.MaxAggregate](raw, a.name)
	if err != nil {
		return MaxResult{}, err
	}
	return MaxResult{value: float64PtrFromFloat64(agg.Value)}, nil
}

// ValueCountAggregation is a typed value_count aggregation definition.
type ValueCountAggregation struct {
	baseAggDefinition
	field estype.Field
}

// ValueCountAgg creates a value_count aggregation definition.
func ValueCountAgg(name string, field estype.Field) ValueCountAggregation {
	return ValueCountAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
	}
}

func (a ValueCountAggregation) build() types.Aggregations {
	agg := types.NewValueCountAggregation()
	field := string(a.field)
	agg.Field = &field
	return types.Aggregations{ValueCount: agg}
}

func (a ValueCountAggregation) parse(raw map[string]types.Aggregate) (ValueCountResult, error) {
	agg, err := requireAggregate[*types.ValueCountAggregate](raw, a.name)
	if err != nil {
		return ValueCountResult{}, err
	}
	return ValueCountResult{value: int64PtrFromFloat64(agg.Value)}, nil
}

// CardinalityAggregation is a typed cardinality aggregation definition.
type CardinalityAggregation struct {
	baseAggDefinition
	field estype.Field
}

// CardinalityAgg creates a cardinality aggregation definition.
func CardinalityAgg(name string, field estype.Field) CardinalityAggregation {
	return CardinalityAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
	}
}

func (a CardinalityAggregation) build() types.Aggregations {
	agg := types.NewCardinalityAggregation()
	field := string(a.field)
	agg.Field = &field
	return types.Aggregations{Cardinality: agg}
}

func (a CardinalityAggregation) parse(raw map[string]types.Aggregate) (CardinalityResult, error) {
	agg, err := requireAggregate[*types.CardinalityAggregate](raw, a.name)
	if err != nil {
		return CardinalityResult{}, err
	}
	return CardinalityResult{value: int64PtrFromInt64(agg.Value)}, nil
}

// StatsAggregation is a typed stats aggregation definition.
type StatsAggregation struct {
	baseAggDefinition
	field estype.Field
}

// StatsAgg creates a stats aggregation definition.
func StatsAgg(name string, field estype.Field) StatsAggregation {
	return StatsAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
	}
}

func (a StatsAggregation) build() types.Aggregations {
	agg := types.NewStatsAggregation()
	field := string(a.field)
	agg.Field = &field
	return types.Aggregations{Stats: agg}
}

func (a StatsAggregation) parse(raw map[string]types.Aggregate) (StatsResult, error) {
	agg, err := requireAggregate[*types.StatsAggregate](raw, a.name)
	if err != nil {
		return StatsResult{}, err
	}
	return StatsResult{
		count: agg.Count,
		min:   float64PtrFromFloat64(agg.Min),
		max:   float64PtrFromFloat64(agg.Max),
		avg:   float64PtrFromFloat64(agg.Avg),
		sum:   float64PtrFromFloat64Value(agg.Sum),
	}, nil
}

// TermsAggregation is a typed string terms aggregation definition.
type TermsAggregation struct {
	baseAggDefinition
	field   estype.Field
	size    *int
	subAggs []AggEntry
}

// TermsAggOption configures a string terms aggregation.
type TermsAggOption func(*TermsAggregation)

// WithTermsSize sets the terms aggregation size.
func WithTermsSize(size int) TermsAggOption {
	return func(a *TermsAggregation) {
		a.size = &size
	}
}

// TermsAggSize sets the terms aggregation size.
func TermsAggSize(size int) TermsAggOption {
	return WithTermsSize(size)
}

// WithSubAggs sets sub-aggregations for a terms aggregation.
func WithSubAggs(defs ...AggEntry) TermsAggOption {
	return func(a *TermsAggregation) {
		a.subAggs = append(a.subAggs, defs...)
	}
}

// TermsAggSubAggs sets sub-aggregations for a terms aggregation.
func TermsAggSubAggs(defs ...AggEntry) TermsAggOption {
	return WithSubAggs(defs...)
}

// StringTermsAgg creates a string terms aggregation definition.
func StringTermsAgg(name string, field estype.Field, opts ...TermsAggOption) TermsAggregation {
	agg := TermsAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
	}
	for _, opt := range opts {
		opt(&agg)
	}
	return agg
}

// TermsAgg creates a string terms aggregation definition.
func TermsAgg(name string, field estype.Field, opts ...TermsAggOption) TermsAggregation {
	return StringTermsAgg(name, field, opts...)
}

func (a TermsAggregation) build() types.Aggregations {
	agg := types.NewTermsAggregation()
	field := string(a.field)
	agg.Field = &field
	if a.size != nil {
		agg.Size = a.size
	}
	return types.Aggregations{
		Terms:        agg,
		Aggregations: Aggs(a.subAggs...).Build(),
	}
}

func (a TermsAggregation) parse(raw map[string]types.Aggregate) (StringTermsResult, error) {
	agg, err := requireAggregate[*types.StringTermsAggregate](raw, a.name)
	if err != nil {
		return StringTermsResult{}, err
	}

	buckets, ok := agg.Buckets.([]types.StringTermsBucket)
	if !ok {
		return StringTermsResult{}, fmt.Errorf("aggregation %q has unexpected buckets type %T", a.name, agg.Buckets)
	}

	out := make([]StringTermsBucket, 0, len(buckets))
	for _, bucket := range buckets {
		key, ok := bucket.Key.(string)
		if !ok {
			return StringTermsResult{}, fmt.Errorf("aggregation %q has unexpected bucket key type %T", a.name, bucket.Key)
		}
		out = append(out, StringTermsBucket{
			key:          key,
			docCount:     bucket.DocCount,
			aggregations: NewAggResults(bucket.Aggregations),
		})
	}
	return StringTermsResult{buckets: out}, nil
}

// DateHistogramAggregation is a typed date histogram aggregation definition.
type DateHistogramAggregation struct {
	baseAggDefinition
	field    estype.Field
	format   *string
	subAggs  []AggEntry
	interval calendarinterval.CalendarInterval
}

// DateHistogramAggOption configures a date histogram aggregation.
type DateHistogramAggOption func(*DateHistogramAggregation)

// DateHistogramAggFormat sets the output date format.
func DateHistogramAggFormat(format string) DateHistogramAggOption {
	return func(a *DateHistogramAggregation) {
		a.format = &format
	}
}

// DateHistogramAggSubAggs sets sub-aggregations.
func DateHistogramAggSubAggs(defs ...AggEntry) DateHistogramAggOption {
	return func(a *DateHistogramAggregation) {
		a.subAggs = append(a.subAggs, defs...)
	}
}

// DateHistogramAgg creates a date histogram aggregation definition.
func DateHistogramAgg(
	name string,
	field estype.Field,
	interval calendarinterval.CalendarInterval,
	opts ...DateHistogramAggOption,
) DateHistogramAggregation {
	agg := DateHistogramAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
		interval:          interval,
	}
	for _, opt := range opts {
		opt(&agg)
	}
	return agg
}

func (a DateHistogramAggregation) build() types.Aggregations {
	agg := types.NewDateHistogramAggregation()
	field := string(a.field)
	agg.Field = &field
	agg.CalendarInterval = &a.interval
	if a.format != nil {
		agg.Format = a.format
	}
	return types.Aggregations{
		DateHistogram: agg,
		Aggregations:  Aggs(a.subAggs...).Build(),
	}
}

func (a DateHistogramAggregation) parse(raw map[string]types.Aggregate) (DateHistogramResult, error) {
	agg, err := requireAggregate[*types.DateHistogramAggregate](raw, a.name)
	if err != nil {
		return DateHistogramResult{}, err
	}

	buckets, ok := agg.Buckets.([]types.DateHistogramBucket)
	if !ok {
		return DateHistogramResult{}, fmt.Errorf("aggregation %q has unexpected buckets type %T", a.name, agg.Buckets)
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

// HistogramAggregation is a typed histogram aggregation definition.
type HistogramAggregation struct {
	baseAggDefinition
	field    estype.Field
	interval float64
	subAggs  []AggEntry
}

// HistogramAggOption configures a histogram aggregation.
type HistogramAggOption func(*HistogramAggregation)

// HistogramAggSubAggs sets sub-aggregations.
func HistogramAggSubAggs(defs ...AggEntry) HistogramAggOption {
	return func(a *HistogramAggregation) {
		a.subAggs = append(a.subAggs, defs...)
	}
}

// HistogramAgg creates a histogram aggregation definition.
func HistogramAgg(name string, field estype.Field, interval float64, opts ...HistogramAggOption) HistogramAggregation {
	agg := HistogramAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		field:             field,
		interval:          interval,
	}
	for _, opt := range opts {
		opt(&agg)
	}
	return agg
}

func (a HistogramAggregation) build() types.Aggregations {
	agg := types.NewHistogramAggregation()
	field := string(a.field)
	agg.Field = &field
	interval := types.Float64(a.interval)
	agg.Interval = &interval
	return types.Aggregations{
		Histogram:    agg,
		Aggregations: Aggs(a.subAggs...).Build(),
	}
}

func (a HistogramAggregation) parse(raw map[string]types.Aggregate) (HistogramResult, error) {
	agg, err := requireAggregate[*types.HistogramAggregate](raw, a.name)
	if err != nil {
		return HistogramResult{}, err
	}

	buckets, ok := agg.Buckets.([]types.HistogramBucket)
	if !ok {
		return HistogramResult{}, fmt.Errorf("aggregation %q has unexpected buckets type %T", a.name, agg.Buckets)
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

// MultiTermsBucketResult is a typed multi_terms bucket.
type MultiTermsBucketResult struct {
	keys         []string
	docCount     int64
	aggregations AggResults
}

// Keys returns the composite bucket keys as strings.
func (b MultiTermsBucketResult) Keys() []string { return b.keys }

// DocCount returns the bucket document count.
func (b MultiTermsBucketResult) DocCount() int64 { return b.docCount }

// Aggregations returns nested aggregation results for this bucket.
func (b MultiTermsBucketResult) Aggregations() AggResults { return b.aggregations }

// MultiTermsResult wraps a multi_terms aggregation result.
type MultiTermsResult struct {
	buckets []MultiTermsBucketResult
}

// Buckets returns the typed multi_terms buckets.
func (r MultiTermsResult) Buckets() []MultiTermsBucketResult { return r.buckets }

// NestedAggregation is a typed nested aggregation definition.
type NestedAggregation struct {
	baseAggDefinition
	path    estype.Field
	subAggs []AggEntry
}

// NestedAggOption configures a nested aggregation.
type NestedAggOption func(*NestedAggregation)

// NestedAggSubAggs sets sub-aggregations for a nested aggregation.
func NestedAggSubAggs(defs ...AggEntry) NestedAggOption {
	return func(a *NestedAggregation) {
		a.subAggs = append(a.subAggs, defs...)
	}
}

// NestedAgg creates a nested aggregation definition.
func NestedAgg(name string, path estype.Field, opts ...NestedAggOption) NestedAggregation {
	agg := NestedAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		path:              path,
	}
	for _, opt := range opts {
		opt(&agg)
	}
	return agg
}

func (a NestedAggregation) build() types.Aggregations {
	nestedAgg := types.NewNestedAggregation()
	path := string(a.path)
	nestedAgg.Path = &path
	return types.Aggregations{
		Nested:       nestedAgg,
		Aggregations: Aggs(a.subAggs...).Build(),
	}
}

func (a NestedAggregation) parse(raw map[string]types.Aggregate) (NestedResult, error) {
	agg, err := requireAggregate[*types.NestedAggregate](raw, a.name)
	if err != nil {
		return NestedResult{}, err
	}
	return NestedResult{
		docCount:     agg.DocCount,
		aggregations: NewAggResults(agg.Aggregations),
	}, nil
}

// FilterAggregation is a typed filter aggregation definition.
type FilterAggregation struct {
	baseAggDefinition
	filter  types.Query
	subAggs []AggEntry
}

// FilterAggOption configures a filter aggregation.
type FilterAggOption func(*FilterAggregation)

// FilterAggSubAggs sets sub-aggregations for a filter aggregation.
func FilterAggSubAggs(defs ...AggEntry) FilterAggOption {
	return func(a *FilterAggregation) {
		a.subAggs = append(a.subAggs, defs...)
	}
}

// FilterAgg creates a filter aggregation definition.
func FilterAgg(name string, filter types.Query, opts ...FilterAggOption) FilterAggregation {
	agg := FilterAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		filter:            filter,
	}
	for _, opt := range opts {
		opt(&agg)
	}
	return agg
}

func (a FilterAggregation) build() types.Aggregations {
	q := a.filter
	return types.Aggregations{
		Filter:       &q,
		Aggregations: Aggs(a.subAggs...).Build(),
	}
}

func (a FilterAggregation) parse(raw map[string]types.Aggregate) (FilterResult, error) {
	agg, err := requireAggregate[*types.FilterAggregate](raw, a.name)
	if err != nil {
		return FilterResult{}, err
	}
	return FilterResult{
		docCount:     agg.DocCount,
		aggregations: NewAggResults(agg.Aggregations),
	}, nil
}

// MultiTermsAggregation is a typed multi_terms aggregation definition.
type MultiTermsAggregation struct {
	baseAggDefinition
	fields  []estype.Field
	size    *int
	subAggs []AggEntry
}

// MultiTermsAggOption configures a multi_terms aggregation.
type MultiTermsAggOption func(*MultiTermsAggregation)

// MultiTermsAggSize sets the number of buckets to return.
func MultiTermsAggSize(size int) MultiTermsAggOption {
	return func(a *MultiTermsAggregation) {
		a.size = &size
	}
}

// MultiTermsAggSubAggs sets sub-aggregations for a multi_terms aggregation.
func MultiTermsAggSubAggs(defs ...AggEntry) MultiTermsAggOption {
	return func(a *MultiTermsAggregation) {
		a.subAggs = append(a.subAggs, defs...)
	}
}

// MultiTermsAgg creates a multi_terms aggregation definition.
func MultiTermsAgg(name string, fields []estype.Field, opts ...MultiTermsAggOption) MultiTermsAggregation {
	agg := MultiTermsAggregation{
		baseAggDefinition: baseAggDefinition{name: name},
		fields:            fields,
	}
	for _, opt := range opts {
		opt(&agg)
	}
	return agg
}

func (a MultiTermsAggregation) build() types.Aggregations {
	multiAgg := types.NewMultiTermsAggregation()
	terms := make([]types.MultiTermLookup, 0, len(a.fields))
	for _, f := range a.fields {
		terms = append(terms, types.MultiTermLookup{Field: string(f)})
	}
	multiAgg.Terms = terms
	if a.size != nil {
		multiAgg.Size = a.size
	}
	return types.Aggregations{
		MultiTerms:   multiAgg,
		Aggregations: Aggs(a.subAggs...).Build(),
	}
}

func (a MultiTermsAggregation) parse(raw map[string]types.Aggregate) (MultiTermsResult, error) {
	agg, err := requireAggregate[*types.MultiTermsAggregate](raw, a.name)
	if err != nil {
		return MultiTermsResult{}, err
	}

	buckets, ok := agg.Buckets.([]types.MultiTermsBucket)
	if !ok {
		return MultiTermsResult{}, fmt.Errorf("aggregation %q has unexpected buckets type %T", a.name, agg.Buckets)
	}

	out := make([]MultiTermsBucketResult, 0, len(buckets))
	for _, bucket := range buckets {
		keys := make([]string, 0, len(bucket.Key))
		for _, k := range bucket.Key {
			keys = append(keys, fmt.Sprintf("%v", k))
		}
		out = append(out, MultiTermsBucketResult{
			keys:         keys,
			docCount:     bucket.DocCount,
			aggregations: NewAggResults(bucket.Aggregations),
		})
	}
	return MultiTermsResult{buckets: out}, nil
}

func requireAggregate[T any](raw map[string]types.Aggregate, name string) (T, error) {
	var zero T
	if raw == nil {
		return zero, fmt.Errorf("aggregation %q not found", name)
	}
	value, ok := raw[name]
	if !ok {
		return zero, fmt.Errorf("aggregation %q not found", name)
	}
	typed, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("aggregation %q has unexpected type %T", name, value)
	}
	return typed, nil
}

func float64PtrFromFloat64(v *types.Float64) *float64 {
	if v == nil {
		return nil
	}
	f := float64(*v)
	return &f
}

func int64PtrFromFloat64(v *types.Float64) *int64 {
	if v == nil {
		return nil
	}
	i := int64(*v)
	return &i
}

func int64PtrFromInt64(v int64) *int64 {
	i := v
	return &i
}

func float64PtrFromFloat64Value(v types.Float64) *float64 {
	f := float64(v)
	return &f
}
