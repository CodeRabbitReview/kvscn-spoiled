GOLINT := golangci-lint
BIN_NAME := storage

cover: ## run all test with coverage out
	go test -v -coverprofile out/cover.out ./...
	go tool cover -html=out/cover.out -o out/cover.html
test: ## run tests
	go test --cover ./...
lint: ## lint the files local env
	$(GOLINT) run --timeout=5m -c .golangci.yml
fmt: ## fmt project
	go fmt ./...
precommit: fmt lint test
build: ## Build the binary file
	go build -o ./bin/${BIN_NAME} -a .