include .env

MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create migrate-up migrate-down migrate-force seed

migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name)

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" up

migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" down

migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" force $(version)

seed:
	go run cmd/seed/main.go