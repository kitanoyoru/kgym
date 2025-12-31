.PHONY: build
build:
	@if [ -z "$(GOOS)" ] || [ -z "$(GOARCH)" ]; then \
		echo "Error: GOOS and GOARCH must be set"; \
		exit 1; \
	fi
	@BINARY_NAME=kgym-sso-service-$(GOOS)-$(GOARCH)$(if $(filter windows,$(GOOS)),.exe,) && \
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $$BINARY_NAME ./cmd/main.go

.PHONY: test
test:
	gotestsum --format short-verbose -- ./internal/... ./migrations/...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: docker-build
docker-build:
	docker build -t kgym-sso-service .

.PHONY: generate-mocks
generate-mocks:
	mockgen -source=internal/repository/token/repository.go -destination=internal/repository/token/mocks/repository_mock.go -package=mocks
	mockgen -source=internal/repository/key/repository.go -destination=internal/repository/key/mocks/repository_mock.go -package=mocks
	mockgen -source=internal/repository/user/repository.go -destination=internal/repository/user/mocks/repository_mock.go -package=mocks
