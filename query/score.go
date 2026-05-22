package query

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/fieldvaluefactormodifier"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/functionboostmode"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/functionscoremode"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// FunctionScoreBuilder constructs a FunctionScoreQuery using method chaining.
type FunctionScoreBuilder struct {
	q types.FunctionScoreQuery
}

// NewFunctionScore returns a new empty FunctionScoreBuilder.
func NewFunctionScore() *FunctionScoreBuilder {
	return &FunctionScoreBuilder{}
}

// Query sets the inner query whose matching documents will be rescored.
func (b *FunctionScoreBuilder) Query(q types.Query) *FunctionScoreBuilder {
	b.q.Query = &q
	return b
}

// Add appends one or more score functions.
func (b *FunctionScoreBuilder) Add(functions ...types.FunctionScore) *FunctionScoreBuilder {
	b.q.Functions = append(b.q.Functions, functions...)
	return b
}

// ScoreMode sets how multiple function scores are combined.
func (b *FunctionScoreBuilder) ScoreMode(mode functionscoremode.FunctionScoreMode) *FunctionScoreBuilder {
	b.q.ScoreMode = &mode
	return b
}

// BoostMode sets how the combined function score interacts with the query score.
func (b *FunctionScoreBuilder) BoostMode(mode functionboostmode.FunctionBoostMode) *FunctionScoreBuilder {
	b.q.BoostMode = &mode
	return b
}

// MaxBoost sets an upper bound for the computed score.
func (b *FunctionScoreBuilder) MaxBoost(v float64) *FunctionScoreBuilder {
	value := types.Float64(v)
	b.q.MaxBoost = &value
	return b
}

// MinScore filters out hits below the given score threshold.
func (b *FunctionScoreBuilder) MinScore(v float64) *FunctionScoreBuilder {
	value := types.Float64(v)
	b.q.MinScore = &value
	return b
}

// Boost applies a boost to the whole function_score query.
func (b *FunctionScoreBuilder) Boost(v float32) *FunctionScoreBuilder {
	b.q.Boost = &v
	return b
}

// Build returns the constructed FunctionScoreQuery.
func (b *FunctionScoreBuilder) Build() *types.FunctionScoreQuery {
	q := b.q
	return &q
}

// BuildQuery returns the constructed function_score wrapped as a Query.
func (b *FunctionScoreBuilder) BuildQuery() types.Query {
	return FunctionScoreQuery(b.Build())
}

// ScoreFunctionBuilder constructs an individual FunctionScore.
type ScoreFunctionBuilder struct {
	fn types.FunctionScore
}

// NewScoreFunction returns a new empty ScoreFunctionBuilder.
func NewScoreFunction() *ScoreFunctionBuilder {
	return &ScoreFunctionBuilder{}
}

// Filter restricts the score function to documents matching the given query.
func (b *ScoreFunctionBuilder) Filter(q types.Query) *ScoreFunctionBuilder {
	b.fn.Filter = &q
	return b
}

// Weight sets a static weight for the score function.
func (b *ScoreFunctionBuilder) Weight(v float64) *ScoreFunctionBuilder {
	value := types.Float64(v)
	b.fn.Weight = &value
	return b
}

// FieldValueFactor sets a field_value_factor score function using a typed field.
func (b *ScoreFunctionBuilder) FieldValueFactor(field estype.Field, opts ...FieldValueFactorOption) *ScoreFunctionBuilder {
	b.fn.FieldValueFactor = FieldValueFactor(field, opts...)
	return b
}

// RandomScore sets a random_score function.
func (b *ScoreFunctionBuilder) RandomScore(opts ...RandomScoreOption) *ScoreFunctionBuilder {
	b.fn.RandomScore = RandomScore(opts...)
	return b
}

// ScriptScore sets a script_score function.
func (b *ScoreFunctionBuilder) ScriptScore(script types.Script) *ScoreFunctionBuilder {
	b.fn.ScriptScore = ScriptScore(script)
	return b
}

// ScriptScoreSource sets an inline script_score function from source text.
func (b *ScoreFunctionBuilder) ScriptScoreSource(source string) *ScoreFunctionBuilder {
	b.fn.ScriptScore = ScriptScore(InlineScript(source))
	return b
}

// Build returns the constructed FunctionScore.
func (b *ScoreFunctionBuilder) Build() types.FunctionScore {
	return b.fn
}

// FieldValueFactorOption configures a field_value_factor score function.
type FieldValueFactorOption func(*types.FieldValueFactorScoreFunction)

// FieldValueFactor creates a field_value_factor score function using a typed field.
func FieldValueFactor(field estype.Field, opts ...FieldValueFactorOption) *types.FieldValueFactorScoreFunction {
	fn := types.NewFieldValueFactorScoreFunction()
	fn.Field = string(field)
	for _, opt := range opts {
		opt(fn)
	}
	return fn
}

// WithFieldValueFactorFactor multiplies the field value by the given factor.
func WithFieldValueFactorFactor(v float64) FieldValueFactorOption {
	return func(fn *types.FieldValueFactorScoreFunction) {
		value := types.Float64(v)
		fn.Factor = &value
	}
}

// WithFieldValueFactorMissing sets the fallback value used when the field is absent.
func WithFieldValueFactorMissing(v float64) FieldValueFactorOption {
	return func(fn *types.FieldValueFactorScoreFunction) {
		value := types.Float64(v)
		fn.Missing = &value
	}
}

// WithFieldValueFactorModifier applies a modifier to the field value.
func WithFieldValueFactorModifier(modifier fieldvaluefactormodifier.FieldValueFactorModifier) FieldValueFactorOption {
	return func(fn *types.FieldValueFactorScoreFunction) {
		fn.Modifier = &modifier
	}
}

// RandomScoreOption configures a random_score function.
type RandomScoreOption func(*types.RandomScoreFunction)

// RandomScore creates a random_score function.
func RandomScore(opts ...RandomScoreOption) *types.RandomScoreFunction {
	fn := types.NewRandomScoreFunction()
	for _, opt := range opts {
		opt(fn)
	}
	return fn
}

// WithRandomScoreField sets the field used to stabilize random scoring.
func WithRandomScoreField(field estype.Field) RandomScoreOption {
	return func(fn *types.RandomScoreFunction) {
		value := string(field)
		fn.Field = &value
	}
}

// WithRandomScoreSeed sets the seed used for random scoring.
func WithRandomScoreSeed(seed string) RandomScoreOption {
	return func(fn *types.RandomScoreFunction) {
		fn.Seed = &seed
	}
}

// InlineScript creates an inline script from source text.
func InlineScript(source string) types.Script {
	return types.Script{Source: &source}
}

// ScriptScore creates a script_score function from a script.
func ScriptScore(script types.Script) *types.ScriptScoreFunction {
	return &types.ScriptScoreFunction{Script: script}
}

// ScriptScoreQueryBuilder constructs a ScriptScoreQuery using method chaining.
type ScriptScoreQueryBuilder struct {
	q types.ScriptScoreQuery
}

// NewScriptScoreQuery returns a new empty ScriptScoreQueryBuilder.
func NewScriptScoreQuery() *ScriptScoreQueryBuilder {
	return &ScriptScoreQueryBuilder{}
}

// Query sets the inner query whose hits will be scored by the script.
func (b *ScriptScoreQueryBuilder) Query(q types.Query) *ScriptScoreQueryBuilder {
	b.q.Query = q
	return b
}

// Script sets the script used to compute scores.
func (b *ScriptScoreQueryBuilder) Script(script types.Script) *ScriptScoreQueryBuilder {
	b.q.Script = script
	return b
}

// ScriptSource sets an inline script from source text.
func (b *ScriptScoreQueryBuilder) ScriptSource(source string) *ScriptScoreQueryBuilder {
	b.q.Script = InlineScript(source)
	return b
}

// MinScore filters out hits below the given score threshold.
func (b *ScriptScoreQueryBuilder) MinScore(v float32) *ScriptScoreQueryBuilder {
	b.q.MinScore = &v
	return b
}

// Boost applies a boost to the whole script_score query.
func (b *ScriptScoreQueryBuilder) Boost(v float32) *ScriptScoreQueryBuilder {
	b.q.Boost = &v
	return b
}

// Build returns the constructed ScriptScoreQuery.
func (b *ScriptScoreQueryBuilder) Build() *types.ScriptScoreQuery {
	q := b.q
	return &q
}

// BuildQuery returns the constructed script_score wrapped as a Query.
func (b *ScriptScoreQueryBuilder) BuildQuery() types.Query {
	return ScriptScoreQuery(b.Build())
}

// ScriptScoreQuery creates a Query with a ScriptScoreQuery.
func ScriptScoreQuery(ssq *types.ScriptScoreQuery) types.Query {
	return types.Query{ScriptScore: ssq}
}
