package esv9

import (
	"context"
	"strconv"
	"strings"

	idxcreate "github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"

	"github.com/tomtwinkle/es-typed-go/estype"
)

// CreateIndexFromDefinitions creates an index from estype-owned settings and mapping
// definitions without requiring callers to depend on Elasticsearch typed request
// structs directly.
func (c *esClient) CreateIndexFromDefinitions(
	ctx context.Context,
	indexName estype.Index,
	settings estype.Settings,
	mapping estype.Mapping,
) (*idxcreate.Response, error) {
	return c.CreateIndex(
		ctx,
		indexName,
		toTypedIndexSettings(settings),
		toTypedTypeMapping(mapping),
	)
}

// CreateIndexFromProviders creates an index from a model/provider-owned estype
// configuration value that provides both settings and mapping definitions.
func (c *esClient) CreateIndexFromProviders(
	ctx context.Context,
	indexName estype.Index,
	provider estype.ESConfig,
) (*idxcreate.Response, error) {
	return c.CreateIndexFromDefinitions(
		ctx,
		indexName,
		provider.Settings(),
		provider.Mapping(),
	)
}

func toTypedIndexSettings(settings estype.Settings) *types.IndexSettings {
	if settings.NumberOfShards == nil &&
		settings.NumberOfReplicas == nil &&
		settings.RefreshInterval == nil {
		return nil
	}

	out := &types.IndexSettings{}

	if settings.NumberOfShards != nil {
		v := strconv.Itoa(*settings.NumberOfShards)
		out.NumberOfShards = &v
	}

	if settings.NumberOfReplicas != nil {
		v := strconv.Itoa(*settings.NumberOfReplicas)
		out.NumberOfReplicas = &v
	}

	if settings.RefreshInterval != nil {
		v := settings.RefreshInterval.String()
		out.RefreshInterval = &v
	}

	return out
}

func toTypedTypeMapping(mapping estype.Mapping) *types.TypeMapping {
	if len(mapping.Fields) == 0 {
		return nil
	}

	props := make(map[string]types.Property, len(mapping.Fields))
	for _, field := range mapping.Fields {
		if field.Path == "" || field.Property == nil {
			continue
		}
		insertTypedProperty(props, strings.Split(field.Path, "."), field.Property)
	}

	if len(props) == 0 {
		return nil
	}

	return &types.TypeMapping{
		Properties: props,
	}
}

func insertTypedProperty(props map[string]types.Property, path []string, property estype.MappingProperty) {
	if len(path) == 0 || property == nil {
		return
	}

	name := path[0]
	if name == "" {
		return
	}

	if len(path) == 1 {
		props[name] = toTypedProperty(property)
		return
	}

	existing, ok := props[name]
	if !ok {
		nested := types.NestedProperty{
			Type:       "nested",
			Properties: map[string]types.Property{},
		}
		insertTypedProperty(nested.Properties, path[1:], property)
		props[name] = nested
		return
	}

	switch p := existing.(type) {
	case types.NestedProperty:
		if p.Properties == nil {
			p.Properties = map[string]types.Property{}
		}
		insertTypedProperty(p.Properties, path[1:], property)
		props[name] = p
	case types.ObjectProperty:
		if p.Properties == nil {
			p.Properties = map[string]types.Property{}
		}
		insertTypedProperty(p.Properties, path[1:], property)
		props[name] = p
	default:
		nested := types.NestedProperty{
			Type:       "nested",
			Properties: map[string]types.Property{},
		}
		insertTypedProperty(nested.Properties, path[1:], property)
		props[name] = nested
	}
}

func toTypedProperty(property estype.MappingProperty) types.Property {
	switch p := property.(type) {
	case estype.FieldType:
		return toTypedPropertyFromTypeName(p.ESTypeName())

	case estype.KeywordProperty:
		out := types.KeywordProperty{Type: p.ESTypeName()}
		if p.IgnoreAbove != nil {
			v := *p.IgnoreAbove
			out.IgnoreAbove = &v
		}
		if p.DocValues != nil {
			v := *p.DocValues
			out.DocValues = &v
		}
		if p.Index != nil {
			v := *p.Index
			out.Index = &v
		}
		if p.Store != nil {
			v := *p.Store
			out.Store = &v
		}
		if p.NullValue != nil {
			v := *p.NullValue
			out.NullValue = &v
		}
		if p.Normalizer != nil {
			v := *p.Normalizer
			out.Normalizer = &v
		}
		if p.Norms != nil {
			v := *p.Norms
			out.Norms = &v
		}
		if p.Similarity != nil {
			v := *p.Similarity
			out.Similarity = &v
		}
		if p.EagerGlobalOrdinals != nil {
			v := *p.EagerGlobalOrdinals
			out.EagerGlobalOrdinals = &v
		}
		if p.SplitQueriesOnWhitespace != nil {
			v := *p.SplitQueriesOnWhitespace
			out.SplitQueriesOnWhitespace = &v
		}
		return out

	case estype.TextProperty:
		out := types.TextProperty{Type: p.ESTypeName()}
		if p.SearchAnalyzer != nil {
			v := p.SearchAnalyzer.String()
			out.SearchAnalyzer = &v
		}
		if p.IndexAnalyzer != nil {
			v := p.IndexAnalyzer.String()
			out.Analyzer = &v
		}
		if p.SearchQuoteAnalyzer != nil {
			v := *p.SearchQuoteAnalyzer
			out.SearchQuoteAnalyzer = &v
		}
		if p.Fielddata != nil {
			v := *p.Fielddata
			out.Fielddata = &v
		}
		if p.Index != nil {
			v := *p.Index
			out.Index = &v
		}
		if p.Store != nil {
			v := *p.Store
			out.Store = &v
		}
		if p.Norms != nil {
			v := *p.Norms
			out.Norms = &v
		}
		if p.Similarity != nil {
			v := *p.Similarity
			out.Similarity = &v
		}
		if p.IndexPhrases != nil {
			v := *p.IndexPhrases
			out.IndexPhrases = &v
		}
		if p.PositionIncrementGap != nil {
			v := *p.PositionIncrementGap
			out.PositionIncrementGap = &v
		}
		if len(p.Fields) > 0 {
			out.Fields = make(map[string]types.Property, len(p.Fields))
			for name, sub := range p.Fields {
				out.Fields[name] = toTypedProperty(sub)
			}
		}
		return out

	case estype.IntegerNumberProperty:
		out := types.IntegerNumberProperty{Type: p.ESTypeName()}
		if p.Coerce != nil {
			v := *p.Coerce
			out.Coerce = &v
		}
		if p.DocValues != nil {
			v := *p.DocValues
			out.DocValues = &v
		}
		if p.IgnoreMalformed != nil {
			v := *p.IgnoreMalformed
			out.IgnoreMalformed = &v
		}
		if p.Index != nil {
			v := *p.Index
			out.Index = &v
		}
		if p.Store != nil {
			v := *p.Store
			out.Store = &v
		}
		if p.NullValue != nil {
			v := *p.NullValue
			out.NullValue = &v
		}
		return out

	case estype.DateProperty:
		out := types.DateProperty{Type: p.ESTypeName()}
		if p.Format != nil {
			v := *p.Format
			out.Format = &v
		}
		if p.DocValues != nil {
			v := *p.DocValues
			out.DocValues = &v
		}
		if p.IgnoreMalformed != nil {
			v := *p.IgnoreMalformed
			out.IgnoreMalformed = &v
		}
		if p.Index != nil {
			v := *p.Index
			out.Index = &v
		}
		if p.Store != nil {
			v := *p.Store
			out.Store = &v
		}
		if p.Locale != nil {
			v := *p.Locale
			out.Locale = &v
		}
		return out

	case estype.NestedProperty:
		out := types.NestedProperty{
			Type: p.ESTypeName(),
		}
		if p.Enabled != nil {
			v := *p.Enabled
			out.Enabled = &v
		}
		if p.IncludeInParent != nil {
			v := *p.IncludeInParent
			out.IncludeInParent = &v
		}
		if p.IncludeInRoot != nil {
			v := *p.IncludeInRoot
			out.IncludeInRoot = &v
		}
		if p.Store != nil {
			v := *p.Store
			out.Store = &v
		}
		if len(p.Properties) > 0 {
			out.Properties = make(map[string]types.Property, len(p.Properties))
			for name, sub := range p.Properties {
				out.Properties[name] = toTypedProperty(sub)
			}
		}
		return out

	case estype.ObjectProperty:
		out := types.ObjectProperty{
			Type: p.ESTypeName(),
		}
		if p.Enabled != nil {
			v := *p.Enabled
			out.Enabled = &v
		}
		if p.Store != nil {
			v := *p.Store
			out.Store = &v
		}
		if len(p.Properties) > 0 {
			out.Properties = make(map[string]types.Property, len(p.Properties))
			for name, sub := range p.Properties {
				out.Properties[name] = toTypedProperty(sub)
			}
		}
		return out
	}

	return toTypedPropertyFromTypeName(property.ESTypeName())
}

func toTypedPropertyFromTypeName(typeName string) types.Property {
	switch typeName {
	case "keyword":
		return types.KeywordProperty{Type: "keyword"}
	case "text":
		return types.TextProperty{Type: "text"}
	case "integer":
		return types.IntegerNumberProperty{Type: "integer"}
	case "date":
		return types.DateProperty{Type: "date"}
	case "nested":
		return types.NestedProperty{Type: "nested"}
	case "object":
		return types.ObjectProperty{Type: "object"}
	default:
		return types.ObjectProperty{Type: typeName}
	}
}
