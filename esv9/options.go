package esv9

import (
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
	core_count "github.com/elastic/go-elasticsearch/v9/typedapi/core/count"
	core_delete_by_query "github.com/elastic/go-elasticsearch/v9/typedapi/core/deletebyquery"
	coredelete "github.com/elastic/go-elasticsearch/v9/typedapi/core/delete"
	core_exists "github.com/elastic/go-elasticsearch/v9/typedapi/core/exists"
	coreget "github.com/elastic/go-elasticsearch/v9/typedapi/core/get"
	coreidx "github.com/elastic/go-elasticsearch/v9/typedapi/core/index"
	core_mget "github.com/elastic/go-elasticsearch/v9/typedapi/core/mget"
	core_msearch "github.com/elastic/go-elasticsearch/v9/typedapi/core/msearch"
	core_open_point_in_time "github.com/elastic/go-elasticsearch/v9/typedapi/core/openpointintime"
	core_reindex "github.com/elastic/go-elasticsearch/v9/typedapi/core/reindex"
	core_scroll "github.com/elastic/go-elasticsearch/v9/typedapi/core/scroll"
	core_update_by_query "github.com/elastic/go-elasticsearch/v9/typedapi/core/updatebyquery"
	esql_query "github.com/elastic/go-elasticsearch/v9/typedapi/esql/query"
	ilm_explain_lifecycle "github.com/elastic/go-elasticsearch/v9/typedapi/ilm/explainlifecycle"
	ilm_get_lifecycle "github.com/elastic/go-elasticsearch/v9/typedapi/ilm/getlifecycle"
	ilm_put_lifecycle "github.com/elastic/go-elasticsearch/v9/typedapi/ilm/putlifecycle"
	indices_analyze "github.com/elastic/go-elasticsearch/v9/typedapi/indices/analyze"
	indices_clear_cache "github.com/elastic/go-elasticsearch/v9/typedapi/indices/clearcache"
	indices_close "github.com/elastic/go-elasticsearch/v9/typedapi/indices/close"
	indices_delete_alias "github.com/elastic/go-elasticsearch/v9/typedapi/indices/deletealias"
	indices_delete_data_stream "github.com/elastic/go-elasticsearch/v9/typedapi/indices/deletedatastream"
	indices_delete_index_template "github.com/elastic/go-elasticsearch/v9/typedapi/indices/deleteindextemplate"
	indices_exists_index_template "github.com/elastic/go-elasticsearch/v9/typedapi/indices/existsindextemplate"
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
	indices_put_settings "github.com/elastic/go-elasticsearch/v9/typedapi/indices/putsettings"
	indices_rollover "github.com/elastic/go-elasticsearch/v9/typedapi/indices/rollover"
	indices_stats "github.com/elastic/go-elasticsearch/v9/typedapi/indices/stats"
	idxdelete "github.com/elastic/go-elasticsearch/v9/typedapi/indices/delete"
	idxrefresh "github.com/elastic/go-elasticsearch/v9/typedapi/indices/refresh"
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
)

// DocumentExistsOption customises a DocumentExists request.
type DocumentExistsOption func(*core_exists.Exists)

// BulkOption customises a Bulk request.
type BulkOption func(*core_bulk.Bulk)

// MgetOption customises a Mget request.
type MgetOption func(*core_mget.Mget)

// MsearchOption customises a Msearch request.
type MsearchOption func(*core_msearch.Msearch)

// CountOption customises a Count request.
type CountOption func(*core_count.Count)

// ScrollOption customises a Scroll request.
type ScrollOption func(*core_scroll.Scroll)

// ClearScrollOption customises a ClearScroll request.
type ClearScrollOption func(*core_clear_scroll.ClearScroll)

// UpdateByQueryOption customises an UpdateByQuery request.
type UpdateByQueryOption func(*core_update_by_query.UpdateByQuery)

// DeleteByQueryOption customises a DeleteByQuery request.
type DeleteByQueryOption func(*core_delete_by_query.DeleteByQuery)

// PutMappingOption customises a PutMapping request.
type PutMappingOption func(*indices_put_mapping.PutMapping)

// PutSettingsOption customises a PutSettings request.
type PutSettingsOption func(*indices_put_settings.PutSettings)

// RolloverOption customises a Rollover request.
type RolloverOption func(*indices_rollover.Rollover)

// PutIndexTemplateOption customises a PutIndexTemplate request.
type PutIndexTemplateOption func(*indices_put_index_template.PutIndexTemplate)

// TasksListOption customises a TasksList request.
type TasksListOption func(*tasks_list.List)

// AnalyzeOption customises an Analyze request.
type AnalyzeOption func(*indices_analyze.Analyze)

// OpenPointInTimeOption customises an OpenPointInTime request.
type OpenPointInTimeOption func(*core_open_point_in_time.OpenPointInTime)

// ClosePointInTimeOption customises a ClosePointInTime request.
type ClosePointInTimeOption func(*core_close_point_in_time.ClosePointInTime)

// PutLifecycleOption customises a PutLifecycle request.
type PutLifecycleOption func(*ilm_put_lifecycle.PutLifecycle)

// PutPipelineOption customises a PutPipeline request.
type PutPipelineOption func(*ingest_put_pipeline.PutPipeline)

// EsqlQueryOption customises an EsqlQuery request.
type EsqlQueryOption func(*esql_query.Query)

// CreateApiKeyOption customises a CreateApiKey request.
type CreateApiKeyOption func(*security_create_api_key.CreateApiKey)

// GetApiKeyOption customises a GetApiKey request.
type GetApiKeyOption func(*security_get_api_key.GetApiKey)

// InvalidateApiKeyOption customises an InvalidateApiKey request.
type InvalidateApiKeyOption func(*security_invalidate_api_key.InvalidateApiKey)

// CreateRepositoryOption customises a CreateRepository request.
type CreateRepositoryOption func(*snapshot_create_repository.CreateRepository)

// CreateSnapshotOption customises a CreateSnapshot request.
type CreateSnapshotOption func(*snapshot_create.Create)

// RestoreSnapshotOption customises a RestoreSnapshot request.
type RestoreSnapshotOption func(*snapshot_restore.Restore)

// PutInferenceOption customises a PutInference request.
type PutInferenceOption func(*inference_put.Put)

// InferenceOption customises an Inference request.
type InferenceOption func(*inference_inference.Inference)

// MlPutJobOption customises an MlPutJob request.
type MlPutJobOption func(*ml_put_job.PutJob)

// MlOpenJobOption customises an MlOpenJob request.
type MlOpenJobOption func(*ml_open_job.OpenJob)

// MlCloseJobOption customises an MlCloseJob request.
type MlCloseJobOption func(*ml_close_job.CloseJob)

// MlPutDatafeedOption customises an MlPutDatafeed request.
type MlPutDatafeedOption func(*ml_put_datafeed.PutDatafeed)

// MlStartDatafeedOption customises an MlStartDatafeed request.
type MlStartDatafeedOption func(*ml_start_datafeed.StartDatafeed)

// MlStopDatafeedOption customises an MlStopDatafeed request.
type MlStopDatafeedOption func(*ml_stop_datafeed.StopDatafeed)

// CcrFollowOption customises a CcrFollow request.
type CcrFollowOption func(*ccr_follow.Follow)

// CcrResumeFollowOption customises a CcrResumeFollow request.
type CcrResumeFollowOption func(*ccr_resume_follow.ResumeFollow)

// PutTransformOption customises a PutTransform request.
type PutTransformOption func(*transform_put_transform.PutTransform)

// CatAllocationOption customises a CatAllocation request.
type CatAllocationOption func(*cat_allocation.Allocation)

// CatCircuitBreakerOption customises a CatCircuitBreaker request.
type CatCircuitBreakerOption func(*cat_circuit_breaker.CircuitBreaker)

// CatComponentTemplatesOption customises a CatComponentTemplates request.
type CatComponentTemplatesOption func(*cat_component_templates.ComponentTemplates)

// CatCountOption customises a CatCount request.
type CatCountOption func(*cat_count.Count)

// CatFielddataOption customises a CatFielddata request.
type CatFielddataOption func(*cat_fielddata.Fielddata)

// CatHealthOption customises a CatHealth request.
type CatHealthOption func(*cat_health.Health)

// CatHelpOption customises a CatHelp request.
type CatHelpOption func(*cat_help.Help)

// CatMasterOption customises a CatMaster request.
type CatMasterOption func(*cat_master.Master)

// CatMlDataFrameAnalyticsOption customises a CatMlDataFrameAnalytics request.
type CatMlDataFrameAnalyticsOption func(*cat_ml_data_frame_analytics.MlDataFrameAnalytics)

// CatMlDatafeedsOption customises a CatMlDatafeeds request.
type CatMlDatafeedsOption func(*cat_ml_datafeeds.MlDatafeeds)

// CatMlJobsOption customises a CatMlJobs request.
type CatMlJobsOption func(*cat_ml_jobs.MlJobs)

// CatMlTrainedModelsOption customises a CatMlTrainedModels request.
type CatMlTrainedModelsOption func(*cat_ml_trained_models.MlTrainedModels)

// CatNodeattrsOption customises a CatNodeattrs request.
type CatNodeattrsOption func(*cat_nodeattrs.Nodeattrs)

// CatPendingTasksOption customises a CatPendingTasks request.
type CatPendingTasksOption func(*cat_pending_tasks.PendingTasks)

// CatPluginsOption customises a CatPlugins request.
type CatPluginsOption func(*cat_plugins.Plugins)

// CatRecoveryOption customises a CatRecovery request.
type CatRecoveryOption func(*cat_recovery.Recovery)

// CatRepositoriesOption customises a CatRepositories request.
type CatRepositoriesOption func(*cat_repositories.Repositories)

// CatSegmentsOption customises a CatSegments request.
type CatSegmentsOption func(*cat_segments.Segments)

// CatShardsOption customises a CatShards request.
type CatShardsOption func(*cat_shards.Shards)

// CatSnapshotsOption customises a CatSnapshots request.
type CatSnapshotsOption func(*cat_snapshots.Snapshots)

// CatTasksOption customises a CatTasks request.
type CatTasksOption func(*cat_tasks.Tasks)

// CatTemplatesOption customises a CatTemplates request.
type CatTemplatesOption func(*cat_templates.Templates)

// CatThreadPoolOption customises a CatThreadPool request.
type CatThreadPoolOption func(*cat_thread_pool.ThreadPool)

// CatTransformsOption customises a CatTransforms request.
type CatTransformsOption func(*cat_transforms.Transforms)

// CreateDocumentOption customises a CreateDocument request.
type CreateDocumentOption func(*coreidx.Index)

// GetDocumentOption customises a GetDocument request.
type GetDocumentOption func(*coreget.Get)

// DeleteDocumentOption customises a DeleteDocument request.
type DeleteDocumentOption func(*coredelete.Delete)

// IndexRefreshOption customises an IndexRefresh request.
type IndexRefreshOption func(*idxrefresh.Refresh)

// DeleteIndexOption customises a DeleteIndex request.
type DeleteIndexOption func(*idxdelete.Delete)

// ReindexOption customises a Reindex request.
type ReindexOption func(*core_reindex.Reindex)

// GetMappingOption customises a GetMapping request.
type GetMappingOption func(*indices_get_mapping.GetMapping)

// GetSettingsOption customises a GetSettings request.
type GetSettingsOption func(*indices_get_settings.GetSettings)

// OpenIndexOption customises an OpenIndex request.
type OpenIndexOption func(*indices_open.Open)

// CloseIndexOption customises a CloseIndex request.
type CloseIndexOption func(*indices_close.Close)

// FlushOption customises a Flush request.
type FlushOption func(*indices_flush.Flush)

// ClearCacheOption customises a ClearCache request.
type ClearCacheOption func(*indices_clear_cache.ClearCache)

// ForceMergeOption customises a ForceMerge request.
type ForceMergeOption func(*indices_forcemerge.Forcemerge)

// IndicesStatsOption customises an IndicesStats request.
type IndicesStatsOption func(*indices_stats.Stats)

// ClusterHealthOption customises a ClusterHealth request.
type ClusterHealthOption func(*cluster_health.Health)

// GetIndexTemplateOption customises a GetIndexTemplate request.
type GetIndexTemplateOption func(*indices_get_index_template.GetIndexTemplate)

// DeleteIndexTemplateOption customises a DeleteIndexTemplate request.
type DeleteIndexTemplateOption func(*indices_delete_index_template.DeleteIndexTemplate)

// ExistsIndexTemplateOption customises an ExistsIndexTemplate request.
type ExistsIndexTemplateOption func(*indices_exists_index_template.ExistsIndexTemplate)

// GetAliasOption customises a GetAlias request.
type GetAliasOption func(*indices_get_alias.GetAlias)

// DeleteAliasOption customises a DeleteAlias request.
type DeleteAliasOption func(*indices_delete_alias.DeleteAlias)

// TasksCancelOption customises a TasksCancel request.
type TasksCancelOption func(*tasks_cancel.Cancel)

// GetDataStreamOption customises a GetDataStream request.
type GetDataStreamOption func(*indices_get_data_stream.GetDataStream)

// DeleteDataStreamOption customises a DeleteDataStream request.
type DeleteDataStreamOption func(*indices_delete_data_stream.DeleteDataStream)

// GetLifecycleOption customises a GetLifecycle request.
type GetLifecycleOption func(*ilm_get_lifecycle.GetLifecycle)

// ExplainLifecycleOption customises an ExplainLifecycle request.
type ExplainLifecycleOption func(*ilm_explain_lifecycle.ExplainLifecycle)

// GetPipelineOption customises a GetPipeline request.
type GetPipelineOption func(*ingest_get_pipeline.GetPipeline)

// DeletePipelineOption customises a DeletePipeline request.
type DeletePipelineOption func(*ingest_delete_pipeline.DeletePipeline)

// GetInferenceOption customises a GetInference request.
type GetInferenceOption func(*inference_get.Get)

// DeleteInferenceOption customises a DeleteInference request.
type DeleteInferenceOption func(*inference_delete.Delete)

// MlGetJobsOption customises an MlGetJobs request.
type MlGetJobsOption func(*ml_get_jobs.GetJobs)

// MlDeleteJobOption customises an MlDeleteJob request.
type MlDeleteJobOption func(*ml_delete_job.DeleteJob)

// MlGetDatafeedsOption customises an MlGetDatafeeds request.
type MlGetDatafeedsOption func(*ml_get_datafeeds.GetDatafeeds)

// MlDeleteDatafeedOption customises an MlDeleteDatafeed request.
type MlDeleteDatafeedOption func(*ml_delete_datafeed.DeleteDatafeed)

// CcrPauseFollowOption customises a CcrPauseFollow request.
type CcrPauseFollowOption func(*ccr_pause_follow.PauseFollow)

// CcrUnfollowOption customises a CcrUnfollow request.
type CcrUnfollowOption func(*ccr_unfollow.Unfollow)

// CcrFollowStatsOption customises a CcrFollowStats request.
type CcrFollowStatsOption func(*ccr_follow_stats.FollowStats)

// GetTransformOption customises a GetTransform request.
type GetTransformOption func(*transform_get_transform.GetTransform)

// DeleteTransformOption customises a DeleteTransform request.
type DeleteTransformOption func(*transform_delete_transform.DeleteTransform)

// StartTransformOption customises a StartTransform request.
type StartTransformOption func(*transform_start_transform.StartTransform)

// StopTransformOption customises a StopTransform request.
type StopTransformOption func(*transform_stop_transform.StopTransform)

// GetTransformStatsOption customises a GetTransformStats request.
type GetTransformStatsOption func(*transform_get_transform_stats.GetTransformStats)

// CatAliasesOption customises a CatAliases request.
type CatAliasesOption func(*cat_aliases.Aliases)

// CatIndicesOption customises a CatIndices request.
type CatIndicesOption func(*cat_indices.Indices)

// CatNodesOption customises a CatNodes request.
type CatNodesOption func(*cat_nodes.Nodes)
