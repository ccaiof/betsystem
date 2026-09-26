COMPOSE_FILE := docker/docker-compose.dev.yml
COMPOSE := docker compose -f $(COMPOSE_FILE)

.PHONY: up restart down delete

up:
	$(COMPOSE) up

restart:
	$(COMPOSE) down
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

delete:
	$(COMPOSE) down -v --remove-orphans