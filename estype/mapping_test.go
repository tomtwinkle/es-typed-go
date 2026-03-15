package estype_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

func TestField_String(t *testing.T) {
	t.Parallel()
	f := estype.Field("status")
	assert.Equal(t, "status", f.String())
}

func TestParseMappingFullFormat(t *testing.T) {
	t.Parallel()
	data := []byte(`{
		"mappings": {
			"properties": {
				"status": { "type": "keyword" },
				"title": {
					"type": "text",
					"fields": {
						"keyword": { "type": "keyword" }
					}
				}
			}
		}
	}`)

	m, err := estype.ParseMapping(data)
	assert.NilError(t, err)
	assert.Assert(t, len(m.Fields) == 3)
	assert.Equal(t, "status", m.Fields[0].Path)
	assert.Equal(t, "keyword", m.Fields[0].Type)
	assert.Equal(t, "title", m.Fields[1].Path)
	assert.Equal(t, "text", m.Fields[1].Type)
	assert.Equal(t, "title.keyword", m.Fields[2].Path)
	assert.Equal(t, "keyword", m.Fields[2].Type)
}

func TestParseMappingAbbreviatedFormat(t *testing.T) {
	t.Parallel()
	data := []byte(`{
		"properties": {
			"id": { "type": "keyword" },
			"price": { "type": "float" }
		}
	}`)

	m, err := estype.ParseMapping(data)
	assert.NilError(t, err)
	assert.Assert(t, len(m.Fields) == 2)
	assert.Equal(t, "id", m.Fields[0].Path)
	assert.Equal(t, "price", m.Fields[1].Path)
}

func TestParseMappingNested(t *testing.T) {
	t.Parallel()
	data := []byte(`{
		"properties": {
			"items": {
				"type": "nested",
				"properties": {
					"color": { "type": "keyword" },
					"date": { "type": "date" }
				}
			}
		}
	}`)

	m, err := estype.ParseMapping(data)
	assert.NilError(t, err)
	assert.Assert(t, len(m.Fields) == 3)
	assert.Equal(t, "items", m.Fields[0].Path)
	assert.Equal(t, "nested", m.Fields[0].Type)
	assert.Equal(t, "items.color", m.Fields[1].Path)
	assert.Equal(t, "items.date", m.Fields[2].Path)
}

func TestParseMappingNoProperties(t *testing.T) {
	t.Parallel()
	data := []byte(`{"settings": {}}`)
	_, err := estype.ParseMapping(data)
	assert.Assert(t, err != nil)
}

func TestParseMappingInvalidJSON(t *testing.T) {
	t.Parallel()
	_, err := estype.ParseMapping([]byte(`{invalid}`))
	assert.Assert(t, err != nil)
}

func TestParseMappingInvalidPropertiesType(t *testing.T) {
	t.Parallel()
	// "properties" is a number rather than an object — the first unmarshal into
	// mappingRoot succeeds (it has no "mappings" key) but the second unmarshal
	// into mappingBody fails because the type is incompatible.
	_, err := estype.ParseMapping([]byte(`{"properties": 123}`))
	assert.ErrorContains(t, err, "failed to parse mapping JSON")
}

func TestParseMappingFieldsSorted(t *testing.T) {
	t.Parallel()
	data := []byte(`{
		"properties": {
			"z_field": { "type": "keyword" },
			"a_field": { "type": "keyword" },
			"m_field": { "type": "keyword" }
		}
	}`)

	m, err := estype.ParseMapping(data)
	assert.NilError(t, err)
	assert.Assert(t, len(m.Fields) == 3)
	assert.Equal(t, "a_field", m.Fields[0].Path)
	assert.Equal(t, "m_field", m.Fields[1].Path)
	assert.Equal(t, "z_field", m.Fields[2].Path)
}
