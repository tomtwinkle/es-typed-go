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

These require running Elasticsearch instances.

```bash
go test -tags=integration -v ./estype/... ./esv8/...
go test -tags=integration -v ./esv9/...
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

A practical local validation sequence is:

```bash
go test ./... -short -count=1
go build ./...
go build -o ./bin/okassertcheck ./tools/okassertcheck
go vet -vettool=$(pwd)/bin/okassertcheck ./...
```
