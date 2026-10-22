-include .env
export

migrate_up:
	@migrate -path ./migrations -database "postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" up
migrate_down:
	@migrate -path ./migrations -database "postgres://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" down
service_run:
	@go run main.go

.PHONY: migrate_up migrate_down service_run run-http-app stop-http-app
