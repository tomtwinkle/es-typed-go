//go:build integration

package esv8_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
	esv8 "github.com/tomtwinkle/es-typed-go/esv8"
	"github.com/tomtwinkle/es-typed-go/query"
)

const (
	typedSearchFieldID       = estype.Field("id")
	typedSearchFieldName     = estype.Field("name")
	typedSearchFieldCategory = estype.Field("category")
)

type typedSearchDoc struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

func createTypedSearchIndex(t *testing.T, ctx context.Context, client esv8.ESClient, idx estype.Index, alias estype.Alias) {
	t.Helper()

	mappings := &types.TypeMapping{
		Properties: map[string]types.Property{
			typedSearchFieldID.String():       types.NewKeywordProperty(),
			typedSearchFieldName.String():     types.NewKeywordProperty(),
			typedSearchFieldCategory.String(): types.NewKeywordProperty(),
		},
	}

	_, err := client.CreateIndex(ctx, idx, noReplicaSettings(), mappings)
	assert.NilError(t, err)
	_, err = client.CreateAlias(ctx, idx, alias, estype.WriteIndexEnabled)
	assert.NilError(t, err)
}

func indexTypedSearchDocs(t *testing.T, ctx context.Context, client esv8.ESClient, alias estype.Alias, docs []typedSearchDoc) {
	t.Helper()

	for _, doc := range docs {
		_, err := client.CreateDocument(ctx, alias, doc.ID, doc)
		assert.NilError(t, err)
	}
}

func TestIntegration_TypedSearch_RequestOptions(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	ctx := context.Background()
	idx := uniqueIndex(t, client)
	alias := uniqueAlias(t)

	createTypedSearchIndex(t, ctx, client, idx, alias)
	indexTypedSearchDocs(t, ctx, client, alias, []typedSearchDoc{
		{ID: "doc-1", Name: "Alpha", Category: "kitchen"},
		{ID: "doc-2", Name: "Beta", Category: "kitchen"},
		{ID: "doc-3", Name: "Gamma", Category: "kitchen"},
		{ID: "doc-4", Name: "Delta", Category: "office"},
	})

	_, err := client.IndexRefresh(ctx, idx)
	assert.NilError(t, err)

	t.Run("typed params zero-sized search does not return parse error", func(t *testing.T) {
		t.Parallel()

		timeout := "30s"
		res, err := esv8.Search[typedSearchDoc](ctx, client, alias, esv8.SearchParams{
			Query:          query.TermValue(typedSearchFieldCategory, "kitchen"),
			Size:           0,
			HasSize:        true,
			From:           0,
			HasFrom:        true,
			TrackTotalHits: 10,
			Source:         false,
			Timeout:        &timeout,
		})
		assert.NilError(t, err)
		assert.Equal(t, int64(3), res.Total)
		assert.Equal(t, 0, len(res.Hits))
	})

	t.Run("builder search_after request does not return parse error", func(t *testing.T) {
		t.Parallel()

		params := query.NewSearch().
			Where(query.TermValue(typedSearchFieldCategory, "kitchen")).
			Sort(query.NewSort().Field(typedSearchFieldID, sortorder.Asc).Build()...).
			Limit(1).
			Offset(0).
			SearchAfter("doc-1").
			TrackTotalHits(true).
			Source(types.SourceFilter{Includes: []string{"id", "name"}}).
			Timeout("30s").
			Build()

		rawRes, err := client.SearchRaw(ctx, alias, params.ToRequest())
		assert.NilError(t, err)
		assert.Assert(t, rawRes.Hits.Total != nil)
		assert.Equal(t, int64(3), rawRes.Hits.Total.Value)
		assert.Equal(t, 1, len(rawRes.Hits.Hits))

		var got typedSearchDoc
		assert.NilError(t, json.Unmarshal(rawRes.Hits.Hits[0].Source_, &got))
		assert.Equal(t, "doc-2", got.ID)
		assert.Equal(t, "Beta", got.Name)
		assert.Equal(t, "", got.Category)
	})
}
