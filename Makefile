include .env

APP_NAME=main

run:
	swag init
	go run main.go

build:
	swag init \
	export GOARCH=$(arch) \
	go build main.go

migrate-up:
	goose -dir db/migrations postgres $(POSTGRES_URL) up

migrate-down:
	goose -dir db/migrations postgres $(POSTGRES_URL) down

migrate-reset:
	goose -dir db/migrations postgres $(POSTGRES_URL) reset

migrate-fresh:
	goose -dir db/migrations postgres $(POSTGRES_URL) reset
	goose -dir db/migrations postgres $(POSTGRES_URL) up

migrate-status:
	goose -dir db/migrations postgres $(POSTGRES_URL) status

migrate-create:
	goose -dir db/migrations create $(name) sql

seeder:
	goose -dir db/seeders -no-versioning postgres $(POSTGRES_URL) up

seeder-nuke:
	goose -dir db/seeders -no-versioning postgres $(POSTGRES_URL) down-to 0

