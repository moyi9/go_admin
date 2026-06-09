APP_NAME := go_web
CONFIG := configs/config.yaml

.PHONY: dev test build tidy swagger docker-up docker-down

dev:
	go run ./cmd/server --config $(CONFIG)

test:
	go test ./...

build:
	go build -o bin/$(APP_NAME) ./cmd/server

tidy:
	go mod tidy

swagger:
	swag init -g cmd/server/main.go -o docs

docker-up:
	docker compose up -d

docker-down:
	docker compose down

