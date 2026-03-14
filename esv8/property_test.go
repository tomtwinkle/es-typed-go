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
	prop := esv8.NewBooleanProperty()
	assert.Assert(t, prop != nil)
}

func TestNewKeywordProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewKeywordProperty(256)
	assert.Assert(t, prop != nil)
	assert.Equal(t, 256, *prop.IgnoreAbove)
}

func TestNewIntegerNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewIntegerNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewLongNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewLongNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewDoubleNumberProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewDoubleNumberProperty()
	assert.Assert(t, prop != nil)
}

func TestNewTextProperty(t *testing.T) {
	t.Parallel()

	t.Run("without raw keyword", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewTextProperty(nil)
		assert.Assert(t, prop != nil)
	})

	t.Run("with raw keyword", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewTextProperty(&estype.RawKeyword{IgnoreAbove: 256})
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Fields != nil)
		_, ok := prop.Fields["keyword"]
		assert.Assert(t, ok)
	})
}

func TestNewDateProperty(t *testing.T) {
	t.Parallel()

	t.Run("with format", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateProperty(estype.DateFormatStrictDate)
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date", *prop.Format)
	})

	t.Run("with multiple formats", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateProperty(estype.DateFormatStrictDateOptionalTime, estype.DateFormatEpochMillis)
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format != nil)
		assert.Equal(t, "strict_date_optional_time||epoch_millis", *prop.Format)
	})

	t.Run("without format", func(t *testing.T) {
		t.Parallel()
		prop := esv8.NewDateProperty()
		assert.Assert(t, prop != nil)
		assert.Assert(t, prop.Format == nil)
	})
}

func TestNewDateNanosProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewDateNanosProperty()
	assert.Assert(t, prop != nil)
}

func TestNewGeoPointProperty(t *testing.T) {
	t.Parallel()
	prop := esv8.NewGeoPointProperty()
	assert.Assert(t, prop != nil)
}

func TestNewNestedProperty(t *testing.T) {
	t.Parallel()
	nestedMapping := types.NewTypeMapping()
	nestedMapping.Properties = map[string]types.Property{
		"name": esv8.NewKeywordProperty(256),
	}
	prop := esv8.NewNestedProperty(nestedMapping)
	assert.Assert(t, prop != nil)
	assert.Assert(t, prop.Properties != nil)
	_, ok := prop.Properties["name"]
	assert.Assert(t, ok)
}
