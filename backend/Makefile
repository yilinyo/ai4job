.PHONY: run test test-int migrate-up migrate-down mock lint docker-up docker-down

run:
	go run ./cmd/server/main.go

test:
	go test ./... -v -run "^Test" -short

test-int:
	go test ./... -v -run "^TestIntegration"

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ai_for?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ai_for?sslmode=disable" down 1

mock:
	go run github.com/vektra/mockery/v2 --all --dir internal/domain --output internal/mocks --outpkg mocks

lint:
	golangci-lint run ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down
