.PHONY: generate clean deps

PROTOC := protoc
PROTO_DIR := .
OUT_DIR := gen/go
GOOGLEAPIS_DIR := deps/googleapis
PROTOBUF_DIR := deps/protobuf/src
VALIDATE_DIR := deps/protoc-gen-validate

PROTO_FILES := $(shell find . -name "*.proto" -not -path "./deps/*" -not -path "./gen/*")

generate: deps
	@mkdir -p $(OUT_DIR)
	@echo "Generating protobuf files..."
	$(PROTOC) \
		-I=$(PROTO_DIR) \
		-I=$(GOOGLEAPIS_DIR) \
		-I=$(PROTOBUF_DIR) \
		-I=$(VALIDATE_DIR) \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) \
		--grpc-gateway_opt=paths=source_relative \
		--validate_out=$(OUT_DIR) \
		--validate_opt=paths=source_relative,lang=go \
		$(PROTO_FILES)

deps:
	@if [ ! -d "$(GOOGLEAPIS_DIR)" ]; then \
		echo "Downloading googleapis..."; \
		mkdir -p deps; \
		git clone --depth 1 https://github.com/googleapis/googleapis.git $(GOOGLEAPIS_DIR) || true; \
	fi
	@if [ ! -d "$(PROTOBUF_DIR)" ]; then \
		echo "Downloading protobuf..."; \
		mkdir -p deps; \
		git clone --depth 1 https://github.com/protocolbuffers/protobuf.git deps/protobuf || true; \
	fi
	@if [ ! -d "$(VALIDATE_DIR)" ]; then \
		echo "Downloading protoc-gen-validate..."; \
		mkdir -p deps; \
		git clone --depth 1 https://github.com/bufbuild/protoc-gen-validate.git $(VALIDATE_DIR) || true; \
	fi

clean:
	@rm -rf $(OUT_DIR)
	@rm -rf deps

