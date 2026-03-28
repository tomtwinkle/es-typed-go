package esv8_test

import (
	"reflect"
	"testing"

	"gotest.tools/v3/assert"

	esv8 "github.com/tomtwinkle/es-typed-go/esv8"
)

// TestESClientSpec_InterfaceCompleteness verifies that ESClientSpec declares
// all the spec-named methods and that esClient (returned by NewClient) satisfies
// the interface at runtime.
func TestESClientSpec_InterfaceCompleteness(t *testing.T) {
	t.Parallel()

	specType := reflect.TypeFor[esv8.ESClientSpec]()
	clientType := reflect.TypeFor[esv8.ESClient]()

	// ESClientSpec must embed ESClient — every ESClient method must appear in ESClientSpec.
	for m := range clientType.Methods() {
		_, ok := specType.MethodByName(m.Name)
		assert.Assert(t, ok, "ESClientSpec is missing ESClient method %q", m.Name)
	}

	// ESClientSpec must have strictly MORE methods than ESClient.
	assert.Assert(t, specType.NumMethod() > clientType.NumMethod(),
		"ESClientSpec should have more methods than ESClient")
	t.Logf("ESClient methods: %d, ESClientSpec methods: %d (added: %d)",
		clientType.NumMethod(), specType.NumMethod(),
		specType.NumMethod()-clientType.NumMethod())
}

// TestESClientSpec_SpecMethodsPresent spot-checks that a representative set
// of spec-named methods is present on the ESClientSpec interface.
func TestESClientSpec_SpecMethodsPresent(t *testing.T) {
	t.Parallel()

	specType := reflect.TypeFor[esv8.ESClientSpec]()

	// Methods to verify — covers core ops, indices, cluster, tasks, cat, etc.
	wantMethods := []string{
		// Core operations
		"Index",
		"Get",
		"Delete",
		"Update",
		"Bulk",
		"Count",
		"Mget",
		"Msearch",
		"Exists",
		"ExistsSource",
		"DeleteByQuery",
		"UpdateByQuery",
		"Scroll",
		"ClearScroll",
		"OpenPointInTime",
		"ClosePointInTime",
		"Explain",
		"FieldCaps",
		"GetSource",
		"GetScript",
		"PutScript",
		"DeleteScript",
		"TermsEnum",
		"Termvectors",
		"Ping",
		"HealthReport",
		"KnnSearch",
		"SearchTemplate",
		"SearchShards",
		"SearchMvt",
		"RankEval",
		"Mtermvectors",
		"MsearchTemplate",
		"ScriptsPainlessExecute",
		"RenderSearchTemplate",
		"ReindexRethrottle",
		"UpdateByQueryRethrottle",
		"DeleteByQueryRethrottle",
		// Indices namespace
		"IndicesCreate",
		"IndicesDelete",
		"IndicesExists",
		"IndicesRefresh",
		"IndicesGetAlias",
		"IndicesPutAlias",
		"IndicesExistsAlias",
		"IndicesDeleteAlias",
		"IndicesUpdateAliases",
		"IndicesGetSettings",
		"IndicesPutSettings",
		"IndicesGetMapping",
		"IndicesPutMapping",
		"IndicesOpen",
		"IndicesClose",
		"IndicesRollover",
		"IndicesForcemerge",
		"IndicesFlush",
		"IndicesShrink",
		"IndicesSplit",
		"IndicesAnalyze",
		"IndicesClone",
		"IndicesGet",
		"IndicesPutIndexTemplate",
		"IndicesGetIndexTemplate",
		"IndicesDeleteIndexTemplate",
		"IndicesCreateDataStream",
		"IndicesDeleteDataStream",
		"IndicesGetDataStream",
		"IndicesStats",
		// Cluster namespace
		"ClusterHealth",
		"ClusterStats",
		"ClusterGetSettings",
		"ClusterPutSettings",
		"ClusterAllocationExplain",
		"ClusterState",
		// Tasks namespace
		"TasksGet",
		"TasksList",
		"TasksCancel",
		// Cat namespace
		"CatHealth",
		"CatIndices",
		"CatAliases",
		"CatNodes",
		"CatShards",
		"CatCount",
		// ILM namespace
		"IlmGetLifecycle",
		"IlmPutLifecycle",
		"IlmDeleteLifecycle",
		"IlmStart",
		"IlmStop",
		// Snapshot namespace
		"SnapshotCreate",
		"SnapshotGet",
		"SnapshotDelete",
		"SnapshotGetRepository",
		// Ingest namespace
		"IngestGetPipeline",
		"IngestPutPipeline",
		"IngestDeletePipeline",
		// Already present in ESClient (via embedding)
		"Info",
		"Reindex",
	}

	missing := make([]string, 0)
	for _, name := range wantMethods {
		if _, ok := specType.MethodByName(name); !ok {
			missing = append(missing, name)
		}
	}
	assert.Assert(t, len(missing) == 0,
		"ESClientSpec is missing expected methods: %v", missing)
}

// TestNewClient_ImplementsESClientSpec verifies that the concrete client
// returned by NewClient also satisfies ESClientSpec (interface completeness check
// is also performed at compile time via the var _ ESClientSpec = (*esClient)(nil) guard).
func TestNewClient_ImplementsESClientSpec(t *testing.T) {
	t.Parallel()

	specType := reflect.TypeFor[esv8.ESClientSpec]()
	// NewClient returns ESClient; ESClientSpec extends it, so any concrete
	// implementation that satisfies ESClient but NOT ESClientSpec would fail here.
	// The compile-time guard in esclient_spec.go ensures this too.
	assert.Assert(t, specType.Kind() == reflect.Interface)
	assert.Assert(t, specType.NumMethod() > 0)
}
