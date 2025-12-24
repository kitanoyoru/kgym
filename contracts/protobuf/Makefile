.PHONY: generate clean deps

PROTOC := protoc

PROTO_DIR := .

OUT_DIR := gen/go

GOOGLEAPIS_DIR := deps/googleapis
PROTOBUF_DIR := deps/protobuf/src
VALIDATE_DIR := deps/protoc-gen-validate
GRPC_GATEWAY_DIR := deps/grpc-gateway

PROTO_FILES := $(shell find . -name "*.proto" -not -path "./deps/*" -not -path "./gen/*")

generate: deps
	@mkdir -p $(OUT_DIR)
	@echo "Generating protobuf files..."
	$(PROTOC) \
		-I=$(PROTO_DIR) \
		-I=$(GOOGLEAPIS_DIR) \
		-I=$(PROTOBUF_DIR) \
		-I=$(VALIDATE_DIR) \
		-I=$(GRPC_GATEWAY_DIR) \
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
	@echo "Initializing git submodules..."
	@cd "$(shell git rev-parse --show-toplevel)" && \
		git submodule update --init --recursive -- contracts/protobuf/deps/googleapis \
		contracts/protobuf/deps/protobuf \
		contracts/protobuf/deps/protoc-gen-validate \
		contracts/protobuf/deps/grpc-gateway

clean:
	@rm -rf $(OUT_DIR)

