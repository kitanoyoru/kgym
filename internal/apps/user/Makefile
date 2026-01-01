.PHONY: build gomod test lint docker-build generate-mocks

gomod:
	go get -u ./...
	go mod tidy

build:
	@if [ -z "$(GOOS)" ] || [ -z "$(GOARCH)" ]; then \
		echo "Error: GOOS and GOARCH must be set"; \
		exit 1; \
	fi
	@BINARY_NAME=kgym-user-service-$(GOOS)-$(GOARCH)$(if $(filter windows,$(GOOS)),.exe,) && \
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $$BINARY_NAME ./cmd/main.go

.PHONY: test
test:
	gotestsum --format short-verbose -- ./internal/... ./migrations/...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: docker-build
docker-build:
	docker build -t kgym-user-service .

.PHONY: generate-mocks
generate-mocks:
	mockgen -source=internal/repository/repository.go -destination=internal/repository/mocks/repository_mock.go -package=mocks
