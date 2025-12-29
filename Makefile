SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)

GORELEASER_VERSION ?= v2.11.0

out:
	@mkdir -p out/build

sdk-dir:
	@rm -rf sdk
	@mkdir -p sdk

build: clean out ## Builds the provider binary
	@VERSION=$$(cat version 2>/dev/null || echo "dev"); \
	go build -ldflags "-X main.version=$$VERSION" -o out/build/pulumi-resource-hcloud-upload-image .

install-plugin: build ## Installs the plugin locally for Pulumi
	@VERSION=$$(cat version 2>/dev/null || echo "dev"); \
	pulumi plugin rm resource hcloud-upload-image --yes; \
	pulumi plugin install resource hcloud-upload-image $$VERSION --file ./out/build/pulumi-resource-hcloud-upload-image

gen-sdk: build sdk-dir ## Generates SDKs for all supported languages
	@pulumi package gen-sdk ./out/build/pulumi-resource-hcloud-upload-image --out sdk

gen-sdk-typescript: build ## Generates TypeScript SDK
	@pulumi package gen-sdk ./out/build/pulumi-resource-hcloud-upload-image --language typescript --out sdk

gen-sdk-python: build ## Generates Python SDK
	@pulumi package gen-sdk ./out/build/pulumi-resource-hcloud-upload-image --language python --out sdk

gen-sdk-go: build ## Generates Go SDK
	@pulumi package gen-sdk ./out/build/pulumi-resource-hcloud-upload-image --language go --out sdk

gen-sdk-csharp: build ## Generates C# SDK
	@pulumi package gen-sdk ./out/build/pulumi-resource-hcloud-upload-image --language csharp --out sdk

gen-sdk-java: build ## Generates Java SDK
	@pulumi package gen-sdk ./out/build/pulumi-resource-hcloud-upload-image --language java --out sdk

schema: build ## Exports the Pulumi schema
	@pulumi package get-schema ./out/build/pulumi-resource-hcloud-upload-image

download: ## Downloads the dependencies
	@go mod download

tidy: ## Cleans up go.mod and go.sum
	@go mod tidy

fmt: ## Formats all code with go fmt
	@go fmt ./...

lint: fmt $(GOLANGCI_LINT) download ## Lints all code with golangci-lint
	@go tool golangci-lint run

test: ## Runs all tests
	@go test $(ARGS) ./...

govulncheck: ## Vulnerability detection using govulncheck
	@go run golang.org/x/vuln/cmd/govulncheck ./...

version: ## Shows the current version
	@cat version 2>/dev/null || echo "dev"

set-version: ## Sets a new version (usage: make set-version VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make set-version VERSION=v1.0.0"; \
		exit 1; \
	fi
	@echo "$(VERSION)" > version
	@echo "Version set to $(VERSION)"

release: clean ## Creates a release using GoReleaser
	@VERSION=$$(cat version 2>/dev/null || echo "dev"); \
	git tag $$VERSION || true; \
	git push origin $$VERSION || true; \
	VERSION=GORELEASER_VERSION curl -sfL https://goreleaser.com/static/run | bash -s -- release --clean

release-snapshot: clean ## Creates a snapshot release using GoReleaser
	@VERSION=$$(cat version 2>/dev/null || echo "dev"); \
	VERSION=GORELEASER_VERSION GORELEASER_CURRENT_TAG=$$VERSION curl -sfL https://goreleaser.com/static/run | bash -s -- release --snapshot --clean

clean: ## Cleans up everything
	@rm -rf bin out dist

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''
