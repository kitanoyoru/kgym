.PHONY: help user-% gym-% subscription-% tools tools-install tools-update ci-test

SERVICES := user gym subscription

TOOLS_DIR := $(CURDIR)/tools
BIN_DIR := $(CURDIR)/bin

help:
	@echo "Available commands:"
	@echo "  make <service>-<target>  - Run target in specific service"
	@echo "  make contracts-generate  - Generate protobuf files"
	@echo "  make contracts-deps      - Update contracts deps submodules"
	@echo "  make tools-install       - Install tools from submodules"
	@echo "  make ci-test             - Test CI pipeline locally with act"
	@echo ""
	@echo "Examples:"
	@echo "  make user-test          - Run tests in user service"
	@echo "  make gym-build          - Build gym service"
	@echo ""
	@echo "Available services: $(SERVICES)"

contracts-generate:
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

define SERVICE_TARGET
$(1)-%:
	@if [ -f internal/apps/$(1)/Makefile ]; then \
		$(MAKE) -C internal/apps/$(1) $$*; \
	else \
		echo "Error: Makefile not found for service $(1)"; \
		exit 1; \
	fi
endef

$(foreach service,$(SERVICES),$(eval $(call SERVICE_TARGET,$(service))))
