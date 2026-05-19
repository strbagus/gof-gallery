include .env

APP_NAME=main

run:
	swag init
	go run main.go

migrate-up:
	goose -dir migrations postgres $(POSTGRES_URL) up

migrate-down:
	goose -dir migrations postgres $(POSTGRES_URL) down

migrate-reset:
	goose -dir migrations postgres $(POSTGRES_URL) reset

migrate-fresh:
	goose -dir migrations postgres $(POSTGRES_URL) reset
	goose -dir migrations postgres $(POSTGRES_URL) up

migrate-status:
	goose -dir migrations postgres $(POSTGRES_URL) status

migrate-create:
	goose -dir migrations create $(name) sql
