// Package esv9 provides a wrapper around the Elasticsearch go-elasticsearch v9
// typed client, offering a type-safe Go API with distinct types for Index names,
// Alias names, and other Elasticsearch concepts to prevent misuse.
// Logging is provided via the standard slog package.
//
//go:generate go run ./generator
package esv9

import (
	"context"
	"fmt"
	"time"

	es9 "github.com/elastic/go-elasticsearch/v9"
	cat_aliases "github.com/elastic/go-elasticsearch/v9/typedapi/cat/aliases"
	cat_allocation "github.com/elastic/go-elasticsearch/v9/typedapi/cat/allocation"
	cat_circuit_breaker "github.com/elastic/go-elasticsearch/v9/typedapi/cat/circuitbreaker"
	cat_component_templates "github.com/elastic/go-elasticsearch/v9/typedapi/cat/componenttemplates"
	cat_count "github.com/elastic/go-elasticsearch/v9/typedapi/cat/count"
	cat_fielddata "github.com/elastic/go-elasticsearch/v9/typedapi/cat/fielddata"
	cat_health "github.com/elastic/go-elasticsearch/v9/typedapi/cat/health"
	cat_help "github.com/elastic/go-elasticsearch/v9/typedapi/cat/help"
	cat_indices "github.com/elastic/go-elasticsearch/v9/typedapi/cat/indices"
	cat_master "github.com/elastic/go-elasticsearch/v9/typedapi/cat/master"
	cat_ml_data_frame_analytics "github.com/elastic/go-elasticsearch/v9/typedapi/cat/mldataframeanalytics"
	cat_ml_datafeeds "github.com/elastic/go-elasticsearch/v9/typedapi/cat/mldatafeeds"
	cat_ml_jobs "github.com/elastic/go-elasticsearch/v9/typedapi/cat/mljobs"
	cat_ml_trained_models "github.com/elastic/go-elasticsearch/v9/typedapi/cat/mltrainedmodels"
	cat_nodeattrs "github.com/elastic/go-elasticsearch/v9/typedapi/cat/nodeattrs"
	cat_nodes "github.com/elastic/go-elasticsearch/v9/typedapi/cat/nodes"
	cat_pending_tasks "github.com/elastic/go-elasticsearch/v9/typedapi/cat/pendingtasks"
	cat_plugins "github.com/elastic/go-elasticsearch/v9/typedapi/cat/plugins"
	cat_recovery "github.com/elastic/go-elasticsearch/v9/typedapi/cat/recovery"
	cat_repositories "github.com/elastic/go-elasticsearch/v9/typedapi/cat/repositories"
	cat_segments "github.com/elastic/go-elasticsearch/v9/typedapi/cat/segments"
	cat_shards "github.com/elastic/go-elasticsearch/v9/typedapi/cat/shards"
	cat_snapshots "github.com/elastic/go-elasticsearch/v9/typedapi/cat/snapshots"
	cat_tasks "github.com/elastic/go-elasticsearch/v9/typedapi/cat/tasks"
	cat_templates "github.com/elastic/go-elasticsearch/v9/typedapi/cat/templates"
	cat_thread_pool "github.com/elastic/go-elasticsearch/v9/typedapi/cat/threadpool"
	cat_transforms "github.com/elastic/go-elasticsearch/v9/typedapi/cat/transforms"
	ccr_follow "github.com/elastic/go-elasticsearch/v9/typedapi/ccr/follow"
	ccr_follow_stats "github.com/elastic/go-elasticsearch/v9/typedapi/ccr/followstats"
	ccr_pause_follow "github.com/elastic/go-elasticsearch/v9/typedapi/ccr/pausefollow"
	ccr_resume_follow "github.com/elastic/go-elasticsearch/v9/typedapi/ccr/resumefollow"
	ccr_unfollow "github.com/elastic/go-elasticsearch/v9/typedapi/ccr/unfollow"
	cluster_health "github.com/elastic/go-elasticsearch/v9/typedapi/cluster/health"
	core_bulk "github.com/elastic/go-elasticsearch/v9/typedapi/core/bulk"
	core_clear_scroll "github.com/elastic/go-elasticsearch/v9/typedapi/core/clearscroll"
	core_close_point_in_time "github.com/elastic/go-elasticsearch/v9/typedapi/core/closepointintime"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/count"
	core_count "github.com/elastic/go-elasticsearch/v9/typedapi/core/count"
	core_delete_by_query "github.com/elastic/go-elasticsearch/v9/typedapi/core/deletebyquery"
	coredelete "github.com/elastic/go-elasticsearch/v9/typedapi/core/delete"
	coreget "github.com/elastic/go-elasticsearch/v9/typedapi/core/get"
	coreidx "github.com/elastic/go-elasticsearch/v9/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/info"
	core_mget "github.com/elastic/go-elasticsearch/v9/typedapi/core/mget"
	core_msearch "github.com/elastic/go-elasticsearch/v9/typedapi/core/msearch"
	core_open_point_in_time "github.com/elastic/go-elasticsearch/v9/typedapi/core/openpointintime"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/reindex"
	core_scroll "github.com/elastic/go-elasticsearch/v9/typedapi/core/scroll"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/update"
	core_update_by_query "github.com/elastic/go-elasticsearch/v9/typedapi/core/updatebyquery"
	esql_query "github.com/elastic/go-elasticsearch/v9/typedapi/esql/query"
	ilm_explain_lifecycle "github.com/elastic/go-elasticsearch/v9/typedapi/ilm/explainlifecycle"
	ilm_get_lifecycle "github.com/elastic/go-elasticsearch/v9/typedapi/ilm/getlifecycle"
	ilm_put_lifecycle "github.com/elastic/go-elasticsearch/v9/typedapi/ilm/putlifecycle"
	indices_analyze "github.com/elastic/go-elasticsearch/v9/typedapi/indices/analyze"
	indices_clear_cache "github.com/elastic/go-elasticsearch/v9/typedapi/indices/clearcache"
	indices_close "github.com/elastic/go-elasticsearch/v9/typedapi/indices/close"
	idxcreate "github.com/elastic/go-elasticsearch/v9/typedapi/indices/create"
	indices_create_data_stream "github.com/elastic/go-elasticsearch/v9/typedapi/indices/createdatastream"
	idxdelete "github.com/elastic/go-elasticsearch/v9/typedapi/indices/delete"
	indices_delete_alias "github.com/elastic/go-elasticsearch/v9/typedapi/indices/deletealias"
	indices_delete_data_stream "github.com/elastic/go-elasticsearch/v9/typedapi/indices/deletedatastream"
	indices_delete_index_template "github.com/elastic/go-elasticsearch/v9/typedapi/indices/deleteindextemplate"
	indices_flush "github.com/elastic/go-elasticsearch/v9/typedapi/indices/flush"
	indices_forcemerge "github.com/elastic/go-elasticsearch/v9/typedapi/indices/forcemerge"
	indices_get_alias "github.com/elastic/go-elasticsearch/v9/typedapi/indices/getalias"
	indices_get_data_stream "github.com/elastic/go-elasticsearch/v9/typedapi/indices/getdatastream"
	indices_get_index_template "github.com/elastic/go-elasticsearch/v9/typedapi/indices/getindextemplate"
	indices_get_mapping "github.com/elastic/go-elasticsearch/v9/typedapi/indices/getmapping"
	indices_get_settings "github.com/elastic/go-elasticsearch/v9/typedapi/indices/getsettings"
	indices_open "github.com/elastic/go-elasticsearch/v9/typedapi/indices/open"
	indices_put_index_template "github.com/elastic/go-elasticsearch/v9/typedapi/indices/putindextemplate"
	indices_put_mapping "github.com/elastic/go-elasticsearch/v9/typedapi/indices/putmapping"
	idxputalias "github.com/elastic/go-elasticsearch/v9/typedapi/indices/putalias"
	idxputsettings "github.com/elastic/go-elasticsearch/v9/typedapi/indices/putsettings"
	indices_put_settings "github.com/elastic/go-elasticsearch/v9/typedapi/indices/putsettings"
	idxrefresh "github.com/elastic/go-elasticsearch/v9/typedapi/indices/refresh"
	indices_rollover "github.com/elastic/go-elasticsearch/v9/typedapi/indices/rollover"
	indices_stats "github.com/elastic/go-elasticsearch/v9/typedapi/indices/stats"
	idxupdatealiases "github.com/elastic/go-elasticsearch/v9/typedapi/indices/updatealiases"
	inference_delete "github.com/elastic/go-elasticsearch/v9/typedapi/inference/delete"
	inference_get "github.com/elastic/go-elasticsearch/v9/typedapi/inference/get"
	inference_inference "github.com/elastic/go-elasticsearch/v9/typedapi/inference/inference"
	inference_put "github.com/elastic/go-elasticsearch/v9/typedapi/inference/put"
	ingest_delete_pipeline "github.com/elastic/go-elasticsearch/v9/typedapi/ingest/deletepipeline"
	ingest_get_pipeline "github.com/elastic/go-elasticsearch/v9/typedapi/ingest/getpipeline"
	ingest_put_pipeline "github.com/elastic/go-elasticsearch/v9/typedapi/ingest/putpipeline"
	ml_close_job "github.com/elastic/go-elasticsearch/v9/typedapi/ml/closejob"
	ml_delete_datafeed "github.com/elastic/go-elasticsearch/v9/typedapi/ml/deletedatafeed"
	ml_delete_job "github.com/elastic/go-elasticsearch/v9/typedapi/ml/deletejob"
	ml_get_datafeeds "github.com/elastic/go-elasticsearch/v9/typedapi/ml/getdatafeeds"
	ml_get_jobs "github.com/elastic/go-elasticsearch/v9/typedapi/ml/getjobs"
	ml_open_job "github.com/elastic/go-elasticsearch/v9/typedapi/ml/openjob"
	ml_put_datafeed "github.com/elastic/go-elasticsearch/v9/typedapi/ml/putdatafeed"
	ml_put_job "github.com/elastic/go-elasticsearch/v9/typedapi/ml/putjob"
	ml_start_datafeed "github.com/elastic/go-elasticsearch/v9/typedapi/ml/startdatafeed"
	ml_stop_datafeed "github.com/elastic/go-elasticsearch/v9/typedapi/ml/stopdatafeed"
	security_create_api_key "github.com/elastic/go-elasticsearch/v9/typedapi/security/createapikey"
	security_get_api_key "github.com/elastic/go-elasticsearch/v9/typedapi/security/getapikey"
	security_invalidate_api_key "github.com/elastic/go-elasticsearch/v9/typedapi/security/invalidateapikey"
	snapshot_create "github.com/elastic/go-elasticsearch/v9/typedapi/snapshot/create"
	snapshot_create_repository "github.com/elastic/go-elasticsearch/v9/typedapi/snapshot/createrepository"
	snapshot_restore "github.com/elastic/go-elasticsearch/v9/typedapi/snapshot/restore"
	tasks_cancel "github.com/elastic/go-elasticsearch/v9/typedapi/tasks/cancel"
	tasks_list "github.com/elastic/go-elasticsearch/v9/typedapi/tasks/list"
	transform_delete_transform "github.com/elastic/go-elasticsearch/v9/typedapi/transform/deletetransform"
	transform_get_transform "github.com/elastic/go-elasticsearch/v9/typedapi/transform/gettransform"
	transform_get_transform_stats "github.com/elastic/go-elasticsearch/v9/typedapi/transform/gettransformstats"
	transform_put_transform "github.com/elastic/go-elasticsearch/v9/typedapi/transform/puttransform"
	transform_start_transform "github.com/elastic/go-elasticsearch/v9/typedapi/transform/starttransform"
	transform_stop_transform "github.com/elastic/go-elasticsearch/v9/typedapi/transform/stoptransform"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// ESClient defines the interface for interacting with Elasticsearch v9.
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
	// taskID is the task identifier string (e.g. "node:task_number").
	WaitForTaskCompletion(ctx context.Context, taskID string, timeout time.Duration) error

	// ---------------------------------------------------------------------------
	// Document & Search API
	// ---------------------------------------------------------------------------

	// DocumentExists reports whether a document with the given ID exists in the alias.
	DocumentExists(ctx context.Context, aliasName estype.Alias, id estype.DocumentID, opts ...DocumentExistsOption) (bool, error)

	// Bulk performs multiple index, create, delete, or update operations in a single request.
	Bulk(ctx context.Context, aliasName estype.Alias, opts ...BulkOption) (*core_bulk.Response, error)

	// Mget retrieves multiple documents by ID in a single request.
	Mget(ctx context.Context, aliasName estype.Alias, opts ...MgetOption) (*core_mget.Response, error)

	// Msearch executes multiple search requests in a single request.
	Msearch(ctx context.Context, aliasName estype.Alias, opts ...MsearchOption) (*core_msearch.Response, error)

	// Count returns the number of documents matching a query against the alias.
	Count(ctx context.Context, aliasName estype.Alias, opts ...CountOption) (*core_count.Response, error)

	// Scroll retrieves the next batch of results from a scroll operation.
	Scroll(ctx context.Context, opts ...ScrollOption) (*core_scroll.Response, error)

	// ClearScroll clears the search context and results for a scroll.
	ClearScroll(ctx context.Context, opts ...ClearScrollOption) (*core_clear_scroll.Response, error)

	// UpdateByQuery updates documents that match the given query in the index.
	UpdateByQuery(ctx context.Context, indexName estype.Index, opts ...UpdateByQueryOption) (*core_update_by_query.Response, error)

	// DeleteByQuery deletes documents that match the given query in the index.
	DeleteByQuery(ctx context.Context, indexName estype.Index, opts ...DeleteByQueryOption) (*core_delete_by_query.Response, error)

	// ---------------------------------------------------------------------------
	// Index Management API
	// ---------------------------------------------------------------------------

	// GetMapping returns the mapping for the specified index.
	GetMapping(ctx context.Context, indexName estype.Index) (indices_get_mapping.Response, error)

	// PutMapping updates the field mappings for the specified index.
	PutMapping(ctx context.Context, indexName estype.Index, opts ...PutMappingOption) (*indices_put_mapping.Response, error)

	// GetSettings returns the settings for the specified index.
	GetSettings(ctx context.Context, indexName estype.Index) (indices_get_settings.Response, error)

	// PutSettings updates the settings for the specified index.
	PutSettings(ctx context.Context, indexName estype.Index, opts ...PutSettingsOption) (*indices_put_settings.Response, error)

	// OpenIndex opens a closed index.
	OpenIndex(ctx context.Context, indexName estype.Index) (*indices_open.Response, error)

	// CloseIndex closes an open index.
	CloseIndex(ctx context.Context, indexName estype.Index) (*indices_close.Response, error)

	// Flush flushes one or more indices.
	Flush(ctx context.Context, indexName estype.Index) (*indices_flush.Response, error)

	// ClearCache clears the caches of one or more indices.
	ClearCache(ctx context.Context, indexName estype.Index) (*indices_clear_cache.Response, error)

	// ForceMerge forces a merge on the shards of one or more indices.
	ForceMerge(ctx context.Context, indexName estype.Index) (*indices_forcemerge.Response, error)

	// Rollover rolls an alias over to a new index when the existing index meets a condition.
	Rollover(ctx context.Context, aliasName estype.Alias, opts ...RolloverOption) (*indices_rollover.Response, error)

	// IndicesStats returns statistics for the specified index.
	IndicesStats(ctx context.Context, indexName estype.Index) (*indices_stats.Response, error)

	// ---------------------------------------------------------------------------
	// Cluster API
	// ---------------------------------------------------------------------------

	// ClusterHealth returns the health status of the specified index.
	ClusterHealth(ctx context.Context, indexName estype.Index) (*cluster_health.Response, error)

	// ---------------------------------------------------------------------------
	// Index Template API
	// ---------------------------------------------------------------------------

	// PutIndexTemplate creates or updates an index template.
	PutIndexTemplate(ctx context.Context, name estype.Template, opts ...PutIndexTemplateOption) (*indices_put_index_template.Response, error)

	// GetIndexTemplate returns an index template.
	GetIndexTemplate(ctx context.Context, name estype.Template) (*indices_get_index_template.Response, error)

	// DeleteIndexTemplate deletes an index template.
	DeleteIndexTemplate(ctx context.Context, name estype.Template) (*indices_delete_index_template.Response, error)

	// ExistsIndexTemplate reports whether an index template exists.
	ExistsIndexTemplate(ctx context.Context, name estype.Template) (bool, error)

	// ---------------------------------------------------------------------------
	// Alias API
	// ---------------------------------------------------------------------------

	// GetAlias returns alias information for the specified alias.
	GetAlias(ctx context.Context, aliasName estype.Alias) (indices_get_alias.Response, error)

	// DeleteAlias removes an alias from the specified index.
	DeleteAlias(ctx context.Context, indexName estype.Index, aliasName estype.Alias) (*indices_delete_alias.Response, error)

	// ---------------------------------------------------------------------------
	// Task Management API
	// ---------------------------------------------------------------------------

	// TasksList returns a list of currently running tasks.
	TasksList(ctx context.Context, opts ...TasksListOption) (*tasks_list.Response, error)

	// TasksCancel cancels a task or a group of tasks.
	TasksCancel(ctx context.Context, taskID estype.TaskID) (*tasks_cancel.Response, error)

	// ---------------------------------------------------------------------------
	// Analysis / Debug API
	// ---------------------------------------------------------------------------

	// Analyze performs analysis on a text string and returns the tokens.
	Analyze(ctx context.Context, indexName estype.Index, opts ...AnalyzeOption) (*indices_analyze.Response, error)

	// ---------------------------------------------------------------------------
	// PIT (Point In Time) API
	// ---------------------------------------------------------------------------

	// OpenPointInTime opens a point in time on the specified index.
	OpenPointInTime(ctx context.Context, indexName estype.Index, keepAlive estype.KeepAlive, opts ...OpenPointInTimeOption) (*core_open_point_in_time.Response, error)

	// ClosePointInTime closes a point in time.
	ClosePointInTime(ctx context.Context, opts ...ClosePointInTimeOption) (*core_close_point_in_time.Response, error)

	// ---------------------------------------------------------------------------
	// Data Stream API
	// ---------------------------------------------------------------------------

	// CreateDataStream creates a data stream.
	CreateDataStream(ctx context.Context, name estype.DataStream) (*indices_create_data_stream.Response, error)

	// GetDataStream returns information about one or more data streams.
	GetDataStream(ctx context.Context, name estype.DataStream) (*indices_get_data_stream.Response, error)

	// DeleteDataStream deletes one or more data streams.
	DeleteDataStream(ctx context.Context, name estype.DataStream) (*indices_delete_data_stream.Response, error)

	// ---------------------------------------------------------------------------
	// ILM (Index Lifecycle Management) API
	// ---------------------------------------------------------------------------

	// PutLifecycle creates or updates a lifecycle policy.
	PutLifecycle(ctx context.Context, policy estype.Policy, opts ...PutLifecycleOption) (*ilm_put_lifecycle.Response, error)

	// GetLifecycle returns lifecycle policy information.
	GetLifecycle(ctx context.Context, policy estype.Policy) (ilm_get_lifecycle.Response, error)

	// ExplainLifecycle returns the current lifecycle status for one or more indices.
	ExplainLifecycle(ctx context.Context, indexName estype.Index) (*ilm_explain_lifecycle.Response, error)

	// ---------------------------------------------------------------------------
	// Ingest Pipeline API
	// ---------------------------------------------------------------------------

	// PutPipeline creates or updates an ingest pipeline.
	PutPipeline(ctx context.Context, id estype.Pipeline, opts ...PutPipelineOption) (*ingest_put_pipeline.Response, error)

	// GetPipeline returns information about one or more ingest pipelines.
	GetPipeline(ctx context.Context, id estype.Pipeline) (ingest_get_pipeline.Response, error)

	// DeletePipeline deletes a pipeline.
	DeletePipeline(ctx context.Context, id estype.Pipeline) (*ingest_delete_pipeline.Response, error)

	// ---------------------------------------------------------------------------
	// ES|QL API
	// ---------------------------------------------------------------------------

	// EsqlQuery executes an ES|QL query.
	EsqlQuery(ctx context.Context, query estype.ESQLQuery, opts ...EsqlQueryOption) (esql_query.Response, error)

	// ---------------------------------------------------------------------------
	// Security & API Key API
	// ---------------------------------------------------------------------------

	// CreateApiKey creates an API key.
	CreateApiKey(ctx context.Context, opts ...CreateApiKeyOption) (*security_create_api_key.Response, error)

	// GetApiKey retrieves information for one or more API keys.
	GetApiKey(ctx context.Context, opts ...GetApiKeyOption) (*security_get_api_key.Response, error)

	// InvalidateApiKey invalidates one or more API keys.
	InvalidateApiKey(ctx context.Context, opts ...InvalidateApiKeyOption) (*security_invalidate_api_key.Response, error)

	// ---------------------------------------------------------------------------
	// Snapshot & Restore API
	// ---------------------------------------------------------------------------

	// CreateRepository creates or updates a snapshot repository.
	CreateRepository(ctx context.Context, repo estype.Repository, opts ...CreateRepositoryOption) (*snapshot_create_repository.Response, error)

	// CreateSnapshot creates a snapshot in a repository.
	CreateSnapshot(ctx context.Context, repo estype.Repository, snap estype.Snapshot, opts ...CreateSnapshotOption) (*snapshot_create.Response, error)

	// RestoreSnapshot restores a snapshot.
	RestoreSnapshot(ctx context.Context, repo estype.Repository, snap estype.Snapshot, opts ...RestoreSnapshotOption) (*snapshot_restore.Response, error)

	// ---------------------------------------------------------------------------
	// Inference API
	// ---------------------------------------------------------------------------

	// PutInference creates or updates an inference endpoint.
	PutInference(ctx context.Context, inferenceId estype.InferenceID, opts ...PutInferenceOption) (*inference_put.Response, error)

	// GetInference returns information about an inference endpoint.
	GetInference(ctx context.Context, inferenceId estype.InferenceID) (*inference_get.Response, error)

	// DeleteInference deletes an inference endpoint.
	DeleteInference(ctx context.Context, inferenceId estype.InferenceID) (*inference_delete.Response, error)

	// Inference performs an inference request against an inference endpoint.
	Inference(ctx context.Context, inferenceId estype.InferenceID, opts ...InferenceOption) (*inference_inference.Response, error)

	// ---------------------------------------------------------------------------
	// Machine Learning (ML) API
	// ---------------------------------------------------------------------------

	// MlPutJob creates an anomaly detection job.
	MlPutJob(ctx context.Context, jobId estype.MLJobID, opts ...MlPutJobOption) (*ml_put_job.Response, error)

	// MlGetJobs returns configuration information for anomaly detection jobs.
	MlGetJobs(ctx context.Context, jobId estype.MLJobID) (*ml_get_jobs.Response, error)

	// MlDeleteJob deletes an anomaly detection job.
	MlDeleteJob(ctx context.Context, jobId estype.MLJobID) (*ml_delete_job.Response, error)

	// MlOpenJob opens an anomaly detection job.
	MlOpenJob(ctx context.Context, jobId estype.MLJobID, opts ...MlOpenJobOption) (*ml_open_job.Response, error)

	// MlCloseJob closes an anomaly detection job.
	MlCloseJob(ctx context.Context, jobId estype.MLJobID, opts ...MlCloseJobOption) (*ml_close_job.Response, error)

	// MlPutDatafeed creates a datafeed.
	MlPutDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlPutDatafeedOption) (*ml_put_datafeed.Response, error)

	// MlGetDatafeeds returns configuration information for datafeeds.
	MlGetDatafeeds(ctx context.Context, datafeedId estype.DatafeedID) (*ml_get_datafeeds.Response, error)

	// MlDeleteDatafeed deletes a datafeed.
	MlDeleteDatafeed(ctx context.Context, datafeedId estype.DatafeedID) (*ml_delete_datafeed.Response, error)

	// MlStartDatafeed starts one or more datafeeds.
	MlStartDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlStartDatafeedOption) (*ml_start_datafeed.Response, error)

	// MlStopDatafeed stops one or more datafeeds.
	MlStopDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlStopDatafeedOption) (*ml_stop_datafeed.Response, error)

	// ---------------------------------------------------------------------------
	// CCR (Cross-Cluster Replication) API
	// ---------------------------------------------------------------------------

	// CcrFollow configures a local index to follow a remote index.
	CcrFollow(ctx context.Context, followerIndex estype.Index, opts ...CcrFollowOption) (*ccr_follow.Response, error)

	// CcrPauseFollow pauses a follower index.
	CcrPauseFollow(ctx context.Context, followerIndex estype.Index) (*ccr_pause_follow.Response, error)

	// CcrResumeFollow resumes a follower index.
	CcrResumeFollow(ctx context.Context, followerIndex estype.Index, opts ...CcrResumeFollowOption) (*ccr_resume_follow.Response, error)

	// CcrUnfollow stops the following task for a follower index.
	CcrUnfollow(ctx context.Context, followerIndex estype.Index) (*ccr_unfollow.Response, error)

	// CcrFollowStats returns cross-cluster replication follower stats.
	CcrFollowStats(ctx context.Context, followerIndex estype.Index) (*ccr_follow_stats.Response, error)

	// ---------------------------------------------------------------------------
	// Transform API
	// ---------------------------------------------------------------------------

	// PutTransform creates or updates a transform.
	PutTransform(ctx context.Context, transformId estype.TransformID, opts ...PutTransformOption) (*transform_put_transform.Response, error)

	// GetTransform returns configuration information for transforms.
	GetTransform(ctx context.Context, transformId estype.TransformID) (*transform_get_transform.Response, error)

	// DeleteTransform deletes a transform.
	DeleteTransform(ctx context.Context, transformId estype.TransformID) (*transform_delete_transform.Response, error)

	// StartTransform starts one or more transforms.
	StartTransform(ctx context.Context, transformId estype.TransformID) (*transform_start_transform.Response, error)

	// StopTransform stops one or more transforms.
	StopTransform(ctx context.Context, transformId estype.TransformID) (*transform_stop_transform.Response, error)

	// GetTransformStats returns usage information for transforms.
	GetTransformStats(ctx context.Context, transformId estype.TransformID) (*transform_get_transform_stats.Response, error)

	// ---------------------------------------------------------------------------
	// Cat API
	// ---------------------------------------------------------------------------

	// CatAliases returns information about aliases.
	CatAliases(ctx context.Context, aliasName estype.Alias) (cat_aliases.Response, error)

	// CatIndices returns high-level information about indices.
	CatIndices(ctx context.Context, indexName estype.Index) (cat_indices.Response, error)

	// CatNodes returns information about the nodes in a cluster.
	CatNodes(ctx context.Context) (cat_nodes.Response, error)

	// CatAllocation provides a snapshot of the number of shards allocated to each data node.
	CatAllocation(ctx context.Context, opts ...CatAllocationOption) (cat_allocation.Response, error)

	// CatCircuitBreaker returns information about circuit breaker status.
	CatCircuitBreaker(ctx context.Context, opts ...CatCircuitBreakerOption) (cat_circuit_breaker.Response, error)

	// CatComponentTemplates returns information about component templates.
	CatComponentTemplates(ctx context.Context, name estype.Template, opts ...CatComponentTemplatesOption) (cat_component_templates.Response, error)

	// CatCount returns document counts for one or more indices.
	CatCount(ctx context.Context, indexName estype.Index, opts ...CatCountOption) (cat_count.Response, error)

	// CatFielddata returns information about the amount of memory used for field data.
	CatFielddata(ctx context.Context, fields []estype.Field, opts ...CatFielddataOption) (cat_fielddata.Response, error)

	// CatHealth returns a concise representation of cluster health.
	CatHealth(ctx context.Context, opts ...CatHealthOption) (cat_health.Response, error)

	// CatHelp returns help for the cat APIs.
	CatHelp(ctx context.Context, opts ...CatHelpOption) (*cat_help.Response, error)

	// CatMaster returns information about the master node.
	CatMaster(ctx context.Context, opts ...CatMasterOption) (cat_master.Response, error)

	// CatMlDataFrameAnalytics returns configuration and usage information about data frame analytics jobs.
	CatMlDataFrameAnalytics(ctx context.Context, analyticsId estype.DataFrameAnalyticsID, opts ...CatMlDataFrameAnalyticsOption) (cat_ml_data_frame_analytics.Response, error)

	// CatMlDatafeeds returns configuration and usage information about datafeeds.
	CatMlDatafeeds(ctx context.Context, datafeedId estype.DatafeedID, opts ...CatMlDatafeedsOption) (cat_ml_datafeeds.Response, error)

	// CatMlJobs returns configuration and usage information about anomaly detection jobs.
	CatMlJobs(ctx context.Context, jobId estype.MLJobID, opts ...CatMlJobsOption) (cat_ml_jobs.Response, error)

	// CatMlTrainedModels returns configuration and usage information about inference trained models.
	CatMlTrainedModels(ctx context.Context, modelId estype.TrainedModelID, opts ...CatMlTrainedModelsOption) (cat_ml_trained_models.Response, error)

	// CatNodeattrs returns information about custom node attributes.
	CatNodeattrs(ctx context.Context, opts ...CatNodeattrsOption) (cat_nodeattrs.Response, error)

	// CatPendingTasks returns cluster-level changes that have not yet been executed.
	CatPendingTasks(ctx context.Context, opts ...CatPendingTasksOption) (cat_pending_tasks.Response, error)

	// CatPlugins returns information about the plugins that are running on each node.
	CatPlugins(ctx context.Context, opts ...CatPluginsOption) (cat_plugins.Response, error)

	// CatRecovery returns information about ongoing and completed shard recoveries.
	CatRecovery(ctx context.Context, indexName estype.Index, opts ...CatRecoveryOption) (cat_recovery.Response, error)

	// CatRepositories returns the snapshot repositories for a cluster.
	CatRepositories(ctx context.Context, opts ...CatRepositoriesOption) (cat_repositories.Response, error)

	// CatSegments returns low-level information about the Lucene segments in index shards.
	CatSegments(ctx context.Context, indexName estype.Index, opts ...CatSegmentsOption) (cat_segments.Response, error)

	// CatShards returns detailed information about shards in the cluster.
	CatShards(ctx context.Context, indexName estype.Index, opts ...CatShardsOption) (cat_shards.Response, error)

	// CatSnapshots returns information about the snapshots stored in one or more repositories.
	CatSnapshots(ctx context.Context, repository estype.Repository, opts ...CatSnapshotsOption) (cat_snapshots.Response, error)

	// CatTasks returns information about currently executing tasks.
	CatTasks(ctx context.Context, opts ...CatTasksOption) (cat_tasks.Response, error)

	// CatTemplates returns information about index templates in a cluster.
	CatTemplates(ctx context.Context, name estype.Template, opts ...CatTemplatesOption) (cat_templates.Response, error)

	// CatThreadPool returns thread pool statistics for each node in a cluster.
	CatThreadPool(ctx context.Context, opts ...CatThreadPoolOption) (cat_thread_pool.Response, error)

	// CatTransforms returns configuration and usage information about transforms.
	CatTransforms(ctx context.Context, transformId estype.TransformID, opts ...CatTransformsOption) (cat_transforms.Response, error)
}

// NewClient constructs an ESClient backed by the Elasticsearch v9 typed client.
func NewClient(config es9.Config) (ESClient, error) {
	typedClient, err := es9.NewTypedClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch TypedClient: %w", err)
	}
	return newESClient(typedClient), nil
}

// NewSpecClient constructs an ESClientSpec backed by the Elasticsearch v9 typed client.
// ESClientSpec is a superset of ESClient that additionally exposes every
// Elasticsearch spec-named endpoint as a typed Go method.
func NewSpecClient(config es9.Config) (ESClientSpec, error) {
	typedClient, err := es9.NewTypedClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch TypedClient: %w", err)
	}
	return newESClient(typedClient), nil
}

// ensure compile-time check that *esClient implements ESClient.
var _ ESClient = (*esClient)(nil)
