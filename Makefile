# Makefile for prefine
# Build variables
BINARY_NAME := prefine
VERSION := $(shell git describe --tags --exact-match 2>/dev/null || git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go variables
GOVERSION := $(shell go version | cut -d ' ' -f 3)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Directories
BUILD_DIR := ./build
DIST_DIR := $(BUILD_DIR)/dist
CMD_DIR := ./cmd/prefine
TEST_IMAGES_DIR := ./test_images

# Compiler flags
LDFLAGS := -w -s \
	-X 'github.com/priory-io/prefine/cmd.Version=$(VERSION)' \
	-X 'github.com/priory-io/prefine/cmd.BuildTime=$(BUILD_TIME)' \
	-X 'github.com/priory-io/prefine/cmd.GitCommit=$(GIT_COMMIT)' \
	-X 'github.com/priory-io/prefine/cmd.GoVersion=$(GOVERSION)'

GCFLAGS := -B -C
CGO_ENABLED := 0

# Release platforms
PLATFORMS := \
	darwin/amd64 \
	darwin/arm64 \
	linux/amd64 \
	linux/arm64 \
	linux/386 \
	windows/amd64 \
	windows/386 \
	freebsd/amd64 \
	openbsd/amd64

# Test images URLs
TEST_IMAGE_URLS := \
	https://picsum.photos/1920/1080.jpg \
	https://picsum.photos/800/600.png \
	https://picsum.photos/1280/720.jpg \
	https://picsum.photos/640/480.png \
	https://picsum.photos/1600/900.jpg \
	https://picsum.photos/1024/768.png

.PHONY: all build clean release test-images help deps lint fmt vet

# Default target
all: build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) $(VERSION) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) go build \
		-ldflags="$(LDFLAGS)" \
		-gcflags="$(GCFLAGS)" \
		-trimpath \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for development (no optimizations, includes debug info)
build-dev:
	@echo "Building $(BINARY_NAME) for development..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build \
		-race \
		-o $(BUILD_DIR)/$(BINARY_NAME)-dev \
		$(CMD_DIR)
	@echo "Development build complete: $(BUILD_DIR)/$(BINARY_NAME)-dev"

# Cross-compile for all platforms
release: clean
	@echo "Building release binaries for $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(DIST_DIR)
	@echo "Building for linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build \
		-ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -trimpath \
		-o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR) && \
		cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz $(BINARY_NAME) && \
		mv $(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz dist/ && rm $(BINARY_NAME)
	@echo "Building for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build \
		-ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -trimpath \
		-o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR) && \
		cd $(BUILD_DIR) && tar -czf $(BINARY_NAME)_$(VERSION)_darwin_amd64.tar.gz $(BINARY_NAME) && \
		mv $(BINARY_NAME)_$(VERSION)_darwin_amd64.tar.gz dist/ && rm $(BINARY_NAME)
	@echo "Building for windows/amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build \
		-ldflags="$(LDFLAGS)" -gcflags="$(GCFLAGS)" -trimpath \
		-o $(BUILD_DIR)/$(BINARY_NAME).exe $(CMD_DIR) && \
		cd $(BUILD_DIR) && zip -q $(BINARY_NAME)_$(VERSION)_windows_amd64.zip $(BINARY_NAME).exe && \
		mv $(BINARY_NAME)_$(VERSION)_windows_amd64.zip dist/ && rm $(BINARY_NAME).exe
	@echo "Release builds complete in $(DIST_DIR)/"
	@ls -la $(DIST_DIR)/

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build for current platform"
	@echo "  build-dev    - Build development version (with race detector)"
	@echo "  release      - Cross-compile for all platforms and create archives"
	@echo "  test-images  - Download test images for testing"
	@echo "  clean        - Remove build artifacts"
	@echo "  clean-test-images - Remove test images"
	@echo "  deps         - Install/update dependencies"
	@echo "  lint         - Run linter (requires golangci-lint)"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  test         - Run tests"
	@echo "  bench        - Run benchmarks"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION      = $(VERSION)"
	@echo "  BUILD_TIME   = $(BUILD_TIME)"
	@echo "  GIT_COMMIT   = $(GIT_COMMIT)"
	@echo "  GOOS         = $(GOOS)"
	@echo "  GOARCH       = $(GOARCH)"