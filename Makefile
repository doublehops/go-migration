
.PHONY: gofmt
gofmt: ## Run gofumpt over the codebase. gofumpt must be installed and in your path.
	gofumpt -l -w .

.PHONY: lint
lint: ## Run golangci-lint. golangci-lint must be installed and in your path.
	golangci-lint run --modules-download-mode vendor

.PHONY: test
test:
	go test ./...
