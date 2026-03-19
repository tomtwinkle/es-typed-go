package esv8

import (
	"context"
	"encoding/json"
	"strings"

	cat_aliases "github.com/elastic/go-elasticsearch/v8/typedapi/cat/aliases"
	cat_allocation "github.com/elastic/go-elasticsearch/v8/typedapi/cat/allocation"
	cat_component_templates "github.com/elastic/go-elasticsearch/v8/typedapi/cat/componenttemplates"
	cat_count "github.com/elastic/go-elasticsearch/v8/typedapi/cat/count"
	cat_fielddata "github.com/elastic/go-elasticsearch/v8/typedapi/cat/fielddata"
	cat_health "github.com/elastic/go-elasticsearch/v8/typedapi/cat/health"
	cat_help "github.com/elastic/go-elasticsearch/v8/typedapi/cat/help"
	cat_indices "github.com/elastic/go-elasticsearch/v8/typedapi/cat/indices"
	cat_master "github.com/elastic/go-elasticsearch/v8/typedapi/cat/master"
	cat_ml_data_frame_analytics "github.com/elastic/go-elasticsearch/v8/typedapi/cat/mldataframeanalytics"
	cat_ml_datafeeds "github.com/elastic/go-elasticsearch/v8/typedapi/cat/mldatafeeds"
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
	core_count "github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	core_delete_by_query "github.com/elastic/go-elasticsearch/v8/typedapi/core/deletebyquery"
	core_mget "github.com/elastic/go-elasticsearch/v8/typedapi/core/mget"
	core_msearch "github.com/elastic/go-elasticsearch/v8/typedapi/core/msearch"
	core_open_point_in_time "github.com/elastic/go-elasticsearch/v8/typedapi/core/openpointintime"
	core_scroll "github.com/elastic/go-elasticsearch/v8/typedapi/core/scroll"
	core_update_by_query "github.com/elastic/go-elasticsearch/v8/typedapi/core/updatebyquery"
	esql_query "github.com/elastic/go-elasticsearch/v8/typedapi/esql/query"
	ilm_explain_lifecycle "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/explainlifecycle"
	ilm_get_lifecycle "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/getlifecycle"
	ilm_put_lifecycle "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/putlifecycle"
	indices_analyze "github.com/elastic/go-elasticsearch/v8/typedapi/indices/analyze"
	indices_clear_cache "github.com/elastic/go-elasticsearch/v8/typedapi/indices/clearcache"
	indices_close "github.com/elastic/go-elasticsearch/v8/typedapi/indices/close"
	indices_create_data_stream "github.com/elastic/go-elasticsearch/v8/typedapi/indices/createdatastream"
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
	indices_put_index_template "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putindextemplate"
	indices_put_mapping "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putmapping"
	indices_put_settings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putsettings"
	indices_rollover "github.com/elastic/go-elasticsearch/v8/typedapi/indices/rollover"
	indices_stats "github.com/elastic/go-elasticsearch/v8/typedapi/indices/stats"
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
	"github.com/tomtwinkle/es-typed-go/estype"
)

// ---------------------------------------------------------------------------
// Document & Search API
// ---------------------------------------------------------------------------

func (c *esClient) DocumentExists(ctx context.Context, aliasName estype.Alias, id estype.DocumentID, opts ...DocumentExistsOption) (bool, error) {
	r := c.typedClient.Exists(aliasName.String(), id.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) Bulk(ctx context.Context, aliasName estype.Alias, opts ...BulkOption) (*core_bulk.Response, error) {
	r := c.typedClient.Bulk().Index(aliasName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) Mget(ctx context.Context, aliasName estype.Alias, opts ...MgetOption) (*core_mget.Response, error) {
	r := c.typedClient.Mget().Index(aliasName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) Msearch(ctx context.Context, aliasName estype.Alias, opts ...MsearchOption) (*core_msearch.Response, error) {
	r := c.typedClient.Msearch().Index(aliasName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) Count(ctx context.Context, aliasName estype.Alias, opts ...CountOption) (*core_count.Response, error) {
	r := c.typedClient.Count().Index(aliasName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) Scroll(ctx context.Context, opts ...ScrollOption) (*core_scroll.Response, error) {
	r := c.typedClient.Scroll()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) ClearScroll(ctx context.Context, opts ...ClearScrollOption) (*core_clear_scroll.Response, error) {
	r := c.typedClient.ClearScroll()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) UpdateByQuery(ctx context.Context, indexName estype.Index, opts ...UpdateByQueryOption) (*core_update_by_query.Response, error) {
	r := c.typedClient.UpdateByQuery(indexName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) DeleteByQuery(ctx context.Context, indexName estype.Index, opts ...DeleteByQueryOption) (*core_delete_by_query.Response, error) {
	r := c.typedClient.DeleteByQuery(indexName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

// ---------------------------------------------------------------------------
// Index Management API
// ---------------------------------------------------------------------------

func (c *esClient) GetMapping(ctx context.Context, indexName estype.Index) (indices_get_mapping.Response, error) {
	return c.typedClient.Indices.GetMapping().Index(indexName.String()).Do(ctx)
}

func (c *esClient) PutMapping(ctx context.Context, indexName estype.Index, opts ...PutMappingOption) (*indices_put_mapping.Response, error) {
	r := c.typedClient.Indices.PutMapping(indexName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) GetSettings(ctx context.Context, indexName estype.Index) (indices_get_settings.Response, error) {
	return c.typedClient.Indices.GetSettings().Index(indexName.String()).Do(ctx)
}

func (c *esClient) PutSettings(ctx context.Context, indexName estype.Index, opts ...PutSettingsOption) (*indices_put_settings.Response, error) {
	r := c.typedClient.Indices.PutSettings().Indices(indexName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) OpenIndex(ctx context.Context, indexName estype.Index) (*indices_open.Response, error) {
	return c.typedClient.Indices.Open(indexName.String()).Do(ctx)
}

func (c *esClient) CloseIndex(ctx context.Context, indexName estype.Index) (*indices_close.Response, error) {
	return c.typedClient.Indices.Close(indexName.String()).Do(ctx)
}

func (c *esClient) Flush(ctx context.Context, indexName estype.Index) (*indices_flush.Response, error) {
	return c.typedClient.Indices.Flush().Index(indexName.String()).Do(ctx)
}

func (c *esClient) ClearCache(ctx context.Context, indexName estype.Index) (*indices_clear_cache.Response, error) {
	return c.typedClient.Indices.ClearCache().Index(indexName.String()).Do(ctx)
}

func (c *esClient) ForceMerge(ctx context.Context, indexName estype.Index) (*indices_forcemerge.Response, error) {
	return c.typedClient.Indices.Forcemerge().Index(indexName.String()).Do(ctx)
}

func (c *esClient) Rollover(ctx context.Context, aliasName estype.Alias, opts ...RolloverOption) (*indices_rollover.Response, error) {
	r := c.typedClient.Indices.Rollover(aliasName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) IndicesStats(ctx context.Context, indexName estype.Index) (*indices_stats.Response, error) {
	return c.typedClient.Indices.Stats().Index(indexName.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// Cluster API
// ---------------------------------------------------------------------------

func (c *esClient) ClusterHealth(ctx context.Context, indexName estype.Index) (*cluster_health.Response, error) {
	return c.typedClient.Cluster.Health().Index(indexName.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// Index Template API
// ---------------------------------------------------------------------------

func (c *esClient) PutIndexTemplate(ctx context.Context, name estype.Template, opts ...PutIndexTemplateOption) (*indices_put_index_template.Response, error) {
	r := c.typedClient.Indices.PutIndexTemplate(name.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) GetIndexTemplate(ctx context.Context, name estype.Template) (*indices_get_index_template.Response, error) {
	return c.typedClient.Indices.GetIndexTemplate().Name(name.String()).Do(ctx)
}

func (c *esClient) DeleteIndexTemplate(ctx context.Context, name estype.Template) (*indices_delete_index_template.Response, error) {
	return c.typedClient.Indices.DeleteIndexTemplate(name.String()).Do(ctx)
}

func (c *esClient) ExistsIndexTemplate(ctx context.Context, name estype.Template) (bool, error) {
	return c.typedClient.Indices.ExistsIndexTemplate(name.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// Alias API
// ---------------------------------------------------------------------------

func (c *esClient) GetAlias(ctx context.Context, aliasName estype.Alias) (indices_get_alias.Response, error) {
	return c.typedClient.Indices.GetAlias().Name(aliasName.String()).Do(ctx)
}

func (c *esClient) DeleteAlias(ctx context.Context, indexName estype.Index, aliasName estype.Alias) (*indices_delete_alias.Response, error) {
	return c.typedClient.Indices.DeleteAlias(indexName.String(), aliasName.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// Task Management API
// ---------------------------------------------------------------------------

func (c *esClient) TasksList(ctx context.Context, opts ...TasksListOption) (*tasks_list.Response, error) {
	r := c.typedClient.Tasks.List()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) TasksCancel(ctx context.Context, taskID estype.TaskID) (*tasks_cancel.Response, error) {
	return c.typedClient.Tasks.Cancel().TaskId(taskID.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// Analysis / Debug API
// ---------------------------------------------------------------------------

func (c *esClient) Analyze(ctx context.Context, indexName estype.Index, opts ...AnalyzeOption) (*indices_analyze.Response, error) {
	r := c.typedClient.Indices.Analyze().Index(indexName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

// ---------------------------------------------------------------------------
// PIT (Point In Time) API
// ---------------------------------------------------------------------------

func (c *esClient) OpenPointInTime(ctx context.Context, indexName estype.Index, keepAlive estype.KeepAlive, opts ...OpenPointInTimeOption) (*core_open_point_in_time.Response, error) {
	r := c.typedClient.OpenPointInTime(indexName.String()).KeepAlive(keepAlive.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) ClosePointInTime(ctx context.Context, opts ...ClosePointInTimeOption) (*core_close_point_in_time.Response, error) {
	r := c.typedClient.ClosePointInTime()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

// ---------------------------------------------------------------------------
// Data Stream API
// ---------------------------------------------------------------------------

func (c *esClient) CreateDataStream(ctx context.Context, name estype.DataStream) (*indices_create_data_stream.Response, error) {
	return c.typedClient.Indices.CreateDataStream(name.String()).Do(ctx)
}

func (c *esClient) GetDataStream(ctx context.Context, name estype.DataStream) (*indices_get_data_stream.Response, error) {
	return c.typedClient.Indices.GetDataStream().Name(name.String()).Do(ctx)
}

func (c *esClient) DeleteDataStream(ctx context.Context, name estype.DataStream) (*indices_delete_data_stream.Response, error) {
	return c.typedClient.Indices.DeleteDataStream(name.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// ILM (Index Lifecycle Management) API
// ---------------------------------------------------------------------------

func (c *esClient) PutLifecycle(ctx context.Context, policy estype.Policy, opts ...PutLifecycleOption) (*ilm_put_lifecycle.Response, error) {
	r := c.typedClient.Ilm.PutLifecycle(policy.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) GetLifecycle(ctx context.Context, policy estype.Policy) (ilm_get_lifecycle.Response, error) {
	return c.typedClient.Ilm.GetLifecycle().Policy(policy.String()).Do(ctx)
}

func (c *esClient) ExplainLifecycle(ctx context.Context, indexName estype.Index) (*ilm_explain_lifecycle.Response, error) {
	return c.typedClient.Ilm.ExplainLifecycle(indexName.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// Ingest Pipeline API
// ---------------------------------------------------------------------------

func (c *esClient) PutPipeline(ctx context.Context, id estype.Pipeline, opts ...PutPipelineOption) (*ingest_put_pipeline.Response, error) {
	r := c.typedClient.Ingest.PutPipeline(id.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) GetPipeline(ctx context.Context, id estype.Pipeline) (ingest_get_pipeline.Response, error) {
	return c.typedClient.Ingest.GetPipeline().Id(id.String()).Do(ctx)
}

func (c *esClient) DeletePipeline(ctx context.Context, id estype.Pipeline) (*ingest_delete_pipeline.Response, error) {
	return c.typedClient.Ingest.DeletePipeline(id.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// ES|QL API
// ---------------------------------------------------------------------------

func (c *esClient) EsqlQuery(ctx context.Context, query estype.ESQLQuery, opts ...EsqlQueryOption) (esql_query.Response, error) {
	req := esql_query.NewRequest()
	req.Query = query.String()
	r := c.typedClient.Esql.Query().Request(req)
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

// ---------------------------------------------------------------------------
// Security & API Key API
// ---------------------------------------------------------------------------

func (c *esClient) CreateApiKey(ctx context.Context, opts ...CreateApiKeyOption) (*security_create_api_key.Response, error) {
	r := c.typedClient.Security.CreateApiKey()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) GetApiKey(ctx context.Context, opts ...GetApiKeyOption) (*security_get_api_key.Response, error) {
	r := c.typedClient.Security.GetApiKey()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) InvalidateApiKey(ctx context.Context, opts ...InvalidateApiKeyOption) (*security_invalidate_api_key.Response, error) {
	r := c.typedClient.Security.InvalidateApiKey()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

// ---------------------------------------------------------------------------
// Snapshot & Restore API
// ---------------------------------------------------------------------------

func (c *esClient) CreateRepository(ctx context.Context, repo estype.Repository, opts ...CreateRepositoryOption) (*snapshot_create_repository.Response, error) {
	r := c.typedClient.Snapshot.CreateRepository(repo.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CreateSnapshot(ctx context.Context, repo estype.Repository, snap estype.Snapshot, opts ...CreateSnapshotOption) (*snapshot_create.Response, error) {
	r := c.typedClient.Snapshot.Create(repo.String(), snap.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) RestoreSnapshot(ctx context.Context, repo estype.Repository, snap estype.Snapshot, opts ...RestoreSnapshotOption) (*snapshot_restore.Response, error) {
	r := c.typedClient.Snapshot.Restore(repo.String(), snap.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

// ---------------------------------------------------------------------------
// Inference API
// ---------------------------------------------------------------------------

func (c *esClient) PutInference(ctx context.Context, inferenceId estype.InferenceID, opts ...PutInferenceOption) (*inference_put.Response, error) {
	r := c.typedClient.Inference.Put(inferenceId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) GetInference(ctx context.Context, inferenceId estype.InferenceID) (*inference_get.Response, error) {
	return c.typedClient.Inference.Get().InferenceId(inferenceId.String()).Do(ctx)
}

func (c *esClient) DeleteInference(ctx context.Context, inferenceId estype.InferenceID) (*inference_delete.Response, error) {
	return c.typedClient.Inference.Delete(inferenceId.String()).Do(ctx)
}

func (c *esClient) Inference(ctx context.Context, inferenceId estype.InferenceID, opts ...InferenceOption) (*inference_inference.Response, error) {
	r := c.typedClient.Inference.Inference(inferenceId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

// ---------------------------------------------------------------------------
// Machine Learning (ML) API
// ---------------------------------------------------------------------------

func (c *esClient) MlPutJob(ctx context.Context, jobId estype.MLJobID, opts ...MlPutJobOption) (*ml_put_job.Response, error) {
	r := c.typedClient.Ml.PutJob(jobId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) MlGetJobs(ctx context.Context, jobId estype.MLJobID) (*ml_get_jobs.Response, error) {
	return c.typedClient.Ml.GetJobs().JobId(jobId.String()).Do(ctx)
}

func (c *esClient) MlDeleteJob(ctx context.Context, jobId estype.MLJobID) (*ml_delete_job.Response, error) {
	return c.typedClient.Ml.DeleteJob(jobId.String()).Do(ctx)
}

func (c *esClient) MlOpenJob(ctx context.Context, jobId estype.MLJobID, opts ...MlOpenJobOption) (*ml_open_job.Response, error) {
	r := c.typedClient.Ml.OpenJob(jobId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) MlCloseJob(ctx context.Context, jobId estype.MLJobID, opts ...MlCloseJobOption) (*ml_close_job.Response, error) {
	r := c.typedClient.Ml.CloseJob(jobId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) MlPutDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlPutDatafeedOption) (*ml_put_datafeed.Response, error) {
	r := c.typedClient.Ml.PutDatafeed(datafeedId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) MlGetDatafeeds(ctx context.Context, datafeedId estype.DatafeedID) (*ml_get_datafeeds.Response, error) {
	return c.typedClient.Ml.GetDatafeeds().DatafeedId(datafeedId.String()).Do(ctx)
}

func (c *esClient) MlDeleteDatafeed(ctx context.Context, datafeedId estype.DatafeedID) (*ml_delete_datafeed.Response, error) {
	return c.typedClient.Ml.DeleteDatafeed(datafeedId.String()).Do(ctx)
}

func (c *esClient) MlStartDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlStartDatafeedOption) (*ml_start_datafeed.Response, error) {
	r := c.typedClient.Ml.StartDatafeed(datafeedId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) MlStopDatafeed(ctx context.Context, datafeedId estype.DatafeedID, opts ...MlStopDatafeedOption) (*ml_stop_datafeed.Response, error) {
	r := c.typedClient.Ml.StopDatafeed(datafeedId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

// ---------------------------------------------------------------------------
// CCR (Cross-Cluster Replication) API
// ---------------------------------------------------------------------------

func (c *esClient) CcrFollow(ctx context.Context, followerIndex estype.Index, opts ...CcrFollowOption) (*ccr_follow.Response, error) {
	r := c.typedClient.Ccr.Follow(followerIndex.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CcrPauseFollow(ctx context.Context, followerIndex estype.Index) (*ccr_pause_follow.Response, error) {
	return c.typedClient.Ccr.PauseFollow(followerIndex.String()).Do(ctx)
}

func (c *esClient) CcrResumeFollow(ctx context.Context, followerIndex estype.Index, opts ...CcrResumeFollowOption) (*ccr_resume_follow.Response, error) {
	r := c.typedClient.Ccr.ResumeFollow(followerIndex.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CcrUnfollow(ctx context.Context, followerIndex estype.Index) (*ccr_unfollow.Response, error) {
	return c.typedClient.Ccr.Unfollow(followerIndex.String()).Do(ctx)
}

func (c *esClient) CcrFollowStats(ctx context.Context, followerIndex estype.Index) (*ccr_follow_stats.Response, error) {
	return c.typedClient.Ccr.FollowStats(followerIndex.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// Transform API
// ---------------------------------------------------------------------------

func (c *esClient) PutTransform(ctx context.Context, transformId estype.TransformID, opts ...PutTransformOption) (*transform_put_transform.Response, error) {
	r := c.typedClient.Transform.PutTransform(transformId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) GetTransform(ctx context.Context, transformId estype.TransformID) (*transform_get_transform.Response, error) {
	return c.typedClient.Transform.GetTransform().TransformId(transformId.String()).Do(ctx)
}

func (c *esClient) DeleteTransform(ctx context.Context, transformId estype.TransformID) (*transform_delete_transform.Response, error) {
	return c.typedClient.Transform.DeleteTransform(transformId.String()).Do(ctx)
}

func (c *esClient) StartTransform(ctx context.Context, transformId estype.TransformID) (*transform_start_transform.Response, error) {
	return c.typedClient.Transform.StartTransform(transformId.String()).Do(ctx)
}

func (c *esClient) StopTransform(ctx context.Context, transformId estype.TransformID) (*transform_stop_transform.Response, error) {
	return c.typedClient.Transform.StopTransform(transformId.String()).Do(ctx)
}

func (c *esClient) GetTransformStats(ctx context.Context, transformId estype.TransformID) (*transform_get_transform_stats.Response, error) {
	return c.typedClient.Transform.GetTransformStats(transformId.String()).Do(ctx)
}

// ---------------------------------------------------------------------------
// Cat API
// ---------------------------------------------------------------------------

func (c *esClient) CatAliases(ctx context.Context, aliasName estype.Alias) (cat_aliases.Response, error) {
	return c.typedClient.Cat.Aliases().Name(aliasName.String()).Do(ctx)
}

func (c *esClient) CatIndices(ctx context.Context, indexName estype.Index) (cat_indices.Response, error) {
	return c.typedClient.Cat.Indices().Index(indexName.String()).Do(ctx)
}

func (c *esClient) CatNodes(ctx context.Context) (cat_nodes.Response, error) {
	return c.typedClient.Cat.Nodes().Do(ctx)
}

func (c *esClient) CatAllocation(ctx context.Context, opts ...CatAllocationOption) (cat_allocation.Response, error) {
	r := c.typedClient.Cat.Allocation()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

// CatCircuitBreaker returns information about circuit breaker status as raw JSON.
// This endpoint is not available in the Elasticsearch v8 TypedAPI.
func (c *esClient) CatCircuitBreaker(ctx context.Context, opts ...CatCircuitBreakerOption) (json.RawMessage, error) {
	return c.performRaw(ctx, "GET", "/_cat/circuit_breaker", nil)
}

func (c *esClient) CatComponentTemplates(ctx context.Context, name estype.Template, opts ...CatComponentTemplatesOption) (cat_component_templates.Response, error) {
	r := c.typedClient.Cat.ComponentTemplates().Name(name.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatCount(ctx context.Context, indexName estype.Index, opts ...CatCountOption) (cat_count.Response, error) {
	r := c.typedClient.Cat.Count().Index(indexName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatFielddata(ctx context.Context, fields []estype.Field, opts ...CatFielddataOption) (cat_fielddata.Response, error) {
	fieldNames := strings.Join(estype.FieldNames(fields...), ",")
	r := c.typedClient.Cat.Fielddata().Fields(fieldNames)
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatHealth(ctx context.Context, opts ...CatHealthOption) (cat_health.Response, error) {
	r := c.typedClient.Cat.Health()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatHelp(ctx context.Context, opts ...CatHelpOption) (*cat_help.Response, error) {
	r := c.typedClient.Cat.Help()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatMaster(ctx context.Context, opts ...CatMasterOption) (cat_master.Response, error) {
	r := c.typedClient.Cat.Master()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatMlDataFrameAnalytics(ctx context.Context, analyticsId estype.DataFrameAnalyticsID, opts ...CatMlDataFrameAnalyticsOption) (cat_ml_data_frame_analytics.Response, error) {
	r := c.typedClient.Cat.MlDataFrameAnalytics().Id(analyticsId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatMlDatafeeds(ctx context.Context, datafeedId estype.DatafeedID, opts ...CatMlDatafeedsOption) (cat_ml_datafeeds.Response, error) {
	r := c.typedClient.Cat.MlDatafeeds().DatafeedId(datafeedId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatMlJobs(ctx context.Context, jobId estype.MLJobID, opts ...CatMlJobsOption) (cat_ml_jobs.Response, error) {
	r := c.typedClient.Cat.MlJobs().JobId(jobId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatMlTrainedModels(ctx context.Context, modelId estype.TrainedModelID, opts ...CatMlTrainedModelsOption) (cat_ml_trained_models.Response, error) {
	r := c.typedClient.Cat.MlTrainedModels().ModelId(modelId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatNodeattrs(ctx context.Context, opts ...CatNodeattrsOption) (cat_nodeattrs.Response, error) {
	r := c.typedClient.Cat.Nodeattrs()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatPendingTasks(ctx context.Context, opts ...CatPendingTasksOption) (cat_pending_tasks.Response, error) {
	r := c.typedClient.Cat.PendingTasks()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatPlugins(ctx context.Context, opts ...CatPluginsOption) (cat_plugins.Response, error) {
	r := c.typedClient.Cat.Plugins()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatRecovery(ctx context.Context, indexName estype.Index, opts ...CatRecoveryOption) (cat_recovery.Response, error) {
	r := c.typedClient.Cat.Recovery().Index(indexName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatRepositories(ctx context.Context, opts ...CatRepositoriesOption) (cat_repositories.Response, error) {
	r := c.typedClient.Cat.Repositories()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatSegments(ctx context.Context, indexName estype.Index, opts ...CatSegmentsOption) (cat_segments.Response, error) {
	r := c.typedClient.Cat.Segments().Index(indexName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatShards(ctx context.Context, indexName estype.Index, opts ...CatShardsOption) (cat_shards.Response, error) {
	r := c.typedClient.Cat.Shards().Index(indexName.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatSnapshots(ctx context.Context, repository estype.Repository, opts ...CatSnapshotsOption) (cat_snapshots.Response, error) {
	r := c.typedClient.Cat.Snapshots().Repository(repository.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatTasks(ctx context.Context, opts ...CatTasksOption) (cat_tasks.Response, error) {
	r := c.typedClient.Cat.Tasks()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatTemplates(ctx context.Context, name estype.Template, opts ...CatTemplatesOption) (cat_templates.Response, error) {
	r := c.typedClient.Cat.Templates().Name(name.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatThreadPool(ctx context.Context, opts ...CatThreadPoolOption) (cat_thread_pool.Response, error) {
	r := c.typedClient.Cat.ThreadPool()
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}

func (c *esClient) CatTransforms(ctx context.Context, transformId estype.TransformID, opts ...CatTransformsOption) (cat_transforms.Response, error) {
	r := c.typedClient.Cat.Transforms().TransformId(transformId.String())
	for _, opt := range opts {
		opt(r)
	}
	return r.Do(ctx)
}
