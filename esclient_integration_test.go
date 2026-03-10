//go:build integration

package estypedgo_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/calendarinterval"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	estypedgo "github.com/tomtwinkle/es-typed-go"
	"github.com/tomtwinkle/es-typed-go/estype"
	"github.com/tomtwinkle/es-typed-go/query"
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

func newTestClient(t *testing.T) estypedgo.ESClient {
	t.Helper()
	client, err := estypedgo.NewClient(es8.Config{
		Addresses: []string{esURL()},
	})
	require.NoError(t, err)
	return client
}

// uniqueIndex returns a test-local index name that is cleaned up after the test.
func uniqueIndex(t *testing.T, client estypedgo.ESClient) estype.Index {
	t.Helper()
	name := fmt.Sprintf("test-%d", time.Now().UnixNano())
	idx, err := estype.ParseESIndex(name)
	require.NoError(t, err)
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
	name := fmt.Sprintf("alias-%d", time.Now().UnixNano())
	alias, err := estype.ParseESAlias(name)
	require.NoError(t, err)
	return alias
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
	client := newTestClient(t)
	ctx := context.Background()

	res, err := client.Info(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.NotEmpty(t, res.ClusterName)
	assert.NotEmpty(t, res.Version.Int)
	t.Logf("Connected to Elasticsearch %s (cluster: %s)", res.Version.Int, res.ClusterName)
}

func TestIntegration_CreateDeleteIndex(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	// Create
	createRes, err := client.CreateIndex(ctx, idx, nil, nil)
	require.NoError(t, err)
	assert.True(t, createRes.Acknowledged)
	assert.Equal(t, idx.String(), createRes.Index)

	// Exists
	exists, err := client.IndexExists(ctx, idx)
	require.NoError(t, err)
	assert.True(t, exists)

	// Delete
	deleteRes, err := client.DeleteIndex(ctx, idx)
	require.NoError(t, err)
	assert.True(t, deleteRes.Acknowledged)

	// No longer exists
	exists, err = client.IndexExists(ctx, idx)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestIntegration_CreateIndexWithMappings(t *testing.T) {
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

	res, err := client.CreateIndex(ctx, idx, nil, mappings)
	require.NoError(t, err)
	assert.True(t, res.Acknowledged)
}

func TestIntegration_AliasLifecycle(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	// Create index first
	_, err := client.CreateIndex(ctx, idx, nil, nil)
	require.NoError(t, err)

	// Alias should not exist yet
	exists, err := client.AliasExists(ctx, alias)
	require.NoError(t, err)
	assert.False(t, exists)

	// GetIndicesForAlias returns empty when alias doesn't exist
	indices, err := client.GetIndicesForAlias(ctx, alias)
	require.NoError(t, err)
	assert.Empty(t, indices)

	// Create alias
	createRes, err := client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)
	assert.True(t, createRes.Acknowledged)

	// Now exists
	exists, err = client.AliasExists(ctx, alias)
	require.NoError(t, err)
	assert.True(t, exists)

	// GetIndicesForAlias returns the index
	indices, err = client.GetIndicesForAlias(ctx, alias)
	require.NoError(t, err)
	require.Len(t, indices, 1)
	assert.Equal(t, idx, indices[0])
}

func TestIntegration_UpdateAliases(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx1 := uniqueIndex(t, client)
	idx2 := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	// Create two indices
	for _, idx := range []estype.Index{idx1, idx2} {
		_, err := client.CreateIndex(ctx, idx, nil, nil)
		require.NoError(t, err)
	}

	// Add idx1 to alias
	addIdx1 := types.IndicesAction{
		Add: &types.AddAction{Index: func() *string { s := idx1.String(); return &s }(), Alias: func() *string { s := alias.String(); return &s }()},
	}
	res, err := client.UpdateAliases(ctx, []types.IndicesAction{addIdx1})
	require.NoError(t, err)
	assert.True(t, res.Acknowledged)

	// Move alias to idx2
	removeIdx1 := types.IndicesAction{
		Remove: &types.RemoveAction{Index: func() *string { s := idx1.String(); return &s }(), Alias: func() *string { s := alias.String(); return &s }()},
	}
	addIdx2 := types.IndicesAction{
		Add: &types.AddAction{Index: func() *string { s := idx2.String(); return &s }(), Alias: func() *string { s := alias.String(); return &s }()},
	}
	res, err = client.UpdateAliases(ctx, []types.IndicesAction{removeIdx1, addIdx2})
	require.NoError(t, err)
	assert.True(t, res.Acknowledged)

	// Alias now points to idx2
	indices, err := client.GetIndicesForAlias(ctx, alias)
	require.NoError(t, err)
	require.Len(t, indices, 1)
	assert.Equal(t, idx2, indices[0])
}

func TestIntegration_DocumentCRUD(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	// Setup: create index and alias
	_, err := client.CreateIndex(ctx, idx, nil, nil)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

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
	require.NoError(t, err)
	assert.Equal(t, docID, createRes.Id_)

	// Get document
	getRes, err := client.GetDocument(ctx, alias, docID)
	require.NoError(t, err)
	assert.True(t, getRes.Found)
	assert.Equal(t, docID, getRes.Id_)

	// Deserialise source and verify
	var got productDoc
	require.NoError(t, json.Unmarshal(getRes.Source_, &got))
	assert.Equal(t, doc.Name, got.Name)
	assert.Equal(t, doc.Category, got.Category)
	assert.InDelta(t, doc.Price, got.Price, 0.001)

	// Update document (partial update via doc field)
	updatedName := "Updated Widget"
	updateDocBytes, err := json.Marshal(map[string]any{"name": updatedName})
	require.NoError(t, err)
	updateReq := update.NewRequest()
	updateReq.Doc = updateDocBytes
	updateRes, err := client.UpdateDocument(ctx, idx, docID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, docID, updateRes.Id_)

	// Verify update
	getRes2, err := client.GetDocument(ctx, alias, docID)
	require.NoError(t, err)
	var got2 productDoc
	require.NoError(t, json.Unmarshal(getRes2.Source_, &got2))
	assert.Equal(t, updatedName, got2.Name)

	// Delete document
	deleteRes, err := client.DeleteDocument(ctx, idx, docID)
	require.NoError(t, err)
	assert.Equal(t, docID, deleteRes.Id_)

	// Verify gone
	getRes3, err := client.GetDocument(ctx, alias, docID)
	require.NoError(t, err)
	assert.False(t, getRes3.Found)
}

func TestIntegration_IndexDocumentCount(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, nil, nil)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	// Index 3 documents
	for i := 1; i <= 3; i++ {
		doc := productDoc{Name: fmt.Sprintf("Product %d", i), Category: "test", Price: float64(i * 10)}
		_, err := client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}

	// Refresh before counting
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	countRes, err := client.IndexDocumentCount(ctx, idx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), countRes.Count)
}

func TestIntegration_AliasRefresh(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, nil, nil)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	// AliasRefresh resolves to the backing index and refreshes it
	refreshRes, err := client.AliasRefresh(ctx, alias)
	require.NoError(t, err)
	require.NotNil(t, refreshRes)
}

func TestIntegration_AliasRefresh_NoIndices(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	alias, _ := estype.ParseESAlias(fmt.Sprintf("alias-notexist-%d", time.Now().UnixNano()))
	_, err := client.AliasRefresh(ctx, alias)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no indices found")
}

func TestIntegration_Search_MatchAll(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, nil, nil)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

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
		require.NoError(t, err)
	}

	// Refresh
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	q := query.New().MatchAll(&types.MatchAllQuery{}).Build()
	res, err := client.Search(ctx, alias, q, 10, 0, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(5), res.Hits.Total.Value)
}

func TestIntegration_Search_TermQuery(t *testing.T) {
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
	_, err := client.CreateIndex(ctx, idx, nil, mappings)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	for i, doc := range []productDoc{
		{Name: "Laptop", Category: "electronics", Price: 999.99},
		{Name: "Phone", Category: "electronics", Price: 699.99},
		{Name: "T-Shirt", Category: "clothing", Price: 29.99},
	} {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	q := query.New().Term("category", types.TermQuery{Value: "electronics"}).Build()
	res, err := client.Search(ctx, alias, q, 10, 0, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(2), res.Hits.Total.Value)
}

func TestIntegration_Search_BoolQuery(t *testing.T) {
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
	_, err := client.CreateIndex(ctx, idx, nil, mappings)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	docs := []productDoc{
		{Name: "Laptop", Category: "electronics", Price: 999.99, InStock: true},
		{Name: "Old Phone", Category: "electronics", Price: 199.99, InStock: false},
		{Name: "T-Shirt", Category: "clothing", Price: 29.99, InStock: true},
	}
	for i, doc := range docs {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	// electronics AND in_stock
	bq := query.NewBoolQuery().
		Must(
			query.New().Term("category", types.TermQuery{Value: "electronics"}).Build(),
			query.New().Term("in_stock", types.TermQuery{Value: true}).Build(),
		).Build()
	q := query.New().Bool(bq).Build()

	res, err := client.Search(ctx, alias, q, 10, 0, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(1), res.Hits.Total.Value)
}

func TestIntegration_Search_WithAggregations(t *testing.T) {
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
	_, err := client.CreateIndex(ctx, idx, nil, mappings)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	docs := []productDoc{
		{Name: "Laptop", Category: "electronics", Price: 999.99},
		{Name: "Phone", Category: "electronics", Price: 699.99},
		{Name: "T-Shirt", Category: "clothing", Price: 29.99},
		{Name: "Jeans", Category: "clothing", Price: 59.99},
		{Name: "Coffee Maker", Category: "kitchen", Price: 79.99},
	}
	for i, doc := range docs {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	// terms aggregation on category + avg price sub-agg
	sub := query.NewAggregations().Avg("avg_price", "price")
	aggs := query.NewAggregations().
		TermsWithSize("by_category", "category", 10).
		SubAggregations("by_category", sub).
		Build()

	q := query.New().MatchAll(&types.MatchAllQuery{}).Build()
	res, err := client.Search(ctx, alias, q, 0, 0, nil, aggs, nil, nil, nil)
	require.NoError(t, err)

	// Verify "by_category" aggregation exists in response
	catAgg, ok := res.Aggregations["by_category"]
	require.True(t, ok, "expected by_category aggregation in response")

	// Cast to StringTermsAggregate to inspect buckets
	termsAgg, ok := catAgg.(*types.StringTermsAggregate)
	require.True(t, ok, "expected StringTermsAggregate")
	buckets, ok := termsAgg.Buckets.([]types.StringTermsBucket)
	require.True(t, ok)
	assert.Len(t, buckets, 3)

	// Find the electronics bucket and verify avg price
	for _, bucket := range buckets {
		if bucket.Key == "electronics" {
			avgAgg, ok := bucket.Aggregations["avg_price"]
			require.True(t, ok)
			avg, ok := avgAgg.(*types.AvgAggregate)
			require.True(t, ok)
			// avg of 999.99 and 699.99 ≈ 849.99
			assert.InDelta(t, 849.99, *avg.Value, 0.1)
		}
	}
}

func TestIntegration_Search_DateHistogramAggregation(t *testing.T) {
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
	_, err := client.CreateIndex(ctx, idx, nil, mappings)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 6; i++ {
		doc := productDoc{
			Name:      fmt.Sprintf("Product %d", i),
			Category:  "test",
			Price:     float64(i+1) * 10,
			CreatedAt: base.AddDate(0, i, 0), // one per month
		}
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	sub := query.NewAggregations().Sum("total_sales", "price")
	aggs := query.NewAggregations().
		DateHistogramWithFormat("by_month", "created_at", "yyyy-MM", calendarinterval.Month).
		SubAggregations("by_month", sub).
		Build()

	q := query.New().MatchAll(&types.MatchAllQuery{}).Build()
	res, err := client.Search(ctx, alias, q, 0, 0, nil, aggs, nil, nil, nil)
	require.NoError(t, err)

	monthAgg, ok := res.Aggregations["by_month"]
	require.True(t, ok)
	dateAgg, ok := monthAgg.(*types.DateHistogramAggregate)
	require.True(t, ok)
	buckets, ok := dateAgg.Buckets.([]types.DateHistogramBucket)
	require.True(t, ok)
	// 6 documents in 6 different months
	assert.Len(t, buckets, 6)
}

func TestIntegration_RefreshInterval(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, nil, nil)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	// Initially not explicitly set
	initial, err := client.GetRefreshInterval(ctx, alias)
	require.NoError(t, err)
	assert.Equal(t, estype.RefreshIntervalNotSet, initial)

	// Disable refresh
	_, err = client.UpdateRefreshInterval(ctx, alias, estype.RefreshIntervalDisable)
	require.NoError(t, err)

	disabled, err := client.GetRefreshInterval(ctx, alias)
	require.NoError(t, err)
	assert.Equal(t, estype.RefreshIntervalDisable, disabled)

	// Restore to 1s
	_, err = client.UpdateRefreshInterval(ctx, alias, estype.RefreshIntervalDefault)
	require.NoError(t, err)

	restored, err := client.GetRefreshInterval(ctx, alias)
	require.NoError(t, err)
	assert.Equal(t, estype.RefreshIntervalDefault, restored)
}

func TestIntegration_Reindex(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	srcIdx := uniqueIndex(t, client)
	dstIdx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	// Create source with documents
	_, err := client.CreateIndex(ctx, srcIdx, nil, nil)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, srcIdx, alias, true)
	require.NoError(t, err)

	for i := 1; i <= 3; i++ {
		doc := productDoc{Name: fmt.Sprintf("Product %d", i), Category: "test", Price: float64(i * 10)}
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}
	_, err = client.IndexRefresh(ctx, srcIdx)
	require.NoError(t, err)

	// Create destination index
	_, err = client.CreateIndex(ctx, dstIdx, nil, nil)
	require.NoError(t, err)

	// Reindex synchronously
	reindexRes, err := client.Reindex(ctx, srcIdx, dstIdx, true)
	require.NoError(t, err)
	require.NotNil(t, reindexRes)
	assert.Empty(t, reindexRes.Failures)

	// Verify destination count
	_, err = client.IndexRefresh(ctx, dstIdx)
	require.NoError(t, err)
	countRes, err := client.IndexDocumentCount(ctx, dstIdx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), countRes.Count)
}

func TestIntegration_DeltaReindex(t *testing.T) {
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
	_, err := client.CreateIndex(ctx, srcIdx, nil, mappings)
	require.NoError(t, err)
	_, err = client.CreateIndex(ctx, dstIdx, nil, mappings)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, srcIdx, alias, true)
	require.NoError(t, err)

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
		require.NoError(t, err)
	}

	// One new doc (after cutoff)
	newDoc := deltaDoc{Name: "New Doc", UpdatedAt: cutoff.Add(1 * time.Second)}
	_, err = client.CreateDocument(ctx, alias, "new-1", newDoc)
	require.NoError(t, err)

	_, err = client.IndexRefresh(ctx, srcIdx)
	require.NoError(t, err)

	// Delta reindex: only docs updated_at >= cutoff
	res, err := client.DeltaReindex(ctx, srcIdx, dstIdx, cutoff, "updated_at", true)
	require.NoError(t, err)
	assert.Empty(t, res.Failures)

	_, err = client.IndexRefresh(ctx, dstIdx)
	require.NoError(t, err)

	countRes, err := client.IndexDocumentCount(ctx, dstIdx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), countRes.Count)
}

func TestIntegration_Search_WithSorting(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			"price": types.NewDoubleNumberProperty(),
		},
	}
	_, err := client.CreateIndex(ctx, idx, nil, mappings)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	for i, doc := range []productDoc{
		{Name: "Expensive", Price: 500.0},
		{Name: "Cheap", Price: 10.0},
		{Name: "Mid", Price: 100.0},
	} {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	sortOrder := sortorder.Asc
	sortField := types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			"price": {Order: &sortOrder},
		},
	}
	q := query.New().MatchAll(&types.MatchAllQuery{}).Build()
	res, err := client.Search(ctx, alias, q, 10, 0, []types.SortCombinations{sortField}, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, res.Hits.Hits, 3)

	// Verify ascending order by price
	var prev float64
	for _, hit := range res.Hits.Hits {
		var doc productDoc
		require.NoError(t, json.Unmarshal(hit.Source_, &doc))
		assert.GreaterOrEqual(t, doc.Price, prev)
		prev = doc.Price
	}
}

func TestIntegration_Search_WithPagination(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, nil, nil)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		doc := productDoc{Name: fmt.Sprintf("Product %d", i), Category: "test", Price: float64(i)}
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	q := query.New().MatchAll(&types.MatchAllQuery{}).Build()

	// Page 1: first 3
	res1, err := client.Search(ctx, alias, q, 3, 0, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	assert.Len(t, res1.Hits.Hits, 3)
	assert.Equal(t, int64(10), res1.Hits.Total.Value)

	// Page 2: next 3
	res2, err := client.Search(ctx, alias, q, 3, 3, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	assert.Len(t, res2.Hits.Hits, 3)
}

func TestIntegration_Search_Request(t *testing.T) {
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
	_, err := client.CreateIndex(ctx, idx, nil, mappings)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	for i, doc := range []productDoc{
		{Name: "A", Category: "cat1", Price: 10},
		{Name: "B", Category: "cat1", Price: 20},
		{Name: "C", Category: "cat2", Price: 30},
	} {
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	// Build a stats aggregation directly via the request
	aggs := query.NewAggregations().Stats("price_stats", "price").Build()
	q := query.New().Term("category", types.TermQuery{Value: "cat1"}).Build()
	res, err := client.Search(ctx, alias, q, 10, 0, nil, aggs, nil, nil, nil)
	require.NoError(t, err)
	assert.Equal(t, int64(2), res.Hits.Total.Value)

	statsRaw, ok := res.Aggregations["price_stats"]
	require.True(t, ok)
	statsAgg, ok := statsRaw.(*types.StatsAggregate)
	require.True(t, ok)
	assert.Equal(t, int64(2), statsAgg.Count)
	assert.InDelta(t, 15.0, *statsAgg.Avg, 0.001)
}

// TestIntegration_SearchWithRequest demonstrates using the lower-level
// search.Request struct directly for scenarios not covered by the high-level Search.
func TestIntegration_SearchWithRequest(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, nil, nil)
	require.NoError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, true)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		doc := productDoc{Name: fmt.Sprintf("Item %d", i), Category: "test", Price: float64(i * 5)}
		_, err = client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc)
		require.NoError(t, err)
	}
	_, err = client.IndexRefresh(ctx, idx)
	require.NoError(t, err)

	// Use SearchWithRequest for a high-level search.Request
	req := search.NewRequest()
	matchAll := types.MatchAllQuery{}
	req.Query = &types.Query{MatchAll: &matchAll}
	size := 5
	req.Size = &size
	req.Source_ = true

	res, err := client.SearchWithRequest(ctx, alias, req)
	require.NoError(t, err)
	assert.Equal(t, int64(5), res.Hits.Total.Value)
}
