.PHONY: help user-% gym-% subscription-%

SERVICES := user gym subscription

help:
	@echo "Available commands:"
	@echo "  make <service>-<target>  - Run target in specific service"
	@echo "  make contracts-generate  - Generate protobuf files"
	@echo ""
	@echo "Examples:"
	@echo "  make user-test          - Run tests in user service"
	@echo "  make gym-build          - Build gym service"
	@echo "  make contracts-generate - Generate protobuf code"
	@echo ""
	@echo "Available services: $(SERVICES)"

contracts-generate:
	@$(MAKE) -C contracts/protobuf generate

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
