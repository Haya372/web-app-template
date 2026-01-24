.PHONY: setup
setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
	go install go.uber.org/mock/mockgen@v0.6.0

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

.PHONY: generate-db-client
generate-db-client:
	sqlc generate -f db/sqlc.yaml

.PHONY: generate generate-db-client
generate: generate-db-client
	go generate ./...
