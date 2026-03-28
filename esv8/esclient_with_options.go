package esv8

import (
	coredelete "github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	coreget "github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	coreidx "github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	core_bulk "github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	core_delete_by_query "github.com/elastic/go-elasticsearch/v8/typedapi/core/deletebyquery"
	core_clear_scroll "github.com/elastic/go-elasticsearch/v8/typedapi/core/clearscroll"
	core_scroll "github.com/elastic/go-elasticsearch/v8/typedapi/core/scroll"
	core_update_by_query "github.com/elastic/go-elasticsearch/v8/typedapi/core/updatebyquery"
	cluster_health "github.com/elastic/go-elasticsearch/v8/typedapi/cluster/health"
	idxdelete "github.com/elastic/go-elasticsearch/v8/typedapi/indices/delete"
	indices_clear_cache "github.com/elastic/go-elasticsearch/v8/typedapi/indices/clearcache"
	indices_forcemerge "github.com/elastic/go-elasticsearch/v8/typedapi/indices/forcemerge"
	indices_stats "github.com/elastic/go-elasticsearch/v8/typedapi/indices/stats"
	idxputalias "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putalias"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	ilm_explain_lifecycle "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/explainlifecycle"
	ml_delete_datafeed "github.com/elastic/go-elasticsearch/v8/typedapi/ml/deletedatafeed"
	ml_delete_job "github.com/elastic/go-elasticsearch/v8/typedapi/ml/deletejob"
	ml_get_datafeeds "github.com/elastic/go-elasticsearch/v8/typedapi/ml/getdatafeeds"
	ml_get_jobs "github.com/elastic/go-elasticsearch/v8/typedapi/ml/getjobs"
	tasks_cancel "github.com/elastic/go-elasticsearch/v8/typedapi/tasks/cancel"
	transform_get_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/gettransform"
	transform_stop_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/stoptransform"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/conflicts"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/healthstatus"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/level"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/optype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/waitforevents"
)

// CreateDocument options

// WithRouting sets the routing value for a CreateDocument request.
func WithRouting(routing string) CreateDocumentOption {
	return func(r *coreidx.Index) { r.Routing(routing) }
}

// WithPipeline sets the ingest pipeline for a CreateDocument request.
func WithPipeline(pipeline string) CreateDocumentOption {
	return func(r *coreidx.Index) { r.Pipeline(pipeline) }
}

// WithRefresh sets the refresh policy for a CreateDocument request.
func WithRefresh(v refresh.Refresh) CreateDocumentOption {
	return func(r *coreidx.Index) { r.Refresh(v) }
}

// WithOpType sets the op_type for a CreateDocument request.
func WithOpType(op optype.OpType) CreateDocumentOption {
	return func(r *coreidx.Index) { r.OpType(op) }
}

// WithDocumentTimeout sets the timeout for a CreateDocument request.
func WithDocumentTimeout(t string) CreateDocumentOption {
	return func(r *coreidx.Index) { r.Timeout(t) }
}

// GetDocument options

// WithGetRouting sets the routing value for a GetDocument request.
func WithGetRouting(routing string) GetDocumentOption {
	return func(r *coreget.Get) { r.Routing(routing) }
}

// WithSourceIncludes specifies fields to include in _source for a GetDocument request.
func WithSourceIncludes(fields ...string) GetDocumentOption {
	return func(r *coreget.Get) { r.SourceIncludes_(fields...) }
}

// WithSourceExcludes specifies fields to exclude from _source for a GetDocument request.
func WithSourceExcludes(fields ...string) GetDocumentOption {
	return func(r *coreget.Get) { r.SourceExcludes_(fields...) }
}

// WithStoredFields specifies stored fields to return for a GetDocument request.
func WithStoredFields(fields ...string) GetDocumentOption {
	return func(r *coreget.Get) { r.StoredFields(fields...) }
}

// DeleteDocument options

// WithDeleteRouting sets the routing value for a DeleteDocument request.
func WithDeleteRouting(routing string) DeleteDocumentOption {
	return func(r *coredelete.Delete) { r.Routing(routing) }
}

// WithDeleteRefresh sets the refresh policy for a DeleteDocument request.
func WithDeleteRefresh(v refresh.Refresh) DeleteDocumentOption {
	return func(r *coredelete.Delete) { r.Refresh(v) }
}

// WithDeleteTimeout sets the timeout for a DeleteDocument request.
func WithDeleteTimeout(t string) DeleteDocumentOption {
	return func(r *coredelete.Delete) { r.Timeout(t) }
}

// Bulk options

// WithBulkRefresh sets the refresh policy for a Bulk request.
func WithBulkRefresh(v refresh.Refresh) BulkOption {
	return func(b *core_bulk.Bulk) { b.Refresh(v) }
}

// WithBulkPipeline sets the ingest pipeline for a Bulk request.
func WithBulkPipeline(p string) BulkOption {
	return func(b *core_bulk.Bulk) { b.Pipeline(p) }
}

// WithBulkRouting sets the routing value for a Bulk request.
func WithBulkRouting(routing string) BulkOption {
	return func(b *core_bulk.Bulk) { b.Routing(routing) }
}

// WithBulkTimeout sets the timeout for a Bulk request.
func WithBulkTimeout(t string) BulkOption {
	return func(b *core_bulk.Bulk) { b.Timeout(t) }
}

// Scroll options

// WithScrollId sets the scroll ID for a Scroll request.
func WithScrollId(id string) ScrollOption {
	return func(r *core_scroll.Scroll) { r.ScrollId(id) }
}

// ClearScroll options

// WithClearScrollId sets the scroll ID(s) to clear for a ClearScroll request.
// Pass "_all" to clear all active scroll contexts.
func WithClearScrollId(ids ...string) ClearScrollOption {
	return func(r *core_clear_scroll.ClearScroll) { r.ScrollId(ids...) }
}

// UpdateByQuery options

// WithUpdateSlices sets the number of slices for an UpdateByQuery request.
func WithUpdateSlices(s string) UpdateByQueryOption {
	return func(r *core_update_by_query.UpdateByQuery) { r.Slices(s) }
}

// WithUpdateConflicts sets the conflict handling for an UpdateByQuery request.
func WithUpdateConflicts(c conflicts.Conflicts) UpdateByQueryOption {
	return func(r *core_update_by_query.UpdateByQuery) { r.Conflicts(c) }
}

// WithUpdateWaitForCompletion sets whether to wait for completion for an UpdateByQuery request.
func WithUpdateWaitForCompletion(b bool) UpdateByQueryOption {
	return func(r *core_update_by_query.UpdateByQuery) { r.WaitForCompletion(b) }
}

// WithUpdateMaxDocs limits the number of documents to update for an UpdateByQuery request.
func WithUpdateMaxDocs(n int64) UpdateByQueryOption {
	return func(r *core_update_by_query.UpdateByQuery) { r.MaxDocs(n) }
}

// DeleteByQuery options

// WithDeleteSlices sets the number of slices for a DeleteByQuery request.
func WithDeleteSlices(s string) DeleteByQueryOption {
	return func(r *core_delete_by_query.DeleteByQuery) { r.Slices(s) }
}

// WithDeleteConflicts sets the conflict handling for a DeleteByQuery request.
func WithDeleteConflicts(c conflicts.Conflicts) DeleteByQueryOption {
	return func(r *core_delete_by_query.DeleteByQuery) { r.Conflicts(c) }
}

// WithDeleteWaitForCompletion sets whether to wait for completion for a DeleteByQuery request.
func WithDeleteWaitForCompletion(b bool) DeleteByQueryOption {
	return func(r *core_delete_by_query.DeleteByQuery) { r.WaitForCompletion(b) }
}

// WithDeleteMaxDocs limits the number of documents to delete for a DeleteByQuery request.
func WithDeleteMaxDocs(n int64) DeleteByQueryOption {
	return func(r *core_delete_by_query.DeleteByQuery) { r.MaxDocs(n) }
}

// ClusterHealth options

// WithWaitForStatus waits until the cluster health reaches the given status.
func WithWaitForStatus(status healthstatus.HealthStatus) ClusterHealthOption {
	return func(r *cluster_health.Health) { r.WaitForStatus(status) }
}

// WithWaitForNodes waits until the specified number of nodes are available.
func WithWaitForNodes(nodes string) ClusterHealthOption {
	return func(r *cluster_health.Health) { r.WaitForNodes(nodes) }
}

// WithWaitForActiveShards waits until the specified number of shards are active.
func WithWaitForActiveShards(shards string) ClusterHealthOption {
	return func(r *cluster_health.Health) { r.WaitForActiveShards(shards) }
}

// WithHealthTimeout sets the timeout for a ClusterHealth request.
func WithHealthTimeout(t string) ClusterHealthOption {
	return func(r *cluster_health.Health) { r.Timeout(t) }
}

// WithWaitForEvents waits until all currently queued events at the given priority level are processed.
func WithWaitForEvents(e waitforevents.WaitForEvents) ClusterHealthOption {
	return func(r *cluster_health.Health) { r.WaitForEvents(e) }
}

// CreateAlias options

// WithCreateAliasIsWriteIndex sets is_write_index for a CreateAlias request.
// If true the index becomes the write target for the alias.
// If false the index is explicitly excluded as the write target.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-add-alias.html
func WithCreateAliasIsWriteIndex(v bool) CreateAliasOption {
	return func(b *idxputalias.PutAlias) { b.IsWriteIndex(v) }
}

// WithCreateAliasFilter sets the filter query for a CreateAlias request.
// The alias will only expose documents matching this query.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-add-alias.html
func WithCreateAliasFilter(filter *types.Query) CreateAliasOption {
	return func(b *idxputalias.PutAlias) { b.Filter(filter) }
}

// WithCreateAliasRouting sets the routing value for both indexing and search on a CreateAlias request.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-add-alias.html
func WithCreateAliasRouting(routing string) CreateAliasOption {
	return func(b *idxputalias.PutAlias) { b.Routing(routing) }
}

// WithCreateAliasIndexRouting sets the routing value for indexing operations on a CreateAlias request.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-add-alias.html
func WithCreateAliasIndexRouting(routing string) CreateAliasOption {
	return func(b *idxputalias.PutAlias) { b.IndexRouting(routing) }
}

// WithCreateAliasSearchRouting sets the routing value for search operations on a CreateAlias request.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-add-alias.html
func WithCreateAliasSearchRouting(routing string) CreateAliasOption {
	return func(b *idxputalias.PutAlias) { b.SearchRouting(routing) }
}

// WithCreateAliasMasterTimeout sets the master_timeout for a CreateAlias request.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-add-alias.html
func WithCreateAliasMasterTimeout(timeout string) CreateAliasOption {
	return func(b *idxputalias.PutAlias) { b.MasterTimeout(timeout) }
}

// WithCreateAliasTimeout sets the timeout for a CreateAlias request.
//
// https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-add-alias.html
func WithCreateAliasTimeout(timeout string) CreateAliasOption {
	return func(b *idxputalias.PutAlias) { b.Timeout(timeout) }
}

// DeleteIndex options

// WithIgnoreUnavailable ignores unavailable indices for a DeleteIndex request.
func WithIgnoreUnavailable(b bool) DeleteIndexOption {
	return func(r *idxdelete.Delete) { r.IgnoreUnavailable(b) }
}

// WithDeleteIndexTimeout sets the timeout for a DeleteIndex request.
func WithDeleteIndexTimeout(t string) DeleteIndexOption {
	return func(r *idxdelete.Delete) { r.Timeout(t) }
}

// ForceMerge options

// WithMaxNumSegments sets the maximum number of segments to merge to for a ForceMerge request.
func WithMaxNumSegments(n string) ForceMergeOption {
	return func(r *indices_forcemerge.Forcemerge) { r.MaxNumSegments(n) }
}

// WithOnlyExpungeDeletes only expunges segments with deletes for a ForceMerge request.
func WithOnlyExpungeDeletes(b bool) ForceMergeOption {
	return func(r *indices_forcemerge.Forcemerge) { r.OnlyExpungeDeletes(b) }
}

// WithForceMergeFlush flushes after the force merge for a ForceMerge request.
func WithForceMergeFlush(b bool) ForceMergeOption {
	return func(r *indices_forcemerge.Forcemerge) { r.Flush(b) }
}

// ClearCache options

// WithFielddataCache clears the fielddata cache for a ClearCache request.
func WithFielddataCache(b bool) ClearCacheOption {
	return func(r *indices_clear_cache.ClearCache) { r.Fielddata(b) }
}

// WithQueryCache clears the query cache for a ClearCache request.
func WithQueryCache(b bool) ClearCacheOption {
	return func(r *indices_clear_cache.ClearCache) { r.Query(b) }
}

// WithRequestCache clears the request cache for a ClearCache request.
func WithRequestCache(b bool) ClearCacheOption {
	return func(r *indices_clear_cache.ClearCache) { r.Request(b) }
}

// WithCacheFields specifies fields whose caches to clear for a ClearCache request.
func WithCacheFields(fields ...string) ClearCacheOption {
	return func(r *indices_clear_cache.ClearCache) { r.Fields(fields...) }
}

// IndicesStats options

// WithStatsLevel sets the level of detail for an IndicesStats request.
func WithStatsLevel(l level.Level) IndicesStatsOption {
	return func(r *indices_stats.Stats) { r.Level(l) }
}

// ExplainLifecycle options

// WithOnlyManaged filters to only managed indices for an ExplainLifecycle request.
func WithOnlyManaged(b bool) ExplainLifecycleOption {
	return func(r *ilm_explain_lifecycle.ExplainLifecycle) { r.OnlyManaged(b) }
}

// WithOnlyErrors filters to only indices with errors for an ExplainLifecycle request.
func WithOnlyErrors(b bool) ExplainLifecycleOption {
	return func(r *ilm_explain_lifecycle.ExplainLifecycle) { r.OnlyErrors(b) }
}

// GetTransform options

// WithTransformAllowNoMatch allows no matching transforms without error for a GetTransform request.
func WithTransformAllowNoMatch(b bool) GetTransformOption {
	return func(r *transform_get_transform.GetTransform) { r.AllowNoMatch(b) }
}

// WithTransformFrom sets the starting offset for a GetTransform request.
func WithTransformFrom(n int) GetTransformOption {
	return func(r *transform_get_transform.GetTransform) { r.From(n) }
}

// WithTransformSize sets the maximum number of transforms to return for a GetTransform request.
func WithTransformSize(n int) GetTransformOption {
	return func(r *transform_get_transform.GetTransform) { r.Size(n) }
}

// StopTransform options

// WithStopTransformWaitForCompletion waits for the transform to stop for a StopTransform request.
func WithStopTransformWaitForCompletion(b bool) StopTransformOption {
	return func(r *transform_stop_transform.StopTransform) { r.WaitForCompletion(b) }
}

// WithStopTransformForce stops the transform immediately for a StopTransform request.
func WithStopTransformForce(b bool) StopTransformOption {
	return func(r *transform_stop_transform.StopTransform) { r.Force(b) }
}

// WithStopTransformTimeout sets the timeout for a StopTransform request.
func WithStopTransformTimeout(t string) StopTransformOption {
	return func(r *transform_stop_transform.StopTransform) { r.Timeout(t) }
}

// ML options

// WithMlDeleteJobForce forcefully deletes a job for an MlDeleteJob request.
func WithMlDeleteJobForce(b bool) MlDeleteJobOption {
	return func(r *ml_delete_job.DeleteJob) { r.Force(b) }
}

// WithMlDeleteJobWaitForCompletion waits for the job deletion to complete for an MlDeleteJob request.
func WithMlDeleteJobWaitForCompletion(b bool) MlDeleteJobOption {
	return func(r *ml_delete_job.DeleteJob) { r.WaitForCompletion(b) }
}

// WithMlDeleteDatafeedForce forcefully deletes a datafeed for an MlDeleteDatafeed request.
func WithMlDeleteDatafeedForce(b bool) MlDeleteDatafeedOption {
	return func(r *ml_delete_datafeed.DeleteDatafeed) { r.Force(b) }
}

// WithMlJobAllowNoMatch allows no matching jobs without error for an MlGetJobs request.
func WithMlJobAllowNoMatch(b bool) MlGetJobsOption {
	return func(r *ml_get_jobs.GetJobs) { r.AllowNoMatch(b) }
}

// WithMlDatafeedAllowNoMatch allows no matching datafeeds without error for an MlGetDatafeeds request.
func WithMlDatafeedAllowNoMatch(b bool) MlGetDatafeedsOption {
	return func(r *ml_get_datafeeds.GetDatafeeds) { r.AllowNoMatch(b) }
}

// TasksCancel options

// WithTasksWaitForCompletion waits for the cancelled task to complete for a TasksCancel request.
func WithTasksWaitForCompletion(b bool) TasksCancelOption {
	return func(r *tasks_cancel.Cancel) { r.WaitForCompletion(b) }
}
