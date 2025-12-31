.PHONY: help user-% gateway-% file-% sso-% tools tools-install tools-update ci-test test-all

SERVICES := user gateway file sso

TOOLS_DIR := $(CURDIR)/tools
BIN_DIR := $(CURDIR)/bin

help:
	@echo "Available commands:"
	@echo "  make <service>-<target>  - Run target in specific service"
	@echo "  make test-all            - Run all tests across all services"
	@echo "  make contracts-generate  - Generate protobuf files"
	@echo "  make contracts-deps      - Update contracts deps submodules"
	@echo "  make tools-install       - Install tools from submodules"
	@echo "  make ci-test             - Test CI pipeline locally with act"
	@echo ""
	@echo "Examples:"
	@echo "  make user-test          - Run tests in user service"
	@echo "  make file-build          - Build file service"
	@echo "  make test-all          - Run all tests"
	@echo ""
	@echo "Available services: $(SERVICES)"

contracts-protobuf-gen-go:
	@$(MAKE) -C contracts/protobuf generate

contracts-deps:
	@$(MAKE) -C contracts/protobuf deps

tools-install: $(BIN_DIR)/mockgen $(BIN_DIR)/gotestsum $(BIN_DIR)/golangci-lint
	@echo "All tools installed successfully in $(BIN_DIR)"
	@echo "Add $(BIN_DIR) to your PATH to use them"

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

$(TOOLS_DIR)/mockgen:
	@git submodule update --init --recursive tools/mockgen

$(BIN_DIR)/mockgen: $(BIN_DIR) $(TOOLS_DIR)/mockgen | tools-update
	@echo "Building mockgen..."
	@cd $(TOOLS_DIR)/mockgen && GOWORK=off go build -o $(BIN_DIR)/mockgen ./mockgen

$(BIN_DIR)/gotestsum: $(BIN_DIR) | tools-update
	@echo "Building gotestsum..."
	@GOWORK=off GOBIN=$(BIN_DIR) go install -ldflags="-s -w" gotest.tools/gotestsum@v1.13.0
	@chmod +x $(BIN_DIR)/gotestsum 2>/dev/null || true

$(TOOLS_DIR)/golangci-lint:
	@git submodule update --init --recursive tools/golangci-lint

$(BIN_DIR)/golangci-lint: $(BIN_DIR) $(TOOLS_DIR)/golangci-lint | tools-update
	@echo "Building golangci-lint..."
	@cd $(TOOLS_DIR)/golangci-lint && GOWORK=off go build -o $(BIN_DIR)/golangci-lint ./cmd/golangci-lint

ci-test:
	@echo "Testing CI pipeline with act..."
	@act push

test-all:
	@echo "Running all tests..."
	@for service in $(SERVICES); do \
		echo ""; \
		echo "=== Testing $$service service ==="; \
		$(MAKE) $$service-test || exit 1; \
	done
	@echo ""
	@echo "All tests completed successfully!"

lint-all:
	@echo "Running all linters..."
	@for service in $(SERVICES); do \
		echo ""; \
		echo "=== Linting $$service service ==="; \
		$(MAKE) $$service-lint || exit 1; \
	done
	@echo ""
	@echo "All linters completed successfully!"

define SERVICE_TARGET
$(1)-%:
	@if [ "$(1)" = "gateway" ]; then \
		if [ -f internal/gateway/Makefile ]; then \
			$(MAKE) -C internal/gateway $$*; \
		else \
			echo "Error: Makefile not found for service $(1)"; \
			exit 1; \
		fi; \
	elif [ -f internal/apps/$(1)/Makefile ]; then \
		$(MAKE) -C internal/apps/$(1) $$*; \
	else \
		echo "Error: Makefile not found for service $(1)"; \
		exit 1; \
	fi
endef

$(foreach service,$(SERVICES),$(eval $(call SERVICE_TARGET,$(service))))
