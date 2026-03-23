package esv8

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/count"
	coredelete "github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	coreget "github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	coreidx "github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/info"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/reindex"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	idxcreate "github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	idxdelete "github.com/elastic/go-elasticsearch/v8/typedapi/indices/delete"
	idxputalias "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putalias"
	idxputsettings "github.com/elastic/go-elasticsearch/v8/typedapi/indices/putsettings"
	idxrefresh "github.com/elastic/go-elasticsearch/v8/typedapi/indices/refresh"
	idxupdatealiases "github.com/elastic/go-elasticsearch/v8/typedapi/indices/updatealiases"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/conflicts"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/optype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"github.com/tomtwinkle/es-typed-go/estype"
)

// esClient is the concrete implementation of ESClient.
type esClient struct {
	typedClient *es8.TypedClient
	logger      *slog.Logger
}

// newESClient creates a new esClient with the given typed client.
// Uses the default slog logger.
func newESClient(typedClient *es8.TypedClient) *esClient {
	return &esClient{
		typedClient: typedClient,
		logger:      slog.Default(),
	}
}

// NewClientWithLogger constructs an ESClient backed by the Elasticsearch v8 typed client
// using a custom slog.Logger.
func NewClientWithLogger(config es8.Config, logger *slog.Logger) (ESClient, error) {
	typedClient, err := es8.NewTypedClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch TypedClient: %w", err)
	}
	c := newESClient(typedClient)
	if logger != nil {
		c.logger = logger
	}
	return c, nil
}

func (c *esClient) Info(ctx context.Context) (*info.Response, error) {
	return c.typedClient.Info().Do(ctx)
}

func (c *esClient) IndexRefresh(ctx context.Context, indexName estype.Index, opts ...IndexRefreshOption) (*idxrefresh.Response, error) {
	b := c.typedClient.Indices.Refresh().Index(indexName.String())
	for _, opt := range opts {
		opt(b)
	}
	return b.Do(ctx)
}

func (c *esClient) AliasRefresh(ctx context.Context, aliasName estype.Alias, opts ...IndexRefreshOption) (*idxrefresh.Response, error) {
	// Elasticsearch does not support Refresh on an alias directly;
	// resolve the alias to its backing index first.
	indices, err := c.GetIndicesForAlias(ctx, aliasName)
	if err != nil {
		return nil, fmt.Errorf("failed to get indices for alias %s: %w", aliasName, err)
	}
	if len(indices) == 0 {
		return nil, fmt.Errorf("no indices found for alias %s", aliasName)
	}
	b := c.typedClient.Indices.Refresh().Index(indices[0].String())
	for _, opt := range opts {
		opt(b)
	}
	return b.Do(ctx)
}

func (c *esClient) IndexDocumentCount(ctx context.Context, indexName estype.Index) (*count.Response, error) {
	return c.typedClient.Count().Index(indexName.String()).Do(ctx)
}

func (c *esClient) CreateIndex(
	ctx context.Context,
	indexName estype.Index,
	settings *types.IndexSettings,
	mappings *types.TypeMapping,
) (*idxcreate.Response, error) {
	req := idxcreate.NewRequest()
	if settings != nil {
		req.Settings = settings
	}
	if mappings != nil {
		req.Mappings = mappings
	}
	return c.typedClient.Indices.Create(indexName.String()).Request(req).Do(ctx)
}

func (c *esClient) DeleteIndex(ctx context.Context, indexName estype.Index, opts ...DeleteIndexOption) (*idxdelete.Response, error) {
	b := c.typedClient.Indices.Delete(indexName.String())
	for _, opt := range opts {
		opt(b)
	}
	return b.Do(ctx)
}

func (c *esClient) IndexExists(ctx context.Context, indexName estype.Index) (bool, error) {
	return c.typedClient.Indices.Exists(indexName.String()).Do(ctx)
}

func (c *esClient) AliasExists(ctx context.Context, aliasName estype.Alias) (bool, error) {
	return c.typedClient.Indices.ExistsAlias(aliasName.String()).Do(ctx)
}

func (c *esClient) GetIndicesForAlias(ctx context.Context, aliasName estype.Alias) ([]estype.Index, error) {
	res, err := c.typedClient.Indices.GetAlias().Name(aliasName.String()).Do(ctx)
	if err != nil {
		// If the alias doesn't exist Elasticsearch returns 404; return an empty slice.
		var esErr *types.ElasticsearchError
		if isElasticsearchError(err, &esErr) && esErr.Status == http.StatusNotFound {
			return []estype.Index{}, nil
		}
		return nil, fmt.Errorf("failed to get indices for alias %s: %w", aliasName, err)
	}

	indices := make([]estype.Index, 0, len(res))
	for indexName := range res {
		esIndex, err := estype.ParseESIndex(indexName)
		if err != nil {
			return nil, fmt.Errorf("failed to parse index name %s: %w", indexName, err)
		}
		indices = append(indices, esIndex)
	}
	return indices, nil
}

func (c *esClient) CreateAlias(
	ctx context.Context,
	indexName estype.Index,
	aliasName estype.Alias,
	isWriteIndex bool,
) (*idxputalias.Response, error) {
	req := idxputalias.NewRequest()
	req.IsWriteIndex = &isWriteIndex
	return c.typedClient.Indices.PutAlias(indexName.String(), aliasName.String()).Request(req).Do(ctx)
}

func (c *esClient) UpdateAliases(
	ctx context.Context, actions []types.IndicesAction,
) (*idxupdatealiases.Response, error) {
	req := idxupdatealiases.NewRequest()
	req.Actions = actions
	return c.typedClient.Indices.UpdateAliases().Request(req).Do(ctx)
}

// GetRefreshInterval returns the current refresh interval for the index backing the alias.
//
// If the alias points to multiple indices, only the first index's setting is returned.
// If refresh_interval is not explicitly set, RefreshIntervalNotSet (0) is returned.
func (c *esClient) GetRefreshInterval(
	ctx context.Context,
	aliasName estype.Alias,
) (estype.RefreshInterval, error) {
	indices, err := c.GetIndicesForAlias(ctx, aliasName)
	if err != nil {
		return 0, fmt.Errorf("failed to get indices for alias %s: %w", aliasName, err)
	}
	if len(indices) == 0 {
		return 0, fmt.Errorf("no indices found for alias %s", aliasName)
	}

	res, err := c.typedClient.Indices.GetSettings().Index(indices[0].String()).Do(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get settings for index %s: %w", indices[0], err)
	}

	indexSettings, ok := res[indices[0].String()]
	if !ok {
		return 0, fmt.Errorf("settings not found for index %s", indices[0])
	}

	if indexSettings.Settings != nil &&
		indexSettings.Settings.Index != nil &&
		indexSettings.Settings.Index.RefreshInterval != nil {
		refreshInterval, ok := indexSettings.Settings.Index.RefreshInterval.(string)
		if !ok {
			return 0, fmt.Errorf("invalid refresh_interval format for index %s", indices[0])
		}
		return estype.ParseRefreshInterval(refreshInterval)
	}

	// Not explicitly set; caller should treat this as "use default".
	return estype.RefreshIntervalNotSet, nil
}

func (c *esClient) UpdateRefreshInterval(
	ctx context.Context,
	aliasName estype.Alias,
	interval estype.RefreshInterval,
) (*idxputsettings.Response, error) {
	indices, err := c.GetIndicesForAlias(ctx, aliasName)
	if err != nil {
		return nil, fmt.Errorf("failed to get indices for alias %s: %w", aliasName, err)
	}
	if len(indices) == 0 {
		return nil, fmt.Errorf("no indices found for alias %s", aliasName)
	}
	res, err := c.typedClient.Indices.PutSettings().
		Indices(indices[0].String()).
		RefreshInterval(interval.ESTypeDuration()).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh interval for alias %s: %w", aliasName, err)
	}
	return res, nil
}

func (c *esClient) CreateDocument(
	ctx context.Context,
	aliasName estype.Alias,
	id string,
	document any,
	opts ...CreateDocumentOption,
) (*coreidx.Response, error) {
	// Set Refresh to WaitFor so the document is visible immediately after the call returns.
	b := c.typedClient.Index(aliasName.String()).Id(id).Document(document).Refresh(refresh.Waitfor)
	for _, opt := range opts {
		opt(b)
	}
	return b.Do(ctx)
}

func (c *esClient) GetDocument(
	ctx context.Context,
	aliasName estype.Alias,
	id string,
	opts ...GetDocumentOption,
) (*coreget.Response, error) {
	b := c.typedClient.Get(aliasName.String(), id)
	for _, opt := range opts {
		opt(b)
	}
	return b.Do(ctx)
}

func (c *esClient) DeleteDocument(ctx context.Context, indexName estype.Index, id string, opts ...DeleteDocumentOption) (*coredelete.Response, error) {
	b := c.typedClient.Delete(indexName.String(), id)
	for _, opt := range opts {
		opt(b)
	}
	return b.Do(ctx)
}

func (c *esClient) UpdateDocument(ctx context.Context, indexName estype.Index, id string, req *update.Request) (*update.Response, error) {
	return c.typedClient.Update(indexName.String(), id).Request(req).Do(ctx)
}

func (c *esClient) SearchRaw(ctx context.Context, aliasName estype.Alias, req *search.Request) (*search.Response, error) {
	c.logger.DebugContext(ctx, "Elasticsearch SearchRaw",
		slog.String("alias", aliasName.String()),
	)
	res, err := c.typedClient.Search().Index(aliasName.String()).Request(req).Do(ctx)
	if err != nil {
		searchErr := c.buildSearchError(err)
		c.logger.ErrorContext(ctx, "Elasticsearch SearchRaw failed",
			slog.String("alias", aliasName.String()),
			slog.Any("error", searchErr),
		)
		return nil, searchErr
	}
	return res, nil
}

// buildSearchError enriches an error from a Search call with Elasticsearch error details.
func (c *esClient) buildSearchError(err error) error {
	if err == nil {
		return nil
	}
	var esErr *types.ElasticsearchError
	if !isElasticsearchError(err, &esErr) {
		return err
	}

	attrs := []any{
		slog.Int("error_status", esErr.Status),
		slog.String("error_type", esErr.ErrorCause.Type),
	}
	if esErr.ErrorCause.Reason != nil {
		attrs = append(attrs, slog.String("error_cause_reason", *esErr.ErrorCause.Reason))
	}

	combinedErr := err
	for _, cause := range esErr.ErrorCause.RootCause {
		if cause.Reason != nil {
			combinedErr = fmt.Errorf("%w; %s", combinedErr, *cause.Reason)
			attrs = append(attrs, slog.String("error_cause_root_reason", *cause.Reason))
		}
	}

	c.logger.Error("Elasticsearch error details", attrs...)
	return combinedErr
}

func (c *esClient) Reindex(
	ctx context.Context, sourceIndex, destIndex estype.Index, waitForCompletion bool, opts ...ReindexOption,
) (*reindex.Response, error) {
	proceed := conflicts.Proceed
	req := &reindex.Request{
		Conflicts: &proceed,
		Source:    types.ReindexSource{Index: []string{sourceIndex.String()}},
		Dest: types.ReindexDestination{
			Index:  destIndex.String(),
			OpType: &optype.Index,
		},
	}
	b := c.typedClient.Reindex().Request(req).WaitForCompletion(waitForCompletion)
	for _, opt := range opts {
		opt(b)
	}
	return b.Do(ctx)
}

// DeltaReindex copies only documents updated since `since` from sourceIndex to destIndex.
func (c *esClient) DeltaReindex(
	ctx context.Context, sourceIndex, destIndex estype.Index, since time.Time, timestampField string,
	waitForCompletion bool, opts ...ReindexOption,
) (*reindex.Response, error) {
	req := reindex.NewRequest()

	sinceStr := since.Format(time.RFC3339Nano)
	rangeQuery := types.NewDateRangeQuery()
	rangeQuery.Gte = &sinceStr

	proceed := conflicts.Proceed
	req.Conflicts = &proceed
	req.Source = types.ReindexSource{
		Index: []string{sourceIndex.String()},
		Query: &types.Query{
			Range: map[string]types.RangeQuery{
				timestampField: rangeQuery,
			},
		},
	}
	req.Dest = types.ReindexDestination{
		Index:  destIndex.String(),
		OpType: &optype.Index,
	}

	b := c.typedClient.Reindex().Request(req).WaitForCompletion(waitForCompletion)
	for _, opt := range opts {
		opt(b)
	}
	return b.Do(ctx)
}

func (c *esClient) WaitForTaskCompletion(ctx context.Context, taskID types.TaskId, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	taskIDStr, err := taskIDToString(taskID)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for task completion: %w", ctx.Err())
		case <-ticker.C:
			res, err := c.typedClient.Tasks.Get(taskIDStr).Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to get task status for task ID %s: %w", taskIDStr, err)
			}
			if res.Completed {
				if res.Error != nil {
					reason := "unknown"
					if res.Error.Reason != nil {
						reason = *res.Error.Reason
					}
					return fmt.Errorf("task %s completed with error: type=%s, reason=%s",
						taskIDStr, res.Error.Type, reason)
				}
				c.logger.InfoContext(ctx, "Task completed successfully",
					slog.String("task_id", taskIDStr))
				return nil
			}
			c.logger.InfoContext(ctx, "Waiting for task to complete...",
				slog.String("task_id", taskIDStr))
		}
	}
}

// taskIDToString converts a types.TaskId (any) to its string representation.
func taskIDToString(taskID types.TaskId) (string, error) {
	switch v := taskID.(type) {
	case string:
		return v, nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case int8:
		return fmt.Sprintf("%d", v), nil
	case int16:
		return fmt.Sprintf("%d", v), nil
	case int32:
		return fmt.Sprintf("%d", v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case uint:
		return fmt.Sprintf("%d", v), nil
	case uint8:
		return fmt.Sprintf("%d", v), nil
	case uint16:
		return fmt.Sprintf("%d", v), nil
	case uint32:
		return fmt.Sprintf("%d", v), nil
	case uint64:
		return fmt.Sprintf("%d", v), nil
	default:
		return "", fmt.Errorf("unsupported task ID type: %T", v)
	}
}

// isElasticsearchError checks if err is an *types.ElasticsearchError and sets target.
func isElasticsearchError(err error, target **types.ElasticsearchError) bool {
	if err == nil {
		return false
	}
	var esErr *types.ElasticsearchError
	ok := false
	// Walk the error chain.
	e := err
	for e != nil {
		ee, cast := e.(*types.ElasticsearchError)
		if !cast {
			e = unwrapErr(e)
			continue
		}
		esErr = ee
		ok = true
		break
	}
	if ok && target != nil {
		*target = esErr
	}
	return ok
}

// unwrapErr calls Unwrap on the error if available.
func unwrapErr(err error) error {
	type unwrapper interface {
		Unwrap() error
	}
	u, ok := err.(unwrapper)
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// schemeFromTransport returns "https" when the transport's connection pool is
// configured with HTTPS addresses, and "http" otherwise. It uses a local
// interface to avoid importing elastictransport directly.
func (c *esClient) schemeFromTransport() string {
	type urlProvider interface {
		URLs() []*url.URL
	}
	up, ok := c.typedClient.Transport.(urlProvider)
	if !ok {
		return "http"
	}
	if urls := up.URLs(); len(urls) > 0 && urls[0].Scheme == "https" {
		return "https"
	}
	return "http"
}

// performRaw executes a raw HTTP request against the Elasticsearch cluster and
// returns the response body as a json.RawMessage.
// path must be a URL path starting with "/" (e.g. "/_internal/desired_balance").
// body may be nil for requests without a request body.
//
// The scheme (http or https) is detected from the transport's connection pool,
// so both plain HTTP and TLS-secured HTTPS clusters are supported.
// The host part of the URL is left empty; the transport fills it in from the
// connection pool at request time, exactly as the go-elasticsearch TypedAPI does.
func (c *esClient) performRaw(ctx context.Context, method, path string, body json.RawMessage) (json.RawMessage, error) {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.schemeFromTransport()+"://"+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("building HTTP request %s %s: %w", method, path, err)
	}
	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/vnd.elasticsearch+json;compatible-with=8")
	}
	req.Header.Set("Accept", "application/vnd.elasticsearch+json;compatible-with=8")

	res, err := c.typedClient.Perform(req)
	if err != nil {
		return nil, fmt.Errorf("performing %s %s: %w", method, path, err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body for %s %s: %w", method, path, err)
	}

	if res.StatusCode >= 400 {
		const maxErrBodyLen = 512
		errBody := data
		if len(errBody) > maxErrBodyLen {
			errBody = errBody[:maxErrBodyLen]
		}
		return nil, fmt.Errorf("unexpected status %d for %s %s: %s", res.StatusCode, method, path, errBody)
	}

	return json.RawMessage(data), nil
}
