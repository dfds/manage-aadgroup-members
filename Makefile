.PHONY: run build start stop

build:
	go mod tidy
	swag init
	echo "" >> docs/swagger.json # In order to add trailing line
	go build .

run:
	swag init
	echo "" >> docs/swagger.json # In order to add trailing line
	go run main.go

start:
	swag init
	docker compose build
	docker compose up -d
stop:
	docker compose down
