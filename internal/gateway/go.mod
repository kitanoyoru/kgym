module github.com/kitanoyoru/kgym/internal/gateway

go 1.25

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/dromara/carbon/v2 v2.6.15
	github.com/go-playground/validator/v10 v10.30.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.4
	github.com/kitanoyoru/kgym/contracts/protobuf v0.0.0-20260101185215-bc756a0f7c9a
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.11.1
	github.com/spf13/cobra v1.10.2
	go.uber.org/multierr v1.11.0
	google.golang.org/grpc v1.78.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
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
)
