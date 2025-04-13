###############################################################################
###                                Tool			                            		###
###############################################################################

mockgen:
	@go install go.uber.org/mock/mockgen@v0.5.1
	@./scripts/mockgen.sh
.PHONY: mocks

###############################################################################
###                                Lint			                            		###
###############################################################################

lint:
	@echo "Running golangci-lint..."
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.1
	@golangci-lint run --timeout 5m
	@echo "golangci-lint completed successfully."
.PHONY: lint

format:
	@echo "Running format..."
	@find . -name '*.go' -type f -not -path "*.git*" -not -name '*_mock.go' | xargs gofmt -w -s
	@find . -name '*.go' -type f -not -path "*.git*" -not -name '*_mock.go' | xargs goimports -w
	@echo "format completed successfully."
.PHONY: format

################################################################################
###                                Test			                            		 ###
################################################################################

test:
	@echo "Running tests..."
	@go test -v `go list ./... | grep -v "/testutil"` -coverprofile=coverage.out
	@echo "Tests completed successfully."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
.PHONY: test

bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./... > benchmarks.out
	@echo "Benchmarks completed successfully."
.PHONY: bench