package esv9_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv9"
)

func TestNewBooleanProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewBooleanProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewBooleanProperty(
			esv9.WithBooleanDocValues(false),
			esv9.WithBooleanIndex(false),
			esv9.WithBooleanStore(true),
			esv9.WithBooleanNullValue(false),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
		assert.Equal(t, false, *prop.NullValue)
	})
}

func TestNewKeywordProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewKeywordProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with ignore above", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewKeywordProperty(esv9.WithKeywordIgnoreAbove(256))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 256, *prop.IgnoreAbove)
	})

	t.Run("with multiple options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewKeywordProperty(
			esv9.WithKeywordIgnoreAbove(256),
			esv9.WithKeywordDocValues(true),
			esv9.WithKeywordNormalizer("lowercase"),
			esv9.WithKeywordStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, 256, *prop.IgnoreAbove)
		assert.Equal(t, true, *prop.DocValues)
		assert.Equal(t, "lowercase", *prop.Normalizer)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewConstantKeywordProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewConstantKeywordProperty()
	assert.Assert(t, prop != nil)
}

func TestNewCountedKeywordProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewCountedKeywordProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with index", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewCountedKeywordProperty(esv9.WithCountedKeywordIndex(true))
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Index)
	})
}

func TestNewWildcardProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewWildcardProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with ignore above", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewWildcardProperty(esv9.WithWildcardIgnoreAbove(256))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 256, *prop.IgnoreAbove)
	})
}

func TestNewTextProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewTextProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with raw keyword", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewTextProperty(esv9.WithTextFields(map[string]types.Property{
			"keyword": esv9.NewKeywordProperty(esv9.WithKeywordIgnoreAbove(256)),
		}))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Fields != nil)
		_, ok := prop.Fields["keyword"]
		assert.Assert(t, ok)
	})

	t.Run("with analyzer", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewTextProperty(
			esv9.WithTextAnalyzer("standard"),
			esv9.WithTextSearchAnalyzer("standard"),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, "standard", *prop.Analyzer)
		assert.Equal(t, "standard", *prop.SearchAnalyzer)
	})
}

func TestNewMatchOnlyTextProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewMatchOnlyTextProperty()
	assert.Assert(t, prop != nil)
}

func TestNewCompletionProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewCompletionProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewCompletionProperty(
			esv9.WithCompletionAnalyzer("standard"),
			esv9.WithCompletionMaxInputLength(50),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, "standard", *prop.Analyzer)
		assert.Equal(t, 50, *prop.MaxInputLength)
	})
}

func TestNewSearchAsYouTypeProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewSearchAsYouTypeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with max shingle size", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewSearchAsYouTypeProperty(esv9.WithSearchAsYouTypeMaxShingleSize(3))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 3, *prop.MaxShingleSize)
	})
}

func TestNewIntegerNumberProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewIntegerNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with coerce", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewIntegerNumberProperty(esv9.WithIntegerNumberCoerce(true))
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
	})
}

func TestNewLongNumberProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewLongNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewLongNumberProperty(
			esv9.WithLongNumberDocValues(true),
			esv9.WithLongNumberStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.DocValues)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewShortNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewShortNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewByteNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewByteNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewDoubleNumberProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDoubleNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with coerce", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDoubleNumberProperty(esv9.WithDoubleNumberCoerce(true))
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
	})
}

func TestNewFloatNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewFloatNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewHalfFloatNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewHalfFloatNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewUnsignedLongNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewUnsignedLongNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewScaledFloatNumberProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewScaledFloatNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with scaling factor", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewScaledFloatNumberProperty(esv9.WithScaledFloatNumberScalingFactor(100))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.ScalingFactor != nil)
	})
}

func TestNewDateProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDateProperty()
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format == nil)
	})

	t.Run("with format", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDateProperty(esv9.WithDateFormat(estype.DateFormatStrictDate))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date", *prop.Format)
	})

	t.Run("with multiple formats", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDateProperty(esv9.WithDateFormat(estype.DateFormatStrictDateOptionalTime, estype.DateFormatEpochMillis))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date_optional_time||epoch_millis", *prop.Format)
	})
}

func TestNewDateNanosProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDateNanosProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with format", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDateNanosProperty(esv9.WithDateNanosFormat(estype.DateFormatStrictDateOptionalTimeNanos))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date_optional_time_nanos", *prop.Format)
	})
}

func TestNewGeoPointProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewGeoPointProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewGeoPointProperty(
			esv9.WithGeoPointIgnoreMalformed(true),
			esv9.WithGeoPointIgnoreZValue(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, true, *prop.IgnoreZValue)
	})
}

func TestNewGeoShapeProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewGeoShapeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewShapeProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewShapeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewPointProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewPointProperty()
	assert.Assert(t, prop != nil)
}

func TestNewIntegerRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewIntegerRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewLongRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewLongRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewFloatRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewFloatRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewDoubleRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewDoubleRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewDateRangeProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDateRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with format", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDateRangeProperty(esv9.WithDateRangeFormat(estype.DateFormatStrictDate))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date", *prop.Format)
	})
}

func TestNewIpRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewIpRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewObjectProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewObjectProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with properties", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewObjectProperty(esv9.WithObjectProperties(map[string]types.Property{
			"name": esv9.NewKeywordProperty(esv9.WithKeywordIgnoreAbove(256)),
		}))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Properties != nil)
		_, ok := prop.Properties["name"]
		assert.Assert(t, ok)
	})

	t.Run("with enabled", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewObjectProperty(esv9.WithObjectEnabled(false))
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.Enabled)
	})
}

func TestNewNestedProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewNestedProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with properties", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewNestedProperty(esv9.WithNestedProperties(map[string]types.Property{
			"name": esv9.NewKeywordProperty(esv9.WithKeywordIgnoreAbove(256)),
		}))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Properties != nil)
		_, ok := prop.Properties["name"]
		assert.Assert(t, ok)
	})

	t.Run("with include in parent", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewNestedProperty(
			esv9.WithNestedIncludeInParent(true),
			esv9.WithNestedIncludeInRoot(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.IncludeInParent)
		assert.Equal(t, true, *prop.IncludeInRoot)
	})
}

func TestNewFlattenedProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewFlattenedProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with depth limit", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewFlattenedProperty(esv9.WithFlattenedDepthLimit(5))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 5, *prop.DepthLimit)
	})
}

func TestNewJoinProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewJoinProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with relations", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewJoinProperty(esv9.WithJoinRelations(map[string][]string{
			"parent": {"child"},
		}))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Relations != nil)
		assert.Assert(t, len(prop.Relations["parent"]) == 1)
		assert.Equal(t, "child", prop.Relations["parent"][0])
	})
}

func TestNewPassthroughObjectProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewPassthroughObjectProperty()
	assert.Assert(t, prop != nil)
}

func TestNewIpProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewIpProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with null value", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewIpProperty(esv9.WithIpNullValue("0.0.0.0"))
		assert.Assert(t, prop != nil)
		assert.Equal(t, "0.0.0.0", *prop.NullValue)
	})
}

func TestNewBinaryProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewBinaryProperty()
	assert.Assert(t, prop != nil)
}

func TestNewTokenCountProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewTokenCountProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with analyzer", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewTokenCountProperty(esv9.WithTokenCountAnalyzer("standard"))
		assert.Assert(t, prop != nil)
		assert.Equal(t, "standard", *prop.Analyzer)
	})
}

func TestNewPercolatorProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewPercolatorProperty()
	assert.Assert(t, prop != nil)
}

func TestNewFieldAliasProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewFieldAliasProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with path", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewFieldAliasProperty(esv9.WithFieldAliasPath("title"))
		assert.Assert(t, prop != nil)
		assert.Equal(t, "title", *prop.Path)
	})
}

func TestNewHistogramProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewHistogramProperty()
	assert.Assert(t, prop != nil)
}

func TestNewVersionProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewVersionProperty()
	assert.Assert(t, prop != nil)
}

func TestNewDenseVectorProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDenseVectorProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with dims", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDenseVectorProperty(
			esv9.WithDenseVectorDims(384),
			esv9.WithDenseVectorIndex(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, 384, *prop.Dims)
		assert.Equal(t, true, *prop.Index)
	})
}

func TestNewSparseVectorProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewSparseVectorProperty()
	assert.Assert(t, prop != nil)
}

func TestNewRankFeatureProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewRankFeatureProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with positive score impact", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewRankFeatureProperty(esv9.WithRankFeaturePositiveScoreImpact(true))
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.PositiveScoreImpact)
	})
}

func TestNewRankFeaturesProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewRankFeaturesProperty()
	assert.Assert(t, prop != nil)
}

func TestNewRankVectorProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewRankVectorProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with dims", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewRankVectorProperty(esv9.WithRankVectorDims(128))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 128, *prop.Dims)
	})
}

func TestNewSemanticTextProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewSemanticTextProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with inference id", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewSemanticTextProperty(esv9.WithSemanticTextInferenceId("my-model"))
		assert.Assert(t, prop != nil)
		assert.Equal(t, "my-model", *prop.InferenceId)
	})
}

func TestNewAggregateMetricDoubleProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewAggregateMetricDoubleProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewAggregateMetricDoubleProperty(
			esv9.WithAggregateMetricDoubleDefaultMetric("max"),
			esv9.WithAggregateMetricDoubleMetrics([]string{"min", "max", "sum", "value_count"}),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, "max", prop.DefaultMetric)
		assert.Equal(t, 4, len(prop.Metrics))
	})
}

func TestNewMurmur3HashProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewMurmur3HashProperty()
	assert.Assert(t, prop != nil)
}

func TestNewIcuCollationProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewIcuCollationProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with language", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewIcuCollationProperty(
			esv9.WithIcuCollationLanguage("ja"),
			esv9.WithIcuCollationCountry("JP"),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, "ja", *prop.Language)
		assert.Equal(t, "JP", *prop.Country)
	})
}

func TestNewDynamicProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewDynamicProperty()
	assert.Assert(t, prop != nil)
}

// TestCustomOption tests that users can write custom functional options.
func TestCustomOption(t *testing.T) {
	t.Parallel()
	prop := esv9.NewKeywordProperty(func(p *types.KeywordProperty) {
		v := true
		p.EagerGlobalOrdinals = &v
	})
	assert.Assert(t, prop != nil)
	assert.Equal(t, true, *prop.EagerGlobalOrdinals)
}

func TestNewExponentialHistogramProperty(t *testing.T) {
	t.Parallel()
	prop := esv9.NewExponentialHistogramProperty()
	assert.Assert(t, prop != nil)
}
