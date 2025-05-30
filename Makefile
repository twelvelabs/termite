.DEFAULT_GOAL := help
SHELL := /bin/bash


##@ App

.PHONY: coverage
coverage: ## Show test coverage
	@make test
	@if ! command -v gocovsh >/dev/null 2>&1; then go install github.com/orlangure/gocovsh@latest; fi
	gocovsh --profile coverage.out

.PHONY: format
format: ## Format the code
	@echo "TODO..."

.PHONY: lint
lint: ## Lint the code
	@if ! command -v golangci-lint >/dev/null 2>&1; then brew install golangci-lint; fi
	golangci-lint run
	@if ! command -v actionlint >/dev/null 2>&1; then go install github.com/rhysd/actionlint/cmd/actionlint@latest; fi
	actionlint
	@if ! command -v cspell >/dev/null 2>&1; then npm install --global cspell; fi
	cspell lint --no-progress --show-suggestions .

.PHONY: test
test: export APP_ENV := test
test: ## Test the code
	go mod tidy
	go test --coverprofile=coverage.out ./...

.PHONY: generate
generate:
	@if ! command -v moq >/dev/null 2>&1; then go install github.com/matryer/moq@latest; fi
	go generate ./...

.PHONY: build
build: ## Build the code
	go mod tidy
	go build

.PHONY: release
release: ## Create a new GitHub release
	git fetch --all --tags
	@if ! command -v svu >/dev/null 2>&1; then echo "Unable to find svu!"; exit 1; fi
	@if [[ "$$(svu next)" == "$$(svu current)" ]]; then echo "Nothing to release!" && exit 1; fi
	git tag -a "$$(svu next)" -m "Release version $$(svu next)" && git push origin --tags

.PHONY: clean
clean: ## Clean build artifacts
	rm ./coverage.out


##@ Other

.PHONY: setup
setup: ## Bootstrap for local development
	@if ! command -v gh >/dev/null 2>&1; then echo "Unable to find gh!"; exit 1; fi
	@if ! command -v git >/dev/null 2>&1; then echo "Unable to find git!"; exit 1; fi
	@if ! command -v go >/dev/null 2>&1; then echo "Unable to find go!"; exit 1; fi

# Via https://www.thapaliya.com/en/writings/well-documented-makefiles/
# Note: The `##@` comments determine grouping
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
