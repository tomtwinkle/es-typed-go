package esv8

import (
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

type createIndexModelStub struct{}

func (createIndexModelStub) Settings() estype.Settings {
	return estype.Settings{
		NumberOfShards:   new(int(3)),
		NumberOfReplicas: new(int(1)),
		RefreshInterval:  new(estype.RefreshInterval(estype.RefreshIntervalDefault)),
	}
}

func (createIndexModelStub) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.NewKeywordProperty()},
			{
				Path: "title",
				Property: estype.NewTextProperty(
					estype.WithField("keyword", estype.NewKeywordProperty(estype.WithIgnoreAbove(256))),
				),
			},
			{Path: "price", Property: estype.NewIntegerNumberProperty()},
			{Path: "date", Property: estype.NewDateProperty()},
			{Path: "items", Property: estype.NewNestedProperty()},
			{Path: "items.name", Property: estype.NewTextProperty()},
			{Path: "items.value", Property: estype.NewIntegerNumberProperty()},
		},
	}
}

func TestToTypedIndexSettings_Empty(t *testing.T) {
	t.Parallel()

	got := toTypedIndexSettings(estype.Settings{})

	assert.Assert(t, got == nil)
}

func TestToTypedIndexSettings_AllFields(t *testing.T) {
	t.Parallel()

	got := toTypedIndexSettings(estype.Settings{
		NumberOfShards:   new(int(3)),
		NumberOfReplicas: new(int(1)),
		RefreshInterval:  new(estype.RefreshInterval(estype.RefreshIntervalDefault)),
	})

	assert.Assert(t, got != nil)
	assert.Assert(t, got.NumberOfShards != nil)
	assert.Equal(t, "3", *got.NumberOfShards)
	assert.Assert(t, got.NumberOfReplicas != nil)
	assert.Equal(t, "1", *got.NumberOfReplicas)
	assert.Assert(t, got.RefreshInterval != nil)
	assert.Equal(t, "1s", got.RefreshInterval)
}

func TestToTypedIndexSettings_DisabledRefresh(t *testing.T) {
	t.Parallel()

	got := toTypedIndexSettings(estype.Settings{
		RefreshInterval: new(estype.RefreshInterval(estype.RefreshIntervalDisable)),
	})

	assert.Assert(t, got != nil)
	assert.Assert(t, got.RefreshInterval != nil)
	assert.Equal(t, "-1", got.RefreshInterval)
}

func TestToTypedTypeMapping_Empty(t *testing.T) {
	t.Parallel()

	got := toTypedTypeMapping(estype.Mapping{})

	assert.Assert(t, got == nil)
}

func TestToTypedTypeMapping_LeafProperties(t *testing.T) {
	t.Parallel()

	got := toTypedTypeMapping(estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.NewKeywordProperty()},
			{Path: "title", Property: estype.NewTextProperty()},
			{Path: "price", Property: estype.NewIntegerNumberProperty()},
			{Path: "date", Property: estype.NewDateProperty()},
		},
	})

	assert.Assert(t, got != nil)
	assert.Assert(t, got.Properties != nil)

	statusProp, ok := got.Properties["status"]
	assert.Assert(t, ok)
	_, ok = statusProp.(types.KeywordProperty)
	assert.Assert(t, ok)

	titleProp, ok := got.Properties["title"]
	assert.Assert(t, ok)
	_, ok = titleProp.(types.TextProperty)
	assert.Assert(t, ok)

	priceProp, ok := got.Properties["price"]
	assert.Assert(t, ok)
	_, ok = priceProp.(types.IntegerNumberProperty)
	assert.Assert(t, ok)

	dateProp, ok := got.Properties["date"]
	assert.Assert(t, ok)
	_, ok = dateProp.(types.DateProperty)
	assert.Assert(t, ok)
}

func TestToTypedTypeMapping_TextMultiField(t *testing.T) {
	t.Parallel()

	got := toTypedTypeMapping(estype.Mapping{
		Fields: []estype.MappingField{
			{
				Path: "title",
				Property: estype.NewTextProperty(
					estype.WithField("keyword", estype.NewKeywordProperty(estype.WithIgnoreAbove(256))),
				),
			},
		},
	})

	assert.Assert(t, got != nil)

	titleProp, ok := got.Properties["title"]
	assert.Assert(t, ok)
	textProp, ok := titleProp.(types.TextProperty)
	assert.Assert(t, ok)
	assert.Assert(t, textProp.Fields != nil)

	keywordField, ok := textProp.Fields["keyword"]
	assert.Assert(t, ok)
	keywordProp, ok := keywordField.(types.KeywordProperty)
	assert.Assert(t, ok)
	assert.Assert(t, keywordProp.IgnoreAbove != nil)
	assert.Equal(t, 256, *keywordProp.IgnoreAbove)
}

func TestToTypedTypeMapping_NestedProperties(t *testing.T) {
	t.Parallel()

	got := toTypedTypeMapping(estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "items", Property: estype.NewNestedProperty()},
			{Path: "items.name", Property: estype.NewTextProperty()},
			{Path: "items.value", Property: estype.NewIntegerNumberProperty()},
		},
	})

	assert.Assert(t, got != nil)

	itemsProp, ok := got.Properties["items"]
	assert.Assert(t, ok)
	nestedProp, ok := itemsProp.(types.NestedProperty)
	assert.Assert(t, ok)
	assert.Assert(t, nestedProp.Properties != nil)

	nameProp, ok := nestedProp.Properties["name"]
	assert.Assert(t, ok)
	_, ok = nameProp.(types.TextProperty)
	assert.Assert(t, ok)

	valueProp, ok := nestedProp.Properties["value"]
	assert.Assert(t, ok)
	_, ok = valueProp.(types.IntegerNumberProperty)
	assert.Assert(t, ok)
}

func TestToTypedTypeMapping_FieldTypeFallbacks(t *testing.T) {
	t.Parallel()

	got := toTypedTypeMapping(estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.FieldType("keyword")},
			{Path: "title", Property: estype.FieldType("text")},
			{Path: "price", Property: estype.FieldType("integer")},
			{Path: "date", Property: estype.FieldType("date")},
			{Path: "items", Property: estype.FieldType("nested")},
		},
	})

	assert.Assert(t, got != nil)

	_, ok := got.Properties["status"].(types.KeywordProperty)
	assert.Assert(t, ok)
	_, ok = got.Properties["title"].(types.TextProperty)
	assert.Assert(t, ok)
	_, ok = got.Properties["price"].(types.IntegerNumberProperty)
	assert.Assert(t, ok)
	_, ok = got.Properties["date"].(types.DateProperty)
	assert.Assert(t, ok)
	_, ok = got.Properties["items"].(types.NestedProperty)
	assert.Assert(t, ok)
}

func TestInsertTypedProperty_CreatesNestedHierarchy(t *testing.T) {
	t.Parallel()

	props := map[string]types.Property{}
	insertTypedProperty(props, []string{"items", "name"}, estype.NewTextProperty())

	itemsProp, ok := props["items"]
	assert.Assert(t, ok)

	nestedProp, ok := itemsProp.(types.NestedProperty)
	assert.Assert(t, ok)
	assert.Assert(t, nestedProp.Properties != nil)

	nameProp, ok := nestedProp.Properties["name"]
	assert.Assert(t, ok)
	_, ok = nameProp.(types.TextProperty)
	assert.Assert(t, ok)
}

func TestCreateIndexModelStub_ProducesConvertibleSettingsAndMapping(t *testing.T) {
	t.Parallel()

	model := createIndexModelStub{}

	settings := toTypedIndexSettings(model.Settings())
	mapping := toTypedTypeMapping(model.Mapping())

	assert.Assert(t, settings != nil)
	assert.Assert(t, mapping != nil)

	assert.Assert(t, settings.NumberOfShards != nil)
	assert.Equal(t, "3", *settings.NumberOfShards)

	assert.Assert(t, settings.NumberOfReplicas != nil)
	assert.Equal(t, "1", *settings.NumberOfReplicas)

	assert.Assert(t, settings.RefreshInterval != nil)
	assert.Equal(t, "1s", settings.RefreshInterval)

	statusProp, ok := mapping.Properties["status"]
	assert.Assert(t, ok)
	_, ok = statusProp.(types.KeywordProperty)
	assert.Assert(t, ok)

	titleProp, ok := mapping.Properties["title"]
	assert.Assert(t, ok)
	textProp, ok := titleProp.(types.TextProperty)
	assert.Assert(t, ok)
	assert.Assert(t, textProp.Fields != nil)

	keywordField, ok := textProp.Fields["keyword"]
	assert.Assert(t, ok)
	_, ok = keywordField.(types.KeywordProperty)
	assert.Assert(t, ok)

	itemsProp, ok := mapping.Properties["items"]
	assert.Assert(t, ok)
	nestedProp, ok := itemsProp.(types.NestedProperty)
	assert.Assert(t, ok)

	nameProp, ok := nestedProp.Properties["name"]
	assert.Assert(t, ok)
	_, ok = nameProp.(types.TextProperty)
	assert.Assert(t, ok)

	valueProp, ok := nestedProp.Properties["value"]
	assert.Assert(t, ok)
	_, ok = valueProp.(types.IntegerNumberProperty)
	assert.Assert(t, ok)
}
