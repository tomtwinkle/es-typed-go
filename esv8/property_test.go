package esv8_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv8"
)

func TestNewBooleanProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewBooleanProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewBooleanProperty(
			esv8.WithBooleanDocValues(false),
			esv8.WithBooleanIndex(false),
			esv8.WithBooleanStore(true),
			esv8.WithBooleanNullValue(false),
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
		prop := esv8.NewKeywordProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with ignore above", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewKeywordProperty(esv8.WithKeywordIgnoreAbove(256))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 256, *prop.IgnoreAbove)
	})

	t.Run("with multiple options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewKeywordProperty(
			esv8.WithKeywordIgnoreAbove(256),
			esv8.WithKeywordDocValues(true),
			esv8.WithKeywordNormalizer("lowercase"),
			esv8.WithKeywordStore(true),
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
	prop := esv8.NewConstantKeywordProperty()
	assert.Assert(t, prop != nil)
}

func TestNewCountedKeywordProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewCountedKeywordProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with index", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewCountedKeywordProperty(esv8.WithCountedKeywordIndex(true))
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Index)
	})
}

func TestNewWildcardProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewWildcardProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with ignore above", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewWildcardProperty(esv8.WithWildcardIgnoreAbove(256))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 256, *prop.IgnoreAbove)
	})
}

func TestNewTextProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewTextProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with raw keyword", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewTextProperty(esv8.WithTextFields(map[string]types.Property{
			"keyword": esv8.NewKeywordProperty(esv8.WithKeywordIgnoreAbove(256)),
		}))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Fields != nil)
		_, ok := prop.Fields["keyword"]
		assert.Assert(t, ok)
	})

	t.Run("with analyzer", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewTextProperty(
			esv8.WithTextAnalyzer("standard"),
			esv8.WithTextSearchAnalyzer("standard"),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, "standard", *prop.Analyzer)
		assert.Equal(t, "standard", *prop.SearchAnalyzer)
	})
}

func TestNewMatchOnlyTextProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewMatchOnlyTextProperty()
	assert.Assert(t, prop != nil)
}

func TestNewCompletionProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewCompletionProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewCompletionProperty(
			esv8.WithCompletionAnalyzer("standard"),
			esv8.WithCompletionMaxInputLength(50),
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
		prop := esv8.NewSearchAsYouTypeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with max shingle size", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewSearchAsYouTypeProperty(esv8.WithSearchAsYouTypeMaxShingleSize(3))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 3, *prop.MaxShingleSize)
	})
}

func TestNewIntegerNumberProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIntegerNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with coerce", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIntegerNumberProperty(esv8.WithIntegerNumberCoerce(true))
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
	})
}

func TestNewLongNumberProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewLongNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewLongNumberProperty(
			esv8.WithLongNumberDocValues(true),
			esv8.WithLongNumberStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.DocValues)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewShortNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewShortNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewByteNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewByteNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewDoubleNumberProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDoubleNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with coerce", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDoubleNumberProperty(esv8.WithDoubleNumberCoerce(true))
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
	})
}

func TestNewFloatNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewFloatNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewHalfFloatNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewHalfFloatNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewUnsignedLongNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewUnsignedLongNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewScaledFloatNumberProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewScaledFloatNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with scaling factor", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewScaledFloatNumberProperty(esv8.WithScaledFloatNumberScalingFactor(100))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.ScalingFactor != nil)
	})
}

func TestNewDateProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateProperty()
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format == nil)
	})

	t.Run("with format", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateProperty(esv8.WithDateFormat(estype.DateFormatStrictDate))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date", *prop.Format)
	})

	t.Run("with multiple formats", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateProperty(esv8.WithDateFormat(estype.DateFormatStrictDateOptionalTime, estype.DateFormatEpochMillis))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date_optional_time||epoch_millis", *prop.Format)
	})
}

func TestNewDateNanosProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateNanosProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with format", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateNanosProperty(esv8.WithDateNanosFormat(estype.DateFormatStrictDateOptionalTimeNanos))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date_optional_time_nanos", *prop.Format)
	})
}

func TestNewGeoPointProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewGeoPointProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewGeoPointProperty(
			esv8.WithGeoPointIgnoreMalformed(true),
			esv8.WithGeoPointIgnoreZValue(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, true, *prop.IgnoreZValue)
	})
}

func TestNewGeoShapeProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewGeoShapeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewShapeProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewShapeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewPointProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewPointProperty()
	assert.Assert(t, prop != nil)
}

func TestNewIntegerRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewIntegerRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewLongRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewLongRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewFloatRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewFloatRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewDoubleRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewDoubleRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewDateRangeProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with format", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateRangeProperty(esv8.WithDateRangeFormat(estype.DateFormatStrictDate))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date", *prop.Format)
	})
}

func TestNewIpRangeProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewIpRangeProperty()
	assert.Assert(t, prop != nil)
}

func TestNewObjectProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewObjectProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with properties", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewObjectProperty(esv8.WithObjectProperties(map[string]types.Property{
			"name": esv8.NewKeywordProperty(esv8.WithKeywordIgnoreAbove(256)),
		}))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Properties != nil)
		_, ok := prop.Properties["name"]
		assert.Assert(t, ok)
	})

	t.Run("with enabled", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewObjectProperty(esv8.WithObjectEnabled(false))
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.Enabled)
	})
}

func TestNewNestedProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewNestedProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with properties", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewNestedProperty(esv8.WithNestedProperties(map[string]types.Property{
			"name": esv8.NewKeywordProperty(esv8.WithKeywordIgnoreAbove(256)),
		}))
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Properties != nil)
		_, ok := prop.Properties["name"]
		assert.Assert(t, ok)
	})

	t.Run("with include in parent", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewNestedProperty(
			esv8.WithNestedIncludeInParent(true),
			esv8.WithNestedIncludeInRoot(true),
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
		prop := esv8.NewFlattenedProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with depth limit", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewFlattenedProperty(esv8.WithFlattenedDepthLimit(5))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 5, *prop.DepthLimit)
	})
}

func TestNewJoinProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewJoinProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with relations", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewJoinProperty(esv8.WithJoinRelations(map[string][]string{
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
	prop := esv8.NewPassthroughObjectProperty()
	assert.Assert(t, prop != nil)
}

func TestNewIpProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIpProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with null value", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIpProperty(esv8.WithIpNullValue("0.0.0.0"))
		assert.Assert(t, prop != nil)
		assert.Equal(t, "0.0.0.0", *prop.NullValue)
	})
}

func TestNewBinaryProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewBinaryProperty()
	assert.Assert(t, prop != nil)
}

func TestNewTokenCountProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewTokenCountProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with analyzer", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewTokenCountProperty(esv8.WithTokenCountAnalyzer("standard"))
		assert.Assert(t, prop != nil)
		assert.Equal(t, "standard", *prop.Analyzer)
	})
}

func TestNewPercolatorProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewPercolatorProperty()
	assert.Assert(t, prop != nil)
}

func TestNewFieldAliasProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewFieldAliasProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with path", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewFieldAliasProperty(esv8.WithFieldAliasPath("title"))
		assert.Assert(t, prop != nil)
		assert.Equal(t, "title", *prop.Path)
	})
}

func TestNewHistogramProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewHistogramProperty()
	assert.Assert(t, prop != nil)
}

func TestNewVersionProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewVersionProperty()
	assert.Assert(t, prop != nil)
}

func TestNewDenseVectorProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDenseVectorProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with dims", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDenseVectorProperty(
			esv8.WithDenseVectorDims(384),
			esv8.WithDenseVectorIndex(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, 384, *prop.Dims)
		assert.Equal(t, true, *prop.Index)
	})
}

func TestNewSparseVectorProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewSparseVectorProperty()
	assert.Assert(t, prop != nil)
}

func TestNewRankFeatureProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewRankFeatureProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with positive score impact", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewRankFeatureProperty(esv8.WithRankFeaturePositiveScoreImpact(true))
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.PositiveScoreImpact)
	})
}

func TestNewRankFeaturesProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewRankFeaturesProperty()
	assert.Assert(t, prop != nil)
}

func TestNewRankVectorProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewRankVectorProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with dims", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewRankVectorProperty(esv8.WithRankVectorDims(128))
		assert.Assert(t, prop != nil)
		assert.Equal(t, 128, *prop.Dims)
	})
}

func TestNewSemanticTextProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewSemanticTextProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with inference id", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewSemanticTextProperty(esv8.WithSemanticTextInferenceId("my-model"))
		assert.Assert(t, prop != nil)
		assert.Equal(t, "my-model", *prop.InferenceId)
	})
}

func TestNewAggregateMetricDoubleProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewAggregateMetricDoubleProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewAggregateMetricDoubleProperty(
			esv8.WithAggregateMetricDoubleDefaultMetric("max"),
			esv8.WithAggregateMetricDoubleMetrics([]string{"min", "max", "sum", "value_count"}),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, "max", prop.DefaultMetric)
		assert.Equal(t, 4, len(prop.Metrics))
	})
}

func TestNewMurmur3HashProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewMurmur3HashProperty()
	assert.Assert(t, prop != nil)
}

func TestNewIcuCollationProperty(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIcuCollationProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with language", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIcuCollationProperty(
			esv8.WithIcuCollationLanguage("ja"),
			esv8.WithIcuCollationCountry("JP"),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, "ja", *prop.Language)
		assert.Equal(t, "JP", *prop.Country)
	})
}

func TestNewDynamicProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewDynamicProperty()
	assert.Assert(t, prop != nil)
}

// TestCustomOption tests that users can write custom functional options.
func TestCustomOption(t *testing.T) {
	t.Parallel()
	prop := esv8.NewKeywordProperty(func(p *types.KeywordProperty) {
		v := true
		p.EagerGlobalOrdinals = &v
	})
	assert.Assert(t, prop != nil)
	assert.Equal(t, true, *prop.EagerGlobalOrdinals)
}
