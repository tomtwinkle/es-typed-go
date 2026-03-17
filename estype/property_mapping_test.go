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
		"field_type_keyword":  {field: estype.MappingField{Path: estype.Field("status"), Property: estype.FieldType("keyword")}, want: "keyword"},
		"field_type_text":     {field: estype.MappingField{Path: estype.Field("title"), Property: estype.FieldType("text")}, want: "text"},
		"field_type_integer":  {field: estype.MappingField{Path: estype.Field("price"), Property: estype.FieldType("integer")}, want: "integer"},
		"field_type_nested":   {field: estype.MappingField{Path: estype.Field("items"), Property: estype.FieldType("nested")}, want: "nested"},
		"text_property":       {field: estype.MappingField{Path: estype.Field("title"), Property: estype.NewTextProperty()}, want: "text"},
		"keyword_property":    {field: estype.MappingField{Path: estype.Field("status"), Property: estype.NewKeywordProperty()}, want: "keyword"},
		"nil_property":        {field: estype.MappingField{Path: estype.Field("status")}, want: ""},
		"field_type_empty":    {field: estype.MappingField{Path: estype.Field("status"), Property: estype.FieldType("")}, want: ""},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := tt.field.TypeName()
			assert.Equal(t, tt.want, got)
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
// estype.Field and behaves like a typed string.
func TestMappingFieldPath(t *testing.T) {
	t.Parallel()
	f := estype.MappingField{Path: estype.Field("items.color"), Property: estype.FieldType("keyword")}
	assert.Equal(t, estype.Field("items.color"), f.Path)
	assert.Equal(t, "items.color", f.Path.String())
}
