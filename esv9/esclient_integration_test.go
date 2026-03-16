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

	esv9 "github.com/tomtwinkle/es-typed-go/esv9"
	"github.com/tomtwinkle/es-typed-go/estype"
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
