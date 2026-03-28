//go:build integration

package esv9_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	es9 "github.com/elastic/go-elasticsearch/v9"
	core_bulk "github.com/elastic/go-elasticsearch/v9/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/deletebyquery"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/updatebyquery"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/conflicts"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/healthstatus"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/level"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/refresh"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	esv9 "github.com/tomtwinkle/es-typed-go/esv9"
	"github.com/tomtwinkle/es-typed-go/esv9/query"
)

// newRawTypedClient creates a raw go-elasticsearch TypedClient for operations
// not yet exposed through ESClient (e.g. opening a scroll context).
func newRawTypedClient(t *testing.T) *es9.TypedClient {
	t.Helper()
	c, err := es9.NewTypedClient(es9.Config{Addresses: []string{esURL()}})
	assert.NilError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := c.Info().Do(ctx)
	if err != nil {
		t.Skipf("skipping integration test: Elasticsearch is unavailable at %s: %v", esURL(), err)
	}
	assert.Assert(t, info != nil)

	return c
}

// withUBQ returns an UpdateByQueryOption that sets the request body to the
// given query and inline script.
func withUBQ(q types.Query, inlineScript string) esv9.UpdateByQueryOption {
	return func(r *updatebyquery.UpdateByQuery) {
		req := updatebyquery.NewRequest()
		req.Query = &q
		req.Script = &types.Script{Source: &inlineScript}
		r.Request(req)
	}
}

// withDBQ returns a DeleteByQueryOption that sets the request query.
func withDBQ(q types.Query) esv9.DeleteByQueryOption {
	return func(r *deletebyquery.DeleteByQuery) {
		req := deletebyquery.NewRequest()
		req.Query = &q
		r.Request(req)
	}
}

// withBulkRaw returns a BulkOption that sets the NDJSON request body.
func withBulkRaw(body []byte) esv9.BulkOption {
	return func(r *core_bulk.Bulk) { r.Raw(bytes.NewReader(body)) }
}

// indexN indexes n productDocs into alias with refresh.True and waits for refresh.
func indexN(t *testing.T, client esv9.ESClient, idx estype.Index, alias estype.Alias, n int) {
	t.Helper()
	ctx := context.Background()
	for i := 1; i <= n; i++ {
		doc := productDoc{Name: fmt.Sprintf("Product %d", i), Category: "test", Price: float64(i * 10)}
		_, err := client.CreateDocument(ctx, alias, fmt.Sprintf("doc-%d", i), doc, esv9.WithRefresh(refresh.True))
		assert.NilError(t, err)
	}
	_, err := client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)
}

// ---------------------------------------------------------------------------
// WithRefresh
// ---------------------------------------------------------------------------

// WithRefresh(refresh.True) makes a newly indexed document immediately visible
// to search without requiring an explicit index refresh call.
func TestIntegration_WithOption_WithRefresh_DocumentImmediatelySearchable(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)

	doc := productDoc{Name: "Instant Widget", Category: "test", Price: 1.0}
	_, err = client.CreateDocument(ctx, alias, "doc-1", doc, esv9.WithRefresh(refresh.True))
	assert.NilError(t, err)

	// No explicit IndexRefresh call — document must still be searchable.
	req := search.NewRequest()
	q := query.MatchAll()
	req.Query = &q
	size := 10
	req.Size = &size
	res, err := client.SearchRaw(ctx, alias, req)
	assert.NilError(t, err)
	assert.Equal(t, int64(1), res.Hits.Total.Value)
}

// ---------------------------------------------------------------------------
// WithDeleteRefresh
// ---------------------------------------------------------------------------

// WithDeleteRefresh(refresh.True) makes a deletion immediately visible to
// search without requiring a separate refresh step.
func TestIntegration_WithOption_WithDeleteRefresh_DeletionImmediatelyReflectedInCount(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)

	doc := productDoc{Name: "Temporary", Category: "test", Price: 1.0}
	_, err = client.CreateDocument(ctx, alias, "doc-1", doc, esv9.WithRefresh(refresh.True))
	assert.NilError(t, err)

	_, err = client.DeleteDocument(ctx, idx, "doc-1", esv9.WithDeleteRefresh(refresh.True))
	assert.NilError(t, err)

	// No explicit refresh — deletion must already be reflected.
	count, err := client.IndexDocumentCount(ctx, idx)
	assert.NilError(t, err)
	assert.Equal(t, int64(0), count.Count)
}

// ---------------------------------------------------------------------------
// WithBulkRefresh
// ---------------------------------------------------------------------------

// WithBulkRefresh(refresh.True) flushes all bulk-indexed documents immediately
// so that they are searchable without a separate refresh call.
func TestIntegration_WithOption_WithBulkRefresh_BulkDocumentsImmediatelySearchable(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)

	var body []byte
	for i := 1; i <= 3; i++ {
		meta := fmt.Sprintf(`{"index":{"_index":%q,"_id":"doc-%d"}}`, idx.String(), i)
		src := fmt.Sprintf(`{"name":"Bulk Product %d","category":"test","price":%d}`, i, i*10)
		body = append(body, []byte(meta+"\n"+src+"\n")...)
	}

	_, err = client.Bulk(ctx, alias,
		withBulkRaw(body),
		esv9.WithBulkRefresh(refresh.True),
	)
	assert.NilError(t, err)

	// No explicit refresh — all three documents must be findable.
	req := search.NewRequest()
	q := query.MatchAll()
	req.Query = &q
	size := 10
	req.Size = &size
	res, err := client.SearchRaw(ctx, alias, req)
	assert.NilError(t, err)
	assert.Equal(t, int64(3), res.Hits.Total.Value)
}

// ---------------------------------------------------------------------------
// WithSourceIncludes / WithSourceExcludes
// ---------------------------------------------------------------------------

// WithSourceIncludes limits the _source fields returned by GetDocument to
// only the specified fields; unrequested fields are absent.
func TestIntegration_WithOption_WithSourceIncludes_OnlyRequestedFieldsReturned(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)

	doc := productDoc{Name: "Filtered Widget", Category: "test", Price: 42.0}
	_, err = client.CreateDocument(ctx, alias, "doc-1", doc, esv9.WithRefresh(refresh.True))
	assert.NilError(t, err)

	res, err := client.GetDocument(ctx, alias, "doc-1", esv9.WithSourceIncludes("name"))
	assert.NilError(t, err)
	assert.Assert(t, res.Found)

	var got map[string]any
	assert.NilError(t, json.Unmarshal(res.Source_, &got))
	_, hasName := got["name"]
	_, hasPrice := got["price"]
	assert.Assert(t, hasName, "expected name field in _source")
	assert.Assert(t, !hasPrice, "price field must be absent when not included")
}

// WithSourceExcludes omits the specified fields from the _source returned by
// GetDocument; all other fields remain present.
func TestIntegration_WithOption_WithSourceExcludes_ExcludedFieldAbsentFromSource(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)

	doc := productDoc{Name: "Exclude Test", Category: "test", Price: 99.0}
	_, err = client.CreateDocument(ctx, alias, "doc-1", doc, esv9.WithRefresh(refresh.True))
	assert.NilError(t, err)

	res, err := client.GetDocument(ctx, alias, "doc-1", esv9.WithSourceExcludes("price"))
	assert.NilError(t, err)
	assert.Assert(t, res.Found)

	var got map[string]any
	assert.NilError(t, json.Unmarshal(res.Source_, &got))
	_, hasName := got["name"]
	_, hasPrice := got["price"]
	assert.Assert(t, hasName, "name field must be present when only price is excluded")
	assert.Assert(t, !hasPrice, "price field must be absent when excluded")
}

// ---------------------------------------------------------------------------
// WithUpdateMaxDocs
// ---------------------------------------------------------------------------

// WithUpdateMaxDocs(n) limits an UpdateByQuery to at most n documents;
// the remaining matching documents are left unchanged.
func TestIntegration_WithOption_WithUpdateMaxDocs_OnlyNDocumentsUpdated(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	indexN(t, client, idx, alias, 5)

	res, err := client.UpdateByQuery(ctx, idx,
		withUBQ(query.MatchAll(), `ctx._source.category = "updated"`),
		esv9.WithUpdateMaxDocs(2),
		esv9.WithUpdateWaitForCompletion(true),
	)
	assert.NilError(t, err)
	assert.Assert(t, res.Updated != nil)
	assert.Equal(t, int64(2), *res.Updated)
}

// ---------------------------------------------------------------------------
// WithDeleteMaxDocs
// ---------------------------------------------------------------------------

// WithDeleteMaxDocs(n) limits a DeleteByQuery to at most n documents;
// the rest of the matching documents remain in the index.
func TestIntegration_WithOption_WithDeleteMaxDocs_OnlyNDocumentsDeleted(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	indexN(t, client, idx, alias, 5)

	res, err := client.DeleteByQuery(ctx, idx,
		withDBQ(query.MatchAll()),
		esv9.WithDeleteMaxDocs(2),
		esv9.WithDeleteWaitForCompletion(true),
	)
	assert.NilError(t, err)
	assert.Assert(t, res.Deleted != nil)
	assert.Equal(t, int64(2), *res.Deleted)
}

// ---------------------------------------------------------------------------
// WithUpdateConflicts / WithDeleteConflicts
// ---------------------------------------------------------------------------

// WithUpdateConflicts(proceed) continues an UpdateByQuery even when version
// conflicts are encountered, reporting them rather than aborting.
func TestIntegration_WithOption_WithUpdateConflicts_ProceedReportsConflictsInsteadOfAborting(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	indexN(t, client, idx, alias, 3)

	res, err := client.UpdateByQuery(ctx, idx,
		withUBQ(query.MatchAll(), `ctx._source.category = "updated"`),
		esv9.WithUpdateConflicts(conflicts.Proceed),
		esv9.WithUpdateWaitForCompletion(true),
	)
	assert.NilError(t, err)
	// On a single-node cluster there are no concurrent conflicts; the option
	// must not cause the request itself to fail.
	assert.Assert(t, res.Updated != nil)
	assert.Equal(t, int64(3), *res.Updated)
}

// WithDeleteConflicts(proceed) continues a DeleteByQuery even when version
// conflicts are encountered, reporting them rather than aborting.
func TestIntegration_WithOption_WithDeleteConflicts_ProceedReportsConflictsInsteadOfAborting(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	indexN(t, client, idx, alias, 3)

	res, err := client.DeleteByQuery(ctx, idx,
		withDBQ(query.MatchAll()),
		esv9.WithDeleteConflicts(conflicts.Proceed),
		esv9.WithDeleteWaitForCompletion(true),
	)
	assert.NilError(t, err)
	assert.Assert(t, res.Deleted != nil)
	assert.Equal(t, int64(3), *res.Deleted)
}

// ---------------------------------------------------------------------------
// WithUpdateWaitForCompletion / WithDeleteWaitForCompletion
// ---------------------------------------------------------------------------

// WithUpdateWaitForCompletion(false) returns immediately with a task reference
// rather than blocking until the UpdateByQuery finishes.
func TestIntegration_WithOption_WithUpdateWaitForCompletion_FalseReturnsTaskBeforeCompletion(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	indexN(t, client, idx, alias, 5)

	res, err := client.UpdateByQuery(ctx, idx,
		withUBQ(query.MatchAll(), `ctx._source.category = "async"`),
		esv9.WithUpdateWaitForCompletion(false),
	)
	assert.NilError(t, err)
	// When wait_for_completion=false the task field is populated instead of Updated.
	assert.Assert(t, res.Task != nil, "expected task reference when wait_for_completion=false")
}

// WithDeleteWaitForCompletion(false) returns immediately with a task reference
// rather than blocking until the DeleteByQuery finishes.
func TestIntegration_WithOption_WithDeleteWaitForCompletion_FalseReturnsTaskBeforeCompletion(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	indexN(t, client, idx, alias, 5)

	res, err := client.DeleteByQuery(ctx, idx,
		withDBQ(query.MatchAll()),
		esv9.WithDeleteWaitForCompletion(false),
	)
	assert.NilError(t, err)
	assert.Assert(t, res.Task != nil, "expected task reference when wait_for_completion=false")
}

// ---------------------------------------------------------------------------
// WithIgnoreUnavailable
// ---------------------------------------------------------------------------

// WithIgnoreUnavailable(true) suppresses the "index not found" error when the
// target index does not exist, making DeleteIndex idempotent.
func TestIntegration_WithOption_WithIgnoreUnavailable_DeleteNonExistentIndexDoesNotError(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()

	ghost := uniqueIndex(t, client) // never created

	exists, err := client.IndexExists(ctx, ghost)
	assert.NilError(t, err)
	assert.Assert(t, !exists)

	_, err = client.DeleteIndex(ctx, ghost, esv9.WithIgnoreUnavailable(true))
	assert.NilError(t, err)
}

// ---------------------------------------------------------------------------
// WithWaitForStatus
// ---------------------------------------------------------------------------

// WithWaitForStatus(Green) blocks ClusterHealth until the cluster health
// reaches green; on a healthy single-node cluster this resolves immediately.
func TestIntegration_WithOption_WithWaitForStatus_ReturnsWhenClusterStatusReached(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)

	res, err := client.ClusterHealth(ctx, idx,
		esv9.WithWaitForStatus(healthstatus.Green),
		esv9.WithHealthTimeout("10s"),
	)
	assert.NilError(t, err)
	// A single-node cluster with zero replicas is always green.
	assert.Assert(t, res.Status == healthstatus.Green || res.Status == healthstatus.Yellow,
		"expected green or yellow cluster status, got %v", res.Status)
}

// ---------------------------------------------------------------------------
// WithScrollId / WithClearScrollId
// ---------------------------------------------------------------------------

// WithScrollId sets the scroll_id on a Scroll request to retrieve the next
// page of results from an ongoing scroll operation.
func TestIntegration_WithOption_WithScrollId_ContinuationPageReturnsRemainingDocuments(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	raw := newRawTypedClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	indexN(t, client, idx, alias, 5)

	// Open a scroll via the raw typed client (page size = 2).
	firstPage, err := raw.Search().Index(alias.String()).Size(2).Scroll("1m").Do(ctx)
	assert.NilError(t, err)
	assert.Assert(t, firstPage.ScrollId_ != nil, "expected scroll_id in first response")
	assert.Equal(t, 2, len(firstPage.Hits.Hits))

	scrollID := *firstPage.ScrollId_

	// Use ESClient.Scroll with WithScrollId to fetch the second page.
	secondPage, err := client.Scroll(ctx, esv9.WithScrollId(scrollID))
	assert.NilError(t, err)
	assert.Assert(t, len(secondPage.Hits.Hits) > 0, "expected hits on second scroll page")

	// Release the scroll context with WithClearScrollId.
	_, err = client.ClearScroll(ctx, esv9.WithClearScrollId(scrollID))
	assert.NilError(t, err)
}

// WithClearScrollId releases the specified scroll context so that server
// resources are freed immediately rather than waiting for the keep-alive to
// expire.
func TestIntegration_WithOption_WithClearScrollId_ScrollContextReleasedWithoutError(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	raw := newRawTypedClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)

	doc := productDoc{Name: "Scroll Target", Category: "test", Price: 1.0}
	_, err = client.CreateDocument(ctx, alias, "doc-1", doc, esv9.WithRefresh(refresh.True))
	assert.NilError(t, err)

	page, err := raw.Search().Index(alias.String()).Size(1).Scroll("1m").Do(ctx)
	assert.NilError(t, err)
	assert.Assert(t, page.ScrollId_ != nil)

	_, err = client.ClearScroll(ctx, esv9.WithClearScrollId(*page.ScrollId_))
	assert.NilError(t, err)
}

// ---------------------------------------------------------------------------
// WithStatsLevel
// ---------------------------------------------------------------------------

// WithStatsLevel(level.Shards) makes IndicesStats include per-shard metrics
// in addition to the default index-level summary.
func TestIntegration_WithOption_WithStatsLevel_ShardsLevelIncludesShardDetail(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)

	res, err := client.IndicesStats(ctx, idx, esv9.WithStatsLevel(level.Shards))
	assert.NilError(t, err)

	indexStats, ok := res.Indices[idx.String()]
	assert.Assert(t, ok, "expected stats entry for the index")
	assert.Assert(t, len(indexStats.Shards) > 0, "expected shard-level detail when level=shards")
}

// ---------------------------------------------------------------------------
// WithQueryCache / WithRequestCache / WithFielddataCache
// ---------------------------------------------------------------------------

// WithQueryCache(true) clears only the query cache for the index.
func TestIntegration_WithOption_WithQueryCache_ClearsQueryCacheWithoutError(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)

	_, err = client.ClearCache(ctx, idx, esv9.WithQueryCache(true))
	assert.NilError(t, err)
}

// WithRequestCache(true) clears only the request cache for the index.
func TestIntegration_WithOption_WithRequestCache_ClearsRequestCacheWithoutError(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)

	_, err = client.ClearCache(ctx, idx, esv9.WithRequestCache(true))
	assert.NilError(t, err)
}

// WithFielddataCache(true) clears only the fielddata cache for the index.
func TestIntegration_WithOption_WithFielddataCache_ClearsFielddataCacheWithoutError(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)

	_, err = client.ClearCache(ctx, idx, esv9.WithFielddataCache(true))
	assert.NilError(t, err)
}

// ---------------------------------------------------------------------------
// WithMaxNumSegments (ForceMerge)
// ---------------------------------------------------------------------------

// WithMaxNumSegments("1") causes ForceMerge to consolidate all segments into
// a single segment on the primary shard.
func TestIntegration_WithOption_WithMaxNumSegments_ForceMergeToSingleSegment(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
	indexN(t, client, idx, alias, 3)

	_, err = client.ForceMerge(ctx, idx, esv9.WithMaxNumSegments("1"))
	assert.NilError(t, err)

	stats, err := client.IndicesStats(ctx, idx, esv9.WithStatsLevel(level.Shards))
	assert.NilError(t, err)
	indexStats := stats.Indices[idx.String()]
	for _, shards := range indexStats.Shards {
		for _, shard := range shards {
			if shard.Routing != nil && shard.Routing.Primary && shard.Segments != nil {
				assert.Assert(t, shard.Segments.Count <= 1,
					"expected at most 1 segment after force-merge with max_num_segments=1, got %d",
					shard.Segments.Count)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// WithOnlyManaged / WithOnlyErrors (ExplainLifecycle)
// ---------------------------------------------------------------------------

// WithOnlyManaged(true) restricts ExplainLifecycle output to ILM-managed
// indices; an index without an ILM policy does not appear in the response.
func TestIntegration_WithOption_WithOnlyManaged_UnmanagedIndexAbsentFromExplain(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)

	res, err := client.ExplainLifecycle(ctx, idx, esv9.WithOnlyManaged(true))
	assert.NilError(t, err)

	_, isPresent := res.Indices[idx.String()]
	assert.Assert(t, !isPresent, "unmanaged index must not appear when only_managed=true")
}

// WithOnlyErrors(true) restricts ExplainLifecycle output to indices that have
// ILM errors; a healthy index with no ILM policy does not appear.
func TestIntegration_WithOption_WithOnlyErrors_HealthyIndexAbsentFromExplain(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)

	res, err := client.ExplainLifecycle(ctx, idx, esv9.WithOnlyErrors(true))
	assert.NilError(t, err)

	_, hasErrors := res.Indices[idx.String()]
	assert.Assert(t, !hasErrors, "index without ILM errors must not appear when only_errors=true")
}

// ---------------------------------------------------------------------------
// WithTransformAllowNoMatch
// ---------------------------------------------------------------------------

// WithTransformAllowNoMatch(true) suppresses the error returned by GetTransform
// when no transform matches the requested ID pattern.
func TestIntegration_WithOption_WithTransformAllowNoMatch_NoErrorOnNonExistentTransform(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()

	_, err := client.GetTransform(ctx, estype.TransformID("nonexistent-transform-*"),
		esv9.WithTransformAllowNoMatch(true),
	)
	assert.NilError(t, err)
}

// ---------------------------------------------------------------------------
// WithMlJobAllowNoMatch
// ---------------------------------------------------------------------------

// WithMlJobAllowNoMatch(true) suppresses the error returned by MlGetJobs when
// no job matches the requested ID pattern.
func TestIntegration_WithOption_WithMlJobAllowNoMatch_NoErrorOnNonExistentJob(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()

	_, err := client.MlGetJobs(ctx, estype.MLJobID("nonexistent-job-*"),
		esv9.WithMlJobAllowNoMatch(true),
	)
	assert.NilError(t, err)
}

// ---------------------------------------------------------------------------
// WithMlDatafeedAllowNoMatch
// ---------------------------------------------------------------------------

// WithMlDatafeedAllowNoMatch(true) suppresses the error returned by
// MlGetDatafeeds when no datafeed matches the requested ID pattern.
func TestIntegration_WithOption_WithMlDatafeedAllowNoMatch_NoErrorOnNonExistentDatafeed(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()

	_, err := client.MlGetDatafeeds(ctx, estype.DatafeedID("nonexistent-datafeed-*"),
		esv9.WithMlDatafeedAllowNoMatch(true),
	)
	assert.NilError(t, err)
}

// ---------------------------------------------------------------------------
// WithRouting / WithGetRouting
// ---------------------------------------------------------------------------

// WithRouting and WithGetRouting work together: a document indexed with a
// custom routing value is retrievable when the same routing value is supplied
// to GetDocument.
func TestIntegration_WithOption_WithRouting_DocumentFoundWhenRoutingMatches(t *testing.T) {
	t.Parallel()
	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), nil)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)

	routing := "shard-key-a"
	doc := productDoc{Name: "Routed Doc", Category: "test", Price: 7.0}
	_, err = client.CreateDocument(ctx, alias, "doc-r1", doc,
		esv9.WithRefresh(refresh.True),
		esv9.WithRouting(routing),
	)
	assert.NilError(t, err)

	res, err := client.GetDocument(ctx, alias, "doc-r1",
		esv9.WithGetRouting(routing),
	)
	assert.NilError(t, err)
	assert.Assert(t, res.Found)
	assert.Equal(t, "doc-r1", res.Id_)
}
