SHELL=/bin/bash -e -o pipefail
PWD = $(shell pwd)

out:
	@mkdir -p out/build

build: out ## Builds the provider binary
	@go build -o bin/pulumi-resource-hcloud-upload-image .

install: build ## Installs the provider binary to /usr/local/bin
	@cp bin/pulumi-resource-hcloud-upload-image /usr/local/bin/

gen-sdk: build ## Generates SDKs for all supported languages
	@pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image

gen-sdk-typescript: build ## Generates TypeScript SDK
	@pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image --language typescript

gen-sdk-python: build ## Generates Python SDK
	@pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image --language python

gen-sdk-go: build ## Generates Go SDK
	@pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image --language go

gen-sdk-csharp: build ## Generates C# SDK
	@pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image --language csharp

gen-sdk-java: build ## Generates Java SDK
	@pulumi package gen-sdk ./bin/pulumi-resource-hcloud-upload-image --language java

schema: build ## Exports the Pulumi schema
	@pulumi package get-schema ./bin/pulumi-resource-hcloud-upload-image

download: ## Downloads the dependencies
	@go mod download

tidy: ## Cleans up go.mod and go.sum
	@go mod tidy

fmt: ## Formats all code with go fmt
	@go fmt ./...

lint: fmt $(GOLANGCI_LINT) download ## Lints all code with golangci-lint
	@go tool -modfile=golangci-lint.mod golangci-lint run

test: ## Runs all tests
	@go test $(ARGS) ./...

govulncheck: ## Vulnerability detection using govulncheck
	@go run golang.org/x/vuln/cmd/govulncheck ./...

clean: ## Cleans up everything
	@rm -rf bin out sdk

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''
