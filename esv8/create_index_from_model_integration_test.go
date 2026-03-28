//go:build integration

package esv8_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	esv8 "github.com/tomtwinkle/es-typed-go/esv8"
)

type createIndexModel struct{}

func (createIndexModel) Settings() estype.Settings {
	return estype.Settings{
		NumberOfShards:   new(int(1)),
		NumberOfReplicas: new(int(0)),
		RefreshInterval:  new(estype.RefreshInterval(estype.RefreshIntervalDefault)),
	}
}

func (createIndexModel) Mapping() estype.Mapping {
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

type createIndexDocument struct {
	Status string       `json:"status"`
	Title  string       `json:"title"`
	Price  int          `json:"price"`
	Date   string       `json:"date"`
	Items  []createItem `json:"items"`
}

type createItem struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestIntegration_CreateIndexFromDefinitions(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	ctx := context.Background()
	index := uniqueIndex(t, client)

	model := createIndexModel{}

	_, err := client.CreateIndexFromDefinitions(ctx, index, model.Settings(), model.Mapping())
	assert.NilError(t, err)

	assertCreateIndexSettingsApplied(t, ctx, client, index)
	assertCreateIndexMappingApplied(t, ctx, client, index)
	assertCreateIndexCanIndexAndReadDocument(t, ctx, client, index)
}

func TestIntegration_CreateIndexFromProviders(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	ctx := context.Background()
	index := uniqueIndex(t, client)

	model := createIndexModel{}

	_, err := client.CreateIndexFromProviders(ctx, index, model)
	assert.NilError(t, err)

	assertCreateIndexSettingsApplied(t, ctx, client, index)
	assertCreateIndexMappingApplied(t, ctx, client, index)
	assertCreateIndexCanIndexAndReadDocument(t, ctx, client, index)
}

func TestIntegration_CreateIndexFromProviders_WithAlias(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	ctx := context.Background()
	index := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	model := createIndexModel{}

	_, err := client.CreateIndexFromProviders(ctx, index, model)
	assert.NilError(t, err)

	createAliasRes, err := client.CreateAlias(ctx, index, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	assert.Assert(t, createAliasRes.Acknowledged)

	exists, err := client.IndexExists(ctx, index)
	assert.NilError(t, err)
	assert.Assert(t, exists)

	aliasExists, err := client.AliasExists(ctx, alias)
	assert.NilError(t, err)
	assert.Assert(t, aliasExists)

	indices, err := client.GetIndicesForAlias(ctx, alias)
	assert.NilError(t, err)
	assert.Assert(t, len(indices) == 1)
	assert.Equal(t, index, indices[0])
}

func assertCreateIndexSettingsApplied(t *testing.T, ctx context.Context, client esv8.ESClient, index estype.Index) {
	t.Helper()

	exists, err := client.IndexExists(ctx, index)
	assert.NilError(t, err)
	assert.Assert(t, exists)

	mappingRes, err := client.GetMapping(ctx, index)
	assert.NilError(t, err)

	indexMapping := requireIndexMapping(t, mappingRes, index)
	assert.Assert(t, indexMapping.Mappings.Properties != nil)
}

func assertCreateIndexMappingApplied(t *testing.T, ctx context.Context, client esv8.ESClient, index estype.Index) {
	t.Helper()

	mappingRes, err := client.GetMapping(ctx, index)
	assert.NilError(t, err)

	indexMapping := requireIndexMapping(t, mappingRes, index)
	assert.Assert(t, indexMapping.Mappings.Properties != nil)

	statusProp, ok := indexMapping.Mappings.Properties["status"]
	assert.Assert(t, ok)
	statusKeyword := requireKeywordProperty(t, statusProp)
	assert.Equal(t, "keyword", statusKeyword.Type)

	titleProp, ok := indexMapping.Mappings.Properties["title"]
	assert.Assert(t, ok)
	titleText := requireTextProperty(t, titleProp)
	assert.Equal(t, "text", titleText.Type)
	assert.Assert(t, titleText.Fields != nil)

	titleKeywordProp, ok := titleText.Fields["keyword"]
	assert.Assert(t, ok)
	titleKeyword := requireKeywordProperty(t, titleKeywordProp)
	assert.Equal(t, "keyword", titleKeyword.Type)
	assert.Assert(t, titleKeyword.IgnoreAbove != nil)
	assert.Equal(t, 256, *titleKeyword.IgnoreAbove)

	priceProp, ok := indexMapping.Mappings.Properties["price"]
	assert.Assert(t, ok)
	priceInteger := requireIntegerNumberProperty(t, priceProp)
	assert.Equal(t, "integer", priceInteger.Type)

	dateProp, ok := indexMapping.Mappings.Properties["date"]
	assert.Assert(t, ok)
	dateField := requireDateProperty(t, dateProp)
	assert.Equal(t, "date", dateField.Type)

	itemsProp, ok := indexMapping.Mappings.Properties["items"]
	assert.Assert(t, ok)
	itemsNested := requireNestedProperty(t, itemsProp)
	assert.Equal(t, "nested", itemsNested.Type)
	assert.Assert(t, itemsNested.Properties != nil)

	itemNameProp, ok := itemsNested.Properties["name"]
	assert.Assert(t, ok)
	itemNameText := requireTextProperty(t, itemNameProp)
	assert.Equal(t, "text", itemNameText.Type)

	itemValueProp, ok := itemsNested.Properties["value"]
	assert.Assert(t, ok)
	itemValueInteger := requireIntegerNumberProperty(t, itemValueProp)
	assert.Equal(t, "integer", itemValueInteger.Type)
}

func requireIndexMapping(t *testing.T, mappingRes map[string]types.IndexMappingRecord, index estype.Index) types.IndexMappingRecord {
	t.Helper()

	if indexMapping, ok := mappingRes[index.String()]; ok {
		return indexMapping
	}

	for key, indexMapping := range mappingRes {
		if key == index.String() {
			return indexMapping
		}
	}

	t.Fatalf("expected mapping response to contain index %q, got keys: %v", index.String(), keysOfIndexMappingRecord(mappingRes))
	return types.IndexMappingRecord{}
}

func requireKeywordProperty(t *testing.T, property types.Property) types.KeywordProperty {
	t.Helper()

	switch p := property.(type) {
	case types.KeywordProperty:
		return p
	case *types.KeywordProperty:
		return *p
	default:
		t.Fatalf("expected keyword property, got %T", property)
		return types.KeywordProperty{}
	}
}

func requireTextProperty(t *testing.T, property types.Property) types.TextProperty {
	t.Helper()

	switch p := property.(type) {
	case types.TextProperty:
		return p
	case *types.TextProperty:
		return *p
	default:
		t.Fatalf("expected text property, got %T", property)
		return types.TextProperty{}
	}
}

func requireIntegerNumberProperty(t *testing.T, property types.Property) types.IntegerNumberProperty {
	t.Helper()

	switch p := property.(type) {
	case types.IntegerNumberProperty:
		return p
	case *types.IntegerNumberProperty:
		return *p
	default:
		t.Fatalf("expected integer property, got %T", property)
		return types.IntegerNumberProperty{}
	}
}

func requireDateProperty(t *testing.T, property types.Property) types.DateProperty {
	t.Helper()

	switch p := property.(type) {
	case types.DateProperty:
		return p
	case *types.DateProperty:
		return *p
	default:
		t.Fatalf("expected date property, got %T", property)
		return types.DateProperty{}
	}
}

func requireNestedProperty(t *testing.T, property types.Property) types.NestedProperty {
	t.Helper()

	switch p := property.(type) {
	case types.NestedProperty:
		return p
	case *types.NestedProperty:
		return *p
	default:
		t.Fatalf("expected nested property, got %T", property)
		return types.NestedProperty{}
	}
}

func keysOfIndexMappingRecord(mappingRes map[string]types.IndexMappingRecord) []string {
	keys := make([]string, 0, len(mappingRes))
	for key := range mappingRes {
		keys = append(keys, key)
	}
	return keys
}

func assertCreateIndexCanIndexAndReadDocument(t *testing.T, ctx context.Context, client esv8.ESClient, index estype.Index) {
	t.Helper()

	aliasName := fmt.Sprintf("alias-%s", uuid.New().String())
	alias, err := estype.ParseESAlias(aliasName)
	assert.NilError(t, err)

	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = client.DeleteAlias(cleanupCtx, index, alias)
	})

	createAliasRes, err := client.CreateAlias(ctx, index, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	assert.Assert(t, createAliasRes.Acknowledged)

	doc := createIndexDocument{
		Status: "active",
		Title:  "sample title",
		Price:  120,
		Date:   "2024-01-02T15:04:05Z",
		Items: []createItem{
			{Name: "first", Value: 1},
			{Name: "second", Value: 2},
		},
	}

	docID := uuid.New().String()

	createDocRes, err := client.CreateDocument(ctx, alias, docID, doc)
	assert.NilError(t, err)
	assert.Assert(t, createDocRes.Result.String() == "created" || createDocRes.Result.String() == "updated")

	getRes, err := client.GetDocument(ctx, alias, docID)
	assert.NilError(t, err)
	assert.Assert(t, getRes.Found)

	var got createIndexDocument
	err = json.Unmarshal(getRes.Source_, &got)
	assert.NilError(t, err)

	assert.Equal(t, doc.Status, got.Status)
	assert.Equal(t, doc.Title, got.Title)
	assert.Equal(t, doc.Price, got.Price)
	assert.Equal(t, doc.Date, got.Date)
	assert.DeepEqual(t, doc.Items, got.Items)
}
