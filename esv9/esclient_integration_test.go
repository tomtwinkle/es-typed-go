//go:build integration

package esv9_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	es9 "github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/calendarinterval"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortorder"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	esv9 "github.com/tomtwinkle/es-typed-go/esv9"
	"github.com/tomtwinkle/es-typed-go/esv9/query"
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

func newTestClient(t *testing.T) esv9.ESClient {
	t.Helper()
	client, err := esv9.NewClient(es9.Config{
		Addresses: []string{esURL()},
	})
	assert.NilError(t, err)
	return client
}

// uniqueIndex returns a test-local index name that is cleaned up after the test.
func uniqueIndex(t *testing.T, client esv9.ESClient) estype.Index {
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

func TestIntegration_AllPropertyMappings_TextFamily(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"title": esv9.NewTextProperty(
				esv9.WithTextAnalyzer("standard"),
				esv9.WithTextSearchAnalyzer("standard"),
				esv9.WithTextSearchQuoteAnalyzer("standard"),
				esv9.WithTextFielddata(false),
				esv9.WithTextIndex(true),
				esv9.WithTextStore(false),
				esv9.WithTextNorms(true),
				esv9.WithTextSimilarity("BM25"),
				esv9.WithTextIndexPhrases(false),
				esv9.WithTextPositionIncrementGap(100),
				esv9.WithTextRawKeyword(256),
				esv9.WithTextFields(map[string]types.Property{
					"raw": esv9.NewKeywordProperty(esv9.WithKeywordIgnoreAbove(256)),
				}),
			),
			"status": esv9.NewKeywordProperty(
				esv9.WithKeywordIgnoreAbove(256),
				esv9.WithKeywordDocValues(true),
				esv9.WithKeywordIndex(true),
				esv9.WithKeywordStore(false),
				esv9.WithKeywordNullValue("N/A"),
				esv9.WithKeywordNorms(false),
				esv9.WithKeywordSimilarity("BM25"),
				esv9.WithKeywordEagerGlobalOrdinals(false),
				esv9.WithKeywordSplitQueriesOnWhitespace(false),
			),
			"type": esv9.NewConstantKeywordProperty(),
			"tags": esv9.NewWildcardProperty(
				esv9.WithWildcardIgnoreAbove(512),
				esv9.WithWildcardDocValues(true),
				esv9.WithWildcardNullValue(""),
			),
			"name": esv9.NewCompletionProperty(
				esv9.WithCompletionAnalyzer("standard"),
				esv9.WithCompletionSearchAnalyzer("standard"),
				esv9.WithCompletionMaxInputLength(50),
				esv9.WithCompletionPreservePositionIncrements(true),
				esv9.WithCompletionPreserveSeparators(true),
			),
			"category": esv9.NewSearchAsYouTypeProperty(
				esv9.WithSearchAsYouTypeAnalyzer("standard"),
				esv9.WithSearchAsYouTypeSearchAnalyzer("standard"),
				esv9.WithSearchAsYouTypeSearchQuoteAnalyzer("standard"),
				esv9.WithSearchAsYouTypeMaxShingleSize(3),
				esv9.WithSearchAsYouTypeIndex(true),
				esv9.WithSearchAsYouTypeStore(false),
				esv9.WithSearchAsYouTypeNorms(true),
				esv9.WithSearchAsYouTypeSimilarity("BM25"),
			),
			"value": esv9.NewMatchOnlyTextProperty(),
			"id":    esv9.NewCountedKeywordProperty(esv9.WithCountedKeywordIndex(true)),
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
			"id": esv9.NewIntegerNumberProperty(
				esv9.WithIntegerNumberCoerce(true),
				esv9.WithIntegerNumberDocValues(true),
				esv9.WithIntegerNumberIgnoreMalformed(false),
				esv9.WithIntegerNumberIndex(true),
				esv9.WithIntegerNumberStore(false),
				esv9.WithIntegerNumberNullValue(0),
			),
			"price": esv9.NewLongNumberProperty(
				esv9.WithLongNumberCoerce(true),
				esv9.WithLongNumberDocValues(true),
				esv9.WithLongNumberIgnoreMalformed(false),
				esv9.WithLongNumberIndex(true),
				esv9.WithLongNumberStore(false),
				esv9.WithLongNumberNullValue(0),
			),
			"type": esv9.NewShortNumberProperty(
				esv9.WithShortNumberCoerce(true),
				esv9.WithShortNumberDocValues(true),
				esv9.WithShortNumberIgnoreMalformed(false),
				esv9.WithShortNumberIndex(true),
				esv9.WithShortNumberStore(false),
				esv9.WithShortNumberNullValue(0),
			),
			"status": esv9.NewByteNumberProperty(
				esv9.WithByteNumberCoerce(true),
				esv9.WithByteNumberDocValues(true),
				esv9.WithByteNumberIgnoreMalformed(false),
				esv9.WithByteNumberIndex(true),
				esv9.WithByteNumberStore(false),
				esv9.WithByteNumberNullValue(0),
			),
			"value": esv9.NewDoubleNumberProperty(
				esv9.WithDoubleNumberCoerce(true),
				esv9.WithDoubleNumberDocValues(true),
				esv9.WithDoubleNumberIgnoreMalformed(false),
				esv9.WithDoubleNumberIndex(true),
				esv9.WithDoubleNumberStore(false),
				esv9.WithDoubleNumberNullValue(0.0),
			),
			"name": esv9.NewFloatNumberProperty(
				esv9.WithFloatNumberCoerce(true),
				esv9.WithFloatNumberDocValues(true),
				esv9.WithFloatNumberIgnoreMalformed(false),
				esv9.WithFloatNumberIndex(true),
				esv9.WithFloatNumberStore(false),
				esv9.WithFloatNumberNullValue(0.0),
			),
			"category": esv9.NewHalfFloatNumberProperty(
				esv9.WithHalfFloatNumberCoerce(true),
				esv9.WithHalfFloatNumberDocValues(true),
				esv9.WithHalfFloatNumberIgnoreMalformed(false),
				esv9.WithHalfFloatNumberIndex(true),
				esv9.WithHalfFloatNumberStore(false),
				esv9.WithHalfFloatNumberNullValue(0.0),
			),
			"enabled": esv9.NewUnsignedLongNumberProperty(
				esv9.WithUnsignedLongNumberDocValues(true),
				esv9.WithUnsignedLongNumberIgnoreMalformed(false),
				esv9.WithUnsignedLongNumberIndex(true),
				esv9.WithUnsignedLongNumberStore(false),
				esv9.WithUnsignedLongNumberNullValue(0),
			),
			"tags": esv9.NewScaledFloatNumberProperty(
				esv9.WithScaledFloatNumberScalingFactor(100),
				esv9.WithScaledFloatNumberCoerce(true),
				esv9.WithScaledFloatNumberDocValues(true),
				esv9.WithScaledFloatNumberIgnoreMalformed(false),
				esv9.WithScaledFloatNumberIndex(true),
				esv9.WithScaledFloatNumberStore(false),
				esv9.WithScaledFloatNumberNullValue(0.0),
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
			"date": esv9.NewDateProperty(
				esv9.WithDateFormat(estype.DateFormatStrictDateOptionalTime, estype.DateFormatEpochMillis),
				esv9.WithDateDocValues(true),
				esv9.WithDateIgnoreMalformed(false),
				esv9.WithDateIndex(true),
				esv9.WithDateStore(false),
				esv9.WithDateLocale("en"),
			),
			"name": esv9.NewDateNanosProperty(
				esv9.WithDateNanosFormat(estype.DateFormatStrictDateOptionalTimeNanos),
				esv9.WithDateNanosDocValues(true),
				esv9.WithDateNanosIgnoreMalformed(false),
				esv9.WithDateNanosIndex(true),
				esv9.WithDateNanosStore(false),
			),
			"enabled": esv9.NewBooleanProperty(
				esv9.WithBooleanDocValues(true),
				esv9.WithBooleanIndex(true),
				esv9.WithBooleanStore(false),
				esv9.WithBooleanNullValue(false),
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
			"name": esv9.NewGeoPointProperty(
				esv9.WithGeoPointIgnoreMalformed(true),
				esv9.WithGeoPointIgnoreZValue(true),
				esv9.WithGeoPointDocValues(true),
				esv9.WithGeoPointIndex(true),
				esv9.WithGeoPointStore(false),
			),
			"status": esv9.NewGeoShapeProperty(
				esv9.WithGeoShapeCoerce(true),
				esv9.WithGeoShapeIgnoreMalformed(true),
				esv9.WithGeoShapeIgnoreZValue(true),
				esv9.WithGeoShapeDocValues(true),
				esv9.WithGeoShapeIndex(true),
				esv9.WithGeoShapeStore(false),
			),
			"category": esv9.NewShapeProperty(
				esv9.WithShapeCoerce(true),
				esv9.WithShapeIgnoreMalformed(true),
				esv9.WithShapeIgnoreZValue(true),
				esv9.WithShapeDocValues(true),
			),
			"value": esv9.NewPointProperty(
				esv9.WithPointIgnoreMalformed(true),
				esv9.WithPointIgnoreZValue(true),
				esv9.WithPointDocValues(true),
				esv9.WithPointNullValue("POINT(0 0)"),
				esv9.WithPointStore(false),
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
			"id": esv9.NewIntegerRangeProperty(
				esv9.WithIntegerRangeCoerce(true),
				esv9.WithIntegerRangeDocValues(true),
				esv9.WithIntegerRangeIndex(true),
				esv9.WithIntegerRangeStore(false),
			),
			"price": esv9.NewLongRangeProperty(
				esv9.WithLongRangeCoerce(true),
				esv9.WithLongRangeDocValues(true),
				esv9.WithLongRangeIndex(true),
				esv9.WithLongRangeStore(false),
			),
			"value": esv9.NewFloatRangeProperty(
				esv9.WithFloatRangeCoerce(true),
				esv9.WithFloatRangeDocValues(true),
				esv9.WithFloatRangeIndex(true),
				esv9.WithFloatRangeStore(false),
			),
			"name": esv9.NewDoubleRangeProperty(
				esv9.WithDoubleRangeCoerce(true),
				esv9.WithDoubleRangeDocValues(true),
				esv9.WithDoubleRangeIndex(true),
				esv9.WithDoubleRangeStore(false),
			),
			"date": esv9.NewDateRangeProperty(
				esv9.WithDateRangeFormat(estype.DateFormatStrictDateOptionalTime),
				esv9.WithDateRangeCoerce(true),
				esv9.WithDateRangeDocValues(true),
				esv9.WithDateRangeIndex(true),
				esv9.WithDateRangeStore(false),
			),
			"status": esv9.NewIpRangeProperty(
				esv9.WithIpRangeCoerce(true),
				esv9.WithIpRangeDocValues(true),
				esv9.WithIpRangeIndex(true),
				esv9.WithIpRangeStore(false),
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
			"items": esv9.NewObjectProperty(
				esv9.WithObjectEnabled(true),
				esv9.WithObjectProperties(map[string]types.Property{
					"name":  esv9.NewKeywordProperty(),
					"price": esv9.NewDoubleNumberProperty(),
				}),
			),
			"tags": esv9.NewNestedProperty(
				esv9.WithNestedEnabled(true),
				esv9.WithNestedIncludeInParent(false),
				esv9.WithNestedIncludeInRoot(false),
				esv9.WithNestedProperties(map[string]types.Property{
					"name":  esv9.NewKeywordProperty(),
					"value": esv9.NewKeywordProperty(),
				}),
			),
			"category": esv9.NewFlattenedProperty(
				esv9.WithFlattenedDepthLimit(5),
				esv9.WithFlattenedDocValues(true),
				esv9.WithFlattenedIndex(true),
				esv9.WithFlattenedIgnoreAbove(1024),
				esv9.WithFlattenedNullValue(""),
				esv9.WithFlattenedSimilarity("BM25"),
				esv9.WithFlattenedEagerGlobalOrdinals(false),
				esv9.WithFlattenedSplitQueriesOnWhitespace(false),
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
			"name": esv9.NewKeywordProperty(),
			"type": esv9.NewJoinProperty(
				esv9.WithJoinRelations(map[string][]string{
					"category": {"items"},
				}),
				esv9.WithJoinEagerGlobalOrdinals(true),
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
			"id": esv9.NewIpProperty(
				esv9.WithIpDocValues(true),
				esv9.WithIpIgnoreMalformed(true),
				esv9.WithIpIndex(true),
				esv9.WithIpNullValue("0.0.0.0"),
				esv9.WithIpStore(false),
			),
			"name": esv9.NewBinaryProperty(
				esv9.WithBinaryDocValues(false),
				esv9.WithBinaryStore(false),
			),
			"title": esv9.NewTokenCountProperty(
				esv9.WithTokenCountAnalyzer("standard"),
				esv9.WithTokenCountDocValues(true),
				esv9.WithTokenCountIndex(true),
				esv9.WithTokenCountStore(false),
				esv9.WithTokenCountEnablePositionIncrements(true),
			),
			"status": esv9.NewHistogramProperty(
				esv9.WithHistogramIgnoreMalformed(true),
			),
			"category": esv9.NewVersionProperty(
				esv9.WithVersionDocValues(true),
			),
			"tags": esv9.NewDenseVectorProperty(
				esv9.WithDenseVectorDims(3),
				esv9.WithDenseVectorIndex(false),
			),
			"value": esv9.NewSparseVectorProperty(),
			"price": esv9.NewRankFeatureProperty(
				esv9.WithRankFeaturePositiveScoreImpact(true),
			),
			"enabled": esv9.NewRankFeaturesProperty(
				esv9.WithRankFeaturesPositiveScoreImpact(true),
			),
			"type": esv9.NewAggregateMetricDoubleProperty(
				esv9.WithAggregateMetricDoubleDefaultMetric("max"),
				esv9.WithAggregateMetricDoubleMetrics([]string{"min", "max", "sum", "value_count"}),
				esv9.WithAggregateMetricDoubleIgnoreMalformed(true),
			),
			"items": esv9.NewPercolatorProperty(),
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
			"status": esv9.NewKeywordProperty(),
			"name":   esv9.NewFieldAliasProperty(esv9.WithFieldAliasPath("status")),
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
			"name": esv9.NewDynamicProperty(
				esv9.WithDynamicAnalyzer("standard"),
				esv9.WithDynamicSearchAnalyzer("standard"),
				esv9.WithDynamicCoerce(true),
				esv9.WithDynamicDocValues(true),
				esv9.WithDynamicEnabled(true),
				esv9.WithDynamicFormat("strict_date_optional_time||epoch_millis"),
				esv9.WithDynamicIgnoreMalformed(true),
				esv9.WithDynamicIndex(true),
				esv9.WithDynamicStore(false),
				esv9.WithDynamicNorms(true),
				esv9.WithDynamicLocale("en"),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_PassthroughObject(t *testing.T) {
	t.Parallel()
	t.Skip("passthrough object type requires time_series index mode")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"items": esv9.NewPassthroughObjectProperty(
				esv9.WithPassthroughObjectEnabled(true),
				esv9.WithPassthroughObjectProperties(map[string]types.Property{
					"name": esv9.NewKeywordProperty(),
				}),
				esv9.WithPassthroughObjectPriority(10),
				esv9.WithPassthroughObjectTimeSeriesDimension(false),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_RankVector(t *testing.T) {
	t.Parallel()
	t.Skip("rank_vectors type requires a specific license")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv9.NewRankVectorProperty(
				esv9.WithRankVectorDims(3),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_SemanticText(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv9.NewSemanticTextProperty(
				esv9.WithSemanticTextInferenceId("my-elser-endpoint"),
				esv9.WithSemanticTextSearchInferenceId("my-elser-endpoint"),
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
			"name": esv9.NewMurmur3HashProperty(
				esv9.WithMurmur3HashDocValues(true),
				esv9.WithMurmur3HashStore(false),
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
			"name": esv9.NewIcuCollationProperty(
				esv9.WithIcuCollationLanguage("en"),
				esv9.WithIcuCollationCountry("US"),
				esv9.WithIcuCollationDocValues(true),
				esv9.WithIcuCollationIndex(true),
				esv9.WithIcuCollationStore(false),
				esv9.WithIcuCollationNullValue(""),
				esv9.WithIcuCollationNorms(true),
				esv9.WithIcuCollationRules(""),
				esv9.WithIcuCollationVariant(""),
				esv9.WithIcuCollationCaseLevel(false),
				esv9.WithIcuCollationNumeric(false),
				esv9.WithIcuCollationHiraganaQuaternaryMode(false),
				esv9.WithIcuCollationVariableTop(""),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}

func TestIntegration_AllPropertyMappings_ExponentialHistogram(t *testing.T) {
	t.Parallel()
	t.Skip("exponential_histogram type is not available in standard Elasticsearch")
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"name": esv9.NewExponentialHistogramProperty(),
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
			"status": esv9.NewKeywordProperty(
				esv9.WithKeywordNormalizer("my_normalizer"),
			),
		},
	}

	res, err := client.CreateIndex(ctx, idx, settings, mappings)
	assert.NilError(t, err)
	assert.Assert(t, res.Acknowledged)
}
