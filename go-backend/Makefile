COVERAGE_IGNORE_FILE=.coverageignore

.PHONY: setup
setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
	go install go.uber.org/mock/mockgen@v0.6.0
	go install github.com/google/wire/cmd/wire@v0.7.0
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

.PHONY: lint
lint:
	golangci-lint run -c .golangci.toml

.PHONY: fmt
fmt:
	golangci-lint run -c .golangci.toml --fix

.PHONY: test-unit
test-unit:
	go test ./...

.PHONY: test-integration
test-integration:
	go test ./... -tags=integration

.PHONY: test-coverage
test-coverage:
	go test ./... -tags=integration -coverprofile=coverage.out
	grep -F -v -f $(COVERAGE_IGNORE_FILE) coverage.out > coverage.tmp
	mv coverage.tmp coverage.out

.PHONY: generate-db-client
generate-db-client:
	sqlc generate -f db/sqlc.yaml

.PHONY: generate-di-container
generate-di-container:
	wire gen ./...

.PHONY: generate-api
generate-api:
	oapi-codegen -config oapi-codegen-server.yaml ../openapi/openapi.yaml
	oapi-codegen -config oapi-codegen-client.yaml ../openapi/openapi.yaml

.PHONY: generate generate-db-client generate-di-container generate-api
generate: generate-db-client generate-di-container generate-api
	go generate ./...

.PHONY: migrate-local
migrate-local:
	psqldef -U postgres -W postgres -h localhost -p 55432 backend --apply < db/schema/schema.sql
	docker exec -it db psql backend -U postgres -f /seeds/master/seed.sql
