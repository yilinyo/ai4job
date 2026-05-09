.PHONY: run test test-int migrate-up migrate-down mock lint docker-up docker-down docker-reset

DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/ai_for?sslmode=disable

run:
	go run ./cmd/server/main.go

test:
	go test ./... -v -run "^Test" -short

test-int:
	go test ./... -v -run "^TestIntegration"

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

mock:
	go run github.com/vektra/mockery/v2 --all --dir internal/domain --output internal/mocks --outpkg mocks

lint:
	golangci-lint run ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-reset:
	docker compose down --volumes
