COVERAGE_IGNORE_FILE=.coverageignore
MOCK_DIR := test/mock

.PHONY: setup
setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
	go install go.uber.org/mock/mockgen@v0.6.0
	go install github.com/google/wire/cmd/wire@v0.7.0
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

.PHONY: clean-generated
clean-generated:
	rm -rf internal/infrastructure/sqlc
	rm -rf internal/infrastructure/http/generated
	rm -rf test/integration/client/generated
	rm -f  internal/infrastructure/di/wire_gen.go
	rm -f  test/integration/wire_gen.go
	git ls-files --others --ignored --exclude-standard $(MOCK_DIR)/ | xargs -r rm -f

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
	rm -rf internal/infrastructure/sqlc
	sqlc generate -f db/sqlc.yaml

.PHONY: generate-di-container
generate-di-container:
	rm -f internal/infrastructure/di/wire_gen.go
	rm -f test/integration/wire_gen.go
	wire gen ./...

.PHONY: generate-api
generate-api:
	rm -rf internal/infrastructure/http/generated
	rm -rf test/integration/client/generated
	mkdir -p internal/infrastructure/http/generated
	mkdir -p test/integration/client/generated
	oapi-codegen -config oapi-codegen-server.yaml ../openapi/openapi.yaml
	oapi-codegen -config oapi-codegen-client.yaml ../openapi/openapi.yaml

# Execution order: sqlc -> oapi-codegen -> mockgen -> wire
# wire_gen.go is deleted before go generate so that its embedded //go:generate
# directive does not trigger a redundant second wire run.
.PHONY: generate generate-db-client generate-api generate-di-container
generate: generate-db-client generate-api
	rm -f internal/infrastructure/di/wire_gen.go
	rm -f test/integration/wire_gen.go
	git ls-files --others --ignored --exclude-standard $(MOCK_DIR)/ | xargs -r rm -f
	go generate ./...
	$(MAKE) generate-di-container

.PHONY: migrate-local
migrate-local:
	psqldef -U postgres -W postgres -h localhost -p 55432 backend --apply < db/schema/schema.sql
	docker exec -it db psql backend -U postgres -f /seeds/master/seed.sql
