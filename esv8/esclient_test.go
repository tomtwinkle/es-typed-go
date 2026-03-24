package esv8

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

func Test_taskIDToString(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		taskID  any
		want    string
		wantErr bool
	}{
		"string":              {taskID: "abc:123", want: "abc:123"},
		"int":                 {taskID: int(42), want: "42"},
		"int8":                {taskID: int8(8), want: "8"},
		"int16":               {taskID: int16(16), want: "16"},
		"int32":               {taskID: int32(32), want: "32"},
		"int64":               {taskID: int64(64), want: "64"},
		"uint":                {taskID: uint(1), want: "1"},
		"uint8":               {taskID: uint8(2), want: "2"},
		"uint16":              {taskID: uint16(3), want: "3"},
		"uint32":              {taskID: uint32(4), want: "4"},
		"uint64":              {taskID: uint64(5), want: "5"},
		"unsupported float64": {taskID: float64(1.5), wantErr: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := taskIDToString(tt.taskID)
			if tt.wantErr {
				assert.Assert(t, err != nil)
				return
			}
			assert.NilError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_isElasticsearchError_nil(t *testing.T) {
	t.Parallel()

	assert.Assert(t, !isElasticsearchError(nil, nil))
}

func Test_isElasticsearchError_direct(t *testing.T) {
	t.Parallel()

	reason := "bad query"
	esErr := &types.ElasticsearchError{
		Status: 400,
		ErrorCause: types.ErrorCause{
			Type:   "query_shard_exception",
			Reason: &reason,
		},
	}

	var target *types.ElasticsearchError
	ok := isElasticsearchError(esErr, &target)

	assert.Assert(t, ok)
	assert.Assert(t, target == esErr)
	assert.Equal(t, 400, target.Status)
	assert.Assert(t, target.ErrorCause.Reason != nil)
	assert.Equal(t, reason, *target.ErrorCause.Reason)
}

func Test_isElasticsearchError_wrapped(t *testing.T) {
	t.Parallel()

	reason := "index not found"
	esErr := &types.ElasticsearchError{
		Status: 404,
		ErrorCause: types.ErrorCause{
			Type:   "index_not_found_exception",
			Reason: &reason,
		},
	}
	err := fmtWrap(esErr)

	var target *types.ElasticsearchError
	ok := isElasticsearchError(err, &target)

	assert.Assert(t, ok)
	assert.Assert(t, target == esErr)
}

func Test_isElasticsearchError_nonElasticsearch(t *testing.T) {
	t.Parallel()

	var target *types.ElasticsearchError
	ok := isElasticsearchError(errors.New("plain error"), &target)

	assert.Assert(t, !ok)
	assert.Assert(t, target == nil)
}

func Test_unwrapErr(t *testing.T) {
	t.Parallel()

	inner := errors.New("inner")
	wrapped := fmtWrap(inner)

	got := unwrapErr(wrapped)

	assert.Assert(t, got == inner)
}

func Test_unwrapErr_withoutUnwrap(t *testing.T) {
	t.Parallel()

	got := unwrapErr(errors.New("plain"))

	assert.Assert(t, got == nil)
}

func Test_newESClient_usesDefaultLogger(t *testing.T) {
	t.Parallel()

	client := newTestESClient(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newJSONResponse(http.StatusOK, `{"version":{"number":"8.19.3"},"tagline":"You Know, for Search"}`), nil
	}))

	assert.Assert(t, client != nil)
	assert.Assert(t, client.logger != nil)
	assert.Assert(t, client.typedClient != nil)
}

func Test_schemeFromTransport_defaultsToHTTPWithoutURLProvider(t *testing.T) {
	t.Parallel()

	client := newTestESClient(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newJSONResponse(http.StatusOK, `{}`), nil
	}))

	assert.Equal(t, "http", client.schemeFromTransport())
}

func Test_schemeFromTransport_usesHTTPWithCustomRoundTripper(t *testing.T) {
	t.Parallel()

	client := newTestESClientWithURLs(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newJSONResponse(http.StatusOK, `{}`), nil
	}), []*url.URL{
		{Scheme: "https", Host: "example.test"},
	})

	assert.Equal(t, "http", client.schemeFromTransport())
}

func Test_schemeFromTransport_defaultsToHTTPWhenNoURLs(t *testing.T) {
	t.Parallel()

	client := newTestESClientWithURLs(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newJSONResponse(http.StatusOK, `{}`), nil
	}), nil)

	assert.Equal(t, "http", client.schemeFromTransport())
}

func Test_performRaw_successWithoutBody(t *testing.T) {
	t.Parallel()

	var seenMethod string
	var seenURL *url.URL
	var seenAccept string
	var seenContentType string
	var seenBody []byte

	client := newTestESClientWithURLs(t, testRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		seenMethod = req.Method
		seenURL = req.URL
		seenAccept = req.Header.Get("Accept")
		seenContentType = req.Header.Get("Content-Type")
		if req.Body != nil {
			body, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			seenBody = body
		}

		return newJSONResponse(http.StatusOK, `{"ok":true}`), nil
	}), []*url.URL{
		{Scheme: "https", Host: "cluster.test"},
	})

	got, err := client.performRaw(context.Background(), http.MethodGet, "/_internal/test", nil)

	assert.NilError(t, err)
	assert.Equal(t, `{"ok":true}`, string(got))
	assert.Equal(t, http.MethodGet, seenMethod)
	assert.Assert(t, seenURL != nil)
	assert.Equal(t, "http", seenURL.Scheme)
	assert.Equal(t, "/_internal/test", seenURL.Path)
	assert.Equal(t, "application/vnd.elasticsearch+json;compatible-with=8", seenAccept)
	assert.Equal(t, "", seenContentType)
	assert.Equal(t, 0, len(seenBody))
}

func Test_performRaw_successWithBodySetsHeaders(t *testing.T) {
	t.Parallel()

	var seenMethod string
	var seenURL *url.URL
	var seenAccept string
	var seenContentType string
	var seenBody []byte

	client := newTestESClientWithURLs(t, testRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		seenMethod = req.Method
		seenURL = req.URL
		seenAccept = req.Header.Get("Accept")
		seenContentType = req.Header.Get("Content-Type")
		body, err := io.ReadAll(req.Body)
		assert.NilError(t, err)
		seenBody = body

		return newJSONResponse(http.StatusOK, `{"acknowledged":true}`), nil
	}), []*url.URL{
		{Scheme: "http", Host: "cluster.test"},
	})

	body := json.RawMessage(`{"query":{"match_all":{}}}`)
	got, err := client.performRaw(context.Background(), http.MethodPost, "/_search", body)

	assert.NilError(t, err)
	assert.Equal(t, `{"acknowledged":true}`, string(got))
	assert.Equal(t, http.MethodPost, seenMethod)
	assert.Assert(t, seenURL != nil)
	assert.Equal(t, "http", seenURL.Scheme)
	assert.Equal(t, "/_search", seenURL.Path)
	assert.Equal(t, "application/vnd.elasticsearch+json;compatible-with=8", seenAccept)
	assert.Equal(t, "application/vnd.elasticsearch+json;compatible-with=8", seenContentType)
	assert.Equal(t, string(body), string(seenBody))
}

func Test_performRaw_performError(t *testing.T) {
	t.Parallel()

	client := newTestESClient(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("network down")
	}))

	_, err := client.performRaw(context.Background(), http.MethodGet, "/_cluster/health", nil)

	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "performing GET /_cluster/health"))
	assert.Assert(t, strings.Contains(err.Error(), "network down"))
}

func Test_performRaw_httpErrorIncludesBody(t *testing.T) {
	t.Parallel()

	client := newTestESClient(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newJSONResponse(http.StatusBadRequest, `{"error":"bad request"}`), nil
	}))

	_, err := client.performRaw(context.Background(), http.MethodGet, "/_bad", nil)

	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "unexpected status 400"))
	assert.Assert(t, strings.Contains(err.Error(), `{"error":"bad request"}`))
}

func Test_performRaw_httpErrorTruncatesLongBody(t *testing.T) {
	t.Parallel()

	longBody := strings.Repeat("x", 600)
	client := newTestESClient(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newTextResponse(http.StatusInternalServerError, longBody), nil
	}))

	_, err := client.performRaw(context.Background(), http.MethodGet, "/_fail", nil)

	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "unexpected status 500"))
	assert.Assert(t, strings.Contains(err.Error(), strings.Repeat("x", 512)))
	assert.Assert(t, !strings.Contains(err.Error(), strings.Repeat("x", 513)))
}

func Test_buildSearchError_nil(t *testing.T) {
	t.Parallel()

	client := newTestESClient(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newJSONResponse(http.StatusOK, `{}`), nil
	}))

	assert.Assert(t, client.buildSearchError(nil) == nil)
}

func Test_buildSearchError_nonElasticsearchError(t *testing.T) {
	t.Parallel()

	client := newTestESClient(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newJSONResponse(http.StatusOK, `{}`), nil
	}))

	err := errors.New("plain error")
	got := client.buildSearchError(err)

	assert.Assert(t, got == err)
}

func Test_buildSearchError_elasticsearchErrorWithoutRootCauseReason(t *testing.T) {
	t.Parallel()

	reason := "all shards failed"
	esErr := &types.ElasticsearchError{
		Status: 400,
		ErrorCause: types.ErrorCause{
			Type:   "search_phase_execution_exception",
			Reason: &reason,
			RootCause: []types.ErrorCause{
				{Type: "search_phase_execution_exception"},
			},
		},
	}

	client := newTestESClient(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newJSONResponse(http.StatusOK, `{}`), nil
	}))

	got := client.buildSearchError(esErr)

	assert.Assert(t, got != nil)
	assert.Equal(t, esErr.Error(), got.Error())
}

func Test_buildSearchError_elasticsearchErrorWithRootCauseReasons(t *testing.T) {
	t.Parallel()

	reason := "all shards failed"
	rootReason1 := "failed to create query"
	rootReason2 := "field [status] not found"

	esErr := &types.ElasticsearchError{
		Status: 400,
		ErrorCause: types.ErrorCause{
			Type:   "search_phase_execution_exception",
			Reason: &reason,
			RootCause: []types.ErrorCause{
				{Type: "query_shard_exception", Reason: &rootReason1},
				{Type: "illegal_argument_exception", Reason: &rootReason2},
			},
		},
	}

	client := newTestESClient(t, testRoundTripperFunc(func(*http.Request) (*http.Response, error) {
		return newJSONResponse(http.StatusOK, `{}`), nil
	}))

	got := client.buildSearchError(esErr)

	assert.Assert(t, got != nil)
	assert.Assert(t, strings.Contains(got.Error(), reason))
	assert.Assert(t, strings.Contains(got.Error(), rootReason1))
	assert.Assert(t, strings.Contains(got.Error(), rootReason2))
}

func Test_CreateIndexFromDefinitions_callsCreateIndex(t *testing.T) {
	t.Parallel()

	var seenMethod string
	var seenPath string
	var payload map[string]any

	client := newTestESClient(t, testRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		seenMethod = req.Method
		seenPath = req.URL.Path

		body, err := io.ReadAll(req.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(body, &payload)
		assert.NilError(t, err)

		return newJSONResponse(http.StatusOK, `{"acknowledged":true,"shards_acknowledged":true,"index":"products-000001"}`), nil
	}))

	settings := estype.Settings{
		NumberOfShards:   intPtr(3),
		NumberOfReplicas: intPtr(1),
		RefreshInterval:  refreshIntervalPtr(estype.RefreshInterval(estype.RefreshIntervalDefault)),
	}
	mapping := estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "status", Property: estype.NewKeywordProperty()},
			{Path: "title", Property: estype.NewTextProperty()},
		},
	}

	res, err := client.CreateIndexFromDefinitions(context.Background(), estype.Index("products-000001"), settings, mapping)

	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	assert.Equal(t, http.MethodPut, seenMethod)
	assert.Equal(t, "/products-000001", seenPath)

	settingsMap, ok := payload["settings"].(map[string]any)
	assert.Assert(t, ok)
	assert.Equal(t, "3", settingsMap["number_of_shards"])
	assert.Equal(t, "1", settingsMap["number_of_replicas"])
	assert.Equal(t, "1s", settingsMap["refresh_interval"])

	mappingsMap, ok := payload["mappings"].(map[string]any)
	assert.Assert(t, ok)
	propertiesMap, ok := mappingsMap["properties"].(map[string]any)
	assert.Assert(t, ok)
	_, ok = propertiesMap["status"]
	assert.Assert(t, ok)
	_, ok = propertiesMap["title"]
	assert.Assert(t, ok)
}

func Test_CreateIndexFromProviders_callsCreateIndex(t *testing.T) {
	t.Parallel()

	var seenPath string
	var payload map[string]any

	client := newTestESClient(t, testRoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		seenPath = req.URL.Path

		body, err := io.ReadAll(req.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(body, &payload)
		assert.NilError(t, err)

		return newJSONResponse(http.StatusOK, `{"acknowledged":true,"shards_acknowledged":true,"index":"products-000002"}`), nil
	}))

	res, err := client.CreateIndexFromProviders(context.Background(), estype.Index("products-000002"), createIndexModelStub{})

	assert.NilError(t, err)
	assert.Assert(t, res != nil)
	assert.Equal(t, "/products-000002", seenPath)

	settingsMap, ok := payload["settings"].(map[string]any)
	assert.Assert(t, ok)
	assert.Equal(t, "3", settingsMap["number_of_shards"])
	assert.Equal(t, "1", settingsMap["number_of_replicas"])

	mappingsMap, ok := payload["mappings"].(map[string]any)
	assert.Assert(t, ok)
	propertiesMap, ok := mappingsMap["properties"].(map[string]any)
	assert.Assert(t, ok)

	itemsValue, ok := propertiesMap["items"]
	assert.Assert(t, ok)
	itemsMap, ok := itemsValue.(map[string]any)
	assert.Assert(t, ok)
	nestedProps, ok := itemsMap["properties"].(map[string]any)
	assert.Assert(t, ok)
	_, ok = nestedProps["name"]
	assert.Assert(t, ok)
	_, ok = nestedProps["value"]
	assert.Assert(t, ok)
}

func Test_toTypedTypeMapping_skipsInvalidFields(t *testing.T) {
	t.Parallel()

	got := toTypedTypeMapping(estype.Mapping{
		Fields: []estype.MappingField{
			{Path: "", Property: estype.NewKeywordProperty()},
			{Path: "status", Property: nil},
		},
	})

	assert.Assert(t, got == nil)
}

func Test_insertTypedProperty_ignoresEmptyInputs(t *testing.T) {
	t.Parallel()

	props := map[string]types.Property{
		"status": types.KeywordProperty{Type: "keyword"},
	}

	insertTypedProperty(props, nil, estype.NewTextProperty())
	insertTypedProperty(props, []string{""}, estype.NewTextProperty())
	insertTypedProperty(props, []string{"title"}, nil)

	assert.Equal(t, 1, len(props))
	_, ok := props["status"]
	assert.Assert(t, ok)
}

func Test_insertTypedProperty_existingObjectProperty(t *testing.T) {
	t.Parallel()

	props := map[string]types.Property{
		"items": types.ObjectProperty{
			Type:       "object",
			Properties: map[string]types.Property{},
		},
	}

	insertTypedProperty(props, []string{"items", "name"}, estype.NewTextProperty())

	itemsProp, ok := props["items"]
	assert.Assert(t, ok)

	objectProp, ok := itemsProp.(types.ObjectProperty)
	assert.Assert(t, ok)
	assert.Assert(t, objectProp.Properties != nil)

	nameProp, ok := objectProp.Properties["name"]
	assert.Assert(t, ok)
	_, ok = nameProp.(types.TextProperty)
	assert.Assert(t, ok)
}

func Test_insertTypedProperty_existingNestedPropertyWithNilMap(t *testing.T) {
	t.Parallel()

	props := map[string]types.Property{
		"items": types.NestedProperty{
			Type: "nested",
		},
	}

	insertTypedProperty(props, []string{"items", "value"}, estype.NewIntegerNumberProperty())

	itemsProp, ok := props["items"]
	assert.Assert(t, ok)

	nestedProp, ok := itemsProp.(types.NestedProperty)
	assert.Assert(t, ok)
	assert.Assert(t, nestedProp.Properties != nil)

	valueProp, ok := nestedProp.Properties["value"]
	assert.Assert(t, ok)
	_, ok = valueProp.(types.IntegerNumberProperty)
	assert.Assert(t, ok)
}

func Test_insertTypedProperty_replacesNonContainerWithNested(t *testing.T) {
	t.Parallel()

	props := map[string]types.Property{
		"items": types.KeywordProperty{Type: "keyword"},
	}

	insertTypedProperty(props, []string{"items", "name"}, estype.NewTextProperty())

	itemsProp, ok := props["items"]
	assert.Assert(t, ok)

	nestedProp, ok := itemsProp.(types.NestedProperty)
	assert.Assert(t, ok)
	assert.Assert(t, nestedProp.Properties != nil)

	nameProp, ok := nestedProp.Properties["name"]
	assert.Assert(t, ok)
	_, ok = nameProp.(types.TextProperty)
	assert.Assert(t, ok)
}

func Test_toTypedProperty_keywordOptions(t *testing.T) {
	t.Parallel()

	ignoreAbove := 256
	docValues := true
	index := true
	store := true
	nullValue := "unknown"
	normalizer := "lowercase"
	norms := false
	similarity := "BM25"
	eagerGlobalOrdinals := true
	splitQueriesOnWhitespace := true

	got := toTypedProperty(estype.KeywordProperty{
		IgnoreAbove:              &ignoreAbove,
		DocValues:                &docValues,
		Index:                    &index,
		Store:                    &store,
		NullValue:                &nullValue,
		Normalizer:               &normalizer,
		Norms:                    &norms,
		Similarity:               &similarity,
		EagerGlobalOrdinals:      &eagerGlobalOrdinals,
		SplitQueriesOnWhitespace: &splitQueriesOnWhitespace,
	})

	prop, ok := got.(types.KeywordProperty)
	assert.Assert(t, ok)
	assert.Assert(t, prop.IgnoreAbove != nil)
	assert.Equal(t, 256, *prop.IgnoreAbove)
	assert.Assert(t, prop.DocValues != nil)
	assert.Equal(t, true, *prop.DocValues)
	assert.Assert(t, prop.Index != nil)
	assert.Equal(t, true, *prop.Index)
	assert.Assert(t, prop.Store != nil)
	assert.Equal(t, true, *prop.Store)
	assert.Assert(t, prop.NullValue != nil)
	assert.Equal(t, "unknown", *prop.NullValue)
	assert.Assert(t, prop.Normalizer != nil)
	assert.Equal(t, "lowercase", *prop.Normalizer)
	assert.Assert(t, prop.Norms != nil)
	assert.Equal(t, false, *prop.Norms)
	assert.Assert(t, prop.Similarity != nil)
	assert.Equal(t, "BM25", *prop.Similarity)
	assert.Assert(t, prop.EagerGlobalOrdinals != nil)
	assert.Equal(t, true, *prop.EagerGlobalOrdinals)
	assert.Assert(t, prop.SplitQueriesOnWhitespace != nil)
	assert.Equal(t, true, *prop.SplitQueriesOnWhitespace)
}

func Test_toTypedProperty_textOptions(t *testing.T) {
	t.Parallel()

	searchAnalyzer := estype.Analyzer("standard")
	indexAnalyzer := estype.Analyzer("kuromoji")
	searchQuoteAnalyzer := "whitespace"
	fielddata := true
	index := true
	store := true
	norms := false
	similarity := "BM25"
	indexPhrases := true
	positionIncrementGap := 42

	got := toTypedProperty(estype.TextProperty{
		SearchAnalyzer:       &searchAnalyzer,
		IndexAnalyzer:        &indexAnalyzer,
		SearchQuoteAnalyzer:  &searchQuoteAnalyzer,
		Fielddata:            &fielddata,
		Index:                &index,
		Store:                &store,
		Norms:                &norms,
		Similarity:           &similarity,
		IndexPhrases:         &indexPhrases,
		PositionIncrementGap: &positionIncrementGap,
		Fields: map[string]estype.MappingProperty{
			"keyword": estype.NewKeywordProperty(),
		},
	})

	prop, ok := got.(types.TextProperty)
	assert.Assert(t, ok)
	assert.Assert(t, prop.SearchAnalyzer != nil)
	assert.Equal(t, "standard", *prop.SearchAnalyzer)
	assert.Assert(t, prop.Analyzer != nil)
	assert.Equal(t, "kuromoji", *prop.Analyzer)
	assert.Assert(t, prop.SearchQuoteAnalyzer != nil)
	assert.Equal(t, "whitespace", *prop.SearchQuoteAnalyzer)
	assert.Assert(t, prop.Fielddata != nil)
	assert.Equal(t, true, *prop.Fielddata)
	assert.Assert(t, prop.Index != nil)
	assert.Equal(t, true, *prop.Index)
	assert.Assert(t, prop.Store != nil)
	assert.Equal(t, true, *prop.Store)
	assert.Assert(t, prop.Norms != nil)
	assert.Equal(t, false, *prop.Norms)
	assert.Assert(t, prop.Similarity != nil)
	assert.Equal(t, "BM25", *prop.Similarity)
	assert.Assert(t, prop.IndexPhrases != nil)
	assert.Equal(t, true, *prop.IndexPhrases)
	assert.Assert(t, prop.PositionIncrementGap != nil)
	assert.Equal(t, 42, *prop.PositionIncrementGap)
	assert.Assert(t, prop.Fields != nil)
	_, ok = prop.Fields["keyword"].(types.KeywordProperty)
	assert.Assert(t, ok)
}

func Test_toTypedProperty_integerOptions(t *testing.T) {
	t.Parallel()

	coerce := true
	docValues := true
	ignoreMalformed := false
	index := true
	store := true
	nullValue := 7

	got := toTypedProperty(estype.IntegerNumberProperty{
		Coerce:          &coerce,
		DocValues:       &docValues,
		IgnoreMalformed: &ignoreMalformed,
		Index:           &index,
		Store:           &store,
		NullValue:       &nullValue,
	})

	prop, ok := got.(types.IntegerNumberProperty)
	assert.Assert(t, ok)
	assert.Assert(t, prop.Coerce != nil)
	assert.Equal(t, true, *prop.Coerce)
	assert.Assert(t, prop.DocValues != nil)
	assert.Equal(t, true, *prop.DocValues)
	assert.Assert(t, prop.IgnoreMalformed != nil)
	assert.Equal(t, false, *prop.IgnoreMalformed)
	assert.Assert(t, prop.Index != nil)
	assert.Equal(t, true, *prop.Index)
	assert.Assert(t, prop.Store != nil)
	assert.Equal(t, true, *prop.Store)
	assert.Assert(t, prop.NullValue != nil)
	assert.Equal(t, 7, *prop.NullValue)
}

func Test_toTypedProperty_dateOptions(t *testing.T) {
	t.Parallel()

	format := "strict_date_optional_time||epoch_millis"
	docValues := true
	ignoreMalformed := false
	index := true
	store := true
	locale := "ja"

	got := toTypedProperty(estype.DateProperty{
		Format:          &format,
		DocValues:       &docValues,
		IgnoreMalformed: &ignoreMalformed,
		Index:           &index,
		Store:           &store,
		Locale:          &locale,
	})

	prop, ok := got.(types.DateProperty)
	assert.Assert(t, ok)
	assert.Assert(t, prop.Format != nil)
	assert.Equal(t, format, *prop.Format)
	assert.Assert(t, prop.DocValues != nil)
	assert.Equal(t, true, *prop.DocValues)
	assert.Assert(t, prop.IgnoreMalformed != nil)
	assert.Equal(t, false, *prop.IgnoreMalformed)
	assert.Assert(t, prop.Index != nil)
	assert.Equal(t, true, *prop.Index)
	assert.Assert(t, prop.Store != nil)
	assert.Equal(t, true, *prop.Store)
	assert.Assert(t, prop.Locale != nil)
	assert.Equal(t, "ja", *prop.Locale)
}

func Test_toTypedProperty_nestedOptions(t *testing.T) {
	t.Parallel()

	enabled := true
	includeInParent := true
	includeInRoot := false
	store := true

	got := toTypedProperty(estype.NestedProperty{
		Enabled:         &enabled,
		IncludeInParent: &includeInParent,
		IncludeInRoot:   &includeInRoot,
		Store:           &store,
		Properties: map[string]estype.MappingProperty{
			"name": estype.NewTextProperty(),
		},
	})

	prop, ok := got.(types.NestedProperty)
	assert.Assert(t, ok)
	assert.Assert(t, prop.Enabled != nil)
	assert.Equal(t, true, *prop.Enabled)
	assert.Assert(t, prop.IncludeInParent != nil)
	assert.Equal(t, true, *prop.IncludeInParent)
	assert.Assert(t, prop.IncludeInRoot != nil)
	assert.Equal(t, false, *prop.IncludeInRoot)
	assert.Assert(t, prop.Store != nil)
	assert.Equal(t, true, *prop.Store)
	assert.Assert(t, prop.Properties != nil)
	_, ok = prop.Properties["name"].(types.TextProperty)
	assert.Assert(t, ok)
}

func Test_toTypedProperty_objectOptions(t *testing.T) {
	t.Parallel()

	enabled := true
	store := true

	got := toTypedProperty(estype.ObjectProperty{
		Enabled: &enabled,
		Store:   &store,
		Properties: map[string]estype.MappingProperty{
			"value": estype.NewIntegerNumberProperty(),
		},
	})

	prop, ok := got.(types.ObjectProperty)
	assert.Assert(t, ok)
	assert.Assert(t, prop.Enabled != nil)
	assert.Equal(t, true, *prop.Enabled)
	assert.Assert(t, prop.Store != nil)
	assert.Equal(t, true, *prop.Store)
	assert.Assert(t, prop.Properties != nil)
	_, ok = prop.Properties["value"].(types.IntegerNumberProperty)
	assert.Assert(t, ok)
}

func Test_toTypedPropertyFromTypeName_defaultFallback(t *testing.T) {
	t.Parallel()

	got := toTypedPropertyFromTypeName("flattened")

	prop, ok := got.(types.ObjectProperty)
	assert.Assert(t, ok)
	assert.Equal(t, "flattened", string(prop.Type))
}

type testRoundTripperFunc func(*http.Request) (*http.Response, error)

func (f testRoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type testTransportWithURLs struct {
	http.RoundTripper
	urls []*url.URL
}

func (t *testTransportWithURLs) URLs() []*url.URL {
	return t.urls
}

func newTestESClient(t *testing.T, rt http.RoundTripper) *esClient {
	t.Helper()

	typedClient, err := es8.NewTypedClient(es8.Config{
		Addresses: []string{"http://example.test"},
		Transport: rt,
	})
	assert.NilError(t, err)

	return newESClient(typedClient)
}

func newTestESClientWithURLs(t *testing.T, rt http.RoundTripper, urls []*url.URL) *esClient {
	t.Helper()

	return newTestESClient(t, &testTransportWithURLs{
		RoundTripper: rt,
		urls:         urls,
	})
}

func newJSONResponse(statusCode int, body string) *http.Response {
	header := make(http.Header)
	header.Set("X-Elastic-Product", "Elasticsearch")

	return &http.Response{
		StatusCode: statusCode,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func newTextResponse(statusCode int, body string) *http.Response {
	header := make(http.Header)
	header.Set("X-Elastic-Product", "Elasticsearch")

	return &http.Response{
		StatusCode: statusCode,
		Header:     header,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

type wrappedErr struct {
	err error
}

func (w wrappedErr) Error() string {
	return "wrapped: " + w.err.Error()
}

func (w wrappedErr) Unwrap() error {
	return w.err
}

func fmtWrap(err error) error {
	return wrappedErr{err: err}
}

func intPtr(v int) *int {
	return &v
}

func refreshIntervalPtr(v estype.RefreshInterval) *estype.RefreshInterval {
	return &v
}
