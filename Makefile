# Description: Makefile for generating protobuf code

# Directories
PROTO_DIR := libs/schema/pkg
BASE_DIR := .

# Get all directories except core
DEPENDENCY_DIRS := $(shell find $(PROTO_DIR) -type d -not -path "$(PROTO_DIR)/core*" -not -path "$(PROTO_DIR)")
DEPENDENCY_PROTOS := $(shell find $(DEPENDENCY_DIRS) -name "*.proto")

# Core protos separately
CORE_PROTOS := $(shell find $(PROTO_DIR)/core -name "*.proto")

.PHONY: all proto proto-dependencies proto-core clean list help gazelle

all: help

# Generate protobuf code for all dependency directories
proto: proto-dependencies proto-core

proto-dependencies:
	@for dir in $(DEPENDENCY_DIRS); do \
		echo "Generating protobuf code for $${dir}..."; \
		PROTOS=$$(find $${dir} -name "*.proto"); \
		if [ ! -z "$${PROTOS}" ]; then \
			protoc --go_out=$(BASE_DIR) \
				--go_opt=paths=source_relative \
				-I$(BASE_DIR) \
				$${PROTOS}; \
		fi; \
	done

# Generate protobuf code for core files last
proto-core:
	@echo "Generating protobuf code for libs/schema/pkg/core..."
	@if [ ! -z "$(CORE_PROTOS)" ]; then \
		protoc --go_out=$(BASE_DIR) \
			--go_opt=paths=source_relative \
			-I$(BASE_DIR) \
			$(CORE_PROTOS); \
	fi

# Clean generated protobuf files
clean:
	@echo "Cleaning generated protobuf files..."
	@find $(PROTO_DIR) -name "*.pb.go" -delete

# Show all proto files that will be processed
list:
	@echo "Dependency directories:"
	@echo "$(DEPENDENCY_DIRS)" | tr ' ' '\n'
	@echo "\nDependency proto files:"
	@echo "$(DEPENDENCY_PROTOS)" | tr ' ' '\n'
	@echo "\nCore proto files:"
	@echo "$(CORE_PROTOS)" | tr ' ' '\n'

gazelle:
	bazel run //:gazelle

help:
	@echo "Available targets:"
	@echo "  proto              - Generate all protobuf code (dependencies first, then core)"
	@echo "  proto-dependencies - Generate protobuf code for all dependency directories"
	@echo "  proto-core         - Generate only core protobuf code"
	@echo "  clean              - Remove all generated protobuf files"
	@echo "  list               - List all proto files that will be processed"
	@echo "  gazelle 			- Run gazelle using bazel"

