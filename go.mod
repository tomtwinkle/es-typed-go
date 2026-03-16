module github.com/tomtwinkle/es-typed-go

go 1.26

require (
	github.com/elastic/go-elasticsearch/v8 v8.19.3
	github.com/elastic/go-elasticsearch/v9 v9.3.1
	github.com/google/uuid v1.6.0
	gotest.tools/v3 v3.5.2
)

require (
	github.com/bitfield/gotestdox v0.2.2 // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.8.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/term v0.35.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	gotest.tools/gotestsum v1.13.0 // indirect
)

tool (
	github.com/tomtwinkle/es-typed-go/cmd/estyped
	gotest.tools/gotestsum
)
