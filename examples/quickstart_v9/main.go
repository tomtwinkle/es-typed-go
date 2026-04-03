package main

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"

	es9 "github.com/elastic/go-elasticsearch/v9"

	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/esv9"
	"github.com/tomtwinkle/es-typed-go/examples/quickstart/esdefinition"
	"github.com/tomtwinkle/es-typed-go/examples/quickstart/esmodel"
	"github.com/tomtwinkle/es-typed-go/query"
)

// This example is the v9 counterpart of examples/quickstart/main.go.
// The only difference is the import of esv9 instead of esv8 and the
// corresponding es9 config type. Query building, field access, and
// aggregation handling are identical — only the import path changes.

type Product struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"`
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Price    float64  `json:"price"`
	Tags     []string `json:"tags"`
	Date     string   `json:"date"`
}

func main() {
	// Only this block differs from the v8 example: es9.Config and esv9.NewClientWithLogger.
	client, err := esv9.NewClientWithLogger(
		es9.Config{
			Addresses: []string{"http://localhost:19201"},
		},
		slog.Default(),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	index := esmodel.Product.Index
	alias := esmodel.Product.Alias

	if err := ensureFreshIndexAndAlias(ctx, client, index, alias); err != nil {
		panic(err)
	}
	if err := seedProducts(ctx, client, alias); err != nil {
		panic(err)
	}

	avgPriceAgg := query.AvgAgg("avg_price", esmodel.Product.Fields.Price)
	byCategoryAgg := query.StringTermsAgg(
		"by_category",
		esmodel.Product.Fields.Category,
		query.WithTermsSize(10),
		query.WithSubAggs(avgPriceAgg),
	)

	params := query.NewSearch().
		Where(query.TermValue(esmodel.Product.Fields.Status, "active")).
		Where(
			query.TermValue(esmodel.Product.Fields.Category, "electronics"),
			query.DateRangeQuery(esmodel.Product.Fields.Date, query.DateRangeGte("2024-06-01"), query.DateRangeLte("2024-12-31")),
		).
		Sort(
			query.NewSort().
				Field(esmodel.Product.Fields.Date, query.SortDesc).
				ScoreDesc().
				Build()...,
		).
		Aggregation(query.Aggs(
			byCategoryAgg,
		).Build()).
		Limit(1).
		Offset(0).
		Build()

	// Only this line differs: esv9.Search instead of esv8.Search.
	resp, err := esv9.Search[Product](ctx, client, alias, params)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Total hits: %d (page size=%d)\n", resp.Total, len(resp.Hits))
	fmt.Printf("Hits (%d):\n", len(resp.Hits))
	for i, hit := range resp.Hits {
		scoreText := "nil"
		if hit.Score != nil {
			scoreText = fmt.Sprintf("%.3f", *hit.Score)
		}
		fmt.Printf(
			"hit[%d]: id=%s index=%s score=%s source=%+v\n",
			i,
			hit.ID,
			hit.Index,
			scoreText,
			hit.Source,
		)
	}

	terms := resp.Aggregations.MustStringTerms(byCategoryAgg)
	for _, bucket := range terms.Buckets() {
		avg := bucket.Aggregations().MustAvg(avgPriceAgg)
		if avg.Value() != nil && math.Abs(*avg.Value()) > 0 {
			fmt.Printf("category=%s avg_price=%.2f\n", bucket.Key(), *avg.Value())
		}
	}
}

func seedProducts(ctx context.Context, client esv9.ESClient, alias estype.Alias) error {
	docs := []Product{
		{ID: "product-1", Status: "active", Title: "Noise Cancelling Headphones", Category: "electronics", Price: 199.99, Tags: []string{"audio", "wireless"}, Date: "2024-11-15"},
		{ID: "product-2", Status: "active", Title: "Go Programming Guide", Category: "books", Price: 39.99, Tags: []string{"programming", "golang"}, Date: "2024-08-10"},
		{ID: "product-3", Status: "inactive", Title: "Vintage Camera", Category: "electronics", Price: 89.50, Tags: []string{"camera", "collectible"}, Date: "2024-10-01"},
		{ID: "product-4", Status: "active", Title: "Mechanical Keyboard", Category: "electronics", Price: 129.00, Tags: []string{"keyboard", "office"}, Date: "2024-07-05"},
		{ID: "product-5", Status: "active", Title: "Standing Desk", Category: "furniture", Price: 499.00, Tags: []string{"office", "desk"}, Date: "2024-09-12"},
		{ID: "product-6", Status: "active", Title: "Wireless Mouse", Category: "electronics", Price: 59.99, Tags: []string{"mouse", "wireless"}, Date: "2024-01-20"},
	}
	for _, doc := range docs {
		if _, err := client.CreateDocument(ctx, alias, doc.ID, doc); err != nil {
			return fmt.Errorf("seed product %s: %w", doc.ID, err)
		}
	}
	fmt.Printf("seeded %d example documents\n", len(docs))
	return nil
}

func ensureFreshIndexAndAlias(
	ctx context.Context,
	client esv9.ESClient,
	index estype.Index,
	alias estype.Alias,
) error {
	exists, err := client.IndexExists(ctx, index)
	if err != nil {
		return err
	}
	if exists {
		if _, err := client.DeleteIndex(ctx, index); err != nil && !strings.Contains(err.Error(), "index_not_found_exception") {
			return fmt.Errorf("delete stale index %s: %w", index, err)
		}
		fmt.Printf("deleted stale index: %s\n", index)
	}

	if _, err := client.CreateIndexFromProviders(ctx, index, esdefinition.Product{}); err != nil {
		return fmt.Errorf("create index with settings and mapping %s: %w", index, err)
	}
	fmt.Printf("created index with settings and mapping: %s\n", index)

	aliasExists, err := client.AliasExists(ctx, alias)
	if err != nil {
		return err
	}
	if !aliasExists {
		if _, err := client.CreateAlias(ctx, index, alias, estype.WriteIndexEnabled); err != nil {
			return fmt.Errorf("create alias %s -> %s: %w", alias, index, err)
		}
		fmt.Printf("created alias: %s -> %s\n", alias, index)
	} else {
		fmt.Printf("alias already exists: %s\n", alias)
	}

	return nil
}
