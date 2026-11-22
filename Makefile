ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: build run test test-integration migrate-up migrate-down check-env

build:
	go build ./...

run:
	go run ./cmd/api

test:
	go test ./...

test-integration:
	go test -tags=integration -v ./...

check-env:
	@test -n "$(DATABASE_URL)" || (echo "DATABASE_URL not set: copy .env.example to .env"; exit 1)

migrate-up: check-env
	migrate -path internal/migrations -database "$(DATABASE_URL)" up

migrate-down: check-env
	migrate -path internal/migrations -database "$(DATABASE_URL)" down 1
