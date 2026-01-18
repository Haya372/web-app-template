.PHONY: lint
lint:
	golangci-lint run -c .golangci.toml

.PHONY: fmt
fmt:
	golangci-lint run -c .golangci.toml --fix

.PHONY: test
test-unit:
	go test ./...

.PHONY: test-integration
test-integration:
	go test ./... -tags=integration
