.PHONY: test quality coverage coverage-html tidy deps clean fmt

help: ## Outputs this help screen
	@grep -E '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

# Default target
all: test ## Run tests (default target)

test: ## Run tests with coverage
	${MAKE} quality
	go test -v -cover ./...

coverage: ## Run tests with race detection and coverage profile
	go test -race -covermode atomic -coverprofile=covprofile ./...

coverage-html: coverage ## Generate HTML coverage report
	@mkdir -p .coverage
	go tool cover -html=covprofile -o .coverage/coverage.html

quality: ## Run linter and security checks
	${MAKE} fmt
	golangci-lint run ./... --fix
	govulncheck ./...
	${MAKE} check-alias

fmt: ## Format code with go fmt, goimports, and gofumpt
	go fmt ./...
	goimports -w .
	gofumpt -l -w .

tidy: ## Clean and tidy dependencies
	go mod tidy -e
	go get -d -v ./...

deps: ## Download dependencies
	go mod download

clean: ## Clean coverage artifacts
	rm -rf .coverage covprofile

check-alias:
	@alias_issues="$$(./.github/alias.sh || true)"; \
	if [ -n "$$alias_issues" ]; then \
		echo "❌ Unnecessary import alias detected:"; \
		echo "$$alias_issues"; \
		exit 1; \
	fi
