COMPOSE_FILE := docker/docker-compose.dev.yml
COMPOSE := docker compose -f $(COMPOSE_FILE)
MIGRATIONS_DIR := internal/infrastructure/migrations
APP_CONTAINER := betsystem-app

EXEC_TEST := docker exec -w /app $(APP_CONTAINER) sh -c 'set -a && . ./.env.test && set +a &&

.PHONY: up restart down delete \
	migrate-up migrate-down migrate-status \
	test test-v test-race test-integration

up:
	$(COMPOSE) up -d

restart:
	$(COMPOSE) down
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

delete:
	$(COMPOSE) down -v --remove-orphans

migrate-up:
	docker exec $(APP_CONTAINER) sh -c \
		'goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" up'

migrate-down:
	docker exec $(APP_CONTAINER) sh -c \
		'goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" down'

migrate-status:
	docker exec $(APP_CONTAINER) sh -c \
		'goose -dir $(MIGRATIONS_DIR) postgres "$$DB_URL" status'

test:
	$(EXEC_TEST) go test -count=1 ./...'

test-v:
	$(EXEC_TEST) go test -count=1 ./... -v'

test-race:
	$(EXEC_TEST) go test -count=1 ./... -race'

test-integration: up migrate-up
	$(EXEC_TEST) go test -count=1 -p 1 ./tests/integration/...'