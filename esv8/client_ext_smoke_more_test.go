package esv8

import (
	"context"
	"io"
	"net/http"
	"testing"

	core_close_point_in_time "github.com/elastic/go-elasticsearch/v8/typedapi/core/closepointintime"
	core_msearch "github.com/elastic/go-elasticsearch/v8/typedapi/core/msearch"
	core_update_by_query "github.com/elastic/go-elasticsearch/v8/typedapi/core/updatebyquery"
	ilm_put_lifecycle "github.com/elastic/go-elasticsearch/v8/typedapi/ilm/putlifecycle"
	indices_analyze "github.com/elastic/go-elasticsearch/v8/typedapi/indices/analyze"
	indices_put_index_template "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putindextemplate"
	indices_put_mapping "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putmapping"
	indices_put_settings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putsettings"
	inference_inference "github.com/elastic/go-elasticsearch/v8/typedapi/inference/inference"
	ingest_put_pipeline "github.com/elastic/go-elasticsearch/v8/typedapi/ingest/putpipeline"
	ml_put_datafeed "github.com/elastic/go-elasticsearch/v8/typedapi/ml/putdatafeed"
	ml_put_job "github.com/elastic/go-elasticsearch/v8/typedapi/ml/putjob"
	transform_put_transform "github.com/elastic/go-elasticsearch/v8/typedapi/transform/puttransform"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

func TestClientExtSmoke_MoreRequestPathsAndMethods(t *testing.T) {
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
			name:       "MsearchMore",
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
			name:       "UpdateByQueryMore",
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
			name:       "PutMappingMore",
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
			name:       "PutSettingsMore",
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
			name:       "PutIndexTemplateMore",
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
			name:       "AnalyzeMore",
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
			name:       "ClosePointInTimeMore",
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
			name:       "PutLifecycleMore",
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
			name:       "PutPipelineMore",
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
			name:       "PutInferenceMore",
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
			name:       "InferenceMore",
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
			name:       "MlPutJobMore",
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
			name:       "MlPutDatafeedMore",
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
			name:       "CcrResumeFollowMore",
			wantMethod: http.MethodPost,
			wantPath:   "/follower-index/_ccr/resume_follow",
			call: func(t *testing.T, client *esClient) {
				t.Helper()

				_, err := client.CcrResumeFollow(
					context.Background(),
					estype.Index("follower-index"),
				)
				assert.NilError(t, err)
			},
		},
		{
			name:       "PutTransformMore",
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
