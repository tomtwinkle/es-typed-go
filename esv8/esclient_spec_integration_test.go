//go:build integration

package esv8_test

import (
"context"
"fmt"
"testing"
"time"

es8 "github.com/elastic/go-elasticsearch/v8"
corebulk "github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/clearscroll"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/closepointintime"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/create"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/deletebyquery"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/mget"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/openpointintime"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/scroll"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
"github.com/elastic/go-elasticsearch/v8/typedapi/core/updatebyquery"
"github.com/elastic/go-elasticsearch/v8/typedapi/types"
"gotest.tools/v3/assert"

esv8 "github.com/tomtwinkle/es-typed-go/esv8"
"github.com/tomtwinkle/es-typed-go/estype"
)

// newSpecClient returns an ESClientSpec for spec-named method integration tests.
func newSpecClient(t *testing.T) esv8.ESClientSpec {
t.Helper()
client, err := esv8.NewSpecClient(es8.Config{
Addresses: []string{esURL()},
})
assert.NilError(t, err)
return client
}

// uniqueSpecIndex creates a fresh test index and registers t.Cleanup to delete it.
func uniqueSpecIndex(t *testing.T, client esv8.ESClientSpec) string {
t.Helper()
ctx := context.Background()
name := fmt.Sprintf("spectest-%d", time.Now().UnixNano())
_, err := client.IndicesCreate(ctx, name, nil)
assert.NilError(t, err)
t.Cleanup(func() {
cctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
_, _ = client.IndicesDelete(cctx, name) //nolint:errcheck
})
return name
}

// ---------------------------------------------------------------------------
// Core spec-named method tests
// ---------------------------------------------------------------------------

func TestIntegration_Spec_Ping(t *testing.T) {
client := newSpecClient(t)
ok, err := client.Ping(context.Background())
assert.NilError(t, err)
assert.Assert(t, ok, "Ping should return true for a running cluster")
}

func TestIntegration_Spec_ClusterHealth(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()

res, err := client.ClusterHealth(ctx)
assert.NilError(t, err)
assert.Assert(t, res != nil)
assert.Assert(t, res.ClusterName != "", "cluster name should not be empty")
t.Logf("Cluster %q status: %s", res.ClusterName, res.Status)
}

func TestIntegration_Spec_IndicesCreateDelete(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
name := fmt.Sprintf("spec-idx-%d", time.Now().UnixNano())
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
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

doc := map[string]any{"title": "hello", "value": 42}

// Index without explicit ID
idxRes, err := client.Index(ctx, idx, &index.Request{Document: doc})
assert.NilError(t, err)
assert.Assert(t, idxRes != nil)

// Create with explicit ID
docID := "spec-doc-explicit"
createRes, err := client.Create(ctx, idx, docID, &create.Request{Document: doc})
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
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

docID := "update-doc"
_, err := client.Create(ctx, idx, docID, &create.Request{
Document: map[string]any{"status": "pending"},
})
assert.NilError(t, err)
_, _ = client.IndicesRefresh(ctx)

updateRes, err := client.Update(ctx, idx, docID, &update.Request{
Doc: map[string]any{"status": "done"},
})
assert.NilError(t, err)
assert.Assert(t, updateRes != nil)
}

func TestIntegration_Spec_Bulk(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

req := &corebulk.Request{}
for i := 1; i <= 5; i++ {
doc := map[string]any{"seq": i, "label": fmt.Sprintf("item-%d", i)}
req.Operations = append(req.Operations,
types.BulkOperationContainer{Index: &types.BulkIndexOperation{Index_: &idx}},
doc)
}

res, err := client.Bulk(ctx, req)
assert.NilError(t, err)
assert.Assert(t, res != nil)
assert.Assert(t, !res.Errors, "bulk should complete without errors")
assert.Assert(t, len(res.Items) == 5, "expected 5 items, got %d", len(res.Items))
}

func TestIntegration_Spec_Count(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

// Index some docs
req := &corebulk.Request{}
for i := 0; i < 3; i++ {
doc := map[string]any{"n": i}
req.Operations = append(req.Operations,
types.BulkOperationContainer{Index: &types.BulkIndexOperation{Index_: &idx}},
doc)
}
_, err := client.Bulk(ctx, req)
assert.NilError(t, err)
_, _ = client.IndicesRefresh(ctx)

countRes, err := client.Count(ctx, &count.Request{})
assert.NilError(t, err)
assert.Assert(t, countRes != nil)
assert.Assert(t, countRes.Count >= 3,
"count should be at least 3, got %d", countRes.Count)
}

func TestIntegration_Spec_Mget(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

for i := 1; i <= 3; i++ {
id := fmt.Sprintf("mget-doc-%d", i)
_, err := client.Create(ctx, idx, id, &create.Request{
Document: map[string]any{"n": i},
})
assert.NilError(t, err)
}
_, _ = client.IndicesRefresh(ctx)

ptr := func(s string) *string { return &s }
res, err := client.Mget(ctx, &mget.Request{
Docs: []types.MgetOperation{
{Index_: &idx, Id_: ptr("mget-doc-1")},
{Index_: &idx, Id_: ptr("mget-doc-2")},
{Index_: &idx, Id_: ptr("mget-doc-3")},
},
})
assert.NilError(t, err)
assert.Assert(t, res != nil)
assert.Assert(t, len(res.Docs) == 3, "expected 3 mget docs, got %d", len(res.Docs))
}

func TestIntegration_Spec_DeleteByQuery(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

req := &corebulk.Request{}
for i := 0; i < 5; i++ {
doc := map[string]any{"tag": "delete-me"}
req.Operations = append(req.Operations,
types.BulkOperationContainer{Index: &types.BulkIndexOperation{Index_: &idx}},
doc)
}
_, err := client.Bulk(ctx, req)
assert.NilError(t, err)
_, _ = client.IndicesRefresh(ctx)

matchAll := types.Query{MatchAll: &types.MatchAllQuery{}}
res, err := client.DeleteByQuery(ctx, idx, &deletebyquery.Request{Query: &matchAll})
assert.NilError(t, err)
assert.Assert(t, res != nil)
t.Logf("DeleteByQuery deleted %d documents", res.Deleted)
}

func TestIntegration_Spec_UpdateByQuery(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

_, err := client.Create(ctx, idx, "ubq-1", &create.Request{
Document: map[string]any{"status": "pending"},
})
assert.NilError(t, err)
_, _ = client.IndicesRefresh(ctx)

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
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

req := &corebulk.Request{}
for i := 0; i < 5; i++ {
doc := map[string]any{"seq": i}
req.Operations = append(req.Operations,
types.BulkOperationContainer{Index: &types.BulkIndexOperation{Index_: &idx}},
doc)
}
_, err := client.Bulk(ctx, req)
assert.NilError(t, err)
_, _ = client.IndicesRefresh(ctx)

// Start a scroll via the semantic Search method, with a scroll window
scrollRes, err := client.SearchWithRequest(ctx, estype.Alias(idx), &search.Request{
Query: &types.Query{MatchAll: &types.MatchAllQuery{}},
Size:  func() *int { n := 2; return &n }(),
Scroll: func() *types.Duration { d := types.Duration("1m"); return &d }(),
})
assert.NilError(t, err)
assert.Assert(t, scrollRes != nil)

if scrollRes.ScrollId != nil {
sid := *scrollRes.ScrollId

// Continue scroll using the spec method
scrollContinue, err := client.Scroll(ctx, &scroll.Request{
ScrollId: &sid,
Scroll:   "1m",
})
assert.NilError(t, err)
assert.Assert(t, scrollContinue != nil)

// Clear the scroll context
clearRes, err := client.ClearScroll(ctx, &clearscroll.Request{
ScrollId: []string{sid},
})
assert.NilError(t, err)
assert.Assert(t, clearRes != nil)
assert.Assert(t, clearRes.Succeeded, "ClearScroll should succeed")
}
}

func TestIntegration_Spec_PointInTime(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

_, err := client.Create(ctx, idx, "pit-doc", &create.Request{
Document: map[string]any{"value": 1},
})
assert.NilError(t, err)
_, _ = client.IndicesRefresh(ctx)

// Open PIT
pitRes, err := client.OpenPointInTime(ctx, idx, &openpointintime.Request{KeepAlive: "1m"})
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
client := newSpecClient(t)
res, err := client.TasksList(context.Background())
assert.NilError(t, err)
assert.Assert(t, res != nil)
t.Logf("TasksList: %d node(s)", len(res.Nodes))
}

func TestIntegration_Spec_CatHealth(t *testing.T) {
client := newSpecClient(t)
res, err := client.CatHealth(context.Background())
assert.NilError(t, err)
assert.Assert(t, res != nil)
t.Logf("CatHealth: %d row(s)", len(res))
}

func TestIntegration_Spec_CatIndices(t *testing.T) {
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
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)
_ = idx

// Global refresh (spec: indices.refresh with no index targets all indices)
res, err := client.IndicesRefresh(ctx)
assert.NilError(t, err)
assert.Assert(t, res != nil)
}

func TestIntegration_Spec_IndicesGetSettings(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
_ = uniqueSpecIndex(t, client)

// GetSettings returns settings for all indices
res, err := client.IndicesGetSettings(ctx)
assert.NilError(t, err)
assert.Assert(t, res != nil)
}

func TestIntegration_Spec_IndicesGetMapping(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
_ = uniqueSpecIndex(t, client)

res, err := client.IndicesGetMapping(ctx)
assert.NilError(t, err)
assert.Assert(t, res != nil)
}

func TestIntegration_Spec_IndicesGetAlias(t *testing.T) {
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
client := newSpecClient(t)
res, err := client.IngestGetPipeline(context.Background())
assert.NilError(t, err)
assert.Assert(t, res != nil)
}

func TestIntegration_Spec_Get_Explicit(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

docID := "explicit-get-doc"
_, err := client.Create(ctx, idx, docID, &create.Request{
Document: map[string]any{"hello": "world"},
})
assert.NilError(t, err)
_, _ = client.IndicesRefresh(ctx)

// Spec-named Get (no request struct needed)
res, err := client.Get(ctx, idx, docID)
assert.NilError(t, err)
assert.Assert(t, res != nil)
assert.Assert(t, res.Found, "document should be found")
assert.Equal(t, docID, res.Id_)
}

func TestIntegration_Spec_GetWithRequest(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

docID := "get-req-doc"
_, err := client.Create(ctx, idx, docID, &create.Request{
Document: map[string]any{"field1": "a", "field2": "b"},
})
assert.NilError(t, err)
_, _ = client.IndicesRefresh(ctx)

// Using the semantic GetDocument for comparison
getRes, err := client.GetDocument(ctx, estype.Alias(idx), docID)
assert.NilError(t, err)
assert.Assert(t, getRes != nil)
assert.Assert(t, getRes.Found)

// Using the spec Get method
specGetRes, err := client.Get(ctx, idx, docID)
assert.NilError(t, err)
assert.Assert(t, specGetRes != nil)
assert.Assert(t, specGetRes.Found)
assert.Equal(t, getRes.Id_, specGetRes.Id_)
}

func TestIntegration_Spec_DeleteWithRequest(t *testing.T) {
client := newSpecClient(t)
ctx := context.Background()
idx := uniqueSpecIndex(t, client)

docID := "del-req-doc"
_, err := client.Create(ctx, idx, docID, &create.Request{
Document: map[string]any{"x": 1},
})
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
