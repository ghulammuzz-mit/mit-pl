BUILD_DIR := bin

VERSION ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w \
	-X 'main.version=$(VERSION)' \
	-X 'main.commit=$(COMMIT)' \
	-X 'main.buildDate=$(DATE)'

# Infisical credentials baked into envctl at build time. Empty by default —
# CI passes real values from secrets. Local builds without these still compile;
# the binary will print a clear error at runtime instead of authenticating.
INFISICAL_PKG := mit/platform/internal/envctl/infisical
ENVCTL_LDFLAGS := $(LDFLAGS) \
	-X '$(INFISICAL_PKG).embeddedClientID=$(INFISICAL_BUILD_CLIENT_ID)' \
	-X '$(INFISICAL_PKG).embeddedClientSecret=$(INFISICAL_BUILD_CLIENT_SECRET)' \
	-X '$(INFISICAL_PKG).embeddedSiteURL=$(INFISICAL_BUILD_SITE_URL)'

.PHONY: all envctl dns mcp-grafana mcp-infisical hpa-metrics build clean run-envctl run-dns run-mcp-grafana run-mcp-infisical run-hpa-metrics release-envctl

all: build

## Build all binaries
build: envctl dns mcp-grafana mcp-infisical hpa-metrics

## Build envctl
envctl:
	@mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags="$(ENVCTL_LDFLAGS)" -o $(BUILD_DIR)/envctl ./cmd/envctl

## Build dns
dns:
	@mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/dns ./cmd/dns

## Build mcp-grafana
mcp-grafana:
	@mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/mcp-grafana ./cmd/mcp-grafana

## Build hpa-metrics
hpa-metrics:
	@mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/hpa-metrics ./cmd/hpa-metrics

## Build mcp-infisical
mcp-infisical:
	@mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/mcp-infisical ./cmd/mcp-infisical

## Run envctl
run-envctl:
	go run ./cmd/envctl

## Run dns
run-dns:
	go run ./cmd/dns

## Run mcp-grafana
run-mcp-grafana:
	go run ./cmd/mcp-grafana

## Run hpa-metrics
run-hpa-metrics:
	go run ./cmd/hpa-metrics

## Run mcp-infisical
run-mcp-infisical:
	go run ./cmd/mcp-infisical

## Build envctl for all platforms (output: dist/)
release-envctl:
	@mkdir -p dist
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="$(ENVCTL_LDFLAGS)" -o dist/envctl-linux-amd64   ./cmd/envctl
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="$(ENVCTL_LDFLAGS)" -o dist/envctl-linux-arm64   ./cmd/envctl
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="$(ENVCTL_LDFLAGS)" -o dist/envctl-darwin-amd64  ./cmd/envctl
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="$(ENVCTL_LDFLAGS)" -o dist/envctl-darwin-arm64  ./cmd/envctl
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="$(ENVCTL_LDFLAGS)" -o dist/envctl-windows-amd64.exe ./cmd/envctl

## Clean
clean:
	rm -rf $(BUILD_DIR) dist

# make
# make envctl
# make dns
# make mcp-grafana
# make run-envctl
# make run-dns
# make run-mcp-grafana
# make clean