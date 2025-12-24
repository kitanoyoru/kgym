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
	mockgen -source=internal/repository/postgres/repository.go -destination=internal/repository/postgres/mocks/repository_mock.go -package=mocks
