package estype

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Mapping represents a parsed Elasticsearch index mapping.
type Mapping struct {
	// Fields contains all leaf field paths found in the mapping.
	Fields []MappingField
}

// MappingField represents a single field in an Elasticsearch mapping.
type MappingField struct {
	// Path is the dot-separated path to the field (e.g. "items.color").
	Path Field
	// Property holds the Elasticsearch property definition for this field.
	// It accepts either a plain type-name string (e.g. "keyword", "text",
	// "integer") or a typed property value such as [TextProperty] or
	// [KeywordProperty] constructed with [NewTextProperty] / [NewKeywordProperty].
	Property any
}

// TypeName returns the Elasticsearch type name for the field.
// If Property is a plain string, that string is returned directly.
// If Property implements ESTypeName() string (e.g. [TextProperty],
// [KeywordProperty]), that method is called.
// Otherwise an empty string is returned.
func (f MappingField) TypeName() string {
	switch v := f.Property.(type) {
	case string:
		return v
	case interface{ ESTypeName() string }:
		return v.ESTypeName()
	}
	return ""
}

// mappingRoot mirrors the JSON structure returned by the ES Get Mapping API.
type mappingRoot struct {
	Mappings *mappingBody `json:"mappings"`
}

// mappingBody represents the "mappings" object in the ES mapping JSON.
type mappingBody struct {
	Properties map[string]mappingProperty `json:"properties"`
}

// mappingProperty represents a single property in an ES mapping.
type mappingProperty struct {
	Type       string                    `json:"type"`
	Properties map[string]mappingProperty `json:"properties"`
	Fields     map[string]mappingProperty `json:"fields"`
}

// ESMappingProvider is implemented by types that describe their Elasticsearch
// field mapping. The estyped generator reads this method when running in struct
// mode to determine field types, so they appear correctly in generated code
// (e.g. "keyword" instead of "unknown").
//
// Place the go:generate directive and the struct definition in the same file,
// then add an ESMapping() method that returns [MappingField] entries using either
// plain type-name strings (e.g. "keyword", "text", "integer") or typed property
// values constructed with [NewTextProperty], [NewKeywordProperty], etc.:
//
//	//go:generate go tool estyped -struct Product -out product_fields.go
//
//	type Product struct {
//		Status string `json:"status"`
//		Title  string `json:"title"`
//	}
//
//	func (Product) ESMapping() Mapping {
//		return Mapping{
//			Fields: []MappingField{
//				{Path: "status", Property: NewKeywordProperty()},
//				{Path: "title",  Property: NewTextProperty()},
//			},
//		}
//	}
type ESMappingProvider interface {
	ESMapping() Mapping
}

// ParseMapping parses an Elasticsearch index mapping JSON and returns the field list.
// It accepts both the full Get Mapping API response format:
//
//	{"mappings": {"properties": { ... }}}
//
// and the abbreviated properties-only format:
//
//	{"properties": { ... }}
func ParseMapping(data []byte) (*Mapping, error) {
	// Try full format first.
	var root mappingRoot
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("failed to parse mapping JSON: %w", err)
	}

	if root.Mappings != nil && root.Mappings.Properties != nil {
		return buildMapping(root.Mappings.Properties), nil
	}

	// Fall back to abbreviated format (properties at top level).
	var body mappingBody
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, fmt.Errorf("failed to parse mapping JSON: %w", err)
	}
	if body.Properties != nil {
		return buildMapping(body.Properties), nil
	}

	return nil, fmt.Errorf("mapping JSON must contain a \"properties\" object")
}

// buildMapping walks the property tree and collects all field paths.
func buildMapping(props map[string]mappingProperty) *Mapping {
	var fields []MappingField
	collectFields("", props, &fields)
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Path < fields[j].Path
	})
	return &Mapping{Fields: fields}
}

// collectFields recursively collects field paths from the property tree.
func collectFields(prefix string, props map[string]mappingProperty, out *[]MappingField) {
	for name, prop := range props {
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}

		*out = append(*out, MappingField{Path: Field(path), Property: prop.Type})

		// Recurse into nested object/nested properties.
		if prop.Properties != nil {
			collectFields(path, prop.Properties, out)
		}

		// Recurse into multi-fields (e.g. "title.keyword", "title.ngram").
		if prop.Fields != nil {
			collectFields(path, prop.Fields, out)
		}
	}
}
