.PHONY: setup
setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
	go install go.uber.org/mock/mockgen@v0.6.0
	go install github.com/google/wire/cmd/wire@v0.7.0

.PHONY: lint
lint:
	golangci-lint run -c .golangci.toml

.PHONY: fmt
fmt:
	golangci-lint run -c .golangci.toml --fix

.PHONY: test
test-unit:
	go test ./...

# TODO: remove -p option
.PHONY: test-integration
test-integration:
	go test ./... -tags=integration -p 1

.PHONY: generate-db-client
generate-db-client:
	sqlc generate -f db/sqlc.yaml

.PHONY: generate-di-container
generate-di-container:
	wire gen ./...

.PHONY: generate generate-db-client generate-di-container
generate: generate-db-client generate-di-container
	go generate ./...

