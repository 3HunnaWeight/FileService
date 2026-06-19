.PHONY: build run migrate-up migrate-down migrate-status

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/file-service ./cmd/api

run: build
	./bin/file-service

migrate-up:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" up

migrate-down:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" down

migrate-status:
	goose -dir ./migrations postgres "$(POSTGRES_DSN)" status

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build
