-include .env
export

migrate_up: postgres-up
	docker compose run --rm migrate
migrate_down: postgres-up
	docker compose run --rm migrate down

postgres-up:
	docker compose up -d postgres
	@attempt=0; \
	until docker compose exec -T postgres sh -c 'pg_isready -h 127.0.0.1 -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"' >/dev/null 2>&1; do \
		attempt=$$((attempt + 1)); \
		if [ "$$attempt" -ge 30 ]; then \
			echo "PostgreSQL did not become ready. Check: docker compose logs postgres"; \
			exit 1; \
		fi; \
		sleep 1; \
	done

run-http-app: migrate_up
	docker compose up -d --build application

service_run:
	@go run main.go
service_deploy: run-http-app
service_undeploy:
	docker compose down

.PHONY: migrate_up migrate_down postgres-up run-http-app service_run service_deploy service_undeploy
