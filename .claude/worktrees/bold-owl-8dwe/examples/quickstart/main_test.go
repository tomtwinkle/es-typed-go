package main

import (
	"math"
	"reflect"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/examples/quickstart/esdefinition"
	"github.com/tomtwinkle/es-typed-go/examples/quickstart/esmodel"
)

func TestProductJSONTagsMatchQuickstartFields(t *testing.T) {
	t.Parallel()

	var product Product

	assert.Equal(t, "id", jsonTagOf(t, product, "ID"))
	assert.Equal(t, "status", jsonTagOf(t, product, "Status"))
	assert.Equal(t, "title", jsonTagOf(t, product, "Title"))
	assert.Equal(t, "category", jsonTagOf(t, product, "Category"))
	assert.Equal(t, "price", jsonTagOf(t, product, "Price"))
	assert.Equal(t, "tags", jsonTagOf(t, product, "Tags"))
	assert.Equal(t, "date", jsonTagOf(t, product, "Date"))
}

func TestProductDefinitionSettings(t *testing.T) {
	t.Parallel()

	settings := esdefinition.Product{}.Settings()

	assert.Assert(t, settings.NumberOfShards != nil)
	assert.Assert(t, settings.NumberOfReplicas != nil)
	assert.Assert(t, settings.RefreshInterval != nil)

	assert.Equal(t, 1, *settings.NumberOfShards)
	assert.Equal(t, 0, *settings.NumberOfReplicas)
	assert.Equal(t, estype.RefreshInterval(estype.RefreshIntervalDefault), *settings.RefreshInterval)
}

func TestProductDefinitionMappingFields(t *testing.T) {
	t.Parallel()

	mapping := esdefinition.Product{}.Mapping()

	assert.Equal(t, 10, len(mapping.Fields))

	got := make(map[string]string, len(mapping.Fields))
	for _, field := range mapping.Fields {
		got[field.Path] = field.TypeName()
	}

	assert.Equal(t, "keyword", got["id"])
	assert.Equal(t, "keyword", got["status"])
	assert.Equal(t, "text", got["title"])
	assert.Equal(t, "keyword", got["category"])
	assert.Equal(t, "integer", got["price"])
	assert.Equal(t, "keyword", got["tags"])
	assert.Equal(t, "date", got["date"])
	assert.Equal(t, "nested", got["items"])
	assert.Equal(t, "text", got["items.name"])
	assert.Equal(t, "integer", got["items.value"])
}

func TestGeneratedProductFields(t *testing.T) {
	t.Parallel()

	assert.Equal(t, estype.Field("category"), esmodel.Product.Category)
	assert.Equal(t, estype.Field("date"), esmodel.Product.Date)
	assert.Equal(t, estype.Field("id"), esmodel.Product.Id)
	assert.Equal(t, estype.Field("items"), esmodel.Product.Items)
	assert.Equal(t, estype.Field("items.name"), esmodel.Product.Items_Name)
	assert.Equal(t, estype.Field("items.value"), esmodel.Product.Items_Value)
	assert.Equal(t, estype.Field("price"), esmodel.Product.Price)
	assert.Equal(t, estype.Field("status"), esmodel.Product.Status)
	assert.Equal(t, estype.Field("tags"), esmodel.Product.Tags)
	assert.Equal(t, estype.Field("title"), esmodel.Product.Title)
}

func TestGeneratedProductFieldsMatchDefinitionPaths(t *testing.T) {
	t.Parallel()

	mapping := esdefinition.Product{}.Mapping()

	want := map[string]bool{
		string(esmodel.Product.Category):    true,
		string(esmodel.Product.Date):        true,
		string(esmodel.Product.Id):          true,
		string(esmodel.Product.Items):       true,
		string(esmodel.Product.Items_Name):  true,
		string(esmodel.Product.Items_Value): true,
		string(esmodel.Product.Price):       true,
		string(esmodel.Product.Status):      true,
		string(esmodel.Product.Tags):        true,
		string(esmodel.Product.Title):       true,
	}

	got := make(map[string]bool, len(mapping.Fields))
	for _, field := range mapping.Fields {
		got[field.Path] = true
	}

	assert.Equal(t, len(want), len(got))
	for path := range want {
		assert.Equal(t, true, got[path])
	}
}

func TestProductTypeShape(t *testing.T) {
	t.Parallel()

	typ := reflectTypeOf(Product{})

	assert.Equal(t, "Product", typ.Name())
	assert.Equal(t, reflect.Struct, typ.Kind())
	assert.Equal(t, 7, typ.NumField())
}

func TestProductFieldKinds(t *testing.T) {
	t.Parallel()

	typ := reflectTypeOf(Product{})

	tests := map[string]reflect.Kind{
		"ID":       reflect.String,
		"Status":   reflect.String,
		"Title":    reflect.String,
		"Category": reflect.String,
		"Price":    reflect.Float64,
		"Tags":     reflect.Slice,
		"Date":     reflect.String,
	}

	for name, wantKind := range tests {
		name := name
		wantKind := wantKind

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			field, ok := typ.FieldByName(name)
			assert.Assert(t, ok)
			assert.Equal(t, wantKind, field.Type.Kind())
		})
	}
}

func TestProductTagsElementType(t *testing.T) {
	t.Parallel()

	typ := reflectTypeOf(Product{})
	field, ok := typ.FieldByName("Tags")
	assert.Assert(t, ok)

	assert.Equal(t, reflect.Slice, field.Type.Kind())
	assert.Equal(t, reflect.String, field.Type.Elem().Kind())
}

func TestReflectTypeOf(t *testing.T) {
	t.Parallel()

	assert.Equal(t, reflect.TypeOf(Product{}), reflectTypeOf(Product{}))
	assert.Equal(t, reflect.TypeOf(&Product{}), reflectTypeOf(&Product{}))
}

func TestMathAbsBehaviorUsedByQuickstartAverageGuard(t *testing.T) {
	t.Parallel()

	assert.Equal(t, 0.0, math.Abs(0))
	assert.Equal(t, 12.5, math.Abs(-12.5))
	assert.Equal(t, math.SmallestNonzeroFloat64, math.Abs(-math.SmallestNonzeroFloat64))
}

func jsonTagOf[T any](t *testing.T, value T, fieldName string) string {
	t.Helper()

	typ := reflectTypeOf(value)
	field, ok := typ.FieldByName(fieldName)
	assert.Assert(t, ok)

	return field.Tag.Get("json")
}

func reflectTypeOf[T any](value T) reflect.Type {
	return reflect.TypeOf(value)
}
