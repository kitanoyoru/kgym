.PHONY: help user-% gateway-% file-% sso-% tools tools-install tools-update ci-test test-all

SERVICES := gateway sso user file

TOOLS_DIR := $(CURDIR)/tools
BIN_DIR := $(CURDIR)/bin

contracts-protobuf-gen-go:
	@$(MAKE) -C contracts/protobuf generate

contracts-deps:
	@$(MAKE) -C contracts/protobuf deps

tools-install: | tools-update
	@mkdir -p $(BIN_DIR)
	@echo "Installing tools..."
	@echo "Building mockgen..."
	@git submodule update --init --recursive tools/mockgen
	@cd $(TOOLS_DIR)/mockgen && GOWORK=off go build -o $(BIN_DIR)/mockgen ./mockgen
	@echo "Building gotestsum..."
	@git submodule update --init --recursive tools/gotestsum
	@cd $(TOOLS_DIR)/gotestsum && GOWORK=off go build -o $(BIN_DIR)/gotestsum -ldflags="-s -w" .
	@chmod +x $(BIN_DIR)/gotestsum 2>/dev/null || true
	@echo "Building golangci-lint..."
	@git submodule update --init --recursive tools/golangci-lint
	@cd $(TOOLS_DIR)/golangci-lint && GOWORK=off go build -o $(BIN_DIR)/golangci-lint ./cmd/golangci-lint
	@echo "Building protoc-gen-go..."
	@git submodule update --init --recursive tools/protobuf-go
	@cd $(TOOLS_DIR)/protobuf-go/cmd/protoc-gen-go && GOWORK=off go build -o $(BIN_DIR)/protoc-gen-go .
	@chmod +x $(BIN_DIR)/protoc-gen-go 2>/dev/null || true
	@echo "Building protoc-gen-go-grpc..."
	@git submodule update --init --recursive tools/grpc-go
	@cd $(TOOLS_DIR)/grpc-go/cmd/protoc-gen-go-grpc && GOWORK=off go build -o $(BIN_DIR)/protoc-gen-go-grpc .
	@chmod +x $(BIN_DIR)/protoc-gen-go-grpc 2>/dev/null || true
	@echo "Building protoc-gen-grpc-gateway..."
	@git submodule update --init --recursive tools/grpc-gateway
	@cd $(TOOLS_DIR)/grpc-gateway/protoc-gen-grpc-gateway && GOWORK=off go build -o $(BIN_DIR)/protoc-gen-grpc-gateway .
	@chmod +x $(BIN_DIR)/protoc-gen-grpc-gateway 2>/dev/null || true
	@echo "Building protoc-gen-validate..."
	@git submodule update --init --recursive tools/protoc-gen-validate
	@cd $(TOOLS_DIR)/protoc-gen-validate && GOWORK=off go build -o $(BIN_DIR)/protoc-gen-validate .
	@chmod +x $(BIN_DIR)/protoc-gen-validate 2>/dev/null || true
	@echo "All tools installed successfully in $(BIN_DIR)"
	@echo "Add $(BIN_DIR) to your PATH to use them"

gomod-all:
	@echo "Running gomod tidy on all services..."
	@for service in $(SERVICES); do \
		echo ""; \
		echo "=== Running gomod on $$service service ==="; \
		$(MAKE) $$service-gomod || exit 1; \
	done
	@echo ""
	@echo "All services gomod completed successfully!"


lint-all:
	@echo "Running all linters..."
	@for service in $(SERVICES); do \
		echo ""; \
		echo "=== Linting $$service service ==="; \
		$(MAKE) $$service-lint || exit 1; \
	done
	@echo ""
	@echo "All linters completed successfully!"

build-all:
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo ""; \
		echo "=== Building $$service service ==="; \
		$(MAKE) $$service-build || exit 1; \
	done
	@echo ""
	@echo "All services built successfully!"

test-all:
	@echo "Running all tests..."
	@for service in $(SERVICES); do \
		echo ""; \
		echo "=== Testing $$service service ==="; \
		$(MAKE) $$service-test || exit 1; \
	done
	@echo ""
	@echo "All tests completed successfully!"

ci-test:
	@echo "Testing CI pipeline with act..."
	@act push

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
