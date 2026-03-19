package esv9

import (
	cat_circuit_breaker "github.com/elastic/go-elasticsearch/v9/typedapi/cat/circuitbreaker"
	cat_allocation "github.com/elastic/go-elasticsearch/v9/typedapi/cat/allocation"
	cat_component_templates "github.com/elastic/go-elasticsearch/v9/typedapi/cat/componenttemplates"
	cat_count "github.com/elastic/go-elasticsearch/v9/typedapi/cat/count"
	cat_fielddata "github.com/elastic/go-elasticsearch/v9/typedapi/cat/fielddata"
	cat_health "github.com/elastic/go-elasticsearch/v9/typedapi/cat/health"
	cat_help "github.com/elastic/go-elasticsearch/v9/typedapi/cat/help"
	cat_master "github.com/elastic/go-elasticsearch/v9/typedapi/cat/master"
	cat_ml_data_frame_analytics "github.com/elastic/go-elasticsearch/v9/typedapi/cat/mldataframeanalytics"
	cat_ml_datafeeds "github.com/elastic/go-elasticsearch/v9/typedapi/cat/mldatafeeds"
	cat_ml_jobs "github.com/elastic/go-elasticsearch/v9/typedapi/cat/mljobs"
	cat_ml_trained_models "github.com/elastic/go-elasticsearch/v9/typedapi/cat/mltrainedmodels"
	cat_nodeattrs "github.com/elastic/go-elasticsearch/v9/typedapi/cat/nodeattrs"
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
	ccr_resume_follow "github.com/elastic/go-elasticsearch/v9/typedapi/ccr/resumefollow"
	core_bulk "github.com/elastic/go-elasticsearch/v9/typedapi/core/bulk"
	core_clear_scroll "github.com/elastic/go-elasticsearch/v9/typedapi/core/clearscroll"
	core_close_point_in_time "github.com/elastic/go-elasticsearch/v9/typedapi/core/closepointintime"
	core_count "github.com/elastic/go-elasticsearch/v9/typedapi/core/count"
	core_delete_by_query "github.com/elastic/go-elasticsearch/v9/typedapi/core/deletebyquery"
	core_exists "github.com/elastic/go-elasticsearch/v9/typedapi/core/exists"
	core_mget "github.com/elastic/go-elasticsearch/v9/typedapi/core/mget"
	core_msearch "github.com/elastic/go-elasticsearch/v9/typedapi/core/msearch"
	core_open_point_in_time "github.com/elastic/go-elasticsearch/v9/typedapi/core/openpointintime"
	core_scroll "github.com/elastic/go-elasticsearch/v9/typedapi/core/scroll"
	core_update_by_query "github.com/elastic/go-elasticsearch/v9/typedapi/core/updatebyquery"
	esql_query "github.com/elastic/go-elasticsearch/v9/typedapi/esql/query"
	ilm_put_lifecycle "github.com/elastic/go-elasticsearch/v9/typedapi/ilm/putlifecycle"
	indices_analyze "github.com/elastic/go-elasticsearch/v9/typedapi/indices/analyze"
	indices_put_index_template "github.com/elastic/go-elasticsearch/v9/typedapi/indices/putindextemplate"
	indices_put_mapping "github.com/elastic/go-elasticsearch/v9/typedapi/indices/putmapping"
	indices_put_settings "github.com/elastic/go-elasticsearch/v9/typedapi/indices/putsettings"
	indices_rollover "github.com/elastic/go-elasticsearch/v9/typedapi/indices/rollover"
	inference_inference "github.com/elastic/go-elasticsearch/v9/typedapi/inference/inference"
	inference_put "github.com/elastic/go-elasticsearch/v9/typedapi/inference/put"
	ingest_put_pipeline "github.com/elastic/go-elasticsearch/v9/typedapi/ingest/putpipeline"
	ml_close_job "github.com/elastic/go-elasticsearch/v9/typedapi/ml/closejob"
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
	tasks_list "github.com/elastic/go-elasticsearch/v9/typedapi/tasks/list"
	transform_put_transform "github.com/elastic/go-elasticsearch/v9/typedapi/transform/puttransform"
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
