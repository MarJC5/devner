# Devner Makefile — convenience targets around `go build` and `goreleaser`.
# The runtime itself uses docker compose directly; this Makefile is only for
# building, testing, and releasing the Go binary.

BIN          := devner
OUT          := bin/$(BIN)
PKG          := ./cmd/devner
VERSION      ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT       := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE         := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS      := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)
INSTALL_DIR  ?= /usr/local/bin

GOOS_LIST    := darwin linux windows
GOARCH_LIST  := amd64 arm64

.PHONY: help build build-menubar build-all-bins run test vet tidy fmt install uninstall clean build-all release-snapshot release gui-dev gui-build gui-clean

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Devner — make targets\n\nUsage: make <target>\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the CLI binary into ./bin/devner
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(OUT) $(PKG)
	@echo "✓ $(OUT) ($(VERSION))"

build-menubar: ## Build the macOS menubar daemon into ./bin/devner-menubar (native only — CGO required for Cocoa/systray)
	@mkdir -p bin
	go build -ldflags="-s -w" -o bin/devner-menubar ./cmd/devner-menubar
	@echo "✓ bin/devner-menubar"

build-all-bins: build build-menubar ## Build CLI + menubar side by side

run: build ## Build and run the TUI
	$(OUT) tui

test: ## Run the full test suite
	go test ./...

vet: ## go vet
	go vet ./...

tidy: ## go mod tidy
	go mod tidy

fmt: ## gofmt -w
	gofmt -w .

install: build ## Build and copy the binary into $(INSTALL_DIR) (may need sudo)
	install -m 0755 $(OUT) $(INSTALL_DIR)/$(BIN)
	@echo "✓ installed to $(INSTALL_DIR)/$(BIN)"

uninstall: ## Remove the installed binary from $(INSTALL_DIR)
	rm -f $(INSTALL_DIR)/$(BIN)

clean: ## Remove build artifacts
	rm -rf bin dist

build-all: ## Cross-compile for darwin/linux/windows × amd64/arm64
	@mkdir -p bin
	@for os in $(GOOS_LIST); do \
		for arch in $(GOARCH_LIST); do \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then continue; fi; \
			ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
			out="bin/$(BIN)-$$os-$$arch$$ext"; \
			echo "→ $$out"; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $$out $(PKG) || exit 1; \
		done; \
	done
	@echo "✓ cross-builds done"

release-snapshot: ## Run GoReleaser in snapshot mode (no tag, no push)
	goreleaser release --snapshot --clean

release: ## Run GoReleaser for real (requires a git tag and GITHUB_TOKEN)
	goreleaser release --clean

# ---- GUI (Wails + React) ------------------------------------------------
# Wails CLI is expected on PATH. If missing:
#   go install github.com/wailsapp/wails/v2/cmd/wails@latest
# And ensure $(go env GOPATH)/bin is on PATH.

gui-dev: ## Launch the Wails GUI in dev mode (hot reload Vite + Go rebuild on save)
	cd gui && PATH="$$(go env GOPATH)/bin:$$PATH" wails dev

gui-build: ## Build the Wails GUI production bundle (gui/build/bin/devner-gui(.app/.exe/…))
	cd gui && PATH="$$(go env GOPATH)/bin:$$PATH" wails build

gui-clean: ## Remove the GUI build artifacts
	rm -rf gui/build/bin gui/frontend/dist gui/frontend/wailsjs
