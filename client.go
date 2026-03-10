// Package estypedgo provides a wrapper around the Elasticsearch go-elasticsearch typed client,
// offering a type-safe Go API with distinct types for Index names, Alias names, and other
// Elasticsearch concepts to prevent misuse. Logging is provided via the standard slog package.
package estypedgo

import (
	"context"
	"fmt"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	coredelete "github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	coreget "github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/info"
	coreidx "github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/reindex"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	idxcreate "github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	idxdelete "github.com/elastic/go-elasticsearch/v8/typedapi/indices/delete"
	idxputalias "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putalias"
	idxputsettings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putsettings"
	idxrefresh "github.com/elastic/go-elasticsearch/v8/typedapi/indices/refresh"
	idxupdatealiases "github.com/elastic/go-elasticsearch/v8/typedapi/indices/updatealiases"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// ESClient defines the interface for interacting with Elasticsearch.
// It is split into Index-oriented and Alias-oriented operations to encourage
// correct use of Index vs Alias types.
type ESClient interface {
	// Info returns information about the Elasticsearch cluster.
	Info(ctx context.Context) (*info.Response, error)

	// IndexRefresh refreshes the specified index.
	IndexRefresh(ctx context.Context, indexName estype.Index) (*idxrefresh.Response, error)

	// AliasRefresh refreshes the index (or indices) backing the specified alias.
	AliasRefresh(ctx context.Context, aliasName estype.Alias) (*idxrefresh.Response, error)

	// IndexDocumentCount returns the number of documents in the specified index.
	IndexDocumentCount(ctx context.Context, indexName estype.Index) (*count.Response, error)

	// CreateIndex creates an index with optional settings and mappings.
	CreateIndex(ctx context.Context, indexName estype.Index, settings *types.IndexSettings, mappings *types.TypeMapping) (*idxcreate.Response, error)

	// DeleteIndex deletes the specified index.
	DeleteIndex(ctx context.Context, indexName estype.Index) (*idxdelete.Response, error)

	// IndexExists reports whether the specified index exists.
	IndexExists(ctx context.Context, indexName estype.Index) (bool, error)

	// AliasExists reports whether the specified alias exists.
	AliasExists(ctx context.Context, aliasName estype.Alias) (bool, error)

	// GetIndicesForAlias returns all index names associated with the given alias.
	GetIndicesForAlias(ctx context.Context, aliasName estype.Alias) ([]estype.Index, error)

	// CreateAlias creates an alias pointing to an index.
	CreateAlias(ctx context.Context, indexName estype.Index, aliasName estype.Alias, isWriteIndex bool) (*idxputalias.Response, error)

	// UpdateAliases performs one or more alias add/remove actions atomically.
	UpdateAliases(ctx context.Context, actions []types.IndicesAction) (*idxupdatealiases.Response, error)

	// GetRefreshInterval returns the current refresh interval for the alias.
	GetRefreshInterval(ctx context.Context, aliasName estype.Alias) (estype.RefreshInterval, error)

	// UpdateRefreshInterval updates the refresh interval for the index backing the alias.
	UpdateRefreshInterval(ctx context.Context, aliasName estype.Alias, interval estype.RefreshInterval) (*idxputsettings.Response, error)

	// CreateDocument indexes (creates or replaces) a document in the alias and waits for refresh.
	CreateDocument(ctx context.Context, aliasName estype.Alias, id string, document any) (*coreidx.Response, error)

	// GetDocument retrieves a document from the alias by its ID.
	GetDocument(ctx context.Context, aliasName estype.Alias, id string) (*coreget.Response, error)

	// DeleteDocument deletes a document from the index by its ID.
	DeleteDocument(ctx context.Context, indexName estype.Index, id string) (*coredelete.Response, error)

	// UpdateDocument partially updates a document in the index.
	UpdateDocument(ctx context.Context, indexName estype.Index, id string, req *update.Request) (*update.Response, error)

	// Search executes a search request against the alias.
	Search(
		ctx context.Context,
		aliasName estype.Alias,
		query types.Query,
		limit int,
		offset int,
		sort []types.SortCombinations,
		aggregations map[string]types.Aggregations,
		highlight *types.Highlight,
		collapse *types.FieldCollapse,
		scriptFields map[string]types.ScriptField,
	) (*search.Response, error)

	// SearchWithRequest executes a search using a fully-constructed search.Request.
	// Use this for advanced scenarios not covered by the Search helper.
	SearchWithRequest(ctx context.Context, aliasName estype.Alias, req *search.Request) (*search.Response, error)

	// Reindex copies documents from sourceIndex to destIndex.
	Reindex(ctx context.Context, sourceIndex, destIndex estype.Index, waitForCompletion bool) (*reindex.Response, error)

	// DeltaReindex copies documents updated since the given time from sourceIndex to destIndex.
	DeltaReindex(
		ctx context.Context,
		sourceIndex, destIndex estype.Index,
		since time.Time,
		timestampField string,
		waitForCompletion bool,
	) (*reindex.Response, error)

	// WaitForTaskCompletion polls the task API until the task finishes or the timeout elapses.
	WaitForTaskCompletion(ctx context.Context, taskID types.TaskId, timeout time.Duration) error
}

// NewClient constructs an ESClient backed by the Elasticsearch typed client.
func NewClient(config es8.Config) (ESClient, error) {
	typedClient, err := es8.NewTypedClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch TypedClient: %w", err)
	}
	return newESClient(typedClient), nil
}

// ensure compile-time check that *esClient implements ESClient.
var _ ESClient = (*esClient)(nil)
