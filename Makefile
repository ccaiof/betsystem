COMPOSE_FILE := docker/docker-compose.dev.yml
COMPOSE := docker compose -f $(COMPOSE_FILE)
MIGRATIONS_DIR := internal/infrastructure/migrations

.PHONY: up restart down delete migrate-up migrate-down migrate-status test test-v test-race

up:
	$(COMPOSE) up

restart:
	$(COMPOSE) down
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

delete:
	$(COMPOSE) down -v --remove-orphans

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

test:
	bash -c 'set -a && source ./.env.test && set +a && go test ./...'

test-v:
	bash -c 'set -a && source ./.env.test && set +a && go test ./... -v'

test-race:
	bash -c 'set -a && source ./.env.test && set +a && go test ./... -race'
