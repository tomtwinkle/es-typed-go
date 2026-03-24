package esv9

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	autoscaling_put_autoscaling_policy "github.com/elastic/go-elasticsearch/v9/typedapi/autoscaling/putautoscalingpolicy"
	connector_post "github.com/elastic/go-elasticsearch/v9/typedapi/connector/post"
	connector_put "github.com/elastic/go-elasticsearch/v9/typedapi/connector/put"
	eql_search "github.com/elastic/go-elasticsearch/v9/typedapi/eql/search"
	esql_async_query "github.com/elastic/go-elasticsearch/v9/typedapi/esql/asyncquery"
	inference_inference "github.com/elastic/go-elasticsearch/v9/typedapi/inference/inference"
	query_rules_put_ruleset "github.com/elastic/go-elasticsearch/v9/typedapi/queryrules/putruleset"
	search_application_render_query "github.com/elastic/go-elasticsearch/v9/typedapi/searchapplication/renderquery"
	search_application_search "github.com/elastic/go-elasticsearch/v9/typedapi/searchapplication/search"
	security_create_api_key "github.com/elastic/go-elasticsearch/v9/typedapi/security/createapikey"
	security_get_api_key "github.com/elastic/go-elasticsearch/v9/typedapi/security/getapikey"
	snapshot_create_repository "github.com/elastic/go-elasticsearch/v9/typedapi/snapshot/createrepository"
	transform_put_transform "github.com/elastic/go-elasticsearch/v9/typedapi/transform/puttransform"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

type recordedRequest struct {
	method string
	path   string
	query  string
	body   string
}

type smokeTransport struct {
	t       *testing.T
	handler func(*http.Request) (status int, body string)
	last    *recordedRequest
}

func (tr *smokeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tr.t.Helper()

	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	tr.last = &recordedRequest{
		method: req.Method,
		path:   req.URL.Path,
		query:  req.URL.RawQuery,
		body:   string(bodyBytes),
	}

	status := http.StatusOK
	body := `{}`

	if tr.handler != nil {
		status, body = tr.handler(req)
	}

	res := &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}
	res.Header.Set("X-Elastic-Product", "Elasticsearch")

	return res, nil
}

func newSmokeClient(t *testing.T, handler func(*http.Request) (int, string)) (*esClient, *smokeTransport) {
	t.Helper()

	tr := &smokeTransport{
		t:       t,
		handler: handler,
	}
	typed := newTestTypedClient(t, tr, "http://example.test")

	return &esClient{
		typedClient: typed,
	}, tr
}

func bodyContainsAll(t *testing.T, body string, want ...string) {
	t.Helper()

	for _, w := range want {
		assert.Assert(t, strings.Contains(body, w))
	}
}

func TestClientExtSmokeStableSubset(t *testing.T) {
	t.Parallel()

	t.Run("EsqlQuery", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"columns":[],"values":[]}`
		})

		_, err := client.EsqlQuery(context.Background(), estype.ESQLQuery("from logs | limit 1"))
		assert.NilError(t, err)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_query", tr.last.path)
		bodyContainsAll(t, tr.last.body, `"query":"from logs | limit 1"`)
	})

	t.Run("GetApiKey", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"api_keys":[]}`
		})

		got, err := client.GetApiKey(
			context.Background(),
			func(r *security_get_api_key.GetApiKey) {
				r.Id("api-id")
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_security/api_key", tr.last.path)
		assert.Assert(t, strings.Contains(tr.last.query, "id=api-id"))
	})

	t.Run("CreateRepository", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.CreateRepository(
			context.Background(),
			estype.Repository("repo-one"),
			func(r *snapshot_create_repository.CreateRepository) {
				r.Raw(strings.NewReader(`{"type":"fs"}`))
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodPut, tr.last.method)
		assert.Equal(t, "/_snapshot/repo-one", tr.last.path)
		bodyContainsAll(t, tr.last.body, `"type":"fs"`)
	})

	t.Run("GetInference", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"endpoints":[]}`
		})

		got, err := client.GetInference(context.Background(), estype.InferenceID("model-1"))
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_inference/model-1", tr.last.path)
	})

	t.Run("DeleteInference", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.DeleteInference(context.Background(), estype.InferenceID("model-1"))
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_inference/model-1", tr.last.path)
	})

	t.Run("Inference", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"inference_results":[]}`
		})

		got, err := client.Inference(
			context.Background(),
			estype.InferenceID("model-1"),
			func(r *inference_inference.Inference) {
				r.Raw(strings.NewReader(`{"input":["hello"]}`))
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_inference/model-1", tr.last.path)
		bodyContainsAll(t, tr.last.body, `"input":["hello"]`)
	})

	t.Run("PutTransform", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.PutTransform(
			context.Background(),
			estype.TransformID("transform-1"),
			func(r *transform_put_transform.PutTransform) {
				r.Raw(strings.NewReader(`{}`))
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodPut, tr.last.method)
		assert.Equal(t, "/_transform/transform-1", tr.last.path)
	})

	t.Run("GetTransform", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"count":0,"transforms":[]}`
		})

		got, err := client.GetTransform(context.Background(), estype.TransformID("transform-1"))
		assert.NilError(t, err)
		assert.Equal(t, int64(0), got.Count)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_transform/transform-1", tr.last.path)
	})

	t.Run("DeleteTransform", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.DeleteTransform(context.Background(), estype.TransformID("transform-1"))
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_transform/transform-1", tr.last.path)
	})

	t.Run("StartTransform", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.StartTransform(context.Background(), estype.TransformID("transform-1"))
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_transform/transform-1/_start", tr.last.path)
	})

	t.Run("StopTransform", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.StopTransform(context.Background(), estype.TransformID("transform-1"))
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_transform/transform-1/_stop", tr.last.path)
	})

	t.Run("GetTransformStats", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"count":0,"transforms":[]}`
		})

		got, err := client.GetTransformStats(context.Background(), estype.TransformID("transform-1"))
		assert.NilError(t, err)
		assert.Equal(t, int64(0), got.Count)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_transform/transform-1/_stats", tr.last.path)
	})

	t.Run("CatAliases", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatAliases(context.Background(), estype.Alias("alias-one"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/aliases/alias-one", tr.last.path)
	})

	t.Run("CatAllocation", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatAllocation(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/allocation", tr.last.path)
	})

	t.Run("CatCircuitBreaker", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatCircuitBreaker(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/circuit_breaker", tr.last.path)
	})

	t.Run("CatComponentTemplates", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatComponentTemplates(context.Background(), estype.Template("tpl-one"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/component_templates/tpl-one", tr.last.path)
	})

	t.Run("CatFielddata", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatFielddata(context.Background(), []estype.Field{"status", "title"})
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/fielddata/status,title", tr.last.path)
	})

	t.Run("CatHealth", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatHealth(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/health", tr.last.path)
	})

	t.Run("CatIndices", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatIndices(context.Background(), estype.Index("logs-0001"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/indices/logs-0001", tr.last.path)
	})

	t.Run("CatMaster", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatMaster(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/master", tr.last.path)
	})

	t.Run("CatMlDataFrameAnalytics", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatMlDataFrameAnalytics(context.Background(), estype.DataFrameAnalyticsID("analytics-1"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/ml/data_frame/analytics/analytics-1", tr.last.path)
	})

	t.Run("CatMlDatafeeds", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatMlDatafeeds(context.Background(), estype.DatafeedID("datafeed-1"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/ml/datafeeds/datafeed-1", tr.last.path)
	})

	t.Run("CatMlJobs", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatMlJobs(context.Background(), estype.MLJobID("job-1"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/ml/anomaly_detectors/job-1", tr.last.path)
	})

	t.Run("CatMlTrainedModels", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatMlTrainedModels(context.Background(), estype.TrainedModelID("model-1"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/ml/trained_models/model-1", tr.last.path)
	})

	t.Run("CatNodeattrs", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatNodeattrs(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/nodeattrs", tr.last.path)
	})

	t.Run("CatNodes", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatNodes(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/nodes", tr.last.path)
	})

	t.Run("CatPendingTasks", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatPendingTasks(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/pending_tasks", tr.last.path)
	})

	t.Run("CatPlugins", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatPlugins(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/plugins", tr.last.path)
	})

	t.Run("CatRecovery", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatRecovery(context.Background(), estype.Index("logs-0001"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/recovery/logs-0001", tr.last.path)
	})

	t.Run("CatRepositories", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatRepositories(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/repositories", tr.last.path)
	})

	t.Run("CatSegments", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatSegments(context.Background(), estype.Index("logs-0001"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/segments/logs-0001", tr.last.path)
	})

	t.Run("CatShards", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatShards(context.Background(), estype.Index("logs-0001"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/shards/logs-0001", tr.last.path)
	})

	t.Run("CatTasks", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatTasks(context.Background())
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/tasks", tr.last.path)
	})

	t.Run("CatTemplates", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatTemplates(context.Background(), estype.Template("tpl-one"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/templates/tpl-one", tr.last.path)
	})

	t.Run("CatTransforms", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `[]`
		})

		got, err := client.CatTransforms(context.Background(), estype.TransformID("transform-1"))
		assert.NilError(t, err)
		assert.Equal(t, 0, len(got))

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_cat/transforms/transform-1", tr.last.path)
	})
}

func TestClientSpecSmokeStableSubset(t *testing.T) {
	t.Parallel()

	t.Run("AutoscalingPutAutoscalingPolicy", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.AutoscalingPutAutoscalingPolicy(context.Background(), "policy-one", &autoscaling_put_autoscaling_policy.Request{})
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodPut, tr.last.method)
		assert.Equal(t, "/_autoscaling/policy/policy-one", tr.last.path)
	})

	t.Run("AutoscalingDeleteAutoscalingPolicy", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.AutoscalingDeleteAutoscalingPolicy(context.Background(), "policy-one")
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_autoscaling/policy/policy-one", tr.last.path)
	})

	t.Run("ConnectorCheckIn", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"result":"updated"}`
		})

		got, err := client.ConnectorCheckIn(context.Background(), "connector-1")
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPut, tr.last.method)
		assert.Equal(t, "/_connector/connector-1/_check_in", tr.last.path)
	})

	t.Run("ConnectorDelete", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"result":"deleted"}`
		})

		got, err := client.ConnectorDelete(context.Background(), "connector-1")
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_connector/connector-1", tr.last.path)
	})

	t.Run("ConnectorGet", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"id":"connector-1"}`
		})

		got, err := client.ConnectorGet(context.Background(), "connector-1")
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_connector/connector-1", tr.last.path)
	})

	t.Run("ConnectorList", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"results":[]}`
		})

		got, err := client.ConnectorList(context.Background())
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_connector", tr.last.path)
	})

	t.Run("ConnectorPost", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"id":"connector-1"}`
		})

		got, err := client.ConnectorPost(context.Background(), &connector_post.Request{})
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_connector", tr.last.path)
	})

	t.Run("ConnectorPut", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"id":"connector-1"}`
		})

		got, err := client.ConnectorPut(context.Background(), &connector_put.Request{})
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPut, tr.last.method)
		assert.Equal(t, "/_connector", tr.last.path)
	})

	t.Run("EqlDelete", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"is_running":false}`
		})

		got, err := client.EqlDelete(context.Background(), "async-id")
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_eql/search/async-id", tr.last.path)
	})

	t.Run("EqlGet", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"is_running":false,"is_partial":false,"hits":{"events":[]}}`
		})

		got, err := client.EqlGet(context.Background(), "async-id")
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_eql/search/async-id", tr.last.path)
	})

	t.Run("EqlGetStatus", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"id":"async-id","is_running":false,"is_partial":false}`
		})

		got, err := client.EqlGetStatus(context.Background(), "async-id")
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_eql/search/status/async-id", tr.last.path)
	})

	t.Run("EqlSearch", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"is_partial":false,"hits":{"events":[]}}`
		})

		got, err := client.EqlSearch(context.Background(), "logs-*", &eql_search.Request{})
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/logs-*/_eql/search", tr.last.path)
	})

	t.Run("EsqlAsyncQuery", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"id":"esql-1","is_running":false,"is_partial":false,"columns":[],"values":[]}`
		})

		_, err := client.EsqlAsyncQuery(context.Background(), &esql_async_query.Request{})
		assert.NilError(t, err)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_query/async", tr.last.path)
	})

	t.Run("EsqlAsyncQueryDelete", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"deleted":true}`
		})

		_, err := client.EsqlAsyncQueryDelete(context.Background(), "esql-1")
		assert.NilError(t, err)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_query/async/esql-1", tr.last.path)
	})

	t.Run("EsqlAsyncQueryGet", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"id":"esql-1","is_running":false,"is_partial":false,"columns":[],"values":[]}`
		})

		_, err := client.EsqlAsyncQueryGet(context.Background(), "esql-1")
		assert.NilError(t, err)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_query/async/esql-1", tr.last.path)
	})

	t.Run("EsqlAsyncQueryStop", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"id":"esql-1","is_running":false,"is_partial":false,"columns":[],"values":[]}`
		})

		_, err := client.EsqlAsyncQueryStop(context.Background(), "esql-1")
		assert.NilError(t, err)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_query/async/esql-1/stop", tr.last.path)
	})

	t.Run("InferenceGet", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"endpoints":[]}`
		})

		got, err := client.InferenceGet(context.Background())
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_inference", tr.last.path)
	})

	t.Run("InferenceDelete", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.InferenceDelete(context.Background(), "model-1")
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_inference/model-1", tr.last.path)
	})

	t.Run("InferenceInference", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"inference_results":[]}`
		})

		got, err := client.InferenceInference(context.Background(), "model-1", &inference_inference.Request{})
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_inference/model-1", tr.last.path)
	})

	t.Run("InferencePut", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"inference_id":"model-1"}`
		})

		got, err := client.InferencePut(context.Background(), "model-1", nil)
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPut, tr.last.method)
		assert.Equal(t, "/_inference/model-1", tr.last.path)
	})

	t.Run("QueryRulesDeleteRuleset", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.QueryRulesDeleteRuleset(context.Background(), "ruleset-1")
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_query_rules/ruleset-1", tr.last.path)
	})

	t.Run("QueryRulesGetRuleset", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"ruleset_id":"ruleset-1","rules":[]}`
		})

		got, err := client.QueryRulesGetRuleset(context.Background(), "ruleset-1")
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_query_rules/ruleset-1", tr.last.path)
	})

	t.Run("QueryRulesListRulesets", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"results":[]}`
		})

		got, err := client.QueryRulesListRulesets(context.Background())
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_query_rules", tr.last.path)
	})

	t.Run("QueryRulesPutRuleset", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"result":"created"}`
		})

		got, err := client.QueryRulesPutRuleset(context.Background(), "ruleset-1", &query_rules_put_ruleset.Request{})
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPut, tr.last.method)
		assert.Equal(t, "/_query_rules/ruleset-1", tr.last.path)
	})

	t.Run("SearchApplicationDelete", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"acknowledged":true}`
		})

		got, err := client.SearchApplicationDelete(context.Background(), "app-one")
		assert.NilError(t, err)
		assert.Assert(t, got.Acknowledged)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_application/search_application/app-one", tr.last.path)
	})

	t.Run("SearchApplicationGet", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"name":"app-one"}`
		})

		got, err := client.SearchApplicationGet(context.Background(), "app-one")
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_application/search_application/app-one", tr.last.path)
	})

	t.Run("SearchApplicationList", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"results":[]}`
		})

		got, err := client.SearchApplicationList(context.Background())
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_application/search_application", tr.last.path)
	})

	t.Run("SearchApplicationRenderQuery", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"query":{}}`
		})

		got, err := client.SearchApplicationRenderQuery(context.Background(), "app-one", &search_application_render_query.Request{})
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_application/search_application/app-one/_render_query", tr.last.path)
	})

	t.Run("SearchApplicationSearch", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"took":1,"timed_out":false,"hits":{"total":{"value":0,"relation":"eq"},"hits":[]}}`
		})

		got, err := client.SearchApplicationSearch(context.Background(), "app-one", &search_application_search.Request{})
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_application/search_application/app-one/_search", tr.last.path)
	})

	t.Run("SecurityAuthenticate", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"username":"elastic","roles":[],"enabled":true,"authentication_realm":{"name":"native","type":"native"},"lookup_realm":{"name":"native","type":"native"}}`
		})

		got, err := client.SecurityAuthenticate(context.Background())
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_security/_authenticate", tr.last.path)
	})

	t.Run("SecurityClearApiKeyCache", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{}`
		})

		got, err := client.SecurityClearApiKeyCache(context.Background(), "api-id")
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodPost, tr.last.method)
		assert.Equal(t, "/_security/api_key/api-id/_clear_cache", tr.last.path)
	})

	t.Run("SecurityCreateApiKey", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"id":"api-id","name":"key-name","api_key":"secret","encoded":"ZW5jb2RlZA=="}`
		})

		got, err := client.SecurityCreateApiKey(context.Background(), &security_create_api_key.Request{})
		assert.NilError(t, err)
		assert.Equal(t, "api-id", got.Id)

		assert.Equal(t, http.MethodPut, tr.last.method)
		assert.Equal(t, "/_security/api_key", tr.last.path)
	})

	t.Run("SecurityDeleteRole", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"found":true}`
		})

		got, err := client.SecurityDeleteRole(context.Background(), "role-one")
		assert.NilError(t, err)
		assert.Assert(t, got.Found)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_security/role/role-one", tr.last.path)
	})

	t.Run("SecurityDeleteRoleMapping", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"found":true}`
		})

		got, err := client.SecurityDeleteRoleMapping(context.Background(), "mapping-one")
		assert.NilError(t, err)
		assert.Assert(t, got.Found)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_security/role_mapping/mapping-one", tr.last.path)
	})

	t.Run("SecurityDeleteServiceToken", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"found":true}`
		})

		got, err := client.SecurityDeleteServiceToken(context.Background(), "ns-one", "svc-one", "token-one")
		assert.NilError(t, err)
		assert.Assert(t, got.Found)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_security/service/ns-one/svc-one/credential/token/token-one", tr.last.path)
	})

	t.Run("SecurityDeleteUser", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"found":true}`
		})

		got, err := client.SecurityDeleteUser(context.Background(), "user-one")
		assert.NilError(t, err)
		assert.Assert(t, got.Found)

		assert.Equal(t, http.MethodDelete, tr.last.method)
		assert.Equal(t, "/_security/user/user-one", tr.last.path)
	})

	t.Run("SecurityEnrollNode", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"http_ca_key":"key","http_ca_cert":"cert","transport_ca_cert":"cert","transport_key":"key","transport_cert":"cert","nodes_addresses":["127.0.0.1:9200"]}`
		})

		got, err := client.SecurityEnrollNode(context.Background())
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_security/enroll/node", tr.last.path)
	})

	t.Run("SecurityGetApiKey", func(t *testing.T) {
		t.Parallel()

		client, tr := newSmokeClient(t, func(req *http.Request) (int, string) {
			return http.StatusOK, `{"api_keys":[]}`
		})

		got, err := client.SecurityGetApiKey(context.Background())
		assert.NilError(t, err)
		assert.Assert(t, got != nil)

		assert.Equal(t, http.MethodGet, tr.last.method)
		assert.Equal(t, "/_security/api_key", tr.last.path)
	})
}
