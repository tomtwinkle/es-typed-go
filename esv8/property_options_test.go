package esv8_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/esv8"
)

func TestNewKeywordProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewKeywordProperty(
		esv8.WithKeywordIndex(false),
		esv8.WithKeywordNullValue("N/A"),
		esv8.WithKeywordNorms(false),
		esv8.WithKeywordSimilarity("BM25"),
		esv8.WithKeywordEagerGlobalOrdinals(true),
		esv8.WithKeywordSplitQueriesOnWhitespace(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, "N/A", *prop.NullValue)
	assert.Equal(t, false, *prop.Norms)
	assert.Equal(t, "BM25", *prop.Similarity)
	assert.Equal(t, true, *prop.EagerGlobalOrdinals)
	assert.Equal(t, true, *prop.SplitQueriesOnWhitespace)
}

func TestNewConstantKeywordProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewConstantKeywordProperty(func(p *types.ConstantKeywordProperty) {
		v := 256
		p.IgnoreAbove = &v
	})
	assert.Assert(t, prop != nil)
	assert.Equal(t, 256, *prop.IgnoreAbove)
}

func TestNewWildcardProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewWildcardProperty(
		esv8.WithWildcardDocValues(false),
		esv8.WithWildcardNullValue(""),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, "", *prop.NullValue)
}

func TestNewTextProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("with extra options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewTextProperty(
			esv8.WithTextSearchQuoteAnalyzer("standard"),
			esv8.WithTextFielddata(true),
			esv8.WithTextIndex(false),
			esv8.WithTextStore(true),
			esv8.WithTextNorms(false),
			esv8.WithTextSimilarity("BM25"),
			esv8.WithTextIndexPhrases(true),
			esv8.WithTextPositionIncrementGap(100),
			esv8.WithTextFields(map[string]types.Property{
				"value": esv8.NewKeywordProperty(),
			}),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, "standard", *prop.SearchQuoteAnalyzer)
		assert.Equal(t, true, *prop.Fielddata)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
		assert.Equal(t, false, *prop.Norms)
		assert.Equal(t, "BM25", *prop.Similarity)
		assert.Equal(t, true, *prop.IndexPhrases)
		assert.Equal(t, 100, *prop.PositionIncrementGap)
		_, ok := prop.Fields["value"]
		assert.Assert(t, ok)
	})

	t.Run("with raw keyword on nil fields", func(t *testing.T) {
		t.Parallel()
		// Use a manually created property with nil Fields to exercise the make() branch.
		p := &types.TextProperty{}
		esv8.WithTextRawKeyword(256)(p)
		assert.Assert(t, p.Fields != nil)
		_, ok := p.Fields["keyword"]
		assert.Assert(t, ok)
	})
}

func TestNewMatchOnlyTextProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewMatchOnlyTextProperty(func(p *types.MatchOnlyTextProperty) {
		p.CopyTo = []string{"title"}
	})
	assert.Assert(t, prop != nil)
	assert.Equal(t, 1, len(prop.CopyTo))
}

func TestNewCompletionProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewCompletionProperty(
		esv8.WithCompletionSearchAnalyzer("standard"),
		esv8.WithCompletionPreservePositionIncrements(true),
		esv8.WithCompletionPreserveSeparators(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, "standard", *prop.SearchAnalyzer)
	assert.Equal(t, true, *prop.PreservePositionIncrements)
	assert.Equal(t, true, *prop.PreserveSeparators)
}

func TestNewSearchAsYouTypeProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewSearchAsYouTypeProperty(
		esv8.WithSearchAsYouTypeAnalyzer("standard"),
		esv8.WithSearchAsYouTypeSearchAnalyzer("standard"),
		esv8.WithSearchAsYouTypeSearchQuoteAnalyzer("standard"),
		esv8.WithSearchAsYouTypeIndex(false),
		esv8.WithSearchAsYouTypeStore(true),
		esv8.WithSearchAsYouTypeNorms(false),
		esv8.WithSearchAsYouTypeSimilarity("BM25"),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, "standard", *prop.Analyzer)
	assert.Equal(t, "standard", *prop.SearchAnalyzer)
	assert.Equal(t, "standard", *prop.SearchQuoteAnalyzer)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
	assert.Equal(t, false, *prop.Norms)
	assert.Equal(t, "BM25", *prop.Similarity)
}

func TestNewIntegerNumberProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewIntegerNumberProperty(
		esv8.WithIntegerNumberDocValues(false),
		esv8.WithIntegerNumberIgnoreMalformed(true),
		esv8.WithIntegerNumberIndex(false),
		esv8.WithIntegerNumberStore(true),
		esv8.WithIntegerNumberNullValue(0),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, true, *prop.IgnoreMalformed)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
	assert.Equal(t, 0, *prop.NullValue)
}

func TestNewLongNumberProperty_Options(t *testing.T) {
	t.Parallel()
	var nullVal int64 = 0
	prop := esv8.NewLongNumberProperty(
		esv8.WithLongNumberCoerce(true),
		esv8.WithLongNumberIgnoreMalformed(true),
		esv8.WithLongNumberIndex(false),
		esv8.WithLongNumberNullValue(nullVal),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, true, *prop.Coerce)
	assert.Equal(t, true, *prop.IgnoreMalformed)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, nullVal, *prop.NullValue)
}

func TestNewShortNumberProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewShortNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewShortNumberProperty(
			esv8.WithShortNumberCoerce(true),
			esv8.WithShortNumberDocValues(false),
			esv8.WithShortNumberIgnoreMalformed(true),
			esv8.WithShortNumberIndex(false),
			esv8.WithShortNumberStore(true),
			esv8.WithShortNumberNullValue(0),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
		assert.Equal(t, 0, *prop.NullValue)
	})
}

func TestNewByteNumberProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewByteNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		var nullVal byte = 0
		prop := esv8.NewByteNumberProperty(
			esv8.WithByteNumberCoerce(true),
			esv8.WithByteNumberDocValues(false),
			esv8.WithByteNumberIgnoreMalformed(true),
			esv8.WithByteNumberIndex(false),
			esv8.WithByteNumberStore(true),
			esv8.WithByteNumberNullValue(nullVal),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
		assert.Equal(t, nullVal, *prop.NullValue)
	})
}

func TestNewDoubleNumberProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewDoubleNumberProperty(
		esv8.WithDoubleNumberDocValues(false),
		esv8.WithDoubleNumberIgnoreMalformed(true),
		esv8.WithDoubleNumberIndex(false),
		esv8.WithDoubleNumberStore(true),
		esv8.WithDoubleNumberNullValue(1.5),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, true, *prop.IgnoreMalformed)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
	assert.Assert(t, prop.NullValue != nil)
}

func TestNewFloatNumberProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewFloatNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		var nullVal float32 = 1.5
		prop := esv8.NewFloatNumberProperty(
			esv8.WithFloatNumberCoerce(true),
			esv8.WithFloatNumberDocValues(false),
			esv8.WithFloatNumberIgnoreMalformed(true),
			esv8.WithFloatNumberIndex(false),
			esv8.WithFloatNumberStore(true),
			esv8.WithFloatNumberNullValue(nullVal),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
		assert.Equal(t, nullVal, *prop.NullValue)
	})
}

func TestNewHalfFloatNumberProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewHalfFloatNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		var nullVal float32 = 1.5
		prop := esv8.NewHalfFloatNumberProperty(
			esv8.WithHalfFloatNumberCoerce(true),
			esv8.WithHalfFloatNumberDocValues(false),
			esv8.WithHalfFloatNumberIgnoreMalformed(true),
			esv8.WithHalfFloatNumberIndex(false),
			esv8.WithHalfFloatNumberStore(true),
			esv8.WithHalfFloatNumberNullValue(nullVal),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
		assert.Equal(t, nullVal, *prop.NullValue)
	})
}

func TestNewUnsignedLongNumberProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewUnsignedLongNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		var nullVal uint64 = 0
		prop := esv8.NewUnsignedLongNumberProperty(
			esv8.WithUnsignedLongNumberDocValues(false),
			esv8.WithUnsignedLongNumberIgnoreMalformed(true),
			esv8.WithUnsignedLongNumberIndex(false),
			esv8.WithUnsignedLongNumberStore(true),
			esv8.WithUnsignedLongNumberNullValue(nullVal),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
		assert.Equal(t, nullVal, *prop.NullValue)
	})
}

func TestNewScaledFloatNumberProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewScaledFloatNumberProperty(
		esv8.WithScaledFloatNumberCoerce(true),
		esv8.WithScaledFloatNumberDocValues(false),
		esv8.WithScaledFloatNumberIgnoreMalformed(true),
		esv8.WithScaledFloatNumberIndex(false),
		esv8.WithScaledFloatNumberStore(true),
		esv8.WithScaledFloatNumberNullValue(1.5),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, true, *prop.Coerce)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, true, *prop.IgnoreMalformed)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
	assert.Assert(t, prop.NullValue != nil)
}

func TestNewDateProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewDateProperty(
		esv8.WithDateDocValues(false),
		esv8.WithDateIgnoreMalformed(true),
		esv8.WithDateIndex(false),
		esv8.WithDateStore(true),
		esv8.WithDateLocale("en"),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, true, *prop.IgnoreMalformed)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
	assert.Equal(t, "en", *prop.Locale)
}

func TestNewDateNanosProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewDateNanosProperty(
		esv8.WithDateNanosDocValues(false),
		esv8.WithDateNanosIgnoreMalformed(true),
		esv8.WithDateNanosIndex(false),
		esv8.WithDateNanosStore(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, true, *prop.IgnoreMalformed)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
}

func TestNewGeoPointProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewGeoPointProperty(
		esv8.WithGeoPointDocValues(false),
		esv8.WithGeoPointIndex(false),
		esv8.WithGeoPointStore(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
}

func TestNewGeoShapeProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewGeoShapeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewGeoShapeProperty(
			esv8.WithGeoShapeCoerce(true),
			esv8.WithGeoShapeIgnoreMalformed(true),
			esv8.WithGeoShapeIgnoreZValue(false),
			esv8.WithGeoShapeDocValues(false),
			esv8.WithGeoShapeIndex(false),
			esv8.WithGeoShapeStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, false, *prop.IgnoreZValue)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewShapeProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewShapeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewShapeProperty(
			esv8.WithShapeCoerce(true),
			esv8.WithShapeIgnoreMalformed(true),
			esv8.WithShapeIgnoreZValue(false),
			esv8.WithShapeDocValues(false),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, false, *prop.IgnoreZValue)
		assert.Equal(t, false, *prop.DocValues)
	})
}

func TestNewPointProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewPointProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewPointProperty(
			esv8.WithPointIgnoreMalformed(true),
			esv8.WithPointIgnoreZValue(false),
			esv8.WithPointDocValues(false),
			esv8.WithPointStore(true),
			esv8.WithPointNullValue("0,0"),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, false, *prop.IgnoreZValue)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.Store)
		assert.Equal(t, "0,0", *prop.NullValue)
	})
}

func TestNewIntegerRangeProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIntegerRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIntegerRangeProperty(
			esv8.WithIntegerRangeCoerce(true),
			esv8.WithIntegerRangeDocValues(false),
			esv8.WithIntegerRangeIndex(false),
			esv8.WithIntegerRangeStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewLongRangeProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewLongRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewLongRangeProperty(
			esv8.WithLongRangeCoerce(true),
			esv8.WithLongRangeDocValues(false),
			esv8.WithLongRangeIndex(false),
			esv8.WithLongRangeStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewFloatRangeProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewFloatRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewFloatRangeProperty(
			esv8.WithFloatRangeCoerce(true),
			esv8.WithFloatRangeDocValues(false),
			esv8.WithFloatRangeIndex(false),
			esv8.WithFloatRangeStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewDoubleRangeProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDoubleRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDoubleRangeProperty(
			esv8.WithDoubleRangeCoerce(true),
			esv8.WithDoubleRangeDocValues(false),
			esv8.WithDoubleRangeIndex(false),
			esv8.WithDoubleRangeStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewDateRangeProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewDateRangeProperty(
		esv8.WithDateRangeCoerce(true),
		esv8.WithDateRangeDocValues(false),
		esv8.WithDateRangeIndex(false),
		esv8.WithDateRangeStore(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, true, *prop.Coerce)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
}

func TestNewIpRangeProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIpRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewIpRangeProperty(
			esv8.WithIpRangeCoerce(true),
			esv8.WithIpRangeDocValues(false),
			esv8.WithIpRangeIndex(false),
			esv8.WithIpRangeStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Coerce)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewObjectProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewObjectProperty()
	assert.Assert(t, prop != nil)
}

func TestNewNestedProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewNestedProperty(
		esv8.WithNestedEnabled(false),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.Enabled)
}

func TestNewFlattenedProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewFlattenedProperty(
		esv8.WithFlattenedDocValues(false),
		esv8.WithFlattenedIndex(false),
		esv8.WithFlattenedIgnoreAbove(256),
		esv8.WithFlattenedNullValue(""),
		esv8.WithFlattenedEagerGlobalOrdinals(true),
		esv8.WithFlattenedSimilarity("BM25"),
		esv8.WithFlattenedSplitQueriesOnWhitespace(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, 256, *prop.IgnoreAbove)
	assert.Equal(t, "", *prop.NullValue)
	assert.Equal(t, true, *prop.EagerGlobalOrdinals)
	assert.Equal(t, "BM25", *prop.Similarity)
	assert.Equal(t, true, *prop.SplitQueriesOnWhitespace)
}

func TestNewJoinProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewJoinProperty(
		esv8.WithJoinEagerGlobalOrdinals(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, true, *prop.EagerGlobalOrdinals)
}

func TestNewPassthroughObjectProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewPassthroughObjectProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewPassthroughObjectProperty(
			esv8.WithPassthroughObjectProperties(map[string]types.Property{
				"value": esv8.NewKeywordProperty(),
			}),
			esv8.WithPassthroughObjectEnabled(false),
			esv8.WithPassthroughObjectPriority(10),
			esv8.WithPassthroughObjectTimeSeriesDimension(true),
		)
		assert.Assert(t, prop != nil)
		_, ok := prop.Properties["value"]
		assert.Assert(t, ok)
		assert.Equal(t, false, *prop.Enabled)
		assert.Equal(t, 10, *prop.Priority)
		assert.Equal(t, true, *prop.TimeSeriesDimension)
	})
}

func TestNewIpProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewIpProperty(
		esv8.WithIpDocValues(false),
		esv8.WithIpIgnoreMalformed(true),
		esv8.WithIpIndex(false),
		esv8.WithIpStore(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, true, *prop.IgnoreMalformed)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
}

func TestNewBinaryProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewBinaryProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewBinaryProperty(
			esv8.WithBinaryDocValues(false),
			esv8.WithBinaryStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewTokenCountProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewTokenCountProperty(
		esv8.WithTokenCountDocValues(false),
		esv8.WithTokenCountIndex(false),
		esv8.WithTokenCountStore(true),
		esv8.WithTokenCountEnablePositionIncrements(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
	assert.Equal(t, true, *prop.EnablePositionIncrements)
}

func TestNewPercolatorProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewPercolatorProperty(func(p *types.PercolatorProperty) {
		v := 256
		p.IgnoreAbove = &v
	})
	assert.Assert(t, prop != nil)
	assert.Equal(t, 256, *prop.IgnoreAbove)
}

func TestNewHistogramProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewHistogramProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewHistogramProperty(
			esv8.WithHistogramIgnoreMalformed(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.IgnoreMalformed)
	})
}

func TestNewVersionProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewVersionProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewVersionProperty(
			esv8.WithVersionDocValues(false),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
	})
}

func TestNewSparseVectorProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewSparseVectorProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewSparseVectorProperty(
			esv8.WithSparseVectorStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewRankFeaturesProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewRankFeaturesProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewRankFeaturesProperty(
			esv8.WithRankFeaturesPositiveScoreImpact(false),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.PositiveScoreImpact)
	})
}

func TestNewSemanticTextProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewSemanticTextProperty(
		esv8.WithSemanticTextSearchInferenceId("my-search-model"),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, "my-search-model", *prop.SearchInferenceId)
}

func TestNewAggregateMetricDoubleProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewAggregateMetricDoubleProperty(
		esv8.WithAggregateMetricDoubleIgnoreMalformed(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, true, *prop.IgnoreMalformed)
}

func TestNewMurmur3HashProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewMurmur3HashProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewMurmur3HashProperty(
			esv8.WithMurmur3HashDocValues(false),
			esv8.WithMurmur3HashStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewIcuCollationProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv8.NewIcuCollationProperty(
		esv8.WithIcuCollationDocValues(false),
		esv8.WithIcuCollationIndex(false),
		esv8.WithIcuCollationStore(true),
		esv8.WithIcuCollationNullValue(""),
		esv8.WithIcuCollationNorms(false),
		esv8.WithIcuCollationRules("&a<b"),
		esv8.WithIcuCollationVariant("@collation=standard"),
		esv8.WithIcuCollationCaseLevel(true),
		esv8.WithIcuCollationNumeric(true),
		esv8.WithIcuCollationHiraganaQuaternaryMode(true),
		esv8.WithIcuCollationVariableTop("!"),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
	assert.Equal(t, "", *prop.NullValue)
	assert.Equal(t, false, *prop.Norms)
	assert.Equal(t, "&a<b", *prop.Rules)
	assert.Equal(t, "@collation=standard", *prop.Variant)
	assert.Equal(t, true, *prop.CaseLevel)
	assert.Equal(t, true, *prop.Numeric)
	assert.Equal(t, true, *prop.HiraganaQuaternaryMode)
	assert.Equal(t, "!", *prop.VariableTop)
}

func TestNewDynamicProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDynamicProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDynamicProperty(
			esv8.WithDynamicAnalyzer("standard"),
			esv8.WithDynamicSearchAnalyzer("standard"),
			esv8.WithDynamicCoerce(true),
			esv8.WithDynamicDocValues(false),
			esv8.WithDynamicEnabled(false),
			esv8.WithDynamicFormat("strict_date"),
			esv8.WithDynamicIgnoreMalformed(true),
			esv8.WithDynamicIndex(false),
			esv8.WithDynamicStore(true),
			esv8.WithDynamicNorms(false),
			esv8.WithDynamicLocale("en"),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, "standard", *prop.Analyzer)
		assert.Equal(t, "standard", *prop.SearchAnalyzer)
		assert.Equal(t, true, *prop.Coerce)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, false, *prop.Enabled)
		assert.Equal(t, "strict_date", *prop.Format)
		assert.Equal(t, true, *prop.IgnoreMalformed)
		assert.Equal(t, false, *prop.Index)
		assert.Equal(t, true, *prop.Store)
		assert.Equal(t, false, *prop.Norms)
		assert.Equal(t, "en", *prop.Locale)
	})
}
