# Makefile for render CLI
# Supports incremental builds, cross-compilation for 40 platforms, and packaging

# ============================================================================
# Configuration
# ============================================================================

BINARY_NAME := render
CMD_PATH := ./cmd/render
BIN_DIR := bin
DIST_DIR := dist
TOOLS_DIR := tools

# Version from git tags, fallback to "dev"
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Build flags
LDFLAGS := -s -w -X main.version=$(VERSION)
GOFLAGS := -ldflags="$(LDFLAGS)"

# Find all Go source files for dependency tracking
GO_FILES := $(shell find . -type f -name '*.go' -not -path './.scratchpad/*' -not -path './tools/*')
BUILD_DEPS := go.mod go.sum Makefile $(GO_FILES)

# Package tool
PACKAGE_TOOL := $(BIN_DIR)/tools/package

# ============================================================================
# Platform Definitions (40 CLI-capable combinations)
# ============================================================================

# Linux platforms (13)
LINUX_ARCHS := 386 amd64 arm arm64 loong64 mips mips64 mips64le mipsle ppc64 ppc64le riscv64 s390x
LINUX_TARGETS := $(addprefix linux/,$(LINUX_ARCHS))

# Darwin platforms (2)
DARWIN_ARCHS := amd64 arm64
DARWIN_TARGETS := $(addprefix darwin/,$(DARWIN_ARCHS))

# Windows platforms (3)
WINDOWS_ARCHS := 386 amd64 arm64
WINDOWS_TARGETS := $(addprefix windows/,$(WINDOWS_ARCHS))

# FreeBSD platforms (5)
FREEBSD_ARCHS := 386 amd64 arm arm64 riscv64
FREEBSD_TARGETS := $(addprefix freebsd/,$(FREEBSD_ARCHS))

# OpenBSD platforms (6)
OPENBSD_ARCHS := 386 amd64 arm arm64 ppc64 riscv64
OPENBSD_TARGETS := $(addprefix openbsd/,$(OPENBSD_ARCHS))

# NetBSD platforms (4)
NETBSD_ARCHS := 386 amd64 arm arm64
NETBSD_TARGETS := $(addprefix netbsd/,$(NETBSD_ARCHS))

# Dragonfly platforms (1)
DRAGONFLY_ARCHS := amd64
DRAGONFLY_TARGETS := $(addprefix dragonfly/,$(DRAGONFLY_ARCHS))

# Solaris platforms (1)
SOLARIS_ARCHS := amd64
SOLARIS_TARGETS := $(addprefix solaris/,$(SOLARIS_ARCHS))

# Illumos platforms (1)
ILLUMOS_ARCHS := amd64
ILLUMOS_TARGETS := $(addprefix illumos/,$(ILLUMOS_ARCHS))

# AIX platforms (1)
AIX_ARCHS := ppc64
AIX_TARGETS := $(addprefix aix/,$(AIX_ARCHS))

# Plan9 platforms (3)
PLAN9_ARCHS := 386 amd64 arm
PLAN9_TARGETS := $(addprefix plan9/,$(PLAN9_ARCHS))

# All BSD targets
BSD_TARGETS := $(FREEBSD_TARGETS) $(OPENBSD_TARGETS) $(NETBSD_TARGETS) $(DRAGONFLY_TARGETS)

# All Unix-like targets (excluding Windows)
UNIX_TARGETS := $(SOLARIS_TARGETS) $(ILLUMOS_TARGETS) $(AIX_TARGETS) $(PLAN9_TARGETS)

# All Unix-like targets (for tar.gz packaging)
ALL_UNIX_TARGETS := $(LINUX_TARGETS) $(DARWIN_TARGETS) $(BSD_TARGETS) $(UNIX_TARGETS)

# All targets
ALL_TARGETS := $(ALL_UNIX_TARGETS) $(WINDOWS_TARGETS)

# Convert targets to binary paths
UNIX_BINARIES := $(addprefix $(BIN_DIR)/,$(addsuffix /$(BINARY_NAME),$(ALL_UNIX_TARGETS)))
WINDOWS_BINARIES := $(addprefix $(BIN_DIR)/,$(addsuffix /$(BINARY_NAME).exe,$(WINDOWS_TARGETS)))
ALL_BINARIES := $(UNIX_BINARIES) $(WINDOWS_BINARIES)

# ============================================================================
# Default Target
# ============================================================================

.PHONY: all
all: build

# ============================================================================
# Help
# ============================================================================

.PHONY: help
help:
	@echo "render Makefile - Build and package for 40 platforms"
	@echo ""
	@echo "Primary targets:"
	@echo "  make build          Build for current OS/arch → bin/"
	@echo "  make build-all      Build all 40 platforms (use -j for parallel)"
	@echo "  make test           Run all tests (unit + acceptance)"
	@echo "  make lint           Run golangci-lint"
	@echo "  make clean          Remove bin/ and dist/"
	@echo ""
	@echo "Platform-specific builds:"
	@echo "  make build-linux    All Linux variants (13 archs)"
	@echo "  make build-darwin   All macOS variants (2 archs)"
	@echo "  make build-windows  All Windows variants (3 archs)"
	@echo "  make build-bsd      FreeBSD, OpenBSD, NetBSD, Dragonfly (16 archs)"
	@echo "  make build-unix     Solaris, Illumos, AIX, Plan9 (6 archs)"
	@echo ""
	@echo "Packaging:"
	@echo "  make package        Create all packages"
	@echo "  make package-tar    tar.gz for Unix systems"
	@echo "  make package-zip    zip for Windows"
	@echo "  make package-deb    .deb for Debian/Ubuntu"
	@echo "  make checksums      Generate SHA256 checksums"
	@echo ""
	@echo "CI/CD helpers:"
	@echo "  make verify         go mod verify"
	@echo "  make download       go mod download"
	@echo "  make release        build-all + package + checksums"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION=$(VERSION)"
	@echo ""
	@echo "Examples:"
	@echo "  make build-all -j8           # Parallel build all platforms"
	@echo "  make release VERSION=v1.0.0  # Full release build"

# ============================================================================
# Tools
# ============================================================================

$(PACKAGE_TOOL): $(TOOLS_DIR)/package/main.go
	@mkdir -p $(dir $@)
	@echo "Building package tool"
	@go build -o $@ ./$(TOOLS_DIR)/package

# ============================================================================
# Build Targets
# ============================================================================

# Build for current OS/arch
.PHONY: build
build: $(BIN_DIR)/$(shell go env GOOS)/$(shell go env GOARCH)/$(BINARY_NAME)$(shell go env GOEXE)

# Build all platforms
.PHONY: build-all
build-all: $(ALL_BINARIES)

# Platform group targets
.PHONY: build-linux
build-linux: $(addprefix $(BIN_DIR)/,$(addsuffix /$(BINARY_NAME),$(LINUX_TARGETS)))

.PHONY: build-darwin
build-darwin: $(addprefix $(BIN_DIR)/,$(addsuffix /$(BINARY_NAME),$(DARWIN_TARGETS)))

.PHONY: build-windows
build-windows: $(addprefix $(BIN_DIR)/,$(addsuffix /$(BINARY_NAME).exe,$(WINDOWS_TARGETS)))

.PHONY: build-bsd
build-bsd: $(addprefix $(BIN_DIR)/,$(addsuffix /$(BINARY_NAME),$(BSD_TARGETS)))

.PHONY: build-unix
build-unix: $(addprefix $(BIN_DIR)/,$(addsuffix /$(BINARY_NAME),$(UNIX_TARGETS)))

# ============================================================================
# Pattern Rules for Building
# ============================================================================

# Generic Unix binary build rule
# Matches: bin/linux/amd64/render, bin/darwin/arm64/render, etc.
$(BIN_DIR)/%/$(BINARY_NAME): $(BUILD_DEPS)
	@mkdir -p $(dir $@)
	@echo "Building $@"
	@GOOS=$(word 1,$(subst /, ,$*)) GOARCH=$(word 2,$(subst /, ,$*)) \
		go build $(GOFLAGS) -o $@ $(CMD_PATH)

# Windows binary build rule (adds .exe extension)
# Matches: bin/windows/amd64/render.exe, etc.
$(BIN_DIR)/windows/%/$(BINARY_NAME).exe: $(BUILD_DEPS)
	@mkdir -p $(dir $@)
	@echo "Building $@"
	@GOOS=windows GOARCH=$* go build $(GOFLAGS) -o $@ $(CMD_PATH)

# ============================================================================
# Packaging Targets
# ============================================================================

.PHONY: package
package: package-tar package-zip package-deb

# Generate tar.gz for all Unix targets
.PHONY: package-tar
package-tar: $(PACKAGE_TOOL) $(UNIX_BINARIES)
	@mkdir -p $(DIST_DIR)
	@for target in $(ALL_UNIX_TARGETS); do \
		os=$$(echo $$target | cut -d/ -f1); \
		arch=$$(echo $$target | cut -d/ -f2); \
		$(PACKAGE_TOOL) tar \
			-binary "$(BIN_DIR)/$$target/$(BINARY_NAME)" \
			-output "$(DIST_DIR)/$(BINARY_NAME)-$${os}-$${arch}-$(VERSION).tar.gz" \
			-name "$(BINARY_NAME)"; \
	done

# Generate zip for all Windows targets
.PHONY: package-zip
package-zip: $(PACKAGE_TOOL) $(WINDOWS_BINARIES)
	@mkdir -p $(DIST_DIR)
	@for arch in $(WINDOWS_ARCHS); do \
		$(PACKAGE_TOOL) zip \
			-binary "$(BIN_DIR)/windows/$$arch/$(BINARY_NAME).exe" \
			-output "$(DIST_DIR)/$(BINARY_NAME)-windows-$${arch}-$(VERSION).zip" \
			-name "$(BINARY_NAME).exe"; \
	done

# Generate DEB packages for Linux amd64/arm64
.PHONY: package-deb
package-deb: $(PACKAGE_TOOL) $(BIN_DIR)/linux/amd64/$(BINARY_NAME) $(BIN_DIR)/linux/arm64/$(BINARY_NAME)
	@mkdir -p $(DIST_DIR)
	@for arch in amd64 arm64; do \
		$(PACKAGE_TOOL) deb \
			-binary "$(BIN_DIR)/linux/$$arch/$(BINARY_NAME)" \
			-output "$(DIST_DIR)/$(BINARY_NAME)_$(VERSION)_$${arch}.deb" \
			-name "$(BINARY_NAME)" \
			-version "$(VERSION)" \
			-arch "$$arch" \
			-maintainer "Werner Strydom <hello@wernerstrydom.com>" \
			-description "Template rendering CLI tool"; \
	done

# ============================================================================
# Checksums
# ============================================================================

.PHONY: checksums
checksums:
	@echo "Generating checksums"
	@cd $(DIST_DIR) && \
		if command -v sha256sum >/dev/null 2>&1; then \
			sha256sum *.tar.gz *.zip *.deb 2>/dev/null | sort > checksums.txt; \
		else \
			shasum -a 256 *.tar.gz *.zip *.deb 2>/dev/null | sort > checksums.txt; \
		fi

# ============================================================================
# Testing and Linting
# ============================================================================

.PHONY: test
test:
	go test -v -race ./internal/...
	go test -v -race ./test/acceptance/...

.PHONY: lint
lint:
	golangci-lint run --timeout=5m

# ============================================================================
# CI/CD Helpers
# ============================================================================

.PHONY: verify
verify:
	go mod verify

.PHONY: download
download:
	go mod download

.PHONY: release
release: build-all package checksums

# ============================================================================
# Cleanup
# ============================================================================

.PHONY: clean
clean:
	rm -rf $(BIN_DIR) $(DIST_DIR)

# ============================================================================
# Debug Targets
# ============================================================================

.PHONY: list-targets
list-targets:
	@echo "All targets ($(words $(ALL_TARGETS))):"
	@for t in $(ALL_TARGETS); do echo "  $$t"; done

.PHONY: list-binaries
list-binaries:
	@echo "All binaries ($(words $(ALL_BINARIES))):"
	@for b in $(ALL_BINARIES); do echo "  $$b"; done
