package esv8

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// NewBooleanProperty creates a new boolean property mapping.
func NewBooleanProperty() *types.BooleanProperty {
	return types.NewBooleanProperty()
}

// NewKeywordProperty creates a new keyword property mapping with the given IgnoreAbove setting.
// IgnoreAbove specifies the maximum string length; strings longer than this value are not indexed.
func NewKeywordProperty(ignoreAbove int) *types.KeywordProperty {
	prop := types.NewKeywordProperty()
	prop.IgnoreAbove = &ignoreAbove
	return prop
}

// NewIntegerNumberProperty creates a new integer number property mapping.
func NewIntegerNumberProperty() *types.IntegerNumberProperty {
	return types.NewIntegerNumberProperty()
}

// NewLongNumberProperty creates a new long number property mapping.
func NewLongNumberProperty() *types.LongNumberProperty {
	return types.NewLongNumberProperty()
}

// NewDoubleNumberProperty creates a new double number property mapping.
func NewDoubleNumberProperty() *types.DoubleNumberProperty {
	return types.NewDoubleNumberProperty()
}

// NewTextProperty creates a new text property mapping.
// If rawKeyword is non-nil, a "keyword" multi-field is added to enable
// exact-match queries alongside full-text search.
func NewTextProperty(rawKeyword *estype.RawKeyword) *types.TextProperty {
	prop := types.NewTextProperty()
	if rawKeyword != nil {
		if prop.Fields == nil {
			prop.Fields = make(map[string]types.Property)
		}
		prop.Fields["keyword"] = NewKeywordProperty(rawKeyword.IgnoreAbove)
	}
	return prop
}

// NewDateProperty creates a new date property mapping with the given formats.
// If no formats are provided, Elasticsearch uses the default format.
// Multiple formats are joined with "||".
//
// https://www.elastic.co/docs/reference/elasticsearch/mapping-reference/mapping-date-format
func NewDateProperty(formats ...estype.DateFormat) *types.DateProperty {
	prop := types.NewDateProperty()
	if len(formats) > 0 {
		format := estype.JoinDateFormats(formats...)
		prop.Format = &format
	}
	return prop
}

// NewDateNanosProperty creates a new date_nanos property mapping.
// Date nanos limits its range of dates from roughly 1970 to 2262.
//
// https://www.elastic.co/docs/reference/elasticsearch/mapping-reference/date_nanos
func NewDateNanosProperty() *types.DateNanosProperty {
	return types.NewDateNanosProperty()
}

// NewGeoPointProperty creates a new geo_point property mapping.
func NewGeoPointProperty() *types.GeoPointProperty {
	return types.NewGeoPointProperty()
}

// NewNestedProperty creates a new nested property mapping from the given type mapping.
// The nested mapping's Properties are copied to the nested property.
func NewNestedProperty(nestedMapping *types.TypeMapping) *types.NestedProperty {
	prop := types.NewNestedProperty()
	prop.Properties = nestedMapping.Properties
	return prop
}
