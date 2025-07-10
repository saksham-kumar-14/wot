include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: test
test:
	@go test -v ./...

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: generate-docs
generate-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt


DB_NAME=wot
DB_USER=admin
DB_HOST=localhost
DB_PORT=5432
DB_PASSWORD=strongpwd
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

db-reset:
	dropdb --if-exists -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME) || true
	createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)
