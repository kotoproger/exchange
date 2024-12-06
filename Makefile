ifneq ($(wildcard .env),)
	include .env
endif
PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
PROJECT_TMP = $(PROJECT_DIR)/tmp
install-deps:
	GOBIN=$(PROJECT_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(PROJECT_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	GOBIN=$(PROJECT_BIN) go install github.com/golang/mock/mockgen@latest
migration-status:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} status -v
migration-add:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} create $(name) sql
migration-up:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v
migration-down:
	$(PROJECT_BIN)/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v
gen-sql:
	rm -f internal/repository/*.go
	$(PROJECT_BIN)/sqlc generate
lint:
	~/go/bin/golangci-lint run --fix

docker-up:
	docker compose up -d

docker-migrate:
	docker compose exec -t app make migration-up

docker-app:
	docker compose exec -t app /application