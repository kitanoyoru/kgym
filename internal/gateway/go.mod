module github.com/kitanoyoru/kgym/internal/gateway

go 1.25

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/dromara/carbon/v2 v2.6.15
	github.com/go-playground/validator/v10 v10.30.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.4
	github.com/kitanoyoru/kgym/contracts/protobuf v0.0.0-20260102173148-e078dcc47ef6
	github.com/kitanoyoru/kgym/pkg/tracing v0.0.0-20260102172404-d6263d0e4ece
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.11.1
	github.com/spf13/cobra v1.10.2
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.64.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0
	go.uber.org/multierr v1.11.0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/caarlos0/env/v10 v10.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/getsentry/sentry-go v0.40.0 // indirect
	github.com/getsentry/sentry-go/otel v0.40.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kitanoyoru/kgym/internal/apps/file v0.0.0-20260102172404-d6263d0e4ece // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/sdk v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251222181119-0a764e51fe1b // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251222181119-0a764e51fe1b // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace (
	github.com/kitanoyoru/kgym/pkg/database => ../../pkg/database
	github.com/kitanoyoru/kgym/pkg/grpc => ../../pkg/grpc
	github.com/kitanoyoru/kgym/pkg/metrics => ../../pkg/metrics
	github.com/kitanoyoru/kgym/pkg/testing => ../../pkg/testing
	google.golang.org/genproto => github.com/googleapis/go-genproto v0.0.0-20240116215550-a9fa1716bcac
)
