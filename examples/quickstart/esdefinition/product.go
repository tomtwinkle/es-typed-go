package esdefinition

import "github.com/tomtwinkle/es-typed-go/estype"

//go:generate go tool estyped -struct Product -package esmodel -out ../esmodel/product_gen.go -group Product

type Product struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"`
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Price    int      `json:"price"`
	Tags     []string `json:"tags"`
	Date     string   `json:"date"`
	Items    []Item   `json:"items"`
}

type Item struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (Product) Alias() estype.Alias { return "product" }
func (Product) Index() estype.Index { return "product-000001" }

func (Product) Settings() estype.Settings {
	return estype.Settings{
		NumberOfShards:   new(int(1)),
		NumberOfReplicas: new(int(0)),
		RefreshInterval:  new(estype.RefreshInterval(estype.RefreshIntervalDefault)),
	}
}

func (Product) Mapping() estype.Mapping {
	return estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "id", Property: estype.NewKeywordProperty()},
			{Path: "status", Property: estype.NewKeywordProperty()},
			{
				Path: "title",
				Property: estype.NewTextProperty(
					estype.WithField("keyword", estype.NewKeywordProperty(estype.WithIgnoreAbove(256))),
				),
			},
			{Path: "category", Property: estype.NewKeywordProperty()},
			{Path: "price", Property: estype.NewIntegerNumberProperty()},
			{Path: "tags", Property: estype.NewKeywordProperty()},
			{Path: "date", Property: estype.NewDateProperty()},
			{Path: "items", Property: estype.NewNestedProperty()},
			{Path: "items.name", Property: estype.NewTextProperty()},
			{Path: "items.value", Property: estype.NewIntegerNumberProperty()},
		},
	}
}
