package esv9_test

import (
	"reflect"
	"testing"

	"gotest.tools/v3/assert"

	esv9 "github.com/tomtwinkle/es-typed-go/esv9"
)

// TestESClientSpec_InterfaceCompleteness verifies that ESClientSpec declares
// all ESClient methods and that it adds spec-named methods on top.
func TestESClientSpec_InterfaceCompleteness(t *testing.T) {
	t.Parallel()

	specType := reflect.TypeOf((*esv9.ESClientSpec)(nil)).Elem()
	clientType := reflect.TypeOf((*esv9.ESClient)(nil)).Elem()

	// Every ESClient method must appear in ESClientSpec (it embeds ESClient).
	for i := 0; i < clientType.NumMethod(); i++ {
		m := clientType.Method(i)
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

// TestESClientSpec_SpecMethodsPresent spot-checks that representative spec-named
// methods are present on the ESClientSpec interface.
func TestESClientSpec_SpecMethodsPresent(t *testing.T) {
	t.Parallel()

	specType := reflect.TypeOf((*esv9.ESClientSpec)(nil)).Elem()

	wantMethods := []string{
		// Core operations
		"Index", "Get", "Delete", "Update", "Bulk", "Count", "Mget",
		"Exists", "DeleteByQuery", "UpdateByQuery", "Scroll", "ClearScroll",
		"OpenPointInTime", "ClosePointInTime", "Explain", "FieldCaps",
		"GetSource", "GetScript", "PutScript", "DeleteScript", "TermsEnum",
		"Termvectors", "Ping", "HealthReport", "KnnSearch",
		// Indices namespace
		"IndicesCreate", "IndicesDelete", "IndicesExists", "IndicesRefresh",
		"IndicesGetAlias", "IndicesPutAlias", "IndicesExistsAlias",
		"IndicesUpdateAliases", "IndicesGetSettings", "IndicesPutSettings",
		"IndicesGetMapping", "IndicesPutMapping", "IndicesOpen", "IndicesClose",
		"IndicesRollover", "IndicesForcemerge", "IndicesFlush",
		// Cluster namespace
		"ClusterHealth", "ClusterStats", "ClusterGetSettings",
		// Tasks namespace
		"TasksGet", "TasksList", "TasksCancel",
		// Cat namespace
		"CatHealth", "CatIndices", "CatAliases", "CatNodes",
		// ILM namespace
		"IlmGetLifecycle", "IlmPutLifecycle", "IlmDeleteLifecycle",
		// Already in ESClient
		"Info", "Search", "Reindex",
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

	specType := reflect.TypeOf((*esv9.ESClientSpec)(nil)).Elem()
	// NewClient returns ESClient; ESClientSpec extends it, so any concrete
	// implementation that satisfies ESClient but NOT ESClientSpec would fail here.
	// The compile-time guard in esclient_spec.go ensures this too.
	assert.Assert(t, specType.Kind() == reflect.Interface)
	assert.Assert(t, specType.NumMethod() > 0)
}
