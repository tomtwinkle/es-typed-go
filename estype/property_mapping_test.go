package estype_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

// TestMappingFieldTypeName verifies that TypeName() returns the correct ES type
// string from both FieldType properties and typed property values.
func TestMappingFieldTypeName(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		field estype.MappingField
		want  string
	}{
		"field_type_keyword":  {field: estype.MappingField{Path: "status", Property: estype.FieldType("keyword")}, want: "keyword"},
		"field_type_text":     {field: estype.MappingField{Path: "title", Property: estype.FieldType("text")}, want: "text"},
		"field_type_integer":  {field: estype.MappingField{Path: "price", Property: estype.FieldType("integer")}, want: "integer"},
		"field_type_nested":   {field: estype.MappingField{Path: "items", Property: estype.FieldType("nested")}, want: "nested"},
		"text_property":       {field: estype.MappingField{Path: "title", Property: estype.NewTextProperty()}, want: "text"},
		"keyword_property":    {field: estype.MappingField{Path: "status", Property: estype.NewKeywordProperty()}, want: "keyword"},
		"nil_property":        {field: estype.MappingField{Path: "status"}, want: ""},
		"field_type_empty":    {field: estype.MappingField{Path: "status", Property: estype.FieldType("")}, want: ""},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tt.field.TypeName()
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestAllPropertyESTypeName verifies that every concrete property type returns
// the correct ES type name string from ESTypeName().
func TestAllPropertyESTypeName(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		property estype.MappingProperty
		want     string
	}{
		"text":                    {property: estype.NewTextProperty(), want: "text"},
		"keyword":                 {property: estype.NewKeywordProperty(), want: "keyword"},
		"constant_keyword":        {property: estype.NewConstantKeywordProperty(), want: "constant_keyword"},
		"counted_keyword":         {property: estype.NewCountedKeywordProperty(), want: "counted_keyword"},
		"wildcard":                {property: estype.NewWildcardProperty(), want: "wildcard"},
		"match_only_text":         {property: estype.NewMatchOnlyTextProperty(), want: "match_only_text"},
		"completion":              {property: estype.NewCompletionProperty(), want: "completion"},
		"search_as_you_type":      {property: estype.NewSearchAsYouTypeProperty(), want: "search_as_you_type"},
		"boolean":                 {property: estype.NewBooleanProperty(), want: "boolean"},
		"integer":                 {property: estype.NewIntegerNumberProperty(), want: "integer"},
		"long":                    {property: estype.NewLongNumberProperty(), want: "long"},
		"short":                   {property: estype.NewShortNumberProperty(), want: "short"},
		"byte":                    {property: estype.NewByteNumberProperty(), want: "byte"},
		"double":                  {property: estype.NewDoubleNumberProperty(), want: "double"},
		"float":                   {property: estype.NewFloatNumberProperty(), want: "float"},
		"half_float":              {property: estype.NewHalfFloatNumberProperty(), want: "half_float"},
		"unsigned_long":           {property: estype.NewUnsignedLongNumberProperty(), want: "unsigned_long"},
		"scaled_float":            {property: estype.NewScaledFloatNumberProperty(), want: "scaled_float"},
		"date":                    {property: estype.NewDateProperty(), want: "date"},
		"date_nanos":              {property: estype.NewDateNanosProperty(), want: "date_nanos"},
		"geo_point":               {property: estype.NewGeoPointProperty(), want: "geo_point"},
		"geo_shape":               {property: estype.NewGeoShapeProperty(), want: "geo_shape"},
		"shape":                   {property: estype.NewShapeProperty(), want: "shape"},
		"point":                   {property: estype.NewPointProperty(), want: "point"},
		"integer_range":           {property: estype.NewIntegerRangeProperty(), want: "integer_range"},
		"long_range":              {property: estype.NewLongRangeProperty(), want: "long_range"},
		"float_range":             {property: estype.NewFloatRangeProperty(), want: "float_range"},
		"double_range":            {property: estype.NewDoubleRangeProperty(), want: "double_range"},
		"date_range":              {property: estype.NewDateRangeProperty(), want: "date_range"},
		"ip_range":                {property: estype.NewIpRangeProperty(), want: "ip_range"},
		"object":                  {property: estype.NewObjectProperty(), want: "object"},
		"nested":                  {property: estype.NewNestedProperty(), want: "nested"},
		"flattened":               {property: estype.NewFlattenedProperty(), want: "flattened"},
		"join":                    {property: estype.NewJoinProperty(), want: "join"},
		"passthrough":             {property: estype.NewPassthroughObjectProperty(), want: "passthrough"},
		"ip":                      {property: estype.NewIpProperty(), want: "ip"},
		"binary":                  {property: estype.NewBinaryProperty(), want: "binary"},
		"token_count":             {property: estype.NewTokenCountProperty(), want: "token_count"},
		"percolator":              {property: estype.NewPercolatorProperty(), want: "percolator"},
		"alias":                   {property: estype.NewFieldAliasProperty(), want: "alias"},
		"histogram":               {property: estype.NewHistogramProperty(), want: "histogram"},
		"version":                 {property: estype.NewVersionProperty(), want: "version"},
		"dense_vector":            {property: estype.NewDenseVectorProperty(), want: "dense_vector"},
		"sparse_vector":           {property: estype.NewSparseVectorProperty(), want: "sparse_vector"},
		"rank_feature":            {property: estype.NewRankFeatureProperty(), want: "rank_feature"},
		"rank_features":           {property: estype.NewRankFeaturesProperty(), want: "rank_features"},
		"rank_vectors":            {property: estype.NewRankVectorProperty(), want: "rank_vectors"},
		"semantic_text":           {property: estype.NewSemanticTextProperty(), want: "semantic_text"},
		"aggregate_metric_double": {property: estype.NewAggregateMetricDoubleProperty(), want: "aggregate_metric_double"},
		"murmur3":                 {property: estype.NewMurmur3HashProperty(), want: "murmur3"},
		"icu_collation_keyword":   {property: estype.NewIcuCollationProperty(), want: "icu_collation_keyword"},
		"dynamic_type":            {property: estype.NewDynamicProperty(), want: "{dynamic_type}"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.property.ESTypeName())
		})
	}
}

// TestAllPropertyMappingPropertyInterface verifies that every property type
// can be assigned to a MappingField.Property (MappingProperty interface).
func TestAllPropertyMappingPropertyInterface(t *testing.T) {
	t.Parallel()
	properties := map[string]estype.MappingProperty{
		"text":                    estype.NewTextProperty(),
		"keyword":                 estype.NewKeywordProperty(),
		"constant_keyword":        estype.NewConstantKeywordProperty(),
		"counted_keyword":         estype.NewCountedKeywordProperty(),
		"wildcard":                estype.NewWildcardProperty(),
		"match_only_text":         estype.NewMatchOnlyTextProperty(),
		"completion":              estype.NewCompletionProperty(),
		"search_as_you_type":      estype.NewSearchAsYouTypeProperty(),
		"boolean":                 estype.NewBooleanProperty(),
		"integer":                 estype.NewIntegerNumberProperty(),
		"long":                    estype.NewLongNumberProperty(),
		"short":                   estype.NewShortNumberProperty(),
		"byte":                    estype.NewByteNumberProperty(),
		"double":                  estype.NewDoubleNumberProperty(),
		"float":                   estype.NewFloatNumberProperty(),
		"half_float":              estype.NewHalfFloatNumberProperty(),
		"unsigned_long":           estype.NewUnsignedLongNumberProperty(),
		"scaled_float":            estype.NewScaledFloatNumberProperty(),
		"date":                    estype.NewDateProperty(),
		"date_nanos":              estype.NewDateNanosProperty(),
		"geo_point":               estype.NewGeoPointProperty(),
		"geo_shape":               estype.NewGeoShapeProperty(),
		"shape":                   estype.NewShapeProperty(),
		"point":                   estype.NewPointProperty(),
		"integer_range":           estype.NewIntegerRangeProperty(),
		"long_range":              estype.NewLongRangeProperty(),
		"float_range":             estype.NewFloatRangeProperty(),
		"double_range":            estype.NewDoubleRangeProperty(),
		"date_range":              estype.NewDateRangeProperty(),
		"ip_range":                estype.NewIpRangeProperty(),
		"object":                  estype.NewObjectProperty(),
		"nested":                  estype.NewNestedProperty(),
		"flattened":               estype.NewFlattenedProperty(),
		"join":                    estype.NewJoinProperty(),
		"passthrough":             estype.NewPassthroughObjectProperty(),
		"ip":                      estype.NewIpProperty(),
		"binary":                  estype.NewBinaryProperty(),
		"token_count":             estype.NewTokenCountProperty(),
		"percolator":              estype.NewPercolatorProperty(),
		"alias":                   estype.NewFieldAliasProperty(),
		"histogram":               estype.NewHistogramProperty(),
		"version":                 estype.NewVersionProperty(),
		"dense_vector":            estype.NewDenseVectorProperty(),
		"sparse_vector":           estype.NewSparseVectorProperty(),
		"rank_feature":            estype.NewRankFeatureProperty(),
		"rank_features":           estype.NewRankFeaturesProperty(),
		"rank_vectors":            estype.NewRankVectorProperty(),
		"semantic_text":           estype.NewSemanticTextProperty(),
		"aggregate_metric_double": estype.NewAggregateMetricDoubleProperty(),
		"murmur3":                 estype.NewMurmur3HashProperty(),
		"icu_collation_keyword":   estype.NewIcuCollationProperty(),
		"dynamic_type":            estype.NewDynamicProperty(),
	}
	for name, prop := range properties {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f := estype.MappingField{Path: "value", Property: prop}
			assert.Assert(t, f.TypeName() != "" || name == "")
		})
	}
}

// TestAnalyzer verifies the Analyzer type and its String() method.
func TestAnalyzer(t *testing.T) {
	t.Parallel()
	a := estype.Analyzer("my_analyzer")
	assert.Equal(t, "my_analyzer", a.String())
	assert.Equal(t, estype.Analyzer("my_analyzer"), a)
}

// TestTextProperty verifies that NewTextProperty and its options produce a
// correctly configured TextProperty.
func TestTextProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewTextProperty()
		assert.Equal(t, "text", p.ESTypeName())
		assert.Assert(t, p.SearchAnalyzer == nil)
		assert.Assert(t, p.IndexAnalyzer == nil)
		assert.Assert(t, p.Fields == nil)
	})

	t.Run("with_search_analyzer", func(t *testing.T) {
		t.Parallel()
		a := estype.Analyzer("my_search")
		p := estype.NewTextProperty(estype.WithSearchAnalyzer(a))
		assert.Assert(t, p.SearchAnalyzer != nil)
		assert.Equal(t, a, *p.SearchAnalyzer)
	})

	t.Run("with_index_analyzer", func(t *testing.T) {
		t.Parallel()
		a := estype.Analyzer("my_index")
		p := estype.NewTextProperty(estype.WithIndexAnalyzer(a))
		assert.Assert(t, p.IndexAnalyzer != nil)
		assert.Equal(t, a, *p.IndexAnalyzer)
	})

	t.Run("with_field", func(t *testing.T) {
		t.Parallel()
		kw := estype.NewKeywordProperty(estype.WithIgnoreAbove())
		p := estype.NewTextProperty(estype.WithField("keyword", kw))
		assert.Assert(t, p.Fields != nil)
		assert.Assert(t, len(p.Fields) == 1)
		sub, ok := p.Fields["keyword"]
		assert.Assert(t, ok, "expected sub-field 'keyword'")
		subKw, ok := sub.(estype.KeywordProperty)
		assert.Assert(t, ok, "expected KeywordProperty sub-field")
		assert.Equal(t, "keyword", subKw.ESTypeName())
	})

	t.Run("multiple_options", func(t *testing.T) {
		t.Parallel()
		sa := estype.Analyzer("search_a")
		ia := estype.Analyzer("index_a")
		p := estype.NewTextProperty(
			estype.WithSearchAnalyzer(sa),
			estype.WithIndexAnalyzer(ia),
			estype.WithField("keyword", estype.NewKeywordProperty()),
		)
		assert.Assert(t, p.SearchAnalyzer != nil)
		assert.Equal(t, sa, *p.SearchAnalyzer)
		assert.Assert(t, p.IndexAnalyzer != nil)
		assert.Equal(t, ia, *p.IndexAnalyzer)
		assert.Assert(t, len(p.Fields) == 1)
	})
}

// TestKeywordProperty verifies that NewKeywordProperty and its options produce
// a correctly configured KeywordProperty.
func TestKeywordProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewKeywordProperty()
		assert.Equal(t, "keyword", p.ESTypeName())
		assert.Assert(t, p.IgnoreAbove == nil)
	})

	t.Run("with_ignore_above_default", func(t *testing.T) {
		t.Parallel()
		p := estype.NewKeywordProperty(estype.WithIgnoreAbove())
		assert.Assert(t, p.IgnoreAbove != nil)
		assert.Equal(t, 256, *p.IgnoreAbove)
	})

	t.Run("with_ignore_above_custom", func(t *testing.T) {
		t.Parallel()
		p := estype.NewKeywordProperty(estype.WithIgnoreAbove(512))
		assert.Assert(t, p.IgnoreAbove != nil)
		assert.Equal(t, 512, *p.IgnoreAbove)
	})
}

// TestMappingFieldPath verifies that the Path field of MappingField is of type
// string and holds the dot-separated field path.
func TestMappingFieldPath(t *testing.T) {
	t.Parallel()
	f := estype.MappingField{Path: "items.color", Property: estype.FieldType("keyword")}
	assert.Equal(t, "items.color", f.Path)
}

// ---------------------------------------------------------------------------
// Constant Keyword
// ---------------------------------------------------------------------------

// TestConstantKeywordProperty verifies that NewConstantKeywordProperty and its
// options produce a correctly configured ConstantKeywordProperty.
func TestConstantKeywordProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewConstantKeywordProperty()
		assert.Equal(t, "constant_keyword", p.ESTypeName())
		assert.Assert(t, p.Value == nil)
	})

	t.Run("with_value", func(t *testing.T) {
		t.Parallel()
		p := estype.NewConstantKeywordProperty(estype.WithConstantKeywordValue("active"))
		assert.Assert(t, p.Value != nil)
		assert.Equal(t, "active", *p.Value)
	})
}

// ---------------------------------------------------------------------------
// Counted Keyword
// ---------------------------------------------------------------------------

// TestCountedKeywordProperty verifies that NewCountedKeywordProperty returns
// a correctly typed CountedKeywordProperty.
func TestCountedKeywordProperty(t *testing.T) {
	t.Parallel()
	p := estype.NewCountedKeywordProperty()
	assert.Equal(t, "counted_keyword", p.ESTypeName())
}

// ---------------------------------------------------------------------------
// Wildcard
// ---------------------------------------------------------------------------

// TestWildcardProperty verifies that NewWildcardProperty and its options
// produce a correctly configured WildcardProperty.
func TestWildcardProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewWildcardProperty()
		assert.Equal(t, "wildcard", p.ESTypeName())
		assert.Assert(t, p.IgnoreAbove == nil)
	})

	t.Run("with_ignore_above", func(t *testing.T) {
		t.Parallel()
		p := estype.NewWildcardProperty(estype.WithWildcardIgnoreAbove(512))
		assert.Assert(t, p.IgnoreAbove != nil)
		assert.Equal(t, 512, *p.IgnoreAbove)
	})
}

// ---------------------------------------------------------------------------
// Match Only Text
// ---------------------------------------------------------------------------

// TestMatchOnlyTextProperty verifies that NewMatchOnlyTextProperty and its
// options produce a correctly configured MatchOnlyTextProperty.
func TestMatchOnlyTextProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewMatchOnlyTextProperty()
		assert.Equal(t, "match_only_text", p.ESTypeName())
		assert.Assert(t, p.Fields == nil)
	})

	t.Run("with_field", func(t *testing.T) {
		t.Parallel()
		p := estype.NewMatchOnlyTextProperty(
			estype.WithMatchOnlyTextField("keyword", estype.NewKeywordProperty()),
		)
		assert.Assert(t, p.Fields != nil)
		assert.Equal(t, 1, len(p.Fields))
		_, ok := p.Fields["keyword"]
		assert.Assert(t, ok)
	})

	t.Run("with_multiple_fields", func(t *testing.T) {
		t.Parallel()
		p := estype.NewMatchOnlyTextProperty(
			estype.WithMatchOnlyTextField("keyword", estype.NewKeywordProperty()),
			estype.WithMatchOnlyTextField("raw", estype.NewKeywordProperty()),
		)
		assert.Equal(t, 2, len(p.Fields))
	})
}

// ---------------------------------------------------------------------------
// Completion
// ---------------------------------------------------------------------------

// TestCompletionProperty verifies that NewCompletionProperty and its options
// produce a correctly configured CompletionProperty.
func TestCompletionProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewCompletionProperty()
		assert.Equal(t, "completion", p.ESTypeName())
		assert.Assert(t, p.Analyzer == nil)
		assert.Assert(t, p.SearchAnalyzer == nil)
	})

	t.Run("with_analyzer", func(t *testing.T) {
		t.Parallel()
		a := estype.Analyzer("my_analyzer")
		p := estype.NewCompletionProperty(estype.WithCompletionAnalyzer(a))
		assert.Assert(t, p.Analyzer != nil)
		assert.Equal(t, a, *p.Analyzer)
	})

	t.Run("with_search_analyzer", func(t *testing.T) {
		t.Parallel()
		a := estype.Analyzer("my_search")
		p := estype.NewCompletionProperty(estype.WithCompletionSearchAnalyzer(a))
		assert.Assert(t, p.SearchAnalyzer != nil)
		assert.Equal(t, a, *p.SearchAnalyzer)
	})

	t.Run("with_all_options", func(t *testing.T) {
		t.Parallel()
		ia := estype.Analyzer("idx_analyzer")
		sa := estype.Analyzer("srch_analyzer")
		p := estype.NewCompletionProperty(
			estype.WithCompletionAnalyzer(ia),
			estype.WithCompletionSearchAnalyzer(sa),
		)
		assert.Assert(t, p.Analyzer != nil)
		assert.Equal(t, ia, *p.Analyzer)
		assert.Assert(t, p.SearchAnalyzer != nil)
		assert.Equal(t, sa, *p.SearchAnalyzer)
	})
}

// ---------------------------------------------------------------------------
// Search As You Type
// ---------------------------------------------------------------------------

// TestSearchAsYouTypeProperty verifies that NewSearchAsYouTypeProperty and its
// options produce a correctly configured SearchAsYouTypeProperty.
func TestSearchAsYouTypeProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewSearchAsYouTypeProperty()
		assert.Equal(t, "search_as_you_type", p.ESTypeName())
		assert.Assert(t, p.Analyzer == nil)
		assert.Assert(t, p.SearchAnalyzer == nil)
		assert.Assert(t, p.MaxShingleSize == nil)
	})

	t.Run("with_analyzer", func(t *testing.T) {
		t.Parallel()
		a := estype.Analyzer("my_analyzer")
		p := estype.NewSearchAsYouTypeProperty(estype.WithSearchAsYouTypeAnalyzer(a))
		assert.Assert(t, p.Analyzer != nil)
		assert.Equal(t, a, *p.Analyzer)
	})

	t.Run("with_search_analyzer", func(t *testing.T) {
		t.Parallel()
		a := estype.Analyzer("my_search")
		p := estype.NewSearchAsYouTypeProperty(estype.WithSearchAsYouTypeSearchAnalyzer(a))
		assert.Assert(t, p.SearchAnalyzer != nil)
		assert.Equal(t, a, *p.SearchAnalyzer)
	})

	t.Run("with_max_shingle_size", func(t *testing.T) {
		t.Parallel()
		p := estype.NewSearchAsYouTypeProperty(estype.WithSearchAsYouTypeMaxShingleSize(3))
		assert.Assert(t, p.MaxShingleSize != nil)
		assert.Equal(t, 3, *p.MaxShingleSize)
	})

	t.Run("with_all_options", func(t *testing.T) {
		t.Parallel()
		a := estype.Analyzer("my_analyzer")
		sa := estype.Analyzer("my_search")
		p := estype.NewSearchAsYouTypeProperty(
			estype.WithSearchAsYouTypeAnalyzer(a),
			estype.WithSearchAsYouTypeSearchAnalyzer(sa),
			estype.WithSearchAsYouTypeMaxShingleSize(4),
		)
		assert.Assert(t, p.Analyzer != nil)
		assert.Equal(t, a, *p.Analyzer)
		assert.Assert(t, p.SearchAnalyzer != nil)
		assert.Equal(t, sa, *p.SearchAnalyzer)
		assert.Equal(t, 4, *p.MaxShingleSize)
	})
}

// ---------------------------------------------------------------------------
// Boolean
// ---------------------------------------------------------------------------

// TestBooleanProperty verifies that NewBooleanProperty and its options produce
// a correctly configured BooleanProperty.
func TestBooleanProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewBooleanProperty()
		assert.Equal(t, "boolean", p.ESTypeName())
		assert.Assert(t, p.NullValue == nil)
	})

	t.Run("with_null_value_true", func(t *testing.T) {
		t.Parallel()
		p := estype.NewBooleanProperty(estype.WithBooleanNullValue(true))
		assert.Assert(t, p.NullValue != nil)
		assert.Assert(t, *p.NullValue == true)
	})

	t.Run("with_null_value_false", func(t *testing.T) {
		t.Parallel()
		p := estype.NewBooleanProperty(estype.WithBooleanNullValue(false))
		assert.Assert(t, p.NullValue != nil)
		assert.Assert(t, *p.NullValue == false)
	})
}

// ---------------------------------------------------------------------------
// Numeric properties
// ---------------------------------------------------------------------------

// TestNumericProperties verifies that all numeric property types return the
// correct ES type name and are constructable.
func TestNumericProperties(t *testing.T) {
	t.Parallel()

	t.Run("integer", func(t *testing.T) {
		t.Parallel()
		p := estype.NewIntegerNumberProperty()
		assert.Equal(t, "integer", p.ESTypeName())
	})

	t.Run("long", func(t *testing.T) {
		t.Parallel()
		p := estype.NewLongNumberProperty()
		assert.Equal(t, "long", p.ESTypeName())
	})

	t.Run("short", func(t *testing.T) {
		t.Parallel()
		p := estype.NewShortNumberProperty()
		assert.Equal(t, "short", p.ESTypeName())
	})

	t.Run("byte", func(t *testing.T) {
		t.Parallel()
		p := estype.NewByteNumberProperty()
		assert.Equal(t, "byte", p.ESTypeName())
	})

	t.Run("double", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDoubleNumberProperty()
		assert.Equal(t, "double", p.ESTypeName())
	})

	t.Run("float", func(t *testing.T) {
		t.Parallel()
		p := estype.NewFloatNumberProperty()
		assert.Equal(t, "float", p.ESTypeName())
	})

	t.Run("half_float", func(t *testing.T) {
		t.Parallel()
		p := estype.NewHalfFloatNumberProperty()
		assert.Equal(t, "half_float", p.ESTypeName())
	})

	t.Run("unsigned_long", func(t *testing.T) {
		t.Parallel()
		p := estype.NewUnsignedLongNumberProperty()
		assert.Equal(t, "unsigned_long", p.ESTypeName())
	})
}

// TestScaledFloatNumberProperty verifies that NewScaledFloatNumberProperty and
// its options produce a correctly configured ScaledFloatNumberProperty.
func TestScaledFloatNumberProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewScaledFloatNumberProperty()
		assert.Equal(t, "scaled_float", p.ESTypeName())
		assert.Assert(t, p.ScalingFactor == nil)
	})

	t.Run("with_scaling_factor", func(t *testing.T) {
		t.Parallel()
		p := estype.NewScaledFloatNumberProperty(estype.WithScalingFactor(100))
		assert.Assert(t, p.ScalingFactor != nil)
		assert.Equal(t, float64(100), *p.ScalingFactor)
	})

	t.Run("with_scaling_factor_decimal", func(t *testing.T) {
		t.Parallel()
		p := estype.NewScaledFloatNumberProperty(estype.WithScalingFactor(10.5))
		assert.Assert(t, p.ScalingFactor != nil)
		assert.Equal(t, 10.5, *p.ScalingFactor)
	})
}

// ---------------------------------------------------------------------------
// Date
// ---------------------------------------------------------------------------

// TestDateProperty verifies that NewDateProperty and its options produce a
// correctly configured DateProperty.
func TestDateProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDateProperty()
		assert.Equal(t, "date", p.ESTypeName())
		assert.Assert(t, p.Format == nil)
	})

	t.Run("with_format", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDateProperty(estype.WithDateFormat("yyyy-MM-dd"))
		assert.Assert(t, p.Format != nil)
		assert.Equal(t, "yyyy-MM-dd", *p.Format)
	})

	t.Run("with_multiple_formats", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDateProperty(estype.WithDateFormat("yyyy-MM-dd||epoch_millis"))
		assert.Assert(t, p.Format != nil)
		assert.Equal(t, "yyyy-MM-dd||epoch_millis", *p.Format)
	})
}

// TestDateNanosProperty verifies that NewDateNanosProperty and its options
// produce a correctly configured DateNanosProperty.
func TestDateNanosProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDateNanosProperty()
		assert.Equal(t, "date_nanos", p.ESTypeName())
		assert.Assert(t, p.Format == nil)
	})

	t.Run("with_format", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDateNanosProperty(estype.WithDateNanosFormat("strict_date_optional_time_nanos"))
		assert.Assert(t, p.Format != nil)
		assert.Equal(t, "strict_date_optional_time_nanos", *p.Format)
	})
}

// ---------------------------------------------------------------------------
// Geo
// ---------------------------------------------------------------------------

// TestGeoProperties verifies that all geo property types return the correct
// ES type name and are constructable.
func TestGeoProperties(t *testing.T) {
	t.Parallel()

	t.Run("geo_point", func(t *testing.T) {
		t.Parallel()
		p := estype.NewGeoPointProperty()
		assert.Equal(t, "geo_point", p.ESTypeName())
	})

	t.Run("geo_shape", func(t *testing.T) {
		t.Parallel()
		p := estype.NewGeoShapeProperty()
		assert.Equal(t, "geo_shape", p.ESTypeName())
	})

	t.Run("shape", func(t *testing.T) {
		t.Parallel()
		p := estype.NewShapeProperty()
		assert.Equal(t, "shape", p.ESTypeName())
	})

	t.Run("point", func(t *testing.T) {
		t.Parallel()
		p := estype.NewPointProperty()
		assert.Equal(t, "point", p.ESTypeName())
	})
}

// ---------------------------------------------------------------------------
// Range
// ---------------------------------------------------------------------------

// TestRangeProperties verifies that all range property types return the correct
// ES type name and are constructable.
func TestRangeProperties(t *testing.T) {
	t.Parallel()

	t.Run("integer_range", func(t *testing.T) {
		t.Parallel()
		p := estype.NewIntegerRangeProperty()
		assert.Equal(t, "integer_range", p.ESTypeName())
	})

	t.Run("long_range", func(t *testing.T) {
		t.Parallel()
		p := estype.NewLongRangeProperty()
		assert.Equal(t, "long_range", p.ESTypeName())
	})

	t.Run("float_range", func(t *testing.T) {
		t.Parallel()
		p := estype.NewFloatRangeProperty()
		assert.Equal(t, "float_range", p.ESTypeName())
	})

	t.Run("double_range", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDoubleRangeProperty()
		assert.Equal(t, "double_range", p.ESTypeName())
	})

	t.Run("ip_range", func(t *testing.T) {
		t.Parallel()
		p := estype.NewIpRangeProperty()
		assert.Equal(t, "ip_range", p.ESTypeName())
	})
}

// TestDateRangeProperty verifies that NewDateRangeProperty and its options
// produce a correctly configured DateRangeProperty.
func TestDateRangeProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDateRangeProperty()
		assert.Equal(t, "date_range", p.ESTypeName())
		assert.Assert(t, p.Format == nil)
	})

	t.Run("with_format", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDateRangeProperty(estype.WithDateRangeFormat("yyyy-MM-dd"))
		assert.Assert(t, p.Format != nil)
		assert.Equal(t, "yyyy-MM-dd", *p.Format)
	})
}

// ---------------------------------------------------------------------------
// Object / Nested
// ---------------------------------------------------------------------------

// TestObjectProperty verifies that NewObjectProperty and its options produce a
// correctly configured ObjectProperty.
func TestObjectProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewObjectProperty()
		assert.Equal(t, "object", p.ESTypeName())
		assert.Assert(t, p.Enabled == nil)
		assert.Assert(t, p.Properties == nil)
	})

	t.Run("with_enabled_true", func(t *testing.T) {
		t.Parallel()
		p := estype.NewObjectProperty(estype.WithObjectEnabled(true))
		assert.Assert(t, p.Enabled != nil)
		assert.Assert(t, *p.Enabled == true)
	})

	t.Run("with_enabled_false", func(t *testing.T) {
		t.Parallel()
		p := estype.NewObjectProperty(estype.WithObjectEnabled(false))
		assert.Assert(t, p.Enabled != nil)
		assert.Assert(t, *p.Enabled == false)
	})

	t.Run("with_property", func(t *testing.T) {
		t.Parallel()
		p := estype.NewObjectProperty(
			estype.WithObjectProperty("status", estype.NewKeywordProperty()),
		)
		assert.Assert(t, p.Properties != nil)
		assert.Equal(t, 1, len(p.Properties))
		sub, ok := p.Properties["status"]
		assert.Assert(t, ok)
		assert.Equal(t, "keyword", sub.ESTypeName())
	})

	t.Run("with_multiple_properties", func(t *testing.T) {
		t.Parallel()
		p := estype.NewObjectProperty(
			estype.WithObjectEnabled(true),
			estype.WithObjectProperty("status", estype.NewKeywordProperty()),
			estype.WithObjectProperty("title", estype.NewTextProperty()),
		)
		assert.Assert(t, *p.Enabled == true)
		assert.Equal(t, 2, len(p.Properties))
	})
}

// TestNestedProperty verifies that NewNestedProperty and its options produce a
// correctly configured NestedProperty.
func TestNestedProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewNestedProperty()
		assert.Equal(t, "nested", p.ESTypeName())
		assert.Assert(t, p.Properties == nil)
	})

	t.Run("with_property", func(t *testing.T) {
		t.Parallel()
		p := estype.NewNestedProperty(
			estype.WithNestedProperty("name", estype.NewKeywordProperty()),
		)
		assert.Assert(t, p.Properties != nil)
		assert.Equal(t, 1, len(p.Properties))
		sub, ok := p.Properties["name"]
		assert.Assert(t, ok)
		assert.Equal(t, "keyword", sub.ESTypeName())
	})

	t.Run("with_multiple_properties", func(t *testing.T) {
		t.Parallel()
		p := estype.NewNestedProperty(
			estype.WithNestedProperty("name", estype.NewKeywordProperty()),
			estype.WithNestedProperty("value", estype.NewIntegerNumberProperty()),
		)
		assert.Equal(t, 2, len(p.Properties))
	})
}

// TestFlattenedProperty verifies that NewFlattenedProperty and its options
// produce a correctly configured FlattenedProperty.
func TestFlattenedProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewFlattenedProperty()
		assert.Equal(t, "flattened", p.ESTypeName())
		assert.Assert(t, p.DepthLimit == nil)
		assert.Assert(t, p.IgnoreAbove == nil)
	})

	t.Run("with_depth_limit", func(t *testing.T) {
		t.Parallel()
		p := estype.NewFlattenedProperty(estype.WithFlattenedDepthLimit(20))
		assert.Assert(t, p.DepthLimit != nil)
		assert.Equal(t, 20, *p.DepthLimit)
	})

	t.Run("with_ignore_above", func(t *testing.T) {
		t.Parallel()
		p := estype.NewFlattenedProperty(estype.WithFlattenedIgnoreAbove(256))
		assert.Assert(t, p.IgnoreAbove != nil)
		assert.Equal(t, 256, *p.IgnoreAbove)
	})

	t.Run("with_all_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewFlattenedProperty(
			estype.WithFlattenedDepthLimit(10),
			estype.WithFlattenedIgnoreAbove(512),
		)
		assert.Equal(t, 10, *p.DepthLimit)
		assert.Equal(t, 512, *p.IgnoreAbove)
	})
}

// TestJoinProperty verifies that NewJoinProperty and its options produce a
// correctly configured JoinProperty.
func TestJoinProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewJoinProperty()
		assert.Equal(t, "join", p.ESTypeName())
		assert.Assert(t, p.Relations == nil)
	})

	t.Run("with_single_relation", func(t *testing.T) {
		t.Parallel()
		p := estype.NewJoinProperty(
			estype.WithJoinRelation("category", "items"),
		)
		assert.Assert(t, p.Relations != nil)
		assert.Equal(t, 1, len(p.Relations))
		children, ok := p.Relations["category"]
		assert.Assert(t, ok)
		assert.Equal(t, 1, len(children))
		assert.Equal(t, "items", children[0])
	})

	t.Run("with_multiple_children", func(t *testing.T) {
		t.Parallel()
		p := estype.NewJoinProperty(
			estype.WithJoinRelation("category", "items", "tags"),
		)
		children := p.Relations["category"]
		assert.Equal(t, 2, len(children))
		assert.Equal(t, "items", children[0])
		assert.Equal(t, "tags", children[1])
	})

	t.Run("with_multiple_relations", func(t *testing.T) {
		t.Parallel()
		p := estype.NewJoinProperty(
			estype.WithJoinRelation("category", "items"),
			estype.WithJoinRelation("items", "tags"),
		)
		assert.Equal(t, 2, len(p.Relations))
	})
}

// TestPassthroughObjectProperty verifies that NewPassthroughObjectProperty and
// its options produce a correctly configured PassthroughObjectProperty.
func TestPassthroughObjectProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewPassthroughObjectProperty()
		assert.Equal(t, "passthrough", p.ESTypeName())
		assert.Assert(t, p.Properties == nil)
	})

	t.Run("with_property", func(t *testing.T) {
		t.Parallel()
		p := estype.NewPassthroughObjectProperty(
			estype.WithPassthroughObjectProperty("status", estype.NewKeywordProperty()),
		)
		assert.Assert(t, p.Properties != nil)
		assert.Equal(t, 1, len(p.Properties))
		sub, ok := p.Properties["status"]
		assert.Assert(t, ok)
		assert.Equal(t, "keyword", sub.ESTypeName())
	})
}

// ---------------------------------------------------------------------------
// Simple marker-type properties
// ---------------------------------------------------------------------------

// TestSimpleMarkerProperties verifies that all marker-type properties (those
// without configuration options) are correctly constructed.
func TestSimpleMarkerProperties(t *testing.T) {
	t.Parallel()

	t.Run("ip", func(t *testing.T) {
		t.Parallel()
		p := estype.NewIpProperty()
		assert.Equal(t, "ip", p.ESTypeName())
	})

	t.Run("binary", func(t *testing.T) {
		t.Parallel()
		p := estype.NewBinaryProperty()
		assert.Equal(t, "binary", p.ESTypeName())
	})

	t.Run("percolator", func(t *testing.T) {
		t.Parallel()
		p := estype.NewPercolatorProperty()
		assert.Equal(t, "percolator", p.ESTypeName())
	})

	t.Run("histogram", func(t *testing.T) {
		t.Parallel()
		p := estype.NewHistogramProperty()
		assert.Equal(t, "histogram", p.ESTypeName())
	})

	t.Run("version", func(t *testing.T) {
		t.Parallel()
		p := estype.NewVersionProperty()
		assert.Equal(t, "version", p.ESTypeName())
	})

	t.Run("sparse_vector", func(t *testing.T) {
		t.Parallel()
		p := estype.NewSparseVectorProperty()
		assert.Equal(t, "sparse_vector", p.ESTypeName())
	})

	t.Run("rank_features", func(t *testing.T) {
		t.Parallel()
		p := estype.NewRankFeaturesProperty()
		assert.Equal(t, "rank_features", p.ESTypeName())
	})

	t.Run("rank_vectors", func(t *testing.T) {
		t.Parallel()
		p := estype.NewRankVectorProperty()
		assert.Equal(t, "rank_vectors", p.ESTypeName())
	})

	t.Run("murmur3", func(t *testing.T) {
		t.Parallel()
		p := estype.NewMurmur3HashProperty()
		assert.Equal(t, "murmur3", p.ESTypeName())
	})

	t.Run("icu_collation_keyword", func(t *testing.T) {
		t.Parallel()
		p := estype.NewIcuCollationProperty()
		assert.Equal(t, "icu_collation_keyword", p.ESTypeName())
	})

	t.Run("dynamic_type", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDynamicProperty()
		assert.Equal(t, "{dynamic_type}", p.ESTypeName())
	})
}

// ---------------------------------------------------------------------------
// Token Count
// ---------------------------------------------------------------------------

// TestTokenCountProperty verifies that NewTokenCountProperty and its options
// produce a correctly configured TokenCountProperty.
func TestTokenCountProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewTokenCountProperty()
		assert.Equal(t, "token_count", p.ESTypeName())
		assert.Assert(t, p.Analyzer == nil)
	})

	t.Run("with_analyzer", func(t *testing.T) {
		t.Parallel()
		a := estype.Analyzer("standard")
		p := estype.NewTokenCountProperty(estype.WithTokenCountAnalyzer(a))
		assert.Assert(t, p.Analyzer != nil)
		assert.Equal(t, a, *p.Analyzer)
	})
}

// ---------------------------------------------------------------------------
// Field Alias
// ---------------------------------------------------------------------------

// TestFieldAliasProperty verifies that NewFieldAliasProperty and its options
// produce a correctly configured FieldAliasProperty.
func TestFieldAliasProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewFieldAliasProperty()
		assert.Equal(t, "alias", p.ESTypeName())
		assert.Assert(t, p.Path == nil)
	})

	t.Run("with_path", func(t *testing.T) {
		t.Parallel()
		p := estype.NewFieldAliasProperty(estype.WithFieldAliasPath("status"))
		assert.Assert(t, p.Path != nil)
		assert.Equal(t, "status", *p.Path)
	})
}

// ---------------------------------------------------------------------------
// Dense Vector
// ---------------------------------------------------------------------------

// TestDenseVectorProperty verifies that NewDenseVectorProperty and its options
// produce a correctly configured DenseVectorProperty.
func TestDenseVectorProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDenseVectorProperty()
		assert.Equal(t, "dense_vector", p.ESTypeName())
		assert.Assert(t, p.Dims == nil)
		assert.Assert(t, p.Similarity == nil)
	})

	t.Run("with_dims", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDenseVectorProperty(estype.WithDenseVectorDims(128))
		assert.Assert(t, p.Dims != nil)
		assert.Equal(t, 128, *p.Dims)
	})

	t.Run("with_similarity", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDenseVectorProperty(estype.WithDenseVectorSimilarity("cosine"))
		assert.Assert(t, p.Similarity != nil)
		assert.Equal(t, "cosine", *p.Similarity)
	})

	t.Run("with_all_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewDenseVectorProperty(
			estype.WithDenseVectorDims(768),
			estype.WithDenseVectorSimilarity("dot_product"),
		)
		assert.Equal(t, 768, *p.Dims)
		assert.Equal(t, "dot_product", *p.Similarity)
	})
}

// ---------------------------------------------------------------------------
// Rank Feature
// ---------------------------------------------------------------------------

// TestRankFeatureProperty verifies that NewRankFeatureProperty and its options
// produce a correctly configured RankFeatureProperty.
func TestRankFeatureProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewRankFeatureProperty()
		assert.Equal(t, "rank_feature", p.ESTypeName())
		assert.Assert(t, p.PositiveScoreImpact == nil)
	})

	t.Run("with_positive_score_impact_true", func(t *testing.T) {
		t.Parallel()
		p := estype.NewRankFeatureProperty(estype.WithRankFeaturePositiveScoreImpact(true))
		assert.Assert(t, p.PositiveScoreImpact != nil)
		assert.Assert(t, *p.PositiveScoreImpact == true)
	})

	t.Run("with_positive_score_impact_false", func(t *testing.T) {
		t.Parallel()
		p := estype.NewRankFeatureProperty(estype.WithRankFeaturePositiveScoreImpact(false))
		assert.Assert(t, p.PositiveScoreImpact != nil)
		assert.Assert(t, *p.PositiveScoreImpact == false)
	})
}

// ---------------------------------------------------------------------------
// Semantic Text
// ---------------------------------------------------------------------------

// TestSemanticTextProperty verifies that NewSemanticTextProperty and its options
// produce a correctly configured SemanticTextProperty.
func TestSemanticTextProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewSemanticTextProperty()
		assert.Equal(t, "semantic_text", p.ESTypeName())
		assert.Assert(t, p.InferenceId == nil)
	})

	t.Run("with_inference_id", func(t *testing.T) {
		t.Parallel()
		p := estype.NewSemanticTextProperty(estype.WithSemanticTextInferenceId("my_model"))
		assert.Assert(t, p.InferenceId != nil)
		assert.Equal(t, "my_model", *p.InferenceId)
	})
}

// ---------------------------------------------------------------------------
// Aggregate Metric Double
// ---------------------------------------------------------------------------

// TestAggregateMetricDoubleProperty verifies that NewAggregateMetricDoubleProperty
// and its options produce a correctly configured AggregateMetricDoubleProperty.
func TestAggregateMetricDoubleProperty(t *testing.T) {
	t.Parallel()

	t.Run("no_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewAggregateMetricDoubleProperty()
		assert.Equal(t, "aggregate_metric_double", p.ESTypeName())
		assert.Assert(t, p.DefaultMetric == nil)
		assert.Assert(t, p.Metrics == nil)
	})

	t.Run("with_default_metric", func(t *testing.T) {
		t.Parallel()
		p := estype.NewAggregateMetricDoubleProperty(
			estype.WithAggregateMetricDoubleDefaultMetric("max"),
		)
		assert.Assert(t, p.DefaultMetric != nil)
		assert.Equal(t, "max", *p.DefaultMetric)
	})

	t.Run("with_metrics", func(t *testing.T) {
		t.Parallel()
		p := estype.NewAggregateMetricDoubleProperty(
			estype.WithAggregateMetricDoubleMetrics("min", "max", "sum", "value_count"),
		)
		assert.Assert(t, p.Metrics != nil)
		assert.Equal(t, 4, len(p.Metrics))
		assert.Equal(t, "min", p.Metrics[0])
		assert.Equal(t, "max", p.Metrics[1])
		assert.Equal(t, "sum", p.Metrics[2])
		assert.Equal(t, "value_count", p.Metrics[3])
	})

	t.Run("with_all_options", func(t *testing.T) {
		t.Parallel()
		p := estype.NewAggregateMetricDoubleProperty(
			estype.WithAggregateMetricDoubleDefaultMetric("max"),
			estype.WithAggregateMetricDoubleMetrics("min", "max"),
		)
		assert.Equal(t, "max", *p.DefaultMetric)
		assert.Equal(t, 2, len(p.Metrics))
	})
}
