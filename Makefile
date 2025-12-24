.PHONY: help user-% gym-% subscription-% tools tools-install tools-update test-ci

SERVICES := user gym subscription

TOOLS_DIR := $(CURDIR)/tools
BIN_DIR := $(CURDIR)/bin

help:
	@echo "Available commands:"
	@echo "  make <service>-<target>  - Run target in specific service"
	@echo "  make contracts-generate  - Generate protobuf files"
	@echo "  make tools-install       - Install tools from submodules"
	@echo "  make tools-update        - Update tool submodules"
	@echo "  make test-ci             - Test CI pipeline locally with act"
	@echo ""
	@echo "Examples:"
	@echo "  make user-test          - Run tests in user service"
	@echo "  make gym-build          - Build gym service"
	@echo "  make contracts-generate  - Generate protobuf code"
	@echo "  make tools-install      - Build and install all tools"
	@echo "  make test-ci             - Test CI pipeline locally"
	@echo ""
	@echo "Available services: $(SERVICES)"

contracts-generate:
	@$(MAKE) -C contracts/protobuf generate

tools-update:
	@echo "Updating tool submodules..."
	@git submodule update --init --recursive

tools-install: $(BIN_DIR)/mockgen $(BIN_DIR)/gotestsum $(BIN_DIR)/golangci-lint
	@echo "All tools installed successfully in $(BIN_DIR)"
	@echo "Add $(BIN_DIR) to your PATH to use them"

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

$(BIN_DIR)/mockgen: $(BIN_DIR) | tools-update
	@echo "Building mockgen..."
	@cd $(TOOLS_DIR)/mockgen && GOWORK=off go build -o $(BIN_DIR)/mockgen ./mockgen

$(BIN_DIR)/gotestsum: $(BIN_DIR) | tools-update
	@echo "Building gotestsum..."
	@GOWORK=off GOBIN=$(BIN_DIR) go install -ldflags="-s -w" gotest.tools/gotestsum@v1.13.0
	@chmod +x $(BIN_DIR)/gotestsum 2>/dev/null || true

$(BIN_DIR)/golangci-lint: $(BIN_DIR) | tools-update
	@echo "Building golangci-lint..."
	@cd $(TOOLS_DIR)/golangci-lint && GOWORK=off go build -o $(BIN_DIR)/golangci-lint ./cmd/golangci-lint

tools: tools-install

test-ci:
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
