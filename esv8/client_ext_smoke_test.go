package esv8

import (
	"context"
	"io"
	"net/http"
	"testing"

	async_search_submit "github.com/elastic/go-elasticsearch/v8/typedapi/asyncsearch/submit"
	autoscaling_put_autoscaling_policy "github.com/elastic/go-elasticsearch/v8/typedapi/autoscaling/putautoscalingpolicy"
	ccr_put_auto_follow_pattern "github.com/elastic/go-elasticsearch/v8/typedapi/ccr/putautofollowpattern"
	cluster_put_settings "github.com/elastic/go-elasticsearch/v8/typedapi/cluster/putsettings"
	connector_last_sync "github.com/elastic/go-elasticsearch/v8/typedapi/connector/lastsync"
	connector_post "github.com/elastic/go-elasticsearch/v8/typedapi/connector/post"
	connector_put "github.com/elastic/go-elasticsearch/v8/typedapi/connector/put"
	connector_sync_job_claim "github.com/elastic/go-elasticsearch/v8/typedapi/connector/syncjobclaim"
	connector_sync_job_error "github.com/elastic/go-elasticsearch/v8/typedapi/connector/syncjoberror"
	connector_sync_job_post "github.com/elastic/go-elasticsearch/v8/typedapi/connector/syncjobpost"
	connector_sync_job_update_stats "github.com/elastic/go-elasticsearch/v8/typedapi/connector/syncjobupdatestats"
	connector_update_api_key_id "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updateapikeyid"
	connector_update_configuration "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updateconfiguration"
	connector_update_error "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updateerror"
	connector_update_features "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updatefeatures"
	connector_update_filtering "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updatefiltering"
	connector_update_filtering_validation "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updatefilteringvalidation"
	connector_update_index_name "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updateindexname"
	connector_update_name "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updatename"
	connector_update_native "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updatenative"
	connector_update_pipeline "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updatepipeline"
	connector_update_scheduling "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updatescheduling"
	connector_update_service_type "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updateservicetype"
	connector_update_status "github.com/elastic/go-elasticsearch/v8/typedapi/connector/updatestatus"
	core_clear_scroll "github.com/elastic/go-elasticsearch/v8/typedapi/core/clearscroll"
	core_close_point_in_time "github.com/elastic/go-elasticsearch/v8/typedapi/core/closepointintime"
	core_count "github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	core_create "github.com/elastic/go-elasticsearch/v8/typedapi/core/create"
	core_delete_by_query "github.com/elastic/go-elasticsearch/v8/typedapi/core/deletebyquery"
	core_field_caps "github.com/elastic/go-elasticsearch/v8/typedapi/core/fieldcaps"
	core_mget "github.com/elastic/go-elasticsearch/v8/typedapi/core/mget"
	core_msearch "github.com/elastic/go-elasticsearch/v8/typedapi/core/msearch"
	core_scroll "github.com/elastic/go-elasticsearch/v8/typedapi/core/scroll"
	core_update_by_query "github.com/elastic/go-elasticsearch/v8/typedapi/core/updatebyquery"
	enrich_put_policy "github.com/elastic/go-elasticsearch/v8/typedapi/enrich/putpolicy"
	eql_search "github.com/elastic/go-elasticsearch/v8/typedapi/eql/search"
	esql_async_query "github.com/elastic/go-elasticsearch/v8/typedapi/esql/asyncquery"
	ilm_put_lifecycle "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/putlifecycle"
	indices_analyze "github.com/elastic/go-elasticsearch/v8/typedapi/indices/analyze"
	indices_put_index_template "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putindextemplate"
	indices_put_mapping "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putmapping"
	indices_put_settings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putsettings"
	indices_put_template "github.com/elastic/go-elasticsearch/v8/typedapi/indices/puttemplate"
	indices_validate_query "github.com/elastic/go-elasticsearch/v8/typedapi/indices/validatequery"
	inference_inference "github.com/elastic/go-elasticsearch/v8/typedapi/inference/inference"
	ingest_put_pipeline "github.com/elastic/go-elasticsearch/v8/typedapi/ingest/putpipeline"
	ml_put_datafeed "github.com/elastic/go-elasticsearch/v8/typedapi/ml/putdatafeed"
	ml_put_job "github.com/elastic/go-elasticsearch/v8/typedapi/ml/putjob"
	query_rules_put_rule "github.com/elastic/go-elasticsearch/v8/typedapi/queryrules/putrule"
	query_rules_put_ruleset "github.com/elastic/go-elasticsearch/v8/typedapi/queryrules/putruleset"
	query_rules_test "github.com/elastic/go-elasticsearch/v8/typedapi/queryrules/test"
	search_application_post_behavioral_analytics_event "github.com/elastic/go-elasticsearch/v8/typedapi/searchapplication/postbehavioralanalyticsevent"
	search_application_put "github.com/elastic/go-elasticsearch/v8/typedapi/searchapplication/put"
	search_application_render_query "github.com/elastic/go-elasticsearch/v8/typedapi/searchapplication/renderquery"
	search_application_search "github.com/elastic/go-elasticsearch/v8/typedapi/searchapplication/search"
	security_activate_user_profile "github.com/elastic/go-elasticsearch/v8/typedapi/security/activateuserprofile"
	security_change_password "github.com/elastic/go-elasticsearch/v8/typedapi/security/changepassword"
	security_create_api_key "github.com/elastic/go-elasticsearch/v8/typedapi/security/createapikey"
	security_create_cross_cluster_api_key "github.com/elastic/go-elasticsearch/v8/typedapi/security/createcrossclusterapikey"
	security_delegate_pki "github.com/elastic/go-elasticsearch/v8/typedapi/security/delegatepki"
	security_get_token "github.com/elastic/go-elasticsearch/v8/typedapi/security/gettoken"
	security_grant_api_key "github.com/elastic/go-elasticsearch/v8/typedapi/security/grantapikey"
	security_invalidate_api_key "github.com/elastic/go-elasticsearch/v8/typedapi/security/invalidateapikey"
	security_put_role "github.com/elastic/go-elasticsearch/v8/typedapi/security/putrole"
	security_put_user "github.com/elastic/go-elasticsearch/v8/typedapi/security/putuser"
	security_query_api_keys "github.com/elastic/go-elasticsearch/v8/typedapi/security/queryapikeys"
	security_update_api_key "github.com/elastic/go-elasticsearch/v8/typedapi/security/updateapikey"
	snapshot_clone "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/clone"
	snapshot_create "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/create"
	snapshot_create_repository "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/createrepository"
	snapshot_restore "github.com/elastic/go-elasticsearch/v8/typedapi/snapshot/restore"
	sql_query "github.com/elastic/go-elasticsearch/v8/typedapi/sql/query"
	sql_translate "github.com/elastic/go-elasticsearch/v8/typedapi/sql/translate"
	transform_preview_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/previewtransform"
	transform_put_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/puttransform"
	transform_update_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/updatetransform"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

func TestClientExtSmoke_RequestPathsAndMethods(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		wantMethod string
		wantPath   string
		wantQuery  string
		call       func(*testing.T, *esClient)
	}

	tests := []testCase{
		{
			name:       "DocumentExists",
			wantMethod: http.MethodHead,
			wantPath:   "/products/_doc/doc-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.DocumentExists(
					context.Background(),
					estype.Alias("products"),
					estype.DocumentID("doc-1"),
				)
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "Bulk",
			wantMethod: http.MethodPost,
			wantPath:   "/logs/_bulk",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Bulk(context.Background(), estype.Alias("logs"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "Mget",
			wantMethod: http.MethodPost,
			wantPath:   "/products/_mget",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Mget(
					context.Background(),
					estype.Alias("products"),
					func(r *core_mget.Mget) {
						req := core_mget.NewRequest()
						req.Ids = []string{"1"}
						r.Request(req)
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "Msearch",
			wantMethod: http.MethodPost,
			wantPath:   "/products/_msearch",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Msearch(
					context.Background(),
					estype.Alias("products"),
					func(r *core_msearch.Msearch) {
						r.Request(&core_msearch.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "Count",
			wantMethod: http.MethodPost,
			wantPath:   "/products/_count",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Count(
					context.Background(),
					estype.Alias("products"),
					func(r *core_count.Count) {
						r.Request(&core_count.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "Scroll",
			wantMethod: http.MethodPost,
			wantPath:   "/_search/scroll",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Scroll(
					context.Background(),
					func(r *core_scroll.Scroll) {
						req := core_scroll.NewRequest()
						req.Scroll = "1m"
						req.ScrollId = "scroll-1"
						r.Request(req)
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClearScroll",
			wantMethod: http.MethodDelete,
			wantPath:   "/_search/scroll",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClearScroll(
					context.Background(),
					func(r *core_clear_scroll.ClearScroll) {
						req := core_clear_scroll.NewRequest()
						req.ScrollId = []string{"scroll-1"}
						r.Request(req)
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "UpdateByQuery",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_update_by_query",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.UpdateByQuery(
					context.Background(),
					estype.Index("products-000001"),
					func(r *core_update_by_query.UpdateByQuery) {
						r.Request(&core_update_by_query.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "DeleteByQuery",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_delete_by_query",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DeleteByQuery(
					context.Background(),
					estype.Index("products-000001"),
					func(r *core_delete_by_query.DeleteByQuery) {
						r.Request(&core_delete_by_query.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetMapping",
			wantMethod: http.MethodGet,
			wantPath:   "/products-000001/_mapping",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetMapping(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutMapping",
			wantMethod: http.MethodPut,
			wantPath:   "/products-000001/_mapping",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutMapping(
					context.Background(),
					estype.Index("products-000001"),
					func(r *indices_put_mapping.PutMapping) {
						r.Request(&indices_put_mapping.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetSettings",
			wantMethod: http.MethodGet,
			wantPath:   "/products-000001/_settings",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetSettings(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutSettings",
			wantMethod: http.MethodPut,
			wantPath:   "/products-000001/_settings",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutSettings(
					context.Background(),
					estype.Index("products-000001"),
					func(r *indices_put_settings.PutSettings) {
						r.Request(&indices_put_settings.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "OpenIndex",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_open",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.OpenIndex(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CloseIndex",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_close",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CloseIndex(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "Flush",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_flush",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Flush(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClearCache",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_cache/clear",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClearCache(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "ForceMerge",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_forcemerge",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ForceMerge(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "Rollover",
			wantMethod: http.MethodPost,
			wantPath:   "/products-write/_rollover",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Rollover(context.Background(), estype.Alias("products-write"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesStats",
			wantMethod: http.MethodGet,
			wantPath:   "/products-000001/_stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesStats(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClusterHealth",
			wantMethod: http.MethodGet,
			wantPath:   "/_cluster/health/products-000001",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterHealth(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutIndexTemplate",
			wantMethod: http.MethodPut,
			wantPath:   "/_index_template/template-one",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutIndexTemplate(
					context.Background(),
					estype.Template("template-one"),
					func(r *indices_put_index_template.PutIndexTemplate) {
						r.Request(&indices_put_index_template.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "Analyze",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_analyze",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Analyze(
					context.Background(),
					estype.Index("products-000001"),
					func(r *indices_analyze.Analyze) {
						r.Request(&indices_analyze.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClosePointInTime",
			wantMethod: http.MethodDelete,
			wantPath:   "/_pit",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClosePointInTime(
					context.Background(),
					func(r *core_close_point_in_time.ClosePointInTime) {
						req := core_close_point_in_time.NewRequest()
						req.Id = "pit-1"
						r.Request(req)
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutLifecycle",
			wantMethod: http.MethodPut,
			wantPath:   "/_ilm/policy/policy-one",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutLifecycle(
					context.Background(),
					estype.Policy("policy-one"),
					func(r *ilm_put_lifecycle.PutLifecycle) {
						r.Request(&ilm_put_lifecycle.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutPipeline",
			wantMethod: http.MethodPut,
			wantPath:   "/_ingest/pipeline/pipeline-one",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutPipeline(
					context.Background(),
					estype.Pipeline("pipeline-one"),
					func(r *ingest_put_pipeline.PutPipeline) {
						r.Request(&ingest_put_pipeline.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "EsqlQuery",
			wantMethod: http.MethodPost,
			wantPath:   "/_query",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EsqlQuery(context.Background(), estype.ESQLQuery("from logs | limit 1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutInference",
			wantMethod: http.MethodPut,
			wantPath:   "/_inference/inference-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutInference(
					context.Background(),
					estype.InferenceID("inference-1"),
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlPutJob",
			wantMethod: http.MethodPut,
			wantPath:   "/_ml/anomaly_detectors/job-one",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlPutJob(
					context.Background(),
					estype.MLJobID("job-one"),
					func(r *ml_put_job.PutJob) {
						r.Request(&ml_put_job.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlPutDatafeed",
			wantMethod: http.MethodPut,
			wantPath:   "/_ml/datafeeds/datafeed-one",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlPutDatafeed(
					context.Background(),
					estype.DatafeedID("datafeed-one"),
					func(r *ml_put_datafeed.PutDatafeed) {
						r.Request(&ml_put_datafeed.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrFollow",
			wantMethod: http.MethodPut,
			wantPath:   "/follower-index/_ccr/follow",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrFollow(
					context.Background(),
					estype.Index("follower-index"),
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrResumeFollow",
			wantMethod: http.MethodPost,
			wantPath:   "/follower-index/_ccr/resume_follow",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrResumeFollow(context.Background(), estype.Index("follower-index"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutTransform",
			wantMethod: http.MethodPut,
			wantPath:   "/_transform/transform-one",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutTransform(
					context.Background(),
					estype.TransformID("transform-one"),
					func(r *transform_put_transform.PutTransform) {
						r.Request(&transform_put_transform.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetIndexTemplate",
			wantMethod: http.MethodGet,
			wantPath:   "/_index_template/products-template",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetIndexTemplate(context.Background(), estype.Template("products-template"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "DeleteIndexTemplate",
			wantMethod: http.MethodDelete,
			wantPath:   "/_index_template/products-template",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DeleteIndexTemplate(context.Background(), estype.Template("products-template"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "ExistsIndexTemplate",
			wantMethod: http.MethodHead,
			wantPath:   "/_index_template/products-template",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.ExistsIndexTemplate(context.Background(), estype.Template("products-template"))
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "GetAlias",
			wantMethod: http.MethodGet,
			wantPath:   "/_alias/products-write",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetAlias(context.Background(), estype.Alias("products-write"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "DeleteAlias",
			wantMethod: http.MethodDelete,
			wantPath:   "/products-000001/_alias/products-write",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DeleteAlias(
					context.Background(),
					estype.Index("products-000001"),
					estype.Alias("products-write"),
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "TasksList",
			wantMethod: http.MethodGet,
			wantPath:   "/_tasks",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TasksList(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "TasksCancel",
			wantMethod: http.MethodPost,
			wantPath:   "/_tasks/node:1/_cancel",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TasksCancel(context.Background(), estype.TaskID("node:1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "OpenPointInTime",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_pit",
			wantQuery:  "keep_alive=1m",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.OpenPointInTime(
					context.Background(),
					estype.Index("products-000001"),
					estype.KeepAlive("1m"),
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "Analyze",
			wantMethod: http.MethodPost,
			wantPath:   "/products-000001/_analyze",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Analyze(
					context.Background(),
					estype.Index("products-000001"),
					func(r *indices_analyze.Analyze) {
						r.Request(&indices_analyze.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClosePointInTime",
			wantMethod: http.MethodDelete,
			wantPath:   "/_pit",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClosePointInTime(
					context.Background(),
					func(r *core_close_point_in_time.ClosePointInTime) {
						req := core_close_point_in_time.NewRequest()
						req.Id = "pit-1"
						r.Request(req)
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "CreateDataStream",
			wantMethod: http.MethodPut,
			wantPath:   "/_data_stream/logs-app",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CreateDataStream(context.Background(), estype.DataStream("logs-app"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetDataStream",
			wantMethod: http.MethodGet,
			wantPath:   "/_data_stream/logs-app",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetDataStream(context.Background(), estype.DataStream("logs-app"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "DeleteDataStream",
			wantMethod: http.MethodDelete,
			wantPath:   "/_data_stream/logs-app",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DeleteDataStream(context.Background(), estype.DataStream("logs-app"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutLifecycle",
			wantMethod: http.MethodPut,
			wantPath:   "/_ilm/policy/logs-policy",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutLifecycle(
					context.Background(),
					estype.Policy("logs-policy"),
					func(r *ilm_put_lifecycle.PutLifecycle) {
						r.Request(&ilm_put_lifecycle.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetLifecycle",
			wantMethod: http.MethodGet,
			wantPath:   "/_ilm/policy/logs-policy",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetLifecycle(context.Background(), estype.Policy("logs-policy"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "ExplainLifecycle",
			wantMethod: http.MethodGet,
			wantPath:   "/products-000001/_ilm/explain",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ExplainLifecycle(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetPipeline",
			wantMethod: http.MethodGet,
			wantPath:   "/_ingest/pipeline/enrich-products",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetPipeline(context.Background(), estype.Pipeline("enrich-products"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "DeletePipeline",
			wantMethod: http.MethodDelete,
			wantPath:   "/_ingest/pipeline/enrich-products",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DeletePipeline(context.Background(), estype.Pipeline("enrich-products"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutPipeline",
			wantMethod: http.MethodPut,
			wantPath:   "/_ingest/pipeline/enrich-products",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutPipeline(
					context.Background(),
					estype.Pipeline("enrich-products"),
					func(r *ingest_put_pipeline.PutPipeline) {
						r.Request(&ingest_put_pipeline.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "EsqlQuery",
			wantMethod: http.MethodPost,
			wantPath:   "/_query",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EsqlQuery(context.Background(), estype.ESQLQuery("from products"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CreateApiKey",
			wantMethod: http.MethodPut,
			wantPath:   "/_security/api_key",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CreateApiKey(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetApiKey",
			wantMethod: http.MethodGet,
			wantPath:   "/_security/api_key",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetApiKey(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "InvalidateApiKey",
			wantMethod: http.MethodDelete,
			wantPath:   "/_security/api_key",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.InvalidateApiKey(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CreateRepository",
			wantMethod: http.MethodPut,
			wantPath:   "/_snapshot/repo-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CreateRepository(context.Background(), estype.Repository("repo-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CreateSnapshot",
			wantMethod: http.MethodPut,
			wantPath:   "/_snapshot/repo-1/snap-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CreateSnapshot(
					context.Background(),
					estype.Repository("repo-1"),
					estype.Snapshot("snap-1"),
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "RestoreSnapshot",
			wantMethod: http.MethodPost,
			wantPath:   "/_snapshot/repo-1/snap-1/_restore",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.RestoreSnapshot(
					context.Background(),
					estype.Repository("repo-1"),
					estype.Snapshot("snap-1"),
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetInference",
			wantMethod: http.MethodGet,
			wantPath:   "/_inference/inference-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetInference(context.Background(), estype.InferenceID("inference-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "DeleteInference",
			wantMethod: http.MethodDelete,
			wantPath:   "/_inference/inference-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DeleteInference(context.Background(), estype.InferenceID("inference-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutInference",
			wantMethod: http.MethodPut,
			wantPath:   "/_inference/inference-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutInference(context.Background(), estype.InferenceID("inference-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "Inference",
			wantMethod: http.MethodPost,
			wantPath:   "/_inference/inference-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Inference(
					context.Background(),
					estype.InferenceID("inference-1"),
					func(r *inference_inference.Inference) {
						r.Request(&inference_inference.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlPutJob",
			wantMethod: http.MethodPut,
			wantPath:   "/_ml/anomaly_detectors/job-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlPutJob(
					context.Background(),
					estype.MLJobID("job-1"),
					func(r *ml_put_job.PutJob) {
						r.Request(&ml_put_job.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlGetJobs",
			wantMethod: http.MethodGet,
			wantPath:   "/_ml/anomaly_detectors/job-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlGetJobs(context.Background(), estype.MLJobID("job-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlDeleteJob",
			wantMethod: http.MethodDelete,
			wantPath:   "/_ml/anomaly_detectors/job-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlDeleteJob(context.Background(), estype.MLJobID("job-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlOpenJob",
			wantMethod: http.MethodPost,
			wantPath:   "/_ml/anomaly_detectors/job-1/_open",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlOpenJob(context.Background(), estype.MLJobID("job-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlCloseJob",
			wantMethod: http.MethodPost,
			wantPath:   "/_ml/anomaly_detectors/job-1/_close",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlCloseJob(context.Background(), estype.MLJobID("job-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlPutDatafeed",
			wantMethod: http.MethodPut,
			wantPath:   "/_ml/datafeeds/datafeed-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlPutDatafeed(
					context.Background(),
					estype.DatafeedID("datafeed-1"),
					func(r *ml_put_datafeed.PutDatafeed) {
						r.Request(&ml_put_datafeed.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlGetDatafeeds",
			wantMethod: http.MethodGet,
			wantPath:   "/_ml/datafeeds/datafeed-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlGetDatafeeds(context.Background(), estype.DatafeedID("datafeed-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlDeleteDatafeed",
			wantMethod: http.MethodDelete,
			wantPath:   "/_ml/datafeeds/datafeed-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlDeleteDatafeed(context.Background(), estype.DatafeedID("datafeed-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlStartDatafeed",
			wantMethod: http.MethodPost,
			wantPath:   "/_ml/datafeeds/datafeed-1/_start",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlStartDatafeed(context.Background(), estype.DatafeedID("datafeed-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "MlStopDatafeed",
			wantMethod: http.MethodPost,
			wantPath:   "/_ml/datafeeds/datafeed-1/_stop",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.MlStopDatafeed(context.Background(), estype.DatafeedID("datafeed-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrFollow",
			wantMethod: http.MethodPut,
			wantPath:   "/follower-idx/_ccr/follow",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrFollow(context.Background(), estype.Index("follower-idx"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrPauseFollow",
			wantMethod: http.MethodPost,
			wantPath:   "/follower-idx/_ccr/pause_follow",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrPauseFollow(context.Background(), estype.Index("follower-idx"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrResumeFollow",
			wantMethod: http.MethodPost,
			wantPath:   "/follower-idx/_ccr/resume_follow",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrResumeFollow(context.Background(), estype.Index("follower-idx"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrUnfollow",
			wantMethod: http.MethodPost,
			wantPath:   "/follower-idx/_ccr/unfollow",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrUnfollow(context.Background(), estype.Index("follower-idx"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrFollowStats",
			wantMethod: http.MethodGet,
			wantPath:   "/follower-idx/_ccr/stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrFollowStats(context.Background(), estype.Index("follower-idx"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutTransform",
			wantMethod: http.MethodPut,
			wantPath:   "/_transform/transform-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.PutTransform(
					context.Background(),
					estype.TransformID("transform-1"),
					func(r *transform_put_transform.PutTransform) {
						r.Request(&transform_put_transform.Request{})
					},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetTransform",
			wantMethod: http.MethodGet,
			wantPath:   "/_transform/transform-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetTransform(context.Background(), estype.TransformID("transform-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "DeleteTransform",
			wantMethod: http.MethodDelete,
			wantPath:   "/_transform/transform-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DeleteTransform(context.Background(), estype.TransformID("transform-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "StartTransform",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/transform-1/_start",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.StartTransform(context.Background(), estype.TransformID("transform-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "StopTransform",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/transform-1/_stop",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.StopTransform(context.Background(), estype.TransformID("transform-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "GetTransformStats",
			wantMethod: http.MethodGet,
			wantPath:   "/_transform/transform-1/_stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.GetTransformStats(context.Background(), estype.TransformID("transform-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatAliases",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/aliases/products-write",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatAliases(context.Background(), estype.Alias("products-write"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatIndices",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/indices/products-000001",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatIndices(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatNodes",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/nodes",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatNodes(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatAllocation",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/allocation",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatAllocation(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatCircuitBreaker",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/circuit_breaker",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatCircuitBreaker(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatComponentTemplates",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/component_templates/products-template",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatComponentTemplates(context.Background(), estype.Template("products-template"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatCount",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/count/products-000001",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatCount(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatFielddata",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/fielddata/status,title",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatFielddata(
					context.Background(),
					[]estype.Field{estype.Field("status"), estype.Field("title")},
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatHealth",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/health",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatHealth(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatHelp",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatHelp(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatMaster",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/master",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatMaster(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatMlDataFrameAnalytics",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/ml/data_frame/analytics/dfa-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatMlDataFrameAnalytics(context.Background(), estype.DataFrameAnalyticsID("dfa-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatMlDatafeeds",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/ml/datafeeds/datafeed-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatMlDatafeeds(context.Background(), estype.DatafeedID("datafeed-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatMlJobs",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/ml/anomaly_detectors/job-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatMlJobs(context.Background(), estype.MLJobID("job-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatMlTrainedModels",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/ml/trained_models/model-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatMlTrainedModels(context.Background(), estype.TrainedModelID("model-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatNodeattrs",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/nodeattrs",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatNodeattrs(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatPendingTasks",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/pending_tasks",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatPendingTasks(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatPlugins",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/plugins",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatPlugins(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatRecovery",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/recovery/products-000001",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatRecovery(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatRepositories",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/repositories",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatRepositories(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatSegments",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/segments/products-000001",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatSegments(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatShards",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/shards/products-000001",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatShards(context.Background(), estype.Index("products-000001"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatSnapshots",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/snapshots/repo-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatSnapshots(context.Background(), estype.Repository("repo-1"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatTasks",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/tasks",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatTasks(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatTemplates",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/templates/products-template",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatTemplates(context.Background(), estype.Template("products-template"))
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatThreadPool",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/thread_pool",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatThreadPool(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CatTransforms",
			wantMethod: http.MethodGet,
			wantPath:   "/_cat/transforms/transform-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CatTransforms(context.Background(), estype.TransformID("transform-1"))
				assert.NilError(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var seenMethod string
			var seenPath string
			var seenQuery string

			client := newTestESClient(t, testRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				seenMethod = req.Method
				seenPath = req.URL.Path
				seenQuery = req.URL.RawQuery

				if req.Body != nil {
					_, err := io.ReadAll(req.Body)
					assert.NilError(t, err)
				}

				return smokeResponseForRequest(req), nil
			}))

			tt.call(t, client)

			assert.Equal(t, tt.wantMethod, seenMethod)
			assert.Equal(t, tt.wantPath, seenPath)
			if tt.wantQuery != "" {
				assert.Equal(t, tt.wantQuery, seenQuery)
			}
		})
	}
}

func TestClientSpecSmoke_SelectedMethods(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		wantMethod string
		wantPath   string
		call       func(*testing.T, *esClient)
	}

	tests := []testCase{
		{
			name:       "AsyncSearchDelete",
			wantMethod: http.MethodDelete,
			wantPath:   "/_async_search/async-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.AsyncSearchDelete(context.Background(), "async-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "AsyncSearchGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_async_search/async-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.AsyncSearchGet(context.Background(), "async-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "AsyncSearchStatus",
			wantMethod: http.MethodGet,
			wantPath:   "/_async_search/status/async-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.AsyncSearchStatus(context.Background(), "async-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "AsyncSearchSubmit",
			wantMethod: http.MethodPost,
			wantPath:   "/_async_search",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.AsyncSearchSubmit(context.Background(), &async_search_submit.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "Capabilities",
			wantMethod: http.MethodGet,
			wantPath:   "/_capabilities",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.Capabilities(context.Background())
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "ClusterDeleteComponentTemplate",
			wantMethod: http.MethodDelete,
			wantPath:   "/_component_template/component-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterDeleteComponentTemplate(context.Background(), "component-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClusterDeleteVotingConfigExclusions",
			wantMethod: http.MethodDelete,
			wantPath:   "/_cluster/voting_config_exclusions",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.ClusterDeleteVotingConfigExclusions(context.Background())
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "ClusterExistsComponentTemplate",
			wantMethod: http.MethodHead,
			wantPath:   "/_component_template/component-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.ClusterExistsComponentTemplate(context.Background(), "component-1")
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "ClusterGetComponentTemplate",
			wantMethod: http.MethodGet,
			wantPath:   "/_component_template",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterGetComponentTemplate(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClusterGetSettings",
			wantMethod: http.MethodGet,
			wantPath:   "/_cluster/settings",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterGetSettings(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClusterInfo",
			wantMethod: http.MethodGet,
			wantPath:   "/_info/ingest",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterInfo(context.Background(), "ingest")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClusterPendingTasks",
			wantMethod: http.MethodGet,
			wantPath:   "/_cluster/pending_tasks",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterPendingTasks(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClusterPostVotingConfigExclusions",
			wantMethod: http.MethodPost,
			wantPath:   "/_cluster/voting_config_exclusions",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.ClusterPostVotingConfigExclusions(context.Background())
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "ClusterRemoteInfo",
			wantMethod: http.MethodGet,
			wantPath:   "/_remote/info",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterRemoteInfo(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClusterState",
			wantMethod: http.MethodGet,
			wantPath:   "/_cluster/state",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterState(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClusterStats",
			wantMethod: http.MethodGet,
			wantPath:   "/_cluster/stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterStats(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "Create",
			wantMethod: http.MethodPut,
			wantPath:   "/products/_create/doc-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				req := core_create.NewRequest()
				*req = core_create.Request(`{"status":"active"}`)

				_, err := client.Create(context.Background(), "products", "doc-1", req)
				assert.NilError(t, err)
			},
		},
		{
			name:       "Delete",
			wantMethod: http.MethodDelete,
			wantPath:   "/products/_doc/doc-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.Delete(context.Background(), "products", "doc-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "Exists",
			wantMethod: http.MethodHead,
			wantPath:   "/products/_doc/doc-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.Exists(context.Background(), "products", "doc-1")
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "ExistsSource",
			wantMethod: http.MethodHead,
			wantPath:   "/products/_source/doc-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.ExistsSource(context.Background(), "products", "doc-1")
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "FieldCaps",
			wantMethod: http.MethodPost,
			wantPath:   "/_field_caps",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.FieldCaps(context.Background(), &core_field_caps.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorList",
			wantMethod: http.MethodGet,
			wantPath:   "/_connector",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorList(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "AutoscalingGetAutoscalingCapacity",
			wantMethod: http.MethodGet,
			wantPath:   "/_autoscaling/capacity",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.AutoscalingGetAutoscalingCapacity(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "AutoscalingGetAutoscalingPolicy",
			wantMethod: http.MethodGet,
			wantPath:   "/_autoscaling/policy/policy-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.AutoscalingGetAutoscalingPolicy(context.Background(), "policy-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "AutoscalingPutAutoscalingPolicy",
			wantMethod: http.MethodPut,
			wantPath:   "/_autoscaling/policy/policy-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.AutoscalingPutAutoscalingPolicy(context.Background(), "policy-1", &autoscaling_put_autoscaling_policy.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "AutoscalingDeleteAutoscalingPolicy",
			wantMethod: http.MethodDelete,
			wantPath:   "/_autoscaling/policy/policy-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.AutoscalingDeleteAutoscalingPolicy(context.Background(), "policy-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrDeleteAutoFollowPattern",
			wantMethod: http.MethodDelete,
			wantPath:   "/_ccr/auto_follow/pattern-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrDeleteAutoFollowPattern(context.Background(), "pattern-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrGetAutoFollowPattern",
			wantMethod: http.MethodGet,
			wantPath:   "/_ccr/auto_follow",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrGetAutoFollowPattern(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrPauseAutoFollowPattern",
			wantMethod: http.MethodPost,
			wantPath:   "/_ccr/auto_follow/pattern-1/pause",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrPauseAutoFollowPattern(context.Background(), "pattern-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrPutAutoFollowPattern",
			wantMethod: http.MethodPut,
			wantPath:   "/_ccr/auto_follow/pattern-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrPutAutoFollowPattern(context.Background(), "pattern-1", &ccr_put_auto_follow_pattern.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrResumeAutoFollowPattern",
			wantMethod: http.MethodPost,
			wantPath:   "/_ccr/auto_follow/pattern-1/resume",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrResumeAutoFollowPattern(context.Background(), "pattern-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "CcrStats",
			wantMethod: http.MethodGet,
			wantPath:   "/_ccr/stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrStats(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "DanglingIndicesDeleteDanglingIndex",
			wantMethod: http.MethodDelete,
			wantPath:   "/_dangling/uuid-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DanglingIndicesDeleteDanglingIndex(context.Background(), "uuid-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "DanglingIndicesImportDanglingIndex",
			wantMethod: http.MethodPost,
			wantPath:   "/_dangling/uuid-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DanglingIndicesImportDanglingIndex(context.Background(), "uuid-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "DanglingIndicesListDanglingIndices",
			wantMethod: http.MethodGet,
			wantPath:   "/_dangling",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.DanglingIndicesListDanglingIndices(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "EnrichDeletePolicy",
			wantMethod: http.MethodDelete,
			wantPath:   "/_enrich/policy/policy-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EnrichDeletePolicy(context.Background(), "policy-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "EnrichExecutePolicy",
			wantMethod: http.MethodPut,
			wantPath:   "/_enrich/policy/policy-1/_execute",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EnrichExecutePolicy(context.Background(), "policy-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "EnrichGetPolicy",
			wantMethod: http.MethodGet,
			wantPath:   "/_enrich/policy",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EnrichGetPolicy(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "EnrichPutPolicy",
			wantMethod: http.MethodPut,
			wantPath:   "/_enrich/policy/policy-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EnrichPutPolicy(context.Background(), "policy-1", &enrich_put_policy.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "EnrichStats",
			wantMethod: http.MethodGet,
			wantPath:   "/_enrich/_stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EnrichStats(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "EqlDelete",
			wantMethod: http.MethodDelete,
			wantPath:   "/_eql/search/async-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EqlDelete(context.Background(), "async-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "EqlGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_eql/search/async-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EqlGet(context.Background(), "async-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "EqlGetStatus",
			wantMethod: http.MethodGet,
			wantPath:   "/_eql/search/status/async-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EqlGetStatus(context.Background(), "async-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "EsqlAsyncQueryDelete",
			wantMethod: http.MethodDelete,
			wantPath:   "/_query/async/async-esql-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EsqlAsyncQueryDelete(context.Background(), "async-esql-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "EsqlAsyncQueryGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_query/async/async-esql-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EsqlAsyncQueryGet(context.Background(), "async-esql-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "EsqlAsyncQueryStop",
			wantMethod: http.MethodPost,
			wantPath:   "/_query/async/async-esql-1/stop",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EsqlAsyncQueryStop(context.Background(), "async-esql-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "FeaturesGetFeatures",
			wantMethod: http.MethodGet,
			wantPath:   "/_features",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.FeaturesGetFeatures(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "FeaturesResetFeatures",
			wantMethod: http.MethodPost,
			wantPath:   "/_features/_reset",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.FeaturesResetFeatures(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesDataStreamsStats",
			wantMethod: http.MethodGet,
			wantPath:   "/_data_stream/_stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesDataStreamsStats(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "ClusterPutSettings",
			wantMethod: http.MethodPut,
			wantPath:   "/_cluster/settings",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ClusterPutSettings(context.Background(), &cluster_put_settings.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "EqlSearch",
			wantMethod: http.MethodPost,
			wantPath:   "/products/_eql/search",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EqlSearch(context.Background(), "products", &eql_search.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "EsqlAsyncQuery",
			wantMethod: http.MethodPost,
			wantPath:   "/_query/async",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.EsqlAsyncQuery(context.Background(), &esql_async_query.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorCheckIn",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_check_in",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorCheckIn(context.Background(), "connector-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorDelete",
			wantMethod: http.MethodDelete,
			wantPath:   "/_connector/connector-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorDelete(context.Background(), "connector-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_connector/connector-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorGet(context.Background(), "connector-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorLastSync",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_last_sync",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorLastSync(context.Background(), "connector-1", &connector_last_sync.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorPost",
			wantMethod: http.MethodPost,
			wantPath:   "/_connector",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorPost(context.Background(), &connector_post.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorPut",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorPut(context.Background(), &connector_put.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSecretDelete",
			wantMethod: http.MethodDelete,
			wantPath:   "/_connector/_secret/secret-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSecretDelete(context.Background(), "secret-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSecretGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_connector/_secret/secret-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSecretGet(context.Background(), "secret-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSecretPost",
			wantMethod: http.MethodPost,
			wantPath:   "/_connector/_secret",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.ConnectorSecretPost(context.Background())
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "ConnectorSecretPut",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/_secret/secret-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSecretPut(context.Background(), "secret-1", []byte(`{}`))
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSyncJobCancel",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/_sync_job/job-1/_cancel",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSyncJobCancel(context.Background(), "job-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSyncJobCheckIn",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/_sync_job/job-1/_check_in",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSyncJobCheckIn(context.Background(), "job-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSyncJobClaim",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/_sync_job/job-1/_claim",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSyncJobClaim(context.Background(), "job-1", &connector_sync_job_claim.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSyncJobDelete",
			wantMethod: http.MethodDelete,
			wantPath:   "/_connector/_sync_job/job-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSyncJobDelete(context.Background(), "job-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSyncJobError",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/_sync_job/job-1/_error",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSyncJobError(context.Background(), "job-1", &connector_sync_job_error.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSyncJobGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_connector/_sync_job/job-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSyncJobGet(context.Background(), "job-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSyncJobList",
			wantMethod: http.MethodGet,
			wantPath:   "/_connector/_sync_job",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSyncJobList(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSyncJobPost",
			wantMethod: http.MethodPost,
			wantPath:   "/_connector/_sync_job",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSyncJobPost(context.Background(), &connector_sync_job_post.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorSyncJobUpdateStats",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/_sync_job/job-1/_stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorSyncJobUpdateStats(context.Background(), "job-1", &connector_sync_job_update_stats.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateActiveFiltering",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_filtering/_activate",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateActiveFiltering(context.Background(), "connector-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateApiKeyId",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_api_key_id",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateApiKeyId(context.Background(), "connector-1", &connector_update_api_key_id.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateConfiguration",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_configuration",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateConfiguration(context.Background(), "connector-1", &connector_update_configuration.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateError",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_error",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateError(context.Background(), "connector-1", &connector_update_error.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateFeatures",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_features",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateFeatures(context.Background(), "connector-1", &connector_update_features.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateFiltering",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_filtering",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateFiltering(context.Background(), "connector-1", &connector_update_filtering.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateFilteringValidation",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_filtering/_validation",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateFilteringValidation(context.Background(), "connector-1", &connector_update_filtering_validation.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateIndexName",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_index_name",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateIndexName(context.Background(), "connector-1", &connector_update_index_name.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateName",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_name",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateName(context.Background(), "connector-1", &connector_update_name.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateNative",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_native",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateNative(context.Background(), "connector-1", &connector_update_native.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdatePipeline",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_pipeline",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdatePipeline(context.Background(), "connector-1", &connector_update_pipeline.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateScheduling",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_scheduling",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateScheduling(context.Background(), "connector-1", &connector_update_scheduling.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateServiceType",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_service_type",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateServiceType(context.Background(), "connector-1", &connector_update_service_type.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "ConnectorUpdateStatus",
			wantMethod: http.MethodPut,
			wantPath:   "/_connector/connector-1/_status",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.ConnectorUpdateStatus(context.Background(), "connector-1", &connector_update_status.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesDelete",
			wantMethod: http.MethodDelete,
			wantPath:   "/products",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesDelete(context.Background(), "products")
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesDeleteTemplate",
			wantMethod: http.MethodDelete,
			wantPath:   "/_template/products-template",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesDeleteTemplate(context.Background(), "products-template")
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesExists",
			wantMethod: http.MethodHead,
			wantPath:   "/products",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.IndicesExists(context.Background(), "products")
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "IndicesExplainDataLifecycle",
			wantMethod: http.MethodGet,
			wantPath:   "/products/_lifecycle/explain",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesExplainDataLifecycle(context.Background(), "products")
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesGet",
			wantMethod: http.MethodGet,
			wantPath:   "/products",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesGet(context.Background(), "products")
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesGetDataLifecycle",
			wantMethod: http.MethodGet,
			wantPath:   "/_data_stream/products/_lifecycle",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesGetDataLifecycle(context.Background(), "products")
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesGetDataLifecycleStats",
			wantMethod: http.MethodGet,
			wantPath:   "/_lifecycle/stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesGetDataLifecycleStats(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesGetFieldMapping",
			wantMethod: http.MethodGet,
			wantPath:   "/_mapping/field/status",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesGetFieldMapping(context.Background(), "status")
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesGetTemplate",
			wantMethod: http.MethodGet,
			wantPath:   "/_template",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesGetTemplate(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesPutTemplate",
			wantMethod: http.MethodPut,
			wantPath:   "/_template/products-template",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesPutTemplate(context.Background(), "products-template", &indices_put_template.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesRecovery",
			wantMethod: http.MethodGet,
			wantPath:   "/_recovery",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesRecovery(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesRefresh",
			wantMethod: http.MethodPost,
			wantPath:   "/_refresh",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesRefresh(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesReloadSearchAnalyzers",
			wantMethod: http.MethodPost,
			wantPath:   "/products/_reload_search_analyzers",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesReloadSearchAnalyzers(context.Background(), "products")
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesResolveCluster",
			wantMethod: http.MethodGet,
			wantPath:   "/_resolve/cluster",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesResolveCluster(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesResolveIndex",
			wantMethod: http.MethodGet,
			wantPath:   "/_resolve/index/products",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesResolveIndex(context.Background(), "products")
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesSegments",
			wantMethod: http.MethodGet,
			wantPath:   "/_segments",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesSegments(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesShardStores",
			wantMethod: http.MethodGet,
			wantPath:   "/_shard_stores",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesShardStores(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "IndicesValidateQuery",
			wantMethod: http.MethodPost,
			wantPath:   "/_validate/query",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.IndicesValidateQuery(context.Background(), &indices_validate_query.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "LicenseGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_license",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.LicenseGet(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "NodesInfo",
			wantMethod: http.MethodGet,
			wantPath:   "/_nodes",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.NodesInfo(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "NodesStats",
			wantMethod: http.MethodGet,
			wantPath:   "/_nodes/stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.NodesStats(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "QueryRulesDeleteRule",
			wantMethod: http.MethodDelete,
			wantPath:   "/_query_rules/ruleset-1/_rule/rule-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.QueryRulesDeleteRule(context.Background(), "ruleset-1", "rule-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "QueryRulesDeleteRuleset",
			wantMethod: http.MethodDelete,
			wantPath:   "/_query_rules/ruleset-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.QueryRulesDeleteRuleset(context.Background(), "ruleset-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "QueryRulesGetRule",
			wantMethod: http.MethodGet,
			wantPath:   "/_query_rules/ruleset-1/_rule/rule-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.QueryRulesGetRule(context.Background(), "ruleset-1", "rule-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "QueryRulesGetRuleset",
			wantMethod: http.MethodGet,
			wantPath:   "/_query_rules/ruleset-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.QueryRulesGetRuleset(context.Background(), "ruleset-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "QueryRulesListRulesets",
			wantMethod: http.MethodGet,
			wantPath:   "/_query_rules",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.QueryRulesListRulesets(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "QueryRulesPutRule",
			wantMethod: http.MethodPut,
			wantPath:   "/_query_rules/ruleset-1/_rule/rule-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.QueryRulesPutRule(context.Background(), "ruleset-1", "rule-1", &query_rules_put_rule.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "QueryRulesPutRuleset",
			wantMethod: http.MethodPut,
			wantPath:   "/_query_rules/ruleset-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.QueryRulesPutRuleset(context.Background(), "ruleset-1", &query_rules_put_ruleset.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "QueryRulesTest",
			wantMethod: http.MethodPost,
			wantPath:   "/_query_rules/ruleset-1/_test",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.QueryRulesTest(context.Background(), "ruleset-1", &query_rules_test.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationDelete",
			wantMethod: http.MethodDelete,
			wantPath:   "/_application/search_application/app-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SearchApplicationDelete(context.Background(), "app-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationDeleteBehavioralAnalytics",
			wantMethod: http.MethodDelete,
			wantPath:   "/_application/analytics/collection-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SearchApplicationDeleteBehavioralAnalytics(context.Background(), "collection-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_application/search_application/app-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SearchApplicationGet(context.Background(), "app-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationGetBehavioralAnalytics",
			wantMethod: http.MethodGet,
			wantPath:   "/_application/analytics",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SearchApplicationGetBehavioralAnalytics(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationList",
			wantMethod: http.MethodGet,
			wantPath:   "/_application/search_application",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SearchApplicationList(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationPostBehavioralAnalyticsEvent",
			wantMethod: http.MethodPost,
			wantPath:   "/_application/analytics/collection-1/event/click",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				req := search_application_post_behavioral_analytics_event.Request([]byte(`{"session":{"id":"session-1"},"user":{"id":"user-1"}}`))

				_, err := client.SearchApplicationPostBehavioralAnalyticsEvent(
					context.Background(),
					"collection-1",
					"click",
					&req,
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationPut",
			wantMethod: http.MethodPut,
			wantPath:   "/_application/search_application/app-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SearchApplicationPut(context.Background(), "app-1", &search_application_put.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationPutBehavioralAnalytics",
			wantMethod: http.MethodPut,
			wantPath:   "/_application/analytics/collection-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SearchApplicationPutBehavioralAnalytics(context.Background(), "collection-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationRenderQuery",
			wantMethod: http.MethodPost,
			wantPath:   "/_application/search_application/app-1/_render_query",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SearchApplicationRenderQuery(context.Background(), "app-1", &search_application_render_query.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SearchApplicationSearch",
			wantMethod: http.MethodPost,
			wantPath:   "/_application/search_application/app-1/_search",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SearchApplicationSearch(context.Background(), "app-1", &search_application_search.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotCleanupRepository",
			wantMethod: http.MethodPost,
			wantPath:   "/_snapshot/repo-1/_cleanup",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotCleanupRepository(context.Background(), "repo-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotClone",
			wantMethod: http.MethodPut,
			wantPath:   "/_snapshot/repo-1/snap-1/_clone/snap-2",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotClone(context.Background(), "repo-1", "snap-1", "snap-2", &snapshot_clone.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotCreate",
			wantMethod: http.MethodPut,
			wantPath:   "/_snapshot/repo-1/snap-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotCreate(context.Background(), "repo-1", "snap-1", &snapshot_create.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotCreateRepository",
			wantMethod: http.MethodPut,
			wantPath:   "/_snapshot/repo-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				var req snapshot_create_repository.Request
				_, err := client.SnapshotCreateRepository(context.Background(), "repo-1", &req)
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotDelete",
			wantMethod: http.MethodDelete,
			wantPath:   "/_snapshot/repo-1/snap-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotDelete(context.Background(), "repo-1", "snap-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotDeleteRepository",
			wantMethod: http.MethodDelete,
			wantPath:   "/_snapshot/repo-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotDeleteRepository(context.Background(), "repo-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_snapshot/repo-1/snap-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotGet(context.Background(), "repo-1", "snap-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotGetRepository",
			wantMethod: http.MethodGet,
			wantPath:   "/_snapshot",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotGetRepository(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotRepositoryAnalyze",
			wantMethod: http.MethodPost,
			wantPath:   "/_snapshot/repo-1/_analyze",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotRepositoryAnalyze(context.Background(), "repo-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotRepositoryVerifyIntegrity",
			wantMethod: http.MethodPost,
			wantPath:   "/_snapshot/repo-1/_verify_integrity",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotRepositoryVerifyIntegrity(context.Background(), "repo-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotRestore",
			wantMethod: http.MethodPost,
			wantPath:   "/_snapshot/repo-1/snap-1/_restore",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotRestore(context.Background(), "repo-1", "snap-1", &snapshot_restore.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotStatus",
			wantMethod: http.MethodGet,
			wantPath:   "/_snapshot/_status",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotStatus(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "SnapshotVerifyRepository",
			wantMethod: http.MethodPost,
			wantPath:   "/_snapshot/repo-1/_verify",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SnapshotVerifyRepository(context.Background(), "repo-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "SqlQuery",
			wantMethod: http.MethodPost,
			wantPath:   "/_sql",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SqlQuery(context.Background(), &sql_query.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SqlTranslate",
			wantMethod: http.MethodPost,
			wantPath:   "/_sql/translate",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SqlTranslate(context.Background(), &sql_translate.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityActivateUserProfile",
			wantMethod: http.MethodPost,
			wantPath:   "/_security/profile/_activate",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityActivateUserProfile(context.Background(), &security_activate_user_profile.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityAuthenticate",
			wantMethod: http.MethodGet,
			wantPath:   "/_security/_authenticate",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityAuthenticate(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityChangePassword",
			wantMethod: http.MethodPut,
			wantPath:   "/_security/user/_password",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityChangePassword(context.Background(), &security_change_password.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityCreateApiKey",
			wantMethod: http.MethodPut,
			wantPath:   "/_security/api_key",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityCreateApiKey(context.Background(), &security_create_api_key.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityCreateCrossClusterApiKey",
			wantMethod: http.MethodPost,
			wantPath:   "/_security/cross_cluster/api_key",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityCreateCrossClusterApiKey(context.Background(), &security_create_cross_cluster_api_key.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityDelegatePki",
			wantMethod: http.MethodPost,
			wantPath:   "/_security/delegate_pki",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityDelegatePki(context.Background(), &security_delegate_pki.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityGetApiKey",
			wantMethod: http.MethodGet,
			wantPath:   "/_security/api_key",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityGetApiKey(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityGetToken",
			wantMethod: http.MethodPost,
			wantPath:   "/_security/oauth2/token",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityGetToken(context.Background(), &security_get_token.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityGrantApiKey",
			wantMethod: http.MethodPost,
			wantPath:   "/_security/api_key/grant",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityGrantApiKey(context.Background(), &security_grant_api_key.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityInvalidateApiKey",
			wantMethod: http.MethodDelete,
			wantPath:   "/_security/api_key",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityInvalidateApiKey(context.Background(), &security_invalidate_api_key.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityPutRole",
			wantMethod: http.MethodPut,
			wantPath:   "/_security/role/role-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityPutRole(context.Background(), "role-1", &security_put_role.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityPutUser",
			wantMethod: http.MethodPut,
			wantPath:   "/_security/user/user-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityPutUser(context.Background(), "user-1", &security_put_user.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityQueryApiKeys",
			wantMethod: http.MethodPost,
			wantPath:   "/_security/_query/api_key",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityQueryApiKeys(context.Background(), &security_query_api_keys.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SecurityUpdateApiKey",
			wantMethod: http.MethodPut,
			wantPath:   "/_security/api_key/api-key-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SecurityUpdateApiKey(context.Background(), "api-key-1", &security_update_api_key.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "SslCertificates",
			wantMethod: http.MethodGet,
			wantPath:   "/_ssl/certificates",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SslCertificates(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "SynonymsGetSynonymsSets",
			wantMethod: http.MethodGet,
			wantPath:   "/_synonyms",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.SynonymsGetSynonymsSets(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "TasksGet",
			wantMethod: http.MethodGet,
			wantPath:   "/_tasks/node:1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TasksGet(context.Background(), "node:1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformDeleteTransform",
			wantMethod: http.MethodDelete,
			wantPath:   "/_transform/transform-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformDeleteTransform(context.Background(), "transform-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformGetNodeStats",
			wantMethod: http.MethodGet,
			wantPath:   "/_transform/_node_stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				got, err := client.TransformGetNodeStats(context.Background())
				assert.NilError(t, err)
				assert.Equal(t, true, got)
			},
		},
		{
			name:       "TransformGetTransform",
			wantMethod: http.MethodGet,
			wantPath:   "/_transform",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformGetTransform(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformGetTransformStats",
			wantMethod: http.MethodGet,
			wantPath:   "/_transform/transform-1/_stats",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformGetTransformStats(context.Background(), "transform-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformPreviewTransform",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/_preview",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformPreviewTransform(context.Background(), &transform_preview_transform.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformPutTransform",
			wantMethod: http.MethodPut,
			wantPath:   "/_transform/transform-1",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformPutTransform(context.Background(), "transform-1", &transform_put_transform.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformResetTransform",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/transform-1/_reset",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformResetTransform(context.Background(), "transform-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformScheduleNowTransform",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/transform-1/_schedule_now",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformScheduleNowTransform(context.Background(), "transform-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformSetUpgradeMode",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/_set_upgrade_mode",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformSetUpgradeMode(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformStartTransform",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/transform-1/_start",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformStartTransform(context.Background(), "transform-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformStopTransform",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/transform-1/_stop",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformStopTransform(context.Background(), "transform-1")
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformUpdateTransform",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/transform-1/_update",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformUpdateTransform(context.Background(), "transform-1", &transform_update_transform.Request{})
				assert.NilError(t, err)
			},
		},
		{
			name:       "TransformUpgradeTransforms",
			wantMethod: http.MethodPost,
			wantPath:   "/_transform/_upgrade",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.TransformUpgradeTransforms(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "XpackInfo",
			wantMethod: http.MethodGet,
			wantPath:   "/_xpack",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.XpackInfo(context.Background())
				assert.NilError(t, err)
			},
		},
		{
			name:       "XpackUsage",
			wantMethod: http.MethodGet,
			wantPath:   "/_xpack/usage",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.XpackUsage(context.Background())
				assert.NilError(t, err)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var seenMethod string
			var seenPath string

			client := newTestESClient(t, testRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				seenMethod = req.Method
				seenPath = req.URL.Path

				if req.Body != nil {
					_, err := io.ReadAll(req.Body)
					assert.NilError(t, err)
				}

				return smokeResponseForRequest(req), nil
			}))

			tt.call(t, client)

			assert.Equal(t, tt.wantMethod, seenMethod)
			assert.Equal(t, tt.wantPath, seenPath)
		})
	}
}

func smokeResponseForRequest(req *http.Request) *http.Response {
	if req.Method == http.MethodHead {
		return newJSONResponse(http.StatusOK, "")
	}

	switch req.URL.Path {
	case "/_capabilities":
		return newJSONResponse(http.StatusOK, "true")
	case "/_cat/circuit_breaker":
		return newJSONResponse(http.StatusOK, `{"ok":true}`)
	default:
		return newJSONResponse(http.StatusOK, smokeJSONBodyForPath(req.URL.Path))
	}
}

func smokeJSONBodyForPath(path string) string {
	switch path {
	case "/products/_mget":
		return `{"docs":[]}`
	case "/products/_msearch":
		return `{"took":1,"responses":[]}`
	case "/products/_count":
		return `{"count":0,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0}}`
	case "/_search/scroll":
		return `{"_scroll_id":"scroll-1","took":1,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":0,"relation":"eq"},"hits":[]}}`
	case "/products-000001/_update_by_query":
		return `{"took":1,"timed_out":false,"total":0,"updated":0,"deleted":0,"batches":0,"version_conflicts":0,"noops":0,"retries":{"bulk":0,"search":0},"throttled_millis":0,"requests_per_second":0,"throttled_until_millis":0,"failures":[]}`
	case "/products-000001/_delete_by_query":
		return `{"took":1,"timed_out":false,"total":0,"deleted":0,"batches":0,"version_conflicts":0,"noops":0,"retries":{"bulk":0,"search":0},"throttled_millis":0,"requests_per_second":0,"throttled_until_millis":0,"failures":[]}`
	case "/products-000001/_mapping":
		return `{}`
	case "/products-000001/_settings":
		return `{}`
	case "/_index_template/products-template":
		return `{"index_templates":[]}`
	case "/products-000001/_analyze":
		return `{"tokens":[]}`
	case "/_pit":
		return `{"succeeded":true,"num_freed":1}`
	case "/_alias/products-write":
		return `{}`
	case "/_tasks":
		return `{"nodes":{}}`
	case "/_data_stream/logs-app":
		return `{"data_streams":[]}`
	case "/_ilm/policy/logs-policy":
		return `{}`
	case "/products-000001/_ilm/explain":
		return `{"indices":{}}`
	case "/_ingest/pipeline/enrich-products":
		return `{}`
	case "/_query":
		return `{"columns":[],"values":[]}`
	case "/_security/api_key":
		return `{}`
	case "/_snapshot/repo-1":
		return `{}`
	case "/_snapshot/repo-1/snap-1":
		return `{"snapshot":{"snapshot":"snap-1","uuid":"u1","version_id":8000000,"version":"8.0.0","indices":[],"include_global_state":false,"state":"SUCCESS","start_time":"2024-01-01T00:00:00.000Z","start_time_in_millis":0,"end_time":"2024-01-01T00:00:00.000Z","end_time_in_millis":0,"duration_in_millis":0,"failures":[],"shards":{"total":0,"failed":0,"successful":0}}}`
	case "/_snapshot/repo-1/snap-1/_restore":
		return `{"accepted":true}`
	case "/_inference/inference-1":
		return `{}`
	case "/_inference/inference-1/_infer":
		return `{"inference_results":[]}`
	case "/_ml/anomaly_detectors/job-1":
		return `{"count":0,"jobs":[]}`
	case "/_ml/datafeeds/datafeed-1":
		return `{"count":0,"datafeeds":[]}`
	case "/follower-idx/_ccr/follow":
		return `{"follow_index_created":true,"follow_index_shards_acked":true,"index_following_started":true}`
	case "/follower-idx/_ccr/resume_follow":
		return `{"acknowledged":true}`
	case "/follower-idx/_ccr/stats":
		return `{}`
	case "/_transform/transform-1":
		return `{"count":0,"transforms":[]}`
	case "/_transform/transform-1/_stats":
		return `{"count":0,"transforms":[]}`
	case "/_cat/aliases/products-write",
		"/_cat/indices/products-000001",
		"/_cat/nodes",
		"/_cat/allocation",
		"/_cat/component_templates/products-template",
		"/_cat/count/products-000001",
		"/_cat/fielddata/status,title",
		"/_cat/health",
		"/_cat/master",
		"/_cat/ml/data_frame/analytics/dfa-1",
		"/_cat/ml/datafeeds/datafeed-1",
		"/_cat/ml/anomaly_detectors/job-1",
		"/_cat/ml/trained_models/model-1",
		"/_cat/nodeattrs",
		"/_cat/pending_tasks",
		"/_cat/plugins",
		"/_cat/recovery/products-000001",
		"/_cat/repositories",
		"/_cat/segments/products-000001",
		"/_cat/shards/products-000001",
		"/_cat/snapshots/repo-1",
		"/_cat/tasks",
		"/_cat/templates/products-template",
		"/_cat/thread_pool",
		"/_cat/transforms/transform-1":
		return `[]`
	case "/_cat":
		return `{}`
	case "/_async_search/async-1":
		return `{"id":"async-1","is_running":false,"is_partial":false,"response":{"took":1,"timed_out":false,"_shards":{"total":1,"successful":1,"skipped":0,"failed":0},"hits":{"total":{"value":0,"relation":"eq"},"hits":[]}}}`
	case "/_async_search/status/async-1":
		return `{"id":"async-1","is_running":false,"is_partial":false,"start_time_in_millis":0,"expiration_time_in_millis":0}`
	case "/_component_template":
		return `{"component_templates":[]}`
	case "/_autoscaling/capacity":
		return `{"policies":{}}`
	case "/_autoscaling/policy/policy-1":
		return `{"policies":[]}`
	case "/_cluster/settings":
		return `{}`
	case "/_info/ingest":
		return `{}`
	case "/_cluster/pending_tasks":
		return `{"tasks":[]}`
	case "/_remote/info":
		return `{}`
	case "/_cluster/state":
		return `{}`
	case "/_cluster/stats":
		return `{}`
	case "/_ccr/auto_follow":
		return `{"patterns":[]}`
	case "/_ccr/auto_follow/pattern-1":
		return `{}`
	case "/_ccr/auto_follow/pattern-1/pause":
		return `{"acknowledged":true}`
	case "/_ccr/auto_follow/pattern-1/resume":
		return `{"acknowledged":true}`
	case "/_ccr/stats":
		return `{}`
	case "/_dangling":
		return `{"dangling_indices":[]}`
	case "/_dangling/uuid-1":
		return `{"acknowledged":true}`
	case "/_enrich/policy":
		return `{"policies":[]}`
	case "/_enrich/policy/policy-1":
		return `{}`
	case "/_enrich/policy/policy-1/_execute":
		return `{"status":{"phase":"COMPLETE"}}`
	case "/_enrich/_stats":
		return `{"executing_policies":[],"coordinator_stats":[]}`
	case "/_eql/search/async-1":
		return `{"id":"async-1","is_running":false,"is_partial":false,"response":{"took":1,"timed_out":false,"hits":{"total":{"value":0,"relation":"eq"},"events":[],"sequences":[]}}}`
	case "/_eql/search/status/async-1":
		return `{"id":"async-1","is_running":false,"is_partial":false,"start_time_in_millis":0,"expiration_time_in_millis":0}`
	case "/products/_eql/search":
		return `{"is_partial":false,"is_running":false,"hits":{"events":[],"sequences":[],"total":{"value":0,"relation":"eq"}}}`
	case "/_query/async":
		return `{"id":"async-esql-1","is_running":false}`
	case "/_query/async/async-esql-1":
		return `{"id":"async-esql-1","is_running":false}`
	case "/_query/async/async-esql-1/_stop":
		return `{"id":"async-esql-1","is_running":false}`
	case "/products/_field_caps":
		return `{"indices":["products"],"fields":{}}`
	case "/_connector":
		return `{"results":[]}`
	case "/_connector/connector-1":
		return `{}`
	case "/_connector/connector-1/_check_in":
		return `{}`
	case "/_connector/connector-1/_last_sync":
		return `{}`
	case "/_connector/connector-1/_active_filtering":
		return `{}`
	case "/_connector/connector-1/_api_key_id":
		return `{}`
	case "/_connector/connector-1/_configuration":
		return `{}`
	case "/_connector/connector-1/_error":
		return `{}`
	case "/_connector/connector-1/_features":
		return `{}`
	case "/_connector/connector-1/_filtering":
		return `{}`
	case "/_connector/connector-1/_filtering_validation":
		return `{}`
	case "/_connector/connector-1/_index_name":
		return `{}`
	case "/_connector/connector-1/_name":
		return `{}`
	case "/_connector/connector-1/_native":
		return `{}`
	case "/_connector/connector-1/_pipeline":
		return `{}`
	case "/_connector/connector-1/_scheduling":
		return `{}`
	case "/_connector/connector-1/_service_type":
		return `{}`
	case "/_connector/connector-1/_status":
		return `{}`
	case "/_connector/_secret":
		return `{}`
	case "/_connector/_secret/secret-1":
		return `{}`
	case "/_connector/_sync_job":
		return `{"results":[]}`
	case "/_connector/_sync_job/job-1":
		return `{}`
	case "/_connector/_sync_job/job-1/_cancel":
		return `{}`
	case "/_connector/_sync_job/job-1/_check_in":
		return `{}`
	case "/_connector/_sync_job/job-1/_claim":
		return `{}`
	case "/_connector/_sync_job/job-1/_error":
		return `{}`
	case "/_connector/_sync_job/job-1/_stats":
		return `{}`
	case "/_features":
		return `{"features":[]}`
	case "/_features/_reset":
		return `{"features":[]}`
	case "/_data_stream/_stats":
		return `{"data_stream_count":0,"backing_indices":0,"total_store_size_bytes":0}`
	case "/products":
		return `{}`
	case "/_template/products-template":
		return `{}`
	case "/products/_lifecycle/explain":
		return `{"indices":{}}`
	case "/_data_stream/products/_lifecycle":
		return `{}`
	case "/_lifecycle/stats":
		return `{"data_retention":"disabled","default_rollover_used":false}
`
	case "/_mapping/field/status":
		return `{}`
	case "/_template":
		return `{}`
	case "/_recovery":
		return `{}`
	case "/_refresh":
		return `{"_shards":{"total":1,"successful":1,"failed":0}}`
	case "/products/_reload_search_analyzers":
		return `{"reload_details":[]}`
	case "/_resolve/cluster":
		return `{}`
	case "/_resolve/index/products":
		return `{"indices":[],"aliases":[],"data_streams":[]}`
	case "/_segments":
		return `{"indices":{}}`
	case "/_shard_stores":
		return `{"indices":{}}`
	case "/_validate/query":
		return `{"valid":true,"_shards":{"total":1,"successful":1,"failed":0}}`
	case "/_license":
		return `{"license":{"uid":"license-1","type":"basic","status":"active"}}`
	case "/_nodes":
		return `{"cluster_name":"test","nodes":{}}`
	case "/_nodes/stats":
		return `{"cluster_name":"test","nodes":{}}`
	case "/_query_rules":
		return `{"results":[]}`
	case "/_query_rules/ruleset-1":
		return `{}`
	case "/_query_rules/ruleset-1/_rule/rule-1":
		return `{}`
	case "/_query_rules/ruleset-1/_test":
		return `{"matched_rules":[]}`
	case "/_application/analytics":
		return `{"collection-1":{"event_data_stream":{"name":"logs-search_application.analytics-default"}}}`
	case "/_application/analytics/collection-1":
		return `{}`
	case "/_application/analytics/collection-1/event/click":
		return `{"accepted":true,"event":{"session":{"id":"session-1"},"user":{"id":"user-1"}}}`
	case "/_application/search_application":
		return `{"results":[]}`
	case "/_application/search_application/app-1":
		return `{}`
	case "/_application/search_application/app-1/_render_query":
		return `{"query":{}}`
	case "/_application/search_application/app-1/_search":
		return `{"took":1,"timed_out":false,"hits":{"total":{"value":0,"relation":"eq"},"hits":[]}}`
	case "/_snapshot":
		return `{}`
	case "/_snapshot/_status":
		return `{"snapshots":[]}`
	case "/_snapshot/repo-1/_analyze":
		return `{}`
	case "/_snapshot/repo-1/_cleanup":
		return `{"results":{"deleted_bytes":0,"deleted_blobs":0}}`
	case "/_snapshot/repo-1/_verify":
		return `{"nodes":{}}`
	case "/_snapshot/repo-1/_verify_integrity":
		return `{"result":"pass"}`
	case "/_snapshot/repo-1/snap-1/_clone/snap-2":
		return `{"accepted":true}`
	case "/_security/_authenticate":
		return `{"username":"user-1","roles":[],"authentication_realm":{"name":"realm","type":"file"},"lookup_realm":{"name":"realm","type":"file"},"authentication_type":"realm"}`
	case "/_security/_query/api_key":
		return `{"count":0,"api_keys":[]}`
	case "/_security/api_key/api-key-1":
		return `{}`
	case "/_security/api_key/grant":
		return `{}`
	case "/_security/cross_cluster/api_key":
		return `{}`
	case "/_security/delegate_pki":
		return `{}`
	case "/_security/oauth2/token":
		return `{}`
	case "/_security/profile/_activate":
		return `{}`
	case "/_security/role/role-1":
		return `{}`
	case "/_security/user/_password":
		return `{}`
	case "/_security/user/user-1":
		return `{}`
	case "/_sql":
		return `{"columns":[],"rows":[]}`
	case "/_sql/translate":
		return `{"size":0}`
	case "/_ssl/certificates":
		return `[]`
	case "/_synonyms":
		return `{"count":0,"results":[]}`
	case "/_tasks/node:1":
		return `{"completed":true,"task":{"node":"node","id":1,"type":"transport","action":"cluster:monitor/task/get","cancellable":false,"headers":{}}}`
	case "/_transform":
		return `{"count":0,"transforms":[]}`
	case "/_transform/_node_stats":
		return `true`
	case "/_transform/_preview":
		return `{"preview":[]}`
	case "/_transform/_set_upgrade_mode":
		return `{}`
	case "/_transform/_upgrade":
		return `{"accepted":true}`
	case "/_transform/transform-1/_reset":
		return `{"acknowledged":true}`
	case "/_transform/transform-1/_schedule_now":
		return `{"acknowledged":true}`
	case "/_transform/transform-1/_start":
		return `{"acknowledged":true}`
	case "/_transform/transform-1/_stop":
		return `{"acknowledged":true}`
	case "/_transform/transform-1/_update":
		return `{}`
	case "/_xpack":
		return `{}`
	case "/_xpack/usage":
		return `{}`
	default:
		return `{}`
	}
}

var _ = types.Query{}
