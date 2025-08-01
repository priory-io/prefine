BINARY_NAME := prefine
VERSION := $(shell git tag --list --sort=-version:refname | head -n1 2>/dev/null || git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

GOVERSION := $(shell go version | cut -d ' ' -f 3)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

BUILD_DIR := ./build
DIST_DIR := $(BUILD_DIR)/dist
CMD_DIR := ./cmd/prefine
TEST_IMAGES_DIR := ./test_images

LDFLAGS := -w -s \
	-X 'github.com/priory-io/prefine/cmd.Version=$(VERSION)' \
	-X 'github.com/priory-io/prefine/cmd.BuildTime=$(BUILD_TIME)' \
	-X 'github.com/priory-io/prefine/cmd.GitCommit=$(GIT_COMMIT)' \
	-X 'github.com/priory-io/prefine/cmd.GoVersion=$(GOVERSION)'

GCFLAGS := -B -C
CGO_ENABLED := 0

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

TEST_IMAGE_URLS := \
	https://picsum.photos/1920/1080.jpg \
	https://picsum.photos/800/600.png \
	https://picsum.photos/1280/720.jpg \
	https://picsum.photos/640/480.png \
	https://picsum.photos/1600/900.jpg \
	https://picsum.photos/1024/768.png

.PHONY: all build clean release test-images help deps lint fmt vet

all: build

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

build-dev:
	@echo "Building $(BINARY_NAME) for development..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build \
		-race \
		-o $(BUILD_DIR)/$(BINARY_NAME)-dev \
		$(CMD_DIR)
	@echo "Development build complete: $(BUILD_DIR)/$(BINARY_NAME)-dev"

release: clean
	@echo "Building release binaries for $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		output_name=$(BINARY_NAME); \
		if [ "$$os" = "windows" ]; then \
			output_name=$(BINARY_NAME).exe; \
		fi; \
		echo "Building for $$os/$$arch..."; \
		if GOOS=$$os GOARCH=$$arch CGO_ENABLED=$(CGO_ENABLED) go build \
			-ldflags="$(LDFLAGS)" \
			-gcflags="$(GCFLAGS)" \
			-trimpath \
			-o $(BUILD_DIR)/$$output_name \
			$(CMD_DIR) 2>/dev/null; then \
			archive_name=$(BINARY_NAME)_$(VERSION)_$${os}_$${arch}; \
			if [ "$$os" = "windows" ]; then \
				(cd $(BUILD_DIR) && zip -q $${archive_name}.zip $$output_name && mv $${archive_name}.zip dist/); \
			else \
				(cd $(BUILD_DIR) && tar -czf $${archive_name}.tar.gz $$output_name && mv $${archive_name}.tar.gz dist/); \
			fi; \
			rm -f $(BUILD_DIR)/$$output_name; \
			echo "✓ Created: $${archive_name}"; \
		else \
			echo "✗ Failed to build for $$os/$$arch"; \
		fi; \
	done
	@echo ""
	@echo "Release builds complete in $(DIST_DIR)/"
	@ls -la $(DIST_DIR)/

