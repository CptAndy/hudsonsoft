include .env

MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) $(word 2,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up: 
	@migrate -path=$(MIGRATIONS_PATH) -database $(DB_MIGRATOR_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database $(DB_MIGRATOR_ADDR) down $(filter_out $@,$(MAKECMDGOALS))