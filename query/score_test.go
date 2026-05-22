package query_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/fieldvaluefactormodifier"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/functionboostmode"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/functionscoremode"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/query"
)

func TestFunctionScoreBuilder(t *testing.T) {
	t.Parallel()

	q := query.NewFunctionScore().
		Query(query.MatchAll()).
		Add(
			query.NewScoreFunction().
				Filter(query.TermValue(FieldCategory, "electronics")).
				Weight(2.5).
				Build(),
		).
		ScoreMode(functionscoremode.Multiply).
		BoostMode(functionboostmode.Multiply).
		MaxBoost(10).
		MinScore(2).
		Boost(1.5).
		BuildQuery()

	assert.Assert(t, q.FunctionScore != nil)
	assert.Assert(t, q.FunctionScore.Query != nil)
	assert.Assert(t, q.FunctionScore.Query.MatchAll != nil)
	assert.Equal(t, 1, len(q.FunctionScore.Functions))
	assert.Assert(t, q.FunctionScore.ScoreMode != nil)
	assert.Equal(t, functionscoremode.Multiply, *q.FunctionScore.ScoreMode)
	assert.Assert(t, q.FunctionScore.BoostMode != nil)
	assert.Equal(t, functionboostmode.Multiply, *q.FunctionScore.BoostMode)
	assert.Assert(t, q.FunctionScore.MaxBoost != nil)
	assert.Equal(t, 10.0, float64(*q.FunctionScore.MaxBoost))
	assert.Assert(t, q.FunctionScore.MinScore != nil)
	assert.Equal(t, 2.0, float64(*q.FunctionScore.MinScore))
	assert.Assert(t, q.FunctionScore.Boost != nil)
	assert.Equal(t, float32(1.5), *q.FunctionScore.Boost)
}

func TestScoreFunctionBuilder_FieldValueFactor(t *testing.T) {
	t.Parallel()

	fn := query.NewScoreFunction().
		Filter(query.TermValue(FieldCategory, "electronics")).
		FieldValueFactor(
			FieldPrice,
			query.WithFieldValueFactorFactor(2),
			query.WithFieldValueFactorMissing(1),
			query.WithFieldValueFactorModifier(fieldvaluefactormodifier.Sqrt),
		).
		Weight(3).
		Build()

	assert.Assert(t, fn.Filter != nil)
	assert.Assert(t, fn.Filter.Term != nil)
	assert.Assert(t, fn.FieldValueFactor != nil)
	assert.Equal(t, string(FieldPrice), fn.FieldValueFactor.Field)
	assert.Assert(t, fn.FieldValueFactor.Factor != nil)
	assert.Equal(t, 2.0, float64(*fn.FieldValueFactor.Factor))
	assert.Assert(t, fn.FieldValueFactor.Missing != nil)
	assert.Equal(t, 1.0, float64(*fn.FieldValueFactor.Missing))
	assert.Assert(t, fn.FieldValueFactor.Modifier != nil)
	assert.Equal(t, fieldvaluefactormodifier.Sqrt, *fn.FieldValueFactor.Modifier)
	assert.Assert(t, fn.Weight != nil)
	assert.Equal(t, 3.0, float64(*fn.Weight))
}

func TestScoreFunctionBuilder_RandomScore(t *testing.T) {
	t.Parallel()

	fn := query.NewScoreFunction().
		RandomScore(
			query.WithRandomScoreField(FieldId),
			query.WithRandomScoreSeed("fixed-seed"),
		).
		Build()

	assert.Assert(t, fn.RandomScore != nil)
	assert.Assert(t, fn.RandomScore.Field != nil)
	assert.Equal(t, string(FieldId), *fn.RandomScore.Field)
	assert.Assert(t, fn.RandomScore.Seed != nil)
	assert.Equal(t, "fixed-seed", *fn.RandomScore.Seed)
}

func TestScoreFunctionBuilder_ScriptScoreSource(t *testing.T) {
	t.Parallel()

	fn := query.NewScoreFunction().
		ScriptScoreSource("doc['price'].value").
		Build()

	assert.Assert(t, fn.ScriptScore != nil)
	assert.Assert(t, fn.ScriptScore.Script.Source != nil)
	assert.Equal(t, "doc['price'].value", *fn.ScriptScore.Script.Source)
}

func TestInlineScript(t *testing.T) {
	t.Parallel()

	script := query.InlineScript("doc['price'].value")

	assert.Assert(t, script.Source != nil)
	assert.Equal(t, "doc['price'].value", *script.Source)
}

func TestScriptScoreQueryBuilder(t *testing.T) {
	t.Parallel()

	q := query.NewScriptScoreQuery().
		Query(query.MatchAll()).
		Script(query.InlineScript("doc['price'].value")).
		MinScore(100).
		Boost(2).
		BuildQuery()

	assert.Assert(t, q.ScriptScore != nil)
	assert.Assert(t, q.ScriptScore.Query.MatchAll != nil)
	assert.Assert(t, q.ScriptScore.Script.Source != nil)
	assert.Equal(t, "doc['price'].value", *q.ScriptScore.Script.Source)
	assert.Assert(t, q.ScriptScore.MinScore != nil)
	assert.Equal(t, float32(100), *q.ScriptScore.MinScore)
	assert.Assert(t, q.ScriptScore.Boost != nil)
	assert.Equal(t, float32(2), *q.ScriptScore.Boost)
}

func TestScriptScoreQueryHelper(t *testing.T) {
	t.Parallel()

	source := "doc['price'].value"
	q := query.ScriptScoreQuery(&types.ScriptScoreQuery{
		Query:  query.MatchAll(),
		Script: types.Script{Source: &source},
	})

	assert.Assert(t, q.ScriptScore != nil)
	assert.Assert(t, q.ScriptScore.Query.MatchAll != nil)
}
