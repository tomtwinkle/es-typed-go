// Package esv8 provides a wrapper around the Elasticsearch go-elasticsearch v8
// typed client, offering a type-safe Go API with distinct types for Index names,
// Alias names, and other Elasticsearch concepts to prevent misuse.
// Logging is provided via the standard slog package.
//
//go:generate go run ./generator
package esv8

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	cat_aliases "github.com/elastic/go-elasticsearch/v8/typedapi/cat/aliases"
	cat_allocation "github.com/elastic/go-elasticsearch/v8/typedapi/cat/allocation"
	cat_component_templates "github.com/elastic/go-elasticsearch/v8/typedapi/cat/componenttemplates"
	cat_count "github.com/elastic/go-elasticsearch/v8/typedapi/cat/count"
	cat_fielddata "github.com/elastic/go-elasticsearch/v8/typedapi/cat/fielddata"
	cat_health "github.com/elastic/go-elasticsearch/v8/typedapi/cat/health"
	cat_help "github.com/elastic/go-elasticsearch/v8/typedapi/cat/help"
	cat_indices "github.com/elastic/go-elasticsearch/v8/typedapi/cat/indices"
	cat_master "github.com/elastic/go-elasticsearch/v8/typedapi/cat/master"
	cat_ml_datafeeds "github.com/elastic/go-elasticsearch/v8/typedapi/cat/mldatafeeds"
	cat_ml_data_frame_analytics "github.com/elastic/go-elasticsearch/v8/typedapi/cat/mldataframeanalytics"
	cat_ml_jobs "github.com/elastic/go-elasticsearch/v8/typedapi/cat/mljobs"
	cat_ml_trained_models "github.com/elastic/go-elasticsearch/v8/typedapi/cat/mltrainedmodels"
	cat_nodeattrs "github.com/elastic/go-elasticsearch/v8/typedapi/cat/nodeattrs"
	cat_nodes "github.com/elastic/go-elasticsearch/v8/typedapi/cat/nodes"
	cat_pending_tasks "github.com/elastic/go-elasticsearch/v8/typedapi/cat/pendingtasks"
	cat_plugins "github.com/elastic/go-elasticsearch/v8/typedapi/cat/plugins"
	cat_recovery "github.com/elastic/go-elasticsearch/v8/typedapi/cat/recovery"
	cat_repositories "github.com/elastic/go-elasticsearch/v8/typedapi/cat/repositories"
	cat_segments "github.com/elastic/go-elasticsearch/v8/typedapi/cat/segments"
	cat_shards "github.com/elastic/go-elasticsearch/v8/typedapi/cat/shards"
	cat_snapshots "github.com/elastic/go-elasticsearch/v8/typedapi/cat/snapshots"
	cat_tasks "github.com/elastic/go-elasticsearch/v8/typedapi/cat/tasks"
	cat_templates "github.com/elastic/go-elasticsearch/v8/typedapi/cat/templates"
	cat_thread_pool "github.com/elastic/go-elasticsearch/v8/typedapi/cat/threadpool"
	cat_transforms "github.com/elastic/go-elasticsearch/v8/typedapi/cat/transforms"
	ccr_follow "github.com/elastic/go-elasticsearch/v8/typedapi/ccr/follow"
	ccr_follow_stats "github.com/elastic/go-elasticsearch/v8/typedapi/ccr/followstats"
	ccr_pause_follow "github.com/elastic/go-elasticsearch/v8/typedapi/ccr/pausefollow"
	ccr_resume_follow "github.com/elastic/go-elasticsearch/v8/typedapi/ccr/resumefollow"
	ccr_unfollow "github.com/elastic/go-elasticsearch/v8/typedapi/ccr/unfollow"
	cluster_health "github.com/elastic/go-elasticsearch/v8/typedapi/cluster/health"
	core_bulk "github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	core_clear_scroll "github.com/elastic/go-elasticsearch/v8/typedapi/core/clearscroll"
	core_close_point_in_time "github.com/elastic/go-elasticsearch/v8/typedapi/core/closepointintime"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	core_count "github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	coredelete "github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	core_delete_by_query "github.com/elastic/go-elasticsearch/v8/typedapi/core/deletebyquery"
	coreget "github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	coreidx "github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/info"
	core_mget "github.com/elastic/go-elasticsearch/v8/typedapi/core/mget"
	core_msearch "github.com/elastic/go-elasticsearch/v8/typedapi/core/msearch"
	core_open_point_in_time "github.com/elastic/go-elasticsearch/v8/typedapi/core/openpointintime"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/reindex"
	core_scroll "github.com/elastic/go-elasticsearch/v8/typedapi/core/scroll"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	core_update_by_query "github.com/elastic/go-elasticsearch/v8/typedapi/core/updatebyquery"
	esql_query "github.com/elastic/go-elasticsearch/v8/typedapi/esql/query"
	ilm_explain_lifecycle "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/explainlifecycle"
	ilm_get_lifecycle "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/getlifecycle"
	ilm_put_lifecycle "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/putlifecycle"
	indices_analyze "github.com/elastic/go-elasticsearch/v8/typedapi/indices/analyze"
	indices_clear_cache "github.com/elastic/go-elasticsearch/v8/typedapi/indices/clearcache"
	indices_close "github.com/elastic/go-elasticsearch/v8/typedapi/indices/close"
	idxcreate "github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	indices_create_data_stream "github.com/elastic/go-elasticsearch/v8/typedapi/indices/createdatastream"
	idxdelete "github.com/elastic/go-elasticsearch/v8/typedapi/indices/delete"
	indices_delete_alias "github.com/elastic/go-elasticsearch/v8/typedapi/indices/deletealias"
	indices_delete_data_stream "github.com/elastic/go-elasticsearch/v8/typedapi/indices/deletedatastream"
	indices_delete_index_template "github.com/elastic/go-elasticsearch/v8/typedapi/indices/deleteindextemplate"
	indices_flush "github.com/elastic/go-elasticsearch/v8/typedapi/indices/flush"
	indices_forcemerge "github.com/elastic/go-elasticsearch/v8/typedapi/indices/forcemerge"
	indices_get_alias "github.com/elastic/go-elasticsearch/v8/typedapi/indices/getalias"
	indices_get_data_stream "github.com/elastic/go-elasticsearch/v8/typedapi/indices/getdatastream"
	indices_get_index_template "github.com/elastic/go-elasticsearch/v8/typedapi/indices/getindextemplate"
	indices_get_mapping "github.com/elastic/go-elasticsearch/v8/typedapi/indices/getmapping"
	indices_get_settings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/getsettings"
	indices_open "github.com/elastic/go-elasticsearch/v8/typedapi/indices/open"
	idxputalias "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putalias"
	indices_put_index_template "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putindextemplate"
	indices_put_mapping "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putmapping"
	idxputsettings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putsettings"
	indices_put_settings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putsettings"
	idxrefresh "github.com/elastic/go-elasticsearch/v8/typedapi/indices/refresh"
	indices_rollover "github.com/elastic/go-elasticsearch/v8/typedapi/indices/rollover"
	indices_stats "github.com/elastic/go-elasticsearch/v8/typedapi/indices/stats"
	idxupdatealiases "github.com/elastic/go-elasticsearch/v8/typedapi/indices/updatealiases"
	inference_delete "github.com/elastic/go-elasticsearch/v8/typedapi/inference/delete"
	inference_get "github.com/elastic/go-elasticsearch/v8/typedapi/inference/get"
	inference_inference "github.com/elastic/go-elasticsearch/v8/typedapi/inference/inference"
	inference_put "github.com/elastic/go-elasticsearch/v8/typedapi/inference/put"
	ingest_delete_pipeline "github.com/elastic/go-elasticsearch/v8/typedapi/ingest/deletepipeline"
	ingest_get_pipeline "github.com/elastic/go-elasticsearch/v8/typedapi/ingest/getpipeline"
	ingest_put_pipeline "github.com/elastic/go-elasticsearch/v8/typedapi/ingest/putpipeline"
	ml_close_job "github.com/elastic/go-elasticsearch/v8/typedapi/ml/closejob"
	ml_delete_datafeed "github.com/elastic/go-elasticsearch/v8/typedapi/ml/deletedatafeed"
	ml_delete_job "github.com/elastic/go-elasticsearch/v8/typedapi/ml/deletejob"
	ml_get_datafeeds "github.com/elastic/go-elasticsearch/v8/typedapi/ml/getdatafeeds"
	ml_get_jobs "github.com/elastic/go-elasticsearch/v8/typedapi/ml/getjobs"
	ml_open_job "github.com/elastic/go-elasticsearch/v8/typedapi/ml/openjob"
	ml_put_datafeed "github.com/elastic/go-elasticsearch/v8/typedapi/ml/putdatafeed"
	ml_put_job "github.com/elastic/go-elasticsearch/v8/typedapi/ml/putjob"
	ml_start_datafeed "github.com/elastic/go-elasticsearch/v8/typedapi/ml/startdatafeed"
	ml_stop_datafeed "github.com/elastic/go-elasticsearch/v8/typedapi/ml/stopdatafeed"
	security_create_api_key "github.com/elastic/go-elasticsearch/v8/typedapi/security/createapikey"
	security_get_api_key "github.com/elastic/go-elasticsearch/v8/typedapi/security/getapikey"
	security_invalidate_api_key "github.com/elastic/go-elasticsearch/v8/typedapi/security/invalidateapikey"
	snapshot_create "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/create"
	snapshot_create_repository "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/createrepository"
	snapshot_restore "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/restore"
	tasks_cancel "github.com/elastic/go-elasticsearch/v8/typedapi/tasks/cancel"
	tasks_list "github.com/elastic/go-elasticsearch/v8/typedapi/tasks/list"
	transform_delete_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/deletetransform"
	transform_get_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/gettransform"
	transform_get_transform_stats "github.com/elastic/go-elasticsearch/v8/typedapi/transform/gettransformstats"
	transform_put_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/puttransform"
	transform_start_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/starttransform"
	transform_stop_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/stoptransform"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// ESClient defines the interface for interacting with Elasticsearch v8.
// It is split into Index-oriented and Alias-oriented operations to encourage
// correct use of Index vs Alias types.
type ESClient interface {
	// Info returns information about the Elasticsearch cluster.
	//
	// Returns basic information about the Elasticsearch cluster, such as its name, version, and build information.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/info-api.html
	Info(ctx context.Context) (*info.Response, error)

	// IndexRefresh refreshes the specified index.
	//
	// Forces a refresh on one or more indices, making all recent write operations available for search. A refresh is expensive relative to indexing; prefer relying on the automatic refresh interval for bulk indexing workloads.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-refresh.html
	IndexRefresh(ctx context.Context, indexName estype.Index, opts ...IndexRefreshOption) (*idxrefresh.Response, error)

	// AliasRefresh refreshes the index (or indices) backing the specified alias.
	//
	// Forces a refresh on the index (or indices) backing the specified alias. Resolves the alias to its backing indices before issuing the refresh.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-refresh.html
	AliasRefresh(ctx context.Context, aliasName estype.Alias, opts ...IndexRefreshOption) (*idxrefresh.Response, error)

	// IndexDocumentCount returns the number of documents in the specified index.
	//
	// Returns the number of documents in the specified index. This is a convenience wrapper around the Count API.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-count.html
	IndexDocumentCount(ctx context.Context, indexName estype.Index) (*count.Response, error)

	// CreateIndex creates an index with optional settings and mappings.
	//
	// Creates a new index with optional settings and mappings. Returns an error if the index already exists.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html
	CreateIndex(ctx context.Context, indexName estype.Index, settings *types.IndexSettings, mappings *types.TypeMapping) (*idxcreate.Response, error)

	// DeleteIndex deletes the specified index.
	//
	// Deletes an index, permanently removing all its documents and metadata.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-delete-index.html
	DeleteIndex(ctx context.Context, indexName estype.Index, opts ...DeleteIndexOption) (*idxdelete.Response, error)

	// IndexExists reports whether the specified index exists.
	//
	// Returns true if the specified index exists, false otherwise.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-exists.html
	IndexExists(ctx context.Context, indexName estype.Index) (bool, error)

	// AliasExists reports whether the specified alias exists.
	//
	// Returns true if the specified alias exists in the cluster, false otherwise.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-exists.html
	AliasExists(ctx context.Context, aliasName estype.Alias) (bool, error)

	// GetIndicesForAlias returns all index names associated with the given alias.
	//
	// Returns the list of index names that the given alias points to.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-get-alias.html
	GetIndicesForAlias(ctx context.Context, aliasName estype.Alias) ([]estype.Index, error)

	// CreateAlias creates an alias pointing to an index.
	//
	// Creates an alias pointing to the specified index. When isWriteIndex is true, the index is designated as the write target for the alias.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-add-alias.html
	CreateAlias(ctx context.Context, indexName estype.Index, aliasName estype.Alias, isWriteIndex bool) (*idxputalias.Response, error)

	// UpdateAliases performs one or more alias add/remove actions atomically.
	//
	// Atomically performs multiple alias add and remove actions. All changes succeed or none are applied.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-aliases.html
	UpdateAliases(ctx context.Context, actions []types.IndicesAction) (*idxupdatealiases.Response, error)

	// GetRefreshInterval returns the current refresh interval for the alias.
	//
	// Returns the current refresh_interval setting for the index backing the alias.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-get-settings.html
	GetRefreshInterval(ctx context.Context, aliasName estype.Alias) (estype.RefreshInterval, error)

	// UpdateRefreshInterval updates the refresh interval for the index backing the alias.
	//
	// Updates the refresh_interval index setting for the index backing the alias. Set the interval to -1 to disable automatic refreshes.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-put-settings.html
	UpdateRefreshInterval(ctx context.Context, aliasName estype.Alias, interval estype.RefreshInterval) (*idxputsettings.Response, error)

	// CreateDocument indexes (creates or replaces) a document in the alias and waits for refresh.
	//
	// Indexes a document in the index backing the alias with the given ID, creating it if it does not exist or replacing it if it does. By default waits for a refresh before returning.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-index_.html
	CreateDocument(ctx context.Context, aliasName estype.Alias, id string, document any, opts ...CreateDocumentOption) (*coreidx.Response, error)

	// GetDocument retrieves a document from the alias by its ID.
	//
	// Retrieves a document by its ID from the index backing the alias.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-get.html
	GetDocument(ctx context.Context, aliasName estype.Alias, id string, opts ...GetDocumentOption) (*coreget.Response, error)

	// DeleteDocument deletes a document from the index by its ID.
	//
	// Deletes a document from the specified index by its ID.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-delete.html
	DeleteDocument(ctx context.Context, indexName estype.Index, id string, opts ...DeleteDocumentOption) (*coredelete.Response, error)

	// UpdateDocument partially updates a document in the index.
	//
	// Partially updates a document in the specified index. The document is merged with the fields provided in the request.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-update.html
	UpdateDocument(ctx context.Context, indexName estype.Index, id string, req *update.Request) (*update.Response, error)

	// SearchRaw executes a search using a fully-constructed search.Request.
	// Use this as the low-level escape hatch beside the high-level Search[T],
	// SearchDocuments[T], and SearchOne[T] helpers.
	//
	// This is intended for advanced scenarios such as kNN search,
	// point-in-time, search_after, or custom source filtering that are not yet
	// modeled by SearchParams.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-search.html
	SearchRaw(ctx context.Context, aliasName estype.Alias, req *search.Request) (*search.Response, error)

	// Reindex copies documents from sourceIndex to destIndex.
	//
	// Copies documents from sourceIndex to destIndex. When waitForCompletion is false, returns immediately with a task ID that can be tracked with WaitForTaskCompletion.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-reindex.html
	Reindex(ctx context.Context, sourceIndex, destIndex estype.Index, waitForCompletion bool, opts ...ReindexOption) (*reindex.Response, error)

	// DeltaReindex copies documents updated since the given time from sourceIndex to destIndex.
	//
	// Copies documents updated after a given point in time from sourceIndex to destIndex, using a range query on timestampField.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-reindex.html
	DeltaReindex(
		ctx context.Context,
		sourceIndex, destIndex estype.Index,
		since time.Time,
		timestampField string,
		waitForCompletion bool,
		opts ...ReindexOption,
	) (*reindex.Response, error)

	// WaitForTaskCompletion polls the task API until the task finishes or the timeout elapses.
	//
	// Polls the Tasks API until the specified task completes or the timeout elapses. Useful for monitoring asynchronous operations such as reindex.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/tasks.html
	WaitForTaskCompletion(ctx context.Context, taskID types.TaskId, timeout time.Duration) error

	// ---------------------------------------------------------------------------
	// Document & Search API
	// ---------------------------------------------------------------------------

	// DocumentExists reports whether a document with the given ID exists in the alias.
	//
	// Returns true if a document with the given ID exists in the alias, false otherwise.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-get.html
	DocumentExists(ctx context.Context, aliasName estype.Alias, id estype.DocumentID, opts ...DocumentExistsOption) (bool, error)

	// Bulk performs multiple index, create, delete, or update operations in a single request.
	//
	// Performs multiple index, create, delete, or update operations in a single API call, reducing network overhead. Use BulkOption to set the request body via r.Request(&req).
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-bulk.html
	Bulk(ctx context.Context, aliasName estype.Alias, opts ...BulkOption) (*core_bulk.Response, error)

	// Mget retrieves multiple documents by ID in a single request.
	//
	// Retrieves multiple documents by ID in a single request. Use MgetOption to set the request body via r.Request(&req).
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-multi-get.html
	Mget(ctx context.Context, aliasName estype.Alias, opts ...MgetOption) (*core_mget.Response, error)

	// Msearch executes multiple search requests in a single request.
	//
	// Executes multiple search requests in a single API call. Use MsearchOption to supply request bodies.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-multi-search.html
	Msearch(ctx context.Context, aliasName estype.Alias, opts ...MsearchOption) (*core_msearch.Response, error)

	// Count returns the number of documents matching a query against the alias.
	//
	// Returns the number of documents matching a query. Without a query set via CountOption, counts all documents in the alias.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-count.html
	Count(ctx context.Context, aliasName estype.Alias, opts ...CountOption) (*core_count.Response, error)

	// Scroll retrieves the next batch of results from a scroll operation.
	//
	// Retrieves the next batch of results for an ongoing scroll operation. Use ScrollOption to set the scroll ID and keep-alive interval.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/paginate-search-results.html#scroll-search-results
	Scroll(ctx context.Context, opts ...ScrollOption) (*core_scroll.Response, error)

	// ClearScroll clears the search context and results for a scroll.
	//
	// Clears one or more scroll contexts, freeing the associated server-side resources.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/clear-scroll-api.html
	ClearScroll(ctx context.Context, opts ...ClearScrollOption) (*core_clear_scroll.Response, error)

	// UpdateByQuery updates documents that match the given query in the index.
	//
	// Updates documents matching a query in-place. Commonly used to apply a script to many documents or to pick up new field mappings.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-update-by-query.html
	UpdateByQuery(ctx context.Context, indexName estype.Index, opts ...UpdateByQueryOption) (*core_update_by_query.Response, error)

	// DeleteByQuery deletes documents that match the given query in the index.
	//
	// Deletes all documents in the index that match the given query. Returns counts of deleted, failed, and version-conflicting documents.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/docs-delete-by-query.html
	DeleteByQuery(ctx context.Context, indexName estype.Index, opts ...DeleteByQueryOption) (*core_delete_by_query.Response, error)

	// ---------------------------------------------------------------------------
	// Index Management API
	// ---------------------------------------------------------------------------

	// GetMapping returns the mapping for the specified index.
	//
	// Returns the field mapping definition for the specified index.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-get-mapping.html
	GetMapping(ctx context.Context, indexName estype.Index, opts ...GetMappingOption) (indices_get_mapping.Response, error)

	// PutMapping updates the field mappings for the specified index.
	//
	// Adds new fields to an existing index or changes the search settings of existing fields. Existing field types cannot be changed.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-put-mapping.html
	PutMapping(ctx context.Context, indexName estype.Index, opts ...PutMappingOption) (*indices_put_mapping.Response, error)

	// GetSettings returns the settings for the specified index.
	//
	// Returns the current index-level settings, such as the number of shards, replicas, and analyzer configuration.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-get-settings.html
	GetSettings(ctx context.Context, indexName estype.Index, opts ...GetSettingsOption) (indices_get_settings.Response, error)

	// PutSettings updates the settings for the specified index.
	//
	// Updates index settings. Dynamic settings can be changed on open indices; static settings require the index to be closed first.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-put-settings.html
	PutSettings(ctx context.Context, indexName estype.Index, opts ...PutSettingsOption) (*indices_put_settings.Response, error)

	// OpenIndex opens a closed index.
	//
	// Reopens a previously closed index, making it available for indexing and search again.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-open-close.html
	OpenIndex(ctx context.Context, indexName estype.Index, opts ...OpenIndexOption) (*indices_open.Response, error)

	// CloseIndex closes an open index.
	//
	// Closes an open index. Closed indices cannot be read or written but still occupy disk space and count towards shard limits.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-open-close.html
	CloseIndex(ctx context.Context, indexName estype.Index, opts ...CloseIndexOption) (*indices_close.Response, error)

	// Flush flushes one or more indices.
	//
	// Flushes one or more indices to durable storage, ensuring that data currently in the transaction log is written to Lucene and the log is cleared.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-flush.html
	Flush(ctx context.Context, indexName estype.Index, opts ...FlushOption) (*indices_flush.Response, error)

	// ClearCache clears the caches of one or more indices.
	//
	// Clears in-memory caches for one or more indices. Supports clearing the query cache, field-data cache, and request cache.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-clearcache.html
	ClearCache(ctx context.Context, indexName estype.Index, opts ...ClearCacheOption) (*indices_clear_cache.Response, error)

	// ForceMerge forces a merge on the shards of one or more indices.
	//
	// Forces a Lucene segment merge on one or more indices, reducing the number of segments. Use on read-only or infrequently updated indices; it is a costly operation.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-forcemerge.html
	ForceMerge(ctx context.Context, indexName estype.Index, opts ...ForceMergeOption) (*indices_forcemerge.Response, error)

	// Rollover rolls an alias over to a new index when the existing index meets a condition.
	//
	// Rolls an alias over to a new index when the current index satisfies the specified conditions (document count, size, or age).
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-rollover-index.html
	Rollover(ctx context.Context, aliasName estype.Alias, opts ...RolloverOption) (*indices_rollover.Response, error)

	// IndicesStats returns statistics for the specified index.
	//
	// Returns statistics for one or more indices including document count, store size, and indexing and search metrics.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-stats.html
	IndicesStats(ctx context.Context, indexName estype.Index, opts ...IndicesStatsOption) (*indices_stats.Response, error)

	// ---------------------------------------------------------------------------
	// Cluster API
	// ---------------------------------------------------------------------------

	// ClusterHealth returns the health status of the specified index.
	//
	// Returns the health status of the cluster or a specific index. Status is green (all shards active), yellow (replica shards unassigned), or red (primary shards unassigned).
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cluster-health.html
	ClusterHealth(ctx context.Context, indexName estype.Index, opts ...ClusterHealthOption) (*cluster_health.Response, error)

	// ---------------------------------------------------------------------------
	// Index Template API
	// ---------------------------------------------------------------------------

	// PutIndexTemplate creates or updates an index template.
	//
	// Creates or updates a composable index template that is applied when new matching indices are created.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-put-template.html
	PutIndexTemplate(ctx context.Context, name estype.Template, opts ...PutIndexTemplateOption) (*indices_put_index_template.Response, error)

	// GetIndexTemplate returns an index template.
	//
	// Returns the configuration of one or more index templates.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-get-template.html
	GetIndexTemplate(ctx context.Context, name estype.Template, opts ...GetIndexTemplateOption) (*indices_get_index_template.Response, error)

	// DeleteIndexTemplate deletes an index template.
	//
	// Deletes one or more index templates.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-delete-template.html
	DeleteIndexTemplate(ctx context.Context, name estype.Template, opts ...DeleteIndexTemplateOption) (*indices_delete_index_template.Response, error)

	// ExistsIndexTemplate reports whether an index template exists.
	//
	// Returns true if an index template with the given name exists.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-exists.html
	ExistsIndexTemplate(ctx context.Context, name estype.Template, opts ...ExistsIndexTemplateOption) (bool, error)

	// ---------------------------------------------------------------------------
	// Alias API
	// ---------------------------------------------------------------------------

	// GetAlias returns alias information for the specified alias.
	//
	// Returns information about one or more aliases, including the indices they point to and any routing or filter settings.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-get-alias.html
	GetAlias(ctx context.Context, aliasName estype.Alias, opts ...GetAliasOption) (indices_get_alias.Response, error)

	// DeleteAlias removes an alias from the specified index.
	//
	// Removes an alias from the specified index.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-delete-alias.html
	DeleteAlias(ctx context.Context, indexName estype.Index, aliasName estype.Alias, opts ...DeleteAliasOption) (*indices_delete_alias.Response, error)

	// ---------------------------------------------------------------------------
	// Task Management API
	// ---------------------------------------------------------------------------

	// TasksList returns a list of currently running tasks.
	//
	// Returns information about tasks currently executing in the cluster. Useful for monitoring long-running operations.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/tasks.html
	TasksList(ctx context.Context, opts ...TasksListOption) (*tasks_list.Response, error)

	// TasksCancel cancels a task or a group of tasks.
	//
	// Attempts to cancel a running task. Not all tasks support cancellation.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/tasks.html
	TasksCancel(ctx context.Context, taskID estype.TaskID, opts ...TasksCancelOption) (*tasks_cancel.Response, error)

	// ---------------------------------------------------------------------------
	// Analysis / Debug API
	// ---------------------------------------------------------------------------

	// Analyze performs analysis on a text string and returns the tokens.
	//
	// Performs text analysis on a string using the specified index's analyzer configuration and returns the resulting tokens.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-analyze.html
	Analyze(ctx context.Context, indexName estype.Index, opts ...AnalyzeOption) (*indices_analyze.Response, error)

	// ---------------------------------------------------------------------------
	// PIT (Point In Time) API
	// ---------------------------------------------------------------------------

	// OpenPointInTime opens a point in time on the specified index.
	//
	// Opens a lightweight point-in-time snapshot that allows consistent pagination across multiple search requests even as the underlying data changes.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/point-in-time-api.html
	OpenPointInTime(ctx context.Context, indexName estype.Index, keepAlive estype.KeepAlive, opts ...OpenPointInTimeOption) (*core_open_point_in_time.Response, error)

	// ClosePointInTime closes a point in time.
	//
	// Closes a point-in-time snapshot and releases its server-side resources.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/point-in-time-api.html
	ClosePointInTime(ctx context.Context, opts ...ClosePointInTimeOption) (*core_close_point_in_time.Response, error)

	// ---------------------------------------------------------------------------
	// Data Stream API
	// ---------------------------------------------------------------------------

	// CreateDataStream creates a data stream.
	//
	// Creates a data stream backed by auto-created time-series indices, suitable for append-only time-stamped data.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-data-stream.html
	CreateDataStream(ctx context.Context, name estype.DataStream) (*indices_create_data_stream.Response, error)

	// GetDataStream returns information about one or more data streams.
	//
	// Returns information about one or more data streams, including their backing indices and lifecycle status.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-get-data-stream.html
	GetDataStream(ctx context.Context, name estype.DataStream, opts ...GetDataStreamOption) (*indices_get_data_stream.Response, error)

	// DeleteDataStream deletes one or more data streams.
	//
	// Deletes one or more data streams and all their backing indices.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-delete-data-stream.html
	DeleteDataStream(ctx context.Context, name estype.DataStream, opts ...DeleteDataStreamOption) (*indices_delete_data_stream.Response, error)

	// ---------------------------------------------------------------------------
	// ILM (Index Lifecycle Management) API
	// ---------------------------------------------------------------------------

	// PutLifecycle creates or updates a lifecycle policy.
	//
	// Creates or updates an ILM lifecycle policy defining how indices transition through hot, warm, cold, and delete phases.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ilm-put-lifecycle.html
	PutLifecycle(ctx context.Context, policy estype.Policy, opts ...PutLifecycleOption) (*ilm_put_lifecycle.Response, error)

	// GetLifecycle returns lifecycle policy information.
	//
	// Returns the configuration of one or more ILM lifecycle policies.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ilm-get-lifecycle.html
	GetLifecycle(ctx context.Context, policy estype.Policy, opts ...GetLifecycleOption) (ilm_get_lifecycle.Response, error)

	// ExplainLifecycle returns the current lifecycle status for one or more indices.
	//
	// Returns the current ILM status for one or more indices, showing their current phase, action, and any policy errors.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ilm-explain-lifecycle.html
	ExplainLifecycle(ctx context.Context, indexName estype.Index, opts ...ExplainLifecycleOption) (*ilm_explain_lifecycle.Response, error)

	// ---------------------------------------------------------------------------
	// Ingest Pipeline API
	// ---------------------------------------------------------------------------

	// PutPipeline creates or updates an ingest pipeline.
	//
	// Creates or updates an ingest pipeline that pre-processes documents with a sequence of processors before they are indexed.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/put-pipeline-api.html
	PutPipeline(ctx context.Context, id estype.Pipeline, opts ...PutPipelineOption) (*ingest_put_pipeline.Response, error)

	// GetPipeline returns information about one or more ingest pipelines.
	//
	// Returns the configuration of one or more ingest pipelines.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/get-pipeline-api.html
	GetPipeline(ctx context.Context, id estype.Pipeline, opts ...GetPipelineOption) (ingest_get_pipeline.Response, error)

	// DeletePipeline deletes a pipeline.
	//
	// Deletes an ingest pipeline.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/delete-pipeline-api.html
	DeletePipeline(ctx context.Context, id estype.Pipeline, opts ...DeletePipelineOption) (*ingest_delete_pipeline.Response, error)

	// ---------------------------------------------------------------------------
	// ES|QL API
	// ---------------------------------------------------------------------------

	// EsqlQuery executes an ES|QL query.
	//
	// Executes an ES|QL query and returns the results in a tabular format. ES|QL provides a pipe-based query language for aggregating and transforming Elasticsearch data.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/esql-query-api.html
	EsqlQuery(ctx context.Context, query estype.ESQLQuery, opts ...EsqlQueryOption) (esql_query.Response, error)

	// ---------------------------------------------------------------------------
	// Security & API Key API
	// ---------------------------------------------------------------------------

	// CreateApiKey creates an API key.
	//
	// Creates an API key for authentication. Keys do not expire unless an expiration is set.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-create-api-key.html
	CreateApiKey(ctx context.Context, opts ...CreateApiKeyOption) (*security_create_api_key.Response, error)

	// GetApiKey retrieves information for one or more API keys.
	//
	// Returns information about one or more API keys, including their ID, name, creation date, and expiration.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-get-api-key.html
	GetApiKey(ctx context.Context, opts ...GetApiKeyOption) (*security_get_api_key.Response, error)

	// InvalidateApiKey invalidates one or more API keys.
	//
	// Invalidates one or more API keys, preventing their use for subsequent requests.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/security-api-invalidate-api-key.html
	InvalidateApiKey(ctx context.Context, opts ...InvalidateApiKeyOption) (*security_invalidate_api_key.Response, error)

	// ---------------------------------------------------------------------------
	// Snapshot & Restore API
	// ---------------------------------------------------------------------------

	// CreateRepository creates or updates a snapshot repository.
	//
	// Registers a snapshot repository (S3, GCS, Azure, HDFS, or shared filesystem) where snapshots will be stored.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/put-snapshot-repo-api.html
	CreateRepository(ctx context.Context, repo estype.Repository, opts ...CreateRepositoryOption) (*snapshot_create_repository.Response, error)

	// CreateSnapshot creates a snapshot in a repository.
	//
	// Creates a snapshot of one or more indices in the specified repository.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/create-snapshot-api.html
	CreateSnapshot(ctx context.Context, repo estype.Repository, snap estype.Snapshot, opts ...CreateSnapshotOption) (*snapshot_create.Response, error)

	// RestoreSnapshot restores a snapshot.
	//
	// Restores indices from a snapshot stored in the specified repository.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/restore-snapshot-api.html
	RestoreSnapshot(ctx context.Context, repo estype.Repository, snap estype.Snapshot, opts ...RestoreSnapshotOption) (*snapshot_restore.Response, error)

	// ---------------------------------------------------------------------------
	// Inference API
	// ---------------------------------------------------------------------------

	// PutInference creates or updates an inference endpoint.
	//
	// Creates or updates an inference endpoint that wraps a machine learning model for use in ingest pipelines or search.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/put-inference-api.html
	PutInference(ctx context.Context, inferenceId estype.InferenceID, opts ...PutInferenceOption) (*inference_put.Response, error)

	// GetInference returns information about an inference endpoint.
	//
	// Returns configuration information for one or more inference endpoints.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/get-inference-api.html
	GetInference(ctx context.Context, inferenceId estype.InferenceID, opts ...GetInferenceOption) (*inference_get.Response, error)

	// DeleteInference deletes an inference endpoint.
	//
	// Deletes an inference endpoint.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/delete-inference-api.html
	DeleteInference(ctx context.Context, inferenceId estype.InferenceID, opts ...DeleteInferenceOption) (*inference_delete.Response, error)

	// Inference performs an inference request against an inference endpoint.
	//
	// Performs an inference request against the specified endpoint, returning results such as embeddings or generated text.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/post-inference-api.html
	Inference(ctx context.Context, inferenceId estype.InferenceID, opts ...InferenceOption) (*inference_inference.Response, error)

	// ---------------------------------------------------------------------------
	// Machine Learning (ML) API
	// ---------------------------------------------------------------------------

	// MlPutJob creates an anomaly detection job.
	//
	// Creates an anomaly detection job that uses machine learning to model data behaviour and detect anomalies.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-put-job.html
	MlPutJob(ctx context.Context, jobId estype.MLJobID, opts ...MlPutJobOption) (*ml_put_job.Response, error)

	// MlGetJobs returns configuration information for anomaly detection jobs.
	//
	// Returns configuration and status information for one or more anomaly detection jobs.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-get-job.html
	MlGetJobs(ctx context.Context, jobId estype.MLJobID, opts ...MlGetJobsOption) (*ml_get_jobs.Response, error)

	// MlDeleteJob deletes an anomaly detection job.
	//
	// Deletes an anomaly detection job along with its model state and results.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-delete-job.html
	MlDeleteJob(ctx context.Context, jobId estype.MLJobID, opts ...MlDeleteJobOption) (*ml_delete_job.Response, error)

	// MlOpenJob opens an anomaly detection job.
	//
	// Opens one or more anomaly detection jobs so they can receive and process data.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-open-job.html
	MlOpenJob(ctx context.Context, jobId estype.MLJobID, opts ...MlOpenJobOption) (*ml_open_job.Response, error)

	// MlCloseJob closes an anomaly detection job.
	//
	// Closes one or more anomaly detection jobs. The job retains its configuration and model state.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-close-job.html
	MlCloseJob(ctx context.Context, jobId estype.MLJobID, opts ...MlCloseJobOption) (*ml_close_job.Response, error)

	// MlPutDatafeed creates a datafeed.
	//
	// Creates a datafeed that retrieves data from Elasticsearch and feeds it to an anomaly detection job.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-put-datafeed.html
	MlPutDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlPutDatafeedOption) (*ml_put_datafeed.Response, error)

	// MlGetDatafeeds returns configuration information for datafeeds.
	//
	// Returns configuration information for one or more datafeeds.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-get-datafeed.html
	MlGetDatafeeds(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlGetDatafeedsOption) (*ml_get_datafeeds.Response, error)

	// MlDeleteDatafeed deletes a datafeed.
	//
	// Deletes an existing datafeed.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-delete-datafeed.html
	MlDeleteDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlDeleteDatafeedOption) (*ml_delete_datafeed.Response, error)

	// MlStartDatafeed starts one or more datafeeds.
	//
	// Starts one or more datafeeds so they begin retrieving data for anomaly detection.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-start-datafeed.html
	MlStartDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlStartDatafeedOption) (*ml_start_datafeed.Response, error)

	// MlStopDatafeed stops one or more datafeeds.
	//
	// Stops one or more running datafeeds.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ml-stop-datafeed.html
	MlStopDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlStopDatafeedOption) (*ml_stop_datafeed.Response, error)

	// ---------------------------------------------------------------------------
	// CCR (Cross-Cluster Replication) API
	// ---------------------------------------------------------------------------

	// CcrFollow configures a local index to follow a remote index.
	//
	// Creates a cross-cluster replication (CCR) follower index that replicates a leader index from a remote cluster.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ccr-put-follow.html
	CcrFollow(ctx context.Context, followerIndex estype.Index, opts ...CcrFollowOption) (*ccr_follow.Response, error)

	// CcrPauseFollow pauses a follower index.
	//
	// Pauses the cross-cluster replication process for a follower index.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ccr-post-pause-follow.html
	CcrPauseFollow(ctx context.Context, followerIndex estype.Index, opts ...CcrPauseFollowOption) (*ccr_pause_follow.Response, error)

	// CcrResumeFollow resumes a follower index.
	//
	// Resumes cross-cluster replication for a paused follower index.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ccr-post-resume-follow.html
	CcrResumeFollow(ctx context.Context, followerIndex estype.Index, opts ...CcrResumeFollowOption) (*ccr_resume_follow.Response, error)

	// CcrUnfollow stops the following task for a follower index.
	//
	// Stops replication and converts a follower index into a regular, standalone index.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ccr-post-unfollow.html
	CcrUnfollow(ctx context.Context, followerIndex estype.Index, opts ...CcrUnfollowOption) (*ccr_unfollow.Response, error)

	// CcrFollowStats returns cross-cluster replication follower stats.
	//
	// Returns cross-cluster replication statistics for one or more follower indices.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/ccr-get-follow-stats.html
	CcrFollowStats(ctx context.Context, followerIndex estype.Index, opts ...CcrFollowStatsOption) (*ccr_follow_stats.Response, error)

	// ---------------------------------------------------------------------------
	// Transform API
	// ---------------------------------------------------------------------------

	// PutTransform creates or updates a transform.
	//
	// Creates or updates a transform that continuously or in batches pivots source data into a destination index.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/put-transform.html
	PutTransform(ctx context.Context, transformId estype.TransformID, opts ...PutTransformOption) (*transform_put_transform.Response, error)

	// GetTransform returns configuration information for transforms.
	//
	// Returns configuration information for one or more transforms.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/get-transform.html
	GetTransform(ctx context.Context, transformId estype.TransformID, opts ...GetTransformOption) (*transform_get_transform.Response, error)

	// DeleteTransform deletes a transform.
	//
	// Deletes a transform.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/delete-transform.html
	DeleteTransform(ctx context.Context, transformId estype.TransformID, opts ...DeleteTransformOption) (*transform_delete_transform.Response, error)

	// StartTransform starts one or more transforms.
	//
	// Starts one or more transforms so they begin indexing transformed data.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/start-transform.html
	StartTransform(ctx context.Context, transformId estype.TransformID, opts ...StartTransformOption) (*transform_start_transform.Response, error)

	// StopTransform stops one or more transforms.
	//
	// Stops one or more running transforms.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/stop-transform.html
	StopTransform(ctx context.Context, transformId estype.TransformID, opts ...StopTransformOption) (*transform_stop_transform.Response, error)

	// GetTransformStats returns usage information for transforms.
	//
	// Returns usage and state information for one or more transforms.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/get-transform-stats.html
	GetTransformStats(ctx context.Context, transformId estype.TransformID, opts ...GetTransformStatsOption) (*transform_get_transform_stats.Response, error)

	// ---------------------------------------------------------------------------
	// Cat API
	// ---------------------------------------------------------------------------

	// CatAliases returns information about aliases.
	//
	// Returns a concise table of alias information, including the index each alias points to and any routing configuration.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-alias.html
	CatAliases(ctx context.Context, aliasName estype.Alias, opts ...CatAliasesOption) (cat_aliases.Response, error)

	// CatIndices returns high-level information about indices.
	//
	// Returns a high-level summary of index information, including health, status, primary shard count, document count, and store size.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-indices.html
	CatIndices(ctx context.Context, indexName estype.Index, opts ...CatIndicesOption) (cat_indices.Response, error)

	// CatNodes returns information about the nodes in a cluster.
	//
	// Returns a table of information about nodes in the cluster, including their roles, memory usage, CPU load, and heap usage.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-nodes.html
	CatNodes(ctx context.Context, opts ...CatNodesOption) (cat_nodes.Response, error)

	// CatAllocation provides a snapshot of the number of shards allocated to each data node.
	//
	// Returns a snapshot of the number of shards allocated to each data node and the disk space used and available.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-allocation.html
	CatAllocation(ctx context.Context, opts ...CatAllocationOption) (cat_allocation.Response, error)

	// CatCircuitBreaker returns information about circuit breaker status.
	// NOTE: This endpoint is not available in the Elasticsearch v8 TypedAPI.
	// The response is returned as raw JSON.
	//
	// Returns information about the field-data circuit breaker for each node.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-fielddata.html
	CatCircuitBreaker(ctx context.Context, opts ...CatCircuitBreakerOption) (json.RawMessage, error)

	// CatComponentTemplates returns information about component templates.
	//
	// Returns a list of component templates that can be referenced in index templates.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-component-templates.html
	CatComponentTemplates(ctx context.Context, name estype.Template, opts ...CatComponentTemplatesOption) (cat_component_templates.Response, error)

	// CatCount returns document counts for one or more indices.
	//
	// Returns the total document count for the cluster or a specific index.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-count.html
	CatCount(ctx context.Context, indexName estype.Index, opts ...CatCountOption) (cat_count.Response, error)

	// CatFielddata returns information about the amount of memory used for field data.
	//
	// Returns information about the heap memory used by field-data on each node.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-fielddata.html
	CatFielddata(ctx context.Context, fields []estype.Field, opts ...CatFielddataOption) (cat_fielddata.Response, error)

	// CatHealth returns a concise representation of cluster health.
	//
	// Returns a concise one-row-per-cluster health summary, equivalent to the cluster health API.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-health.html
	CatHealth(ctx context.Context, opts ...CatHealthOption) (cat_health.Response, error)

	// CatHelp returns help for the cat APIs.
	//
	// Returns the list of available cat APIs.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat.html
	CatHelp(ctx context.Context, opts ...CatHelpOption) (*cat_help.Response, error)

	// CatMaster returns information about the master node.
	//
	// Returns information about the elected master node.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-master.html
	CatMaster(ctx context.Context, opts ...CatMasterOption) (cat_master.Response, error)

	// CatMlDataFrameAnalytics returns configuration and usage information about data frame analytics jobs.
	//
	// Returns configuration and usage information for data frame analytics jobs.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-dfanalytics.html
	CatMlDataFrameAnalytics(ctx context.Context, analyticsId estype.DataFrameAnalyticsID, opts ...CatMlDataFrameAnalyticsOption) (cat_ml_data_frame_analytics.Response, error)

	// CatMlDatafeeds returns configuration and usage information about datafeeds.
	//
	// Returns configuration and usage statistics for ML datafeeds.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-datafeeds.html
	CatMlDatafeeds(ctx context.Context, datafeedId estype.DatafeedID, opts ...CatMlDatafeedsOption) (cat_ml_datafeeds.Response, error)

	// CatMlJobs returns configuration and usage information about anomaly detection jobs.
	//
	// Returns configuration and usage information for anomaly detection jobs.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-anomaly-detectors.html
	CatMlJobs(ctx context.Context, jobId estype.MLJobID, opts ...CatMlJobsOption) (cat_ml_jobs.Response, error)

	// CatMlTrainedModels returns configuration and usage information about inference trained models.
	//
	// Returns configuration and usage information for trained inference models.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-trained-model.html
	CatMlTrainedModels(ctx context.Context, modelId estype.TrainedModelID, opts ...CatMlTrainedModelsOption) (cat_ml_trained_models.Response, error)

	// CatNodeattrs returns information about custom node attributes.
	//
	// Returns information about custom node attributes.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-nodeattrs.html
	CatNodeattrs(ctx context.Context, opts ...CatNodeattrsOption) (cat_nodeattrs.Response, error)

	// CatPendingTasks returns cluster-level changes that have not yet been executed.
	//
	// Returns cluster-level changes that are queued but have not yet been applied.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-pending-tasks.html
	CatPendingTasks(ctx context.Context, opts ...CatPendingTasksOption) (cat_pending_tasks.Response, error)

	// CatPlugins returns information about the plugins that are running on each node.
	//
	// Returns a list of plugins running on each node.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-plugins.html
	CatPlugins(ctx context.Context, opts ...CatPluginsOption) (cat_plugins.Response, error)

	// CatRecovery returns information about ongoing and completed shard recoveries.
	//
	// Returns information about ongoing and completed shard recovery processes.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-recovery.html
	CatRecovery(ctx context.Context, indexName estype.Index, opts ...CatRecoveryOption) (cat_recovery.Response, error)

	// CatRepositories returns the snapshot repositories for a cluster.
	//
	// Returns the list of snapshot repositories registered in the cluster.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-repositories.html
	CatRepositories(ctx context.Context, opts ...CatRepositoriesOption) (cat_repositories.Response, error)

	// CatSegments returns low-level information about the Lucene segments in index shards.
	//
	// Returns low-level information about the Lucene segments in index shards, such as segment count and size.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-segments.html
	CatSegments(ctx context.Context, indexName estype.Index, opts ...CatSegmentsOption) (cat_segments.Response, error)

	// CatShards returns detailed information about shards in the cluster.
	//
	// Returns detailed information about shard allocation across the cluster.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-shards.html
	CatShards(ctx context.Context, indexName estype.Index, opts ...CatShardsOption) (cat_shards.Response, error)

	// CatSnapshots returns information about the snapshots stored in one or more repositories.
	//
	// Returns information about snapshots stored in one or more repositories.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-snapshots.html
	CatSnapshots(ctx context.Context, repository estype.Repository, opts ...CatSnapshotsOption) (cat_snapshots.Response, error)

	// CatTasks returns information about currently executing tasks.
	//
	// Returns information about tasks currently executing in the cluster.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-tasks.html
	CatTasks(ctx context.Context, opts ...CatTasksOption) (cat_tasks.Response, error)

	// CatTemplates returns information about index templates in a cluster.
	//
	// Returns information about index templates in the cluster.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-templates.html
	CatTemplates(ctx context.Context, name estype.Template, opts ...CatTemplatesOption) (cat_templates.Response, error)

	// CatThreadPool returns thread pool statistics for each node in a cluster.
	//
	// Returns thread pool statistics for each node, including queue size, active threads, and rejected operations.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-thread-pool.html
	CatThreadPool(ctx context.Context, opts ...CatThreadPoolOption) (cat_thread_pool.Response, error)

	// CatTransforms returns configuration and usage information about transforms.
	//
	// Returns configuration and usage statistics for transforms.
	//
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/cat-transforms.html
	CatTransforms(ctx context.Context, transformId estype.TransformID, opts ...CatTransformsOption) (cat_transforms.Response, error)
}

// NewClient constructs an ESClient backed by the Elasticsearch v8 typed client.
func NewClient(config es8.Config) (ESClient, error) {
	typedClient, err := es8.NewTypedClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch TypedClient: %w", err)
	}
	return newESClient(typedClient), nil
}

// NewSpecClient constructs an ESClientSpec backed by the Elasticsearch v8 typed client.
// ESClientSpec is a superset of ESClient that additionally exposes every
// Elasticsearch spec-named endpoint as a typed Go method.
func NewSpecClient(config es8.Config) (ESClientSpec, error) {
	typedClient, err := es8.NewTypedClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch TypedClient: %w", err)
	}
	return newESClient(typedClient), nil
}

// ensure compile-time check that *esClient implements ESClient.
var _ ESClient = (*esClient)(nil)
