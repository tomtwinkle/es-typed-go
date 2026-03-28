package esv9

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	es9 "github.com/elastic/go-elasticsearch/v9"
	es9types "github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

type wrappedErr struct {
	err error
}

func (e wrappedErr) Error() string {
	return "wrapped: " + e.err.Error()
}

func (e wrappedErr) Unwrap() error {
	return e.err
}

type plainErr struct{}

func (plainErr) Error() string {
	return "plain"
}

type stubTransportWithURLs struct {
	urls []*url.URL
}

func (s stubTransportWithURLs) Perform(*http.Request) (*http.Response, error) {
	return nil, errors.New("not implemented")
}

func (s stubTransportWithURLs) URLs() []*url.URL {
	return s.urls
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type customProperty string

func (p customProperty) ESTypeName() string {
	return string(p)
}

func newTestTypedClient(t *testing.T, transport http.RoundTripper, addresses ...string) *es9.TypedClient {
	t.Helper()

	cfg := es9.Config{
		Transport: transport,
	}
	if len(addresses) > 0 {
		cfg.Addresses = addresses
	}

	client, err := es9.NewTypedClient(cfg)
	assert.NilError(t, err)

	return client
}

func newStringPtr[T ~string](v T) *T {
	return &v
}

func newBoolPtr(v bool) *bool {
	return &v
}

func newIntPtr(v int) *int {
	return &v
}

func newLoggerBuffer() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	return logger, &buf
}

func Test_isElasticsearchError_nil(t *testing.T) {
	t.Parallel()

	assert.Assert(t, !isElasticsearchError(nil, nil))
}

func Test_isElasticsearchError_direct(t *testing.T) {
	t.Parallel()

	reason := "index not found"
	esErr := &es9types.ElasticsearchError{
		Status: 404,
		ErrorCause: es9types.ErrorCause{
			Type:   "index_not_found_exception",
			Reason: &reason,
		},
	}

	var target *es9types.ElasticsearchError
	ok := isElasticsearchError(esErr, &target)

	assert.Assert(t, ok)
	assert.Assert(t, target == esErr)
	assert.Equal(t, 404, target.Status)
	assert.Equal(t, "index_not_found_exception", target.ErrorCause.Type)
	assert.Assert(t, target.ErrorCause.Reason != nil)
	assert.Equal(t, reason, *target.ErrorCause.Reason)
}

func Test_isElasticsearchError_wrapped(t *testing.T) {
	t.Parallel()

	reason := "bad query"
	esErr := &es9types.ElasticsearchError{
		Status: 400,
		ErrorCause: es9types.ErrorCause{
			Type:   "parsing_exception",
			Reason: &reason,
		},
	}

	err := wrappedErr{err: wrappedErr{err: esErr}}

	var target *es9types.ElasticsearchError
	ok := isElasticsearchError(err, &target)

	assert.Assert(t, ok)
	assert.Assert(t, target == esErr)
}

func Test_isElasticsearchError_nonElasticsearchError(t *testing.T) {
	t.Parallel()

	var target *es9types.ElasticsearchError
	ok := isElasticsearchError(wrappedErr{err: errors.New("boom")}, &target)

	assert.Assert(t, !ok)
	assert.Assert(t, target == nil)
}

func Test_unwrapErr_WithUnwrap(t *testing.T) {
	t.Parallel()

	inner := errors.New("inner")
	got := unwrapErr(wrappedErr{err: inner})

	assert.Assert(t, got == inner)
}

func Test_unwrapErr_WithoutUnwrap(t *testing.T) {
	t.Parallel()

	got := unwrapErr(plainErr{})

	assert.Assert(t, got == nil)
}

func Test_newESClient_UsesTypedClientAndDefaultLogger(t *testing.T) {
	t.Parallel()

	typedClient := newTestTypedClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{}`)),
			Request:    req,
		}, nil
	}), "http://example.test")

	client := newESClient(typedClient)

	assert.Assert(t, client != nil)
	assert.Assert(t, client.typedClient == typedClient)
	assert.Assert(t, client.logger == slog.Default())
}

func Test_esClient_schemeFromTransport_defaultsToHTTPWhenTransportHasNoURLs(t *testing.T) {
	t.Parallel()

	client := &esClient{
		typedClient: &es9.TypedClient{},
	}

	got := client.schemeFromTransport()

	assert.Equal(t, "http", got)
}

func Test_esClient_schemeFromTransport_detectsHTTPS(t *testing.T) {
	t.Parallel()

	client := &esClient{
		typedClient: &es9.TypedClient{
			BaseClient: es9.BaseClient{
				Transport: stubTransportWithURLs{
					urls: []*url.URL{
						{Scheme: "https", Host: "example.test"},
					},
				},
			},
		},
	}

	got := client.schemeFromTransport()

	assert.Equal(t, "https", got)
}

func Test_esClient_schemeFromTransport_returnsHTTPForHTTPURLs(t *testing.T) {
	t.Parallel()

	client := &esClient{
		typedClient: &es9.TypedClient{
			BaseClient: es9.BaseClient{
				Transport: stubTransportWithURLs{
					urls: []*url.URL{
						{Scheme: "http", Host: "example.test"},
					},
				},
			},
		},
	}

	got := client.schemeFromTransport()

	assert.Equal(t, "http", got)
}

func Test_esClient_buildSearchError_nil(t *testing.T) {
	t.Parallel()

	logger, _ := newLoggerBuffer()
	client := &esClient{logger: logger}

	got := client.buildSearchError(nil)

	assert.Assert(t, got == nil)
}

func Test_esClient_buildSearchError_nonElasticsearchError(t *testing.T) {
	t.Parallel()

	logger, _ := newLoggerBuffer()
	client := &esClient{logger: logger}
	original := errors.New("boom")

	got := client.buildSearchError(original)

	assert.Assert(t, got == original)
}

func Test_esClient_buildSearchError_enrichesElasticsearchError(t *testing.T) {
	t.Parallel()

	logger, buf := newLoggerBuffer()
	client := &esClient{logger: logger}

	reason := "query malformed"
	root1 := "failed to parse [status]"
	root2 := "unknown token"
	original := &es9types.ElasticsearchError{
		Status: 400,
		ErrorCause: es9types.ErrorCause{
			Type:   "search_phase_execution_exception",
			Reason: &reason,
			RootCause: []es9types.ErrorCause{
				{Type: "parsing_exception", Reason: &root1},
				{Type: "x_content_parse_exception", Reason: &root2},
				{Type: "ignored_without_reason"},
			},
		},
	}

	got := client.buildSearchError(original)

	assert.Assert(t, got != nil)
	assert.Assert(t, got != original)
	assert.Assert(t, strings.Contains(got.Error(), original.Error()))
	assert.Assert(t, strings.Contains(got.Error(), root1))
	assert.Assert(t, strings.Contains(got.Error(), root2))

	logOutput := buf.String()
	assert.Assert(t, strings.Contains(logOutput, "Elasticsearch error details"))
	assert.Assert(t, strings.Contains(logOutput, "error_status=400"))
	assert.Assert(t, strings.Contains(logOutput, "error_type=search_phase_execution_exception"))
	assert.Assert(t, strings.Contains(logOutput, "error_cause_reason=\"query malformed\""))
	assert.Assert(t, strings.Contains(logOutput, "error_cause_root_reason=\"failed to parse [status]\""))
	assert.Assert(t, strings.Contains(logOutput, "error_cause_root_reason=\"unknown token\""))
}

func Test_esClient_performRaw_successWithoutBody(t *testing.T) {
	t.Parallel()

	var gotMethod string
	var gotPath string
	var gotAccept string
	var gotContentType string
	var gotBody string

	typedClient := newTestTypedClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		gotAccept = req.Header.Get("Accept")
		gotContentType = req.Header.Get("Content-Type")

		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			gotBody = string(bodyBytes)
		}

		res := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"ok":true}`)),
			Request:    req,
		}
		res.Header.Set("X-Elastic-Product", "Elasticsearch")
		return res, nil
	}), "http://example.test")

	client := &esClient{typedClient: typedClient, logger: slog.Default()}

	got, err := client.performRaw(context.Background(), http.MethodGet, "/_internal/test", nil)

	assert.NilError(t, err)
	assert.DeepEqual(t, []byte(`{"ok":true}`), []byte(got))
	assert.Equal(t, http.MethodGet, gotMethod)
	assert.Equal(t, "/_internal/test", gotPath)
	assert.Equal(t, "application/vnd.elasticsearch+json;compatible-with=9", gotAccept)
	assert.Equal(t, "", gotContentType)
	assert.Equal(t, "", gotBody)
}

func Test_esClient_performRaw_successWithBody(t *testing.T) {
	t.Parallel()

	var gotMethod string
	var gotPath string
	var gotAccept string
	var gotContentType string
	var gotBody string

	typedClient := newTestTypedClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		gotMethod = req.Method
		gotPath = req.URL.Path
		gotAccept = req.Header.Get("Accept")
		gotContentType = req.Header.Get("Content-Type")

		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			assert.NilError(t, err)
			gotBody = string(bodyBytes)
		}

		res := &http.Response{
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"created":true}`)),
			Request:    req,
		}
		res.Header.Set("X-Elastic-Product", "Elasticsearch")
		return res, nil
	}), "https://example.test")

	client := &esClient{typedClient: typedClient, logger: slog.Default()}

	got, err := client.performRaw(
		context.Background(),
		http.MethodPost,
		"/_internal/test",
		json.RawMessage(`{"name":"value"}`),
	)

	assert.NilError(t, err)
	assert.DeepEqual(t, []byte(`{"created":true}`), []byte(got))
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/_internal/test", gotPath)
	assert.Equal(t, "application/vnd.elasticsearch+json;compatible-with=9", gotAccept)
	assert.Equal(t, "application/vnd.elasticsearch+json;compatible-with=9", gotContentType)
	assert.Equal(t, `{"name":"value"}`, gotBody)
}

func Test_esClient_performRaw_invalidPath(t *testing.T) {
	t.Parallel()

	typedClient := newTestTypedClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("should not be called")
	}), "http://example.test")

	client := &esClient{typedClient: typedClient, logger: slog.Default()}

	got, err := client.performRaw(context.Background(), http.MethodGet, "%", nil)

	assert.Assert(t, got == nil)
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "building HTTP request GET %"))
}

func Test_esClient_performRaw_transportError(t *testing.T) {
	t.Parallel()

	typedClient := newTestTypedClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network down")
	}), "http://example.test")

	client := &esClient{typedClient: typedClient, logger: slog.Default()}

	got, err := client.performRaw(context.Background(), http.MethodDelete, "/_internal/test", nil)

	assert.Assert(t, got == nil)
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "performing DELETE /_internal/test"))
	assert.Assert(t, strings.Contains(err.Error(), "network down"))
}

type errReadCloser struct {
	err error
}

func (e errReadCloser) Read([]byte) (int, error) {
	return 0, e.err
}

func (e errReadCloser) Close() error {
	return nil
}

func Test_esClient_performRaw_readBodyError(t *testing.T) {
	t.Parallel()

	typedClient := newTestTypedClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		res := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       errReadCloser{err: errors.New("read failed")},
			Request:    req,
		}
		res.Header.Set("X-Elastic-Product", "Elasticsearch")
		return res, nil
	}), "http://example.test")

	client := &esClient{typedClient: typedClient, logger: slog.Default()}

	got, err := client.performRaw(context.Background(), http.MethodGet, "/_internal/test", nil)

	assert.Assert(t, got == nil)
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "reading response body for GET /_internal/test"))
	assert.Assert(t, strings.Contains(err.Error(), "read failed"))
}

func Test_esClient_performRaw_httpErrorStatusIncludesTrimmedBody(t *testing.T) {
	t.Parallel()

	longBody := strings.Repeat("x", 600)

	typedClient := newTestTypedClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		res := &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(longBody)),
			Request:    req,
		}
		res.Header.Set("X-Elastic-Product", "Elasticsearch")
		return res, nil
	}), "http://example.test")

	client := &esClient{typedClient: typedClient, logger: slog.Default()}

	got, err := client.performRaw(context.Background(), http.MethodGet, "/_internal/test", nil)

	assert.Assert(t, got == nil)
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "unexpected status 400 for GET /_internal/test"))
	assert.Assert(t, strings.Contains(err.Error(), strings.Repeat("x", 512)))
	assert.Assert(t, !strings.Contains(err.Error(), strings.Repeat("x", 513)))
}

func Test_taskIDToString_supportedTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   any
		want string
	}{
		{name: "string", in: "node:123", want: "node:123"},
		{name: "int", in: int(1), want: "1"},
		{name: "int8", in: int8(2), want: "2"},
		{name: "int16", in: int16(3), want: "3"},
		{name: "int32", in: int32(4), want: "4"},
		{name: "int64", in: int64(5), want: "5"},
		{name: "uint", in: uint(6), want: "6"},
		{name: "uint8", in: uint8(7), want: "7"},
		{name: "uint16", in: uint16(8), want: "8"},
		{name: "uint32", in: uint32(9), want: "9"},
		{name: "uint64", in: uint64(10), want: "10"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := taskIDToString(tt.in)

			assert.NilError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_taskIDToString_unsupportedType(t *testing.T) {
	t.Parallel()

	got, err := taskIDToString(struct{}{})

	assert.Equal(t, "", got)
	assert.Assert(t, err != nil)
	assert.Assert(t, strings.Contains(err.Error(), "unsupported task ID type: struct {}"))
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

func Test_insertTypedProperty_ignoresEmptyPathAndNilProperty(t *testing.T) {
	t.Parallel()

	props := map[string]es9types.Property{
		"status": es9types.KeywordProperty{Type: "keyword"},
	}

	insertTypedProperty(props, nil, estype.NewTextProperty())
	insertTypedProperty(props, []string{""}, estype.NewTextProperty())
	insertTypedProperty(props, []string{"title"}, nil)

	assert.Equal(t, 1, len(props))
	_, ok := props["status"].(es9types.KeywordProperty)
	assert.Assert(t, ok)
}

func Test_insertTypedProperty_extendsExistingObjectProperty(t *testing.T) {
	t.Parallel()

	props := map[string]es9types.Property{
		"items": es9types.ObjectProperty{
			Type:       "object",
			Properties: map[string]es9types.Property{},
		},
	}

	insertTypedProperty(props, []string{"items", "name"}, estype.NewTextProperty())

	itemsProp, ok := props["items"]
	assert.Assert(t, ok)

	objectProp, ok := itemsProp.(es9types.ObjectProperty)
	assert.Assert(t, ok)
	assert.Assert(t, objectProp.Properties != nil)

	nameProp, ok := objectProp.Properties["name"]
	assert.Assert(t, ok)
	_, ok = nameProp.(es9types.TextProperty)
	assert.Assert(t, ok)
}

func Test_insertTypedProperty_replacesLeafWithNestedHierarchy(t *testing.T) {
	t.Parallel()

	props := map[string]es9types.Property{
		"items": es9types.KeywordProperty{Type: "keyword"},
	}

	insertTypedProperty(props, []string{"items", "name"}, estype.NewTextProperty())

	itemsProp, ok := props["items"]
	assert.Assert(t, ok)

	nestedProp, ok := itemsProp.(es9types.NestedProperty)
	assert.Assert(t, ok)
	assert.Assert(t, nestedProp.Properties != nil)

	nameProp, ok := nestedProp.Properties["name"]
	assert.Assert(t, ok)
	_, ok = nameProp.(es9types.TextProperty)
	assert.Assert(t, ok)
}

func Test_toTypedProperty_keywordCopiesAllSupportedFields(t *testing.T) {
	t.Parallel()

	prop := estype.KeywordProperty{
		IgnoreAbove:              newIntPtr(256),
		DocValues:                newBoolPtr(true),
		Index:                    newBoolPtr(false),
		Store:                    newBoolPtr(true),
		NullValue:                newStringPtr("missing"),
		Normalizer:               newStringPtr("lowercase"),
		Norms:                    newBoolPtr(true),
		Similarity:               newStringPtr("BM25"),
		EagerGlobalOrdinals:      newBoolPtr(true),
		SplitQueriesOnWhitespace: newBoolPtr(true),
	}

	got := toTypedProperty(prop)

	keywordProp, ok := got.(es9types.KeywordProperty)
	assert.Assert(t, ok)
	assert.Assert(t, keywordProp.IgnoreAbove != nil)
	assert.Equal(t, 256, *keywordProp.IgnoreAbove)
	assert.Assert(t, keywordProp.DocValues != nil)
	assert.Equal(t, true, *keywordProp.DocValues)
	assert.Assert(t, keywordProp.Index != nil)
	assert.Equal(t, false, *keywordProp.Index)
	assert.Assert(t, keywordProp.Store != nil)
	assert.Equal(t, true, *keywordProp.Store)
	assert.Assert(t, keywordProp.NullValue != nil)
	assert.Equal(t, "missing", *keywordProp.NullValue)
	assert.Assert(t, keywordProp.Normalizer != nil)
	assert.Equal(t, "lowercase", *keywordProp.Normalizer)
	assert.Assert(t, keywordProp.Norms != nil)
	assert.Equal(t, true, *keywordProp.Norms)
	assert.Assert(t, keywordProp.Similarity != nil)
	assert.Equal(t, "BM25", *keywordProp.Similarity)
	assert.Assert(t, keywordProp.EagerGlobalOrdinals != nil)
	assert.Equal(t, true, *keywordProp.EagerGlobalOrdinals)
	assert.Assert(t, keywordProp.SplitQueriesOnWhitespace != nil)
	assert.Equal(t, true, *keywordProp.SplitQueriesOnWhitespace)
}

func Test_toTypedProperty_textCopiesAllSupportedFields(t *testing.T) {
	t.Parallel()

	prop := estype.TextProperty{
		SearchAnalyzer:       newStringPtr(estype.Analyzer("standard")),
		IndexAnalyzer:        newStringPtr(estype.Analyzer("kuromoji")),
		SearchQuoteAnalyzer:  newStringPtr("whitespace"),
		Fielddata:            newBoolPtr(true),
		Index:                newBoolPtr(false),
		Store:                newBoolPtr(true),
		Norms:                newBoolPtr(false),
		Similarity:           newStringPtr("boolean"),
		IndexPhrases:         newBoolPtr(true),
		PositionIncrementGap: newIntPtr(7),
		Fields: map[string]estype.MappingProperty{
			"keyword": estype.NewKeywordProperty(estype.WithIgnoreAbove(128)),
		},
	}

	got := toTypedProperty(prop)

	textProp, ok := got.(es9types.TextProperty)
	assert.Assert(t, ok)
	assert.Assert(t, textProp.SearchAnalyzer != nil)
	assert.Equal(t, "standard", *textProp.SearchAnalyzer)
	assert.Assert(t, textProp.Analyzer != nil)
	assert.Equal(t, "kuromoji", *textProp.Analyzer)
	assert.Assert(t, textProp.SearchQuoteAnalyzer != nil)
	assert.Equal(t, "whitespace", *textProp.SearchQuoteAnalyzer)
	assert.Assert(t, textProp.Fielddata != nil)
	assert.Equal(t, true, *textProp.Fielddata)
	assert.Assert(t, textProp.Index != nil)
	assert.Equal(t, false, *textProp.Index)
	assert.Assert(t, textProp.Store != nil)
	assert.Equal(t, true, *textProp.Store)
	assert.Assert(t, textProp.Norms != nil)
	assert.Equal(t, false, *textProp.Norms)
	assert.Assert(t, textProp.Similarity != nil)
	assert.Equal(t, "boolean", *textProp.Similarity)
	assert.Assert(t, textProp.IndexPhrases != nil)
	assert.Equal(t, true, *textProp.IndexPhrases)
	assert.Assert(t, textProp.PositionIncrementGap != nil)
	assert.Equal(t, 7, *textProp.PositionIncrementGap)
	assert.Assert(t, textProp.Fields != nil)

	fieldProp, ok := textProp.Fields["keyword"]
	assert.Assert(t, ok)
	keywordProp, ok := fieldProp.(es9types.KeywordProperty)
	assert.Assert(t, ok)
	assert.Assert(t, keywordProp.IgnoreAbove != nil)
	assert.Equal(t, 128, *keywordProp.IgnoreAbove)
}

func Test_toTypedProperty_integerDateNestedObjectAndFallback(t *testing.T) {
	t.Parallel()

	t.Run("integer", func(t *testing.T) {
		t.Parallel()

		prop := estype.IntegerNumberProperty{
			Coerce:          newBoolPtr(true),
			DocValues:       newBoolPtr(false),
			IgnoreMalformed: newBoolPtr(true),
			Index:           newBoolPtr(false),
			Store:           newBoolPtr(true),
			NullValue:       func() *int { v := 42; return &v }(),
		}

		got := toTypedProperty(prop)

		intProp, ok := got.(es9types.IntegerNumberProperty)
		assert.Assert(t, ok)
		assert.Assert(t, intProp.Coerce != nil)
		assert.Equal(t, true, *intProp.Coerce)
		assert.Assert(t, intProp.DocValues != nil)
		assert.Equal(t, false, *intProp.DocValues)
		assert.Assert(t, intProp.IgnoreMalformed != nil)
		assert.Equal(t, true, *intProp.IgnoreMalformed)
		assert.Assert(t, intProp.Index != nil)
		assert.Equal(t, false, *intProp.Index)
		assert.Assert(t, intProp.Store != nil)
		assert.Equal(t, true, *intProp.Store)
		assert.Assert(t, intProp.NullValue != nil)
		assert.Equal(t, 42, *intProp.NullValue)
	})

	t.Run("date", func(t *testing.T) {
		t.Parallel()

		prop := estype.DateProperty{
			Format:          newStringPtr("strict_date_optional_time"),
			DocValues:       newBoolPtr(true),
			IgnoreMalformed: newBoolPtr(false),
			Index:           newBoolPtr(true),
			Store:           newBoolPtr(false),
			Locale:          newStringPtr("ja-JP"),
		}

		got := toTypedProperty(prop)

		dateProp, ok := got.(es9types.DateProperty)
		assert.Assert(t, ok)
		assert.Assert(t, dateProp.Format != nil)
		assert.Equal(t, "strict_date_optional_time", *dateProp.Format)
		assert.Assert(t, dateProp.DocValues != nil)
		assert.Equal(t, true, *dateProp.DocValues)
		assert.Assert(t, dateProp.IgnoreMalformed != nil)
		assert.Equal(t, false, *dateProp.IgnoreMalformed)
		assert.Assert(t, dateProp.Index != nil)
		assert.Equal(t, true, *dateProp.Index)
		assert.Assert(t, dateProp.Store != nil)
		assert.Equal(t, false, *dateProp.Store)
		assert.Assert(t, dateProp.Locale != nil)
		assert.Equal(t, "ja-JP", *dateProp.Locale)
	})

	t.Run("nested", func(t *testing.T) {
		t.Parallel()

		prop := estype.NestedProperty{
			Enabled:         newBoolPtr(true),
			IncludeInParent: newBoolPtr(true),
			IncludeInRoot:   newBoolPtr(false),
			Store:           newBoolPtr(true),
			Properties: map[string]estype.MappingProperty{
				"name": estype.NewTextProperty(),
			},
		}

		got := toTypedProperty(prop)

		nestedProp, ok := got.(es9types.NestedProperty)
		assert.Assert(t, ok)
		assert.Assert(t, nestedProp.Enabled != nil)
		assert.Equal(t, true, *nestedProp.Enabled)
		assert.Assert(t, nestedProp.IncludeInParent != nil)
		assert.Equal(t, true, *nestedProp.IncludeInParent)
		assert.Assert(t, nestedProp.IncludeInRoot != nil)
		assert.Equal(t, false, *nestedProp.IncludeInRoot)
		assert.Assert(t, nestedProp.Store != nil)
		assert.Equal(t, true, *nestedProp.Store)
		assert.Assert(t, nestedProp.Properties != nil)

		nameProp, ok := nestedProp.Properties["name"]
		assert.Assert(t, ok)
		_, ok = nameProp.(es9types.TextProperty)
		assert.Assert(t, ok)
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()

		prop := estype.ObjectProperty{
			Enabled: newBoolPtr(false),
			Store:   newBoolPtr(true),
			Properties: map[string]estype.MappingProperty{
				"status": estype.NewKeywordProperty(),
			},
		}

		got := toTypedProperty(prop)

		objectProp, ok := got.(es9types.ObjectProperty)
		assert.Assert(t, ok)
		assert.Assert(t, objectProp.Enabled != nil)
		assert.Equal(t, false, *objectProp.Enabled)
		assert.Assert(t, objectProp.Store != nil)
		assert.Equal(t, true, *objectProp.Store)
		assert.Assert(t, objectProp.Properties != nil)

		statusProp, ok := objectProp.Properties["status"]
		assert.Assert(t, ok)
		_, ok = statusProp.(es9types.KeywordProperty)
		assert.Assert(t, ok)
	})

	t.Run("fallback from unknown mapping property", func(t *testing.T) {
		t.Parallel()

		got := toTypedProperty(customProperty("flattened"))

		objectProp, ok := got.(es9types.ObjectProperty)
		assert.Assert(t, ok)
		assert.Equal(t, "flattened", objectProp.Type)
	})
}

func Test_toTypedPropertyFromTypeName_knownAndUnknown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		typeName string
		wantType reflect.Type
		wantES   string
	}{
		{name: "keyword", typeName: "keyword", wantType: reflect.TypeOf(es9types.KeywordProperty{}), wantES: "keyword"},
		{name: "text", typeName: "text", wantType: reflect.TypeOf(es9types.TextProperty{}), wantES: "text"},
		{name: "integer", typeName: "integer", wantType: reflect.TypeOf(es9types.IntegerNumberProperty{}), wantES: "integer"},
		{name: "date", typeName: "date", wantType: reflect.TypeOf(es9types.DateProperty{}), wantES: "date"},
		{name: "nested", typeName: "nested", wantType: reflect.TypeOf(es9types.NestedProperty{}), wantES: "nested"},
		{name: "object", typeName: "object", wantType: reflect.TypeOf(es9types.ObjectProperty{}), wantES: "object"},
		{name: "unknown", typeName: "flattened", wantType: reflect.TypeOf(es9types.ObjectProperty{}), wantES: "flattened"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := toTypedPropertyFromTypeName(tt.typeName)

			assert.Equal(t, tt.wantType, reflect.TypeOf(got))

			switch v := got.(type) {
			case es9types.KeywordProperty:
				assert.Equal(t, tt.wantES, v.Type)
			case es9types.TextProperty:
				assert.Equal(t, tt.wantES, v.Type)
			case es9types.IntegerNumberProperty:
				assert.Equal(t, tt.wantES, v.Type)
			case es9types.DateProperty:
				assert.Equal(t, tt.wantES, v.Type)
			case es9types.NestedProperty:
				assert.Equal(t, tt.wantES, v.Type)
			case es9types.ObjectProperty:
				assert.Equal(t, tt.wantES, v.Type)
			default:
				t.Fatalf("unexpected property type %T", got)
			}
		})
	}
}
