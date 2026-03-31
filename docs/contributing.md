# Contributing

This document is for contributors to `es-typed-go`.

## Development prerequisites

- Go `1.26` or later
- Docker and Docker Compose
- A local Elasticsearch environment for integration tests

## Build and test commands

### Build all packages

```bash
go build ./...
```

### Run unit tests

```bash
go test ./... -short -count=1
```

### Run unit tests with verbose output

```bash
go test -v -count=1 -timeout=60s ./...
```

### Run integration tests

These require running Elasticsearch instances (start them with `docker compose up -d` first).

```bash
# v8 integration tests (ES on port 19200)
go test -tags=integration -v -timeout=120s ./estype/... ./esv8/...

# v9 integration tests (ES on port 19201)
go test -tags=integration -v -timeout=120s ./esv9/...
```

Override the Elasticsearch URL with the `ES_URL` environment variable if needed:

```bash
ES_URL=http://localhost:19200 go test -tags=integration -v -timeout=120s ./esv8/...
```

### Regenerate generated files

```bash
go generate ./...
```

## Local Elasticsearch for integration tests

The repository includes `compose.yaml` for local development. Start the services before running integration tests.

```bash
docker compose up -d
docker compose ps
```

The default local ports are:

- Elasticsearch v8: `http://localhost:19200`
- Elasticsearch v9: `http://localhost:19201`

Both services have `cluster.routing.allocation.disk.threshold_enabled: "false"` set to prevent the cluster from going read-only on machines with limited disk space.

## Linting and custom vet checks

This repository includes a custom analyzer for unchecked type assertions.

Build it:

```bash
go build -o ./bin/okassertcheck ./tools/okassertcheck
```

Run it with `go vet`:

```bash
go vet -vettool=$(pwd)/bin/okassertcheck ./...
```

A practical local validation sequence:

```bash
go test ./... -short -count=1
go build ./...
go build -o ./bin/okassertcheck ./tools/okassertcheck
go vet -vettool=$(pwd)/bin/okassertcheck ./...
```

## Adding a new Elasticsearch version

The repository follows a mirrored structure: `esv8/` and `esv9/` expose the same exported symbols. When adding a new version:

1. Copy `esv9/` to `esv{N}/` and update import paths.
2. Update `go.mod` with the new `go-elasticsearch/v{N}` dependency.
3. Add a `SearchParams.ToV{N}Request()` method to `query/search.go`.
4. Add the new `esv{N}.SearchRequest` interface and update `typed_search.go`.
5. Add integration test containers to `compose.yaml`.
6. Mirror all `esv9/` tests in `esv{N}/`.
