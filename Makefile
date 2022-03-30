GOLINT := golangci-lint
BIN_NAME := storage

test: ## run all test
	go test -cover ./...
lint: ## lint the files local env
	$(GOLINT) run --timeout=5m -c .golangci.yml
fmt: ## fmt project
	go fmt ./...
precommit: fmt lint test
build: ## Build the binary file
	CGO_ENABLED=1 go build -o ./bin/${BIN_NAME} -a .