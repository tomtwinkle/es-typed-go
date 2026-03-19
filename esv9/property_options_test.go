package esv9_test

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/esv9"
)

func TestNewKeywordProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewKeywordProperty(
		esv9.WithKeywordIndex(false),
		esv9.WithKeywordNullValue("N/A"),
		esv9.WithKeywordNorms(false),
		esv9.WithKeywordSimilarity("BM25"),
		esv9.WithKeywordEagerGlobalOrdinals(true),
		esv9.WithKeywordSplitQueriesOnWhitespace(true),
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
	prop := esv9.NewConstantKeywordProperty(func(p *types.ConstantKeywordProperty) {
		v := 256
		p.IgnoreAbove = &v
	})
	assert.Assert(t, prop != nil)
	assert.Equal(t, 256, *prop.IgnoreAbove)
}

func TestNewWildcardProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewWildcardProperty(
		esv9.WithWildcardDocValues(false),
		esv9.WithWildcardNullValue(""),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, "", *prop.NullValue)
}

func TestNewTextProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("with extra options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewTextProperty(
			esv9.WithTextSearchQuoteAnalyzer("standard"),
			esv9.WithTextFielddata(true),
			esv9.WithTextIndex(false),
			esv9.WithTextStore(true),
			esv9.WithTextNorms(false),
			esv9.WithTextSimilarity("BM25"),
			esv9.WithTextIndexPhrases(true),
			esv9.WithTextPositionIncrementGap(100),
			esv9.WithTextFields(map[string]types.Property{
				"value": esv9.NewKeywordProperty(),
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
		esv9.WithTextRawKeyword(256)(p)
		assert.Assert(t, p.Fields != nil)
		_, ok := p.Fields["keyword"]
		assert.Assert(t, ok)
	})
}

func TestNewMatchOnlyTextProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewMatchOnlyTextProperty(func(p *types.MatchOnlyTextProperty) {
		p.CopyTo = []string{"title"}
	})
	assert.Assert(t, prop != nil)
	assert.Equal(t, 1, len(prop.CopyTo))
}

func TestNewCompletionProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewCompletionProperty(
		esv9.WithCompletionSearchAnalyzer("standard"),
		esv9.WithCompletionPreservePositionIncrements(true),
		esv9.WithCompletionPreserveSeparators(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, "standard", *prop.SearchAnalyzer)
	assert.Equal(t, true, *prop.PreservePositionIncrements)
	assert.Equal(t, true, *prop.PreserveSeparators)
}

func TestNewSearchAsYouTypeProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewSearchAsYouTypeProperty(
		esv9.WithSearchAsYouTypeAnalyzer("standard"),
		esv9.WithSearchAsYouTypeSearchAnalyzer("standard"),
		esv9.WithSearchAsYouTypeSearchQuoteAnalyzer("standard"),
		esv9.WithSearchAsYouTypeIndex(false),
		esv9.WithSearchAsYouTypeStore(true),
		esv9.WithSearchAsYouTypeNorms(false),
		esv9.WithSearchAsYouTypeSimilarity("BM25"),
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
	prop := esv9.NewIntegerNumberProperty(
		esv9.WithIntegerNumberDocValues(false),
		esv9.WithIntegerNumberIgnoreMalformed(true),
		esv9.WithIntegerNumberIndex(false),
		esv9.WithIntegerNumberStore(true),
		esv9.WithIntegerNumberNullValue(0),
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
	prop := esv9.NewLongNumberProperty(
		esv9.WithLongNumberCoerce(true),
		esv9.WithLongNumberIgnoreMalformed(true),
		esv9.WithLongNumberIndex(false),
		esv9.WithLongNumberNullValue(nullVal),
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
		prop := esv9.NewShortNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewShortNumberProperty(
			esv9.WithShortNumberCoerce(true),
			esv9.WithShortNumberDocValues(false),
			esv9.WithShortNumberIgnoreMalformed(true),
			esv9.WithShortNumberIndex(false),
			esv9.WithShortNumberStore(true),
			esv9.WithShortNumberNullValue(0),
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
		prop := esv9.NewByteNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		var nullVal byte = 0
		prop := esv9.NewByteNumberProperty(
			esv9.WithByteNumberCoerce(true),
			esv9.WithByteNumberDocValues(false),
			esv9.WithByteNumberIgnoreMalformed(true),
			esv9.WithByteNumberIndex(false),
			esv9.WithByteNumberStore(true),
			esv9.WithByteNumberNullValue(nullVal),
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
	prop := esv9.NewDoubleNumberProperty(
		esv9.WithDoubleNumberDocValues(false),
		esv9.WithDoubleNumberIgnoreMalformed(true),
		esv9.WithDoubleNumberIndex(false),
		esv9.WithDoubleNumberStore(true),
		esv9.WithDoubleNumberNullValue(1.5),
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
		prop := esv9.NewFloatNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		var nullVal float32 = 1.5
		prop := esv9.NewFloatNumberProperty(
			esv9.WithFloatNumberCoerce(true),
			esv9.WithFloatNumberDocValues(false),
			esv9.WithFloatNumberIgnoreMalformed(true),
			esv9.WithFloatNumberIndex(false),
			esv9.WithFloatNumberStore(true),
			esv9.WithFloatNumberNullValue(nullVal),
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
		prop := esv9.NewHalfFloatNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		var nullVal float32 = 1.5
		prop := esv9.NewHalfFloatNumberProperty(
			esv9.WithHalfFloatNumberCoerce(true),
			esv9.WithHalfFloatNumberDocValues(false),
			esv9.WithHalfFloatNumberIgnoreMalformed(true),
			esv9.WithHalfFloatNumberIndex(false),
			esv9.WithHalfFloatNumberStore(true),
			esv9.WithHalfFloatNumberNullValue(nullVal),
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
		prop := esv9.NewUnsignedLongNumberProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		var nullVal uint64 = 0
		prop := esv9.NewUnsignedLongNumberProperty(
			esv9.WithUnsignedLongNumberDocValues(false),
			esv9.WithUnsignedLongNumberIgnoreMalformed(true),
			esv9.WithUnsignedLongNumberIndex(false),
			esv9.WithUnsignedLongNumberStore(true),
			esv9.WithUnsignedLongNumberNullValue(nullVal),
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
	prop := esv9.NewScaledFloatNumberProperty(
		esv9.WithScaledFloatNumberCoerce(true),
		esv9.WithScaledFloatNumberDocValues(false),
		esv9.WithScaledFloatNumberIgnoreMalformed(true),
		esv9.WithScaledFloatNumberIndex(false),
		esv9.WithScaledFloatNumberStore(true),
		esv9.WithScaledFloatNumberNullValue(1.5),
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
	prop := esv9.NewDateProperty(
		esv9.WithDateDocValues(false),
		esv9.WithDateIgnoreMalformed(true),
		esv9.WithDateIndex(false),
		esv9.WithDateStore(true),
		esv9.WithDateLocale("en"),
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
	prop := esv9.NewDateNanosProperty(
		esv9.WithDateNanosDocValues(false),
		esv9.WithDateNanosIgnoreMalformed(true),
		esv9.WithDateNanosIndex(false),
		esv9.WithDateNanosStore(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, true, *prop.IgnoreMalformed)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
}

func TestNewGeoPointProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewGeoPointProperty(
		esv9.WithGeoPointDocValues(false),
		esv9.WithGeoPointIndex(false),
		esv9.WithGeoPointStore(true),
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
		prop := esv9.NewGeoShapeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewGeoShapeProperty(
			esv9.WithGeoShapeCoerce(true),
			esv9.WithGeoShapeIgnoreMalformed(true),
			esv9.WithGeoShapeIgnoreZValue(false),
			esv9.WithGeoShapeDocValues(false),
			esv9.WithGeoShapeIndex(false),
			esv9.WithGeoShapeStore(true),
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
		prop := esv9.NewShapeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewShapeProperty(
			esv9.WithShapeCoerce(true),
			esv9.WithShapeIgnoreMalformed(true),
			esv9.WithShapeIgnoreZValue(false),
			esv9.WithShapeDocValues(false),
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
		prop := esv9.NewPointProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewPointProperty(
			esv9.WithPointIgnoreMalformed(true),
			esv9.WithPointIgnoreZValue(false),
			esv9.WithPointDocValues(false),
			esv9.WithPointStore(true),
			esv9.WithPointNullValue("0,0"),
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
		prop := esv9.NewIntegerRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewIntegerRangeProperty(
			esv9.WithIntegerRangeCoerce(true),
			esv9.WithIntegerRangeDocValues(false),
			esv9.WithIntegerRangeIndex(false),
			esv9.WithIntegerRangeStore(true),
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
		prop := esv9.NewLongRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewLongRangeProperty(
			esv9.WithLongRangeCoerce(true),
			esv9.WithLongRangeDocValues(false),
			esv9.WithLongRangeIndex(false),
			esv9.WithLongRangeStore(true),
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
		prop := esv9.NewFloatRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewFloatRangeProperty(
			esv9.WithFloatRangeCoerce(true),
			esv9.WithFloatRangeDocValues(false),
			esv9.WithFloatRangeIndex(false),
			esv9.WithFloatRangeStore(true),
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
		prop := esv9.NewDoubleRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDoubleRangeProperty(
			esv9.WithDoubleRangeCoerce(true),
			esv9.WithDoubleRangeDocValues(false),
			esv9.WithDoubleRangeIndex(false),
			esv9.WithDoubleRangeStore(true),
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
	prop := esv9.NewDateRangeProperty(
		esv9.WithDateRangeCoerce(true),
		esv9.WithDateRangeDocValues(false),
		esv9.WithDateRangeIndex(false),
		esv9.WithDateRangeStore(true),
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
		prop := esv9.NewIpRangeProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewIpRangeProperty(
			esv9.WithIpRangeCoerce(true),
			esv9.WithIpRangeDocValues(false),
			esv9.WithIpRangeIndex(false),
			esv9.WithIpRangeStore(true),
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
	prop := esv9.NewObjectProperty()
	assert.Assert(t, prop != nil)
}

func TestNewNestedProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewNestedProperty(
		esv9.WithNestedEnabled(false),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.Enabled)
}

func TestNewFlattenedProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewFlattenedProperty(
		esv9.WithFlattenedDocValues(false),
		esv9.WithFlattenedIndex(false),
		esv9.WithFlattenedIgnoreAbove(256),
		esv9.WithFlattenedNullValue(""),
		esv9.WithFlattenedEagerGlobalOrdinals(true),
		esv9.WithFlattenedSimilarity("BM25"),
		esv9.WithFlattenedSplitQueriesOnWhitespace(true),
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
	prop := esv9.NewJoinProperty(
		esv9.WithJoinEagerGlobalOrdinals(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, true, *prop.EagerGlobalOrdinals)
}

func TestNewPassthroughObjectProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewPassthroughObjectProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewPassthroughObjectProperty(
			esv9.WithPassthroughObjectProperties(map[string]types.Property{
				"value": esv9.NewKeywordProperty(),
			}),
			esv9.WithPassthroughObjectEnabled(false),
			esv9.WithPassthroughObjectPriority(10),
			esv9.WithPassthroughObjectTimeSeriesDimension(true),
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
	prop := esv9.NewIpProperty(
		esv9.WithIpDocValues(false),
		esv9.WithIpIgnoreMalformed(true),
		esv9.WithIpIndex(false),
		esv9.WithIpStore(true),
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
		prop := esv9.NewBinaryProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewBinaryProperty(
			esv9.WithBinaryDocValues(false),
			esv9.WithBinaryStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewTokenCountProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewTokenCountProperty(
		esv9.WithTokenCountDocValues(false),
		esv9.WithTokenCountIndex(false),
		esv9.WithTokenCountStore(true),
		esv9.WithTokenCountEnablePositionIncrements(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, false, *prop.DocValues)
	assert.Equal(t, false, *prop.Index)
	assert.Equal(t, true, *prop.Store)
	assert.Equal(t, true, *prop.EnablePositionIncrements)
}

func TestNewPercolatorProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewPercolatorProperty(func(p *types.PercolatorProperty) {
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
		prop := esv9.NewHistogramProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewHistogramProperty(
			esv9.WithHistogramIgnoreMalformed(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.IgnoreMalformed)
	})
}

func TestNewVersionProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewVersionProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewVersionProperty(
			esv9.WithVersionDocValues(false),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
	})
}

func TestNewSparseVectorProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewSparseVectorProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewSparseVectorProperty(
			esv9.WithSparseVectorStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewRankFeaturesProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewRankFeaturesProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewRankFeaturesProperty(
			esv9.WithRankFeaturesPositiveScoreImpact(false),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.PositiveScoreImpact)
	})
}

func TestNewSemanticTextProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewSemanticTextProperty(
		esv9.WithSemanticTextSearchInferenceId("my-search-model"),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, "my-search-model", *prop.SearchInferenceId)
}

func TestNewAggregateMetricDoubleProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewAggregateMetricDoubleProperty(
		esv9.WithAggregateMetricDoubleIgnoreMalformed(true),
	)
	assert.Assert(t, prop != nil)
	assert.Equal(t, true, *prop.IgnoreMalformed)
}

func TestNewMurmur3HashProperty_Options(t *testing.T) {
	t.Parallel()

	t.Run("no options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewMurmur3HashProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewMurmur3HashProperty(
			esv9.WithMurmur3HashDocValues(false),
			esv9.WithMurmur3HashStore(true),
		)
		assert.Assert(t, prop != nil)
		assert.Equal(t, false, *prop.DocValues)
		assert.Equal(t, true, *prop.Store)
	})
}

func TestNewIcuCollationProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewIcuCollationProperty(
		esv9.WithIcuCollationDocValues(false),
		esv9.WithIcuCollationIndex(false),
		esv9.WithIcuCollationStore(true),
		esv9.WithIcuCollationNullValue(""),
		esv9.WithIcuCollationNorms(false),
		esv9.WithIcuCollationRules("&a<b"),
		esv9.WithIcuCollationVariant("@collation=standard"),
		esv9.WithIcuCollationCaseLevel(true),
		esv9.WithIcuCollationNumeric(true),
		esv9.WithIcuCollationHiraganaQuaternaryMode(true),
		esv9.WithIcuCollationVariableTop("!"),
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
		prop := esv9.NewDynamicProperty()
		assert.Assert(t, prop != nil)
	})

	t.Run("with options", func(t *testing.T) {
		t.Parallel()
		prop := esv9.NewDynamicProperty(
			esv9.WithDynamicAnalyzer("standard"),
			esv9.WithDynamicSearchAnalyzer("standard"),
			esv9.WithDynamicCoerce(true),
			esv9.WithDynamicDocValues(false),
			esv9.WithDynamicEnabled(false),
			esv9.WithDynamicFormat("strict_date"),
			esv9.WithDynamicIgnoreMalformed(true),
			esv9.WithDynamicIndex(false),
			esv9.WithDynamicStore(true),
			esv9.WithDynamicNorms(false),
			esv9.WithDynamicLocale("en"),
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

func TestNewExponentialHistogramProperty_Options(t *testing.T) {
	t.Parallel()
	prop := esv9.NewExponentialHistogramProperty(func(p *types.ExponentialHistogramProperty) {
		v := 256
		p.IgnoreAbove = &v
	})
	assert.Assert(t, prop != nil)
	assert.Equal(t, 256, *prop.IgnoreAbove)
}
