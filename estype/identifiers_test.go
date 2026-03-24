package estype_test

import (
	"testing"

	"github.com/tomtwinkle/es-typed-go/estype"
	"gotest.tools/v3/assert"
)

func TestIdentifierStringMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "DocumentID",
			got:  estype.DocumentID("doc-1").String(),
			want: "doc-1",
		},
		{
			name: "DataStream",
			got:  estype.DataStream("logs-app").String(),
			want: "logs-app",
		},
		{
			name: "Policy",
			got:  estype.Policy("daily-retention").String(),
			want: "daily-retention",
		},
		{
			name: "Pipeline",
			got:  estype.Pipeline("ingest-main").String(),
			want: "ingest-main",
		},
		{
			name: "Template",
			got:  estype.Template("products-template").String(),
			want: "products-template",
		},
		{
			name: "Repository",
			got:  estype.Repository("snapshot-repo").String(),
			want: "snapshot-repo",
		},
		{
			name: "Snapshot",
			got:  estype.Snapshot("snapshot-20250101").String(),
			want: "snapshot-20250101",
		},
		{
			name: "TaskID",
			got:  estype.TaskID("node:123").String(),
			want: "node:123",
		},
		{
			name: "InferenceID",
			got:  estype.InferenceID("elser-endpoint").String(),
			want: "elser-endpoint",
		},
		{
			name: "MLJobID",
			got:  estype.MLJobID("job-1").String(),
			want: "job-1",
		},
		{
			name: "DatafeedID",
			got:  estype.DatafeedID("datafeed-1").String(),
			want: "datafeed-1",
		},
		{
			name: "TransformID",
			got:  estype.TransformID("transform-1").String(),
			want: "transform-1",
		},
		{
			name: "DataFrameAnalyticsID",
			got:  estype.DataFrameAnalyticsID("dfa-1").String(),
			want: "dfa-1",
		},
		{
			name: "TrainedModelID",
			got:  estype.TrainedModelID("model-1").String(),
			want: "model-1",
		},
		{
			name: "KeepAlive",
			got:  estype.KeepAlive("1m").String(),
			want: "1m",
		},
		{
			name: "ESQLQuery",
			got:  estype.ESQLQuery("FROM my-index | LIMIT 10").String(),
			want: "FROM my-index | LIMIT 10",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.got)
		})
	}
}
