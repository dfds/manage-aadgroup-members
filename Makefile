.PHONY: run build start stop

build:
	go mod tidy
	swag init
	go build .

run:
	swag init
	go run main.go

start:
	swag init
	docker compose build
	docker compose up -d

stop:
	docker compose down
