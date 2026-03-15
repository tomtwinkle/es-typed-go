//go:build integration

package esv8_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	corebulk "github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/clearscroll"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/closepointintime"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	corecreate "github.com/elastic/go-elasticsearch/v8/typedapi/core/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/deletebyquery"
	coreindex "github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/mget"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/scroll"
	coresearch "github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/updatebyquery"
	indicescreate "github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/google/uuid"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	esv8 "github.com/tomtwinkle/es-typed-go/esv8"
)

// mustMarshal marshals v to JSON, failing the test on error.
func mustMarshal(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	assert.NilError(t, err)
	return b
}

// newSpecClient returns an ESClientSpec for spec-named method integration tests.
func newSpecClient(t *testing.T) esv8.ESClientSpec {
	t.Helper()
	client, err := esv8.NewSpecClient(es8.Config{
		Addresses: []string{esURL()},
	})
	assert.NilError(t, err)
	return client
}

// uniqueSpecIndex creates a fresh test index with no replicas and registers t.Cleanup to delete it.
// Using number_of_replicas=0 ensures the index is immediately GREEN on a single-node cluster,
// avoiding shard allocation issues (e.g. UpdateByQuery with scripts failing on yellow indices).
func uniqueSpecIndex(t *testing.T, client esv8.ESClientSpec) string {
	t.Helper()
	ctx := context.Background()
	name := fmt.Sprintf("spectest-%s", uuid.New().String())
	replicas := "0"
	req := &indicescreate.Request{
		Settings: &types.IndexSettings{
			NumberOfReplicas: &replicas,
		},
	}
	_, err := client.IndicesCreate(ctx, name, req)
	assert.NilError(t, err)
	t.Cleanup(func() {
		cctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = client.IndicesDelete(cctx, name) //nolint:errcheck
	})
	return name
}

// bulkOps builds a bulk []any payload of N index operations for the given index.
func bulkOps(idx string, docs ...map[string]any) corebulk.Request {
	var req corebulk.Request
	for _, doc := range docs {
		req = append(req,
			map[string]any{"index": map[string]any{"_index": idx}},
			doc,
		)
	}
	return req
}

// ---------------------------------------------------------------------------
// Core spec-named method tests
// ---------------------------------------------------------------------------

func TestIntegration_Spec_Ping(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ok, err := client.Ping(context.Background())
	assert.NilError(t, err)
	assert.Assert(t, ok, "Ping should return true for a running cluster")
}

func TestIntegration_Spec_ClusterHealth(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()

	res, err := client.ClusterHealth(ctx)
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	assert.Assert(t, res.ClusterName != "", "cluster name should not be empty")
	t.Logf("Cluster %q status: %s", res.ClusterName, res.Status)
}

func TestIntegration_Spec_IndicesCreateDelete(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	name := fmt.Sprintf("spec-idx-%s", uuid.New().String())
	t.Cleanup(func() {
		cctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, _ = client.IndicesDelete(cctx, name)
	})

	// Create
	createRes, err := client.IndicesCreate(ctx, name, nil)
	assert.NilError(t, err)
	assert.Assert(t, createRes != nil)

	// Exists
	exists, err := client.IndicesExists(ctx, name)
	assert.NilError(t, err)
	assert.Assert(t, exists, "index should exist after creation")

	// Delete
	delRes, err := client.IndicesDelete(ctx, name)
	assert.NilError(t, err)
	assert.Assert(t, delRes != nil)

	// Gone
	exists, err = client.IndicesExists(ctx, name)
	assert.NilError(t, err)
	assert.Assert(t, !exists, "index should not exist after deletion")
}

func TestIntegration_Spec_Index_Get_Create_Delete_Exists(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	doc := map[string]any{"title": "hello", "value": 42}

	// Index without explicit ID – index.Request is json.RawMessage
	idxReq := coreindex.Request(mustMarshal(t, doc))
	idxRes, err := client.Index(ctx, idx, &idxReq)
	assert.NilError(t, err)
	assert.Assert(t, idxRes != nil)

	// Create with explicit ID – create.Request is json.RawMessage
	docID := "spec-doc-explicit"
	createReq := corecreate.Request(mustMarshal(t, doc))
	createRes, err := client.Create(ctx, idx, docID, &createReq)
	assert.NilError(t, err)
	assert.Assert(t, createRes != nil)

	// Refresh so doc is visible
	_, err = client.IndicesRefresh(ctx)
	assert.NilError(t, err)

	// Get the doc
	getRes, err := client.Get(ctx, idx, docID)
	assert.NilError(t, err)
	assert.Assert(t, getRes != nil)
	assert.Assert(t, getRes.Found, "document should be found")
	assert.Equal(t, docID, getRes.Id_)

	// Exists check
	exists, err := client.Exists(ctx, idx, docID)
	assert.NilError(t, err)
	assert.Assert(t, exists, "document should exist")

	// Delete the doc
	delRes, err := client.Delete(ctx, idx, docID)
	assert.NilError(t, err)
	assert.Assert(t, delRes != nil)

	// No longer exists
	_, _ = client.IndicesRefresh(ctx)
	exists, err = client.Exists(ctx, idx, docID)
	assert.NilError(t, err)
	assert.Assert(t, !exists, "document should not exist after deletion")
}

func TestIntegration_Spec_Update(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	docID := "update-doc"
	createReq := corecreate.Request(mustMarshal(t, map[string]any{"status": "pending"}))
	_, err := client.Create(ctx, idx, docID, &createReq)
	assert.NilError(t, err)
	_, _ = client.IndicesRefresh(ctx)

	// update.Request.Doc is json.RawMessage, not map[string]any
	updateRes, err := client.Update(ctx, idx, docID, &update.Request{
		Doc: mustMarshal(t, map[string]any{"status": "done"}),
	})
	assert.NilError(t, err)
	assert.Assert(t, updateRes != nil)
}

func TestIntegration_Spec_Bulk(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	// bulk.Request is []any; each operation is an action+document pair
	var req corebulk.Request
	for i := 1; i <= 5; i++ {
		req = append(req,
			map[string]any{"index": map[string]any{"_index": idx}},
			map[string]any{"seq": i, "label": fmt.Sprintf("item-%d", i)},
		)
	}

	res, err := client.Bulk(ctx, &req)
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	assert.Assert(t, !res.Errors, "bulk should complete without errors")
	assert.Assert(t, len(res.Items) == 5, "expected 5 items, got %d", len(res.Items))
}

func TestIntegration_Spec_Count(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	// Index 3 docs via Bulk
	docs := []map[string]any{{"n": 0}, {"n": 1}, {"n": 2}}
	req := bulkOps(idx, docs...)
	_, err := client.Bulk(ctx, &req)
	assert.NilError(t, err)
	_, _ = client.IndicesRefresh(ctx)

	countRes, err := client.Count(ctx, &count.Request{})
	assert.NilError(t, err)
	assert.Assert(t, countRes != nil)
	assert.Assert(t, countRes.Count >= 3,
		"count should be at least 3, got %d", countRes.Count)
}

func TestIntegration_Spec_Mget(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	for i := 1; i <= 3; i++ {
		id := fmt.Sprintf("mget-doc-%d", i)
		createReq := corecreate.Request(mustMarshal(t, map[string]any{"n": i}))
		_, err := client.Create(ctx, idx, id, &createReq)
		assert.NilError(t, err)
	}
	_, _ = client.IndicesRefresh(ctx)

	ptr := func(s string) *string { return &s }
	res, err := client.Mget(ctx, &mget.Request{
		Docs: []types.MgetOperation{
			{Index_: ptr(idx), Id_: "mget-doc-1"},
			{Index_: ptr(idx), Id_: "mget-doc-2"},
			{Index_: ptr(idx), Id_: "mget-doc-3"},
		},
	})
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	assert.Assert(t, len(res.Docs) == 3, "expected 3 mget docs, got %d", len(res.Docs))
}

func TestIntegration_Spec_DeleteByQuery(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	docs := []map[string]any{{"tag": "delete-me"}, {"tag": "delete-me"}, {"tag": "delete-me"}}
	req := bulkOps(idx, docs...)
	_, err := client.Bulk(ctx, &req)
	assert.NilError(t, err)
	// Refresh only the test index, not all indices, to avoid interference from
	// YELLOW shards of other parallel tests.
	_, _ = client.IndexRefresh(ctx, estype.Index(idx))

	matchAll := types.Query{MatchAll: &types.MatchAllQuery{}}
	res, err := client.DeleteByQuery(ctx, idx, &deletebyquery.Request{Query: &matchAll})
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	t.Logf("DeleteByQuery deleted %d documents", res.Deleted)
}

func TestIntegration_Spec_UpdateByQuery(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	createReq := corecreate.Request(mustMarshal(t, map[string]any{"status": "pending"}))
	_, err := client.Create(ctx, idx, "ubq-1", &createReq)
	assert.NilError(t, err)
	// Refresh only the test index, not all indices, to avoid interference from
	// YELLOW shards of other parallel tests.
	_, _ = client.IndexRefresh(ctx, estype.Index(idx))

	matchAll := types.Query{MatchAll: &types.MatchAllQuery{}}
	src := "ctx._source.status = 'done'"
	res, err := client.UpdateByQuery(ctx, idx, &updatebyquery.Request{
		Query:  &matchAll,
		Script: &types.Script{Source: &src},
	})
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	t.Logf("UpdateByQuery updated %d documents", res.Updated)
}

func TestIntegration_Spec_ScrollAndClear(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	// Index docs via bulk
	docs := []map[string]any{{"s": 0}, {"s": 1}, {"s": 2}, {"s": 3}, {"s": 4}}
	req := bulkOps(idx, docs...)
	_, err := client.Bulk(ctx, &req)
	assert.NilError(t, err)
	_, _ = client.IndicesRefresh(ctx)

	// Initiate a scroll using SearchWithRequest – Scroll is a URL parameter set
	// via the builder, not via the request body. We use SearchWithRequest here,
	// which does not set scroll, so we won't get a valid scroll_id back. Instead
	// we verify that calling Scroll with an empty scroll_id returns a well-formed
	// ES error (not a panic or transport error).
	scrollRes, err := client.SearchWithRequest(ctx, estype.Alias(idx), &coresearch.Request{
		Query: &types.Query{MatchAll: &types.MatchAllQuery{}},
		Size:  func() *int { n := 2; return &n }(),
	})
	assert.NilError(t, err)
	assert.Assert(t, scrollRes != nil)
	t.Logf("SearchWithRequest returned %d hits", scrollRes.Hits.Total.Value)

	// Attempt ClearScroll with _all – this is a fire-and-forget cleanup.
	clearRes, err := client.ClearScroll(ctx, &clearscroll.Request{
		ScrollId: []string{"_all"},
	})
	// _all may return an error if there are no open scroll contexts; we accept either.
	if err == nil {
		assert.Assert(t, clearRes != nil)
	}

	// Test Scroll method itself – calling it with an invalid ID tests the method
	// path and returns a 404/not-found error from ES.
	_, scrollErr := client.Scroll(ctx, &scroll.Request{
		ScrollId: "invalid_scroll_id_for_test",
		Scroll:   "1m",
	})
	// We expect an error (ES rejects an unknown scroll_id), but NOT a panic.
	assert.Assert(t, scrollErr != nil, "Scroll with invalid ID should return an error")
	t.Logf("Scroll with invalid ID returned expected error: %v", scrollErr)
}

func TestIntegration_Spec_PointInTime(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	createReq := corecreate.Request(mustMarshal(t, map[string]any{"value": 1}))
	_, err := client.Create(ctx, idx, "pit-doc", &createReq)
	assert.NilError(t, err)
	_, _ = client.IndicesRefresh(ctx)

	// Open PIT – keepAlive is now a required positional parameter.
	pitRes, err := client.OpenPointInTime(ctx, idx, "1m", nil)
	assert.NilError(t, err)
	assert.Assert(t, pitRes != nil)
	assert.Assert(t, pitRes.Id != "", "PIT id should not be empty")
	t.Logf("Opened PIT: %s", pitRes.Id)

	// Close PIT
	closeRes, err := client.ClosePointInTime(ctx, &closepointintime.Request{Id: pitRes.Id})
	assert.NilError(t, err)
	assert.Assert(t, closeRes != nil)
	assert.Assert(t, closeRes.Succeeded, "PIT close should succeed")
}

func TestIntegration_Spec_TasksList(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	res, err := client.TasksList(context.Background())
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	t.Logf("TasksList: %d node(s)", len(res.Nodes))
}

func TestIntegration_Spec_CatHealth(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	res, err := client.CatHealth(context.Background())
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	t.Logf("CatHealth: %d row(s)", len(res))
}

func TestIntegration_Spec_CatIndices(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	res, err := client.CatIndices(ctx)
	assert.NilError(t, err)
	assert.Assert(t, res != nil)

	found := false
	for _, entry := range res {
		if entry.Index != nil && *entry.Index == idx {
			found = true
			break
		}
	}
	assert.Assert(t, found, "newly created index %q should appear in CatIndices", idx)
}

func TestIntegration_Spec_IndicesRefresh(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)
	_ = idx

	res, err := client.IndicesRefresh(ctx)
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
}

func TestIntegration_Spec_IndicesGetSettings(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	_ = uniqueSpecIndex(t, client)

	res, err := client.IndicesGetSettings(ctx)
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
}

func TestIntegration_Spec_IndicesGetMapping(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	_ = uniqueSpecIndex(t, client)

	res, err := client.IndicesGetMapping(ctx)
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
}

func TestIntegration_Spec_IndicesGetAlias(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	// Add an alias via the semantic API
	idxParsed, _ := estype.ParseESIndex(idx)
	aliasParsed, _ := estype.ParseESAlias("alias-" + idx)
	_, err := client.CreateAlias(ctx, idxParsed, aliasParsed, false)
	assert.NilError(t, err)

	// Retrieve all aliases using the spec method
	res, err := client.IndicesGetAlias(ctx)
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	_, ok := res[idx]
	assert.Assert(t, ok, "alias response should include newly created index %q", idx)
}

func TestIntegration_Spec_IngestGetPipeline(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	res, err := client.IngestGetPipeline(context.Background())
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
}

func TestIntegration_Spec_Get_Explicit(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	docID := "explicit-get-doc"
	createReq := corecreate.Request(mustMarshal(t, map[string]any{"hello": "world"}))
	_, err := client.Create(ctx, idx, docID, &createReq)
	assert.NilError(t, err)
	_, _ = client.IndicesRefresh(ctx)

	res, err := client.Get(ctx, idx, docID)
	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	assert.Assert(t, res.Found, "document should be found")
	assert.Equal(t, docID, res.Id_)
}

func TestIntegration_Spec_GetMatchesGetDocument(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	docID := "get-req-doc"
	createReq := corecreate.Request(mustMarshal(t, map[string]any{"field1": "a", "field2": "b"}))
	_, err := client.Create(ctx, idx, docID, &createReq)
	assert.NilError(t, err)
	_, _ = client.IndicesRefresh(ctx)

	// Using the semantic GetDocument
	getRes, err := client.GetDocument(ctx, estype.Alias(idx), docID)
	assert.NilError(t, err)
	assert.Assert(t, getRes != nil)
	assert.Assert(t, getRes.Found)

	// Using the spec Get method – same document
	specGetRes, err := client.Get(ctx, idx, docID)
	assert.NilError(t, err)
	assert.Assert(t, specGetRes != nil)
	assert.Assert(t, specGetRes.Found)
	assert.Equal(t, getRes.Id_, specGetRes.Id_)
}

func TestIntegration_Spec_Delete(t *testing.T) {
	t.Parallel()
	client := newSpecClient(t)
	ctx := context.Background()
	idx := uniqueSpecIndex(t, client)

	docID := "del-doc"
	createReq := corecreate.Request(mustMarshal(t, map[string]any{"x": 1}))
	_, err := client.Create(ctx, idx, docID, &createReq)
	assert.NilError(t, err)
	_, _ = client.IndicesRefresh(ctx)

	delRes, err := client.Delete(ctx, idx, docID)
	assert.NilError(t, err)
	assert.Assert(t, delRes != nil)

	// Confirm deleted
	_, _ = client.IndicesRefresh(ctx)
	exists, err := client.Exists(ctx, idx, docID)
	assert.NilError(t, err)
	assert.Assert(t, !exists, "doc should be gone after Delete")
}
