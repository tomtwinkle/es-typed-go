//go:build integration

package esv8_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	esv8 "github.com/tomtwinkle/es-typed-go/esv8"
	"github.com/tomtwinkle/es-typed-go/esv8/query"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func esURL() string {
	if u := os.Getenv("ES_URL"); u != "" {
		return u
	}
	return "http://localhost:9200"
}

func newTestClient(t *testing.T) esv8.ESClient {
	t.Helper()
	client, err := esv8.NewClient(es8.Config{
		Addresses: []string{esURL()},
	})
	assert.NilError(t, err)
	return client
}

// uniqueIndex returns a test-local index name that is cleaned up after the test.
func uniqueIndex(t *testing.T, client esv8.ESClient) estype.Index {
	t.Helper()
	name := fmt.Sprintf("test-%s", uuid.New().String())
	idx, err := estype.ParseESIndex(name)
	assert.NilError(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Best-effort cleanup; ignore errors (index may already be deleted).
		_, _ = client.DeleteIndex(ctx, idx) //nolint:errcheck
	})
	return idx
}

// uniqueAlias returns an alias name unique to the test.
func uniqueAlias(t *testing.T) estype.Alias {
	t.Helper()
	name := fmt.Sprintf("alias-%s", uuid.New().String())
	alias, err := estype.ParseESAlias(name)
	assert.NilError(t, err)
	return alias
}

// noReplicaSettings returns index settings with zero replicas to keep the cluster
// GREEN on a single-node CI cluster and avoid interfering with parallel tests
// that use UpdateByQuery with scripts.
func noReplicaSettings() *types.IndexSettings {
	replicas := "0"
	return &types.IndexSettings{NumberOfReplicas: &replicas}
}

// productDoc is a sample document used across integration tests.
type productDoc struct {
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	InStock   bool      `json:"in_stock"`
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestIntegration_Info(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()

	res, err := client.Info(ctx)
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	assert.Assert(t, res.ClusterName != "")
	assert.Assert(t, res.Version.Int != "")
	t.Logf("Connected to Elasticsearch %s (cluster: %s)", res.Version.Int, res.ClusterName)
}

func TestIntegration_CreateDeleteIndex(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	// Create
	createRes, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	assert.Assert(t, createRes.Acknowledged)
	assert.Equal(t, idx.String(), createRes.Index)

	// Exists
	exists, err := client.IndexExists(ctx, idx)
	assert.NilError(t, err)
	assert.Assert(t, exists)

	// Delete
	deleteRes, err := client.DeleteIndex(ctx, idx)
	assert.NilError(t, err)
	assert.Assert(t, deleteRes.Acknowledged)

	// No longer exists
	exists, err = client.IndexExists(ctx, idx)
	assert.NilError(t, err)
	assert.Assert(t, !exists)
}

func TestIntegration_CreateIndexWithMappings(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name":     types.NewKeywordProperty(),
			"category": types.NewKeywordProperty(),
			"price":    types.NewDoubleNumberProperty(),
			"created_at": &types.DateProperty{
				Format: func() *string { s := "strict_date_optional_time"; return &s }(),
			},
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AliasLifecycle(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	// Create index first
	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)

	// Alias should not exist yet
	exists, err := client.AliasExists(ctx, alias)
	assert.NilError(t, err)
	assert.Assert(t, !exists)

	// GetIndicesForAlias returns empty when alias doesn't exist
	indices, err := client.GetIndicesForAlias(ctx, alias)
	assert.NilError(t, err)
	assert.Assert(t, len(indices) == 0)

	// Create alias
	createRes, err := client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)
	assert.Assert(t, createRes.Acknowledged)

	// Now exists
	exists, err = client.AliasExists(ctx, alias)
	assert.NilError(t, err)
	assert.Assert(t, exists)

	// GetIndicesForAlias returns the index
	indices, err = client.GetIndicesForAlias(ctx, alias)
	assert.NilError(t, err)
	assert.Assert(t, len(indices) == 1)
	assert.Equal(t, idx, indices[0])
}

func TestIntegration_UpdateAliases(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx1 := uniqueIndex(t, client)
	idx2 := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	// Create two indices
	for _, idx := range []estype.Index{idx1, idx2} {
		_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
		assert.NilError(t, err)
	}

	// Add idx1 to alias
	addIdx1 := types.IndicesAction{
		Add: &types.AddAction{Index: func() *string { s := idx1.String(); return &s }(), Alias: func() *string { s := alias.String(); return &s }()},
	}
	res, err := client.UpdateAliases(ctx, []types.IndicesAction{addIdx1})
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)

	// Move alias to idx2
	removeIdx1 := types.IndicesAction{
		Remove: &types.RemoveAction{Index: func() *string { s := idx1.String(); return &s }(), Alias: func() *string { s := alias.String(); return &s }()},
	}
	addIdx2 := types.IndicesAction{
		Add: &types.AddAction{Index: func() *string { s := idx2.String(); return &s }(), Alias: func() *string { s := alias.String(); return &s }()},
	}
	res, err = client.UpdateAliases(ctx, []types.IndicesAction{removeIdx1, addIdx2})
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)

	// Alias now points to idx2
	indices, err := client.GetIndicesForAlias(ctx, alias)
	assert.NilError(t, err)
	assert.Assert(t, len(indices) == 1)
	assert.Equal(t, idx2, indices[0])
}

func TestIntegration_DocumentCRUD(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	// Setup: create index and alias
	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	doc := productDoc{
		Name:      "Test Widget",
		Category:  "electronics",
		Price:     29.99,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		InStock:   true,
	}
	docID := "doc-1"

	// Create (index) document
	createRes, err := client.CreateDocument(ctx, alias, docID, doc)
	assert.NilError(t, err)
	assert.Equal(t, docID, createRes.Id_)

	// Get document
	getRes, err := client.GetDocument(ctx, alias, docID)
	assert.NilError(t, err)
	assert.Assert(t, getRes.Found)
	assert.Equal(t, docID, getRes.Id_)

	// Deserialise source and verify
	var got productDoc
	assert.NilError(t, json.Unmarshal(getRes.Source_, &got))
	assert.Equal(t, doc.Name, got.Name)
	assert.Equal(t, doc.Category, got.Category)
	assert.Assert(t, math.Abs(doc.Price-got.Price) < 0.001)

	// Update document (partial update via doc field)
	updatedName := "Updated Widget"
	updateDocBytes, err := json.Marshal(map[string]any{"name": updatedName})
	assert.NilError(t, err)
	updateReq := update.NewRequest()
	updateReq.Doc = updateDocBytes
	updateRes, err := client.UpdateDocument(ctx, idx, docID, updateReq)
	assert.NilError(t, err)
	assert.Equal(t, docID, updateRes.Id_)

	// Verify update
	getRes2, err := client.GetDocument(ctx, alias, docID)
	assert.NilError(t, err)
	var got2 productDoc
	assert.NilError(t, json.Unmarshal(getRes2.Source_, &got2))
	assert.Equal(t, updatedName, got2.Name)

	// Delete document
	deleteRes, err := client.DeleteDocument(ctx, idx, docID)
	assert.NilError(t, err)
	assert.Equal(t, docID, deleteRes.Id_)

	// Verify gone
	getRes3, err := client.GetDocument(ctx, alias, docID)
	assert.NilError(t, err)
	assert.Assert(t, !getRes3.Found)
}

func TestIntegration_IndexDocumentCount(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	// Index 3 documents
	for i := 1; i <= 3; i++ {
		doc := productDoc{Name: fmt.Sprintf("Product %d", i), Category: "test", Price: float64(i * 10)}
		_, err := client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}

	// Refresh before counting
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	countRes, err := client.IndexDocumentCount(ctx, idx)
	assert.NilError(t, err)
	assert.Equal(t, int64(3), countRes.Count)
}

func TestIntegration_AliasRefresh(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	// AliasRefresh resolves to the backing index and refreshes it
	refreshRes, err := client.AliasRefresh(ctx, alias)
	assert.NilError(t, err)
	assert.Assert(t, refreshRes != nil)
}

func TestIntegration_AliasRefresh_NoIndices(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()

	alias, _ := estype.ParseESAlias(fmt.Sprintf("alias-notexist-%s", uuid.New().String()))
	_, err := client.AliasRefresh(ctx, alias)
	assert.Assert(t, err != nil)
	assert.ErrorContains(t, err, "no indices found")
}

func TestIntegration_Search_MatchAll(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	// Index documents
	docs := []productDoc{
		{Name: "Laptop", Category: "electronics", Price: 999.99},
		{Name: "Phone", Category: "electronics", Price: 699.99},
		{Name: "T-Shirt", Category: "clothing", Price: 29.99},
		{Name: "Jeans", Category: "clothing", Price: 59.99},
		{Name: "Coffee Maker", Category: "kitchen", Price: 79.99},
	}
	for i, doc := range docs {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}

	// Refresh
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	q := query.MatchAll()
	res, err := client.Search(ctx, alias, q, 10, 0, nil, nil, nil, nil, nil)
	assert.NilError(t, err)
	assert.Equal(t, int64(5), res.Hits.Total.Value)
}

func TestIntegration_Search_TermQuery(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name":     types.NewKeywordProperty(),
			"category": types.NewKeywordProperty(),
			"price":    types.NewDoubleNumberProperty(),
		},
	}
	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	for i, doc := range []productDoc{
		{Name: "Laptop", Category: "electronics", Price: 999.99},
		{Name: "Phone", Category: "electronics", Price: 699.99},
		{Name: "T-Shirt", Category: "clothing", Price: 29.99},
	} {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	q := query.TermValue("category", "electronics")
	res, err := client.Search(ctx, alias, q, 10, 0, nil, nil, nil, nil, nil)
	assert.NilError(t, err)
	assert.Equal(t, int64(2), res.Hits.Total.Value)
}

func TestIntegration_Search_BoolQuery(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"category": types.NewKeywordProperty(),
			"in_stock": types.NewBooleanProperty(),
			"price":    types.NewDoubleNumberProperty(),
		},
	}
	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	docs := []productDoc{
		{Name: "Laptop", Category: "electronics", Price: 999.99, InStock: true},
		{Name: "Old Phone", Category: "electronics", Price: 199.99, InStock: false},
		{Name: "T-Shirt", Category: "clothing", Price: 29.99, InStock: true},
	}
	for i, doc := range docs {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	// electronics AND in_stock
	bq := query.NewBoolQuery().
		Must(
			query.TermValue("category", "electronics"),
			query.TermValue("in_stock", true),
		).Build()
	q := query.BoolQuery(bq)

	res, err := client.Search(ctx, alias, q, 10, 0, nil, nil, nil, nil, nil)
	assert.NilError(t, err)
	assert.Equal(t, int64(1), res.Hits.Total.Value)
}

func TestIntegration_Search_WithAggregations(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"category": types.NewKeywordProperty(),
			"price":    types.NewDoubleNumberProperty(),
		},
	}
	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	docs := []productDoc{
		{Name: "Laptop", Category: "electronics", Price: 999.99},
		{Name: "Phone", Category: "electronics", Price: 699.99},
		{Name: "T-Shirt", Category: "clothing", Price: 29.99},
		{Name: "Jeans", Category: "clothing", Price: 59.99},
		{Name: "Coffee Maker", Category: "kitchen", Price: 79.99},
	}
	for i, doc := range docs {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	// terms aggregation on category + avg price sub-agg
	sub := query.NewAggregations().Avg("avg_price", "price")
	aggs := query.NewAggregations().
		TermsWithSize("by_category", "category", 10).
		SubAggregations("by_category", sub).
		Build()

	q := query.MatchAll()
	res, err := client.Search(ctx, alias, q, 0, 0, nil, aggs, nil, nil, nil)
	assert.NilError(t, err)

	// Verify "by_category" aggregation exists in response
	catAgg, ok := res.Aggregations["by_category"]
	assert.Assert(t, ok, "expected by_category aggregation in response")

	// Cast to StringTermsAggregate to inspect buckets
	termsAgg, ok := catAgg.(*types.StringTermsAggregate)
	assert.Assert(t, ok, "expected StringTermsAggregate")
	buckets, ok := termsAgg.Buckets.([]types.StringTermsBucket)
	assert.Assert(t, ok)
	assert.Assert(t, len(buckets) == 3)

	// Find the electronics bucket and verify avg price
	for _, bucket := range buckets {
		if bucket.Key == "electronics" {
			avgAgg, ok := bucket.Aggregations["avg_price"]
			assert.Assert(t, ok)
			avg, ok := avgAgg.(*types.AvgAggregate)
			assert.Assert(t, ok)
			// avg of 999.99 and 699.99 ≈ 849.99
			assert.Assert(t, math.Abs(849.99-float64(*avg.Value)) < 0.1)
		}
	}
}

func TestIntegration_Search_DateHistogramAggregation(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"created_at": &types.DateProperty{
				Format: func() *string { s := "strict_date_optional_time"; return &s }(),
			},
			"price": types.NewDoubleNumberProperty(),
		},
	}
	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 6; i++ {
		doc := productDoc{
			Name:      fmt.Sprintf("Product %d", i),
			Category:  "test",
			Price:     float64(i+1) * 10,
			CreatedAt: base.AddDate(0, i, 0), // one per month
		}
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	sub := query.NewAggregations().Sum("total_sales", "price")
	aggs := query.NewAggregations().
		DateHistogramWithFormat("by_month", "created_at", "yyyy-MM", calendarinterval.Month).
		SubAggregations("by_month", sub).
		Build()

	q := query.MatchAll()
	res, err := client.Search(ctx, alias, q, 0, 0, nil, aggs, nil, nil, nil)
	assert.NilError(t, err)

	monthAgg, ok := res.Aggregations["by_month"]
	assert.Assert(t, ok)
	dateAgg, ok := monthAgg.(*types.DateHistogramAggregate)
	assert.Assert(t, ok)
	buckets, ok := dateAgg.Buckets.([]types.DateHistogramBucket)
	assert.Assert(t, ok)
	// 6 documents in 6 different months
	assert.Assert(t, len(buckets) == 6)
}

func TestIntegration_RefreshInterval(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	// Initially not explicitly set
	initial, err := client.GetRefreshInterval(ctx, alias)
	assert.NilError(t, err)
	assert.Equal(t, estype.RefreshIntervalNotSet, initial)

	// Disable refresh
	_, err = client.UpdateRefreshInterval(ctx, alias, estype.RefreshIntervalDisable)
	assert.NilError(t, err)

	disabled, err := client.GetRefreshInterval(ctx, alias)
	assert.NilError(t, err)
	assert.Equal(t, estype.RefreshIntervalDisable, disabled)

	// Restore to 1s
	_, err = client.UpdateRefreshInterval(ctx, alias, estype.RefreshIntervalDefault)
	assert.NilError(t, err)

	restored, err := client.GetRefreshInterval(ctx, alias)
	assert.NilError(t, err)
	assert.Equal(t, estype.RefreshIntervalDefault, restored)
}

func TestIntegration_Reindex(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	srcIdx := uniqueIndex(t, client)
	dstIdx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	// Create source with documents
	_, err := client.CreateIndex(ctx, srcIdx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, srcIdx, alias, true)
	assert.NilError(t, err)

	for i := 1; i <= 3; i++ {
		doc := productDoc{Name: fmt.Sprintf("Product %d", i), Category: "test", Price: float64(i * 10)}
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}
	_, err = client.IndexRefresh(ctx, srcIdx)
	assert.NilError(t, err)

	// Create destination index
	_, err = client.CreateIndex(ctx, dstIdx, noReplicaSettings(), nil)
	assert.NilError(t, err)

	// Reindex synchronously
	reindexRes, err := client.Reindex(ctx, srcIdx, dstIdx, true)
	assert.NilError(t, err)
	assert.Assert(t, reindexRes != nil)
	assert.Assert(t, len(reindexRes.Failures) == 0)

	// Verify destination count
	_, err = client.IndexRefresh(ctx, dstIdx)
	assert.NilError(t, err)
	countRes, err := client.IndexDocumentCount(ctx, dstIdx)
	assert.NilError(t, err)
	assert.Equal(t, int64(3), countRes.Count)
}

func TestIntegration_DeltaReindex(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	srcIdx := uniqueIndex(t, client)
	dstIdx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"updated_at": &types.DateProperty{
				Format: func() *string { s := "strict_date_optional_time"; return &s }(),
			},
		},
	}
	_, err := client.CreateIndex(ctx, srcIdx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	_, err = client.CreateIndex(ctx, dstIdx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, srcIdx, alias, true)
	assert.NilError(t, err)

	cutoff := time.Now().UTC()

	// Two old docs (before cutoff)
	type deltaDoc struct {
		Name      string    `json:"name"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	for i, doc := range []deltaDoc{
		{Name: "Old Doc 1", UpdatedAt: cutoff.Add(-2 * time.Hour)},
		{Name: "Old Doc 2", UpdatedAt: cutoff.Add(-1 * time.Hour)},
	} {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("old-%d", i), doc)
		assert.NilError(t, err)
	}

	// One new doc (after cutoff)
	newDoc := deltaDoc{Name: "New Doc", UpdatedAt: cutoff.Add(1 * time.Second)}
	_, err = client.CreateDocument(ctx, alias, "new-1", newDoc)
	assert.NilError(t, err)

	_, err = client.IndexRefresh(ctx, srcIdx)
	assert.NilError(t, err)

	// Delta reindex: only docs updated_at >= cutoff
	res, err := client.DeltaReindex(ctx, srcIdx, dstIdx, cutoff, "updated_at", true)
	assert.NilError(t, err)
	assert.Assert(t, len(res.Failures) == 0)

	_, err = client.IndexRefresh(ctx, dstIdx)
	assert.NilError(t, err)

	countRes, err := client.IndexDocumentCount(ctx, dstIdx)
	assert.NilError(t, err)
	assert.Equal(t, int64(1), countRes.Count)
}

func TestIntegration_Search_WithSorting(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"price": types.NewDoubleNumberProperty(),
		},
	}
	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	for i, doc := range []productDoc{
		{Name: "Expensive", Price: 500.0},
		{Name: "Cheap", Price: 10.0},
		{Name: "Mid", Price: 100.0},
	} {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	sortOrder := sortorder.Asc
	sortField := types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			"price": {Order: &sortOrder},
		},
	}
	q := query.MatchAll()
	res, err := client.Search(ctx, alias, q, 10, 0, []types.SortCombinations{sortField}, nil, nil, nil, nil)
	assert.NilError(t, err)
	assert.Assert(t, len(res.Hits.Hits) == 3)

	// Verify ascending order by price
	var prev float64
	for _, hit := range res.Hits.Hits {
		var doc productDoc
		assert.NilError(t, json.Unmarshal(hit.Source_, &doc))
		assert.Assert(t, doc.Price >= prev)
		prev = doc.Price
	}
}

func TestIntegration_Search_WithPagination(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	for i := 0; i < 10; i++ {
		doc := productDoc{Name: fmt.Sprintf("Product %d", i), Category: "test", Price: float64(i)}
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	q := query.MatchAll()

	// Page 1: first 3
	res1, err := client.Search(ctx, alias, q, 3, 0, nil, nil, nil, nil, nil)
	assert.NilError(t, err)
	assert.Assert(t, len(res1.Hits.Hits) == 3)
	assert.Equal(t, int64(10), res1.Hits.Total.Value)

	// Page 2: next 3
	res2, err := client.Search(ctx, alias, q, 3, 3, nil, nil, nil, nil, nil)
	assert.NilError(t, err)
	assert.Assert(t, len(res2.Hits.Hits) == 3)
}

func TestIntegration_Search_Request(t *testing.T) {
	t.Parallel()
	// Tests the search.Request-based helper (SearchWithRequest) that would be
	// a lower-level alternative to the high-level Search method.
	// Here we verify the existing Search method produces the same results
	// when using a more complex query.
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"category": types.NewKeywordProperty(),
			"price":    types.NewDoubleNumberProperty(),
		},
	}
	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	for i, doc := range []productDoc{
		{Name: "A", Category: "cat1", Price: 10},
		{Name: "B", Category: "cat1", Price: 20},
		{Name: "C", Category: "cat2", Price: 30},
	} {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	// Build a stats aggregation directly via the request
	aggs := query.NewAggregations().Stats("price_stats", "price").Build()
	q := query.TermValue("category", "cat1")
	res, err := client.Search(ctx, alias, q, 10, 0, nil, aggs, nil, nil, nil)
	assert.NilError(t, err)
	assert.Equal(t, int64(2), res.Hits.Total.Value)

	statsRaw, ok := res.Aggregations["price_stats"]
	assert.Assert(t, ok)
	statsAgg, ok := statsRaw.(*types.StatsAggregate)
	assert.Assert(t, ok)
	assert.Equal(t, int64(2), statsAgg.Count)
	assert.Assert(t, math.Abs(15.0-float64(*statsAgg.Avg)) < 0.001)
}

// TestIntegration_SearchWithRequest demonstrates using the lower-level
// search.Request struct directly for scenarios not covered by the high-level Search.
func TestIntegration_SearchWithRequest(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	assert.NilError(t, err)

	for i := 0; i < 5; i++ {
		doc := productDoc{Name: fmt.Sprintf("Item %d", i), Category: "test", Price: float64(i * 5)}
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		assert.NilError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	// Use SearchWithRequest for a high-level search.Request
	req := search.NewRequest()
	matchAll := types.MatchAllQuery{}
	req.Query = &types.Query{MatchAll: &matchAll}
	size := 5
	req.Size = &size
	req.Source_ = true

	res, err := client.SearchWithRequest(ctx, alias, req)
	assert.NilError(t, err)
	assert.Equal(t, int64(5), res.Hits.Total.Value)
}

// ---------------------------------------------------------------------------
// Property Mapping Integration Tests
// ---------------------------------------------------------------------------

func TestIntegration_AllPropertyMappings_TextFamily(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"title": esv8.NewTextProperty(
				esv8.WithTextAnalyzer("standard"),
				esv8.WithTextSearchAnalyzer("standard"),
				esv8.WithTextSearchQuoteAnalyzer("standard"),
				esv8.WithTextFielddata(false),
				esv8.WithTextIndex(true),
				esv8.WithTextStore(false),
				esv8.WithTextNorms(true),
				esv8.WithTextSimilarity("BM25"),
				esv8.WithTextIndexPhrases(false),
				esv8.WithTextPositionIncrementGap(100),
				esv8.WithTextRawKeyword(256),
				esv8.WithTextFields(map[string]types.Property{
					"raw": esv8.NewKeywordProperty(esv8.WithKeywordIgnoreAbove(256)),
				}),
			),
			"status": esv8.NewKeywordProperty(
				esv8.WithKeywordIgnoreAbove(256),
				esv8.WithKeywordDocValues(true),
				esv8.WithKeywordIndex(true),
				esv8.WithKeywordStore(false),
				esv8.WithKeywordNullValue("N/A"),
				esv8.WithKeywordNorms(false),
				esv8.WithKeywordSimilarity("BM25"),
				esv8.WithKeywordEagerGlobalOrdinals(false),
				esv8.WithKeywordSplitQueriesOnWhitespace(false),
			),
			"type": esv8.NewConstantKeywordProperty(),
			"tags": esv8.NewWildcardProperty(
				esv8.WithWildcardIgnoreAbove(512),
				esv8.WithWildcardNullValue(""),
			),
			"name": esv8.NewCompletionProperty(
				esv8.WithCompletionAnalyzer("standard"),
				esv8.WithCompletionSearchAnalyzer("standard"),
				esv8.WithCompletionMaxInputLength(50),
				esv8.WithCompletionPreservePositionIncrements(true),
				esv8.WithCompletionPreserveSeparators(true),
			),
			"category": esv8.NewSearchAsYouTypeProperty(
				esv8.WithSearchAsYouTypeAnalyzer("standard"),
				esv8.WithSearchAsYouTypeSearchAnalyzer("standard"),
				esv8.WithSearchAsYouTypeSearchQuoteAnalyzer("standard"),
				esv8.WithSearchAsYouTypeMaxShingleSize(3),
				esv8.WithSearchAsYouTypeIndex(true),
				esv8.WithSearchAsYouTypeStore(false),
				esv8.WithSearchAsYouTypeNorms(true),
				esv8.WithSearchAsYouTypeSimilarity("BM25"),
			),
			"value": esv8.NewMatchOnlyTextProperty(),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_Numeric(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"id": esv8.NewIntegerNumberProperty(
				esv8.WithIntegerNumberCoerce(true),
				esv8.WithIntegerNumberDocValues(true),
				esv8.WithIntegerNumberIgnoreMalformed(false),
				esv8.WithIntegerNumberIndex(true),
				esv8.WithIntegerNumberStore(false),
				esv8.WithIntegerNumberNullValue(0),
			),
			"price": esv8.NewLongNumberProperty(
				esv8.WithLongNumberCoerce(true),
				esv8.WithLongNumberDocValues(true),
				esv8.WithLongNumberIgnoreMalformed(false),
				esv8.WithLongNumberIndex(true),
				esv8.WithLongNumberStore(false),
				esv8.WithLongNumberNullValue(0),
			),
			"type": esv8.NewShortNumberProperty(
				esv8.WithShortNumberCoerce(true),
				esv8.WithShortNumberDocValues(true),
				esv8.WithShortNumberIgnoreMalformed(false),
				esv8.WithShortNumberIndex(true),
				esv8.WithShortNumberStore(false),
				esv8.WithShortNumberNullValue(0),
			),
			"status": esv8.NewByteNumberProperty(
				esv8.WithByteNumberCoerce(true),
				esv8.WithByteNumberDocValues(true),
				esv8.WithByteNumberIgnoreMalformed(false),
				esv8.WithByteNumberIndex(true),
				esv8.WithByteNumberStore(false),
				esv8.WithByteNumberNullValue(0),
			),
			"value": esv8.NewDoubleNumberProperty(
				esv8.WithDoubleNumberCoerce(true),
				esv8.WithDoubleNumberDocValues(true),
				esv8.WithDoubleNumberIgnoreMalformed(false),
				esv8.WithDoubleNumberIndex(true),
				esv8.WithDoubleNumberStore(false),
				esv8.WithDoubleNumberNullValue(0.0),
			),
			"name": esv8.NewFloatNumberProperty(
				esv8.WithFloatNumberCoerce(true),
				esv8.WithFloatNumberDocValues(true),
				esv8.WithFloatNumberIgnoreMalformed(false),
				esv8.WithFloatNumberIndex(true),
				esv8.WithFloatNumberStore(false),
				esv8.WithFloatNumberNullValue(0.0),
			),
			"category": esv8.NewHalfFloatNumberProperty(
				esv8.WithHalfFloatNumberCoerce(true),
				esv8.WithHalfFloatNumberDocValues(true),
				esv8.WithHalfFloatNumberIgnoreMalformed(false),
				esv8.WithHalfFloatNumberIndex(true),
				esv8.WithHalfFloatNumberStore(false),
				esv8.WithHalfFloatNumberNullValue(0.0),
			),
			"enabled": esv8.NewUnsignedLongNumberProperty(
				esv8.WithUnsignedLongNumberDocValues(true),
				esv8.WithUnsignedLongNumberIgnoreMalformed(false),
				esv8.WithUnsignedLongNumberIndex(true),
				esv8.WithUnsignedLongNumberStore(false),
				esv8.WithUnsignedLongNumberNullValue(0),
			),
			"tags": esv8.NewScaledFloatNumberProperty(
				esv8.WithScaledFloatNumberScalingFactor(100),
				esv8.WithScaledFloatNumberCoerce(true),
				esv8.WithScaledFloatNumberDocValues(true),
				esv8.WithScaledFloatNumberIgnoreMalformed(false),
				esv8.WithScaledFloatNumberIndex(true),
				esv8.WithScaledFloatNumberStore(false),
				esv8.WithScaledFloatNumberNullValue(0.0),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_DateAndBoolean(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"date": esv8.NewDateProperty(
				esv8.WithDateFormat(estype.DateFormatStrictDateOptionalTime, estype.DateFormatEpochMillis),
				esv8.WithDateDocValues(true),
				esv8.WithDateIgnoreMalformed(false),
				esv8.WithDateIndex(true),
				esv8.WithDateStore(false),
				esv8.WithDateLocale("en"),
			),
			"name": esv8.NewDateNanosProperty(
				esv8.WithDateNanosFormat(estype.DateFormatStrictDateOptionalTimeNanos),
				esv8.WithDateNanosDocValues(true),
				esv8.WithDateNanosIgnoreMalformed(false),
				esv8.WithDateNanosIndex(true),
				esv8.WithDateNanosStore(false),
			),
			"enabled": esv8.NewBooleanProperty(
				esv8.WithBooleanDocValues(true),
				esv8.WithBooleanIndex(true),
				esv8.WithBooleanStore(false),
				esv8.WithBooleanNullValue(false),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_Geo(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv8.NewGeoPointProperty(
				esv8.WithGeoPointIgnoreMalformed(true),
				esv8.WithGeoPointIgnoreZValue(true),
				esv8.WithGeoPointDocValues(true),
				esv8.WithGeoPointIndex(true),
				esv8.WithGeoPointStore(false),
			),
			"status": esv8.NewGeoShapeProperty(
				esv8.WithGeoShapeCoerce(true),
				esv8.WithGeoShapeIgnoreMalformed(true),
				esv8.WithGeoShapeIgnoreZValue(true),
				esv8.WithGeoShapeDocValues(true),
				esv8.WithGeoShapeIndex(true),
				esv8.WithGeoShapeStore(false),
			),
			"category": esv8.NewShapeProperty(
				esv8.WithShapeCoerce(true),
				esv8.WithShapeIgnoreMalformed(true),
				esv8.WithShapeIgnoreZValue(true),
				esv8.WithShapeDocValues(true),
			),
			"value": esv8.NewPointProperty(
				esv8.WithPointIgnoreMalformed(true),
				esv8.WithPointIgnoreZValue(true),
				esv8.WithPointDocValues(true),
				esv8.WithPointNullValue("POINT(0 0)"),
				esv8.WithPointStore(false),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_Range(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"id": esv8.NewIntegerRangeProperty(
				esv8.WithIntegerRangeCoerce(true),
				esv8.WithIntegerRangeDocValues(true),
				esv8.WithIntegerRangeIndex(true),
				esv8.WithIntegerRangeStore(false),
			),
			"price": esv8.NewLongRangeProperty(
				esv8.WithLongRangeCoerce(true),
				esv8.WithLongRangeDocValues(true),
				esv8.WithLongRangeIndex(true),
				esv8.WithLongRangeStore(false),
			),
			"value": esv8.NewFloatRangeProperty(
				esv8.WithFloatRangeCoerce(true),
				esv8.WithFloatRangeDocValues(true),
				esv8.WithFloatRangeIndex(true),
				esv8.WithFloatRangeStore(false),
			),
			"name": esv8.NewDoubleRangeProperty(
				esv8.WithDoubleRangeCoerce(true),
				esv8.WithDoubleRangeDocValues(true),
				esv8.WithDoubleRangeIndex(true),
				esv8.WithDoubleRangeStore(false),
			),
			"date": esv8.NewDateRangeProperty(
				esv8.WithDateRangeFormat(estype.DateFormatStrictDateOptionalTime),
				esv8.WithDateRangeCoerce(true),
				esv8.WithDateRangeDocValues(true),
				esv8.WithDateRangeIndex(true),
				esv8.WithDateRangeStore(false),
			),
			"status": esv8.NewIpRangeProperty(
				esv8.WithIpRangeCoerce(true),
				esv8.WithIpRangeDocValues(true),
				esv8.WithIpRangeIndex(true),
				esv8.WithIpRangeStore(false),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_ObjectAndNested(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"items": esv8.NewObjectProperty(
				esv8.WithObjectEnabled(true),
				esv8.WithObjectProperties(map[string]types.Property{
					"name":  esv8.NewKeywordProperty(),
					"price": esv8.NewDoubleNumberProperty(),
				}),
			),
			"tags": esv8.NewNestedProperty(
				esv8.WithNestedEnabled(true),
				esv8.WithNestedIncludeInParent(false),
				esv8.WithNestedIncludeInRoot(false),
				esv8.WithNestedProperties(map[string]types.Property{
					"name":  esv8.NewKeywordProperty(),
					"value": esv8.NewKeywordProperty(),
				}),
			),
			"category": esv8.NewFlattenedProperty(
				esv8.WithFlattenedDepthLimit(5),
				esv8.WithFlattenedDocValues(true),
				esv8.WithFlattenedIndex(true),
				esv8.WithFlattenedIgnoreAbove(1024),
				esv8.WithFlattenedNullValue(""),
				esv8.WithFlattenedSimilarity("BM25"),
				esv8.WithFlattenedEagerGlobalOrdinals(false),
				esv8.WithFlattenedSplitQueriesOnWhitespace(false),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_Join(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv8.NewKeywordProperty(),
			"type": esv8.NewJoinProperty(
				esv8.WithJoinRelations(map[string][]string{
					"category": {"items"},
				}),
				esv8.WithJoinEagerGlobalOrdinals(true),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_Special(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"id": esv8.NewIpProperty(
				esv8.WithIpDocValues(true),
				esv8.WithIpIgnoreMalformed(true),
				esv8.WithIpIndex(true),
				esv8.WithIpNullValue("0.0.0.0"),
				esv8.WithIpStore(false),
			),
			"name": esv8.NewBinaryProperty(
				esv8.WithBinaryDocValues(false),
				esv8.WithBinaryStore(false),
			),
			"title": esv8.NewTokenCountProperty(
				esv8.WithTokenCountAnalyzer("standard"),
				esv8.WithTokenCountDocValues(true),
				esv8.WithTokenCountIndex(true),
				esv8.WithTokenCountStore(false),
				esv8.WithTokenCountEnablePositionIncrements(true),
			),
			"status": esv8.NewHistogramProperty(
				esv8.WithHistogramIgnoreMalformed(true),
			),
			"category": esv8.NewVersionProperty(),
			"tags": esv8.NewDenseVectorProperty(
				esv8.WithDenseVectorDims(3),
				esv8.WithDenseVectorIndex(false),
			),
			"price": esv8.NewRankFeatureProperty(
				esv8.WithRankFeaturePositiveScoreImpact(true),
			),
			"enabled": esv8.NewRankFeaturesProperty(
				esv8.WithRankFeaturesPositiveScoreImpact(true),
			),
			"type": esv8.NewAggregateMetricDoubleProperty(
				esv8.WithAggregateMetricDoubleDefaultMetric("max"),
				esv8.WithAggregateMetricDoubleMetrics([]string{"min", "max", "sum", "value_count"}),
				esv8.WithAggregateMetricDoubleIgnoreMalformed(true),
			),
			"items": esv8.NewPercolatorProperty(),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_Alias(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"status": esv8.NewKeywordProperty(),
			"name":   esv8.NewFieldAliasProperty(esv8.WithFieldAliasPath("status")),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_Dynamic(t *testing.T) {
	t.Parallel()
	t.Skip("DynamicProperty is only valid within dynamic_templates, not as a standalone property")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv8.NewDynamicProperty(
				esv8.WithDynamicAnalyzer("standard"),
				esv8.WithDynamicSearchAnalyzer("standard"),
				esv8.WithDynamicCoerce(true),
				esv8.WithDynamicDocValues(true),
				esv8.WithDynamicEnabled(true),
				esv8.WithDynamicFormat("strict_date_optional_time||epoch_millis"),
				esv8.WithDynamicIgnoreMalformed(true),
				esv8.WithDynamicIndex(true),
				esv8.WithDynamicStore(false),
				esv8.WithDynamicNorms(true),
				esv8.WithDynamicLocale("en"),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_PassthroughObject(t *testing.T) {
	t.Parallel()
	t.Skip("passthrough object type is not available in Elasticsearch v8")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"items": esv8.NewPassthroughObjectProperty(
				esv8.WithPassthroughObjectEnabled(true),
				esv8.WithPassthroughObjectProperties(map[string]types.Property{
					"name": esv8.NewKeywordProperty(),
				}),
				esv8.WithPassthroughObjectPriority(10),
				esv8.WithPassthroughObjectTimeSeriesDimension(false),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_RankVector(t *testing.T) {
	t.Parallel()
	t.Skip("rank_vectors type is not available in Elasticsearch v8")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv8.NewRankVectorProperty(
				esv8.WithRankVectorDims(3),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_SparseVector(t *testing.T) {
	t.Parallel()
	t.Skip("sparse_vector type is not supported in Elasticsearch 8.0-8.10")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"value": esv8.NewSparseVectorProperty(
				esv8.WithSparseVectorStore(true),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_SemanticText(t *testing.T) {
	t.Parallel()
	t.Skip("semantic_text type is not available in Elasticsearch v8")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv8.NewSemanticTextProperty(
				esv8.WithSemanticTextInferenceId("my-elser-endpoint"),
				esv8.WithSemanticTextSearchInferenceId("my-elser-endpoint"),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_Murmur3Hash(t *testing.T) {
	t.Parallel()
	t.Skip("murmur3 type requires the mapper-murmur3 plugin")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv8.NewMurmur3HashProperty(
				esv8.WithMurmur3HashDocValues(true),
				esv8.WithMurmur3HashStore(false),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_IcuCollation(t *testing.T) {
	t.Parallel()
	t.Skip("icu_collation_keyword type requires the analysis-icu plugin")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv8.NewIcuCollationProperty(
				esv8.WithIcuCollationLanguage("en"),
				esv8.WithIcuCollationCountry("US"),
				esv8.WithIcuCollationDocValues(true),
				esv8.WithIcuCollationIndex(true),
				esv8.WithIcuCollationStore(false),
				esv8.WithIcuCollationNullValue(""),
				esv8.WithIcuCollationNorms(true),
				esv8.WithIcuCollationRules(""),
				esv8.WithIcuCollationVariant(""),
				esv8.WithIcuCollationCaseLevel(false),
				esv8.WithIcuCollationNumeric(false),
				esv8.WithIcuCollationHiraganaQuaternaryMode(false),
				esv8.WithIcuCollationVariableTop(""),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_KeywordNormalizer(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	settings := &types.IndexSettings{
		NumberOfReplicas: func() *string { s := "0"; return &s }(),
		Analysis: &types.IndexSettingsAnalysis{
			Normalizer: map[string]types.Normalizer{
				"my_normalizer": types.LowercaseNormalizer{
					Type: "lowercase",
				},
			},
		},
	}

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"status": esv8.NewKeywordProperty(
				esv8.WithKeywordNormalizer("my_normalizer"),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, settings, mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}
