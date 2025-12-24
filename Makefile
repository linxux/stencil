# Define the Go compiler
GO = go
GOFLAGS = -ldflags="-s -w"

# Define the build directory
BUILD_DIR = ./bin

# Define the binary name
BINARY = stencil

# Define the source files
SOURCES = ./cmd/stencil/main.go

.PHONY: all init build dev run clean update-deps test

# Default target
all: build

# Init project convertional-commits
init:
	chmod +x .bin/conventional-commits/setup.sh && .bin/conventional-commits/setup.sh

# Version
generate-version:
	chmod +x .bin/version/generate-version-info.sh && .bin/version/generate-version-info.sh

# Build the binary
build:
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY) $(SOURCES)

# Development target
dev:
	$(GO) run $(SOURCES)

# Run the binary
run:
	$(BUILD_DIR)/$(BINARY)

# Clean the build directory
clean:
	rm -rf $(BUILD_DIR)

# Install the dependencies
update-deps:
	$(GO) get -u all
	$(GO) mod tidy

# Test the project
test:
	$(GO) test ./...
